<template>
  <main class="chunk-upload-view">
    <section class="upload-panel" aria-labelledby="chunk-upload-title">
      <div class="page-head">
        <div>
          <p class="page-kicker">附件管理</p>
          <h1 id="chunk-upload-title" class="page-title">大文件分片上传</h1>
          <p class="page-description">
            选择 MP4 视频后自动计算文件指纹并按分片上传，适用于较大视频附件。
          </p>
        </div>
        <div class="upload-badge">
          <t-icon name="cloud-upload" />
          <span>Chunk Upload</span>
        </div>
      </div>

      <div class="upload-layout">
        <label
          class="upload-dropzone"
          :class="{ 'is-uploading': isUploading }"
          :aria-disabled="isUploading"
        >
          <input
            ref="fileInput"
            type="file"
            class="input-file"
            accept=".mp4,video/mp4"
            :disabled="isUploading"
            @change="onChange"
          />
          <span class="dropzone-icon" aria-hidden="true">
            <t-icon name="upload" />
          </span>
          <span class="dropzone-title">
            {{ isUploading ? "文件上传处理中" : "选择 MP4 视频文件" }}
          </span>
          <span class="dropzone-desc">
            {{
              isUploading ? "请等待当前任务完成后再选择新文件" : "点击此区域选择文件，最大支持 10GB"
            }}
          </span>
          <span class="upload-action">浏览文件</span>
        </label>

        <aside class="rules-panel" aria-label="上传规则">
          <h2 class="section-title">上传规则</h2>
          <ul class="rule-list">
            <li>
              <t-icon name="check-circle" />
              <span>仅支持 MP4 视频格式</span>
            </li>
            <li>
              <t-icon name="check-circle" />
              <span>单文件大小不超过 10GB</span>
            </li>
            <li>
              <t-icon name="check-circle" />
              <span>上传前会先生成文件 hash</span>
            </li>
          </ul>
        </aside>
      </div>

      <section class="status-panel" aria-live="polite">
        <div class="status-head">
          <div>
            <h2 class="section-title">上传状态</h2>
            <p class="status-subtitle">{{ statusMeta.description }}</p>
          </div>
          <span class="status-pill" :class="`is-${statusMeta.tone}`">
            <t-icon :name="statusMeta.icon" />
            <span>{{ statusMeta.label }}</span>
          </span>
        </div>

        <div
          class="progress-track"
          role="progressbar"
          :aria-valuenow="uploadProgress"
          aria-valuemin="0"
          aria-valuemax="100"
        >
          <div
            class="progress-bar"
            :class="`is-${statusMeta.tone}`"
            :style="{ width: `${uploadProgress}%` }"
          />
        </div>
        <div class="progress-meta">
          <span>Hash {{ progressInfo.hash }}%</span>
          <span>分片 {{ progressInfo.chunk }}%</span>
        </div>

        <dl class="file-detail">
          <div class="detail-item">
            <dt>文件名称</dt>
            <dd>{{ fileInfo.name || "尚未选择文件" }}</dd>
          </div>
          <div class="detail-item">
            <dt>文件大小</dt>
            <dd>{{ fileSizeText || "-" }}</dd>
          </div>
          <div v-if="status === 'done'" class="detail-item is-success">
            <dt>文件路径</dt>
            <dd>
              <a :href="fileURL" target="_blank" rel="noreferrer">{{ fileURL }}</a>
            </dd>
          </div>
          <div v-else-if="status === 'error'" class="detail-item is-error">
            <dt>失败原因</dt>
            <dd>{{ errorMessage }}</dd>
          </div>
        </dl>
      </section>
    </section>
  </main>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
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

const isUploading = computed(() => {
  return status.value === "hashing" || status.value === "chunking";
});

const uploadProgress = computed(() => {
  if (status.value === "done") {
    return 100;
  }
  if (status.value === "hashing") {
    return progressInfo.value.hash;
  }
  if (status.value === "chunking") {
    return progressInfo.value.chunk;
  }
  return 0;
});

const statusMeta = computed(() => {
  if (status.value === "done") {
    return {
      label: "上传完成",
      description: "文件已合并完成，可通过下方路径访问。",
      icon: "check-circle",
      tone: "success",
    };
  }
  if (status.value === "error") {
    return {
      label: "上传失败",
      description: "请检查文件或网络状态后重新上传。",
      icon: "error-circle",
      tone: "error",
    };
  }
  if (status.value === "hashing") {
    return {
      label: "计算 hash",
      description: `正在生成文件指纹，进度 ${progressInfo.value.hash}%。`,
      icon: "loading",
      tone: "processing",
    };
  }
  if (status.value === "chunking") {
    return {
      label: "分片上传",
      description: `正在上传文件分片，进度 ${progressInfo.value.chunk}%。`,
      icon: "loading",
      tone: "processing",
    };
  }
  return {
    label: "等待文件",
    description: "选择文件后会自动开始上传任务。",
    icon: "info-circle",
    tone: "neutral",
  };
});

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
    MessagePlugin.warning("上传文件不超过10G");
    return false;
  }
  if (file.name.replace(/.+\./, "") !== "mp4") {
    MessagePlugin.warning("上传文件仅支持MP4格式");
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
.chunk-upload-view {
  min-height: calc(100vh - 60px);
  padding: 32px;
  background:
    linear-gradient(180deg, rgba(99, 102, 241, 0.07), rgba(99, 102, 241, 0) 240px), #f5f7fa;
}

.upload-panel {
  max-width: 1080px;
  padding: 28px;
  background: #ffffff;
  border: 1px solid #e7eaf3;
  border-radius: 8px;
}

.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 24px;
}

.page-kicker {
  margin: 0 0 6px;
  color: #6366f1;
  font-size: 13px;
  font-weight: 600;
  line-height: 20px;
}

