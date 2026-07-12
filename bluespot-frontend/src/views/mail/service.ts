import { post } from '@/utils/request'

export interface SendMailRequest {
  to: string
  subject: string
  body: string
}

export function sendMail(data: SendMailRequest) {
  return post<void>('/mail', data)
}
