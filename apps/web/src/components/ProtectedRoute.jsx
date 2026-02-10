import { useAuth } from '../context/AuthContext';
import { Navigate, useLocation } from 'react-router-dom';

/**
 * ProtectedRoute ensures that the wrapped component is only rendered for
 * authenticated users. If no user is present, it redirects to the login page.
 */
const ProtectedRoute = ({ children }) => {
    const { user, loading } = useAuth();
    const location = useLocation();

    if (loading) {
        // You could render a spinner here; keep it simple for now.
        return null;
    }

    if (!user) {
        // Preserve the attempted location so we can redirect back after login.
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    return children;
};

export default ProtectedRoute;
