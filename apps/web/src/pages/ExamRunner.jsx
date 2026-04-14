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
    const [publicTestCasesMap, setPublicTestCasesMap] = useState({});

    // Run Test Modal State
    const [showRunModal, setShowRunModal] = useState(false);
    const [runTab, setRunTab] = useState('public');
    const [selectedPublicCase, setSelectedPublicCase] = useState('');
    const [customInputs, setCustomInputs] = useState([]);
    const [customOutput, setCustomOutput] = useState('');
    const [customOutputVarName, setCustomOutputVarName] = useState('out');
    const [customOutputVarType, setCustomOutputVarType] = useState('string');
    const [customSampleName, setCustomSampleName] = useState('custom_sample');

    // Pre-fill custom inputs based on the first public testcase
    useEffect(() => {
        if (!currentChallenge) return;
        const cases = publicTestCasesMap[currentChallenge.id];
        if (cases && cases.length > 0) {
            const firstCase = Object.assign({}, cases[0]);
            const inputs = firstCase.input || firstCase.Input || [];
            const out = firstCase.expected_output || firstCase.expectedOutput || firstCase.ExpectedOutput || {};
            
            // Set input shapes without values
            if (Array.isArray(inputs)) {
                setCustomInputs(inputs.map(i => ({ name: i.name || i.Name, type: i.type || i.Type, value: '' })));
            }
            if (out.name || out.Name) {
                setCustomOutputVarName(out.name || out.Name);
                setCustomOutputVarType(out.type || out.Type);
            }
        }
    }, [currentIndex, challenges, publicTestCasesMap]);

    // Session state
    const [sessionId, setSessionId] = useState(null);
    const [sessionStatus, setSessionStatus] = useState(null);
    const [timeLeft, setTimeLeft] = useState(null); // in seconds, null = unlimited
    const [examFinished, setExamFinished] = useState(false);
    const [attemptMap, setAttemptMap] = useState(() => {
        // Restore attempt map from localStorage for persistence across refreshes
        try {
            const stored = localStorage.getItem('exam_attempt_map');
            return stored ? JSON.parse(stored) : {};
        } catch { return {}; }
    }); // { challengeId: count }
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

        const userId = user?.id || user?.ID || '';

        // Step 1: Try to fetch the active session (bypassing the cached one entirely)
        try {
            const activeRes = await client.get(`/submissions/sessions/active?user_id=${userId}`);
            const activeSession = activeRes?.data;
            const sid = activeSession?.id || activeSession?.ID;
            const sessionExamId = activeSession?.exam_id || activeSession?.ExamID || activeSession?.examID || '';

            if (sid) {
                // Step 2: If session is for THIS exam, reuse it (e.g. student refreshed or navigated back)
                if (sessionExamId === examId || String(sessionExamId) === String(examId)) {
                    applySession(activeSession, sid);
                    return { id: sid, status: 'active', ...activeSession };
                } else {
                    // Step 3: Session is for a DIFFERENT exam -> close it and create a new one
                    try {
                        await client.post(`/submissions/sessions/${sid}/close`);
                        localStorage.removeItem('session_id');
                    } catch (closeErr) {
                        console.warn('Failed to close old active session:', closeErr);
                    }
                }
            }
        } catch (err) {
            // No active session or error fetching it, that's fine, we will create one below
            console.warn('No active session found or error fetching it:', err);
        }

        // Step 4: Try to create a new session
        try {
            const newSession = await tryCreate();
            const sid = newSession?.id || newSession?.ID;
            if (sid) {
                applySession(newSession, sid);
                return newSession;
            }
        } catch (err) {
            const apiMsg = err?.response?.data?.error || err?.message || '';
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
                        let parsedTemplates = ch.code_templates || ch.CodeTemplates || {};
                        while (typeof parsedTemplates === 'string') {
                            try { parsedTemplates = JSON.parse(parsedTemplates); } catch (e) { parsedTemplates = {}; break; }
                        }
                        if (typeof parsedTemplates !== 'object' || parsedTemplates === null || Array.isArray(parsedTemplates)) {
                            parsedTemplates = {};
                        }
                        return {
                            ...ch,
                            id: ch.id || ch.ID || item.challenge_id || item.challengeID,
                            title: ch.title || ch.Title || 'Reto',
                            description: ch.description || ch.Description || '',
                            difficulty: ch.difficulty || ch.Difficulty || 'medium',
                            constraints: ch.constraints || ch.Constraints || '',
                            points: item.points || item.Points || 0,
                            order: item.order || item.Order || 0,
                            code_templates: parsedTemplates
                        };
                    })
                    .filter(ch => ch.id)
                    .sort((a, b) => a.order - b.order);

                setChallenges(challengeList);

                // Initialize code templates
                const initialCode = {};
                const tcPromises = [];
                challengeList.forEach(ch => {
                    const templates = ch.code_templates || ch.CodeTemplates || {};
                    const langs = Object.keys(templates);
                    if (langs.length > 0) {
                        initialCode[ch.id] = templates[langs[0]];
                    } else {
                        initialCode[ch.id] = '# Escribe tu solución aquí\ndef solve():\n    pass\n';
                    }

                    tcPromises.push(client.get(`/test-cases/challenge/${ch.id}?exam_id=${id}`).then(res => ({
                        id: ch.id,
                        cases: res.data.filter(tc => 
                            (tc.type === 'public' || tc.is_sample || tc.isSample) &&
                            tc.title !== 'Custom Test Case' && tc.Title !== 'Custom Test Case'
                        )
                    })).catch(() => ({ id: ch.id, cases: [] })));
                });
                setCodeMap(initialCode);

                if (challengeList.length > 0) {
                    const templates = challengeList[0].code_templates || challengeList[0].CodeTemplates || {};
                    const langs = Object.keys(templates);
                    if (langs.length > 0) setLanguage(langs[0]);
                }

                Promise.all(tcPromises).then(results => {
                    const tcMap = {};
                    results.forEach(r => tcMap[r.id] = r.cases);
                    setPublicTestCasesMap(tcMap);
                });

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
        heartbeatRef.current = setInterval(sendHeartbeat, 15000);

        return () => {
            if (heartbeatRef.current) clearInterval(heartbeatRef.current);
        };
    }, [sessionId, navigate]);

    // We no longer close session on browser close/refresh. Let it stay active.

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
        const ch = challenges[idx];
        if (ch) {
            const templates = ch.code_templates || ch.CodeTemplates || {};
            const langs = Object.keys(templates);
            if (langs.length > 0 && !langs.includes(language)) {
                setLanguage(langs[0]);
            }
        }
    };

    // Submit solution
    const handleSubmit = async () => {
        if (examFinished) {
            setOutput('El examen ha finalizado. No se pueden realizar más envíos.');
            return;
        }

        // Enforce 1 attempt per challenge (persisted across refreshes)
        const challengeId = currentChallenge?.id;
        const currentAttempts = attemptMap[challengeId] || 0;
        if (currentAttempts >= 1) {
            setOutput('Ya has utilizado tu único intento para este reto.');
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
                score: 0,
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

                        // Track attempt count for this challenge and persist
                        setAttemptMap(prev => {
                            const updated = { ...prev, [challengeId]: (prev[challengeId] || 0) + 1 };
                            localStorage.setItem('exam_attempt_map', JSON.stringify(updated));
                            return updated;
                        });

                        // Build output
                        const lines = [`Resultado: ${score}% (${accepted}/${results.length} casos correctos) | Intento: 1/1\n`];
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
                padding: '0.75rem 1.5rem', background: '#c8102e', color: 'white', flexShrink: 0,
                boxShadow: '0 2px 8px rgba(200,16,46,0.3)'
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <Target size={20} style={{ color: '#fff' }} />
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
                    width: '260px', background: '#ffffff', color: '#1f2937',
                    display: 'flex', flexDirection: 'column', overflowY: 'auto', flexShrink: 0,
                    borderRight: '1px solid #e5e7eb'
                }}>
                    <div style={{ padding: '1rem', borderBottom: '1px solid #e5e7eb' }}>
                        <h3 style={{ margin: 0, fontSize: '0.9rem', fontWeight: 700, color: '#c8102e', letterSpacing: '0.5px' }}>RETOS DEL EXAMEN</h3>
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
                                    background: isActive ? 'rgba(200,16,46,0.08)' : 'transparent',
                                    color: '#1f2937', cursor: 'pointer', width: '100%',
                                    borderLeft: isActive ? '3px solid #c8102e' : '3px solid transparent',
                                    transition: 'all 0.2s'
                                }}
                            >
                                <div style={{
                                    width: '28px', height: '28px', borderRadius: '50%', flexShrink: 0,
                                    display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '0.75rem', fontWeight: 700,
                                    background: isSolved ? '#10b981' : isFailed ? '#ef4444' : isActive ? '#c8102e' : '#f3f4f6',
                                    color: isSolved || isFailed || isActive ? 'white' : '#6b7280'
                                }}>
                                    {isSolved ? <CheckCircle2 size={14} /> : isFailed ? <XCircle size={14} /> : idx + 1}
                                </div>
                                <div style={{ overflow: 'hidden' }}>
                                    <div style={{ fontSize: '0.85rem', fontWeight: 600, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                        {ch.title}
                                    </div>
                                    <div style={{ fontSize: '0.7rem', color: '#9ca3af' }}>{ch.points} pts • {ch.difficulty === 'easy' ? 'Fácil' : ch.difficulty === 'hard' ? 'Difícil' : 'Medio'}</div>
                                </div>
                            </button>
                        );
                    })}
                </div>

                {/* CENTER + RIGHT (Problem, Editor, Console) */}
                {currentChallenge ? (
                    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>

                        {/* TOP SECTION: Description + Editor */}
                        <div style={{ flex: 1, display: 'flex', overflow: 'hidden' }}>

                            {/* CENTER: Problem Description */}
                            <div className="problem-description" style={{
                                width: '40%', overflowY: 'auto', padding: '1.5rem', background: '#ffffff',
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

                                {publicTestCasesMap[currentChallenge.id]?.length > 0 && (
                                    <div style={{ marginTop: '2rem' }}>
                                        <h4 style={{ margin: '0 0 1rem', fontSize: '0.95rem', color: '#1f2937', fontWeight: 800 }}> Casos de Prueba </h4>
                                        {publicTestCasesMap[currentChallenge.id].map((tc, idx) => (
                                            <div key={idx} style={{ background: '#f9fafb', border: '1px solid #e5e7eb', borderRadius: '8px', padding: '1rem', marginBottom: '1rem' }}>
                                                <div style={{ fontWeight: 700, fontSize: '0.85rem', color: '#4b5563', marginBottom: '0.5rem' }}>Entrada:</div>
                                                <pre style={{ background: '#f3f4f6', padding: '0.5rem', borderRadius: '4px', fontSize: '0.85rem', color: '#1f2937', margin: '0 0 1rem 0' }}>{Array.isArray(tc.input) ? tc.input.map(i => i ? `${i.name} = ${i.value}` : 'nil').join(', ') : JSON.stringify(tc.input)}</pre>
                                                <div style={{ fontWeight: 700, fontSize: '0.85rem', color: '#4b5563', marginBottom: '0.5rem' }}>Salida Esperada:</div>
                                                <pre style={{ background: '#f3f4f6', padding: '0.5rem', borderRadius: '4px', fontSize: '0.85rem', color: '#1f2937', margin: 0 }}>{tc.expected_output?.value || tc.ExpectedOutput?.value || tc.expectedOutput?.value || ''}</pre>
                                            </div>
                                        ))}
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

                            {/* RIGHT: Code Editor */}
                            <div style={{ flex: 1, display: 'flex', flexDirection: 'column', background: '#1e1e1e', minWidth: 0, minHeight: 0 }}>
                                {/* Editor toolbar */}
                                <div style={{
                                    display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                                    padding: '0.5rem 1rem', background: '#f9fafb', borderBottom: '1px solid #e5e7eb'
                                }}>
                                    <select value={language} onChange={(e) => setLanguage(e.target.value)}
                                        style={{ background: '#ffffff', color: '#1f2937', border: '1px solid #d1d5db', borderRadius: '6px', padding: '0.35rem 0.75rem', fontSize: '0.85rem' }}>
                                        {(() => {
                                            const templates = currentChallenge?.code_templates || currentChallenge?.CodeTemplates || {};
                                            const langs = Object.keys(templates).filter(l => ['python', 'javascript', 'java', 'cpp', 'go'].includes(l));
                                            return langs.length > 0
                                                ? langs.map(l => <option key={l} value={l}>{l === 'cpp' ? 'C++' : l.charAt(0).toUpperCase() + l.slice(1)}</option>)
                                                : <option value="python">Python</option>;
                                        })()}
                                    </select>
                                    {(() => {
                                        const chAttempts = currentChallenge ? (attemptMap[currentChallenge.id] || 0) : 0;
                                        const limitReached = chAttempts >= 1;
                                        const isSubmitDisabled = submitting || examFinished || limitReached;
                                        const isTestDisabled = submitting || examFinished;

                                        let btnText = 'Enviar Solución';
                                        if (examFinished) btnText = 'Examen Finalizado';
                                        else if (limitReached) btnText = 'Ya enviado';
                                        else if (submitting) btnText = 'Evaluando...';

                                        return (
                                            <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                                <span style={{
                                                    fontSize: '0.75rem', fontWeight: 700, color: limitReached ? '#ef4444' : '#6b7280',
                                                    whiteSpace: 'nowrap'
                                                }}>
                                                    {chAttempts}/1 intentos
                                                </span>
                                                <button onClick={() => setShowRunModal(true)} disabled={isTestDisabled}
                                                    style={{
                                                        background: isTestDisabled ? '#e5e7eb' : '#ffffff',
                                                        color: isTestDisabled ? '#9ca3af' : '#c8102e', border: '1px solid #c8102e', borderRadius: '10px',
                                                        padding: '0.5rem 1.25rem', fontWeight: 700, fontSize: '0.85rem',
                                                        cursor: isTestDisabled ? 'not-allowed' : 'pointer',
                                                        display: 'flex', alignItems: 'center', gap: '6px',
                                                        opacity: isTestDisabled ? 0.6 : 1, transition: 'background 0.2s'
                                                    }}>
                                                    <Timer size={14} /> Probar Código
                                                </button>
                                                <button onClick={handleSubmit} disabled={isSubmitDisabled}
                                                    style={{
                                                        background: isSubmitDisabled ? '#e5e7eb' : '#c8102e',
                                                        color: isSubmitDisabled ? '#9ca3af' : 'white', border: 'none', borderRadius: '10px',
                                                        padding: '0.5rem 1.25rem', fontWeight: 700, fontSize: '0.85rem',
                                                        cursor: isSubmitDisabled ? 'not-allowed' : 'pointer',
                                                        display: 'flex', alignItems: 'center', gap: '6px',
                                                        opacity: isSubmitDisabled ? 0.6 : 1
                                                    }}>
                                                    <Send size={14} /> {btnText}
                                                </button>
                                            </div>
                                        );
                                    })()}
                                </div>

                                {/* Monaco Editor */}
                                <div style={{ flex: 1, minHeight: 0 }}>
                                    <Editor
                                        height="100%"
                                        theme="vs-dark"
                                        language={['python', 'javascript', 'java', 'cpp', 'go'].includes(language) ? language : 'python'}
                                        value={currentCode}
                                        onChange={handleCodeChange}
                                        options={{ minimap: { enabled: false }, fontSize: 14, padding: { top: 12 } }}
                                    />
                                </div>
                            </div>

                        </div>

                        {/* BOTTOM SECTION: Output panel */}
                        <div style={{
                            height: output ? '300px' : '50px',
                            flexShrink: 0,
                            background: '#1a1a2e',
                            borderTop: '2px solid #c8102e',
                            transition: 'height 0.3s ease-in-out',
                            overflow: 'hidden',
                            display: 'flex',
                            flexDirection: 'column'
                        }}>
                            <div style={{
                                display: 'flex',
                                justifyContent: 'space-between',
                                alignItems: 'center',
                                padding: '0.75rem 1.5rem',
                                background: '#111827',
                                borderBottom: output ? '1px solid rgba(200,16,46,0.3)' : 'none',
                                cursor: 'pointer'
                            }} onClick={() => setOutput(output ? '' : ' Esperando envío...')}>
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
                                        CONSOLA DE RESULTADOS {output ? '▼' : '▲'}
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
                            <div style={{ flex: 1, overflowY: 'auto', padding: '1rem 1.5rem' }}>
                                {output && (
                                    <pre style={{
                                        color: '#d1d5db',
                                        fontSize: '1rem',
                                        lineHeight: '1.6',
                                        margin: 0,
                                        whiteSpace: 'pre-wrap',
                                        fontFamily: '"Fira Code", "JetBrains Mono", monospace'
                                    }}>
                                        {output}
                                    </pre>
                                )}
                            </div>
                        </div>

                    </div>
                ) : (
                    <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#fff' }}>
                        <div style={{ textAlign: 'center', color: '#9ca3af' }}>
                            <Target size={48} style={{ marginBottom: '1rem', opacity: 0.3, color: '#c8102e' }} />
                            <h3>Este examen no tiene retos asignados</h3>
                            <p>Contacta a tu profesor para más información.</p>
                        </div>
                    </div>
                )}
            </div>

            {/* RUN TEST MODAL */}
            {showRunModal && (
                <div style={{
                    position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, 
                    background: 'rgba(0,0,0,0.6)', backdropFilter: 'blur(4px)',
                    display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000
                }}>
                    <div style={{
                        background: '#ffffff', borderRadius: '16px', width: '500px', maxWidth: '90%',
                        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.4)', overflow: 'hidden'
                    }}>
                        <div style={{ padding: '1.25rem 1.5rem', borderBottom: '1px solid #e5e7eb', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <h3 style={{ margin: 0, fontSize: '1.1rem', fontWeight: 700, color: '#111827', display: 'flex', alignItems: 'center', gap: '8px' }}>
                                <Timer size={18} style={{ color: '#c8102e' }} /> Probar Ejecución 
                                <span style={{ fontSize: '0.7rem', padding: '2px 8px', background: '#fef2f2', color: '#b91c1c', borderRadius: '12px' }}>Modo Prueba</span>
                            </h3>
                            <button onClick={() => setShowRunModal(false)} style={{ background: 'transparent', border: 'none', cursor: 'pointer', color: '#9ca3af' }}>
                                <XCircle size={20} />
                            </button>
                        </div>
                        <div style={{ padding: '1.5rem' }}>
                            <p style={{ margin: '0 0 1rem 0', fontSize: '0.9rem', color: '#4b5563' }}>
                                Ejecuta tu código localmente sin gastar intentos. Selecciona un caso de prueba existente o define uno nuevo.
                            </p>
                            
                            <div style={{ display: 'flex', gap: '1rem', marginBottom: '1.5rem', borderBottom: '1px solid #f3f4f6' }}>
                                <button 
                                    onClick={() => setRunTab('public')}
                                    style={{ 
                                        padding: '0.5rem 0', background: 'none', border: 'none', borderBottom: runTab === 'public' ? '2px solid #c8102e' : '2px solid transparent',
                                        color: runTab === 'public' ? '#111827' : '#6b7280', fontWeight: 600, fontSize: '0.9rem', cursor: 'pointer', transition: 'all 0.2s'
                                    }}>
                                    Casos Públicos
                                </button>
                                <button 
                                    onClick={() => setRunTab('custom')}
                                    style={{ 
                                        padding: '0.5rem 0', background: 'none', border: 'none', borderBottom: runTab === 'custom' ? '2px solid #c8102e' : '2px solid transparent',
                                        color: runTab === 'custom' ? '#111827' : '#6b7280', fontWeight: 600, fontSize: '0.9rem', cursor: 'pointer', transition: 'all 0.2s'
                                    }}>
                                    Caso Personalizado
                                </button>
                            </div>

                            {runTab === 'public' && (
                                <div style={{ marginBottom: '1rem' }}>
                                    <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: '#374151', marginBottom: '0.5rem' }}>Seleccionar Caso de Prueba:</label>
                                    <select 
                                        value={selectedPublicCase} 
                                        onChange={(e) => setSelectedPublicCase(e.target.value)}
                                        style={{ width: '100%', padding: '0.75rem', borderRadius: '8px', border: '1px solid #d1d5db', fontSize: '0.9rem' }}
                                    >
                                        <option value="">Selecciona un caso...</option>
                                        {(publicTestCasesMap[currentChallenge?.id] || []).map((tc, idx) => (
                                            <option key={idx} value={tc.id || idx}>Caso #{idx + 1} ({tc.expected_output?.value || tc.expectedOutput?.value || 'n/a'})</option>
                                        ))}
                                    </select>
                                </div>
                            )}

                            {runTab === 'custom' && (
                                <div style={{ background: '#f9fafb', padding: '1rem', borderRadius: '8px', border: '1px solid #e5e7eb' }}>
                                    <div style={{ marginBottom: '0.75rem' }}>
                                        <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: '#374151', marginBottom: '0.25rem' }}>Nombre del Caso</label>
                                        <input type="text" value={customSampleName} onChange={e => setCustomSampleName(e.target.value)} style={{ width: '100%', padding: '0.5rem', borderRadius: '6px', border: '1px solid #d1d5db', fontSize: '0.85rem' }} />
                                    </div>
                                    <div style={{ marginBottom: '0.75rem' }}>
                                        <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: '#374151', marginBottom: '0.25rem' }}>Inputs (Variables)</label>
                                        {customInputs.map((inp, i) => (
                                            <div key={i} style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.5rem', alignItems: 'center' }}>
                                                <span style={{ fontSize: '0.85rem', fontWeight: 600, color: '#4b5563', width: '80px' }}>{inp.name} ({inp.type}):</span>
                                                <input placeholder="Valor..." value={inp.value} onChange={e => { const n = [...customInputs]; n[i].value = e.target.value; setCustomInputs(n); }} style={{ flex: 1, padding: '0.5rem', borderRadius: '6px', border: '1px solid #d1d5db', fontSize: '0.85rem' }} />
                                            </div>
                                        ))}
                                    </div>
                                    <div style={{ marginBottom: '0.5rem' }}>
                                        <label style={{ display: 'block', fontSize: '0.85rem', fontWeight: 600, color: '#374151', marginBottom: '0.25rem' }}>Salida Esperada ({customOutputVarName})</label>
                                        <input type="text" value={customOutput} onChange={e => setCustomOutput(e.target.value)} placeholder="Valor de salida esperado" style={{ width: '100%', padding: '0.5rem', borderRadius: '6px', border: '1px solid #d1d5db', fontSize: '0.85rem' }} />
                                    </div>
                                </div>
                            )}

                        </div>
                        <div style={{ padding: '1rem 1.5rem', background: '#f9fafb', borderTop: '1px solid #e5e7eb', display: 'flex', justifyContent: 'flex-end', gap: '1rem' }}>
                            <button onClick={() => setShowRunModal(false)} style={{ background: 'transparent', border: '1px solid #d1d5db', borderRadius: '8px', padding: '0.5rem 1rem', fontWeight: 600, color: '#374151', cursor: 'pointer' }}>
                                Cancelar
                            </button>
                            <button onClick={async () => {
                                setShowRunModal(false);
                                setOutput('⏳ Ejecutando prueba local en el servidor...\n\n');
                                
                                try {
                                    let res;
                                    if (runTab === 'public') {
                                        // Find the selected public case and run it via execute-custom
                                        const allCases = publicTestCasesMap[currentChallenge?.id] || [];
                                        const selected = allCases.find(tc => (tc.id || '') === selectedPublicCase) || allCases[parseInt(selectedPublicCase)] || allCases[0];
                                        if (!selected) {
                                            setOutput('No hay caso de prueba público seleccionado.');
                                            return;
                                        }
                                        const inputs = (selected.input || selected.Input || []).map(v => ({
                                            name: v.name || v.Name,
                                            type: v.type || v.Type,
                                            value: String(v.value ?? v.Value ?? '')
                                        }));
                                        const expectedOut = selected.expected_output || selected.expectedOutput || selected.ExpectedOutput || {};
                                        res = await client.post('/submissions/execute-custom', {
                                            code: currentCode,
                                            language: language,
                                            challenge_id: currentChallenge?.id,
                                            session_id: sessionId,
                                            input_variables: inputs,
                                            output_variable: {
                                                name: expectedOut.name || expectedOut.Name || 'out',
                                                type: expectedOut.type || expectedOut.Type || 'string',
                                                value: String(expectedOut.value ?? expectedOut.Value ?? '')
                                            }
                                        });
                                    } else {
                                        res = await client.post('/submissions/execute-custom', {
                                            code: currentCode,
                                            language: language,
                                            challenge_id: currentChallenge?.id,
                                            session_id: sessionId,
                                            input_variables: customInputs,
                                            output_variable: {
                                                name: customOutputVarName,
                                                type: customOutputVarType,
                                                value: customOutput
                                            }
                                        });
                                    }

                                    let submissionId = res.data?.id || res.data?.ID;
                                    if (!submissionId && res.data?.Submission) {
                                        submissionId = res.data.Submission.id || res.data.Submission.ID;
                                    }

                                    if (!submissionId) {
                                        setOutput('La API no retornó un ID de ejecución válido.\n\nResultados crudos:\n' + JSON.stringify(res.data, null, 2));
                                        return;
                                    }

                                    setOutput('Prueba encolada. Ejecutando...');
                                    
                                    // Poll for results
                                    for (let attempt = 0; attempt < 40; attempt++) {
                                        const pollRes = await client.get(`/submissions/${submissionId}`);
                                        const results = pollRes?.data?.Results || pollRes?.data?.results || [];

                                        if (Array.isArray(results) && results.length > 0) {
                                            const hasPending = results.some(r => {
                                                const s = String(r?.Status || r?.status || '').toLowerCase();
                                                return s === 'queued' || s === 'running';
                                            });

                                            if (!hasPending) {
                                                const accepted = results.filter(r => String(r?.Status || r?.status || '').toLowerCase() === 'accepted').length;
                                                const score = Math.round((accepted / results.length) * 100);
                                                const lines = [`Resultado de Prueba: ${score}% (${accepted}/${results.length} casos correctos)\n`];
                                                results.forEach((r, i) => {
                                                    const st = (r.Status || r.status || 'unknown').toLowerCase();
                                                    const err = r.ErrorMessage || r.errorMessage || '';
                                                    lines.push(`  Caso ${i + 1}: ${st === 'accepted' ? '✅' : '❌'} ${st}${err ? ' - ' + err : ''}`);
                                                });
                                                lines.push(`\n(Resultados no afectan tus puntajes)`);
                                                setOutput(lines.join('\n'));
                                                return;
                                            }
                                        }
                                        await sleep(1000);
                                        setOutput(`Ejecutando pruebas... (${attempt + 1}s)`);
                                    }
                                    
                                    setOutput('Tiempo de espera agotado para la ejecución local.');
                                } catch (err) {
                                    const errorMsg = err?.response?.data?.error || err.message;
                                    setOutput((prev) => prev + `Status: ❌ Error en la ejecución.\nMotivo: ${errorMsg}`);
                                }
                            }} style={{ background: '#c8102e', color: 'white', border: 'none', borderRadius: '8px', padding: '0.5rem 1.25rem', fontWeight: 700, cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '6px' }}>
                                <Timer size={16} /> Ejecutar Prueba
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ExamRunner;
