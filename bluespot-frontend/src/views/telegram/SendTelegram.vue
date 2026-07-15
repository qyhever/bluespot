<template>
  <main class="p-5">
    <section class="bg-white rounded-lg p-6">
      <div class="page-heading">
        <div>
          <h1 class="text-lg font-bold">发送 Telegram 消息</h1>
          <p class="text-sm text-gray-600">填写消息内容，调用后端接口发送到已配置的会话。</p>
        </div>
      </div>

      <div class="w-120 mt-6">

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
  
          <div class="flex justify-end gap-4 mt-6">
            <t-button theme="default" variant="base" :disabled="submitting" @click="resetForm">
              重置
            </t-button>
            <t-button theme="primary" type="submit" :loading="submitting">发送消息</t-button>
          </div>
        </t-form>
      </div>
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

</style>
