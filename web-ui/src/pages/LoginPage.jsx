import React, { useState } from 'react';
import { HiShieldCheck } from 'react-icons/hi';
import api from '../services/api';

const LoginPage = ({ onLogin }) => {
  const [token, setToken] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e) => {
    e.preventDefault();
    if (!token.trim()) {
      alert('Please enter a token');
      return;
    }

    setLoading(true);
    try {
      api.setToken(token);
      await api.getStatus();
      onLogin();
    } catch (error) {
      api.clearToken();
      alert('Invalid token: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-vault-dark-bg px-8">
      <div className="w-full max-w-md bg-vault-dark-card p-10 rounded border border-vault-border-dark text-center">
        <div className="flex justify-center mb-6">
          <HiShieldCheck className="w-16 h-16 text-white" />
        </div>
        <h1 className="text-2xl font-semibold text-white mb-3">Sign in to Vault</h1>
        <p className="text-base text-vault-text-muted mb-8 leading-relaxed">
          Enter your authentication token to access the vault.
        </p>

        <form onSubmit={handleLogin} className="text-left">
          <div className="mb-6">
            <label htmlFor="token" className="block text-sm font-medium text-vault-text-on-dark mb-2">
              Token
            </label>
            <input
              id="token"
              type="password"
              className="w-full px-4 py-3 text-sm bg-vault-dark-bg border border-vault-border-dark text-white rounded focus:outline-none focus:border-vault-primary focus:ring-2 focus:ring-vault-primary/20 transition-all"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder="Enter your vault token"
              disabled={loading}
            />
          </div>

          <button
            type="submit"
            className="w-full px-6 py-3 text-base font-medium text-white bg-vault-primary rounded hover:bg-vault-primary-hover transition-colors disabled:opacity-50"
            disabled={loading}
          >
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default LoginPage;
