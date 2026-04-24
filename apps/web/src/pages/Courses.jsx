import { useState, useEffect, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import Swal from 'sweetalert2';
import { 
    BookOpen, 
    Plus, 
    Compass, 
    Hash, 
    Calendar, 
    Settings, 
    ChevronRight,
    Users,
    Key,
    AlertCircle,
    RotateCcw,
    Trash2
} from 'lucide-react';
import PageLoader from '../components/PageLoader';
import './Courses.css';

const Courses = () => {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const { user } = useContext(AuthContext);
    const navigate = useNavigate();

    const fetchCourses = async () => {
        try {
            const scope = (user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin') ? '?scope=owned' : '?scope=enrolled';
            const { data } = await client.get(`/courses${scope}`);
            setCourses(Array.isArray(data) ? data : (data?.items || []));
        } catch (err) {
            console.error('Error loading courses:', err);
            setError('No se pudieron cargar los cursos. Por favor, intenta de nuevo.');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user) fetchCourses();
    }, [user]);

    const handleDeleteCourse = async (e, id) => {
        e.preventDefault();
        e.stopPropagation();
        
        const { isConfirmed } = await Swal.fire({
            title: '¿Eliminar curso?',
            text: 'Esta acción borrará todos los datos asociados. Es irreversible.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            confirmButtonText: 'Sí, eliminar curso',
            cancelButtonText: 'Cancelar'
        });

        if (!isConfirmed) return;
        
        try {
            await client.delete(`/courses/${id}`);
            Swal.fire({
                icon: 'success',
                title: 'Curso Eliminado',
                timer: 1000,
                toast: true,
                position: 'top-end',
                showConfirmButton: false
            });
            await fetchCourses();
        } catch (err) {
            console.error('Error deleting course:', err);
            Swal.fire({ 
                icon: 'error', 
                title: 'Error', 
                text: 'No se pudo eliminar el curso. Verifica si tiene estudiantes inscritos.' 
            });
        }
    };

    if (loading) return (
        <div className="courses-page-new">
            <PageLoader message="Cargando cursos..." minHeight="240px" />
            <div className="courses-grid-new">
                {[...Array(3)].map((_, i) => (
                    <div key={i} className="course-card-new skeleton-card">
                        <div className="skeleton card-content-skeleton"></div>
                    </div>
                ))}
            </div>
        </div>
    );

    if (error) return (
        <div className="courses-page-new">
            <div className="error-container">
                <AlertCircle size={48} />
                <h3>Error en la carga</h3>
                <p>{error}</p>
                <button onClick={() => window.location.reload()} className="btn-retry">
                    <RotateCcw size={16} /> Reintentar carga
                </button>
            </div>
        </div>
    );

    return (
        <div className="courses-page-new">
            <header className="courses-header-new">
                <div className="header-info-new">
                    <h1>Mis Cursos</h1>
                    <p>Gestiona tus asignaturas y accede a los retos programados</p>
                </div>
                
                <div className="header-actions-new">
                    {(user?.role === 'student') && (
                        <>
                            <button onClick={() => navigate('/courses/browse')} className="btn-action-outline">
                                <Compass size={18} /> Explorar Cursos
                            </button>
                            <button onClick={() => navigate('/courses/join')} className="btn-action-filled">
                                <Key size={18} /> Unirse con Código
                            </button>
                        </>
                    )}
                    {(user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin') && (
                        <button onClick={() => navigate('/courses/create')} className="btn-action-filled">
                            <Plus size={18} /> Crear Nuevo Curso
                        </button>
                    )}
                </div>
            </header>

            {courses.length === 0 ? (
                <div className="empty-state-new">
                    <div className="icon-circle">
                        <BookOpen size={40} />
                    </div>
                    <h3>Aún no tienes cursos</h3>
                    <p>
                        {user?.role === 'student' 
                            ? 'No te has unido a ningún curso. Revisa la sección de explorar o solicita un código a tu docente.' 
                            : 'No tienes cursos asignados todavía. Contacta al administrador para que te asigne uno.'}
                    </p>
                    {user?.role === 'student' && (
                        <button onClick={() => navigate('/courses/browse')} className="btn-cta-link">
                            Empezar ahora <ChevronRight size={16} />
                        </button>
                    )}
                </div>
            ) : (
                <div className="courses-grid-new">
                    {courses.map((course) => {
                        const accentColor = course.visual_identity || '#c8102e';
                        return (
                        <div key={course.id} className="course-card-new">
                            <Link to={`/courses/${course.id}`} className="card-clickable-area">
                                <div className="card-icon-area" style={{ backgroundColor: accentColor + '26', color: accentColor }}>
                                    <BookOpen className="course-icon" />
                                </div>
                                <div className="card-info-area">
                                    <h3>{course.name || 'Sin nombre'}</h3>
                                    <div className="course-metadata-new">
                                        <div className="meta-item-new">
                                            <Hash size={14} />
                                            <span>{course.code || 'S/N'}</span>
                                        </div>
                                        <div className="meta-item-new">
                                            <Calendar size={14} />
                                            <span>{course.period ? `${course.period.year}-${course.period.semester}` : 'S/P'}</span>
                                        </div>
                                        <div className="meta-item-new">
                                            <Users size={14} />
                                            <span>{course.studentCount || 0} inscritos</span>
                                        </div>
                                    </div>
                                </div>
                            </Link>
                            
                            {(user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin') && (
                                <div className="card-admin-actions">
                                    <button 
                                        onClick={(e) => {
                                            e.preventDefault();
                                            e.stopPropagation();
                                            navigate(`/courses/edit/${course.id}`);
                                        }}
                                        className="btn-settings"
                                        title="Configurar curso"
                                    >
                                        <Settings size={18} />
                                    </button>
                                    <button 
                                        onClick={(e) => handleDeleteCourse(e, course.id)}
                                        className="btn-delete-course"
                                        title="Eliminar curso"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </div>
                            )}
                            <Link to={`/courses/${course.id}`} className="card-arrow-link">
                                <ChevronRight size={20} />
                            </Link>
                        </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};

export default Courses;
