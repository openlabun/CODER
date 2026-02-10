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
            if (token) {
                try {
                    const { data } = await client.get('/auth/me');
                    setUser({ token, username: data.username, role: data.role });
                } catch (error) {
                    localStorage.removeItem('token');
                }
            }
            setLoading(false);
        };
        checkAuth();
    }, []);

    const login = async (username, password) => {
        const { data } = await client.post('/auth/login', { username, password });
        localStorage.setItem('token', data.accessToken);

        // Fetch user details including role
        const userResponse = await client.get('/auth/me');
        setUser({ token: data.accessToken, username: userResponse.data.username, role: userResponse.data.role });
    };

    const register = async (username, password, role) => {
        const { data } = await client.post('/auth/register', { username, password, role });
        localStorage.setItem('token', data.accessToken);
        setUser({ token: data.accessToken, username, role });
    };

    const logout = () => {
        localStorage.removeItem('token');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, register, logout, loading }}>
            {!loading && children}
        </AuthContext.Provider>
    );
};
