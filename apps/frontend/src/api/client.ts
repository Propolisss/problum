import axios, { AxiosError, type AxiosInstance } from 'axios';
import { getAccessToken, setAccessToken, clearAccessToken } from '../features/auth/token';

type FailedReq = {
  resolve: (value?: any) => void;
  reject: (err: any) => void;
  config: any;
};

const api: AxiosInstance = axios.create({
  baseURL: '/api',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

let isRefreshing = false;
let failedQueue: FailedReq[] = [];
let onLogout: (() => void) | null = null;

export function setOnLogout(cb: () => void) {
  onLogout = cb;
}

function processQueue(error: any, token: string | null = null) {
  failedQueue.forEach((p) => {
    if (error) {
      p.reject(error);
    } else {
      p.config.headers['Authorization'] = `Bearer ${token}`;
      p.resolve(api(p.config));
    }
  });
  failedQueue = [];
}

api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (r) => r,
  async (error: AxiosError & { config?: any }) => {
    const originalRequest = error.config;
    if (!originalRequest) return Promise.reject(error);

    const status = error.response?.status;
    if (status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject, config: originalRequest });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const resp = await axios.post('/auth/refresh', null, { withCredentials: true });
        const newToken = resp.data?.access_token;
        if (!newToken) throw new Error('no access token on refresh');

        setAccessToken(newToken);
        processQueue(null, newToken);

        originalRequest.headers = originalRequest.headers ?? {};
        originalRequest.headers['Authorization'] = `Bearer ${newToken}`;

        return api(originalRequest);
      } catch (e) {
        processQueue(e, null);
        clearAccessToken();
        if (onLogout) onLogout();
        return Promise.reject(e);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  },
);

export default api;
