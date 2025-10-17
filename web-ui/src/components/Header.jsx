import React, { useState, useEffect } from 'react';
import { HiShieldCheck, HiLockClosed, HiLogout } from 'react-icons/hi';
import api from '../services/api';

const Header = ({ onLogout }) => {
  const [status, setStatus] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  const fetchStatus = async () => {
    try {
      const data = await api.getStatus();
      setStatus(data);
    } catch (error) {
      console.error('Failed to fetch status:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSeal = async () => {
    if (window.confirm('Are you sure you want to seal the vault?')) {
      try {
        await api.seal();
        await fetchStatus();
      } catch (error) {
        alert('Failed to seal vault: ' + error.message);
      }
    }
  };

  return (
    <header className="h-14 bg-vault-dark-sidebar border-b border-vault-border-dark flex items-center justify-between px-6">
      <div className="flex items-center gap-6">
        <div className="flex items-center gap-2.5 pr-6 border-r border-vault-border-dark">
          <HiShieldCheck className="w-6 h-6 text-white" />
          <span className="text-white font-medium text-base">Vault</span>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-vault-text-on-dark text-sm px-2 py-1 rounded hover:bg-white/10 transition-colors">
            secrets
          </span>
        </div>
      </div>
      <div className="flex items-center gap-3">
        {!loading && status && (
          <div className="flex items-center gap-2">
            <span className={`px-3 py-1.5 rounded text-sm font-medium border ${
              status.sealed
                ? 'bg-vault-danger/15 text-vault-danger border-vault-danger'
                : 'bg-vault-success/15 text-vault-success border-vault-success'
            }`}>
              {status.sealed ? 'Sealed' : 'Unsealed'}
            </span>
          </div>
        )}
        {!status?.sealed && (
          <button
            onClick={handleSeal}
            className="flex items-center gap-2 px-3.5 py-2 text-sm font-medium text-white bg-white/10 border border-vault-border-dark rounded hover:bg-white/15 transition-colors"
          >
            <HiLockClosed className="w-4 h-4" />
            Seal
          </button>
        )}
        {onLogout && (
          <button
            onClick={onLogout}
            className="flex items-center gap-2 px-3.5 py-2 text-sm font-medium text-white bg-white/10 border border-vault-border-dark rounded hover:bg-white/15 transition-colors"
          >
            <HiLogout className="w-4 h-4" />
            Sign out
          </button>
        )}
      </div>
    </header>
  );
};

export default Header;
