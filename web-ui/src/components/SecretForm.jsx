import React, { useState } from 'react';
import { HiPlus, HiX } from 'react-icons/hi';

const SecretForm = ({ onSave, onCancel }) => {
  const [path, setPath] = useState('');
  const [fields, setFields] = useState([{ key: '', value: '' }]);

  const addField = () => {
    setFields([...fields, { key: '', value: '' }]);
  };

  const removeField = (index) => {
    setFields(fields.filter((_, i) => i !== index));
  };

  const updateField = (index, field, value) => {
    const newFields = [...fields];
    newFields[index][field] = value;
    setFields(newFields);
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (!path.trim()) {
      alert('Please enter a secret path');
      return;
    }

    const data = {};
    let hasData = false;

    fields.forEach((field) => {
      if (field.key.trim() && field.value.trim()) {
        data[field.key] = field.value;
        hasData = true;
      }
    });

    if (!hasData) {
      alert('Please add at least one key-value pair');
      return;
    }

    onSave(path, data);
  };

  return (
    <div className="p-8">
      <div className="mb-8 pb-6 border-b border-vault-border-light">
        <h2 className="text-xl font-semibold text-vault-text-primary">Create Secret</h2>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="mb-6">
          <label htmlFor="path" className="block text-sm font-medium text-vault-text-primary mb-2">
            Path
          </label>
          <input
            id="path"
            type="text"
            className="w-full px-4 py-2.5 text-sm border border-vault-border-light rounded focus:outline-none focus:border-vault-primary focus:ring-2 focus:ring-vault-primary/20 transition-all"
            value={path}
            onChange={(e) => setPath(e.target.value)}
            placeholder="secret/myapp"
          />
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-vault-text-primary mb-2">Secret Data</label>
          <div className="space-y-3 mb-4">
            {fields.map((field, index) => (
              <div key={index} className="flex gap-3">
                <input
                  type="text"
                  className="flex-1 px-4 py-2.5 text-sm border border-vault-border-light rounded focus:outline-none focus:border-vault-primary focus:ring-2 focus:ring-vault-primary/20 transition-all"
                  value={field.key}
                  onChange={(e) => updateField(index, 'key', e.target.value)}
                  placeholder="Key"
                />
                <input
                  type="text"
                  className="flex-1 px-4 py-2.5 text-sm border border-vault-border-light rounded focus:outline-none focus:border-vault-primary focus:ring-2 focus:ring-vault-primary/20 transition-all"
                  value={field.value}
                  onChange={(e) => updateField(index, 'value', e.target.value)}
                  placeholder="Value"
                />
                {fields.length > 1 && (
                  <button
                    type="button"
                    className="p-2.5 text-white bg-vault-danger rounded hover:bg-red-600 transition-colors"
                    onClick={() => removeField(index)}
                  >
                    <HiX className="w-5 h-5" />
                  </button>
                )}
              </div>
            ))}
          </div>
          <button
            type="button"
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-vault-text-primary bg-white border border-vault-border-light rounded hover:bg-vault-light-bg transition-colors"
            onClick={addField}
          >
            <HiPlus className="w-4 h-4" />
            Add Field
          </button>
        </div>

        <div className="flex gap-3 justify-end pt-6 border-t border-vault-border-light">
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 text-sm font-medium text-vault-text-primary bg-white border border-vault-border-light rounded hover:bg-vault-light-bg transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="px-4 py-2 text-sm font-medium text-white bg-vault-primary rounded hover:bg-vault-primary-hover transition-colors"
          >
            Save Secret
          </button>
        </div>
      </form>
    </div>
  );
};

export default SecretForm;
