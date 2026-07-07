/**
 * 版本更新检测工具
 */

interface MetaInfo {
  buildHash: string
  buildTime: string
}

export interface UpdateInfo {
  oldBuildHash: string
  oldBuildTime: string
  newBuildHash: string
  newBuildTime: string
}

export class VersionChecker {
  // private currentVersion: string
  private metaInfo: MetaInfo = { buildHash: '', buildTime: '' }
  private checkInterval: number
  private timerId: number | null = null
  private onUpdateCallback?: (info: UpdateInfo) => void

  constructor(checkInterval = 30000) {
    this.checkInterval = checkInterval
  }

  /**
   * 开始版本检查
   */
  async start(onUpdate?: (info: UpdateInfo) => void) {
    this.onUpdateCallback = onUpdate
    // 首次获取版本信息
    this.metaInfo = {
      buildHash: LOCAL_BUILD_HASH,
      buildTime: LOCAL_BUILD_TIME,
    }
    // 定期检查版本
    this.runLoopTask()
    window.addEventListener('visibilitychange', () => {
      if (document.visibilityState === 'visible') {
        this.checkVersion()
        this.runLoopTask()
      }
    })
    document.addEventListener('focus', () => {
      this.checkVersion()
      this.runLoopTask()
    })
  }

  /**
   * 停止版本检查
   */
  stop() {
    if (this.timerId) {
      clearInterval(this.timerId)
      this.timerId = null
    }
  }

  // 定时检查
  private async runLoopTask() {
    if (this.timerId) {
      clearInterval(this.timerId)
      this.timerId = null
    }
    this.timerId = window.setInterval(() => {
      this.checkVersion()
    }, this.checkInterval)
  }

  private async fetchMetaInfo() {
    try {
      const baseURL = import.meta.env.BASE_URL
      const url = `${baseURL}meta.json?t=${Date.now()}`
      const response = await fetch(url, {
        cache: 'no-cache',
      })

      if (!response.ok) {
        console.log('fetchMetaInfo failed:', response.status)
        return
      }

      const versionInfo: MetaInfo = await response.json()
      return versionInfo
    } catch (error) {
      console.log('fetchMetaInfo failed:', error)
    }
  }

  /**
   * 检查版本
   */
  private async checkVersion() {
    try {
      const versionInfo = await this.fetchMetaInfo()

      if (
        versionInfo &&
        versionInfo.buildHash &&
        this.metaInfo.buildHash &&
        versionInfo.buildHash !== this.metaInfo.buildHash
      ) {
        // console.log('检测到新版本:', versionInfo.buildHash, '当前版本:', this.metaInfo.buildHash)
        this.onUpdateCallback?.({
          oldBuildHash: this.metaInfo.buildHash,
          oldBuildTime: this.metaInfo.buildTime,
          newBuildHash: versionInfo.buildHash,
          newBuildTime: versionInfo.buildTime,
        })
      }
      return versionInfo
    } catch (error) {
      console.log('版本检查出错:', error)
    }
  }

  /**
   * 刷新页面
   */
  refresh() {
    window.location.reload()
  }
}

export const versionChecker = new VersionChecker()
