import { useState, useContext } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import './Auth.css';

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [role, setRole] = useState('');
    const { register } = useContext(AuthContext);
    const navigate = useNavigate();
    const [error, setError] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await register(username, password, role);
            navigate('/dashboard');
        } catch (err) {
            setError('Registration failed. Try again.');
        }
    };

    return (
        <div className="auth-container">
            <div className="auth-card">
                <h2 className="auth-title">Regístrate</h2>
                <p className="auth-subtitle">Crea tu cuenta para comenzar a programar</p>

                {error && <div className="auth-error">{error}</div>}

                <form onSubmit={handleSubmit} className="auth-form">
                    <div className="form-group">
                        <label>Usuario</label>
                        <input
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="Neo"
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label>Contraseña</label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="••••••••"
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label>Rol</label>
                        <select
                            value={role}
                            onChange={(e) => setRole(e.target.value)}
                            className="role-select"
                            required
                        >
                            <option value="">Selecciona tu rol</option>
                            <option value="student">Estudiante</option>
                            <option value="professor">Profesor</option>
                        </select>
                    </div>
                    <button type="submit" className="btn-auth">Registrarse</button>
                </form>

                <p className="auth-footer">
                    ¿Ya tienes una cuenta? <Link to="/login">Ingresa aquí</Link>
                </p>
            </div>
        </div>
    );
};

export default Register;
