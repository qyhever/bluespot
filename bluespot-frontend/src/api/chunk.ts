import { post } from '@/utils/request'

export interface UploadVerifyRequest {
  fileMd5: string
  fileName: string
  fileSize: number
}

export interface UploadVerifyResponse {
  isExists: boolean
  url: string
  uploadId: string
  uploadedChunks: number[]
}

export interface UploadChunkRequest {
  chunk: Blob
  uploadId: string
  fileMd5: string
  chunkIndex: number
  fileName: string
  chunkTotal: number
}

export interface MergeChunksRequest {
  uploadId: string
  fileMd5: string
  chunkLength: number
}

export interface MergeChunksResponse {
  url: string
  msg: string
}

/**
 * 上传文件预检(秒传检测)
 * @param params 文件信息
 * @returns 预检结果
 */
export function uploadVerify(params: UploadVerifyRequest) {
  return post<UploadVerifyResponse>('/upload/verify', params)
}

/**
 * 上传文件分片
 * @param params 分片信息
 * @returns 上传结果
 */
export function uploadChunk(params: UploadChunkRequest) {
  const {
    chunkIndex,
    chunkTotal,
  } = params
  const fd = new FormData()
  fd.append('chunk', params.chunk)
  fd.append('uploadId', params.uploadId)
  fd.append('fileMd5', params.fileMd5)
  fd.append('chunkIndex', chunkIndex.toString())
  fd.append('fileName', params.fileName)
  return post<void>(`/upload/chunk?chunkIndex=${chunkIndex}&chunkTotal=${chunkTotal}`, fd)
}

/**
 * 合并分片
 * @param params 合并请求
 * @returns 合并结果
 */
export function mergeChunks(params: MergeChunksRequest) {
  return post<MergeChunksResponse>('/upload/merge', params)
}
