self.importScripts('./spark-md5.js')

self.addEventListener('message', async (event) => {
  const {chunks} = event.data
  const start = new Date().getTime()
  const md5 = await createMd5(chunks)
  const end = new Date().getTime()
  const interval = (end - start) / 1000
  console.log('calc hash time: ', interval + 's')
  self.postMessage({
    type: 'hash',
    data: md5
  })
})

const createMd5 = (chunks) => {
  const spark = new SparkMD5()
  return new Promise(resolve => {
    function _read (i) {
      self.postMessage({
        type: 'progress',
        data: parseInt(i / chunks.length * 100)
      })
      if (i >= chunks.length) {
        const md5 = spark.end()
        resolve(md5)
        return
      }
      const blob = chunks[i].raw
      const reader = new FileReader()

      reader.onload = (event) => {
        const bytes = event?.target?.result
        spark.append(bytes)
        _read(i + 1)
      }
      reader.readAsArrayBuffer(blob)
    }
    _read(0)
  })
}

// 大文件会计算异常
// const getHash = (file) => {
//   return new Promise((resolve) => {
//     const fileReader = new FileReader()
//     fileReader.readAsArrayBuffer(file)
//     fileReader.onload = function (e) {
//       const fileMd5 = SparkMD5.ArrayBuffer.hash(e.target.result)
//       resolve(fileMd5)
//     }
//   })
// }