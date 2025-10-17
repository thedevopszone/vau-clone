import React, { useState } from 'react';
import { HiShieldCheck } from 'react-icons/hi';
import api from '../services/api';

const InitPage = ({ onInitialized }) => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);

  const handleInitialize = async () => {
    setLoading(true);
    try {
      const data = await api.initialize();
      setResult(data);
    } catch (error) {
      alert('Failed to initialize vault: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleContinue = () => {
    if (result?.root_token) {
      api.setToken(result.root_token);
      onInitialized();
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-vault-dark-bg px-8">
      <div className="w-full max-w-2xl bg-vault-dark-card p-10 rounded border border-vault-border-dark text-center">
        <div className="flex justify-center mb-6 text-white">
          <HiShieldCheck className="w-16 h-16" />
        </div>
        <h1 className="text-2xl font-semibold text-white mb-3">Initialize Vault</h1>
        <p className="text-base text-vault-text-muted mb-8 leading-relaxed">
          Initialize your vault to generate the root token and unseal key.
          These credentials will only be shown once, so make sure to save them securely.
        </p>

        {!result ? (
          <button
            className="w-full max-w-md mx-auto px-6 py-3 text-base font-medium text-white bg-vault-primary rounded hover:bg-vault-primary-hover transition-colors disabled:opacity-50"
            onClick={handleInitialize}
            disabled={loading}
          >
            {loading ? 'Initializing...' : 'Initialize Vault'}
          </button>
        ) : (
          <div className="text-left">
            <div className="mb-6 p-4 bg-vault-success/15 border border-vault-success rounded">
              <h3 className="text-lg font-semibold text-vault-success mb-2">Vault Initialized Successfully!</h3>
              <p className="text-sm text-vault-text-on-dark">Save these credentials securely. They will not be shown again.</p>
            </div>

            <div className="space-y-4 mb-8">
              <div className="p-4 bg-vault-dark-bg border border-vault-border-dark rounded">
                <label className="block text-sm font-medium text-vault-text-on-dark mb-2">Root Token:</label>
                <div className="flex gap-2">
                  <code className="flex-1 p-3 bg-black/30 border border-vault-border-dark text-white rounded text-sm break-all">{result.root_token}</code>
                  <button
                    className="px-4 py-2 text-sm font-medium text-white bg-white/10 border border-vault-border-dark rounded hover:bg-white/15 transition-colors"
                    onClick={() => navigator.clipboard.writeText(result.root_token)}
                  >
                    Copy
                  </button>
                </div>
              </div>

              <div className="p-4 bg-vault-dark-bg border border-vault-border-dark rounded">
                <label className="block text-sm font-medium text-vault-text-on-dark mb-2">Unseal Key:</label>
                <div className="flex gap-2">
                  <code className="flex-1 p-3 bg-black/30 border border-vault-border-dark text-white rounded text-sm break-all">{result.unseal_key}</code>
                  <button
                    className="px-4 py-2 text-sm font-medium text-white bg-white/10 border border-vault-border-dark rounded hover:bg-white/15 transition-colors"
                    onClick={() => navigator.clipboard.writeText(result.unseal_key)}
                  >
                    Copy
                  </button>
                </div>
              </div>
            </div>

            <button className="w-full px-6 py-3 text-base font-medium text-white bg-vault-primary rounded hover:bg-vault-primary-hover transition-colors" onClick={handleContinue}>
              Continue to Unseal
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default InitPage;
