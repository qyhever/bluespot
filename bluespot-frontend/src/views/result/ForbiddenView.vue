<template>
  <div class="h-screen flex flex-col items-center justify-center">
    <img src="@/assets/result/403.svg" alt="pic" />
    <div class="text-2xl mt-4">403 Forbidden</div>
    <div class="mt-4 mb-4 text-gray-500">抱歉，您无权限访问此页面</div>
    <router-link to="/"><t-button>返回首页</t-button></router-link>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'ForbiddenView',
})
</script>

<style scoped></style>
