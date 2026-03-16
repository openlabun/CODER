import { useState, useEffect, useContext } from 'react';
import { Link } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import './Dashboard.css';

const Dashboard = () => {
    const { user } = useContext(AuthContext);
    const [stats, setStats] = useState({
        totalSubmissions: 0,
        acceptedSubmissions: 0,
        activeChallenges: 0,
        recentSubmissions: []
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchDashboardData = async () => {
            try {
                // Fetch submissions for stats
                const { data } = await client.get('/submissions');
                const submissions = data.items || [];

                setStats({
                    totalSubmissions: submissions.length,
                    acceptedSubmissions: submissions.filter(s => s.status === 'accepted').length,
                    activeChallenges: 0, // Will be populated when challenges endpoint is ready
                    recentSubmissions: submissions.slice(0, 5)
                });
            } catch (err) {
                console.error('Error fetching dashboard data:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchDashboardData();
    }, []);

    if (loading) {
        return <div className="loading">Cargando panel...</div>;
    }

    return (
        <div className="dashboard-page">
            <div className="page-header">
                <h1>Panel Principal</h1>
                <p className="welcome-text">¡Bienvenido de nuevo, {user?.username}!</p>
            </div>

            <div className="stats-grid">
                <div className="stat-card">
                    <div className="stat-icon">📝</div>
                    <div className="stat-content">
                        <h3>Envíos Totales</h3>
                        <div className="stat-value">{stats.totalSubmissions}</div>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon">✅</div>
                    <div className="stat-content">
                        <h3>Aceptados</h3>
                        <div className="stat-value">{stats.acceptedSubmissions}</div>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon">🎯</div>
                    <div className="stat-content">
                        <h3>Retos Activos</h3>
                        <div className="stat-value">{stats.activeChallenges}</div>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon">📊</div>
                    <div className="stat-content">
                        <h3>Tasa de Éxito</h3>
                        <div className="stat-value">
                            {stats.totalSubmissions > 0
                                ? Math.round((stats.acceptedSubmissions / stats.totalSubmissions) * 100)
                                : 0}%
                        </div>
                    </div>
                </div>
            </div>

            <div className="dashboard-content">
                <div className="quick-actions-section">
                    <h2>Acciones Rápidas</h2>
                    <div className="actions-grid">
                        <Link to="/challenges" className="action-card">
                            <span className="action-icon">🚀</span>
                            <h3>Explorar Retos</h3>
                            <p>Explora y resuelve retos de programación</p>
                        </Link>
                        <Link to="/courses" className="action-card">
                            <span className="action-icon">📚</span>
                            <h3>Mis Cursos</h3>
                            <p>Mira tus cursos inscritos</p>
                        </Link>
                        <Link to="/leaderboard" className="action-card">
                            <span className="action-icon">🏆</span>
                            <h3>Clasificación</h3>
                            <p>Revisa tu posición en el ranking</p>
                        </Link>
                        <Link to="/submissions" className="action-card">
                            <span className="action-icon">📋</span>
                            <h3>Mis Envíos</h3>
                            <p>Revisa tu historial de envíos</p>
                        </Link>
                    </div>
                </div>

                {stats.recentSubmissions.length > 0 && (
                    <div className="recent-activity-section">
                        <h2>Envíos Recientes</h2>
                        <div className="submissions-list">
                            {stats.recentSubmissions.map(sub => (
                                <div key={sub.id} className="submission-item">
                                    <span className={`status-badge ${sub.status}`}>{sub.status}</span>
                                    <span className="submission-challenge">{sub.challengeId}</span>
                                    <span className="submission-lang">{sub.language}</span>
                                    <span className="submission-date">
                                        {new Date(sub.createdAt).toLocaleDateString()}
                                    </span>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {stats.totalSubmissions === 0 && (
                    <div className="empty-state">
                        <h3>🎯 ¿Listo para empezar?</h3>
                        <p>Aún no has enviado ninguna solución. ¡Explora los retos y empieza a programar!</p>
                        <Link to="/challenges" className="btn-primary">Explorar Retos</Link>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Dashboard;
