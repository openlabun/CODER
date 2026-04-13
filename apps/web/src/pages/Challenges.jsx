import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import Swal from 'sweetalert2';
import { 
    AlertCircle, 
    Search, 
    Code, 
    Target,
    RotateCcw,
    Edit2,
    Trash2,
    Send,
    Archive,
    PlusCircle,
    Clock3,
    Eye,
    EyeOff,
    Filter
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
    const [visibilityFilter, setVisibilityFilter] = useState('all');
    const hasFetchedInitialData = useRef(false);

    const isTeacher = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    const fetchChallenges = async () => {
        setLoading(true);
        setError('');

        try {
            const endpoint = isTeacher ? '/challenges' : '/challenges/public';
            const { data } = await client.get(endpoint);
            const list = Array.isArray(data) ? data : (data.items || []);
            setChallenges(list);
        } catch (err) {
            console.error('Error loading challenges:', err);
            setError('Error al conectar con el servidor.');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user && user.role === 'student') {
            navigate('/public-exams');
            return;
        }

        if (!user || hasFetchedInitialData.current) {
            return;
        }

        hasFetchedInitialData.current = true;
        fetchChallenges();
    }, [user, navigate]);

    const updateChallengeStatus = (challengeId, nextStatus) => {
        setChallenges((prev) => prev.map((challenge) => {
            const currentId = challenge.id || challenge.ID;
            if (currentId !== challengeId) return challenge;

            return {
                ...challenge,
                status: nextStatus,
                Status: nextStatus,
            };
        }));
    };

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
            updateChallengeStatus(id, 'published');
            Swal.fire({ icon: 'success', title: 'Publicado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
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
            updateChallengeStatus(id, 'archived');
            Swal.fire({ icon: 'success', title: 'Archivado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
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
            setChallenges((prev) => prev.filter((challenge) => {
                const challengeId = challenge.id || challenge.ID;
                return challengeId !== id;
            }));
            Swal.fire({ icon: 'success', title: 'Eliminado', timer: 1000, toast: true, position: 'top-end', showConfirmButton: false });
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

    const normalizeStatus = (status) => String(status || 'draft').toLowerCase();

    const filteredChallenges = challenges.filter((c) => {
        const title = String(c.title || c.Title || '').toLowerCase();
        const matchesSearch = title.includes(searchTerm.toLowerCase());
        const status = normalizeStatus(c.status || c.Status);
        const matchesVisibility = visibilityFilter === 'all' ? true : status === visibilityFilter;
        return matchesSearch && matchesVisibility;
    });

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

    const getVisibilityMeta = (status) => {
        const s = normalizeStatus(status);
        if (s === 'published') return { label: 'Publicado', class: 'published', icon: <Eye size={14} /> };
        if (s === 'private') return { label: 'Privado', class: 'private', icon: <EyeOff size={14} /> };
        if (s === 'archived') return { label: 'Archivado', class: 'archived', icon: <Archive size={14} /> };
        return { label: 'Borrador', class: 'draft', icon: <Edit2 size={14} /> };
    };

    const getWorkerTimeLimit = (challenge) => (
        challenge.workerTimeLimit || challenge.WorkerTimeLimit || challenge.worker_time_limit || 1000
    );

    const visibilityFilters = [
        { key: 'all', label: 'Todos' },
        { key: 'draft', label: 'Borrador' },
        { key: 'published', label: 'Publicado' },
        { key: 'private', label: 'Privado' },
        { key: 'archived', label: 'Archivado' },
    ];

    return (
        <div className="challenges-page">
            <header className="page-header-compact">
                <div className="header-info">
                    <h1>Repositorio de Retos</h1>
                    <p>Gestiona tus desafíos de programación y comparte con la comunidad docente</p>
                </div>
                
                <div className="header-actions-mini">
                    <div className="search-bar-mini">
                        <Search size={18} />
                        <input 
                            type="text" 
                            placeholder="Buscar por título..." 
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                    {isTeacher && (
                        <button className="btn-create-mini" onClick={() => navigate('/challenges/create')}>
                            <PlusCircle size={18} />
                            <span>Nuevo Reto</span>
                        </button>
                    )}
                </div>
            </header>

            {isTeacher && (
                <section className="visibility-filter-row">
                    <div className="filter-label">
                        <Filter size={14} />
                        <span>Filtrar por visibilidad:</span>
                    </div>
                    <div className="filter-chips">
                        {visibilityFilters.map((filterItem) => (
                            <button
                                key={filterItem.key}
                                type="button"
                                className={`filter-chip ${visibilityFilter === filterItem.key ? 'active' : ''}`}
                                onClick={() => setVisibilityFilter(filterItem.key)}
                            >
                                {filterItem.label}
                            </button>
                        ))}
                    </div>
                </section>
            )}

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
                        const diff = getDifficultyLabel(challenge.difficulty || challenge.Difficulty);
                        const status = normalizeStatus(challenge.status || challenge.Status);
                        const isPublished = status === 'published';
                        const isArchived = status === 'archived';
                        const visibilityMeta = getVisibilityMeta(status);
                        const challengeId = challenge.id || challenge.ID;
                        const timeLimit = getWorkerTimeLimit(challenge);
                        
                        return (
                            <div key={challengeId} className={`challenge-card-mini ${isArchived ? 'archived' : ''}`}>
                                <div className={`card-accent ${diff.class}`}></div>
                                <div className="card-main">
                                    <div className="card-top">
                                        <div className="title-area">
                                            <Code size={16} className="title-icon" />
                                            <h3>{challenge.title || challenge.Title || 'Reto sin título'}</h3>
                                        </div>
                                    </div>

                                    <div className="meta-badges-row">
                                        {isTeacher && (
                                            <span className={`status-badge ${visibilityMeta.class}`}>
                                                {visibilityMeta.icon}
                                                {visibilityMeta.label}
                                            </span>
                                        )}
                                        <span className={`diff-pill ${diff.class}`}>
                                            {diff.label}
                                        </span>
                                    </div>

                                    {isTeacher && (
                                        <div className="teacher-actions teacher-actions-top">
                                            {!isPublished && (
                                                <button 
                                                    className="action-btn publish" 
                                                    title={isArchived ? 'Republicar' : 'Publicar'}
                                                    onClick={(e) => handlePublish(e, challengeId)}
                                                    disabled={processingId === challengeId}
                                                >
                                                    {isArchived ? <RotateCcw size={14} /> : <Send size={14} />}
                                                </button>
                                            )}
                                            {!isArchived && (
                                                <button 
                                                    className="action-btn archive" 
                                                    title="Archivar"
                                                    onClick={(e) => handleArchive(e, challengeId)}
                                                    disabled={processingId === challengeId}
                                                >
                                                    <Archive size={14} />
                                                </button>
                                            )}
                                            <button 
                                                className="action-btn edit" 
                                                title="Editar"
                                                onClick={(e) => handleEdit(e, challengeId)}
                                                disabled={processingId === challengeId}
                                            >
                                                <Edit2 size={14} />
                                            </button>
                                            <button 
                                                className="action-btn delete" 
                                                title="Eliminar"
                                                onClick={(e) => handleDelete(e, challengeId)}
                                                disabled={processingId === challengeId}
                                            >
                                                <Trash2 size={14} />
                                            </button>
                                        </div>
                                    )}

                                    <p className="description-text">{challenge.description || challenge.Description || 'Sin descripción disponible.'}</p>
                                    
                                    <div className="card-footer-mini">
                                        <div className="stats-mini">
                                            <div className="stat">
                                                <Clock3 size={14} />
                                                <span>{timeLimit} ms</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};

export default Challenges;
