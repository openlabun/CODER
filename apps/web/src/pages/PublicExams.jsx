import { useState, useEffect, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { AuthContext } from '../context/AuthContext';
import { 
    ShieldCheck, 
    Search, 
    Trophy, 
    Clock, 
    Calendar,
    ChevronRight,
    ArrowRight,
    Sparkles,
    Layout,
    PlusCircle
} from 'lucide-react';
import './Challenges.css';

const PublicExams = () => {
    const { user } = useContext(AuthContext);
    const navigate = useNavigate();
    const [exams, setExams] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(true);

    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

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

    const filteredExams = exams.filter(e => 
        (e.title || e.Title || '').toLowerCase().includes(searchTerm.toLowerCase())
    );

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

    return (
        <div className="challenges-page">
            <header className="page-header-compact">
                <div className="header-info">
                    <div className="breadcrumb-mini">
                        <ShieldCheck size={14} />
                        <span>Comunidad</span>
                        <ChevronRight size={12} />
                        <span>Exámenes Públicos</span>
                    </div>
                    <h1>Evaluaciones Públicas</h1>
                    <p>Encuentra retos y exámenes creados por la comunidad RobleCode</p>
                </div>
                
                <div className="header-actions-mini">
                    <div className="search-bar-mini">
                        <Search size={18} />
                        <input 
                            type="text" 
                            placeholder="Buscar examen..." 
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

            {filteredExams.length === 0 ? (
                <div className="empty-state-mini">
                    <div className="icon-circle">
                        <Layout size={32} />
                    </div>
                    <h3>No hay exámenes públicos disponibles</h3>
                    <p>{searchTerm ? 'Prueba con otros términos de búsqueda.' : 'Vuelve más tarde para ver nuevas evaluaciones.'}</p>
                </div>
            ) : (
                <div className="challenges-grid-compact">
                    {filteredExams.map((exam) => {
                        const examId = exam.id || exam.ID;
                        const title = exam.title || exam.Title;
                        const desc = exam.description || exam.Description || 'Sin descripción disponible.';
                        const timeLimit = exam.timeLimit || exam.TimeLimit || 3600;
                        const startTime = exam.startTime || exam.StartTime;

                        return (
                            <div key={examId} className="challenge-card-mini public-exam-card">
                                <div className="card-accent public"></div>
                                <div className="card-main">
                                    <div className="card-top">
                                        <div className="title-area">
                                            <Trophy size={16} className="title-icon highlight" />
                                            <h3>{title}</h3>
                                        </div>
                                        <div className="badge-group">
                                            <span className="status-badge public">Público</span>
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
                                                <Link to={`/exam/${examId}/edit`} className="btn-action-mini primary">
                                                    Editar Examen <ChevronRight size={16} />
                                                </Link>
                                            ) : (
                                                <Link to={`/exam/${examId}`} className="btn-action-mini primary">
                                                    Iniciar Examen <ArrowRight size={16} />
                                                </Link>
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
                <p>¿Eres profesor? Crea exámenes públicos para que toda la comunidad pueda resolver tus retos.</p>
            </div>
        </div>
    );
};

export default PublicExams;
