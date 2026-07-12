<template>
  <div class="simple-upload-view">
    <div class="simple-upload-panel">
      <h1 class="page-title">附件上传</h1>
      <p class="page-tips">选择文件后自动上传，当前示例使用后端 /attach/upload 接口。</p>

      <t-upload
        v-model="files"
        theme="file"
        :max="1"
        :multiple="false"
        :request-method="uploadFile"
        tips="仅支持单文件上传"
        placeholder="请选择要上传的文件"
      />

      <div v-if="uploadResult" class="upload-result">
        <div class="result-title">上传结果</div>
        <dl class="result-list">
          <div class="result-item">
            <dt>存储文件名</dt>
            <dd>{{ uploadResult.fileName }}</dd>
          </div>
          <div class="result-item">
            <dt>原始文件名</dt>
            <dd>{{ uploadResult.originName }}</dd>
          </div>
          <div class="result-item">
            <dt>访问地址</dt>
            <dd>
              <a :href="uploadResult.url" target="_blank" rel="noreferrer">
                {{ uploadResult.url }}
              </a>
            </dd>
          </div>
        </dl>
      </div>
      <t-progress :percentage="progressValue"></t-progress>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { post } from '@/utils/request'
import type { RequestMethodResponse, UploadFile } from 'tdesign-vue-next'
import FakeProgress from 'fake-progress'

const fp = ref<FakeProgress | null>(null)

const progressValue = computed(() => {
  if (!fp.value) return 0
  return parseInt(String(fp.value.progress * 100), 10)
})

interface AttachUploadResponse {
  fileName: string
  originName: string
  url: string
}

const files = ref<UploadFile[]>([])
const uploadResult = ref<AttachUploadResponse | null>(null)

defineOptions({
  name: 'SimpleUploadView',
})

async function uploadFile(uploadFiles: UploadFile | UploadFile[]): Promise<RequestMethodResponse> {
  const currentFile = Array.isArray(uploadFiles) ? uploadFiles[0] : uploadFiles
  const rawFile = currentFile?.raw

  uploadResult.value = null

  if (!rawFile) {
    const error = '未获取到待上传文件'
    MessagePlugin.error(error)
    return {
      status: 'fail',
      error,
      response: {},
    }
  }

  const formData = new FormData()
  formData.append('file', rawFile)

  try {
    fp.value?.start()
    const response = await post<AttachUploadResponse>('/attach/upload', formData, {
      onUploadProgress(event) {
        console.log('event: ', event)
        if (!event.lengthComputable) return
        currentFile.percent = Math.round((event.loaded / event.total) * 100)
      },
    })
    fp.value?.end()
    uploadResult.value = response

    return {
      status: 'success',
      response: {
        url: response.url,
        fileName: response.fileName,
        originName: response.originName,
      },
    }
  } catch (err) {
    const error = err instanceof Error ? err.message : '上传失败'
    MessagePlugin.error(error)

    return {
      status: 'fail',
      error,
      response: {},
    }
  }
}

onMounted(() => {
  fp.value = new FakeProgress({
    timeConstant: 20000, // 总模拟时长 ms
    autoStart: true, // 实例化自动开始
  })
})
</script>

<style scoped>
.simple-upload-view {
  min-height: 100vh;
  padding: 32px;
  background: #f5f7fa;
}

.simple-upload-panel {
  max-width: 760px;
  padding: 24px;
  background: #ffffff;
  border: 1px solid #e7e7e7;
  border-radius: 8px;
}

.page-title {
  margin: 0;
  color: #1f2329;
  font-size: 22px;
  font-weight: 600;
  line-height: 30px;
}

.page-tips {
  margin: 8px 0 20px;
  color: #646a73;
  font-size: 14px;
  line-height: 22px;
}

.upload-result {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #e7e7e7;
}

.result-title {
  margin-bottom: 12px;
  color: #1f2329;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

.result-list {
  margin: 0;
}

.result-item {
  display: grid;
  grid-template-columns: 96px minmax(0, 1fr);
  gap: 12px;
  margin-bottom: 10px;
  color: #1f2329;
  font-size: 14px;
  line-height: 22px;
}

.result-item dt {
  color: #646a73;
}

.result-item dd {
  min-width: 0;
  margin: 0;
  overflow-wrap: anywhere;
}

.result-item a {
  color: #0052d9;
  text-decoration: none;
}

.result-item a:hover {
  text-decoration: underline;
}

@media (max-width: 640px) {
  .simple-upload-view {
    padding: 16px;
  }

  .simple-upload-panel {
    padding: 18px;
  }

  .result-item {
    grid-template-columns: 1fr;
    gap: 2px;
  }
}
</style>
