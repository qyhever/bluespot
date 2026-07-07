import { post } from '@/utils/request'

export interface LoginRequest {
  username: string
  password: string
}

export interface AuthTokenResponse {
  accessToken: string
  refreshToken: string
}

export function login(data: LoginRequest) {
  return post<AuthTokenResponse>('/auth/login', data)
}

export function refreshAccessToken(refreshToken: string) {
  return post<AuthTokenResponse>('/auth/refresh', {
    refreshToken,
  })
}
