<template>
  <div class="p-20">
    <div class="upload-file">
      <el-button>点击上传</el-button>
      <input ref="fileInput" type="file" class="input-file" @change="onChange" />
    </div>
    <div class="text-color-regular flex items-center pt-10">
      <i class="el-icon-info fs-16" />
      <span class="fs-12 ml-6">只能上传: mp4</span>
    </div>
    <div class="pt-10">
      <div class="text-color-regular">
        {{ fileInfo.name }}
      </div>
      <div v-if="status" class="text-color-secondary pt-6">
        <div v-if="status === 'done'">
          <div>
            <span>上传完成: </span>
            <span>文件路径 {{ fileURL }}</span>
          </div>
          <span>计算hash耗时: {{ timeInfo.hash }}s</span>
          <span class="ml-6">上传切片耗时: {{ timeInfo.chunk }}s</span>
          <span class="ml-6">合并切片耗时: {{ timeInfo.merge }}s</span>
          <span class="ml-6">总耗时: {{ totalTime }}s</span>
        </div>
        <template v-else>
          <i class="el-icon-loading" />
          <span v-if="status === 'hashing'">生成文件hash中 {{ progressInfo.hash }}%</span>
          <span v-else-if="status === 'chunking'">文件切片上传中 {{ progressInfo.chunk }}%</span>
          <span v-else-if="status === 'merging'">文件切片合并中...</span>
          <span v-else>上传中...</span>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { verifyFile, createChunks, uploadChunks, mergeChunks } from '@/utils/upload'
import { toFixed } from '@/utils/index'

defineOptions({
  name: 'SliceUploadView',
})

type UploadStatus = '' | 'hashing' | 'chunking' | 'merging' | 'done'

const fileInput = ref<HTMLInputElement | null>(null)
const fileInfo = ref({
  name: '',
})
const timeInfo = ref({
  hash: 0,
  chunk: 0,
  merge: 0,
})
const progressInfo = ref({
  hash: 0,
  chunk: 0,
})
const status = ref<UploadStatus>('')
const fileURL = ref('')

const totalTime = computed(() => {
  return toFixed(timeInfo.value.hash + timeInfo.value.chunk + timeInfo.value.merge)
})

function getInterval(start: number) {
  const end = new Date().getTime()
  const interval = (end - start) / 1000
  return toFixed(interval)
}

function checkFile(file: File) {
  if (file.size > 10 * 1024 * 1024 * 1024) {
    alert('上传文件不超过10G')
    return false
  }
  if (file.name.replace(/.+\./, '') !== 'mp4') {
    alert('上传文件仅支持MP4格式')
    return false
  }
  return true
}

function reset() {
  timeInfo.value = {
    hash: 0,
    chunk: 0,
    merge: 0,
  }
  progressInfo.value = {
    hash: 0,
    chunk: 0,
  }
}

function onChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  if (!file) {
    return
  }
  if (!checkFile(file)) {
    return
  }
  reset()
  fileInfo.value = file
  const chunks = createChunks(file)
  console.log('chunks: ', chunks)
  const hashStart = new Date().getTime()
  status.value = 'hashing'
  const worker = new Worker('./hash.js')
  worker.postMessage({
    chunks,
    file,
  })
  worker.onerror = (err) => {
    console.log('worker.onerror', err)
    worker.terminate()
  }
  worker.onmessage = async (event) => {
    timeInfo.value.hash = +getInterval(hashStart)
    const { type } = event.data
    if (type === 'progress') {
      const percent = event.data.data
      progressInfo.value.hash = percent
      return
    }
    const fileHash = event.data.data
    console.log('fileHash: ', fileHash)
    worker.terminate()
    const res = await verifyFile(fileHash, chunks, file.name)
    if (res.fileURL) {
      // 已上传过，直接返回
      fileURL.value = res.fileURL
      status.value = 'done'
      return
    }
    // console.log('uploadChunks: ', uploadChunks)
    // console.log('mergeChunks: ', mergeChunks)
    const chunkStart = new Date().getTime()
    status.value = 'chunking'
    await uploadChunks(chunks, fileHash, (val) => {
      progressInfo.value.chunk = val
    })
    timeInfo.value.chunk = +getInterval(chunkStart)
    const mergeStart = new Date().getTime()
    status.value = 'merging'
    const url = await mergeChunks(fileHash, file.name)
    status.value = 'done'
    timeInfo.value.merge = +getInterval(mergeStart)
    fileURL.value = url
    // verifyFile(fileHash, chunks, file)
  }
}
</script>

<style scoped>
.upload-file {
  position: relative;
  display: inline-block;
  cursor: pointer;
  &:hover {
    opacity: 0.8;
  }
}
.input-file {
  opacity: 0;
  position: absolute;
  top: 0;
  right: 0;
  height: 100%;
  width: 100%;
  cursor: pointer;
}
</style>
