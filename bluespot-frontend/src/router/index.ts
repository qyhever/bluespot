import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import LoginView from '../views/LoginView.vue'

const HOME_PATH = '/home'
const LOGIN_PATH = '/signin'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: HOME_PATH,
    },
    {
      path: '/home',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: LOGIN_PATH,
      name: 'signin',
      component: LoginView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

router.beforeEach((to) => {
  const userStore = useUserStore()
  const hasToken = Boolean(userStore.accessToken)
  const isLoginPage = to.path === LOGIN_PATH

  if (isLoginPage) {
    return hasToken ? HOME_PATH : true
  }

  return hasToken
    ? true
    : {
        path: LOGIN_PATH,
        query: to.fullPath === HOME_PATH ? undefined : { redirect: to.fullPath },
      }
})

export default router
