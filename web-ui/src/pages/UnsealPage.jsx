import React, { useState } from 'react';
import { HiShieldCheck } from 'react-icons/hi';
import api from '../services/api';

const UnsealPage = ({ onUnsealed }) => {
  const [unsealKey, setUnsealKey] = useState('');
  const [loading, setLoading] = useState(false);

  const handleUnseal = async (e) => {
    e.preventDefault();
    if (!unsealKey.trim()) {
      alert('Please enter the unseal key');
      return;
    }

    setLoading(true);
    try {
      const data = await api.unseal(unsealKey);
      if (!data.sealed) {
        onUnsealed();
      } else {
        alert('Vault is still sealed. Please check your unseal key.');
      }
    } catch (error) {
      alert('Failed to unseal vault: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-vault-dark-bg px-8">
      <div className="w-full max-w-lg bg-vault-dark-card p-10 rounded border border-vault-border-dark text-center">
        <div className="flex justify-center mb-6">
          <HiShieldCheck className="w-16 h-16 text-white" />
        </div>
        <h1 className="text-2xl font-semibold text-white mb-3">Unseal Vault</h1>
        <p className="text-base text-vault-text-muted mb-8 leading-relaxed">
          The vault is currently sealed. Enter your unseal key to unlock it.
        </p>

        <form onSubmit={handleUnseal} className="text-left">
          <div className="mb-6">
            <label htmlFor="unseal-key" className="block text-sm font-medium text-vault-text-on-dark mb-2">
              Unseal Key
            </label>
            <input
              id="unseal-key"
              type="password"
              className="w-full px-4 py-3 text-sm bg-vault-dark-bg border border-vault-border-dark text-white rounded focus:outline-none focus:border-vault-primary focus:ring-2 focus:ring-vault-primary/20 transition-all"
              value={unsealKey}
              onChange={(e) => setUnsealKey(e.target.value)}
              placeholder="Enter your unseal key"
              disabled={loading}
            />
          </div>

          <button
            type="submit"
            className="w-full px-6 py-3 text-base font-medium text-white bg-vault-primary rounded hover:bg-vault-primary-hover transition-colors disabled:opacity-50"
            disabled={loading}
          >
            {loading ? 'Unsealing...' : 'Unseal Vault'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default UnsealPage;
