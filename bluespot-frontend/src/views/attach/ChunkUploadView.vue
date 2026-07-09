<template>
  <div class="p-20">
    <div class="upload-file">
      <t-button>点击上传</t-button>
      <input ref="fileInput" type="file" class="input-file" @change="onChange" />
    </div>
    <div class="flex items-center pt-10 text-color-regular">
      <t-icon name="info-circle" />
      <span class="ml-6 fs-12">只能上传: mp4</span>
    </div>
    <div class="pt-10">
      <div class="text-color-regular">
        {{ fileInfo.name }}
      </div>
      <div v-if="fileInfo.size" class="pt-6 text-color-secondary">文件大小: {{ fileSizeText }}</div>
      <div v-if="status" class="pt-6 text-color-secondary">
        <div v-if="status === 'done'">
          <div>
            <span>上传完成: </span>
            <span>文件路径 {{ fileURL }}</span>
          </div>
        </div>
        <div v-else-if="status === 'error'" class="text-error">
          <span>上传失败: {{ errorMessage }}</span>
        </div>
        <template v-else>
          <t-icon name="loading" />
          <span v-if="status === 'hashing'">生成文件hash中 {{ progressInfo.hash }}%</span>
          <span v-else-if="status === 'chunking'">文件切片上传中 {{ progressInfo.chunk }}%</span>
          <span v-else>上传中...</span>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue";
import ChunkUploader from "@/utils/chunk-upload";

defineOptions({
  name: "ChunkUploadView",
});

type UploadStatus = "" | "hashing" | "chunking" | "done" | "error";

let uploader: ChunkUploader;

const fileInput = ref<HTMLInputElement | null>(null);
const fileInfo = ref({
  name: "",
  size: 0,
});
const progressInfo = ref({
  hash: 0,
  chunk: 0,
});
const status = ref<UploadStatus>("");
const fileURL = ref("");
const errorMessage = ref("");

const fileSizeText = computed(() => {
  const size = fileInfo.value.size;
  if (!size) {
    return "";
  }
  if (size >= 1024 * 1024 * 1024) {
    return `${(size / 1024 / 1024 / 1024).toFixed(2)}GB`;
  }
  if (size >= 1024 * 1024) {
    return `${(size / 1024 / 1024).toFixed(2)}MB`;
  }
  return `${(size / 1024).toFixed(2)}KB`;
});

function checkFile(file: File) {
  if (file.size > 10 * 1024 * 1024 * 1024) {
    alert("上传文件不超过10G");
    return false;
  }
  if (file.name.replace(/.+\./, "") !== "mp4") {
    alert("上传文件仅支持MP4格式");
    return false;
  }
  return true;
}

function reset() {
  progressInfo.value = {
    hash: 0,
    chunk: 0,
  };
  fileURL.value = "";
  errorMessage.value = "";
}

function onChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (fileInput.value) {
    fileInput.value.value = "";
  }
  if (!file) {
    return;
  }
  if (!checkFile(file)) {
    return;
  }
  reset();
  fileInfo.value = {
    name: file.name,
    size: file.size,
  };
  status.value = "hashing";
  uploader.start(file);
}

onBeforeMount(() => {
  uploader = new ChunkUploader({
    onHashProgress: (percent) => {
      progressInfo.value.hash = percent;
      status.value = "hashing";
    },
    onProgress: (percent) => {
      progressInfo.value.chunk = percent;
      status.value = "chunking";
    },
    onSuccess: (res) => {
      progressInfo.value.chunk = 100;
      status.value = "done";
      fileURL.value = res.url;
    },
    onError: (err) => {
      status.value = "error";
      errorMessage.value = err instanceof Error ? err.message : "请稍后重试";
      console.error("上传失败", err);
    },
  });
});
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
.text-error {
  color: var(--td-error-color, #d54941);
}
</style>
