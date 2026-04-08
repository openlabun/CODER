import { Link } from 'react-router-dom';
import { 
    Code2, 
    Zap, 
    Award, 
    CheckCircle2, 
    ArrowRight, 
    Terminal, 
    ShieldCheck, 
    Users,
    ChevronRight
} from 'lucide-react';
import './Home.css';

const Home = () => {
    return (
        <div className="home-container">
            {/* Hero Section */}
            <header className="hero-section">
                <div className="hero-badge">
                    <span className="badge-dot"></span>
                    Plataforma Oficial de Evaluación
                </div>
                <h1 className="hero-title">
                    Domina el arte de la <br />
                    <span className="text-gradient">Programación Algorítmica</span>
                </h1>
                <p className="hero-subtitle">
                    La plataforma de la Universidad del Norte diseñada para potenciar 
                    tus habilidades lógicas a través de retos reales y evaluación automática.
                </p>
                <div className="hero-actions">
                    <Link to="/register" className="btn btn-primary btn-with-icon">
                        Comenzar Ahora <ArrowRight size={20} />
                    </Link>
                    <Link to="/login" className="btn btn-outline">
                        Entrar a mi Cuenta
                    </Link>
                </div>
                
                <div className="hero-stats">
                    <div className="stat-item">
                        <strong>+50</strong>
                        <span>Retos Activos</span>
                    </div>
                    <div className="stat-separator"></div>
                    <div className="stat-item">
                        <strong>4</strong>
                        <span>Lenguajes</span>
                    </div>
                    <div className="stat-separator"></div>
                    <div className="stat-item">
                        <strong>100%</strong>
                        <span>Seguro</span>
                    </div>
                </div>
            </header>

            {/* Features Section */}
            <section className="features-section">
                <div className="section-header">
                    <h2 className="section-title">¿Por qué usar RobleCode?</h2>
                    <p className="section-subtitle">Herramientas diseñadas para la excelencia académica</p>
                </div>
                
                <div className="features-grid">
                    <div className="feature-card">
                        <div className="feature-icon-wrapper">
                            <Code2 className="feature-icon" />
                        </div>
                        <h3>Soporte Multilenguaje</h3>
                        <p>Resuelve ejercicios en Python, Java, C++, o Node.js con entornos de ejecución optimizados.</p>
                        <ul className="feature-list">
                            <li><CheckCircle2 size={14} /> Sandboxing seguro</li>
                            <li><CheckCircle2 size={14} /> Versiones actualizadas</li>
                        </ul>
                    </div>

                    <div className="feature-card highlighted">
                        <div className="feature-icon-wrapper">
                            <Zap className="feature-icon" />
                        </div>
                        <h3>Evaluación Instantánea</h3>
                        <p>Recibe retroalimentación en segundos con un análisis detallado de tus casos de prueba.</p>
                        <ul className="feature-list">
                            <li><CheckCircle2 size={14} /> Veredictos en tiempo real</li>
                            <li><CheckCircle2 size={14} /> Análisis de eficiencia</li>
                        </ul>
                    </div>

                    <div className="feature-card">
                        <div className="feature-icon-wrapper">
                            <Award className="feature-icon" />
                        </div>
                        <h3>Ranking y Analítica</h3>
                        <p>Mide tu progreso académico y compite sanamente con tus compañeros de curso.</p>
                        <ul className="feature-list">
                            <li><CheckCircle2 size={14} /> Tableros de posiciones</li>
                            <li><CheckCircle2 size={14} /> Historial de intentos</li>
                        </ul>
                    </div>
                </div>
            </section>

            {/* How it Works */}
            <section className="steps-section">
                <div className="steps-header">
                    <h2 className="section-title">Tu camino al éxito</h2>
                    <div className="steps-progress-bar"></div>
                </div>

                <div className="steps-grid">
                    <div className="step-card">
                        <span className="step-number">01</span>
                        <div className="step-icon"><Users size={24} /></div>
                        <h3>Crea tu Perfil</h3>
                        <p>Regístrate con tus datos institucionales para acceder a tus cursos.</p>
                    </div>
                    <div className="step-card">
                        <span className="step-number">02</span>
                        <div className="step-icon"><Terminal size={24} /></div>
                        <h3>Elige un Reto</h3>
                        <p>Explora la biblioteca de ejercicios y selecciona tu próximo desafío.</p>
                    </div>
                    <div className="step-card">
                        <span className="step-number">03</span>
                        <div className="step-icon"><ShieldCheck size={24} /></div>
                        <h3>Envía y Valida</h3>
                        <p>Envía tu solución y deja que nuestro motor de evaluación verifique tu lógica.</p>
                    </div>
                </div>
            </section>

            {/* CTA Section */}
            <section className="cta-container">
                <div className="cta-content">
                    <h2>Preparado para el siguiente nivel?</h2>
                    <p>Únete a cientos de estudiantes que ya están mejorando sus habilidades hoy.</p>
                    <div className="cta-buttons">
                        <Link to="/register" className="btn btn-primary btn-large">
                            Registrarme ahora
                        </Link>
                        <Link to="/about" className="btn-link">
                            Aprender más <ChevronRight size={16} />
                        </Link>
                    </div>
                </div>
                <div className="cta-decoration">
                    <Code2 size={120} opacity={0.05} />
                </div>
            </section>

            <footer className="home-footer">
                <p>&copy; 2026 Universidad del Norte - Departamento de Ingeniería de Sistemas</p>
            </footer>
        </div>
    );
};

export default Home;
