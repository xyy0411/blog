import {createRouter, createWebHistory} from 'vue-router';
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes :  [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/home.vue')
    }
    ,
    {
      path: '/about',
      name: 'About',
      component: () => import('@/views/about.vue')
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import("@/views/login.vue")
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/register.vue')
    },
  ]
})

export default router;

