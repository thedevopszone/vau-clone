import React from 'react';
import { HiKey, HiTrash } from 'react-icons/hi';

const SecretsList = ({ secrets, selectedPath, onSelect, onDelete, loading }) => {
  if (loading) {
    return (
      <div className="w-72 bg-white border border-vault-border-light rounded p-8 text-center">
        <div className="text-base text-vault-text-secondary">Loading secrets...</div>
      </div>
    );
  }

  if (!secrets || secrets.length === 0) {
    return (
      <div className="w-72 bg-white border border-vault-border-light rounded p-8 text-center">
        <div className="text-base text-vault-text-secondary">No secrets found</div>
      </div>
    );
  }

  return (
    <div className="w-72 bg-white border border-vault-border-light rounded p-2 overflow-y-auto">
      {secrets.map((secret) => (
        <div
          key={secret}
          className={`
            flex items-center justify-between px-4 py-2.5 mb-0.5 rounded cursor-pointer transition-all border
            ${selectedPath === secret
              ? 'bg-vault-primary/10 text-vault-primary border-vault-primary font-medium'
              : 'border-transparent hover:bg-vault-light-bg'
            }
          `}
        >
          <div className="flex items-center gap-2 flex-1" onClick={() => onSelect(secret)}>
            <HiKey className={`w-5 h-5 ${selectedPath === secret ? 'opacity-100' : 'opacity-60'}`} />
            <span className="text-sm truncate">{secret}</span>
          </div>
          <button
            className="p-1 opacity-0 hover:opacity-100 group-hover:opacity-50 transition-opacity"
            onClick={(e) => {
              e.stopPropagation();
              onDelete(secret);
            }}
            title="Delete secret"
          >
            <HiTrash className="w-4 h-4 text-vault-danger" />
          </button>
        </div>
      ))}
    </div>
  );
};

export default SecretsList;
