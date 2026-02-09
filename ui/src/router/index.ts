import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import Login from '../views/Login.vue';

const router = createRouter({
  history: createWebHistory('/_/'),
  routes: [
    { path: '/login', name: 'Login', component: Login },
    { 
      path: '/', 
      name: 'Dashboard', 
      component: () => import('../views/Dashboard.vue'),
      meta: { requiresAuth: true }
    },
  ],
});

router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next({ name: 'Login' });
  } else {
    next();
  }
});

export default router;
