import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import Login from '../views/Login.vue';

const router = createRouter({
  history: createWebHistory('/'),
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: () => import('../views/Dashboard.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/collections',
      name: 'Collections',
      component: () => import('../views/Collections.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/collections/new',
      name: 'CollectionNew',
      component: () => import('../views/CollectionNew.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/collections/:name',
      name: 'CollectionView',
      component: () => import('../views/CollectionView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/collections/:name/settings',
      name: 'CollectionSettings',
      component: () => import('../views/CollectionSettings.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/collections/:name/new',
      name: 'RecordNew',
      component: () => import('../views/RecordNew.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/collections/:name/edit/:id',
      name: 'RecordEdit',
      component: () => import('../views/RecordEdit.vue'),
      meta: { requiresAuth: true },
    },
    { path: '/login', name: 'Login', component: Login },
  ],
});

router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    next({ name: 'Login' });
  } else if (to.name === 'Login' && auth.isAuthenticated) {
    next({ name: 'Dashboard' });
  } else {
    next();
  }
});

export default router;
