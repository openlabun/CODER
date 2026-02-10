import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Navbar.css';

const Navbar = () => {
    const { user, logout } = useAuth();

    return (
        <nav className="navbar">
            <div className="navbar-brand">
                <Link to="/">Juez Online</Link>
            </div>
            <div className="navbar-links">
                <Link to="/challenges">Challenges</Link>
                {user ? (
                    <>
                        <span style={{ color: 'var(--text-secondary)' }}>Hello, {user.username || 'User'}!</span>
                        <button onClick={logout} className="btn-logout">Logout</button>
                    </>
                ) : (
                    <>
                        <Link to="/login">Login</Link>
                        <Link to="/register">Register</Link>
                    </>
                )}
            </div>
        </nav>
    );
};

export default Navbar;
