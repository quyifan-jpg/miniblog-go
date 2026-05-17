import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const ProtectedRoute = ({ children }) => {
   const { token, loading } = useAuth();
   const location = useLocation();

   if (loading) {
      return (
         <div className="min-h-screen flex items-center justify-center bg-gray-900 text-gray-400">
            Loading…
         </div>
      );
   }

   if (!token) {
      const redirect = encodeURIComponent(location.pathname + location.search);
      return <Navigate to={`/login?redirect=${redirect}`} replace />;
   }

   return children;
};

export default ProtectedRoute;
