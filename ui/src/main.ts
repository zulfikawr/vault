import { createApp } from 'vue';
import { createPinia } from 'pinia';
import './style.css';
import App from './App.vue';
import router from './router';
import axios from 'axios';

const pinia = createPinia();
const app = createApp(App);

// Initialize axios
const token = localStorage.getItem('token');
if (token) {
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
}

app.use(pinia);
app.use(router);
app.mount('#app');