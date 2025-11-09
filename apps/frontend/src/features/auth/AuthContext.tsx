import React, { createContext, useContext, useEffect, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import api, { setOnLogout } from '../../api/client';
import axios from 'axios';
import { setAccessToken as setToken, clearAccessToken } from './token';

type AuthContextValue = {
    accessToken: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;

    login: (login: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
    tryRefresh: () => Promise<boolean>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [accessToken, setAccessTokenState] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const queryClient = useQueryClient();

    useEffect(() => {
        setToken(accessToken);
    }, [accessToken]);

    useEffect(() => {
        setOnLogout(() => {
            setAccessTokenState(null);
        });
    }, []);

    const tryRefresh = async (): Promise<boolean> => {
        try {
            const resp = await axios.post('/api/auth/refresh', null, { withCredentials: true });
            const token = resp.data?.access_token;
            if (!token) throw new Error('no token');
            setAccessTokenState(token);
            return true;
        } catch (e) {
            setAccessTokenState(null);
            return false;
        }
    };

    const login = async (loginVal: string, password: string) => {
        const resp = await api.post('/auth/login', { login: loginVal, password });
        const token = resp.data?.access_token;
        if (!token) throw new Error('no access token');
        setAccessTokenState(token);
    };

    const logout = async () => {
        try {
            await api.post('/auth/logout', null);
        } catch (e) {
        }
        setAccessTokenState(null);
        clearAccessToken();
        queryClient.clear();
    };

    const value: AuthContextValue = {
        accessToken,
        isAuthenticated: !!accessToken,
        isLoading,
        login,
        logout,
        tryRefresh,
    };

    useEffect(() => {
        const initializeAuth = async () => {
            try {
                await tryRefresh();
            } catch (e) {
                console.error('Auth initialization failed', e);
            } finally {
                setIsLoading(false);
            }
        };
        initializeAuth();
    }, []);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export function useAuthContext() {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error('useAuthContext must be used inside AuthProvider');
    return ctx;
}
