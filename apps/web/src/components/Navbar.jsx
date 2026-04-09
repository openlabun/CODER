import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Trophy, LogIn, UserPlus, LogOut, User, LayoutDashboard, BookOpen, Medal, CircleUser, Menu, X } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useState, useEffect } from 'react';
import './Navbar.css';

const Navbar = () => {
    const { user, logout } = useAuth();
    const location = useLocation();
    const [scrolled, setScrolled] = useState(false);
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    const isActive = (path) => location.pathname === path;

    useEffect(() => {
        const handleScroll = () => {
            setScrolled(window.scrollY > 20);
        };
        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, []);

    const navLinks = user ? [
        { path: '/dashboard', label: 'PANEL PRINCIPAL', icon: <LayoutDashboard size={18} /> },
        { path: '/courses', label: 'Cursos', icon: <BookOpen size={18} /> },
        { path: '/challenges', label: 'Retos', icon: <Trophy size={18} /> },
        { path: '/leaderboard', label: 'Ranking', icon: <Medal size={18} /> },
    ] : [];

    return (
        <nav className={`navbar ${scrolled ? 'scrolled' : ''}`}>
            <div className="navbar-container">
                <Link to="/" className="navbar-brand">
                    <div className="brand-logo-wrapper">
                        <img src="/logo.png" alt="RobleCode Logo" className="brand-logo-img" />
                    </div>
                    <div className="brand-text-container">
                        <span className="brand-logo-text">RobleCode</span>
                        <span className="brand-tag">UNINORTE ACADEMY</span>
                    </div>
                </Link>

                {/* Desktop Navigation */}
                <div className="navbar-desktop">
                    <div className="nav-main-links">
                        {navLinks.map((link) => (
                            <Link 
                                key={link.path}
                                to={link.path} 
                                className={`nav-item ${isActive(link.path) ? 'active' : ''}`}
                            >
                                {link.icon}
                                <span>{link.label}</span>
                                {isActive(link.path) && (
                                    <motion.div layoutId="nav-active" className="nav-active-pill" />
                                )}
                            </Link>
                        ))}
                    </div>

                    <div className="nav-divider"></div>

                    <div className="user-section">
                        {user ? (
                            <div className="user-profile-group">
                                <div className="user-info">
                                    <span className="user-name">{user.username || 'Usuario'}</span>
                                    <span className={`user-rank ${user.role || 'student'}`}>
                                        {user.role === 'teacher' || user.role === 'professor' ? 'Profesor' : 
                                         user.role === 'admin' ? 'Administrador' : 'Estudiante'}
                                    </span>
                                </div>
                                <div className="user-avatar-circle">
                                    <User size={18} />
                                </div>
                                <button onClick={logout} className="btn-logout-minimal" title="Cerrar Sesión">
                                    <LogOut size={18} />
                                </button>
                            </div>
                        ) : (
                            <div className="auth-actions">
                                <Link to="/login" className="btn-login-minimal">Entrar</Link>
                                <Link to="/register" className="btn-register-solid">Registro</Link>
                            </div>
                        )}
                    </div>
                </div>

                {/* Mobile Toggle */}
                <button className="mobile-toggle" onClick={() => setMobileMenuOpen(!mobileMenuOpen)}>
                    {mobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
                </button>
            </div>

            {/* Mobile Menu */}
            <AnimatePresence>
                {mobileMenuOpen && (
                    <motion.div 
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="navbar-mobile"
                    >
                        {user && (
                            <div className="mobile-user-section">
                                <div className="user-avatar-circle">
                                    <User size={18} />
                                </div>
                                <div className="user-info">
                                    <span className="user-name">{user.username || 'Usuario'}</span>
                                    <span className={`user-rank ${user.role || 'student'}`}>
                                        {user.role === 'teacher' || user.role === 'professor' ? 'Profesor' : 
                                         user.role === 'admin' ? 'Administrador' : 'Estudiante'}
                                    </span>
                                </div>
                            </div>
                        )}
                        {navLinks.map((link) => (
                            <Link 
                                key={link.path}
                                to={link.path} 
                                className={`mobile-nav-item ${isActive(link.path) ? 'active' : ''}`}
                                onClick={() => setMobileMenuOpen(false)}
                            >
                                {link.icon}
                                <span>{link.label}</span>
                            </Link>
                        ))}
                        {user ? (
                            <button onClick={() => { logout(); setMobileMenuOpen(false); }} className="mobile-nav-item logout">
                                <LogOut size={18} />
                                <span>Cerrar Sesión</span>
                            </button>
                        ) : (
                            <>
                                <Link to="/login" className="mobile-nav-item" onClick={() => setMobileMenuOpen(false)}>Ingresar</Link>
                                <Link to="/register" className="mobile-nav-item highlight" onClick={() => setMobileMenuOpen(false)}>Registrarse</Link>
                            </>
                        )}
                    </motion.div>
                )}
            </AnimatePresence>
        </nav>
    );
};

export default Navbar;
