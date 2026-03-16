import { Link } from 'react-router-dom';
import './Home.css';

const Home = () => {
    return (
        <div className="home-container">
            <div className="hero-section">
                <h1 className="hero-title">
                    Evaluación de Algoritmos <br />
                    <span style={{ color: 'var(--primary-color)' }}>Universidad del Norte</span>
                </h1>
                <p className="hero-subtitle">
                    La plataforma institucional para el fortalecimiento de habilidades de programación.
                    Resuelve retos, compite sanamente y prepárate para los desafíos de la ingeniería.
                </p>
                <div className="hero-actions">
                    <Link to="/register" className="btn btn-primary">Comenzar Ahora</Link>
                    <Link to="/login" className="btn btn-secondary">Ingresar</Link>
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

            <div className="info-section">
                <h2 className="section-title">¿Cómo Funciona?</h2>
                <div className="steps-container">
                    <div className="step-box">
                        <div className="step-number">1</div>
                        <h3>Regístrate</h3>
                        <p>Únete con tu cuenta institucional y configura tu perfil.</p>
                    </div>
                    <div className="step-box">
                        <div className="step-number">2</div>
                        <h3>Elige un Reto</h3>
                        <p>Selecciona ejercicios acordes a tu nivel de programación.</p>
                    </div>
                    <div className="step-box">
                        <div className="step-number">3</div>
                        <h3>Escribe tu Código</h3>
                        <p>Desarrolla la solución en el lenguaje de tu preferencia y envíalo.</p>
                    </div>
                    <div className="step-box">
                        <div className="step-number">4</div>
                        <h3>Obtén Resultados</h3>
                        <p>El sistema verificará la eficiencia y exactitud en tiempo real.</p>
                    </div>
                </div>
            </div>

            <div className="cta-section">
                <h2>¿Listo para mejorar tu lógica algorítmica?</h2>
                <p>Únete a la comunidad de estudiantes y comienza a practicar hoy mismo.</p>
                <Link to="/register" className="btn btn-primary">Crear Cuenta Gratuita</Link>
            </div>
        </div>
    );
};

export default Home;
