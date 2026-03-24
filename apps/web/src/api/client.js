import axios from 'axios';

const client = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:3000',
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  const email = localStorage.getItem('user_email');
  
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  if (email) {
    config.headers['X-User-Email'] = email;
  }
  
  return config;
});

export default client;
