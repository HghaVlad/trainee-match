import axios from 'axios';

export const api = axios.create({
  baseURL: 'http://localhost:8000/api/v1',
  withCredentials: true,
});

api.interceptors.response.use(
  (res) => res,
  async (error) => {
    if (error.response?.status === 401) {
      try {
        await api.post('/auth/refresh');
        return api.request(error.config);
      } catch {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);
