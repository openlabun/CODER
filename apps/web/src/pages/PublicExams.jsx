import { useState, useEffect, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { AuthContext } from '../context/AuthContext';
import { 
    ShieldCheck, 
    Search, 
    Filter,
    Trophy, 
    Clock, 
    Calendar,
    ChevronRight,
    Play,
    Edit2,
    AlertTriangle,
    Lock,
    Eye,
    EyeOff,
    Users,
    BookOpen,
    Sparkles,
    Layout,
    PlusCircle,
    Target,
    ArrowRight
} from 'lucide-react';
import PageLoader from '../components/PageLoader';
import './Challenges.css';

const PublicExams = () => {
    const { user } = useContext(AuthContext);
    const navigate = useNavigate();
    const [exams, setExams] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [visibilityFilter, setVisibilityFilter] = useState('all');
    const [loading, setLoading] = useState(true);

    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';
    const currentUserId = String(user?.id || user?.ID || '');

    const fetchPublicExams = async () => {
        setLoading(true);
        try {
            const { data } = await client.get('/exams/public');
            setExams(Array.isArray(data) ? data : (data.items || []));
        } catch (err) {
            console.error('Error loading public exams:', err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPublicExams();
    }, []);

    const getExamVisibility = (exam) => String(exam.visibility || exam.Visibility || 'public').toLowerCase();

    const visibleExams = isProfessor
        ? exams
        : exams.filter((exam) => getExamVisibility(exam) === 'public');

    const filteredExams = visibleExams.filter((e) => {
        const title = (e.title || e.Title || '').toLowerCase();
        const matchesSearch = title.includes(searchTerm.toLowerCase());

        if (!isProfessor) return matchesSearch;

        const visibility = getExamVisibility(e);
        const matchesVisibility = visibilityFilter === 'all' ? true : visibility === visibilityFilter;
        return matchesSearch && matchesVisibility;
    });

    const getVisibilityMeta = (exam) => {
        const visibility = String(exam.visibility || exam.Visibility || 'public').toLowerCase();

        if (visibility === 'private') {
            return { label: 'Privado', className: 'private', icon: <EyeOff size={14} /> };
        }

        if (visibility === 'teachers') {
            return { label: 'Docentes', className: 'teachers', icon: <Users size={14} /> };
        }

        if (visibility === 'course') {
            return { label: 'Curso', className: 'course', icon: <BookOpen size={14} /> };
        }

        return { label: 'Público', className: 'public', icon: <Eye size={14} /> };
    };

    const parseDateSafe = (value) => {
        if (!value) return null;
        const parsed = new Date(value);
        return Number.isNaN(parsed.getTime()) ? null : parsed;
    };

    const getFilterLabel = (filterValue) => {
        const labels = {
            'all': 'Todos',
            'public': 'Para todos',
            'teachers': 'Solo profesores',
            'course': 'Para mis Cursos',
            'private': 'Privados/Borradores'
        };
        return labels[filterValue] || filterValue;
    };

    if (loading) return (
        <div className="challenges-page">
            <PageLoader message="Cargando exámenes..." />
        </div>
    );

    return (
        <div className="challenges-page">
            <header className="page-header-compact">
                <div className="header-info">
                    <div className="breadcrumb-mini">
                        <ShieldCheck size={14} />
                        <span>Comunidad</span>
                        <ChevronRight size={12} />
                        <span>Actividades Públicas</span>
                    </div>
                    <h1>Actividades Públicas</h1>
                    <p>Encuentra retos y evaluaciones creadas por la comunidad RobleCode</p>
                </div>
                
                <div className="header-actions-mini">
                    <div className="search-bar-mini">
                        <Search size={18} />
                        <input 
                            type="text" 
                            placeholder="Buscar actividad..." 
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                    {isProfessor && (
                        <button className="btn-create-mini" onClick={() => navigate('/exams/create')}>
                            <PlusCircle size={18} />
                            <span>Nuevo Examen</span>
                        </button>
                    )}
                </div>
            </header>

            {isProfessor && (
                <section className="visibility-filter-row">
                    <div className="filter-label">
                        <Filter size={14} />
                        <span>Filtrar por visibilidad:</span>
                    </div>
                    <div className="filter-chips">
                        <button
                            type="button"
                            className={`filter-chip ${visibilityFilter === 'all' ? 'active' : ''}`}
                            onClick={() => setVisibilityFilter('all')}
                        >
                            Todos
                        </button>
                        <button
                            type="button"
                            className={`filter-chip ${visibilityFilter === 'public' ? 'active' : ''}`}
                            onClick={() => setVisibilityFilter('public')}
                        >
                            Para todos
                        </button>
                        <button
                            type="button"
                            className={`filter-chip ${visibilityFilter === 'teachers' ? 'active' : ''}`}
                            onClick={() => setVisibilityFilter('teachers')}
                        >
                            Solo profesores
                        </button>
                        <button
                            type="button"
                            className={`filter-chip ${visibilityFilter === 'course' ? 'active' : ''}`}
                            onClick={() => setVisibilityFilter('course')}
                        >
                            Para mis Cursos
                        </button>
                        <button
                            type="button"
                            className={`filter-chip ${visibilityFilter === 'private' ? 'active' : ''}`}
                            onClick={() => setVisibilityFilter('private')}
                        >
                            Privados/Borradores
                        </button>
                    </div>
                </section>
            )}

            {filteredExams.length === 0 ? (
                <div className="empty-state-mini">
                    <div className="icon-circle">
                        <Layout size={32} />
                    </div>
                    <h3>No hay exámenes disponibles</h3>
                    <p>
                        {searchTerm 
                            ? 'Prueba con otros términos de búsqueda.' 
                            : `${isProfessor && visibilityFilter !== 'all' ? `No hay exámenes con visibilidad "${getFilterLabel(visibilityFilter)}".` : 'Vuelve más tarde para ver nuevas evaluaciones.'}`
                        }
                    </p>
                </div>
            ) : (
                <div className="challenges-grid-compact">
                    {filteredExams.map((exam) => {
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
                        const timeLimit = exam.timeLimit || exam.TimeLimit || 3600;
                        const startTime = exam.startTime || exam.StartTime || exam.start_time;
                        const endTime = exam.endTime || exam.EndTime || exam.end_time;
                        const tryLimit = exam.try_limit ?? exam.tryLimit ?? exam.TryLimit ?? 1;
                        const limitText = tryLimit === -1 ? 'Ilimitados' : tryLimit;
                        const formattedAvailability = (!startTime && !endTime) ? 'Siempre' :
                            (startTime && endTime) ? `${new Date(startTime).toLocaleDateString()} al ${new Date(endTime).toLocaleDateString()}` :
                            (startTime ? `Desde ${new Date(startTime).toLocaleDateString()}` : `Hasta ${new Date(endTime).toLocaleDateString()}`);
                        const allowLateValue = exam.allowLateSubmissions ?? exam.AllowLateSubmissions ?? exam.allow_late_submissions;
                        const allowLateSubmissions =
                            allowLateValue === true ||
                            allowLateValue === 1 ||
                            allowLateValue === '1' ||
                            String(allowLateValue).toLowerCase() === 'true';
                        const hasEnded = Boolean(parseDateSafe(endTime) && parseDateSafe(endTime) <= new Date());
                        const studentCanStart = !hasEnded || allowLateSubmissions;
                        const visibilityMeta = getVisibilityMeta(exam);

                        return (
                            <div key={examId} className="challenge-card-mini public-exam-card">
                                <div className="card-accent public"></div>
                                <div className="card-main">
                                    <div className="card-top">
                                        <div className="title-area">
                                            <Trophy size={16} className="title-icon highlight" />
                                            <h3>{title}</h3>
                                        </div>
                                        <div className="meta-badges-row">
                                            <span className={`status-badge ${visibilityMeta.className}`}>
                                                {visibilityMeta.icon}
                                                {visibilityMeta.label}
                                            </span>
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

                                        <div className="meta-badges-row exam-timing-flags">
                                            {hasEnded && !allowLateSubmissions && (
                                                <span className="status-badge closed">
                                                    <AlertTriangle size={14} />
                                                    Examen cerrado
                                                </span>
                                            )}
                                            {hasEnded && allowLateSubmissions && (
                                                <span className="status-badge overtime">
                                                    <AlertTriangle size={14} />
                                                    Hora de cierre superada
                                                </span>
                                            )}
                                            {allowLateSubmissions && (
                                                <span className="status-badge late">
                                                    <Clock size={14} />
                                                    Admite entregas tardías
                                                </span>
                                            )}
                                        </div>
                                        
                                        <div className="actions-wrapper">
                                            {isProfessor ? (
                                                canEditExam ? (
                                                    <Link to={`/exam/${examId}/edit`} className="btn-action-mini primary">
                                                        Editar Actividad <ChevronRight size={16} />
                                                    </Link>
                                                ) : (
                                                    <Link to={`/exam/${examId}`} className="btn-action-mini">
                                                        Abrir Actividad <ArrowRight size={16} />
                                                    </Link>
                                                )
                                            ) : (
                                                studentCanStart ? (
                                                    <Link to={`/exam/${examId}`} className="btn-action-mini primary">
                                                        Iniciar Actividad <ArrowRight size={16} />
                                                    </Link>
                                                ) : (
                                                    <button type="button" className="btn-action-mini" disabled>
                                                        <Lock size={14} />
                                                        <span>Examen Cerrado</span>
                                                    </button>
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
                <p>¿Eres profesor? Crea actividades públicas para que toda la comunidad pueda resolver tus retos.</p>
            </div>
        </div>
    );
};

export default PublicExams;
