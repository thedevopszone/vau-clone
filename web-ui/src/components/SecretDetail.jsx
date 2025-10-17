import React, { useState } from 'react';
import { HiEye, HiEyeOff, HiClipboardCopy, HiPencil, HiTrash } from 'react-icons/hi';

const SecretDetail = ({ secret, onEdit, onDelete }) => {
  const [showValues, setShowValues] = useState({});

  const toggleValue = (key) => {
    setShowValues((prev) => ({
      ...prev,
      [key]: !prev[key],
    }));
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  return (
    <div className="p-8">
      <div className="flex justify-between items-start mb-8 pb-6 border-b border-vault-border-light">
        <div>
          <h2 className="text-xl font-semibold text-vault-text-primary mb-2 break-all">{secret.path}</h2>
          <p className="text-sm text-vault-text-secondary">
            Version: {secret.version || 1} | Created: {new Date().toLocaleDateString()}
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={onEdit}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-vault-text-primary bg-white border border-vault-border-light rounded hover:bg-vault-light-bg transition-colors"
          >
            <HiPencil className="w-4 h-4" />
            Edit
          </button>
          <button
            onClick={onDelete}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-vault-danger rounded hover:bg-red-600 transition-colors"
          >
            <HiTrash className="w-4 h-4" />
            Delete
          </button>
        </div>
      </div>

      <div>
        <h3 className="text-xs font-bold uppercase tracking-wider text-vault-text-secondary mb-4">Secret Data</h3>
        {Object.entries(secret.data || {}).map(([key, value]) => (
          <div key={key} className="mb-4 p-4 bg-vault-light-bg border border-vault-border-light rounded">
            <div className="flex justify-between items-center mb-2">
              <label className="text-sm font-semibold text-vault-text-primary">{key}</label>
              <div className="flex gap-1">
                <button
                  onClick={() => toggleValue(key)}
                  className="p-1.5 text-vault-text-secondary hover:text-vault-text-primary transition-colors"
                  title={showValues[key] ? 'Hide' : 'Show'}
                >
                  {showValues[key] ? <HiEyeOff className="w-5 h-5" /> : <HiEye className="w-5 h-5" />}
                </button>
                <button
                  onClick={() => copyToClipboard(value)}
                  className="p-1.5 text-vault-text-secondary hover:text-vault-text-primary transition-colors"
                  title="Copy"
                >
                  <HiClipboardCopy className="w-5 h-5" />
                </button>
              </div>
            </div>
            <div className="p-3 bg-white border border-vault-border-light rounded">
              <code className="text-sm text-vault-text-primary break-all">
                {showValues[key] ? value : '••••••••'}
              </code>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SecretDetail;
