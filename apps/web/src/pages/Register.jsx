import { useState, useContext } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { Mail, Lock, User, Code, AlertCircle, UserPlus, Zap, Trophy, ShieldCheck, Eye, EyeOff } from 'lucide-react';
import './Auth.css';

const Register = () => {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const { register } = useContext(AuthContext);
    const navigate = useNavigate();
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setIsLoading(true);
        const trimmedPassword = password.trim();
        try {
            await register(name, email, trimmedPassword);
            navigate('/dashboard');
        } catch (err) {
            setError(err.message || 'Error al crear la cuenta. Por favor intenta de nuevo.');
        } finally {
            setIsLoading(false);
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
                            <p className="auth-subtitle">Únete a la comunidad de RobleCode</p>
                        </div>

                        {error && (
                            <div className="auth-error">
                                <AlertCircle size={18} />
                                <span>{error}</span>
                            </div>
                        )}

                        <form onSubmit={handleSubmit} className="auth-form">
                            <div className="form-group">
                                <label>Nombre Completo</label>
                                <div className="input-wrapper">
                                    <User className="input-icon" size={20} />
                                    <input
                                        type="text"
                                        value={name}
                                        onChange={(e) => setName(e.target.value)}
                                        placeholder="Ej. Juan Pérez"
                                        required
                                    />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Email Institucional</label>
                                <div className="input-wrapper">
                                    <Mail className="input-icon" size={20} />
                                    <input
                                        type="email"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        placeholder="usuario@uninorte.edu.co"
                                        required
                                    />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Contraseña</label>
                                <div className="input-wrapper">
                                    <Lock className="input-icon" size={20} />
                                    <input
                                        type={showPassword ? "text" : "password"}
                                        value={password}
                                        onChange={(e) => setPassword(e.target.value)}
                                        placeholder="Mínimo 8 caracteres"
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
                                {isLoading ? 'Procesando...' : (
                                    <>
                                        <UserPlus size={20} style={{ marginRight: '12px' }} />
                                        Registrarme
                                    </>
                                )}
                            </button>
                        </form>

                        <p className="auth-footer">
                            ¿Ya tienes cuenta? <Link to="/login">Inicia sesión</Link>
                        </p>
                    </div>
                </div>

                {/* Right Side: Brand/Promo */}
                <div className="auth-side-brand">
                    <div className="brand-content">
                        <h1 className="brand-name">RobleCode</h1>
                        <p className="brand-tagline">
                            Lleva tu lógica al siguiente nivel. Resuelve desafíos, compite con tus compañeros y domina el arte de programar.
                        </p>
                        
                        <div className="brand-features">
                            <div className="feature-tag">
                                <Zap size={14} style={{ marginRight: '6px' }} />
                                Enfoque Académico
                            </div>
                            <div className="feature-tag">
                                <Trophy size={14} style={{ marginRight: '6px' }} />
                                Ranking Real
                            </div>
                            <div className="feature-tag">
                                <ShieldCheck size={14} style={{ marginRight: '6px' }} />
                                Verificado
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Register;
