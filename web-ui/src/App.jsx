import React, { useState, useEffect } from 'react';
import api from './services/api';
import Layout from './components/Layout';
import InitPage from './pages/InitPage';
import UnsealPage from './pages/UnsealPage';
import LoginPage from './pages/LoginPage';
import SecretsPage from './pages/SecretsPage';

function App() {
  const [vaultStatus, setVaultStatus] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    checkVaultStatus();
    const interval = setInterval(checkVaultStatus, 10000);
    return () => clearInterval(interval);
  }, []);

  const checkVaultStatus = async () => {
    try {
      const status = await api.getStatus();
      setVaultStatus(status);

      // Check if we have a valid token
      const token = api.getToken();
      if (token && !status.sealed) {
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(false);
      }
    } catch (error) {
      console.error('Failed to check vault status:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    api.clearToken();
    setIsAuthenticated(false);
  };

  const handleLogin = () => {
    setIsAuthenticated(true);
    checkVaultStatus();
  };

  const handleInitialized = () => {
    checkVaultStatus();
  };

  const handleUnsealed = () => {
    checkVaultStatus();
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-vault-light-bg">
        <div className="w-12 h-12 border-4 border-vault-border-light border-t-vault-primary rounded-full animate-spin"></div>
        <p className="mt-4 text-base text-vault-text-secondary">Loading Vault...</p>
      </div>
    );
  }

  // Vault not initialized
  if (!vaultStatus?.initialized) {
    return <InitPage onInitialized={handleInitialized} />;
  }

  // Vault is sealed
  if (vaultStatus?.sealed) {
    return <UnsealPage onUnsealed={handleUnsealed} />;
  }

  // Not authenticated
  if (!isAuthenticated) {
    return <LoginPage onLogin={handleLogin} />;
  }

  // Main application
  return (
    <Layout onLogout={handleLogout}>
      <SecretsPage />
    </Layout>
  );
}

export default App;
