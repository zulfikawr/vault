import { defineStore } from 'pinia';
import axios from 'axios';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: JSON.parse(localStorage.getItem('user') || 'null'),
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
  },
  actions: {
    async login(identity: string, password: string) {
      try {
        const response = await axios.post('/api/collections/users/auth-with-password', {
          identity,
          password,
        });
        const { token, record } = response.data.data;
        this.token = token;
        this.user = record;
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(record));
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        return true;
      } catch (error) {
        console.error('Login failed', error);
        return false;
      }
    },
    logout() {
      this.token = null;
      this.user = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      delete axios.defaults.headers.common['Authorization'];
    },
  },
});
