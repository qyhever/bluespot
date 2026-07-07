import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import {
  refreshAccessToken as refreshAccessTokenApi,
  type AuthTokenResponse,
} from '@/api/global'
import { getUserInfo, type UserInfo } from '@/api/user'

const ACCESS_TOKEN_KEY = 'accessToken'
const REFRESH_TOKEN_KEY = 'refreshToken'

export const useUserStore = defineStore('user', () => {
  const accessToken = ref(localStorage.getItem(ACCESS_TOKEN_KEY) || '')
  const refreshToken = ref(localStorage.getItem(REFRESH_TOKEN_KEY) || '')
  const userInfo = ref<UserInfo | null>(null)
  const userInfoLoading = ref(false)
  const displayName = computed(() => userInfo.value?.nickname || userInfo.value?.username || '用户')
  let userInfoPromise: Promise<UserInfo> | null = null

  function setTokens(tokens: AuthTokenResponse) {
    accessToken.value = tokens.accessToken
    refreshToken.value = tokens.refreshToken
    localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken)
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) {
      return Promise.reject(new Error('Refresh token missing'))
    }

    const response = await refreshAccessTokenApi(refreshToken.value)

    setTokens(response)
    return response.accessToken
  }

  function fetchUserInfo(force = false) {
    if (userInfo.value && !force) {
      return Promise.resolve(userInfo.value)
    }
    if (userInfoPromise) {
      return userInfoPromise
    }

    userInfoLoading.value = true
    userInfoPromise = getUserInfo()
      .then((info) => {
        userInfo.value = info
        return info
      })
      .finally(() => {
        userInfoLoading.value = false
        userInfoPromise = null
      })

    return userInfoPromise
  }

  function logout() {
    accessToken.value = ''
    refreshToken.value = ''
    userInfo.value = null
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  }

  return {
    accessToken,
    refreshToken,
    userInfo,
    userInfoLoading,
    displayName,
    setTokens,
    refreshAccessToken,
    fetchUserInfo,
    logout,
  }
})
