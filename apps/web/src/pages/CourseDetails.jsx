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
import { Eye, EyeOff, Lock, Trash2, Calendar, Clock, Trophy, ChevronRight, Edit, ArrowRight, Users, PlusCircle, BookOpen, Sparkles, Target } from 'lucide-react';
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

    const fetchData = async () => {
        setLoading(true);
        try {
            const results = await Promise.allSettled([
                client.get(`/courses/${id}`),
                getCourseExams(id)
            ]);

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
                <div className="challenges-grid-compact">
                    {exams.map(exam => {
                        const examId = exam.id || exam.ID;
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

            <div className="info-footer-compact">
                <Sparkles size={16} className="icon-sparkle" />
                <p>{isProfessor ? 'Gestiona tus exámenes y retos desde esta vista de curso.' : 'Resuelve los exámenes asignados por tu profesor.'}</p>
            </div>
        </div>
    );
};

export default CourseDetails;
