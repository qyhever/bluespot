<template>
  <main class="send-mail-view">
    <section class="send-mail-panel">
      <div class="page-heading">
        <div>
          <h1 class="page-title">发送邮件</h1>
          <p class="page-desc">填写收件人、主题和正文，通过后端 /mail 接口发送邮件。</p>
        </div>
      </div>

      <t-form
        class="mail-form"
        :data="formData"
        :rules="rules"
        label-align="top"
        required-mark
        @submit="handleSubmit"
      >
        <t-form-item label="收件人" name="to">
          <t-input
            v-model="formData.to"
            clearable
            placeholder="name@example.com"
            :disabled="submitting"
          />
        </t-form-item>

        <t-form-item label="邮件主题" name="subject">
          <t-input
            v-model="formData.subject"
            clearable
            placeholder="请输入邮件主题"
            :disabled="submitting"
          />
        </t-form-item>

        <t-form-item label="邮件正文" name="body">
          <t-textarea
            v-model="formData.body"
            placeholder="请输入邮件正文，支持 HTML 内容"
            :autosize="{ minRows: 8, maxRows: 14 }"
            :disabled="submitting"
          />
        </t-form-item>

        <div class="form-actions">
          <t-button theme="default" variant="base" :disabled="submitting" @click="resetForm">
            重置
          </t-button>
          <t-button theme="primary" type="submit" :loading="submitting">发送邮件</t-button>
        </div>
      </t-form>
    </section>
  </main>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { sendMail } from './service'
import type { FormProps, SubmitContext } from 'tdesign-vue-next'

interface SendMailForm {
  to: string
  subject: string
  body: string
}

defineOptions({
  name: 'SendMail',
})

const initialFormData: SendMailForm = {
  to: 'arlong@qyhever.cn',
  subject: 'hi',
  body: '光轮冰棍发起旋毛自来也双式之丸',
}

const formData = reactive<SendMailForm>({ ...initialFormData })
const submitting = ref(false)

const rules: FormProps<SendMailForm>['rules'] = {
  to: [
    { required: true, message: '请输入收件人邮箱', trigger: 'blur' },
    { email: true, message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  subject: [{ required: true, message: '请输入邮件主题', trigger: 'blur' }],
  body: [{ required: true, message: '请输入邮件正文', trigger: 'blur' }],
}

function resetForm() {
  Object.assign(formData, initialFormData)
}

async function handleSubmit(context: SubmitContext<SendMailForm>) {
  if (context.validateResult !== true || submitting.value) {
    return
  }

  submitting.value = true
  try {
    await sendMail({
      to: formData.to.trim(),
      subject: formData.subject.trim(),
      body: formData.body.trim(),
    })
    MessagePlugin.success('邮件发送成功')
    resetForm()
  } catch {
    // 统一错误提示由 request 封装处理。
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.send-mail-view {
  min-height: calc(100vh - 110px);
  padding: 0 32px 32px;
  background: #f5f7fa;
}

.send-mail-panel {
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

.mail-form {
  max-width: 640px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 8px;
}

@media (max-width: 768px) {
  .send-mail-view {
    padding: 0 16px 24px;
  }

  .send-mail-panel {
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
