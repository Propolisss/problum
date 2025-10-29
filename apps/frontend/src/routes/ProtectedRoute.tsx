import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../features/auth/hooks';

function LoadingScreen() {
    return (
        <div className="flex items-center justify-center h-screen bg-gray-100">
            <div className="text-lg text-gray-600">Загрузка платформы...</div>
        </div>
    )
}

export const ProtectedRoute: React.FC<{ children: JSX.Element }> = ({ children }) => {
    const { isAuthenticated, isLoading } = useAuth();
    const location = useLocation();

    if (isLoading) {
        return <LoadingScreen />;
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" state={{ from: location }} replace />;
    }

    return children;
};
