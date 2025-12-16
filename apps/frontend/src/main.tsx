import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import {
  QueryClient,
  QueryClientProvider,
  QueryCache,
  MutationCache
} from '@tanstack/react-query';
import { Toaster, toast } from 'react-hot-toast';
import App from './App';
import { AuthProvider } from './features/auth/AuthContext';
import './index.css';

const handleGlobalError = (error: any) => {
  if (error?.response?.status === 401) {
    return;
  }

  toast.error('Произошла ошибка. Попробуйте позже.', {
    duration: 4000,
    position: 'top-right',
    style: {
      background: '#333',
      color: '#fff',
    },
  });
};

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: handleGlobalError,
  }),
  mutationCache: new MutationCache({
    onError: handleGlobalError,
  }),
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: false,
    },
  },
});

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <App />
          <Toaster />
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  </React.StrictMode>,
);
