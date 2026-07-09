<template>
  <div class="h-screen flex flex-col items-center justify-center">
    <img src="@/assets/result/404.svg" alt="pic" />
    <div class="text-2xl mt-4">404 Not Found</div>
    <div class="mt-4 mb-4 text-gray-500">抱歉，您访问的页面不存在</div>
    <router-link to="/"><t-button>返回首页</t-button></router-link>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'NotFoundView',
})
</script>

<style scoped></style>
