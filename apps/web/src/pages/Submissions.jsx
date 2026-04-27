import { useState, useEffect, useContext, useCallback } from 'react';
import {
    CheckCircle, XCircle, Clock, Code,
    Calendar, ChevronRight, ChevronDown, AlertCircle,
    Trophy, RotateCcw, Target, Users,
    Hash, User, Layers, BookOpen
} from 'lucide-react';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import PageLoader from '../components/PageLoader';
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
    // Student / Professor: expanded challenge within detail
    const [expandedChallengeId, setExpandedChallengeId] = useState(null);
    // Professor: expanded attempt (submission) within a student drill-down
    const [expandedAttemptId, setExpandedAttemptId] = useState(null);
    const [selectedAttempt, setSelectedAttempt] = useState(null);

    // Fetch exams on mount
    useEffect(() => {
        const fetchExams = async () => {
            if (!user) return;
            try {
                // Fetch public exams
                const pubRes = await client.get('/exams/public');
                const pubExams = Array.isArray(pubRes.data) ? pubRes.data : (pubRes.data?.items || pubRes.data?.exams || []);
                
                // Fetch course exams
                const courseScope = isProfessor ? '?scope=owned' : '?scope=enrolled';
                const courseRes = await client.get(`/courses${courseScope}`);
                const courses = Array.isArray(courseRes.data) ? courseRes.data : (courseRes.data?.items || courseRes.data?.courses || []);
                
                const courseExamPromises = courses.map(c => client.get(`/exams/course/${c.id || c.ID}`));
                const courseExamResults = await Promise.allSettled(courseExamPromises);
                
                let courseExams = [];
                courseExamResults.forEach(res => {
                    if (res.status === 'fulfilled') {
                        const exams = Array.isArray(res.value.data) ? res.value.data : (res.value.data?.items || res.value.data?.exams || []);
                        courseExams = [...courseExams, ...exams];
                    }
                });

                // Merge and deduplicate
                const allExamsMap = new Map();
                pubExams.forEach(e => allExamsMap.set(e.id || e.ID, e));
                courseExams.forEach(e => allExamsMap.set(e.id || e.ID, e));
                
                setExams(Array.from(allExamsMap.values()));
            } catch (err) {
                console.error('Error loading exams:', err);
                setError('No se pudieron cargar los exámenes.');
            } finally {
                setLoading(false);
            }
        };
        if (user) fetchExams();
    }, [user, isProfessor]);

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
            setExpandedChallengeId(null);
            setExpandedAttemptId(null);
        } else {
            setExpandedExamId(examId);
            setExpandedStudentId(null);
            setExpandedChallengeId(null);
            setExpandedAttemptId(null);
            loadExamDetails(examId);
        }
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

            // Calculate weighted score
            let totalMaxPoints = 0;
            let totalEarnedPoints = 0;
            let totalChallenges = 0;
            let solvedChallenges = 0;

            if (details?.loaded) {
                const challengeBestScores = {};
                details.submissions.forEach(sub => {
                    const cid = sub.challengeId || sub.challenge_id || sub.ChallengeID;
                    const sc = sub.score || sub.Score || 0;
                    if (!challengeBestScores[cid] || sc > challengeBestScores[cid]) {
                        challengeBestScores[cid] = sc;
                    }
                });

                totalChallenges = details.items.length;
                details.items.forEach(item => {
                    const cid = item.challenge?.id || item.challenge?.ID || item.challengeID || item.challenge_id;
                    const pts = item.points || item.Points || 0;
                    totalMaxPoints += pts;
                    const best = challengeBestScores[cid] || 0;
                    totalEarnedPoints += Math.round((best / 100) * pts);
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
                                        {totalEarnedPoints}/{totalMaxPoints} pts
                                    </div>
                                    <span className="attempt-count">{solvedChallenges}/{totalChallenges} resueltos</span>
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
                                    <div className="rc-results-loading-shell">
                                        <PageLoader
                                            message="Cargando resultados del examen..."
                                            compact
                                            minHeight="0"
                                            size={16}
                                        />
                                        <div className="rc-results-skeleton" aria-hidden="true">
                                            <div className="rc-results-skeleton-row"></div>
                                            <div className="rc-results-skeleton-row"></div>
                                            <div className="rc-results-skeleton-row"></div>
                                        </div>
                                    </div>
                                </div>
                            )}
                            {details?.error && (
                                <div className="detail-error">{details.error}</div>
                            )}
                            {details?.loaded && (
                                <>
                                    {/* Score summary bar */}
                                    {totalMaxPoints > 0 && (
                                        <div style={{ padding: '1rem 0 0.5rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
                                            <div style={{ flex: 1, height: '8px', background: '#e2e8f0', borderRadius: '100px', overflow: 'hidden' }}>
                                                <div style={{
                                                    width: `${Math.round((totalEarnedPoints / totalMaxPoints) * 100)}%`,
                                                    height: '100%',
                                                    background: totalEarnedPoints === totalMaxPoints
                                                        ? 'linear-gradient(90deg, #16a34a, #22c55e)'
                                                        : totalEarnedPoints > 0
                                                            ? 'linear-gradient(90deg, #f59e0b, #fbbf24)'
                                                            : '#ef4444',
                                                    borderRadius: '100px',
                                                    transition: 'width 0.5s ease'
                                                }} />
                                            </div>
                                            <span style={{ fontWeight: 900, fontSize: '0.85rem', color: '#1e293b', whiteSpace: 'nowrap' }}>
                                                {Math.round((totalEarnedPoints / totalMaxPoints) * 100)}%
                                            </span>
                                        </div>
                                    )}

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

                                                const mySubs = details.submissions
                                                    .filter(s => (s.challengeId || s.challenge_id || s.ChallengeID) === cid)
                                                    .sort((a, b) => new Date(b.created_at || b.CreatedAt || 0) - new Date(a.created_at || a.CreatedAt || 0));

                                                const bestScore = mySubs.length > 0 ? Math.max(...mySubs.map(s => s.score || s.Score || 0)) : null;
                                                const isSolved = bestScore === 100;
                                                const earned = bestScore !== null ? Math.round((bestScore / 100) * pts) : 0;
                                                const isChExpanded = expandedChallengeId === cid;

                                                return (
                                                    <div key={cid || idx} className={`challenge-result-card ${isSolved ? 'solved' : mySubs.length > 0 ? 'attempted' : 'unattempted'}`}
                                                        onClick={() => setExpandedChallengeId(isChExpanded ? null : cid)}
                                                        style={{ cursor: 'pointer' }}
                                                    >
                                                        <div className="cr-header">
                                                            <div className="cr-title-row">
                                                                <Code size={15} className="cr-icon" />
                                                                <span className="cr-title">{title}</span>
                                                            </div>
                                                            <div className="cr-badges">
                                                                <span className={`diff-badge ${diff}`}>
                                                                    {diff === 'easy' ? 'Fácil' : diff === 'hard' ? 'Difícil' : 'Medio'}
                                                                </span>
                                                                <span className="pts-badge">{earned}/{pts} pts</span>
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
                                                                        </div>
                                                                    </div>
                                                                    {/* Expandable detail: list attempts */}
                                                                    {isChExpanded && (
                                                                        <div style={{ marginTop: '0.75rem', borderTop: '1px solid #e2e8f0', paddingTop: '0.75rem' }}>
                                                                            {mySubs.map((sub, si) => {
                                                                                const sc = sub.score || sub.Score || 0;
                                                                                return (
                                                                                    <div key={sub.id || sub.ID || si} style={{
                                                                                        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                                                                                        padding: '0.4rem 0.5rem', borderRadius: '8px', marginBottom: '3px',
                                                                                        background: sc === 100 ? '#f0fdf4' : '#f8fafc', fontSize: '0.78rem', fontWeight: 700
                                                                                    }}>
                                                                                        <span style={{ display: 'flex', alignItems: 'center', gap: '6px', color: '#475569' }}>
                                                                                            {sc === 100 ? <CheckCircle size={12} style={{ color: '#16a34a' }} /> : <XCircle size={12} style={{ color: '#ef4444' }} />}
                                                                                            Intento #{si + 1}
                                                                                        </span>
                                                                                        <span style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                                                                                            <span style={{
                                                                                                padding: '2px 8px', borderRadius: '100px', fontSize: '0.7rem', fontWeight: 900,
                                                                                                background: sc === 100 ? '#dcfce7' : sc >= 50 ? '#fef3c7' : '#fee2e2',
                                                                                                color: sc === 100 ? '#16a34a' : sc >= 50 ? '#d97706' : '#ef4444'
                                                                                            }}>{sc}%</span>
                                                                                            <span style={{ color: '#94a3b8', fontSize: '0.7rem' }}>
                                                                                                {formatDate(sub.created_at || sub.CreatedAt)} {formatTime(sub.created_at || sub.CreatedAt)}
                                                                                            </span>
                                                                                        </span>
                                                                                    </div>
                                                                                );
                                                                            })}
                                                                        </div>
                                                                    )}
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
                                    <div className="rc-results-loading-shell">
                                        <PageLoader
                                            message="Cargando resultados del examen..."
                                            compact
                                            minHeight="0"
                                            size={16}
                                        />
                                        <div className="rc-results-skeleton" aria-hidden="true">
                                            <div className="rc-results-skeleton-row"></div>
                                            <div className="rc-results-skeleton-row"></div>
                                            <div className="rc-results-skeleton-row"></div>
                                            <div className="rc-results-skeleton-row"></div>
                                        </div>
                                    </div>
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
                                        <div className="st-col best">SCORE GENERAL</div>
                                        <div className="st-col challenges">RETOS RESUELTOS</div>
                                        <div className="st-col last">ÚLTIMO ENVÍO</div>
                                        <div className="st-col action"></div>
                                    </div>
                                    {Object.entries(studentMap).map(([uid, subs]) => {
                                        const isStudentExpanded = expandedStudentId === uid;
                                        // Per-challenge best scores for this student
                                        const challengeScores = {};
                                        subs.forEach(s => {
                                            const cid = s.challengeId || s.challenge_id || s.ChallengeID;
                                            const sc = s.score || s.Score || 0;
                                            if (!challengeScores[cid] || sc > challengeScores[cid]) {
                                                challengeScores[cid] = sc;
                                            }
                                        });

                                        // Weighted score for this student
                                        let studentTotalMax = 0;
                                        let studentTotalEarned = 0;
                                        const solvedCount = Object.values(challengeScores).filter(sc => sc === 100).length;

                                        details.items.forEach(item => {
                                            const cid = item.challenge?.id || item.challenge?.ID || item.challengeID || item.challenge_id;
                                            const pts = item.points || item.Points || 0;
                                            studentTotalMax += pts;
                                            const best = challengeScores[cid] || 0;
                                            studentTotalEarned += Math.round((best / 100) * pts);
                                        });

                                        const studentScorePct = studentTotalMax > 0 ? Math.round((studentTotalEarned / studentTotalMax) * 100) : 0;

                                        const lastSub = subs.sort((a, b) =>
                                            new Date(b.created_at || b.CreatedAt || 0) - new Date(a.created_at || a.CreatedAt || 0)
                                        )[0];

                                        return (
                                            <div key={uid} className="student-row-wrapper">
                                                <div
                                                    className={`student-row ${isStudentExpanded ? 'expanded' : ''}`}
                                                    onClick={() => {
                                                        setExpandedStudentId(isStudentExpanded ? null : uid);
                                                        setExpandedAttemptId(null);
                                                    }}
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
                                                        <span className={`score-pill ${studentScorePct === 100 ? 'perfect' : studentScorePct >= 50 ? 'partial' : 'low'}`}>
                                                            {studentTotalEarned}/{studentTotalMax} pts ({studentScorePct}%)
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
                                                        {/* Per-challenge breakdown for this student */}
                                                        {details.items.map((item, idx) => {
                                                            const ch = item.challenge || {};
                                                            const cid = ch.id || ch.ID || item.challengeID || item.challenge_id;
                                                            const title = ch.title || ch.Title || `Reto #${idx + 1}`;
                                                            const pts = item.points || item.Points || 0;
                                                            const chSubs = subs
                                                                .filter(s => (s.challengeId || s.challenge_id || s.ChallengeID) === cid)
                                                                .sort((a, b) => new Date(b.created_at || b.CreatedAt || 0) - new Date(a.created_at || a.CreatedAt || 0));
                                                            const bestSc = chSubs.length > 0 ? Math.max(...chSubs.map(s => s.score || s.Score || 0)) : 0;
                                                            const earned = Math.round((bestSc / 100) * pts);
                                                            const isAttemptExpanded = expandedAttemptId === `${uid}-${cid}`;

                                                            return (
                                                                <div key={cid || idx}>
                                                                    <div
                                                                        className="sub-detail-row"
                                                                        style={{ cursor: 'pointer' }}
                                                                        onClick={(e) => {
                                                                            e.stopPropagation();
                                                                            setExpandedAttemptId(isAttemptExpanded ? null : `${uid}-${cid}`);
                                                                        }}
                                                                    >
                                                                        <div className="sub-detail-challenge">
                                                                            <Code size={13} />
                                                                            <span>{title}</span>
                                                                        </div>
                                                                        <div className="sub-detail-status">
                                                                            <span className={`status-pill-mini ${bestSc === 100 ? 'accepted' : bestSc > 0 ? 'rejected' : 'pending'}`}>
                                                                                {bestSc === 100 ? <><CheckCircle size={13} /> Aceptado</> : bestSc > 0 ? <><XCircle size={13} /> Parcial</> : <><Clock size={13} /> Sin intentos</>}
                                                                            </span>
                                                                        </div>
                                                                        <div className="sub-detail-score">
                                                                            <Trophy size={12} />
                                                                            <span>{earned}/{pts} pts</span>
                                                                        </div>
                                                                        <div className="sub-detail-lang">
                                                                            <span>{chSubs.length} intento{chSubs.length !== 1 ? 's' : ''}</span>
                                                                        </div>
                                                                        <div className="sub-detail-date">
                                                                            <div className={`expand-chevron-small ${isAttemptExpanded ? 'open' : ''}`}>
                                                                                <ChevronDown size={12} />
                                                                            </div>
                                                                        </div>
                                                                    </div>
                                                                    {isAttemptExpanded && chSubs.length > 0 && (
                                                                        <div style={{ padding: '0.25rem 0.5rem 0.5rem 2rem' }}>
                                                                            {chSubs.map((sub, si) => {
                                                                                const sc = sub.score || sub.Score || 0;
                                                                                return (
                                                                                    <div key={sub.id || sub.ID || si} style={{
                                                                                        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                                                                                        padding: '0.35rem 0.6rem', borderRadius: '8px', marginBottom: '2px',
                                                                                        background: sc === 100 ? '#f0fdf4' : '#f8fafc', fontSize: '0.73rem', fontWeight: 700,
                                                                                        cursor: 'pointer', border: '1px solid transparent'
                                                                                    }}
                                                                                    onClick={() => setSelectedAttempt({ ...sub, attemptNumber: si + 1 })}
                                                                                    onMouseEnter={(e) => e.currentTarget.style.borderColor = '#cbd5e1'}
                                                                                    onMouseLeave={(e) => e.currentTarget.style.borderColor = 'transparent'}
                                                                                    >
                                                                                        <span style={{ display: 'flex', alignItems: 'center', gap: '5px', color: '#475569' }}>
                                                                                            {sc === 100 ? <CheckCircle size={11} style={{ color: '#16a34a' }} /> : <XCircle size={11} style={{ color: '#ef4444' }} />}
                                                                                            Intento #{si + 1} — {sub.language || sub.Language || '—'}
                                                                                        </span>
                                                                                        <span style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                                                                            <span style={{
                                                                                                padding: '1px 7px', borderRadius: '100px', fontSize: '0.65rem', fontWeight: 900,
                                                                                                background: sc === 100 ? '#dcfce7' : sc >= 50 ? '#fef3c7' : '#fee2e2',
                                                                                                color: sc === 100 ? '#16a34a' : sc >= 50 ? '#d97706' : '#ef4444'
                                                                                            }}>{sc}%</span>
                                                                                            <span style={{ color: '#94a3b8', fontSize: '0.65rem' }}>
                                                                                                {formatDate(sub.created_at || sub.CreatedAt)} {formatTime(sub.created_at || sub.CreatedAt)}
                                                                                            </span>
                                                                                        </span>
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
            <PageLoader message="Cargando envíos..." minHeight="240px" />
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
                        : 'Consulta tu progreso y resultados en cada examen'}</p>
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
                        : 'Los exámenes públicos que estén disponibles aparecerán aquí.'}</p>
                </div>
            ) : (
                <div className="exams-list">
                    {isProfessor ? renderProfessorView() : renderStudentView()}
                </div>
            )}

            {/* ATTEMPT MODAL (CODE & RESULTS) */}
            {selectedAttempt && (
                <div style={{
                    position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
                    background: 'rgba(0,0,0,0.6)', backdropFilter: 'blur(4px)',
                    display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000,
                    padding: '2rem'
                }}>
                    <div style={{
                        background: '#ffffff', borderRadius: '12px', width: '800px', maxWidth: '100%',
                        maxHeight: '90vh', display: 'flex', flexDirection: 'column',
                        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.4)', overflow: 'hidden'
                    }}>
                        <div style={{ padding: '1.25rem 1.5rem', borderBottom: '1px solid #e5e7eb', display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: '#f8fafc' }}>
                            <div>
                                <h3 style={{ margin: 0, fontSize: '1.1rem', fontWeight: 700, color: '#0f172a', display: 'flex', alignItems: 'center', gap: '8px' }}>
                                    <Code size={18} style={{ color: '#c8102e' }} /> Detalles del Envío
                                </h3>
                                <div style={{ fontSize: '0.8rem', color: '#64748b', marginTop: '4px', display: 'flex', gap: '1rem' }}>
                                    <span>Intento #{selectedAttempt.attemptNumber}</span>
                                    <span>•</span>
                                    <span>{selectedAttempt.language || selectedAttempt.Language || '—'}</span>
                                    <span>•</span>
                                    <span style={{ color: (selectedAttempt.score || selectedAttempt.Score) === 100 ? '#16a34a' : '#ef4444', fontWeight: 'bold' }}>
                                        {selectedAttempt.score || selectedAttempt.Score}%
                                    </span>
                                </div>
                            </div>
                            <button onClick={() => setSelectedAttempt(null)} style={{ background: 'transparent', border: 'none', cursor: 'pointer', color: '#94a3b8' }}>
                                <XCircle size={24} />
                            </button>
                        </div>
                        <div style={{ display: 'flex', flex: 1, minHeight: 0 }}>
                            <div style={{ flex: 1, display: 'flex', flexDirection: 'column', borderRight: '1px solid #e5e7eb' }}>
                                <div style={{ padding: '0.5rem 1rem', background: '#f1f5f9', borderBottom: '1px solid #e5e7eb', fontSize: '0.8rem', fontWeight: 600, color: '#475569' }}>
                                    CÓDIGO ENVIADO
                                </div>
                                <div style={{ flex: 1, overflowY: 'auto', padding: '1rem', background: '#1e1e1e' }}>
                                    <pre style={{ margin: 0, color: '#d4d4d4', fontFamily: 'monospace', fontSize: '0.85rem', whiteSpace: 'pre-wrap' }}>
                                        {selectedAttempt.code || selectedAttempt.Code || 'Sin código'}
                                    </pre>
                                </div>
                            </div>
                            <div style={{ width: '350px', display: 'flex', flexDirection: 'column', background: '#f8fafc' }}>
                                <div style={{ padding: '0.5rem 1rem', background: '#f1f5f9', borderBottom: '1px solid #e5e7eb', fontSize: '0.8rem', fontWeight: 600, color: '#475569' }}>
                                    CASOS DE PRUEBA
                                </div>
                                <div style={{ flex: 1, overflowY: 'auto', padding: '1rem' }}>
                                    {(!selectedAttempt.results || selectedAttempt.results.length === 0) ? (
                                        <div style={{ color: '#64748b', fontSize: '0.85rem', textAlign: 'center', marginTop: '2rem' }}>
                                            No hay detalles de resultados.
                                        </div>
                                    ) : (
                                        selectedAttempt.results.map((r, i) => {
                                            const st = (r.Status || r.status || 'unknown').toLowerCase();
                                            const isAcc = st === 'accepted';
                                            return (
                                                <div key={i} style={{
                                                    marginBottom: '0.75rem', padding: '0.75rem', borderRadius: '8px',
                                                    background: '#fff', border: `1px solid ${isAcc ? '#bbf7d0' : '#fecaca'}`,
                                                    boxShadow: '0 1px 2px rgba(0,0,0,0.05)'
                                                }}>
                                                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
                                                        <span style={{ fontSize: '0.8rem', fontWeight: 700, color: '#334155' }}>Caso #{i + 1}</span>
                                                        <span style={{
                                                            fontSize: '0.7rem', fontWeight: 800, padding: '2px 8px', borderRadius: '100px',
                                                            background: isAcc ? '#dcfce7' : '#fee2e2', color: isAcc ? '#16a34a' : '#ef4444'
                                                        }}>
                                                            {isAcc ? 'ACEPTADO' : 'FALLIDO'}
                                                        </span>
                                                    </div>
                                                    {!isAcc && (r.error_message || r.ErrorMessage || r.errorMessage) && (
                                                        <div style={{ fontSize: '0.75rem', color: '#ef4444', background: '#fef2f2', padding: '0.5rem', borderRadius: '4px', marginTop: '0.5rem', whiteSpace: 'pre-wrap' }}>
                                                            {r.error_message || r.ErrorMessage || r.errorMessage}
                                                        </div>
                                                    )}
                                                </div>
                                            );
                                        })
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Submissions;
