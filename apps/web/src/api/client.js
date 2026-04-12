import axios from 'axios';

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

const client = axios.create({
  baseURL,
});

let refreshPromise = null;

const clearAuthStorage = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('refresh_token');
  localStorage.removeItem('user_email');
  localStorage.removeItem('session_id');
  window.dispatchEvent(new Event('auth:logout'));
};

export const refreshAccessToken = async () => {
  if (refreshPromise) {
    return refreshPromise;
  }

  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) {
    throw new Error('Refresh token not found');
  }

  refreshPromise = axios
    .post(`${baseURL}/auth/refresh-token`, { refresh_token: refreshToken })
    .then(({ data }) => {
      const newAccessToken = data?.access_token;
      const newRefreshToken = data?.refresh_token;

      if (!newAccessToken) {
        throw new Error('Refresh did not return access token');
      }

      localStorage.setItem('token', newAccessToken);
      if (newRefreshToken) {
        localStorage.setItem('refresh_token', newRefreshToken);
      }

      return newAccessToken;
    })
    .finally(() => {
      refreshPromise = null;
    });

  return refreshPromise;
};

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

  client.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error?.config;
      const status = error?.response?.status;

      if (!originalRequest || status !== 401 || originalRequest._retry) {
        return Promise.reject(error);
      }

      const requestUrl = String(originalRequest.url || '');
      if (
        requestUrl.includes('/auth/login') ||
        requestUrl.includes('/auth/register') ||
        requestUrl.includes('/auth/refresh-token')
      ) {
        return Promise.reject(error);
      }

      try {
        originalRequest._retry = true;
        const newAccessToken = await refreshAccessToken();
        originalRequest.headers = originalRequest.headers || {};
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        return client(originalRequest);
      } catch (refreshError) {
        clearAuthStorage();
        return Promise.reject(refreshError);
      }
    }
  );

export default client;
