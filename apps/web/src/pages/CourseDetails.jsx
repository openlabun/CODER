import { useState, useEffect, useContext } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { 
    getCourseExams, 
    toggleExamVisibility, 
    closeExam, 
    deleteExam,
    createExamSession
} from '../api/exams';
import { AuthContext } from '../context/AuthContext';
import { Eye, EyeOff, Lock, Trash2, Calendar, Clock, Trophy, Target, ChevronRight, Code, Edit, ArrowRight, Loader2 } from 'lucide-react';
import Swal from 'sweetalert2';
import './Courses.css';
import './CourseActions.css';
import './Challenges.css';

const CourseDetails = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);
    const [course, setCourse] = useState(null);
    const [exams, setExams] = useState([]);
    const [challenges, setChallenges] = useState([]);
    const [loading, setLoading] = useState(true);
    const [processingId, setProcessingId] = useState(null);

    const fetchData = async () => {
        setLoading(true);
        try {
            // Fetch non-blocking results separately to avoid one failing the others
            const results = await Promise.allSettled([
                client.get(`/courses/${id}`),
                getCourseExams(id),
                client.get(`/courses/${id}/challenges`)
            ]);

            // [0] Course Details
            if (results[0].status === 'fulfilled') {
                setCourse(results[0].value.data);
            } else {
                console.error("Course error:", results[0].reason);
                Swal.fire({ icon: 'error', title: 'Error cargando curso', text: results[0].reason.response?.data?.error || 'No se pudo obtener la información del curso.' });
            }

            // [1] Exams
            if (results[1].status === 'fulfilled') {
                setExams(results[1].value);
            } else {
                console.error("Exams error:", results[1].reason);
                setExams([]);
            }

            // [2] Challenges
            if (results[2].status === 'fulfilled') {
                setChallenges(results[2].value.data.challenges || []);
            } else {
                console.error("Challenges error:", results[2].reason);
                setChallenges([]);
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
    }, [id]);

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



    if (loading) return (
        <div className="page-loader">
            <Loader2 className="page-loader-spinner" size={48} />
            <p className="page-loader-text">Cargando curso...</p>
        </div>
    );
    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    return (
        <div className="course-details-page">
            <div className="course-header">
                <div>
                    <h1>{course?.name}</h1>
                    <p className="course-meta">
                        {course?.code} - {course?.period ? `${course.period.year}-${course.period.semester}` : 'S/P'} - Grupo {course?.groupNumber}
                    </p>
                </div>

                {isProfessor && (
                    <div className="course-actions">
                        <button
                            onClick={() => navigate(`/courses/${id}/students`)}
                            className="btn-action btn-students"
                        >
                            <span>👥</span> Ver Estudiantes
                        </button>
                        <button
                            onClick={() => navigate(`/challenges/create?courseId=${id}`)}
                            className="btn-action btn-create-challenge"
                        >
                            <span>➕</span> Crear Reto
                        </button>
                    </div>
                )}
            </div>

            <section className="exams-section-new">
                <div className="section-header">
                    <h2>📚 Exámenes</h2>
                    {isProfessor && (
                        <button className="btn-add-mini" onClick={() => navigate(`/exams/create?courseId=${id}`)}>
                            Nuevo Examen
                        </button>
                    )}
                </div>
                {exams.length === 0 ? (
                    <div className="empty-state-mini-alt">
                        <div className="empty-state-icon">📝</div>
                        <h3>No hay exámenes programados</h3>
                        <p>Los exámenes aparecerán aquí cuando sean creados.</p>
                    </div>
                ) : (
                    <div className="challenges-grid-compact">
                                {exams.map(exam => {
                                    const examId = exam.id || exam.ID;
                                    const title = exam.title || exam.Title;
                                    const desc = exam.description || exam.Description || 'Sin descripción disponible.';
                                    const visibility = String(exam.visibility || exam.Visibility || 'private').toLowerCase();
                                    const endTime = exam.endTime || exam.EndTime;
                                    const isClosed = Boolean(endTime && new Date(endTime) <= new Date());
                                    const isVisible = visibility !== 'private';
                                    const isStudentVisible = visibility === 'public' || visibility === 'course';
                                    const timeLimit = exam.timeLimit || exam.TimeLimit || 3600;
                                    const startTime = exam.startTime || exam.StartTime;

                                    return (
                                        <div key={examId} className={`challenge-card-mini public-exam-card ${isClosed ? 'closed' : ''}`}>
                                            <div className="card-accent public"></div>
                                            <div className="card-main">
                                                <div className="card-top">
                                                    <div className="title-area">
                                                        <Trophy size={16} className="title-icon highlight" />
                                                        <h3>{title}</h3>
                                                    </div>
                                                    <div className="badge-group">
                                                        {isClosed && <span className="status-badge closed">Cerrado</span>}
                                                        {!isVisible && <span className="status-badge private">Privado</span>}
                                                        {isVisible && <span className="status-badge public">{visibility === 'course' ? 'Curso' : 'Público'}</span>}
                                                    </div>
                                                </div>
                                                <p className="description-text">{desc}</p>
                                                
                                                <div className="card-footer-mini">
                                                    <div className="stats-mini">
                                                        <div className="stat">
                                                            <Clock size={14} />
                                                            <span>{Math.floor(timeLimit / 60)} min</span>
                                                        </div>
                                                        <div className="stat">
                                                            <Calendar size={14} />
                                                            <span>{startTime ? new Date(startTime).toLocaleDateString() : 'Siempre'}</span>
                                                        </div>
                                                    </div>
                                                    
                                                    <div className="actions-wrapper">
                                                        {isProfessor ? (
                                                            <div className="exam-admin-group">
                                                                <button onClick={() => navigate(`/exam/${examId}/edit`)} className="btn-action-mini primary" title="Editar">
                                                                    <Edit size={14} /> Editar
                                                                </button>
                                                                <button onClick={() => handleToggleVisibility(examId)} className="btn-action-mini" disabled={processingId === examId} title="Visibilidad">
                                                                    {isVisible ? <Eye size={14} /> : <EyeOff size={14} />}
                                                                </button>
                                                                {!isClosed && (
                                                                    <button onClick={() => handleCloseExam(examId)} className="btn-action-mini" disabled={processingId === examId} title="Cerrar Examen">
                                                                        <Lock size={14} />
                                                                    </button>
                                                                )}
                                                                <button onClick={() => handleDeleteExam(examId)} className="btn-action-mini delete" disabled={processingId === examId} title="Eliminar Examen">
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
                                            </div>
                                        </div>
                                    );
                                })}
                    </div>
                )}
            </section>
        </div>
    );
};

export default CourseDetails;
