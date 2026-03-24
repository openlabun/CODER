import { createContext, useState, useEffect, useContext } from 'react';
import client from '../api/client';

export const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const checkAuth = async () => {
            const token = localStorage.getItem('token');
            const email = localStorage.getItem('user_email');
            if (token && email) {
                try {
                    const { data } = await client.get('/auth/me', {
                        headers: { 'X-User-Email': email }
                    });
                    setUser({ 
                        id: data.ID || data.UserID || email,
                        token, 
                        username: data.Username || data.Email, 
                        role: String(data.Role).toLowerCase() || 'student', 
                        email: data.Email 
                    });
                } catch (error) {
                    localStorage.removeItem('token');
                    localStorage.removeItem('user_email');
                }
            }
            setLoading(false);
        };
        checkAuth();
    }, []);

    const login = async (email, password) => {
        try {
            const { data } = await client.post('/auth/login', { email, password });
            
            // Go API structure: data.token.access_token and data.user_data
            const token = data.token?.access_token;
            const userData = data.user_data;

            if (!token) throw new Error('Token no recibido del servidor');

            localStorage.setItem('token', token);
            localStorage.setItem('user_email', email);

            // Fetch user details including role using the required header
            const { data: profile } = await client.get('/auth/me', {
                headers: { 'X-User-Email': email }
            });
            
            setUser({ 
                id: profile.ID || profile.UserID || userData?.ID || email,
                token: token, 
                username: profile.Username || profile.Email || userData?.Username, 
                role: String(profile.Role || userData?.Role).toLowerCase() || 'student',
                email: email
            });
        } catch (error) {
            if (!error.response) {
                throw new Error('Error de conexión: No se pudo contactar con el servidor.');
            }
            const data = error.response.data;
            const message = data?.error || data?.message || 'Error al iniciar sesión. Verifica tus credenciales.';
            throw new Error(message);
        }
    };

    const register = async (name, email, password) => {
        try {
            const { data } = await client.post('/auth/register', { name, email, password });
            
            const token = data.token?.access_token;
            const userData = data.user_data;
            if (!token) throw new Error('Token no recibido tras el registro');

            localStorage.setItem('token', token);
            localStorage.setItem('user_email', email);
            
            setUser({ 
                id: userData?.ID || userData?.UserID || email,
                token, 
                username: userData?.Username || name, 
                role: String(userData?.Role || 'student').toLowerCase(), 
                email: email 
            });
        } catch (error) {
            if (!error.response) {
                throw new Error('Error de conexión: No se pudo contactar con el servidor.');
            }
            const data = error.response.data;
            const message = data?.error || data?.message || 'Error al registrar el usuario.';
            throw new Error(message);
        }
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user_email');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, register, logout, loading }}>
            {!loading && children}
        </AuthContext.Provider>
    );
};
