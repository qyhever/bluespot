import { post } from '@/utils/request'

export interface SendTelegramMessageRequest {
  text: string
}

export function sendTelegramMessage(data: SendTelegramMessageRequest) {
  return post<void>('/telegram', data)
}
