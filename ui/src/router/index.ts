import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import Login from '../views/Login.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
    {
      path: '/logs',
      name: 'Logs',
      component: () => import('../views/Logs.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../views/Settings.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/storage',
      name: 'Storage',
      component: () => import('../views/Storage.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/sql-editor',
      name: 'SqlEditor',
      component: () => import('../views/SqlEditor.vue'),
      meta: { requiresAuth: true },
    },
    { path: '/login', name: 'Login', component: Login },
  ],
});

router.beforeEach((to, _, next) => {
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
