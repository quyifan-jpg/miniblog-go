import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Register = () => {
   const { register } = useAuth();
   const navigate = useNavigate();

   const [email, setEmail] = useState('');
   const [username, setUsername] = useState('');
   const [password, setPassword] = useState('');
   const [submitting, setSubmitting] = useState(false);
   const [error, setError] = useState(null);

   const handleSubmit = async e => {
      e.preventDefault();
      setError(null);
      if (password.length < 8) {
         setError('Password must be at least 8 characters');
         return;
      }
      setSubmitting(true);
      try {
         await register({ email, username, password });
         navigate('/studio', { replace: true });
      } catch (err) {
         const msg =
            err?.response?.data?.detail ||
            err?.response?.data?.message ||
            err?.message ||
            'Registration failed';
         setError(msg);
      } finally {
         setSubmitting(false);
      }
   };

   return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-b from-gray-900 to-black p-4">
         <div className="w-full max-w-md bg-gradient-to-br from-gray-800 to-gray-900 border border-gray-700 rounded-md shadow-2xl p-8">
            <h1 className="text-2xl font-bold text-white mb-2 text-center">Create account</h1>
            <p className="text-gray-400 text-sm mb-6 text-center">
               Start building podcasts in minutes
            </p>
            {error && (
               <div className="mb-4 p-3 bg-red-900/30 border border-red-800/50 rounded text-red-400 text-sm">
                  {error}
               </div>
            )}
            <form onSubmit={handleSubmit} className="space-y-4">
               <div>
                  <label className="block text-sm text-gray-300 mb-1">Email</label>
                  <input
                     type="email"
                     required
                     value={email}
                     onChange={e => setEmail(e.target.value)}
                     className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white focus:outline-none focus:border-emerald-600"
                     autoComplete="email"
                  />
               </div>
               <div>
                  <label className="block text-sm text-gray-300 mb-1">Username</label>
                  <input
                     type="text"
                     required
                     minLength={3}
                     maxLength={50}
                     value={username}
                     onChange={e => setUsername(e.target.value)}
                     className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white focus:outline-none focus:border-emerald-600"
                     autoComplete="username"
                  />
                  <p className="mt-1 text-xs text-gray-500">
                     Letters, digits, _ and - only. 3–50 chars.
                  </p>
               </div>
               <div>
                  <label className="block text-sm text-gray-300 mb-1">Password</label>
                  <input
                     type="password"
                     required
                     minLength={8}
                     value={password}
                     onChange={e => setPassword(e.target.value)}
                     className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white focus:outline-none focus:border-emerald-600"
                     autoComplete="new-password"
                  />
                  <p className="mt-1 text-xs text-gray-500">Minimum 8 characters.</p>
               </div>
               <button
                  type="submit"
                  disabled={submitting}
                  className={`w-full py-2.5 bg-gradient-to-r from-emerald-700 to-emerald-800 hover:from-emerald-600 hover:to-emerald-700 text-white font-medium rounded transition ${
                     submitting ? 'opacity-70 cursor-not-allowed' : ''
                  }`}
               >
                  {submitting ? 'Creating…' : 'Create account'}
               </button>
            </form>
            <p className="mt-6 text-sm text-gray-400 text-center">
               Already have an account?{' '}
               <Link to="/login" className="text-emerald-400 hover:text-emerald-300">
                  Sign in
               </Link>
            </p>
         </div>
      </div>
   );
};

export default Register;
