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
        return <div className="loading">Loading dashboard...</div>;
    }

    return (
        <div className="dashboard-page">
            <div className="page-header">
                <h1>Dashboard</h1>
                <p className="welcome-text">Welcome back, {user?.username}!</p>
            </div>

            <div className="stats-grid">
                <div className="stat-card">
                    <div className="stat-icon">📝</div>
                    <div className="stat-content">
                        <h3>Total Submissions</h3>
                        <div className="stat-value">{stats.totalSubmissions}</div>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon">✅</div>
                    <div className="stat-content">
                        <h3>Accepted</h3>
                        <div className="stat-value">{stats.acceptedSubmissions}</div>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon">🎯</div>
                    <div className="stat-content">
                        <h3>Active Challenges</h3>
                        <div className="stat-value">{stats.activeChallenges}</div>
                    </div>
                </div>
                <div className="stat-card">
                    <div className="stat-icon">📊</div>
                    <div className="stat-content">
                        <h3>Success Rate</h3>
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
                    <h2>Quick Actions</h2>
                    <div className="actions-grid">
                        <Link to="/challenges" className="action-card">
                            <span className="action-icon">🚀</span>
                            <h3>Browse Challenges</h3>
                            <p>Explore and solve coding challenges</p>
                        </Link>
                        <Link to="/courses" className="action-card">
                            <span className="action-icon">📚</span>
                            <h3>My Courses</h3>
                            <p>View your enrolled courses</p>
                        </Link>
                        <Link to="/leaderboard" className="action-card">
                            <span className="action-icon">🏆</span>
                            <h3>Leaderboard</h3>
                            <p>Check your ranking</p>
                        </Link>
                        <Link to="/submissions" className="action-card">
                            <span className="action-icon">📋</span>
                            <h3>My Submissions</h3>
                            <p>Review your submission history</p>
                        </Link>
                    </div>
                </div>

                {stats.recentSubmissions.length > 0 && (
                    <div className="recent-activity-section">
                        <h2>Recent Submissions</h2>
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
                        <h3>🎯 Ready to start?</h3>
                        <p>You haven't submitted any solutions yet. Browse challenges and start coding!</p>
                        <Link to="/challenges" className="btn-primary">Browse Challenges</Link>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Dashboard;
