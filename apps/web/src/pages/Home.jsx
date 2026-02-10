import { Link } from 'react-router-dom';
import './Home.css';

const Home = () => {
    return (
        <div className="home-container">
            <div className="hero-section">
                <h1 className="hero-title">
                    Master Algorithms <br />
                    <span className="neon-text">Level Up Your Code</span>
                </h1>
                <p className="hero-subtitle">
                    The ultimate platform for competitive programming and algorithmic challenges.
                    Compete, learn, and dominate the leaderboard.
                </p>
                <div className="hero-actions">
                    <Link to="/register" className="btn btn-primary glow-effect">Get Started</Link>
                    <Link to="/login" className="btn btn-secondary">Login</Link>
                </div>
            </div>

            <div className="features-grid">
                <div className="feature-card">
                    <h3>🚀 Multi-Language Support</h3>
                    <p>Code in Python, Java, C++, or Node.js. Our isolated runners handle it all.</p>
                </div>
                <div className="feature-card">
                    <h3>⚡ Real-time Evaluation</h3>
                    <p>Instant feedback on your submissions. Pass the test cases and climb the ranks.</p>
                </div>
                <div className="feature-card">
                    <h3>🏆 Global Leaderboard</h3>
                    <p>Compete with students worldwide. Prove you are the best algorithmist.</p>
                </div>
            </div>
        </div>
    );
};

export default Home;
