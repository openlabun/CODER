import { Link } from 'react-router-dom';
import './Home.css';

const Home = () => {
    return (
        <div className="home-container">
            <div className="hero-section">
                <h1 className="hero-title">
                    Evaluación de Algoritmos <br />
                    <span className="neon-text" style={{ color: 'var(--primary-color)' }}>Universidad del Norte</span>
                </h1>
                <p className="hero-subtitle">
                    La plataforma institucional para el fortalecimiento de habilidades de programación.
                    Resuelve retos, compite sanamente y prepárate para los desafíos de la ingeniería.
                </p>
                <div className="hero-actions">
                    <Link to="/register" className="btn btn-primary glow-effect">Comenzar Ahora</Link>
                    <Link to="/login" className="btn btn-secondary" style={{ backgroundColor: 'var(--surface-color)', color: 'var(--text-color)', border: '1px solid var(--border-color)' }}>Ingresar</Link>
                </div>
            </div>

            <div className="features-grid">
                <div className="feature-card" style={{ backgroundColor: 'var(--surface-color)', border: 'var(--glass-border)', boxShadow: 'var(--shadow-glow)' }}>
                    <h3 style={{ color: 'var(--primary-color)' }}>🚀 Soporte Multilenguaje</h3>
                    <p>Programá en Python, Java, C++, o Node.js. Nuestra infraestructura aislada se encarga del resto.</p>
                </div>
                <div className="feature-card" style={{ backgroundColor: 'var(--surface-color)', border: 'var(--glass-border)', boxShadow: 'var(--shadow-glow)' }}>
                    <h3 style={{ color: 'var(--primary-color)' }}>⚡ Evaluación Instantánea</h3>
                    <p>Retroalimentación inmediata en tus entregas. Supera los casos de prueba y escala en el ranking.</p>
                </div>
                <div className="feature-card" style={{ backgroundColor: 'var(--surface-color)', border: 'var(--glass-border)', boxShadow: 'var(--shadow-glow)' }}>
                    <h3 style={{ color: 'var(--primary-color)' }}>🏆 Ranking Académico</h3>
                    <p>Mide tu desempeño frente a otros estudiantes. Demuestra tus capacidades algorítmicas.</p>
                </div>
            </div>
        </div>
    );
};

export default Home;
