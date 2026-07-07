/**
 * 全局类型声明文件
 *
 */

/**
 * 分页查询通用参数
 */
interface IPaginationQuery {
  currentPage: number
  pageSize: number
  sortField?: string
  sortValue?: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  [key: string]: any // 允许扩展其他查询参数
}

/**
 * 分页响应通用类型
 */
interface IPaginationResponse<T> {
  list: T[]
  total: number
}

/**
 * API 响应通用结构
 */
interface IApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

/**
 * 布尔枚举类型（0 否，1 是）
 */
type BooleanEnum = 0 | 1

/**
 * 排序方向
 */
type SortOrder = 'asc' | 'desc'
