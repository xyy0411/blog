import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login.vue'),
      meta: { hideLayout: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/register.vue'),
      meta: { hideLayout: true },
    },
    {
      path: '/',
      name: 'HomeRedirect',
      redirect: '/matching-application',
    },
    {
      path: '/about',
      name: 'About',
      component: () => import('@/views/about.vue'),
    },
    {
      path: '/matching-stats',
      name: 'MatchingStats',
      component: () => import('@/views/matching-stats.vue'),
    },
    {
      path: '/matching-application',
      name: 'MatchingApplication',
      component: () => import('@/views/matching-application.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/not-found.vue'),
    },
  ],
});

export default router;
