/* eslint-disable @typescript-eslint/no-explicit-any */

export interface RequestOptions {
  url: string
  baseURL?: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS'
  params?: Record<string, any>
  data?: any
  headers?: Record<string, string>
  onUploadProgress?: (event: ProgressEvent) => void
  onDownloadProgress?: (event: ProgressEvent) => void
  responseType?: 'json' | 'blob' | 'text'
  signal?: AbortSignal
  timeout?: number // in milliseconds
}

export interface XhrRequestOptions extends RequestOptions {
  onUploadProgress?: (event: ProgressEvent) => void
  onDownloadProgress?: (event: ProgressEvent) => void
}

/**
 * 封装网络请求函数
 * @param options - 请求配置项
 * @returns 返回一个解析为响应数据的 Promise。
 */
export function request<T = any>(options: RequestOptions): Promise<T> {
  const {
    url,
    baseURL = '',
    method: inputMethod = 'GET',
    params,
    data,
    headers = {},
    onUploadProgress,
    onDownloadProgress,
    responseType,
    signal,
    timeout,
  } = options

  const method = inputMethod.toUpperCase()

  const base = baseURL.startsWith('/') ? window.location.origin + baseURL : baseURL
  const fullUrl = new URL(base + url)
  if (params) {
    for (const key in params) {
      if (Object.prototype.hasOwnProperty.call(params, key)) {
        fullUrl.searchParams.append(key, params[key])
      }
    }
  }

  if ((onUploadProgress && data) || onDownloadProgress) {
    return sendReqByXhr<T>({
      ...options,
      method: method as any,
      url: fullUrl.toString(),
    })
  }

  return new Promise<T>((resolve, reject) => {
    const controller = new AbortController()
    const fetchSignal = controller.signal
    let timedOut = false
    let timeoutId: number | null = null

    const abortHandler = () => {
      controller.abort()
    }

    if (signal) {
      if (signal.aborted) {
        return reject(new Error('请求已取消'))
      }
      signal.addEventListener('abort', abortHandler)
    }

    if (timeout) {
      timeoutId = setTimeout(() => {
        timedOut = true
        controller.abort()
      }, timeout)
    }

    const cleanup = () => {
      if (timeoutId) clearTimeout(timeoutId)
      if (signal) signal.removeEventListener('abort', abortHandler)
    }

    const fetchOptions: RequestInit = {
      method,
      headers,
      signal: fetchSignal,
    }

    if (method !== 'GET' && method !== 'HEAD') {
      if (data instanceof FormData || typeof data === 'string' || data instanceof URLSearchParams) {
        fetchOptions.body = data
      } else if (data) {
        const contentTypeKey = Object.keys(headers).find(
          (key) => key.toLowerCase() === 'content-type',
        )
        const userContentType = contentTypeKey ? headers[contentTypeKey] : null

        if (userContentType && userContentType.includes('application/x-www-form-urlencoded')) {
          fetchOptions.body = new URLSearchParams(data)
        } else {
          fetchOptions.body = JSON.stringify(data)
          if (!userContentType) {
            headers['Content-Type'] = 'application/json'
          }
        }
      }
    }

    fetch(fullUrl.toString(), fetchOptions)
      .then((response) => {
        let bodyParser: 'json' | 'blob' | 'text'
        const contentType = response.headers.get('content-type')
        const isJson = contentType && contentType.includes('application/json')

        if (responseType === 'blob') {
          bodyParser = 'blob'
        } else if (responseType === 'json') {
          bodyParser = 'json'
        } else {
          bodyParser = isJson ? 'json' : 'text'
        }

        if (response.status === 204) {
          return Promise.resolve(null)
        }
        const bodyPromise = response[bodyParser]()

        return bodyPromise.then((body) => {
          if (response.ok) {
            return body
          } else {
            let error: any
            if (isJson && body && (body as any).message) {
              error = new Error((body as any).message)
              error.data = body
            } else if (typeof body === 'string' && body) {
              error = new Error(body)
            } else {
              error = new Error(`HTTP 错误! 状态码: ${response.status}`)
            }
            error.response = response
            throw error
          }
        })
      })
      .then((data) => {
        cleanup()
        resolve(data as T)
      })
      .catch((error) => {
        cleanup()
        if (error.name === 'AbortError') {
          if (timedOut) {
            reject(new Error('请求超时'))
          } else {
            reject(new Error('请求已取消'))
          }
        } else {
          reject(error)
        }
      })
  })
}

/**
 * 使用 XMLHttpRequest 发送请求
 * @param options - 配置项
 * @returns 返回一个解析为响应数据的 Promise。
 */
