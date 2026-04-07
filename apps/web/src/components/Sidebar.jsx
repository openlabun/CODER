import { Link, useLocation } from 'react-router-dom';
import { useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { 
    LayoutDashboard, 
    Trophy, 
    BookOpen, 
    FileText, 
    BarChart3, 
    PlusCircle, 
    LogOut,
    Menu,
    ChevronRight,
    Settings,
    ShieldCheck
} from 'lucide-react';
import './Sidebar.css';

const Sidebar = () => {
    const location = useLocation();
    const { user, logout } = useContext(AuthContext);

    const isActive = (path) => location.pathname === path;

    const navItems = user?.role === 'student' ? [
        { path: '/dashboard', label: 'PANEL PRINCIPAL', icon: <LayoutDashboard size={20} />, id: 'tour-nav-dashboard' },
        { path: '/challenges', label: 'Retos', icon: <Trophy size={20} />, id: 'tour-nav-challenges' },
        { path: '/courses', label: 'Mis Cursos', icon: <BookOpen size={20} />, id: 'tour-nav-courses' },
        { path: '/submissions', label: 'Envíos', icon: <FileText size={20} />, id: 'tour-nav-submissions' },
        { path: '/leaderboard', label: 'Clasificación', icon: <BarChart3 size={20} />, id: 'tour-nav-leaderboard' },
    ] : [
        { path: '/dashboard', label: 'PANEL PRINCIPAL', icon: <LayoutDashboard size={20} />, id: 'tour-nav-dashboard' },
        { path: '/challenges', label: 'Retos', icon: <Trophy size={20} />, id: 'tour-nav-challenges' },
        { path: '/challenges/create', label: 'Crear Reto', icon: <PlusCircle size={20} />, id: 'tour-nav-challenges-create' },
        { path: '/courses', label: 'Cursos', icon: <BookOpen size={20} />, id: 'tour-nav-courses' },
        { path: '/submissions', label: 'Envíos', icon: <FileText size={20} />, id: 'tour-nav-submissions' },
        { path: '/leaderboard', label: 'Clasificación', icon: <BarChart3 size={20} />, id: 'tour-nav-leaderboard' },
    ];

    return (
        <aside className="sidebar">
            <div className="sidebar-header">
                <Link to={user ? "/dashboard" : "/"} className="sidebar-brand-minimal">
                    <div className="brand-logo">
                        <img src="/logo.png" alt="RobleCode" />
                    </div>
                </Link>
            </div>

            <nav className="sidebar-nav">
                <div className="nav-section">
                    <ul className="nav-list">
                        {navItems.map((item) => (
                            <li key={item.path}>
                                <Link 
                                    to={item.path} 
                                    className={`nav-link ${isActive(item.path) ? 'active' : ''}`}
                                    id={item.id}
                                >
                                    <span className="nav-icon">{item.icon}</span>
                                    <span className="nav-label">{item.label}</span>
                                    {isActive(item.path) && <ChevronRight size={14} className="active-arrow" />}
                                </Link>
                            </li>
                        ))}
                    </ul>
                </div>

                {user?.role === 'admin' && (
                    <div className="nav-section">
                        <span className="section-title">Administración</span>
                        <ul className="nav-list">
                            <li>
                                <Link to="/settings" className="nav-link">
                                    <span className="nav-icon"><Settings size={20} /></span>
                                    <span className="nav-label">Configuración</span>
                                </Link>
                            </li>
                        </ul>
                    </div>
                )}
            </nav>

            <div className="sidebar-footer">
                <div className="user-card">
                    <div className="user-avatar">
                        <div className="avatar-placeholder">
                            {user?.username?.charAt(0).toUpperCase() || 'U'}
                        </div>
                    </div>
                    <div className="user-details">
                        <span className="username" title={user?.username || 'Usuario'}>
                            {user?.username || 'Usuario'}
                        </span>
                        <div className={`role-badge ${user?.role || 'student'}`}>
                            <ShieldCheck size={10} />
                            <span>
                                {user?.role === 'professor' || user?.role === 'teacher' ? 'Profesor' : 
                                 user?.role === 'admin' ? 'Administrador' : 'Estudiante'}
                            </span>
                        </div>
                    </div>
                </div>
                <button onClick={() => {
                    import('../tours/tourConfig').then(({ startTour }) => {
                        startTour(user?.role);
                    });
                }} className="logout-button" style={{marginBottom: '10px', background: 'transparent', border: '1px solid rgba(255,255,255,0.2)', opacity: 0.8}}>
                    <BookOpen size={18} />
                    <span>Hacer Tour</span>
                </button>
                <button onClick={logout} className="logout-button">
                    <LogOut size={18} />
                    <span>Cerrar Sesión</span>
                </button>
            </div>
        </aside>
    );
};

export default Sidebar;
