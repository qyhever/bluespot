<template>
  <div
    v-if="visible"
    class="fixed top-4 right-4 z-[9999] bg-white border border-gray-200 rounded-lg shadow-lg p-4 max-w-sm"
  >
    <div class="flex items-start gap-3">
      <div class="flex-shrink-0">
        <svg class="w-6 h-6 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
      </div>
      <div class="flex-1">
        <h4 class="text-md font-bold text-gray-900">发现新版本</h4>
        <template v-if="versionExtraInfoVisible">
          <div class="mt-1 text-sm text-gray-900">
            <div class="font-bold">本地版本</div>
            <div class="text-gray-600">
              <div>hash: {{ visibleInfo.oldBuildHash }}</div>
              <div>构建时间: {{ visibleInfo.oldBuildTime }}</div>
            </div>
          </div>
          <div class="mt-1 text-sm text-gray-900">
            <div class="font-bold">新版本</div>
            <div class="text-gray-600">
              <div>hash: {{ visibleInfo.newBuildHash }}</div>
              <div>构建时间: {{ visibleInfo.newBuildTime }}</div>
            </div>
          </div>
        </template>
        <p class="mt-1 text-sm text-gray-600">检测到新版本已发布，建议立即更新以获得最佳体验。</p>
        <div class="mt-3 flex gap-2">
          <button
            class="cursor-pointer inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            @click="handleUpdate"
          >
            立即更新
          </button>
          <button
            class="cursor-pointer inline-flex items-center px-3 py-1.5 border border-gray-300 text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            @click="handleLater"
          >
            稍后提醒
          </button>
          <button
            class="cursor-pointer inline-flex items-center px-3 py-1.5 text-xs font-medium text-gray-400 hover:text-gray-600"
            @click="handleClose"
          >
            忽略
          </button>
        </div>
      </div>
      <button class="flex-shrink-0 text-gray-400 hover:text-gray-600" @click="handleClose">
        <svg class="w-[16px] h-[14px]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { versionChecker, type UpdateInfo } from '@/utils/version-checker'

defineOptions({
  name: 'VersionUpdatePopup',
})

const visible = ref(false)
const versionExtraInfoVisible = ref(Boolean(localStorage.getItem('versionExtraInfoVisible')))
const visibleInfo = ref({} as UpdateInfo)

const open = (info: UpdateInfo) => {
  visible.value = true
  visibleInfo.value = info
}

const close = () => {
  visible.value = false
}

const handleUpdate = () => {
  versionChecker.refresh()
}

const handleLater = () => {
  close()
  // 5分钟后再次提醒
  setTimeout(
    () => {
      open(visibleInfo.value)
    },
    5 * 60 * 1000,
  )
}

const handleClose = () => {
  close()
}

// 暴露方法给父组件
defineExpose({
  open,
  close,
})
</script>
