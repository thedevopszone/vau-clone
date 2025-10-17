import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_VAULT_ADDR || 'http://127.0.0.1:8200';

class VaultAPI {
  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add token to requests if available
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('vault_token');
      if (token) {
        config.headers['X-Vault-Token'] = token;
      }
      return config;
    });
  }

  // System Operations
  async getHealth() {
    const response = await this.client.get('/v1/sys/health');
    return response.data;
  }

  async getStatus() {
    const response = await this.client.get('/v1/sys/status');
    return response.data;
  }

  async initialize() {
    const response = await this.client.post('/v1/sys/init');
    return response.data;
  }

  async unseal(key) {
    const response = await this.client.post('/v1/sys/unseal', { key });
    return response.data;
  }

  async seal() {
    const response = await this.client.post('/v1/sys/seal');
    return response.data;
  }

  // Secret Operations
  async writeSecret(path, data) {
    const response = await this.client.post(`/v1/secret/${path}`, { data });
    return response.data;
  }

  async readSecret(path) {
    const response = await this.client.get(`/v1/secret/${path}`);
    return response.data;
  }

  async deleteSecret(path) {
    const response = await this.client.delete(`/v1/secret/${path}`);
    return response.data;
  }

  async listSecrets(prefix = '') {
    const response = await this.client.get('/v1/secrets/list', {
      params: { prefix },
    });
    return response.data;
  }

  // Authentication
  async createToken(ttl = '24h') {
    const response = await this.client.post('/v1/auth/token/create', { ttl });
    return response.data;
  }

  // Token management
  setToken(token) {
    localStorage.setItem('vault_token', token);
  }

  getToken() {
    return localStorage.getItem('vault_token');
  }

  clearToken() {
    localStorage.removeItem('vault_token');
  }
}

export default new VaultAPI();
