import { useState, useContext } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { Mail, Lock, Code, AlertCircle, ArrowRight, Zap, Trophy, ShieldCheck, Eye, EyeOff } from 'lucide-react';
import './Auth.css';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const { login } = useContext(AuthContext);
    const navigate = useNavigate();
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e) => {
        if (e) e.preventDefault();
        setError('');
        setIsLoading(true);

        const trimmedPassword = password.trim();

        try {
            await login(email, trimmedPassword);
            navigate('/dashboard');
        } catch (err) {
            // Log the technical error for developers
            console.log('Login Error Context:', {
                message: err.message,
                technical: err
            });
            
            // Map technical/backend errors to friendly user messages
            let userMessage = 'No pudimos verificar tus datos. Revisa tu correo e intenta nuevamente.';
            
            const errorStr = (err.message || '').toLowerCase();
            
            if (errorStr.includes('401') || errorStr.includes('password') || errorStr.includes('incorrecta')) {
                userMessage = 'Tu correo o contraseña no son correctos. Por favor, verifícalos.';
            } else if (errorStr.includes('404') || errorStr.includes('not found') || errorStr.includes('no registrado')) {
                userMessage = 'Esta cuenta no está registrada. ¿Deseas crear una nueva?';
            } else if (errorStr.includes('network') || errorStr.includes('500') || errorStr.includes('conn')) {
                userMessage = 'Estamos teniendo problemas de conexión. Intenta de nuevo en unos momentos.';
            } else if (errorStr.includes('timeout')) {
                userMessage = 'La conexión ha tardado demasiado. Revisa tu internet.';
            }

            setError(userMessage);
        } finally {
            setIsLoading(false);
        }
    };

    const handleEmailKeyDown = (e) => {
        if (e.key === 'Enter') {
            const currentEmail = email.trim();
            if (currentEmail !== '' && !currentEmail.includes('@')) {
                e.preventDefault();
                setEmail(`${currentEmail}@uninorte.edu.co`);
            }
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-split">
                {/* Left Side: Form */}
                <div className="auth-side-form">
                    <div className="auth-card">
                        <div className="auth-header">
                            <h2 className="auth-title">RobleCode</h2>
                            <p className="auth-subtitle">Accede a tus cursos y desafíos</p>
                        </div>

                        {error && (
                            <div className="auth-error">
                                <AlertCircle size={18} />
                                <span>{error}</span>
                            </div>
                        )}

                        <form onSubmit={handleSubmit} className="auth-form">
                            <div className="form-group">
                                <label>Email Institucional</label>
                                 <div className="input-wrapper">
                                    <Mail className="input-icon" size={20} />
                                    <input
                                        type="email"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        onKeyDown={handleEmailKeyDown}
                                        placeholder="usuario@uninorte.edu.co"
                                        required
                                    />
                                </div>
                                <button 
                                    type="button" 
                                    className="btn-domain-helper"
                                    onClick={() => {
                                        if (!email.includes('@')) {
                                            setEmail(prev => `${prev}@uninorte.edu.co`);
                                        }
                                    }}
                                >
                                    <Zap size={12} /> Añadir @uninorte.edu.co
                                </button>
                            </div>
                            <div className="form-group">
                                <label>Contraseña</label>
                                <div className="input-wrapper">
                                    <Lock className="input-icon" size={20} />
                                    <input
                                        type={showPassword ? "text" : "password"}
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        placeholder="••••••••"
                                        required
                                    />
                                    <button 
                                        type="button"
                                        className="password-toggle-btn"
                                        onClick={() => setShowPassword(!showPassword)}
                                        style={{ background: 'none', border: 'none', cursor: 'pointer', position: 'absolute', right: '1rem', top: '50%', transform: 'translateY(-50%)', color: '#6b7280', display: 'flex', alignItems: 'center', justifyContent: 'center' }}
                                    >
                                        {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                                    </button>
                                </div>
                            </div>
                            
                            <button type="submit" className="btn-auth" disabled={isLoading}>
                                {isLoading ? 'Validando...' : (
                                    <>
                                        Entrar Ahora <ArrowRight size={20} style={{ marginLeft: '12px' }} />
                                    </>
                                )}
                            </button>
                        </form>

                        <p className="auth-footer">
                            ¿Aún no tienes cuenta? <Link to="/register">Crea una gratis</Link>
                        </p>
                    </div>
                </div>

                {/* Right Side: Brand/Promo */}
                <div className="auth-side-brand">
                    <div className="brand-content">
                        <h1 className="brand-name">RobleCode</h1>
                        <p className="brand-tagline">
                            La plataforma académica oficial de Uninorte para potenciar tus habilidades de programación competitiva.
                        </p>
                        
                        <div className="brand-features">
                            <div className="feature-tag">
                                <Zap size={14} style={{ marginRight: '6px' }} />
                                Ejecución Rápida
                            </div>
                            <div className="feature-tag">
                                <Trophy size={14} style={{ marginRight: '6px' }} />
                                Ranking Real
                            </div>
                            <div className="feature-tag">
                                <ShieldCheck size={14} style={{ marginRight: '6px' }} />
                                Seguro
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login;