function sendReqByXhr<T = any>(options: XhrRequestOptions): Promise<T> {
  return new Promise<T>((resolve, reject) => {
    const {
      method = 'GET',
      url,
      headers = {},
      data,
      onUploadProgress,
      onDownloadProgress,
      responseType,
      signal,
      timeout,
    } = options

    const xhr = new XMLHttpRequest()
    xhr.open(method, url, true)

    if (timeout) {
      xhr.timeout = timeout
    }

    if (signal) {
      signal.addEventListener('abort', () => {
        xhr.abort()
      })
    }

    if (responseType === 'blob' || responseType === 'json') {
      xhr.responseType = responseType
    }

    let requestData: any = data
    if (
      typeof data === 'object' &&
      data !== null &&
      !(data instanceof FormData) &&
      !(data instanceof Blob) &&
      !(data instanceof ArrayBuffer) &&
      !(data instanceof URLSearchParams)
    ) {
      const contentTypeKey = Object.keys(headers).find(
        (key) => key.toLowerCase() === 'content-type',
      )
      const userContentType = contentTypeKey ? headers[contentTypeKey] : null

      if (userContentType && userContentType.includes('application/x-www-form-urlencoded')) {
        requestData = new URLSearchParams(data).toString()
      } else {
        requestData = JSON.stringify(data)
        if (!userContentType) {
          headers['Content-Type'] = 'application/json'
        }
      }
    }

    for (const key in headers) {
      if (Object.prototype.hasOwnProperty.call(headers, key)) {
        const value = headers[key]
        if (value !== undefined) {
          xhr.setRequestHeader(key, value)
        }
      }
    }

    if (onUploadProgress) {
      xhr.upload.onprogress = onUploadProgress
    }

    if (onDownloadProgress) {
      xhr.onprogress = onDownloadProgress
    }

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        if (xhr.responseType === 'blob' || xhr.responseType === 'json') {
          resolve(xhr.response)
        } else {
          try {
            resolve(JSON.parse(xhr.responseText))
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
          } catch (err) {
            resolve(xhr.responseText as any)
          }
        }
      } else {
        let error: any
        try {
          const errorResponse = JSON.parse(xhr.responseText)
          if (errorResponse && errorResponse.message) {
            error = new Error(errorResponse.message)
            error.data = errorResponse
          } else {
            error = new Error(`请求失败，状态码: ${xhr.status}`)
          }
          // eslint-disable-next-line @typescript-eslint/no-unused-vars
        } catch (err) {
          if (xhr.responseText) {
            error = new Error(xhr.responseText)
          } else {
            error = new Error(`请求失败，状态码: ${xhr.status}`)
          }
        }
        error.response = xhr
        reject(error)
      }
    }

    xhr.onerror = () => {
      reject(new Error('网络错误'))
    }

    xhr.onabort = () => {
      reject(new Error('请求已取消'))
    }

    xhr.ontimeout = () => {
      reject(new Error('请求超时'))
    }

    xhr.send(requestData)
  })
}

/**
 * 发送 GET 请求
 * @param url - 请求地址
 * @param params - URL 查询参数
 * @param options - 其他请求配置项
 * @returns
 */
export function get<T = any>(
  url: string,
  params?: Record<string, any>,
  options?: Omit<RequestOptions, 'url' | 'params' | 'method'>,
): Promise<T> {
  return request<T>({
    method: 'GET',
    url,
    params,
    ...options,
  })
}

/**
 * 发送 POST 请求
 * @param url - 请求地址
 * @param data - 请求体
 * @param options - 其他请求配置项
 * @returns
 */
export function post<T = any>(
  url: string,
  data?: any,
  options?: Omit<RequestOptions, 'url' | 'data' | 'method'>,
): Promise<T> {
  return request<T>({
    method: 'POST',
    url,
    data,
    ...options,
  })
}

/**
 * 发送 PUT 请求
 * @param url - 请求地址
 * @param data - 请求体
 * @param options - 其他请求配置项
 * @returns
 */
export function put<T = any>(
  url: string,
  data?: any,
  options?: Omit<RequestOptions, 'url' | 'data' | 'method'>,
): Promise<T> {
  return request<T>({
    method: 'PUT',
    url,
    data,
    ...options,
  })
}

/**
 * 发送 DELETE 请求
 * @param url - 请求地址
 * @param options - 其他请求配置项
 * @returns
 */
export function del<T = any>(
  url: string,
  options?: Omit<RequestOptions, 'url' | 'method'>,
): Promise<T> {
  return request<T>({
    method: 'DELETE',
    url,
    ...options,
  })
}

/**
 * 发送 PATCH 请求
 * @param url - 请求地址
 * @param data - 请求体
 * @param options - 其他请求配置项
 * @returns
 */
export function patch<T = any>(
  url: string,
  data?: any,
  options?: Omit<RequestOptions, 'url' | 'data' | 'method'>,
): Promise<T> {
  return request<T>({
    method: 'PATCH',
    url,
    data,
    ...options,
  })
}
