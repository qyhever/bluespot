/* eslint-disable @typescript-eslint/no-explicit-any */
import { MessagePlugin } from 'tdesign-vue-next'
import { useUserStore } from '@/stores/user'
import router from '@/router'
import { request as http } from './fetch'
import type { RequestOptions } from './fetch'

const defaultOptions = {
  baseURL: '/sfr/api',
}
const SUCCESS_CODE = 1000

// Token 刷新状态管理
let refreshPromise: Promise<string> | null = null

/**
 * 获取或创建刷新 token 的 Promise
 * 确保多个并发请求只触发一次 token 刷新
 */
function getRefreshTokenPromise(): Promise<string> {
  if (!refreshPromise) {
    const userStore = useUserStore()
    refreshPromise = userStore.refreshAccessToken().finally(() => {
      refreshPromise = null
    })
  }
  return refreshPromise
}

const codeMessage: Record<string, string> = {
  400: '请求错误',
  401: '登录状态失效，请重新登录',
  403: '禁止访问',
  404: '请求地址不存在',
  500: '服务器繁忙，请稍后再试',
  502: '网关错误',
  503: '服务不可用，服务器暂时过载或维护',
  504: '网关超时',
}

function handleRequestHeader(configHeaders?: RequestOptions['headers']) {
  const headers = {
    ...(configHeaders || {}),
  }
  const userStore = useUserStore()
  const token = userStore.accessToken
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  return headers
}

function handleRequestParams(configParams?: RequestOptions['params']) {
  const params = {
    ...(configParams || {}),
    t: new Date().getTime(), // 避免缓存
  }
  return params
}

function requestStart(config: RequestOptions) {
  const headers = handleRequestHeader(config.headers)
  const params = handleRequestParams(config.params)
  return {
    ...config,
    headers,
    params,
  }
}

interface IRequestThenEndParam<T = any> {
  response: T
}

function requestThenEnd(param: IRequestThenEndParam) {
  const { response: responseData } = param
  // console.log('responseData: ', responseData)
  // success code
  if (responseData.code === SUCCESS_CODE) {
    return responseData.data
  }
  MessagePlugin.closeAll()
  MessagePlugin.warning(responseData.message || '操作失败')
  return Promise.reject(responseData)
}

interface IRequestCatchEndParam {
  error: {
    response: Response
  }
}

async function requestCatchEnd(
  param: IRequestCatchEndParam,
  originalOptions?: RequestOptions,
): Promise<any> {
  console.log('requestCatchEnd param', param)
  const { error } = param

  if (error.response?.status) {
    const status = error.response.status

    // 401 错误：尝试刷新 token 并重试
    if (status === 401 && originalOptions) {
      // 判断是否是刷新 token 接口本身失败
      if (originalOptions.url?.includes('/auth/refresh')) {
        MessagePlugin.closeAll()
        MessagePlugin.error('登录已过期，请重新登录')
        const userStore = useUserStore()
        userStore.logout()
        router.replace('/signin')
        return Promise.reject(new Error('Refresh token failed'))
      }

      try {
        // 等待 token 刷新完成（如果正在刷新，会复用同一个 Promise）
        await getRefreshTokenPromise()

        // Token 刷新成功，重试原请求
        console.log('Token refreshed, retrying request:', originalOptions.url)
        return request(originalOptions)
      } catch (refreshError) {
        // Token 刷新失败，跳转登录
        console.error('Token refresh failed:', refreshError)
        MessagePlugin.closeAll()
        MessagePlugin.error('登录已过期，请重新登录')
        const userStore = useUserStore()
        userStore.logout()
        router.replace('/signin')
        return Promise.reject(refreshError)
      }
    }

    // 其他错误码的处理
    const msg = codeMessage[status]
    MessagePlugin.closeAll()
    MessagePlugin.error(msg || '操作失败')
  }

  return Promise.reject(new Error('Request failed'))
}

export async function request<T = any>(opts: RequestOptions): Promise<T> {
  const options = {
    ...defaultOptions,
    ...opts,
  }
  try {
    const config = requestStart(options)
    const response = await http<T>(config)
    return requestThenEnd({
      response,
    })
  } catch (error: unknown) {
    return requestCatchEnd(
      {
        error: error as { response: Response },
      },
      options, // 传递原始请求配置，用于重试
    )
  }
}

export async function get<T = any>(
  url: string,
  params?: Record<string, any>,
  options?: Omit<RequestOptions, 'url' | 'params' | 'method'>,
) {
  return request<T>({
    method: 'GET',
    url,
    params,
    ...options,
  })
}

export async function post<T = any>(
  url: string,
  data?: any,
  options?: Omit<RequestOptions, 'url' | 'data' | 'method'>,
) {
  return request<T>({
    method: 'POST',
    url,
    data,
    ...options,
  })
}

export async function put<T = any>(
  url: string,
  data?: any,
  options?: Omit<RequestOptions, 'url' | 'data' | 'method'>,
) {
  return request<T>({
    method: 'PUT',
    url,
    data,
    ...options,
  })
}

export async function del<T = any>(
  url: string,
  params?: Record<string, any>,
  options?: Omit<RequestOptions, 'url' | 'params' | 'method'>,
) {
  return request<T>({
    method: 'DELETE',
    url,
    params,
    ...options,
  })
}

export async function patch<T = any>(
  url: string,
  data?: any,
  options?: Omit<RequestOptions, 'url' | 'data' | 'method'>,
) {
  return request<T>({
    method: 'PATCH',
    url,
    data,
    ...options,
  })
}
