import SparkMD5 from 'spark-md5'

interface WorkerInput {
  file: File
  chunkSize: number
}

interface WorkerProgress {
  type: 'progress'
  progress: number
}

interface WorkerSuccess {
  type: 'success'
  md5: string
}

interface WorkerError {
  type: 'error'
  message: string
}

self.onmessage = async (event: MessageEvent<WorkerInput>) => {
  const { file, chunkSize } = event.data

  try {
    const spark = new SparkMD5.ArrayBuffer()
    const chunks = Math.ceil(file.size / chunkSize)

    for (let index = 0; index < chunks; index++) {
      const start = index * chunkSize
      const end = Math.min(start + chunkSize, file.size)
      const blob = file.slice(start, end)

      const buffer = await blob.arrayBuffer()

      spark.append(buffer)

      const progress = Number((((index + 1) / chunks) * 100).toFixed(2))

      self.postMessage({
        type: 'progress',
        progress
      } satisfies WorkerProgress)
    }

    const md5 = spark.end()

    self.postMessage({
      type: 'success',
      md5
    } satisfies WorkerSuccess)
  } catch (error) {
    self.postMessage({
      type: 'error',
      message: error instanceof Error ? error.message : 'MD5 计算失败'
    } satisfies WorkerError)
  }
}