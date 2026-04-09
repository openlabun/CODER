import { useState, useEffect, useContext } from 'react';
import { 
    FileText, 
    CheckCircle, 
    XCircle, 
    Clock, 
    Code, 
    Cpu, 
    Calendar,
    ChevronRight,
    AlertCircle,
    Terminal,
    Trophy,
    RotateCcw
} from 'lucide-react';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import './Submissions.css';

const Submissions = () => {
    const [submissions, setSubmissions] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const { user } = useContext(AuthContext);
    const isProfessor = user?.role === 'professor' || user?.role === 'teacher' || user?.role === 'admin';

    useEffect(() => {
        const fetchSubmissions = async () => {
            if (!user) return;
            const userId = user.id || user.ID;
            if (!userId) {
                setLoading(false);
                return;
            }
            try {
                // Students use /submissions/{id}, professors use /submissions/user/{id}
                const endpoint = isProfessor 
                    ? `/submissions/user/${userId}` 
                    : `/submissions/${userId}`;
                const { data } = await client.get(endpoint);
                setSubmissions(Array.isArray(data) ? data : (data.items || []));
            } catch (err) {
                console.error('Error loading submissions:', err);
                // Don't show error for students if endpoint returns empty
                if (!isProfessor) {
                    setSubmissions([]);
                } else {
                    setError('No se pudieron cargar los envíos.');
                }
            } finally {
                setLoading(false);
            }
        };
        if (user) fetchSubmissions();
    }, [user]);

    const getStatusLabel = (status) => {
        const s = (status || 'pending').toLowerCase();
        if (s === 'accepted' || s === 'success') return { label: 'Aceptado', class: 'accepted', icon: <CheckCircle size={14} /> };
        if (s === 'wrong_answer' || s === 'rejected' || s === 'failed') return { label: 'Rechazado', class: 'rejected', icon: <XCircle size={14} /> };
        if (s === 'runtime_error') return { label: 'Error', class: 'error', icon: <AlertCircle size={14} /> };
        return { label: 'Pendiente', class: 'pending', icon: <Clock size={14} /> };
    };

    if (loading) return (
        <div className="submissions-page-mini">
            <header className="page-header-mini">
                <div className="skeleton title-skeleton"></div>
            </header>
            <div className="skeleton-table-mini">
                {[...Array(6)].map((_, i) => (
                    <div key={i} className="skeleton-row-mini shimmer"></div>
                ))}
            </div>
        </div>
    );

    return (
        <div className="submissions-page-mini">
            <header className="page-header-mini">
                <div className="header-info-mini">
                    <h1>Historial de Envíos</h1>
                    <p>Revisa el progreso de tus soluciones y el feedback del sistema</p>
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
            ) : submissions.length === 0 ? (
                <div className="empty-state-mini">
                    <div className="icon-circle-mini">
                        <Terminal size={32} />
                    </div>
                    <h3>{isProfessor ? 'Sin envíos todavía' : 'Historial de Envíos'}</h3>
                    <p>{isProfessor 
                        ? '¡Los envíos de tus estudiantes aparecerán aquí!' 
                        : 'Tus envíos se registran dentro de cada examen que resuelves. Accede a un examen desde "Exámenes Públicos" o desde tu curso para ver tu progreso.'}
                    </p>
                </div>
            ) : (
                <div className="submissions-list-mini">
                    <div className="list-columns-header">
                        <div className="col-challenge">DESAFÍO</div>
                        <div className="col-status">ESTADO</div>
                        <div className="col-score">PUNTAJE</div>
                        <div className="col-meta">DATOS TÉCNICOS</div>
                        <div className="col-date">FECHA</div>
                        <div className="col-action"></div>
                    </div>
                    {submissions.map((sub) => {
                        const status = getStatusLabel(sub.status);
                        return (
                            <div key={sub.id} className="submission-item-mini">
                                <div className="col-challenge">
                                    <div className="challenge-info-mini">
                                        <Code size={16} className="icon-muted" />
                                        <div className="text-wrap">
                                            <span className="challenge-name">{sub.challengeTitle || sub.challengeId || 'Cargando...'}</span>
                                            <span className="sub-id">#{sub.id.slice(0, 8)}</span>
                                        </div>
                                    </div>
                                </div>

                                <div className="col-status">
                                    <span className={`status-pill-mini ${status.class}`}>
                                        {status.icon}
                                        {status.label}
                                    </span>
                                </div>

                                <div className="col-score">
                                    <div className="score-badge-mini">
                                        <Trophy size={14} />
                                        <span>{sub.score || 0}%</span>
                                    </div>
                                </div>

                                <div className="col-meta">
                                    <div className="technical-meta-mini">
                                        <div className="meta-pair">
                                            <Cpu size={12} />
                                            <span>{sub.memoryMbTotal || 0} MB</span>
                                        </div>
                                        <div className="meta-pair">
                                            <Clock size={12} />
                                            <span>{sub.timeMsTotal || 0} ms</span>
                                        </div>
                                        <div className="meta-pair capitalize">
                                            <span>{sub.language}</span>
                                        </div>
                                    </div>
                                </div>

                                <div className="col-date">
                                    <div className="date-box-mini">
                                        <Calendar size={14} />
                                        <span>{new Date(sub.createdAt).toLocaleDateString()}</span>
                                    </div>
                                </div>

                                <div className="col-action">
                                    <button className="btn-view-details">
                                        Detalles <ChevronRight size={14} />
                                    </button>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};

export default Submissions;
