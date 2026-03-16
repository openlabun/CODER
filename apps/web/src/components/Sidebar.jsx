import { Link, useLocation } from 'react-router-dom';
import { useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import './Sidebar.css';

const Sidebar = () => {
    const location = useLocation();
    const { user, logout } = useContext(AuthContext);

    const isActive = (path) => location.pathname.startsWith(path);

    return (
        <aside className="sidebar">
            <div className="sidebar-header">
                <Link to={user ? "/dashboard" : "/"} className="logo-link">
                    <div className="logo">Uninorte<span style={{ color: 'var(--secondary-color)' }}>Coder</span></div>
                </Link>
            </div>

            <nav className="sidebar-nav">
                <ul>
                    {user?.role === 'student' && (
                        <>
                            <li>
                                <Link to="/" className={isActive('/') && location.pathname === '/' ? 'active' : ''}>
                                    🏠 Inicio
                                </Link>
                            </li>
                            <li>
                                <Link to="/dashboard" className={isActive('/dashboard') ? 'active' : ''}>
                                    📊 Panel principal
                                </Link>
                            </li>
                            <li>
                                <Link to="/challenges" className={isActive('/challenges') ? 'active' : ''}>
                                    🚀 Retos
                                </Link>
                            </li>
                            <li>
                                <Link to="/courses" className={isActive('/courses') ? 'active' : ''}>
                                    📚 Mis Cursos
                                </Link>
                            </li>
                            <li>
                                <Link to="/submissions" className={isActive('/submissions') ? 'active' : ''}>
                                    📝 Envíos
                                </Link>
                            </li>
                            <li>
                                <Link to="/leaderboard" className={isActive('/leaderboard') ? 'active' : ''}>
                                    🏆 Clasificación
                                </Link>
                            </li>
                        </>
                    )}

                    {(user?.role === 'professor' || user?.role === 'admin') && (
                        <>
                            <li>
                                <Link to="/" className={isActive('/') && location.pathname === '/' ? 'active' : ''}>
                                    🏠 Inicio
                                </Link>
                            </li>
                            <li>
                                <Link to="/dashboard" className={isActive('/dashboard') ? 'active' : ''}>
                                    📊 Panel principal
                                </Link>
                            </li>
                            <li>
                                <Link to="/challenges" className={isActive('/challenges') ? 'active' : ''}>
                                    🚀 Retos
                                </Link>
                            </li>
                            <li>
                                <Link to="/challenges/create" className={isActive('/challenges/create') ? 'active' : ''}>
                                    ➕ Crear Reto
                                </Link>
                            </li>
                            <li>
                                <Link to="/courses" className={isActive('/courses') ? 'active' : ''}>
                                    📚 Cursos
                                </Link>
                            </li>
                            <li>
                                <Link to="/courses/create" className={isActive('/courses/create') ? 'active' : ''}>
                                    ➕ Crear Curso
                                </Link>
                            </li>
                            <li>
                                <Link to="/submissions" className={isActive('/submissions') ? 'active' : ''}>
                                    📝 Envíos
                                </Link>
                            </li>
                            <li>
                                <Link to="/leaderboard" className={isActive('/leaderboard') ? 'active' : ''}>
                                    🏆 Clasificación
                                </Link>
                            </li>
                        </>
                    )}
                </ul>
            </nav>

            <div className="sidebar-footer">
                <div className="user-profile">
                    <div className="user-avatar">
                        {user?.username?.charAt(0).toUpperCase()}
                    </div>
                    <div className="user-details">
                        <span className="username">{user?.username}</span>
                        <span className="role">{user?.role}</span>
                    </div>
                </div>
                <button onClick={logout} className="btn-logout">
                    <span>🚪</span> Cerrar Sesión
                </button>
            </div>
        </aside>
    );
};

export default Sidebar;
