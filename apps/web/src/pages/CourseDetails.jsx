import { useState, useEffect, useContext } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { 
    getCourseExams, 
    toggleExamVisibility, 
    closeExam, 
    deleteExam
} from '../api/exams';
import { AuthContext } from '../context/AuthContext';
import { Eye, EyeOff, Lock, Trash2, Calendar, Clock, Trophy, ChevronRight, Edit, ArrowRight, Users, PlusCircle, BookOpen, Sparkles, Target, BarChart3 } from 'lucide-react';
import Swal from 'sweetalert2';
import './Challenges.css';

const CourseDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);
    const [course, setCourse] = useState(null);
    const [exams, setExams] = useState([]);
    const [loading, setLoading] = useState(true);
    const [processingId, setProcessingId] = useState(null);
    const [expandedResultsExamId, setExpandedResultsExamId] = useState(null);
    const [examResultsMap, setExamResultsMap] = useState({});
    const currentUserId = String(user?.id || user?.ID || '');
    const [studentNameById, setStudentNameById] = useState({});

    const fetchData = async () => {
        setLoading(true);
        try {
            const isProfessorRole = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';
            const requests = [
                client.get(`/courses/${id}`),
                getCourseExams(id)
            ];
            if (isProfessorRole) {
                requests.push(client.get(`/courses/${id}/students`));
            }

            const results = await Promise.allSettled(requests);

            if (results[0].status === 'fulfilled') {
                setCourse(results[0].value.data);
            } else {
                console.error("Course error:", results[0].reason);
                Swal.fire({ icon: 'error', title: 'Error cargando curso', text: results[0].reason.response?.data?.error || 'No se pudo obtener la información del curso.' });
            }

            if (results[1].status === 'fulfilled') {
                setExams(results[1].value);
            } else {
                console.error("Exams error:", results[1].reason);
                setExams([]);
            }

            if (isProfessorRole && results[2]?.status === 'fulfilled') {
                const students = Array.isArray(results[2].value.data)
                    ? results[2].value.data
                    : (results[2].value.data?.students || []);
                const mappedStudents = {};
                students.forEach(student => {
                    const sid = String(student.id || student.ID || student.user_id || student.UserID || '');
                    if (!sid) return;
                    mappedStudents[sid] = student.username || student.name || student.email || sid;
                });
                setStudentNameById(mappedStudents);
            }
        } catch (err) {
            console.error("Fetch data error:", err);
            Swal.fire({ icon: 'error', title: 'Error de conexión', text: 'No se pudo comunicar con el servidor.' });
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [id, user?.role]);

    const handleToggleVisibility = async (examId) => {
        setProcessingId(examId);
        try {
            await toggleExamVisibility(examId);
            Swal.fire({ icon: 'success', title: 'Visibilidad Actualizada', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
            await fetchData();
        } catch (err) {
            console.error(err);
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo cambiar la visibilidad.' });
        } finally {
            setProcessingId(null);
        }
    };

    const handleCloseExam = async (examId) => {
        const { isConfirmed } = await Swal.fire({
            title: '¿Cerrar Examen?',
            text: 'Los estudiantes ya no podrán realizar más envíos.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Sí, cerrar',
            cancelButtonText: 'Cancelar'
        });

        if (!isConfirmed) return;

        setProcessingId(examId);
        try {
            await closeExam(examId);
            Swal.fire({ icon: 'success', title: 'Examen Cerrado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
            await fetchData();
        } catch (err) {
            console.error(err);
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo cerrar el examen.' });
        } finally {
            setProcessingId(null);
        }
    };

    const handleDeleteExam = async (examId) => {
        const { isConfirmed } = await Swal.fire({
            title: '¿Eliminar Examen?',
            text: 'Toda la información del examen se perderá. Esta acción es irreversible.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            confirmButtonText: 'Sí, eliminar',
            cancelButtonText: 'Cancelar'
        });

        if (!isConfirmed) return;

        setProcessingId(examId);
        try {
            await deleteExam(examId);
            Swal.fire({ icon: 'success', title: 'Examen Eliminado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
            await fetchData();
        } catch (err) {
            console.error(err);
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo eliminar el examen.' });
        } finally {
            setProcessingId(null);
        }
    };

    const loadExamResults = async (examId) => {
        if (examResultsMap[examId]?.loaded || examResultsMap[examId]?.loading) return;
        setExamResultsMap(prev => ({ ...prev, [examId]: { ...prev[examId], loading: true } }));
        try {
            const itemsRes = await client.get(`/exams/${examId}/items`);
            const items = Array.isArray(itemsRes.data) ? itemsRes.data : (itemsRes.data?.items || []);
            const allSubmissions = [];

            for (const item of items) {
                const challengeId = item.challenge?.id || item.challenge?.ID || item.challenge_id || item.challengeID;
                if (!challengeId) continue;
                try {
                    const subRes = await client.get(`/submissions/challenge/${challengeId}`);
                    const subs = Array.isArray(subRes.data) ? subRes.data : (subRes.data?.items || []);
                    subs.forEach(sub => {
                        const submission = sub?.Submission || sub?.submission || sub;
                        const score = submission?.score || submission?.Score || 0;
                        const userId = submission?.user_id || submission?.UserID || submission?.userId || 'desconocido';
                        allSubmissions.push({ userId, challengeId, score, createdAt: submission?.created_at || submission?.CreatedAt || null });
                    });
                } catch (subErr) {
                    console.warn(`No submissions for challenge ${challengeId}:`, subErr);
                }
            }

            const byStudent = {};
            allSubmissions.forEach(s => {
                if (!byStudent[s.userId]) byStudent[s.userId] = { submissions: [], bestByChallenge: {} };
                byStudent[s.userId].submissions.push(s);
                const prevBest = byStudent[s.userId].bestByChallenge[s.challengeId] || 0;
                byStudent[s.userId].bestByChallenge[s.challengeId] = Math.max(prevBest, s.score);
            });

            const studentRows = Object.entries(byStudent).map(([userId, data]) => {
                let totalMax = 0;
                let totalEarned = 0;
                let solved = 0;
                items.forEach(item => {
                    const challengeId = item.challenge?.id || item.challenge?.ID || item.challenge_id || item.challengeID;
                    const points = item.points || item.Points || 0;
                    const bestScore = data.bestByChallenge[challengeId] || 0;
                    totalMax += points;
                    totalEarned += Math.round((bestScore / 100) * points);
                    if (bestScore === 100) solved += 1;
                });
                const latest = data.submissions
                    .map(s => s.createdAt)
                    .filter(Boolean)
                    .sort((a, b) => new Date(b) - new Date(a))[0];
                return {
                    userId,
                    submissionsCount: data.submissions.length,
                    solved,
                    totalChallenges: items.length,
                    totalEarned,
                    totalMax,
                    scorePct: totalMax > 0 ? Math.round((totalEarned / totalMax) * 100) : 0,
                    latest
                };
            }).sort((a, b) => b.scorePct - a.scorePct);

            setExamResultsMap(prev => ({
                ...prev,
                [examId]: { loading: false, loaded: true, studentRows }
            }));
        } catch (err) {
            console.error('Error loading exam results:', err);
            setExamResultsMap(prev => ({
                ...prev,
                [examId]: { loading: false, loaded: false, error: 'No se pudieron cargar los resultados.' }
            }));
        }
    };

    const toggleExamResults = async (examId) => {
        const nextId = expandedResultsExamId === examId ? null : examId;
        setExpandedResultsExamId(nextId);
        if (nextId) await loadExamResults(examId);
    };

    if (loading) return (
        <div className="challenges-page">
            <header className="page-header-compact">
                <div className="skeleton title-skeleton"></div>
            </header>
            <div className="challenges-grid-compact">
                {[...Array(3)].map((_, i) => (
                    <div key={i} className="challenge-card-mini skeleton-card">
                        <div className="skeleton card-content-skeleton"></div>
                    </div>
                ))}
            </div>
        </div>
    );

    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    return (
        <div className="challenges-page">
            <header className="page-header-compact">
                <div className="header-info">
                    <div className="breadcrumb-mini">
                        <BookOpen size={14} />
                        <span>Cursos</span>
                        <ChevronRight size={12} />
                        <span>{course?.name || 'Curso'}</span>
                    </div>
                    <h1>{course?.name}</h1>
                    <p>
                        {course?.code} — {course?.period ? `${course.period.year}-${course.period.semester}` : 'S/P'} — Grupo {course?.groupNumber}
                    </p>
                </div>

                <div className="header-actions-mini">
                    {isProfessor && (
                        <>
                            <button className="btn-create-mini" onClick={() => navigate(`/courses/${id}/students`)}>
                                <Users size={18} />
                                <span>Estudiantes</span>
                            </button>
                            <button className="btn-create-mini" onClick={() => navigate(`/exams/create?courseId=${id}`)}>
                                <PlusCircle size={18} />
                                <span>Nuevo Examen</span>
                            </button>
                        </>
                    )}
                </div>
            </header>

            {exams.length === 0 ? (
                <div className="empty-state-mini">
                    <div className="icon-circle">
                        <Trophy size={32} />
                    </div>
                    <h3>No hay exámenes en este curso</h3>
                    <p>{isProfessor ? 'Crea un examen para que tus estudiantes lo resuelvan.' : 'Los exámenes aparecerán aquí cuando sean creados.'}</p>
                </div>
            ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                    {exams.map(exam => {
                        const examId = exam.id || exam.ID;
                        const examOwnerId = String(
                            exam.user_id ||
                            exam.UserID ||
                            exam.userId ||
                            exam.created_by ||
                            exam.createdBy ||
                            exam.CreatedBy ||
                            exam.owner_id ||
                            exam.ownerId ||
                            exam.OwnerID ||
                            exam.professor_id ||
                            exam.professorId ||
                            exam.ProfessorID ||
                            exam.teacher_id ||
                            exam.teacherId ||
                            exam.TeacherID ||
                            ''
                        );
                        const canEditExam = isProfessor && (!examOwnerId || examOwnerId === currentUserId);
                        const title = exam.title || exam.Title;
                        const desc = exam.description || exam.Description || 'Sin descripción disponible.';
                        const visibility = String(exam.visibility || exam.Visibility || 'private').toLowerCase();
                        const timeLimit = exam.timeLimit || exam.TimeLimit || 3600;
                        const startTime = exam.start_time || exam.startTime || exam.StartTime;
                        const endTime = exam.end_time || exam.endTime || exam.EndTime;
                        const isClosed = Boolean(endTime && new Date(endTime) <= new Date());
                        const isVisible = visibility !== 'private';
                        const isStudentVisible = visibility === 'public' || visibility === 'course';
                        const tryLimit = exam.try_limit ?? exam.tryLimit ?? exam.TryLimit ?? 1;

                        const limitText = tryLimit === -1 ? 'Ilimitados' : tryLimit;

                        const formattedAvailability = (!startTime && !endTime) ? 'Siempre' :
                            (startTime && endTime) ? `${new Date(startTime).toLocaleDateString()} al ${new Date(endTime).toLocaleDateString()}` :
                            (startTime ? `Desde ${new Date(startTime).toLocaleDateString()}` : `Hasta ${new Date(endTime).toLocaleDateString()}`);

                        return (
                            <div
                                key={examId}
                                style={{
                                    border: '1px solid #e5e7eb',
                                    borderRadius: '12px',
                                    background: isClosed ? '#f8fafc' : '#ffffff',
                                    opacity: isClosed ? 0.9 : 1,
                                    padding: '0.9rem 1rem',
                                    display: 'flex',
                                    flexDirection: 'column',
                                    gap: '0.65rem',
                                    cursor: 'default'
                                }}
                            >
                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: '0.75rem' }}>
                                    <div style={{ minWidth: 0 }}>
                                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.45rem' }}>
                                            <Trophy size={16} style={{ color: '#c8102e', flexShrink: 0 }} />
                                            <h3 style={{ margin: 0, fontSize: '1rem', color: '#111827', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{title}</h3>
                                        </div>
                                        <p style={{ margin: '0.35rem 0 0', fontSize: '0.85rem', color: '#6b7280', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                                            {desc}
                                        </p>
                                    </div>
                                    <div className="badge-group" style={{ flexShrink: 0 }}>
                                        {isClosed && <span className="status-badge closed">Cerrado</span>}
                                        {!isVisible && <span className="status-badge private">Privado</span>}
                                        {isVisible && <span className="status-badge public">{visibility === 'course' ? 'Curso' : 'Público'}</span>}
                                    </div>
                                </div>

                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: '0.75rem', flexWrap: 'wrap' }}>
                                    <div className="stats-mini">
                                        <div className="stat">
                                            <Clock size={14} />
                                            <span>{Math.floor(timeLimit / 60)} min</span>
                                        </div>
                                        <div className="stat">
                                            <Calendar size={14} />
                                            <span>{formattedAvailability}</span>
                                        </div>
                                        {!isProfessor && (
                                            <div className="stat" style={{ color: '#4b5563' }}>
                                                <Target size={14} />
                                                <span>Límite: {limitText}</span>
                                            </div>
                                        )}
                                    </div>

                                    <div className="actions-wrapper">
                                            {isProfessor ? (
                                                <div className="exam-admin-group">
                                                    {canEditExam && (
                                                        <button onClick={(e) => { e.stopPropagation(); navigate(`/exam/${examId}/edit`); }} className="btn-action-mini primary" title="Editar">
                                                            <Edit size={14} /> Editar
                                                        </button>
                                                    )}
                                                    <button onClick={(e) => { e.stopPropagation(); toggleExamResults(examId); }} className="btn-action-mini" title="Resultados">
                                                        <BarChart3 size={14} />
                                                    </button>
                                                    <button onClick={(e) => { e.stopPropagation(); handleToggleVisibility(examId); }} className="btn-action-mini" disabled={processingId === examId} title="Visibilidad">
                                                        {isVisible ? <Eye size={14} /> : <EyeOff size={14} />}
                                                    </button>
                                                    {!isClosed && (
                                                        <button onClick={(e) => { e.stopPropagation(); handleCloseExam(examId); }} className="btn-action-mini" disabled={processingId === examId} title="Cerrar Examen">
                                                            <Lock size={14} />
                                                        </button>
                                                    )}
                                                    <button onClick={(e) => { e.stopPropagation(); handleDeleteExam(examId); }} className="btn-action-mini delete" disabled={processingId === examId} title="Eliminar Examen">
                                                        <Trash2 size={14} />
                                                    </button>
                                                </div>
                                            ) : (
                                                !isClosed && isStudentVisible && (
                                                    <Link to={`/exam/${examId}`} className="btn-action-mini primary">
                                                        Iniciar Examen <ArrowRight size={16} />
                                                    </Link>
                                                )
                                            )}
                                    </div>
                                </div>

                                {isProfessor && expandedResultsExamId === examId && (
                                    <div style={{ marginTop: '0.5rem', borderTop: '1px solid #e5e7eb', paddingTop: '0.75rem' }}>
                                        {examResultsMap[examId]?.loading && (
                                            <p style={{ margin: 0, fontSize: '0.85rem', color: '#6b7280' }}>Cargando resultados...</p>
                                        )}
                                        {examResultsMap[examId]?.error && (
                                            <p style={{ margin: 0, fontSize: '0.85rem', color: '#b91c1c' }}>{examResultsMap[examId].error}</p>
                                        )}
                                        {examResultsMap[examId]?.loaded && examResultsMap[examId].studentRows?.length === 0 && (
                                            <p style={{ margin: 0, fontSize: '0.85rem', color: '#6b7280' }}>Aún no hay envíos para este examen.</p>
                                        )}
                                        {examResultsMap[examId]?.loaded && examResultsMap[examId].studentRows?.length > 0 && (
                                            <div style={{ display: 'grid', gap: '0.4rem' }}>
                                                {examResultsMap[examId].studentRows.map(row => (
                                                    <div key={row.userId} style={{ display: 'grid', gridTemplateColumns: '1.3fr .6fr .8fr .8fr .9fr', gap: '0.6rem', alignItems: 'center', fontSize: '0.8rem', background: '#f9fafb', border: '1px solid #e5e7eb', borderRadius: '8px', padding: '0.5rem 0.65rem' }}>
                                                        <span style={{ fontWeight: 700, color: '#1f2937', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }} title={studentNameById[row.userId] || row.userId}>
                                                            {studentNameById[row.userId] || row.userId}
                                                        </span>
                                                        <span style={{ color: '#4b5563' }}>{row.submissionsCount} env.</span>
                                                        <span style={{ color: '#4b5563' }}>{row.solved}/{row.totalChallenges}</span>
                                                        <span style={{ fontWeight: 700, color: row.scorePct >= 70 ? '#15803d' : row.scorePct >= 40 ? '#b45309' : '#b91c1c' }}>
                                                            {row.totalEarned}/{row.totalMax} ({row.scorePct}%)
                                                        </span>
                                                        <span style={{ color: '#6b7280' }}>{row.latest ? new Date(row.latest).toLocaleDateString() : '—'}</span>
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                )}
                            </div>
                        );
                    })}
                </div>
            )}

            <div className="info-footer-compact">
                <Sparkles size={16} className="icon-sparkle" />
                <p>{isProfessor ? 'Gestiona tus exámenes y retos desde esta vista de curso.' : 'Resuelve los exámenes asignados por tu profesor.'}</p>
            </div>
        </div>
    );
};

export default CourseDetails;
