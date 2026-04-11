import { useState, useEffect, useContext, useRef, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import Editor from '@monaco-editor/react';
import { getExamDetails } from '../api/exams';
import client from '../api/client';
import { AuthContext } from '../context/AuthContext';
import Swal from 'sweetalert2';
import {
    CheckCircle2, XCircle, Clock, Target, Code,
    ChevronRight, ChevronLeft, Send, LogOut, AlertCircle, Timer
} from 'lucide-react';
import './ChallengeSolver.css';

const ExamRunner = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);

    const [exam, setExam] = useState(null);
    const [challenges, setChallenges] = useState([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [codeMap, setCodeMap] = useState({});         // { challengeId: code }
    const [resultMap, setResultMap] = useState({});      // { challengeId: result }
    const [language, setLanguage] = useState('python');
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [output, setOutput] = useState('');

    // Session state
    const [sessionId, setSessionId] = useState(null);
    const [sessionStatus, setSessionStatus] = useState(null);
    const [timeLeft, setTimeLeft] = useState(null); // in seconds, null = unlimited
    const [examFinished, setExamFinished] = useState(false);
    const [attemptMap, setAttemptMap] = useState({}); // { challengeId: count }
    const heartbeatRef = useRef(null);
    const timerRef = useRef(null);

    // Professors should use the editor, not the runner
    useEffect(() => {
        const role = user?.role;
        if (role === 'professor' || role === 'teacher' || role === 'admin') {
            navigate(`/exam/${id}/edit`, { replace: true });
        }
    }, [user, id, navigate]);

    const sleep = (ms) => new Promise(r => setTimeout(r, ms));

    // Format seconds as MM:SS or HH:MM:SS
    const formatTime = (seconds) => {
        if (seconds == null || seconds < 0) return '∞';
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = seconds % 60;
        if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
        return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
    };
    /**
     * Ensure there's an active session for this exam.
     * Note: The backend's close-session endpoint is professor-only,
     * so students can't close their own sessions. Instead we let old sessions
     * freeze (no heartbeat → backend marks them frozen after 60s inactivity),
     * and CreateSession auto-expires frozen sessions before creating new ones.
     *
     * Strategy:
     *   1. Try POST to create a new session.
     *   2. If "active session" conflict → discover it via heartbeat.
     *   3. If the session is for THIS exam on a fresh mount → reuse it (page refresh).
     *   4. If the session is old/stale → wait for it to freeze and retry creation.
     */
    const ensureSession = useCallback(async (examId) => {
        // Helper to apply session data to state
        const applySession = (sessionData, sid) => {
            localStorage.setItem('session_id', sid);
            setSessionId(sid);
            setSessionStatus(sessionData?.status || sessionData?.Status || 'active');
            const tl = sessionData?.time_left ?? sessionData?.TimeLeft ?? sessionData?.timeLeft;
            if (tl != null && tl > 0) {
                setTimeLeft(tl);
            } else if (tl === -1) {
                setTimeLeft(null); // unlimited
            }
        };

        const tryCreate = async () => {
            const userId = user?.id || user?.ID || '';
            const createRes = await client.post('/submissions/sessions', {
                user_id: userId,
                exam_id: examId,
            });
            return createRes?.data;
        };

        // Step 1: Try to create a new session
        try {
            const newSession = await tryCreate();
            const sid = newSession?.id || newSession?.ID;
            if (sid) {
                applySession(newSession, sid);
                return newSession;
            }
        } catch (err) {
            const apiMsg = err?.response?.data?.error || err?.message || '';

            // Step 2: Backend says "already has active session"
            if (apiMsg.toLowerCase().includes('active session')) {
                // Discover the real active session via heartbeat probe
                const probeId = localStorage.getItem('session_id') || 'probe';
                let realSid = null;
                let sessionExamId = '';

                try {
                    const hbRes = await client.post(`/submissions/sessions/${probeId}/heartbeat`);
                    const activeSession = hbRes?.data;
                    realSid = activeSession?.id || activeSession?.ID;
                    sessionExamId = activeSession?.exam_id || activeSession?.ExamID || activeSession?.examID || '';

                    // Step 3: If session is for THIS exam, reuse it (student refreshed the page)
                    if (realSid && sessionExamId === examId) {
                        applySession(activeSession, realSid);
                        return { id: realSid, status: 'active', ...activeSession };
                    }

                    // Step 4: Session is for a DIFFERENT exam → close it and create new
                    if (realSid) {
                        try {
                            await client.post(`/submissions/sessions/${realSid}/close`);
                            localStorage.removeItem('session_id');
                        } catch (closeErr) {
                            console.warn('Failed to close old session:', closeErr);
                        }

                        // Retry creating a new session
                        try {
                            const newSession = await tryCreate();
                            const sid = newSession?.id || newSession?.ID;
                            if (sid) {
                                applySession(newSession, sid);
                                return newSession;
                            }
                        } catch (retryErr) {
                            console.error('Retry after close failed:', retryErr);
                        }
                    }
                } catch (hbErr) {
                    // Heartbeat failed — session may already be frozen/expired
                    console.warn('Heartbeat probe failed (session may be frozen):', hbErr);
                    localStorage.removeItem('session_id');

                    // Try creating immediately — frozen session will be auto-expired by backend
                    try {
                        const newSession = await tryCreate();
                        const sid = newSession?.id || newSession?.ID;
                        if (sid) {
                            applySession(newSession, sid);
                            return newSession;
                        }
                    } catch (retryErr) {
                        console.error('Retry after frozen failed:', retryErr);
                    }
                }

                // All recovery attempts failed
                Swal.fire({
                    icon: 'info',
                    title: 'Sesión activa existente',
                    html: 'No se pudo recuperar ni cerrar la sesión anterior.<br/>Contacta a tu profesor para asistencia.',
                });
                return null;
            }

            console.error('Failed to create session:', err);
            Swal.fire({
                icon: 'error',
                title: 'Error al iniciar sesión de examen',
                text: apiMsg || 'No se pudo crear la sesión. Inténtalo de nuevo.',
            });
        }

        return null;
    }, [user]);

    // Fetch exam data and create/retrieve session
    useEffect(() => {
        const fetchExam = async () => {
            try {
                const data = await getExamDetails(id);
                setExam(data);

                const itemsRes = await client.get(`/exams/${id}/items`);
                const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);
                const challengeList = items
                    .map(item => {
                        const ch = item.challenge || item.Challenge || {};
                        return {
                            ...ch,
                            id: ch.id || ch.ID || item.challenge_id || item.challengeID,
                            title: ch.title || ch.Title || 'Reto',
                            description: ch.description || ch.Description || '',
                            difficulty: ch.difficulty || ch.Difficulty || 'medium',
                            constraints: ch.constraints || ch.Constraints || '',
                            points: item.points || item.Points || 0,
                            order: item.order || item.Order || 0,
                        };
                    })
                    .filter(ch => ch.id)
                    .sort((a, b) => a.order - b.order);

                setChallenges(challengeList);

                // Initialize code templates
                const initialCode = {};
                challengeList.forEach(ch => {
                    initialCode[ch.id] = '# Escribe tu solución aquí\ndef solve():\n    pass\n';
                });
                setCodeMap(initialCode);

                // --- Create or retrieve the exam session ---
                await ensureSession(id);
            } catch (err) {
                console.error(err);
                Swal.fire({ icon: 'error', title: 'Error', text: err?.response?.data?.error || 'No se pudo cargar el examen.' });
                navigate('/public-exams');
            } finally {
                setLoading(false);
            }
        };

        fetchExam();
    }, [id, navigate, ensureSession]);

    // Heartbeat: send every 60 seconds to keep session alive and avoid frozen state
    useEffect(() => {
        if (!sessionId) return;

        const sendHeartbeat = async () => {
            try {
                await client.post(`/submissions/sessions/${sessionId}/heartbeat`);
            } catch (err) {
                console.warn('Heartbeat failed:', err);
                const status = err?.response?.status;
                const msg = err?.response?.data?.error || '';
                // If session is no longer active, notify user
                if (status === 400 || status === 403 || msg.toLowerCase().includes('blocked') || msg.toLowerCase().includes('expired')) {
                    Swal.fire({
                        icon: 'warning',
                        title: 'Sesión finalizada',
                        text: 'Tu sesión de examen ha terminado.',
                        confirmButtonText: 'Aceptar'
                    }).then(() => navigate('/public-exams'));
                }
            }
        };

        // Send first heartbeat immediately
        sendHeartbeat();
        heartbeatRef.current = setInterval(sendHeartbeat, 30000);

        return () => {
            if (heartbeatRef.current) clearInterval(heartbeatRef.current);
        };
    }, [sessionId, navigate]);

    // Close session when browser is closed or refreshed
    useEffect(() => {
        if (!sessionId) return;

        const handleBeforeUnload = () => {
            const sid = sessionId || localStorage.getItem('session_id');
            if (!sid) return;

            const baseURL = client.defaults.baseURL || '';
            const token = localStorage.getItem('token');
            const email = localStorage.getItem('user_email');
            const url = `${baseURL}/submissions/sessions/${sid}/close`;

            try {
                fetch(url, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'X-User-Email': email || '',
                        'Content-Type': 'application/json',
                    },
                    keepalive: true,
                    body: '{}',
                });
            } catch (e) {
                // Best-effort
            }
            localStorage.removeItem('session_id');
        };

        window.addEventListener('beforeunload', handleBeforeUnload);

        return () => {
            window.removeEventListener('beforeunload', handleBeforeUnload);
        };
    }, [sessionId]);

    // Countdown timer
    useEffect(() => {
        if (timeLeft == null || timeLeft <= 0) return;

        timerRef.current = setInterval(() => {
            setTimeLeft(prev => {
                if (prev == null) return null;
                if (prev <= 1) {
                    clearInterval(timerRef.current);
                    if (heartbeatRef.current) clearInterval(heartbeatRef.current);

                    // Close session on the backend
                    const sid = sessionId || localStorage.getItem('session_id');
                    if (sid) {
                        client.post(`/submissions/sessions/${sid}/close`).catch(() => { });
                    }
                    localStorage.removeItem('session_id');
                    setSessionId(null);

                    setExamFinished(true);
                    Swal.fire({
                        icon: 'warning',
                        title: '⏰ Tiempo agotado',
                        html: 'El tiempo del examen ha finalizado.',
                        confirmButtonText: 'Ver resultados'
                    });
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => {
            if (timerRef.current) clearInterval(timerRef.current);
        };
    }, [timeLeft != null && timeLeft > 0]); // re-run only when timer starts

    const currentChallenge = challenges[currentIndex] || null;
    const currentCode = currentChallenge ? (codeMap[currentChallenge.id] || '') : '';
    const currentResult = currentChallenge ? (resultMap[currentChallenge.id] || null) : null;

    const handleCodeChange = (value) => {
        if (!currentChallenge) return;
        setCodeMap(prev => ({ ...prev, [currentChallenge.id]: value || '' }));
    };

    const handleSelectChallenge = (idx) => {
        setCurrentIndex(idx);
        setOutput('');
    };

    // Submit solution
    const handleSubmit = async () => {
        if (examFinished) {
            setOutput('El examen ha finalizado. No se pueden realizar más envíos.');
            return;
        }

        // Check try limit per challenge
        const tryLimit = exam?.tryLimit || exam?.TryLimit || exam?.try_limit || -1;
        const challengeId = currentChallenge?.id;
        const currentAttempts = attemptMap[challengeId] || 0;
        if (tryLimit > 0 && currentAttempts >= tryLimit) {
            setOutput(`Has alcanzado el límite de ${tryLimit} intento(s) para este reto.`);
            return;
        }

        if (!currentCode.trim()) {
            setOutput('No puedes enviar código vacío.');
            return;
        }

        if (!sessionId) {
            // Try to create/retrieve session one more time
            const session = await ensureSession(id);
            if (!session) {
                Swal.fire({ icon: 'warning', title: 'Sesión no activa', text: 'No se pudo obtener una sesión activa. Reintenta o vuelve a entrar al examen.', customClass: { container: 'swal-ultra-high-z' } });
                return;
            }
        }

        const activeSessionId = sessionId || localStorage.getItem('session_id');
        if (!activeSessionId) {
            Swal.fire({ icon: 'warning', title: 'Sesión no activa', text: 'No hay una sesión de examen activa. Vuelve a entrar al examen.', customClass: { container: 'swal-ultra-high-z' } });
            return;
        }


        setSubmitting(true);
        setOutput('Enviando solución...');

        try {
            // Send heartbeat right before submission to keep session active
            try {
                await client.post(`/submissions/sessions/${activeSessionId}/heartbeat`);
            } catch (hbErr) {
                console.warn('Pre-submit heartbeat failed:', hbErr);
            }

            const { data } = await client.post('/submissions', {
                code: currentCode,
                language: language,
                challenge_id: challengeId,
                session_id: activeSessionId
            });

            const submissionId = data?.id || data?.ID;
            if (!submissionId) {
                setOutput('La API no retornó un ID de submission.');
                setSubmitting(false);
                return;
            }

            setOutput('Solución enviada. Ejecutando pruebas...');

            // Poll for results
            for (let attempt = 0; attempt < 40; attempt++) {
                const res = await client.get(`/submissions/${submissionId}`);
                const submission = res?.data?.Submission || res?.data?.submission;
                const results = res?.data?.Results || res?.data?.results || [];

                if (Array.isArray(results) && results.length > 0) {
                    const hasPending = results.some(r => {
                        const s = String(r?.Status || r?.status || '').toLowerCase();
                        return s === 'queued' || s === 'running';
                    });

                    if (!hasPending) {
                        const accepted = results.filter(r => String(r?.Status || r?.status || '').toLowerCase() === 'accepted').length;
                        const score = Math.round((accepted / results.length) * 100);
                        const status = score === 100 ? 'accepted' : 'wrong_answer';

                        setResultMap(prev => ({ ...prev, [challengeId]: { status, score, results } }));

                        // Track attempt count for this challenge
                        setAttemptMap(prev => ({ ...prev, [challengeId]: (prev[challengeId] || 0) + 1 }));

                        // Build output
                        const attemptsUsed = (attemptMap[challengeId] || 0) + 1;
                        const tryLimitVal = exam?.tryLimit || exam?.TryLimit || exam?.try_limit || -1;
                        const attemptsInfo = tryLimitVal > 0 ? ` | Intentos: ${attemptsUsed}/${tryLimitVal}` : '';
                        const lines = [`Resultado: ${score}% (${accepted}/${results.length} casos correctos)${attemptsInfo}\n`];
                        results.forEach((r, i) => {
                            const st = (r.Status || r.status || 'unknown').toLowerCase();
                            const err = r.ErrorMessage || r.errorMessage || '';
                            lines.push(`  Caso ${i + 1}: ${st === 'accepted' ? '✅' : '❌'} ${st}${err ? ' - ' + err : ''}`);
                        });
                        setOutput(lines.join('\n'));

                        if (score === 100) {
                            Swal.fire({
                                icon: 'success',
                                title: '¡Correcto!',
                                text: `Todos los casos de prueba pasaron. (+${currentChallenge.points} pts)`,
                                confirmButtonText: 'Genial',
                                customClass: { container: 'swal-ultra-high-z' },
                                backdrop: `rgba(0,0,0,0.4)`
                            });
                        }

                        setSubmitting(false);
                        return;
                    }
                }
                await sleep(1500);
            }

            setOutput('La evaluación sigue en proceso. Revisa más tarde.');
        } catch (err) {
            const msg = err?.response?.data?.error || err?.message || 'Error al enviar';
            setOutput(`Error: ${msg}`);
            Swal.fire({
                icon: 'error',
                title: 'Error de Envío',
                text: msg,
                customClass: { container: 'swal-ultra-high-z' }
            });
        } finally {
            setSubmitting(false);
        }
    };

    // Finish exam — stop heartbeat so session freezes on backend
    const handleFinishExam = async () => {
        const solved = Object.values(resultMap).filter(r => r.status === 'accepted').length;
        const { isConfirmed } = await Swal.fire({
            title: '¿Terminar Examen?',
            html: `Has resuelto correctamente <strong>${solved} de ${challenges.length}</strong> retos.<br/>Esta acción es definitiva.`,
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Sí, terminar',
            cancelButtonText: 'Seguir resolviendo',
            confirmButtonColor: '#c8102e'
        });

        if (isConfirmed) {
            // Close session on the backend
            const activeSessionId = sessionId || localStorage.getItem('session_id');
            if (activeSessionId) {
                try {
                    await client.post(`/submissions/sessions/${activeSessionId}/close`);
                } catch (err) {
                    console.warn('Failed to close session:', err);
                }
            }
            if (heartbeatRef.current) clearInterval(heartbeatRef.current);
            if (timerRef.current) clearInterval(timerRef.current);
            localStorage.removeItem('session_id');
            setSessionId(null);

            setExamFinished(true);
            Swal.fire({ icon: 'success', title: 'Examen Finalizado', text: `Puntuación: ${solved}/${challenges.length} retos correctos.` });
        }
    };

    // --- RENDER ---
    if (loading) return (
        <div className="dashboard-loading">
            <div className="loader-orbit"><div className="orbit-dot"></div></div>
            <p>Cargando examen...</p>
        </div>
    );

    if (!exam) return (
        <div className="dashboard-loading error">
            <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>📝</div>
            <h2>Examen no encontrado</h2>
            <button onClick={() => navigate('/public-exams')} className="btn-retry" style={{ marginTop: '2rem' }}>Volver</button>
        </div>
    );

    return (
        <div className="solver-container" style={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
            {/* TOP BAR */}
            <div className="solver-header" style={{
                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                padding: '0.75rem 1.5rem', background: '#1e1e2e', color: 'white', flexShrink: 0
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <Target size={20} style={{ color: '#c8102e' }} />
                    <h2 style={{ margin: 0, fontSize: '1.1rem', fontWeight: 700 }}>{exam.title || exam.Title}</h2>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <span style={{ fontSize: '0.85rem', opacity: 0.7 }}>
                        {Object.values(resultMap).filter(r => r.status === 'accepted').length}/{challenges.length} resueltos
                    </span>
                    {/* Countdown Timer */}
                    {timeLeft != null && (
                        <div style={{
                            display: 'flex', alignItems: 'center', gap: '6px',
                            padding: '0.4rem 0.85rem', borderRadius: '10px',
                            background: timeLeft > 300 ? 'rgba(16,185,129,0.15)' : timeLeft > 60 ? 'rgba(245,158,11,0.2)' : 'rgba(239,68,68,0.25)',
                            border: `1px solid ${timeLeft > 300 ? '#10b981' : timeLeft > 60 ? '#f59e0b' : '#ef4444'}`,
                            fontVariantNumeric: 'tabular-nums',
                            animation: timeLeft <= 60 ? 'pulse 1s ease-in-out infinite' : 'none'
                        }}>
                            <Timer size={15} style={{ color: timeLeft > 300 ? '#10b981' : timeLeft > 60 ? '#f59e0b' : '#ef4444' }} />
                            <span style={{
                                fontSize: '0.9rem', fontWeight: 800, letterSpacing: '0.5px',
                                color: timeLeft > 300 ? '#10b981' : timeLeft > 60 ? '#f59e0b' : '#ef4444'
                            }}>
                                {formatTime(timeLeft)}
                            </span>
                        </div>
                    )}
                    {timeLeft == null && (
                        <div style={{
                            display: 'flex', alignItems: 'center', gap: '6px',
                            padding: '0.4rem 0.85rem', borderRadius: '10px',
                            background: 'rgba(139,92,246,0.12)', border: '1px solid rgba(139,92,246,0.3)',
                        }}>
                            <Clock size={15} style={{ color: '#a78bfa' }} />
                            <span style={{ fontSize: '0.85rem', fontWeight: 600, color: '#a78bfa' }}>Sin límite</span>
                        </div>
                    )}
                    {!examFinished ? (
                        <button onClick={handleFinishExam} style={{
                            background: '#c8102e', color: 'white', border: 'none', borderRadius: '10px',
                            padding: '0.5rem 1rem', fontWeight: 700, fontSize: '0.85rem', cursor: 'pointer',
                            display: 'flex', alignItems: 'center', gap: '6px'
                        }}>
                            <LogOut size={16} /> Terminar Examen
                        </button>
                    ) : (
                        <button onClick={() => navigate('/public-exams')} style={{
                            background: '#374151', color: 'white', border: 'none', borderRadius: '10px',
                            padding: '0.5rem 1rem', fontWeight: 700, fontSize: '0.85rem', cursor: 'pointer',
                            display: 'flex', alignItems: 'center', gap: '6px'
                        }}>
                            <LogOut size={16} /> Volver a Exámenes
                        </button>
                    )}
                </div>
            </div>

            {/* MAIN AREA */}
            <div style={{ display: 'flex', flex: 1, overflow: 'hidden' }}>
                {/* LEFT SIDEBAR: Challenge List */}
                <div style={{
                    width: '260px', background: '#252536', color: 'white',
                    display: 'flex', flexDirection: 'column', overflowY: 'auto', flexShrink: 0
                }}>
                    <div style={{ padding: '1rem', borderBottom: '1px solid rgba(255,255,255,0.08)' }}>
                        <h3 style={{ margin: 0, fontSize: '0.9rem', fontWeight: 700, opacity: 0.7 }}>RETOS DEL EXAMEN</h3>
                    </div>
                    {challenges.map((ch, idx) => {
                        const chResult = resultMap[ch.id];
                        const isActive = idx === currentIndex;
                        const isSolved = chResult?.status === 'accepted';
                        const isFailed = chResult && chResult.status !== 'accepted';
                        return (
                            <button
                                key={ch.id}
                                onClick={() => handleSelectChallenge(idx)}
                                style={{
                                    display: 'flex', alignItems: 'center', gap: '0.75rem',
                                    padding: '0.85rem 1rem', border: 'none', textAlign: 'left',
                                    background: isActive ? 'rgba(200,16,46,0.2)' : 'transparent',
                                    color: 'white', cursor: 'pointer', width: '100%',
                                    borderLeft: isActive ? '3px solid #c8102e' : '3px solid transparent',
                                    transition: 'all 0.2s'
                                }}
                            >
                                <div style={{
                                    width: '28px', height: '28px', borderRadius: '50%', flexShrink: 0,
                                    display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '0.75rem', fontWeight: 700,
                                    background: isSolved ? '#10b981' : isFailed ? '#ef4444' : isActive ? '#c8102e' : 'rgba(255,255,255,0.1)',
                                }}>
                                    {isSolved ? <CheckCircle2 size={14} /> : isFailed ? <XCircle size={14} /> : idx + 1}
                                </div>
                                <div style={{ overflow: 'hidden' }}>
                                    <div style={{ fontSize: '0.85rem', fontWeight: 600, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                        {ch.title}
                                    </div>
                                    <div style={{ fontSize: '0.7rem', opacity: 0.5 }}>{ch.points} pts • {ch.difficulty === 'easy' ? 'Fácil' : ch.difficulty === 'hard' ? 'Difícil' : 'Medio'}</div>
                                </div>
                            </button>
                        );
                    })}
                </div>

                {/* CENTER: Problem Description */}
                {currentChallenge && (
                    <>
                        <div className="problem-description" style={{
                            width: '35%', overflowY: 'auto', padding: '1.5rem', background: '#fafafa',
                            borderRight: '1px solid #e5e7eb'
                        }}>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '1rem' }}>
                                <Code size={18} style={{ color: '#c8102e' }} />
                                <h2 style={{ margin: 0, fontSize: '1.3rem', fontWeight: 800 }}>{currentChallenge.title}</h2>
                            </div>
                            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1.5rem' }}>
                                <span style={{
                                    padding: '3px 10px', borderRadius: '6px', fontSize: '0.7rem', fontWeight: 700,
                                    background: currentChallenge.difficulty === 'easy' ? '#dcfce7' : currentChallenge.difficulty === 'hard' ? '#fee2e2' : '#fef3c7',
                                    color: currentChallenge.difficulty === 'easy' ? '#15803d' : currentChallenge.difficulty === 'hard' ? '#b91c1c' : '#b45309'
                                }}>
                                    {currentChallenge.difficulty === 'easy' ? 'Fácil' : currentChallenge.difficulty === 'hard' ? 'Difícil' : 'Medio'}
                                </span>
                                <span style={{ padding: '3px 10px', borderRadius: '6px', fontSize: '0.7rem', fontWeight: 700, background: '#e0e7ff', color: '#3730a3' }}>
                                    {currentChallenge.points} pts
                                </span>
                            </div>
                            <p style={{ lineHeight: 1.7, color: '#444', fontSize: '0.95rem', whiteSpace: 'pre-wrap' }}>
                                {currentChallenge.description}
                            </p>
                            {currentChallenge.constraints && (
                                <div style={{ marginTop: '1.5rem', padding: '1rem', background: '#fff7ed', borderRadius: '10px', border: '1px solid #fed7aa' }}>
                                    <h4 style={{ margin: '0 0 0.5rem', fontSize: '0.85rem', color: '#9a3412' }}>⚡ Restricciones</h4>
                                    <p style={{ margin: 0, fontSize: '0.85rem', color: '#78350f' }}>{currentChallenge.constraints}</p>
                                </div>
                            )}

                            {/* Nav buttons */}
                            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '2rem', gap: '0.5rem' }}>
                                <button
                                    disabled={currentIndex === 0}
                                    onClick={() => handleSelectChallenge(currentIndex - 1)}
                                    style={{
                                        flex: 1, padding: '0.6rem', border: '1px solid #ddd', borderRadius: '10px',
                                        background: 'white', cursor: currentIndex === 0 ? 'not-allowed' : 'pointer',
                                        opacity: currentIndex === 0 ? 0.4 : 1, fontWeight: 600, fontSize: '0.85rem',
                                        display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '4px'
                                    }}
                                >
                                    <ChevronLeft size={16} /> Anterior
                                </button>
                                <button
                                    disabled={currentIndex === challenges.length - 1}
                                    onClick={() => handleSelectChallenge(currentIndex + 1)}
                                    style={{
                                        flex: 1, padding: '0.6rem', border: '1px solid #ddd', borderRadius: '10px',
                                        background: 'white', cursor: currentIndex === challenges.length - 1 ? 'not-allowed' : 'pointer',
                                        opacity: currentIndex === challenges.length - 1 ? 0.4 : 1, fontWeight: 600, fontSize: '0.85rem',
                                        display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '4px'
                                    }}
                                >
                                    Siguiente <ChevronRight size={16} />
                                </button>
                            </div>
                        </div>

                        {/* RIGHT: Code Editor + Output */}
                        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', background: '#1e1e1e' }}>
                            {/* Editor toolbar */}
                            <div style={{
                                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                                padding: '0.5rem 1rem', background: '#2d2d3d', borderBottom: '1px solid #3d3d50'
                            }}>
                                <select value={language} onChange={(e) => setLanguage(e.target.value)}
                                    style={{ background: '#1e1e2e', color: 'white', border: '1px solid #555', borderRadius: '6px', padding: '0.35rem 0.75rem', fontSize: '0.85rem' }}>
                                    <option value="python">Python</option>
                                </select>
                                {(() => {
                                    const tryLimitVal = exam?.tryLimit || exam?.TryLimit || exam?.try_limit || -1;
                                    const chAttempts = currentChallenge ? (attemptMap[currentChallenge.id] || 0) : 0;
                                    const limitReached = tryLimitVal > 0 && chAttempts >= tryLimitVal;
                                    const isDisabled = submitting || examFinished || limitReached;

                                    let btnText = 'Enviar Solución';
                                    if (examFinished) btnText = 'Examen Finalizado';
                                    else if (limitReached) btnText = 'Límite alcanzado';
                                    else if (submitting) btnText = 'Evaluando...';

                                    return (
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                            {tryLimitVal > 0 && (
                                                <span style={{
                                                    fontSize: '0.75rem', fontWeight: 700, color: limitReached ? '#ef4444' : '#9ca3af',
                                                    whiteSpace: 'nowrap'
                                                }}>
                                                    {chAttempts}/{tryLimitVal} intentos
                                                </span>
                                            )}
                                            <button onClick={handleSubmit} disabled={isDisabled}
                                                style={{
                                                    background: isDisabled ? '#555' : 'linear-gradient(135deg, #c8102e, #a00d25)',
                                                    color: 'white', border: 'none', borderRadius: '10px',
                                                    padding: '0.5rem 1.25rem', fontWeight: 700, fontSize: '0.85rem',
                                                    cursor: isDisabled ? 'not-allowed' : 'pointer',
                                                    display: 'flex', alignItems: 'center', gap: '6px',
                                                    opacity: isDisabled ? 0.6 : 1
                                                }}>
                                                <Send size={14} /> {btnText}
                                            </button>
                                        </div>
                                    );
                                })()}
                            </div>

                            {/* Monaco Editor */}
                            <div style={{ flex: 1 }}>
                                <Editor
                                    height="100%"
                                    theme="vs-dark"
                                    language={language}
                                    value={currentCode}
                                    onChange={handleCodeChange}
                                    options={{ minimap: { enabled: false }, fontSize: 14, padding: { top: 12 } }}
                                />
                            </div>

                            {/* Output panel: Redesigned for maximum visibility */}
                            <div style={{
                                height: output ? '250px' : '50px',
                                background: '#0d1117',
                                borderTop: '2px solid #30363d',
                                transition: 'height 0.3s ease-in-out',
                                overflow: 'auto',
                                padding: '1rem 1.5rem',
                                boxShadow: '0 -4px 15px rgba(0,0,0,0.5)'
                            }}>
                                <div style={{
                                    display: 'flex',
                                    justifyContent: 'space-between',
                                    alignItems: 'center',
                                    marginBottom: output ? '1rem' : 0,
                                    paddingBottom: output ? '0.75rem' : 0,
                                    borderBottom: output ? '1px solid rgba(255,255,255,0.1)' : 'none'
                                }}>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                        <div style={{
                                            width: '8px',
                                            height: '8px',
                                            borderRadius: '50%',
                                            background: currentResult?.status === 'accepted' ? '#10b981' : currentResult ? '#ef4444' : '#8b949e',
                                            boxShadow: currentResult ? `0 0 8px ${currentResult?.status === 'accepted' ? '#10b981' : '#ef4444'}` : 'none'
                                        }}></div>
                                        <span style={{
                                            color: '#e6edf3',
                                            fontSize: '0.9rem',
                                            fontWeight: 800,
                                            letterSpacing: '1px'
                                        }}>
                                            CONSOLA DE RESULTADOS
                                        </span>
                                    </div>

                                    {currentResult && (
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                            <span style={{
                                                fontSize: '0.85rem',
                                                fontWeight: 700,
                                                color: currentResult.status === 'accepted' ? '#10b981' : '#ef4444'
                                            }}>
                                                {currentResult.status === 'accepted' ? 'ACEPTADO' : 'FALLIDO'}
                                            </span>
                                            <span style={{
                                                fontSize: '1rem',
                                                fontWeight: 900,
                                                padding: '4px 12px',
                                                borderRadius: '8px',
                                                background: currentResult.status === 'accepted' ? 'rgba(16, 185, 129, 0.2)' : 'rgba(239, 68, 68, 0.2)',
                                                color: currentResult.status === 'accepted' ? '#10b981' : '#ef4444',
                                                border: `1px solid ${currentResult.status === 'accepted' ? '#10b981' : '#ef4444'}`
                                            }}>
                                                {currentResult.score}%
                                            </span>
                                        </div>
                                    )}
                                </div>
                                {output && (
                                    <pre style={{
                                        color: '#d1d5db',
                                        fontSize: '1rem',
                                        lineHeight: '1.6',
                                        marginTop: '0.5rem',
                                        whiteSpace: 'pre-wrap',
                                        fontFamily: '"Fira Code", "JetBrains Mono", monospace',
                                        padding: '0.5rem 0'
                                    }}>
                                        {output}
                                    </pre>
                                )}
                            </div>
                        </div>
                    </>
                )}

                {!currentChallenge && (
                    <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#fafafa' }}>
                        <div style={{ textAlign: 'center', color: '#999' }}>
                            <Target size={48} style={{ marginBottom: '1rem', opacity: 0.3 }} />
                            <h3>Este examen no tiene retos asignados</h3>
                            <p>Contacta a tu profesor para más información.</p>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ExamRunner;
