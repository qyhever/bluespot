<template>
  <div id="app">
    <router-view />
    <VersionUpdatePopup  ref="versionUpdatePopupRef" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import VersionUpdatePopup from '@/components/VersionUpdatePopup.vue'
import { versionChecker } from '@/utils/version-checker'

defineOptions({
  name: 'App',
})

const versionUpdatePopupRef = ref<InstanceType<typeof VersionUpdatePopup> | null>(null)

onMounted(() => {
  if (import.meta.env.VITE_APP_MODE_ENV !== 'dev') {
    // 启动版本检测
    versionChecker.start((info) => {
      console.log('检测到新版本:', info)
      versionUpdatePopupRef.value?.open(info)
    })
  }
})

onBeforeUnmount(() => {
  // 停止版本检测
  versionChecker.stop()
})
</script>
