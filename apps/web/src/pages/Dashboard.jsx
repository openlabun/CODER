import { useState, useEffect, useContext } from 'react';
import { Link } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import client from '../api/client';
import { 
    Trophy, 
    BookOpen, 
    FileText, 
    CheckCircle, 
    Target, 
    BarChart3, 
    History, 
    Rocket,
    ArrowRight,
    Clock,
    Zap,
    TrendingUp,
    Star
} from 'lucide-react';
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
            if (!user?.id) return;
            
            try {
                // Fetch user-specific submissions to avoid 400 error on broad list
                const submissionsRes = await client.get(`/submissions/user/${user.id}`);
                const submissions = Array.isArray(submissionsRes.data) ? submissionsRes.data : (submissionsRes.data.items || []);

                const challengesRes = await client.get('/challenges/mine');
                const challenges = Array.isArray(challengesRes.data) ? challengesRes.data : (challengesRes.data.items || []);

                setStats({
                    totalSubmissions: submissions.length,
                    acceptedSubmissions: submissions.filter(s => s.status === 'accepted').length,
                    activeChallenges: challenges.filter(c => c.status === 'published').length, 
                    recentSubmissions: submissions.slice(0, 4)
                });
            } catch (err) {
                console.error('Error fetching dashboard data:', err);
            } finally {
                setLoading(false);
            }
        };

        if (user) {
            fetchDashboardData();
        }
    }, [user]);

    if (loading) {
        return (
            <div className="dashboard-loading">
                <div className="loader-orbit">
                    <div className="orbit-dot"></div>
                </div>
                <p>Cargando panel...</p>
            </div>
        );
    }

    const successRate = stats.totalSubmissions > 0
        ? Math.round((stats.acceptedSubmissions / stats.totalSubmissions) * 100)
        : 0;

    return (
        <div className="dashboard-compact">
            {/* Horizontal Hero & Metrics Combined for Compactness */}
            <div className="dashboard-top-section">
                <section className="dashboard-hero-min">
                    <div className="hero-content">
                        <div className="hero-badge">
                            <Zap size={12} fill="#ffc72c" color="#ffc72c" />
                            <span>Dashboard Académico</span>
                        </div>
                        <h1>Hola, <span className="text-highlight">{user?.username}</span></h1>
                        <p>Continúa mejorando tus habilidades de programación.</p>
                    </div>
                    <div className="hero-progress">
                        <div className="progress-mini" style={{ '--progress': `${successRate}%` }}>
                            <span className="val">{successRate}%</span>
                        </div>
                    </div>
                </section>

                <div className="metrics-compact-bar">
                    <div className="m-item">
                        <FileText size={18} className="blue" />
                        <div className="m-data">
                            <span className="v">{stats.totalSubmissions}</span>
                            <span className="l">Envíos</span>
                        </div>
                    </div>
                    <div className="m-item">
                        <CheckCircle size={18} className="green" />
                        <div className="m-data">
                            <span className="v">{stats.acceptedSubmissions}</span>
                            <span className="l">Aceptados</span>
                        </div>
                    </div>
                    <div className="m-item">
                        <Target size={18} className="red" />
                        <div className="m-data">
                            <span className="v">{stats.activeChallenges}</span>
                            <span className="l">Retos</span>
                        </div>
                    </div>
                    <div className="m-item">
                        <TrendingUp size={18} className="purple" />
                        <div className="m-data">
                            <span className="v">0</span>
                            <span className="l">Puntos</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* Main Content Areas */}
            <div className="dashboard-main-columns">
                
                {/* Activity List Section */}
                <section className="dashboard-card recent-activity">
                    <div className="card-header">
                        <div className="header-label">
                            <Clock size={16} />
                            <h2>Envíos Recientes</h2>
                        </div>
                        <Link to="/submissions" className="link-more">
                            Ver historial <ArrowRight size={12} />
                        </Link>
                    </div>

                    <div className="compact-list">
                        {stats.recentSubmissions.length > 0 ? (
                            stats.recentSubmissions.map((sub, idx) => (
                                <div key={sub.id || idx} className="c-row">
                                    <div className={`c-dot ${sub.status}`}></div>
                                    <div className="c-info">
                                        <span className="c-name">{sub.challengeId || 'Desafío'}</span>
                                        <span className="c-meta">{sub.language} | {new Date(sub.createdAt).toLocaleDateString()}</span>
                                    </div>
                                    <span className={`c-status ${sub.status}`}>{sub.status}</span>
                                </div>
                            ))
                        ) : (
                            <div className="empty-compact">
                                <Rocket size={24} />
                                <p>Sin actividad reciente</p>
                            </div>
                        )}
                    </div>
                </section>

                {/* Quick Side Actions */}
                <div className="side-actions-group">
                    <section className="dashboard-card quick-links">
                        <div className="header-label smaller">
                            <Star size={14} />
                            <h2>Acceso Rápido</h2>
                        </div>
                        <div className="links-grid">
                            <Link to="/challenges" className="l-box brand">
                                <Zap size={18} />
                                <span>Retos</span>
                            </Link>
                            <Link to="/courses" className="l-box yellow">
                                <BookOpen size={18} />
                                <span>Cursos</span>
                            </Link>
                            <Link to="/leaderboard" className="l-box purple">
                                <Trophy size={18} />
                                <span>Ranking</span>
                            </Link>
                            <Link to="/submissions" className="l-box blue">
                                <History size={18} />
                                <span>Envíos</span>
                            </Link>
                        </div>
                    </section>

                    <section className="tip-compact">
                        <p><Zap size={14} /> Modulariza tu código para mejorar la legibilidad y facilitar pruebas.</p>
                    </section>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
