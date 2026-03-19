import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Navbar.css';

const Navbar = () => {
    const { user, logout } = useAuth();

    return (
        <nav className="navbar">
            <div className="navbar-brand">
                <Link to="/">UninorteCoder</Link>
            </div>
            <div className="navbar-links">
                <Link to="/challenges">Retos</Link>
                {user ? (
                    <>
                        <span style={{ color: 'var(--text-color)' }}>¡Hola, {user.username || 'Usuario'}!</span>
                        <button onClick={logout} className="btn-logout">Cerrar Sesión</button>
                    </>
                ) : (
                    <>
                        <Link to="/login">Ingresar</Link>
                        <Link to="/register">Registrarse</Link>
                    </>
                )}
            </div>
        </nav>
    );
};

export default Navbar;