.page-title {
  margin: 0;
  color: #1e1b4b;
  font-size: 24px;
  font-weight: 650;
  line-height: 34px;
}

.page-description {
  max-width: 620px;
  margin: 8px 0 0;
  color: #646a73;
  font-size: 14px;
  line-height: 22px;
}

.upload-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex: 0 0 auto;
  padding: 8px 12px;
  color: #3730a3;
  font-size: 13px;
  font-weight: 600;
  line-height: 20px;
  background: #eef2ff;
  border: 1px solid #dfe3ff;
  border-radius: 999px;
}

.upload-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 20px;
  align-items: stretch;
}

.upload-dropzone {
  position: relative;
  display: flex;
  min-height: 260px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  background: #f8fafc;
  border: 1px dashed #9aa6ff;
  border-radius: 8px;
  transition:
    border-color 0.18s ease,
    background-color 0.18s ease;
}

.upload-dropzone:hover {
  background: #f4f7ff;
  border-color: #6366f1;
}

.upload-dropzone.is-uploading {
  cursor: not-allowed;
  background: #f8fafc;
  border-color: #cbd5e1;
}

.dropzone-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  margin-bottom: 16px;
  color: #ffffff;
  font-size: 28px;
  background: #6366f1;
  border-radius: 8px;
}

.dropzone-title {
  color: #1f2329;
  font-size: 18px;
  font-weight: 650;
  line-height: 28px;
}

.dropzone-desc {
  max-width: 360px;
  margin-top: 6px;
  color: #646a73;
  font-size: 14px;
  line-height: 22px;
}

.upload-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 112px;
  height: 36px;
  margin-top: 20px;
  padding: 0 18px;
  color: #ffffff;
  font-size: 14px;
  font-weight: 600;
  background: #10b981;
  border-radius: 6px;
  transition: background-color 0.18s ease;
}

.upload-dropzone:hover .upload-action {
  background: #059669;
}

.input-file {
  opacity: 0;
  position: absolute;
  inset: 0;
  height: 100%;
  width: 100%;
  cursor: pointer;
}

.input-file:disabled {
  cursor: not-allowed;
}

.rules-panel {
  padding: 20px;
  background: #fbfcff;
  border: 1px solid #e7eaf3;
  border-radius: 8px;
}

.section-title {
  margin: 0;
  color: #1f2329;
  font-size: 16px;
  font-weight: 650;
  line-height: 24px;
}

.rule-list {
  display: grid;
  gap: 14px;
  margin: 18px 0 0;
  padding: 0;
  list-style: none;
}

.rule-list li {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  color: #4e5969;
  font-size: 14px;
  line-height: 22px;
}

.rule-list .t-icon {
  flex: 0 0 auto;
  margin-top: 3px;
  color: #10b981;
}

.status-panel {
  margin-top: 20px;
  padding: 20px;
  background: #ffffff;
  border: 1px solid #e7eaf3;
  border-radius: 8px;
}

.status-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.status-subtitle {
  margin: 4px 0 0;
  color: #646a73;
  font-size: 14px;
  line-height: 22px;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex: 0 0 auto;
  padding: 6px 10px;
  font-size: 13px;
  font-weight: 600;
  line-height: 20px;
  border-radius: 999px;
}

.status-pill.is-neutral {
  color: #4e5969;
  background: #f2f3f5;
}

.status-pill.is-processing {
  color: #3730a3;
  background: #eef2ff;
}

.status-pill.is-success {
  color: #047857;
  background: #ecfdf5;
}

.status-pill.is-error {
  color: #b91c1c;
  background: #fef2f2;
}

.status-pill.is-processing .t-icon {
  animation: spin 1s linear infinite;
}

.progress-track {
  width: 100%;
  height: 10px;
  overflow: hidden;
  background: #edf0f7;
  border-radius: 999px;
}

.progress-bar {
  height: 100%;
  background: #6366f1;
  border-radius: inherit;
  transition: width 0.2s ease;
}

.progress-bar.is-success {
  background: #10b981;
}

.progress-bar.is-error {
  background: #dc2626;
}

.progress-meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 8px;
  color: #86909c;
  font-size: 12px;
  line-height: 20px;
}

.file-detail {
  display: grid;
  gap: 10px;
  margin: 18px 0 0;
}

.detail-item {
  display: grid;
  grid-template-columns: 96px minmax(0, 1fr);
  gap: 12px;
  color: #1f2329;
  font-size: 14px;
  line-height: 22px;
}

.detail-item dt {
  color: #646a73;
}

.detail-item dd {
  min-width: 0;
  margin: 0;
  overflow-wrap: anywhere;
}

.detail-item a {
  color: #0052d9;
  text-decoration: none;
}

.detail-item a:hover {
  text-decoration: underline;
}

.detail-item.is-success dd {
  color: #047857;
}

.detail-item.is-error dd,
.text-error {
  color: var(--td-error-color, #d54941);
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .status-pill.is-processing .t-icon {
    animation: none;
  }

  .progress-bar,
  .upload-dropzone,
  .upload-action {
    transition: none;
  }
}

@media (max-width: 900px) {
  .upload-layout {
    grid-template-columns: 1fr;
  }

  .rules-panel {
    padding: 18px;
  }
}

@media (max-width: 640px) {
  .chunk-upload-view {
    padding: 16px;
  }

  .upload-panel {
    padding: 18px;
  }

  .page-head,
  .status-head {
    flex-direction: column;
  }

  .upload-badge,
  .status-pill {
    align-self: flex-start;
  }

  .upload-dropzone {
    min-height: 220px;
    padding: 24px 18px;
  }

  .detail-item {
    grid-template-columns: 1fr;
    gap: 2px;
  }
}
</style>
