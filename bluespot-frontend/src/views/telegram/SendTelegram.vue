<template>
  <main class="send-telegram-view">
    <section class="send-telegram-panel">
      <div class="page-heading">
        <div>
          <h1 class="page-title">发送 Telegram 消息</h1>
          <p class="page-desc">填写消息内容，通过后端 /telegram 接口发送到已配置的会话。</p>
        </div>
      </div>

      <t-form
        class="telegram-form"
        :data="formData"
        :rules="rules"
        label-align="top"
        required-mark
        @submit="handleSubmit"
      >
        <t-form-item label="消息内容" name="text">
          <t-textarea
            v-model="formData.text"
            placeholder="请输入要发送的 Telegram 消息"
            :autosize="{ minRows: 8, maxRows: 14 }"
            :disabled="submitting"
          />
        </t-form-item>

        <div class="form-actions">
          <t-button theme="default" variant="base" :disabled="submitting" @click="resetForm">
            重置
          </t-button>
          <t-button theme="primary" type="submit" :loading="submitting">发送消息</t-button>
        </div>
      </t-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { sendTelegramMessage } from './service'
import type { FormProps, SubmitContext } from 'tdesign-vue-next'

interface SendTelegramForm {
  text: string
}

defineOptions({
  name: 'SendTelegram',
})

const initialFormData: SendTelegramForm = {
  text: '',
}

const formData = reactive<SendTelegramForm>({ ...initialFormData })
const submitting = ref(false)

const rules: FormProps<SendTelegramForm>['rules'] = {
  text: [
    {
      validator: (value: string) => value.trim().length > 0,
      message: '请输入消息内容',
      trigger: 'blur',
    },
  ],
}

function resetForm() {
  Object.assign(formData, initialFormData)
}

async function handleSubmit(context: SubmitContext<SendTelegramForm>) {
  if (context.validateResult !== true || submitting.value) {
    return
  }

  submitting.value = true
  try {
    await sendTelegramMessage({
      text: formData.text.trim(),
    })
    MessagePlugin.success('Telegram 消息发送成功')
    resetForm()
  } catch {
    // 统一错误提示由 request 封装处理。
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.send-telegram-view {
  min-height: calc(100vh - 110px);
  padding: 0 32px 32px;
  background: #f5f7fa;
}

.send-telegram-panel {
  max-width: 760px;
  padding: 24px;
  background: #ffffff;
  border: 1px solid #e7e7e7;
  border-radius: 8px;
}

.page-heading {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
}

.page-title {
  margin: 0;
  color: #1f2329;
  font-size: 22px;
  font-weight: 600;
  line-height: 30px;
}

.page-desc {
  margin: 8px 0 0;
  color: #646a73;
  font-size: 14px;
  line-height: 22px;
}

.telegram-form {
  max-width: 640px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
}

@media (max-width: 768px) {
  .send-telegram-view {
    padding: 0 16px 24px;
  }

  .send-telegram-panel {
    padding: 20px;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .form-actions :deep(.t-button) {
    width: 100%;
  }
}
</style>
