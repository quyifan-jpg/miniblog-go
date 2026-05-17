import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import api, { TOKEN_STORAGE_KEY } from '../services/api';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
   const [user, setUser] = useState(null);
   const [token, setToken] = useState(() => localStorage.getItem(TOKEN_STORAGE_KEY));
   const [loading, setLoading] = useState(Boolean(localStorage.getItem(TOKEN_STORAGE_KEY)));

   const persistToken = useCallback(nextToken => {
      if (nextToken) {
         localStorage.setItem(TOKEN_STORAGE_KEY, nextToken);
      } else {
         localStorage.removeItem(TOKEN_STORAGE_KEY);
      }
      setToken(nextToken);
   }, []);

   const fetchMe = useCallback(async () => {
      try {
         const res = await api.auth.me();
         setUser(res.data);
         return res.data;
      } catch (err) {
         setUser(null);
         persistToken(null);
         throw err;
      }
   }, [persistToken]);

   useEffect(() => {
      if (!token) {
         setLoading(false);
         return;
      }
      let cancelled = false;
      (async () => {
         try {
            await fetchMe();
         } catch {
            // interceptor will redirect on 401; nothing else to do
         } finally {
            if (!cancelled) setLoading(false);
         }
      })();
      return () => {
         cancelled = true;
      };
   }, [token, fetchMe]);

   const login = async ({ email, password }) => {
      const res = await api.auth.login({ email, password });
      persistToken(res.data.access_token);
      const me = await api.auth.me();
      setUser(me.data);
      return me.data;
   };

   const register = async ({ email, username, password }) => {
      const res = await api.auth.register({ email, username, password });
      persistToken(res.data.access_token);
      const me = await api.auth.me();
      setUser(me.data);
      return me.data;
   };

   const logout = () => {
      persistToken(null);
      setUser(null);
   };

   const value = {
      user,
      token,
      isAuthenticated: Boolean(token && user),
      loading,
      login,
      register,
      logout,
   };

   return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
   const ctx = useContext(AuthContext);
   if (!ctx) {
      throw new Error('useAuth must be used inside <AuthProvider>');
   }
   return ctx;
};
