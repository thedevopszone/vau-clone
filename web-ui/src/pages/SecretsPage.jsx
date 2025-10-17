import React, { useState, useEffect } from 'react';
import { HiPlus } from 'react-icons/hi';
import api from '../services/api';
import SecretsList from '../components/SecretsList';
import SecretDetail from '../components/SecretDetail';
import SecretForm from '../components/SecretForm';

const SecretsPage = () => {
  const [secrets, setSecrets] = useState([]);
  const [selectedSecret, setSelectedSecret] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadSecrets();
  }, []);

  const loadSecrets = async () => {
    setLoading(true);
    try {
      const data = await api.listSecrets();
      setSecrets(data.keys || data.secrets || []);
    } catch (error) {
      console.error('Failed to load secrets:', error);
      setSecrets([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectSecret = async (path) => {
    try {
      const data = await api.readSecret(path);
      setSelectedSecret({ path, ...data });
      setShowForm(false);
    } catch (error) {
      alert('Failed to read secret: ' + error.message);
    }
  };

  const handleDeleteSecret = async (path) => {
    if (window.confirm(`Are you sure you want to delete "${path}"?`)) {
      try {
        await api.deleteSecret(path);
        await loadSecrets();
        if (selectedSecret?.path === path) {
          setSelectedSecret(null);
        }
      } catch (error) {
        alert('Failed to delete secret: ' + error.message);
      }
    }
  };

  const handleCreateNew = () => {
    setSelectedSecret(null);
    setShowForm(true);
  };

  const handleSaveSecret = async (path, data) => {
    try {
      await api.writeSecret(path, data);
      await loadSecrets();
      setShowForm(false);
      handleSelectSecret(path);
    } catch (error) {
      alert('Failed to save secret: ' + error.message);
    }
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold text-vault-text-primary">Secrets</h1>
        <button
          onClick={handleCreateNew}
          className="flex items-center gap-2 px-4 py-2.5 text-sm font-medium text-white bg-vault-primary rounded hover:bg-vault-primary-hover transition-colors"
        >
          <HiPlus className="w-5 h-5" />
          Create Secret
        </button>
      </div>

      <div className="flex gap-6 flex-1 overflow-hidden">
        <SecretsList
          secrets={secrets}
          selectedPath={selectedSecret?.path}
          onSelect={handleSelectSecret}
          onDelete={handleDeleteSecret}
          loading={loading}
        />

        <div className="flex-1 bg-white border border-vault-border-light rounded overflow-y-auto">
          {showForm ? (
            <SecretForm
              onSave={handleSaveSecret}
              onCancel={() => setShowForm(false)}
            />
          ) : selectedSecret ? (
            <SecretDetail
              secret={selectedSecret}
              onEdit={() => setShowForm(true)}
              onDelete={() => handleDeleteSecret(selectedSecret.path)}
            />
          ) : (
            <div className="flex items-center justify-center h-full text-base text-vault-text-secondary p-8 text-center">
              <p>Select a secret to view its details or create a new one.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default SecretsPage;
