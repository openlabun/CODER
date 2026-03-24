import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import client from '../api/client';
import { 
    AlertCircle, 
    Search, 
    Trophy, 
    Code, 
    ChevronRight, 
    Zap,
    Target,
    RotateCcw
} from 'lucide-react';
import './Challenges.css';

const Challenges = () => {
    const [challenges, setChallenges] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchChallenges = async () => {
            try {
                const { data } = await client.get('/challenges');
                setChallenges(Array.isArray(data) ? data : (data.items || []));
            } catch (err) {
                console.error('Error loading challenges:', err);
                setError('Error al conectar con el servidor.');
            } finally {
                setLoading(false);
            }
        };
        fetchChallenges();
    }, []);

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
                        return (
                            <Link to={`/challenge/${challenge.id}`} key={challenge.id} className="challenge-card-mini">
                                <div className={`card-accent ${diff.class}`}></div>
                                <div className="card-main">
                                    <div className="card-top">
                                        <div className="title-area">
                                            <Code size={16} className="title-icon" />
                                            <h3>{challenge.title}</h3>
                                        </div>
                                        <span className={`diff-pill ${diff.class}`}>
                                            {diff.label}
                                        </span>
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
                                        <div className="btn-action-mini">
                                            Resolver <ChevronRight size={16} />
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
