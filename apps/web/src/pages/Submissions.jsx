import { useState, useEffect, useContext, useCallback } from 'react';
import {
    CheckCircle, XCircle, Clock, Code,
    Calendar, ChevronRight, ChevronDown, AlertCircle,
    Trophy, RotateCcw, Target, Users,
    Hash, User, Layers, BookOpen
} from 'lucide-react';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import './Submissions.css';

const Submissions = () => {
    const { user } = useContext(AuthContext);
    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';
    const userId = user?.id || user?.ID;

    const [exams, setExams] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Expanded exam state
    const [expandedExamId, setExpandedExamId] = useState(null);
    const [examDetails, setExamDetails] = useState({}); // { examId: { items, submissions, loading, loaded } }

    // Professor: expanded student within an exam
    const [expandedStudentId, setExpandedStudentId] = useState(null);

    // Fetch exams on mount
    useEffect(() => {
        const fetchExams = async () => {
            if (!user) return;
            try {
                const { data } = await client.get('/exams/public');
                const examList = Array.isArray(data) ? data : (data?.items || data?.exams || []);
                setExams(examList);
            } catch (err) {
                console.error('Error loading exams:', err);
                setError('No se pudieron cargar los exámenes.');
            } finally {
                setLoading(false);
            }
        };
        if (user) fetchExams();
    }, [user]);

    // Load exam details (items + submissions) lazily
    const loadExamDetails = useCallback(async (examId) => {
        if (examDetails[examId]?.loaded) return;

        setExamDetails(prev => ({
            ...prev,
            [examId]: { ...prev[examId], loading: true }
        }));

        try {
            // Fetch exam items (challenges)
            const itemsRes = await client.get(`/exams/${examId}/items`);
            const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);

            // For each challenge, fetch submissions
            const allSubmissions = [];
            for (const item of items) {
                const challengeId = item.challenge?.id || item.challenge?.ID || item.challengeID || item.challenge_id;
                if (!challengeId) continue;
                try {
                    const subRes = await client.get(`/submissions/challenge/${challengeId}`);
                    const subs = Array.isArray(subRes.data) ? subRes.data : (subRes.data?.items || []);
                    subs.forEach(sub => {
                        // Attach challenge info to each submission for display
                        const s = sub?.Submission || sub?.submission || sub;
                        const results = sub?.Results || sub?.results || [];
                        allSubmissions.push({
                            ...s,
                            challengeTitle: item.challenge?.title || item.challenge?.Title || 'Reto',
                            challengeId,
                            points: item.points || item.Points || 0,
                            results,
                        });
                    });
                } catch (subErr) {
                    // Some challenges may not have submissions
                    console.warn(`No submissions for challenge ${challengeId}:`, subErr.message);
                }
            }

            setExamDetails(prev => ({
                ...prev,
                [examId]: {
                    items,
                    submissions: allSubmissions,
                    loading: false,
                    loaded: true
                }
            }));
        } catch (err) {
            console.error('Error loading exam details:', err);
            setExamDetails(prev => ({
                ...prev,
                [examId]: { ...prev[examId], loading: false, error: 'Error al cargar detalles' }
            }));
        }
    }, [examDetails]);

    const toggleExam = (examId) => {
        if (expandedExamId === examId) {
            setExpandedExamId(null);
            setExpandedStudentId(null);
        } else {
            setExpandedExamId(examId);
            setExpandedStudentId(null);
            loadExamDetails(examId);
        }
    };

    const getStatusInfo = (status) => {
        const s = (status || 'pending').toLowerCase();
        if (s === 'accepted' || s === 'success') return { label: 'Aceptado', cls: 'accepted', icon: <CheckCircle size={13} /> };
        if (s === 'wrong_answer' || s === 'rejected' || s === 'failed') return { label: 'Rechazado', cls: 'rejected', icon: <XCircle size={13} /> };
        if (s === 'runtime_error' || s === 'error') return { label: 'Error', cls: 'error', icon: <AlertCircle size={13} /> };
        return { label: 'Pendiente', cls: 'pending', icon: <Clock size={13} /> };
    };

    const formatDate = (dateStr) => {
        if (!dateStr) return '—';
        const d = new Date(dateStr);
        return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' });
    };

    const formatTime = (dateStr) => {
        if (!dateStr) return '';
        const d = new Date(dateStr);
        return d.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
    };

    // ============ STUDENT VIEW ============
    const renderStudentView = () => {
        return exams.map(exam => {
            const examId = exam.id || exam.ID;
            const isExpanded = expandedExamId === examId;
            const details = examDetails[examId];

            // Calculate summary from loaded submissions
            let totalChallenges = 0;
            let solvedChallenges = 0;
            let totalAttempts = 0;
            let bestScoreMap = {};

            if (details?.loaded) {
                // Group submissions by challenge, find best score per challenge
                const challengeMap = {};
                details.submissions.forEach(sub => {
                    const cid = sub.challengeId || sub.challenge_id || sub.ChallengeID;
                    if (!challengeMap[cid]) challengeMap[cid] = [];
                    challengeMap[cid].push(sub);
                });

                totalChallenges = details.items.length;
                Object.entries(challengeMap).forEach(([cid, subs]) => {
                    totalAttempts += subs.length;
                    const best = Math.max(...subs.map(s => s.score || s.Score || 0));
                    bestScoreMap[cid] = best;
                    if (best === 100) solvedChallenges++;
                });
            }

            return (
                <div key={examId} className={`exam-card-wrapper ${isExpanded ? 'expanded' : ''}`}>
                    <div className="exam-card" onClick={() => toggleExam(examId)}>
                        <div className="exam-card-left">
                            <div className="exam-icon-box">
                                <Target size={20} />
                            </div>
                            <div className="exam-info">
                                <h3 className="exam-title">{exam.title || exam.Title || 'Examen'}</h3>
                                <div className="exam-meta-row">
                                    <span className="exam-meta-item">
                                        <Calendar size={12} />
                                        {formatDate(exam.startTime || exam.StartTime || exam.start_time || exam.created_at || exam.CreatedAt)}
                                    </span>
                                    {(exam.timeLimit || exam.TimeLimit || exam.time_limit) > 0 && (
                                        <span className="exam-meta-item">
                                            <Clock size={12} />
                                            {Math.floor((exam.timeLimit || exam.TimeLimit || exam.time_limit) / 60)} min
                                        </span>
                                    )}
                                </div>
                            </div>
                        </div>
                        <div className="exam-card-right">
                            {details?.loaded && (
                                <div className="exam-score-summary">
                                    <div className="score-chip">
                                        <Trophy size={14} />
                                        {solvedChallenges}/{totalChallenges}
                                    </div>
                                    <span className="attempt-count">{totalAttempts} envío{totalAttempts !== 1 ? 's' : ''}</span>
                                </div>
                            )}
                            <div className={`expand-chevron ${isExpanded ? 'open' : ''}`}>
                                <ChevronRight size={18} />
                            </div>
                        </div>
                    </div>

                    {isExpanded && (
                        <div className="exam-detail-panel">
                            {details?.loading && (
                                <div className="detail-loading">
                                    <div className="spinner-small"></div>
                                    <span>Cargando resultados...</span>
                                </div>
                            )}
                            {details?.error && (
                                <div className="detail-error">{details.error}</div>
                            )}
                            {details?.loaded && (
                                <>
                                    {details.items.length === 0 ? (
                                        <div className="detail-empty">Este examen no tiene retos asignados.</div>
                                    ) : (
                                        <div className="challenge-results-grid">
                                            {details.items.map((item, idx) => {
                                                const ch = item.challenge || {};
                                                const cid = ch.id || ch.ID || item.challengeID || item.challenge_id;
                                                const title = ch.title || ch.Title || `Reto #${idx + 1}`;
                                                const pts = item.points || item.Points || 0;
                                                const diff = (ch.difficulty || ch.Difficulty || 'medium').toLowerCase();

                                                // Find submissions for this challenge
                                                const mySubs = details.submissions
                                                    .filter(s => (s.challengeId || s.challenge_id || s.ChallengeID) === cid)
                                                    .sort((a, b) => new Date(b.created_at || b.CreatedAt || 0) - new Date(a.created_at || a.CreatedAt || 0));

                                                const bestScore = mySubs.length > 0 ? Math.max(...mySubs.map(s => s.score || s.Score || 0)) : null;
                                                const isSolved = bestScore === 100;
                                                const lastSub = mySubs[0];

                                                return (
                                                    <div key={cid || idx} className={`challenge-result-card ${isSolved ? 'solved' : mySubs.length > 0 ? 'attempted' : 'unattempted'}`}>
                                                        <div className="cr-header">
                                                            <div className="cr-title-row">
                                                                <Code size={15} className="cr-icon" />
                                                                <span className="cr-title">{title}</span>
                                                            </div>
                                                            <div className="cr-badges">
                                                                <span className={`diff-badge ${diff}`}>
                                                                    {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                                </span>
                                                                <span className="pts-badge">{pts} pts</span>
                                                            </div>
                                                        </div>
                                                        <div className="cr-body">
                                                            {mySubs.length === 0 ? (
                                                                <div className="cr-no-submissions">Sin intentos</div>
                                                            ) : (
                                                                <>
                                                                    <div className="cr-score-row">
                                                                        <div className={`cr-score-circle ${isSolved ? 'perfect' : bestScore >= 50 ? 'partial' : 'low'}`}>
                                                                            {bestScore}%
                                                                        </div>
                                                                        <div className="cr-score-details">
                                                                            <span className="cr-attempts">
                                                                                <Hash size={12} /> {mySubs.length} intento{mySubs.length !== 1 ? 's' : ''}
                                                                            </span>
                                                                            {lastSub && (
                                                                                <span className="cr-last-date">
                                                                                    <Calendar size={11} /> {formatDate(lastSub.created_at || lastSub.CreatedAt)}
                                                                                </span>
                                                                            )}
                                                                        </div>
                                                                    </div>
                                                                    {/* Mini timeline of attempts */}
                                                                    <div className="cr-attempts-timeline">
                                                                        {mySubs.slice(0, 5).map((sub, si) => {
                                                                            const sc = sub.score || sub.Score || 0;
                                                                            return (
                                                                                <div key={sub.id || sub.ID || si}
                                                                                    className={`attempt-dot ${sc === 100 ? 'perfect' : sc >= 50 ? 'partial' : 'low'}`}
                                                                                    title={`${sc}% — ${formatDate(sub.created_at || sub.CreatedAt)} ${formatTime(sub.created_at || sub.CreatedAt)}`}
                                                                                />
                                                                            );
                                                                        })}
                                                                        {mySubs.length > 5 && (
                                                                            <span className="more-attempts">+{mySubs.length - 5}</span>
                                                                        )}
                                                                    </div>
                                                                </>
                                                            )}
                                                        </div>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    )}
                                </>
                            )}
                        </div>
                    )}
                </div>
            );
        });
    };

    // ============ PROFESSOR VIEW ============
    const renderProfessorView = () => {
        return exams.map(exam => {
            const examId = exam.id || exam.ID;
            const isExpanded = expandedExamId === examId;
            const details = examDetails[examId];

            // Group submissions by student
            let studentMap = {};
            let totalStudents = 0;
            if (details?.loaded) {
                details.submissions.forEach(sub => {
                    const uid = sub.user_id || sub.UserID || sub.userId || 'unknown';
                    if (!studentMap[uid]) studentMap[uid] = [];
                    studentMap[uid].push(sub);
                });
                totalStudents = Object.keys(studentMap).length;
            }

            return (
                <div key={examId} className={`exam-card-wrapper ${isExpanded ? 'expanded' : ''}`}>
                    <div className="exam-card" onClick={() => toggleExam(examId)}>
                        <div className="exam-card-left">
                            <div className="exam-icon-box professor">
                                <BookOpen size={20} />
                            </div>
                            <div className="exam-info">
                                <h3 className="exam-title">{exam.title || exam.Title || 'Examen'}</h3>
                                <div className="exam-meta-row">
                                    <span className="exam-meta-item">
                                        <Calendar size={12} />
                                        {formatDate(exam.startTime || exam.StartTime || exam.start_time || exam.created_at || exam.CreatedAt)}
                                    </span>
                                    {details?.loaded && (
                                        <span className="exam-meta-item highlight">
                                            <Users size={12} />
                                            {totalStudents} estudiante{totalStudents !== 1 ? 's' : ''}
                                        </span>
                                    )}
                                </div>
                            </div>
                        </div>
                        <div className="exam-card-right">
                            {details?.loaded && (
                                <div className="exam-score-summary">
                                    <span className="attempt-count">{details.submissions.length} envío{details.submissions.length !== 1 ? 's' : ''} total</span>
                                </div>
                            )}
                            <div className={`expand-chevron ${isExpanded ? 'open' : ''}`}>
                                <ChevronRight size={18} />
                            </div>
                        </div>
                    </div>

                    {isExpanded && (
                        <div className="exam-detail-panel">
                            {details?.loading && (
                                <div className="detail-loading">
                                    <div className="spinner-small"></div>
                                    <span>Cargando resultados...</span>
                                </div>
                            )}
                            {details?.error && (
                                <div className="detail-error">{details.error}</div>
                            )}
                            {details?.loaded && totalStudents === 0 && (
                                <div className="detail-empty">Ningún estudiante ha realizado envíos en este examen.</div>
                            )}
                            {details?.loaded && totalStudents > 0 && (
                                <div className="students-table">
                                    <div className="students-table-header">
                                        <div className="st-col student">ESTUDIANTE</div>
                                        <div className="st-col subs">ENVÍOS</div>
                                        <div className="st-col best">MEJOR SCORE</div>
                                        <div className="st-col challenges">RETOS RESUELTOS</div>
                                        <div className="st-col last">ÚLTIMO ENVÍO</div>
                                        <div className="st-col action"></div>
                                    </div>
                                    {Object.entries(studentMap).map(([uid, subs]) => {
                                        const isStudentExpanded = expandedStudentId === uid;
                                        // Per-challenge best scores
                                        const challengeScores = {};
                                        subs.forEach(s => {
                                            const cid = s.challengeId || s.challenge_id || s.ChallengeID;
                                            const sc = s.score || s.Score || 0;
                                            if (!challengeScores[cid] || sc > challengeScores[cid]) {
                                                challengeScores[cid] = sc;
                                            }
                                        });
                                        const solvedCount = Object.values(challengeScores).filter(sc => sc === 100).length;
                                        const bestOverall = subs.length > 0 ? Math.max(...subs.map(s => s.score || s.Score || 0)) : 0;
                                        const lastSub = subs.sort((a, b) =>
                                            new Date(b.created_at || b.CreatedAt || 0) - new Date(a.created_at || a.CreatedAt || 0)
                                        )[0];

                                        return (
                                            <div key={uid} className="student-row-wrapper">
                                                <div
                                                    className={`student-row ${isStudentExpanded ? 'expanded' : ''}`}
                                                    onClick={() => setExpandedStudentId(isStudentExpanded ? null : uid)}
                                                >
                                                    <div className="st-col student">
                                                        <div className="student-avatar">
                                                            <User size={14} />
                                                        </div>
                                                        <span className="student-id" title={uid}>
                                                            {uid.length > 16 ? uid.slice(0, 8) + '…' : uid}
                                                        </span>
                                                    </div>
                                                    <div className="st-col subs">
                                                        <span className="sub-count-badge">{subs.length}</span>
                                                    </div>
                                                    <div className="st-col best">
                                                        <span className={`score-pill ${bestOverall === 100 ? 'perfect' : bestOverall >= 50 ? 'partial' : 'low'}`}>
                                                            {bestOverall}%
                                                        </span>
                                                    </div>
                                                    <div className="st-col challenges">
                                                        <span>{solvedCount}/{details.items.length}</span>
                                                    </div>
                                                    <div className="st-col last">
                                                        <span className="date-small">{formatDate(lastSub?.created_at || lastSub?.CreatedAt)}</span>
                                                    </div>
                                                    <div className="st-col action">
                                                        <div className={`expand-chevron-small ${isStudentExpanded ? 'open' : ''}`}>
                                                            <ChevronDown size={14} />
                                                        </div>
                                                    </div>
                                                </div>

                                                {isStudentExpanded && (
                                                    <div className="student-submissions-detail">
                                                        {subs.sort((a, b) =>
                                                            new Date(b.created_at || b.CreatedAt || 0) - new Date(a.created_at || a.CreatedAt || 0)
                                                        ).map((sub, si) => {
                                                            const sc = sub.score || sub.Score || 0;
                                                            const st = getStatusInfo(sc === 100 ? 'accepted' : sc > 0 ? 'wrong_answer' : 'pending');
                                                            return (
                                                                <div key={sub.id || sub.ID || si} className="sub-detail-row">
                                                                    <div className="sub-detail-challenge">
                                                                        <Code size={13} />
                                                                        <span>{sub.challengeTitle || 'Reto'}</span>
                                                                    </div>
                                                                    <div className="sub-detail-status">
                                                                        <span className={`status-pill-mini ${st.cls}`}>
                                                                            {st.icon}
                                                                            {st.label}
                                                                        </span>
                                                                    </div>
                                                                    <div className="sub-detail-score">
                                                                        <Trophy size={12} />
                                                                        <span>{sc}%</span>
                                                                    </div>
                                                                    <div className="sub-detail-lang">
                                                                        <span>{sub.language || sub.Language || '—'}</span>
                                                                    </div>
                                                                    <div className="sub-detail-date">
                                                                        <Calendar size={11} />
                                                                        <span>{formatDate(sub.created_at || sub.CreatedAt)} {formatTime(sub.created_at || sub.CreatedAt)}</span>
                                                                    </div>
                                                                </div>
                                                            );
                                                        })}
                                                    </div>
                                                )}
                                            </div>
                                        );
                                    })}
                                </div>
                            )}
                        </div>
                    )}
                </div>
            );
        });
    };

    // ============ RENDER ============
    if (loading) return (
        <div className="submissions-page-mini">
            <header className="page-header-mini">
                <div className="skeleton title-skeleton"></div>
            </header>
            <div className="skeleton-table-mini">
                {[...Array(4)].map((_, i) => (
                    <div key={i} className="skeleton-row-mini shimmer"></div>
                ))}
            </div>
        </div>
    );

    return (
        <div className="submissions-page-mini">
            <header className="page-header-mini">
                <div className="header-info-mini">
                    <h1>{isProfessor ? 'Resultados de Exámenes' : 'Mi Historial de Exámenes'}</h1>
                    <p>{isProfessor
                        ? 'Revisa los envíos de tus estudiantes por examen'
                        : 'Consulta tu progreso y resultados en cada examen'}
                    </p>
                </div>
            </header>

            {error ? (
                <div className="error-state-mini">
                    <AlertCircle size={40} />
                    <h3>Oops! Algo salió mal</h3>
                    <p>{error}</p>
                    <button onClick={() => window.location.reload()} className="btn-retry">
                        <RotateCcw size={16} /> Reintentar carga
                    </button>
                </div>
            ) : exams.length === 0 ? (
                <div className="empty-state-mini">
                    <div className="icon-circle-mini">
                        <Layers size={32} />
                    </div>
                    <h3>Sin exámenes disponibles</h3>
                    <p>{isProfessor
                        ? 'Los exámenes que crees aparecerán aquí con los resultados de tus estudiantes.'
                        : 'Los exámenes públicos que estén disponibles aparecerán aquí.'}
                    </p>
                </div>
            ) : (
                <div className="exams-list">
                    {isProfessor ? renderProfessorView() : renderStudentView()}
                </div>
            )}
        </div>
    );
};

export default Submissions;
