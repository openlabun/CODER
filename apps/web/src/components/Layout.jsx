import { Outlet, useLocation } from 'react-router-dom';
import Sidebar from './Sidebar';
import Navbar from './Navbar'; // Keep Navbar for mobile or top actions if needed
import './Layout.css';

const Layout = () => {
    const location = useLocation();
    // Hide sidebar on public pages if desired, or keep it consistent
    const isPublic = ['/', '/login', '/register'].includes(location.pathname);
    const showNavbar = location.pathname === '/';

    return (
        <div className="app-layout">
            {!isPublic && <Sidebar />}
            <div className={`main-content ${!isPublic ? 'with-sidebar' : ''}`}>
                {showNavbar && <Navbar />} 
                <main>
                    <Outlet />
                </main>
            </div>
        </div>
    );
};

export default Layout;
