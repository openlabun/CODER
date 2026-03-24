import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import Swal from 'sweetalert2';
import { 
    AlertCircle, 
    Search, 
    Trophy, 
    Code, 
    ChevronRight, 
    Zap,
    Target,
    RotateCcw,
    Edit2,
    Trash2,
    Send,
    Archive,
    MoreVertical
} from 'lucide-react';
import './Challenges.css';

const Challenges = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const [challenges, setChallenges] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [processingId, setProcessingId] = useState(null);

    const isTeacher = user?.role === 'professor' || user?.role === 'admin';

    const fetchChallenges = async () => {
        try {
            const { data } = await client.get('/challenges/mine');
            setChallenges(Array.isArray(data) ? data : (data.items || []));
        } catch (err) {
            console.error('Error loading challenges:', err);
            setError('Error al conectar con el servidor.');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchChallenges();
    }, []);

    const handlePublish = async (e, id) => {
        e.preventDefault();
        
        const { isConfirmed } = await Swal.fire({
            title: '¿Publicar Reto?',
            text: 'Será visible para todos los estudiantes.',
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Sí, publicar',
            cancelButtonText: 'Cancelar'
        });

        if (!isConfirmed) return;
        
        setProcessingId(id);
        try {
            await client.post(`/challenges/${id}/publish`);
            Swal.fire({ icon: 'success', title: 'Publicado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
            fetchChallenges();
        } catch (err) {
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo publicar el reto.' });
        } finally {
            setProcessingId(null);
        }
    };

    const handleArchive = async (e, id) => {
        e.preventDefault();
        setProcessingId(id);
        try {
            await client.post(`/challenges/${id}/archive`);
            Swal.fire({ icon: 'success', title: 'Archivado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
            fetchChallenges();
        } catch (err) {
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo archivar.' });
        } finally {
            setProcessingId(null);
        }
    };

    const handleDelete = async (e, id) => {
        e.preventDefault();
        
        const { isConfirmed } = await Swal.fire({
            title: '¿Eliminar Reto?',
            text: 'Esta acción es permanente.',
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            confirmButtonText: 'Sí, eliminar',
            cancelButtonText: 'Cancelar'
        });

        if (!isConfirmed) return;

        setProcessingId(id);
        try {
            await client.delete(`/challenges/${id}`);
            Swal.fire({ icon: 'success', title: 'Eliminado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
            fetchChallenges();
        } catch (err) {
            Swal.fire({ icon: 'error', title: 'Error', text: 'No se pudo eliminar el reto.' });
        } finally {
            setProcessingId(null);
        }
    };

    const handleEdit = (e, id) => {
        e.preventDefault();
        navigate(`/challenges/edit/${id}`);
    };

    const filteredChallenges = challenges.filter(c => 
        c.title?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (loading) return (
        <div className="challenges-page">
            <div className="page-header-compact">
                <div className="skeleton title-skeleton"></div>
            </div>
            <div className="challenges-grid-compact">
                {[...Array(6)].map((_, i) => (
                    <div key={i} className="challenge-card-mini skeleton-card">
                        <div className="skeleton card-content-skeleton"></div>
                    </div>
                ))}
            </div>
        </div>
    );

    if (error) return (
        <div className="challenges-page">
            <div className="error-container">
                <AlertCircle size={48} />
                <h3>Error al cargar retos</h3>
                <p>{error}</p>
                <button onClick={() => window.location.reload()} className="btn-retry">
                    <RotateCcw size={16} /> Reintentar carga
                </button>
            </div>
        </div>
    );

    const getDifficultyLabel = (diff) => {
        const d = (diff || 'medium').toLowerCase();
        if (d === 'easy') return { label: 'Fácil', class: 'easy' };
        if (d === 'hard') return { label: 'Difícil', class: 'hard' };
        return { label: 'Medio', class: 'medium' };
    };

    return (
        <div className="challenges-page">
            <header className="page-header-compact">
                <div className="header-info">
                    <h1>Desafíos de Programación</h1>
                    <p>Pon a prueba tu lógica y sube en el ranking institucional</p>
                </div>
                
                <div className="search-bar-mini">
                    <Search size={18} />
                    <input 
                        type="text" 
                        placeholder="Buscar por título..." 
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
            </header>

            {filteredChallenges.length === 0 ? (
                <div className="empty-state-mini">
                    <div className="icon-circle">
                        <Target size={32} />
                    </div>
                    <h3>No se encontraron desafíos</h3>
                    <p>{searchTerm ? 'Prueba con otros términos de búsqueda.' : 'Aún no hay retos publicados.'}</p>
                </div>
            ) : (
                <div className="challenges-grid-compact">
                    {filteredChallenges.map((challenge) => {
                        const diff = getDifficultyLabel(challenge.difficulty);
                        const isPublished = challenge.status === 'published';
                        const isArchived = challenge.status === 'archived';
                        
                        return (
                            <Link to={`/challenge/${challenge.id}`} key={challenge.id} className={`challenge-card-mini ${isArchived ? 'archived' : ''}`}>
                                <div className={`card-accent ${diff.class}`}></div>
                                <div className="card-main">
                                    <div className="card-top">
                                        <div className="title-area">
                                            <Code size={16} className="title-icon" />
                                            <h3>{challenge.title}</h3>
                                        </div>
                                        <div className="badge-group">
                                            {isTeacher && (
                                                <span className={`status-badge ${challenge.status}`}>
                                                    {challenge.status}
                                                </span>
                                            )}
                                            <span className={`diff-pill ${diff.class}`}>
                                                {diff.label}
                                            </span>
                                        </div>
                                    </div>
                                    <p className="description-text">{challenge.description}</p>
                                    
                                    <div className="card-footer-mini">
                                        <div className="stats-mini">
                                            <div className="stat">
                                                <Trophy size={14} />
                                                <span>100 pts</span>
                                            </div>
                                            <div className="stat">
                                                <Zap size={14} />
                                                <span>{challenge.attempts || 0} envíos</span>
                                            </div>
                                        </div>
                                        
                                        <div className="actions-wrapper">
                                            {isTeacher && (
                                                <div className="teacher-actions">
                                                    {!isPublished && !isArchived && (
                                                        <button 
                                                            className="action-btn publish" 
                                                            title="Publicar"
                                                            onClick={(e) => handlePublish(e, challenge.id)}
                                                            disabled={processingId === challenge.id}
                                                        >
                                                            <Send size={14} />
                                                        </button>
                                                    )}
                                                    {!isArchived && (
                                                        <button 
                                                            className="action-btn archive" 
                                                            title="Archivar"
                                                            onClick={(e) => handleArchive(e, challenge.id)}
                                                            disabled={processingId === challenge.id}
                                                        >
                                                            <Archive size={14} />
                                                        </button>
                                                    )}
                                                    <button 
                                                        className="action-btn edit" 
                                                        title="Editar"
                                                        onClick={(e) => handleEdit(e, challenge.id)}
                                                        disabled={processingId === challenge.id}
                                                    >
                                                        <Edit2 size={14} />
                                                    </button>
                                                    <button 
                                                        className="action-btn delete" 
                                                        title="Eliminar"
                                                        onClick={(e) => handleDelete(e, challenge.id)}
                                                        disabled={processingId === challenge.id}
                                                    >
                                                        <Trash2 size={14} />
                                                    </button>
                                                </div>
                                            )}
                                            <div className="btn-action-mini">
                                                Resolver <ChevronRight size={16} />
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        );
                    })}
                </div>
            )}
        </div>
    );
};

export default Challenges;
