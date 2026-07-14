import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import LoginView from '../views/LoginView.vue'
import BasicLayout from '../layouts/BasicLayout.vue'

const HOME_PATH = '/home'
const LOGIN_PATH = '/signin'

declare module 'vue-router' {
  interface RouteMeta {
    authKey?: string
    activeMenu?: string
  }
}

const basicRoutes = [
  {
    path: HOME_PATH,
    name: 'home',
    component: () => import('../views/HomeView.vue'),
  },
  {
    path: '/attach/simple-upload',
    name: 'simple-upload',
    component: () => import('../views/attach/SimpleUploadView.vue'),
  },
  {
    path: '/attach/chunk-upload',
    name: 'chunk-upload',
    component: () => import('../views/attach/ChunkUploadView.vue'),
  },
  {
    path: '/mail/send',
    name: 'mail-send',
    component: () => import('../views/mail/SendMail.vue'),
  },
  {
    path: '/notification/tg',
    name: 'telegram-send',
    component: () => import('../views/telegram/SendTelegram.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: HOME_PATH,
    },
    {
      path: LOGIN_PATH,
      name: 'signin',
      component: LoginView,
    },
    {
      path: '/',
      component: BasicLayout,
      children: basicRoutes,
    },
    {
      path: '/forbidden',
      component: () => import('@/views/result/ForbiddenView.vue'),
    },
    {
      path: '/notfound',
      component: () => import('@/views/result/NotfoundView.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      component: () => import('@/views/result/NotfoundView.vue'),
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

router.beforeEach(async (to) => {
  const userStore = useUserStore()
  const hasToken = Boolean(userStore.accessToken)
  const isLoginPage = to.path === LOGIN_PATH

  if (isLoginPage) {
    return hasToken ? HOME_PATH : true
  }

  if (hasToken && !userStore.userInfo) {
    await userStore.fetchUserInfo()
  }

  return hasToken
    ? true
    : {
        path: LOGIN_PATH,
        query: to.fullPath === HOME_PATH ? undefined : { redirect: to.fullPath },
      }
})

export const getActive = (maxLevel = 3): string => {
  // 非组件内调用必须通过Router实例获取当前路由
  const route = router.currentRoute.value

  if (!route.path) {
    return ''
  }
  if (route.meta?.activeMenu) {
    return route.meta.activeMenu
  }
  const ret = route.path
    .split('/')
    .filter((_item: string, index: number) => index <= maxLevel && index > 0)
    .map((item: string) => `/${item}`)
    .join('')

  return ret
}

export default router
