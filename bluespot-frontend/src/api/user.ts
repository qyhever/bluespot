import { get } from '@/utils/request'

export interface UserInfo {
  userId: number
  username: string
  nickname: string
}

export function getUserInfo() {
  return get<UserInfo>('/user/info')
}
