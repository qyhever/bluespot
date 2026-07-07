<template>
  <main class="login-page">

    <section class="login-panel" aria-label="登录">
      <div class="login-card">
        <div class="card-head">
          <h2>登录</h2>
        </div>

        <form class="login-form" @submit.prevent="submit">
          <div class="field">
            <label for="account">账号</label>
            <div class="control">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path
                  d="M20 21a8 8 0 0 0-16 0M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Z"
                  stroke="currentColor"
                  stroke-width="1.9"
                  stroke-linecap="round"
                />
              </svg>
              <input v-model="account" id="account" type="text" autocomplete="username" />
            </div>
          </div>

          <div class="field">
            <label for="password">密码</label>
            <div class="control">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                <path
                  d="M7 10V8a5 5 0 0 1 10 0v2M6 10h12a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1v-8a1 1 0 0 1 1-1Z"
                  stroke="currentColor"
                  stroke-width="1.9"
                  stroke-linecap="round"
                />
              </svg>
              <input
                v-model="password"
                id="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
              />
              <button
                class="toggle-pass"
                type="button"
                title="显示或隐藏密码"
                aria-label="显示或隐藏密码"
                @click="showPassword = !showPassword"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
                  <path
                    d="M2.5 12s3.5-6 9.5-6 9.5 6 9.5 6-3.5 6-9.5 6-9.5-6-9.5-6Z"
                    stroke="currentColor"
                    stroke-width="1.8"
                  />
                  <path
                    d="M12 15a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z"
                    stroke="currentColor"
                    stroke-width="1.8"
                  />
                </svg>
              </button>
            </div>
          </div>

          <div class="flex justify-between items-center gap-3">
            <label class="check flex items-center">
              <input v-model="remember" type="checkbox" />
              <span>保持登录</span>
            </label>
            <button class="link-button" type="button" @click="reset">重置输入</button>
          </div>

          <button class="primary" type="submit" :disabled="submitting">
            <span>{{ submitting ? "正在验证" : "登录" }}</span>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M5 12h14m-6-6 6 6-6 6"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>

          <div v-if="notice" class="notice" :class="noticeType">{{ notice }}</div>
        </form>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { login } from "@/api/global";
import { useUserStore } from "@/stores/user";

defineOptions({
  name: "LoginView",
});

const router = useRouter();
const userStore = useUserStore();
const account = ref("");
const password = ref("");
const remember = ref(true);
const showPassword = ref(false);
const submitting = ref(false);
const notice = ref("");
const noticeType = ref("");

function reset() {
  account.value = "";
  password.value = "";
  noticeType.value = "";
  notice.value = "";
}

async function submit() {
  if (submitting.value) return;

  const username = account.value.trim();
  if (!username || !password.value) {
    noticeType.value = "error";
    notice.value = "请输入账号和密码";
    return;
  }

  submitting.value = true;
  noticeType.value = "";
  notice.value = "正在验证...";

  try {
    const tokens = await login({
      username,
      password: password.value,
    });
    userStore.setTokens(tokens);
    void userStore.fetchUserInfo(true).catch((error) => {
      console.warn("Failed to load user info after login.", error);
    });
    if (remember.value) {
      localStorage.setItem("albumSession", "remote");
    } else {
      localStorage.removeItem("albumSession");
    }
    noticeType.value = "success";
    notice.value = "验证通过，正在进入相册";
    await router.push("/");
  } catch (error) {
    console.warn("Failed to login.", error);
    noticeType.value = "error";
    notice.value = "账号或密码不正确";
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped>
.login-page {
  --bg: #f6f5f2;
  --panel: #ffffff;
  --ink: #101013;
  --muted: #68665f;
  --line: #d9d6cf;
  --line-strong: #b9b5aa;
  --dark: #171719;
  --accent: #c43f2f;
  --teal: #1f7a78;
  --good: #18794e;
  --shadow: 0 22px 70px rgba(18, 18, 20, 0.16);
  --radius: 8px;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--ink);
  background:
    linear-gradient(90deg, rgba(16, 16, 19, 0.045) 1px, transparent 1px) 0 0 / 44px 44px,
    linear-gradient(rgba(16, 16, 19, 0.035) 1px, transparent 1px) 0 0 / 44px 44px,
    var(--bg);
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
}

button,
input {
  font: inherit;
}

button {
  cursor: pointer;
}

.login-panel {
  padding: 34px;
  width: 480px;
}

.login-card {
  width: min(100%, 470px);
  border-radius: var(--radius);
  background: rgba(255, 255, 255, 0.82);
  box-shadow: var(--shadow);
  overflow: hidden;
}

.card-head {
  text-align: center;
  padding: 24px 24px 18px;
  border-bottom: 1px solid var(--line);
}

.card-head h2 {
  font-size: 30px;
  line-height: 1;
}

.card-head p {
  margin: 10px 0 0;
  color: var(--muted);
  line-height: 1.55;
}

.login-form {
  padding: 22px 24px 24px;
  display: grid;
  gap: 16px;
}

.field {
  display: grid;
  gap: 8px;
}

.field label {
  color: var(--dark);
  font-size: 13px;
  font-weight: 700;
}

.control {
  position: relative;
}

.control svg {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--muted);
  pointer-events: none;
}

.control input {
  width: 100%;
  height: 48px;
  border-radius: var(--radius);
  border: 1px solid var(--line-strong);
  background: #fff;
  color: var(--ink);
  outline: none;
  padding: 0 48px 0 44px;
}

.control input:focus {
  border-color: var(--dark);
  box-shadow: 0 0 0 3px rgba(16, 16, 19, 0.12);
}

.toggle-pass {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  width: 38px;
  height: 38px;
  border: 0;
  border-radius: 7px;
  color: var(--dark);
  background: transparent;
  display: grid;
  place-items: center;
}

.toggle-pass svg {
  position: static;
  transform: none;
  color: inherit;
}

.toggle-pass:hover {
  background: rgba(16, 16, 19, 0.07);
}

.check {
  gap: 8px;
  color: var(--muted);
  font-size: 13px;
}

.check input {
  width: 16px;
  height: 16px;
  accent-color: var(--dark);
}

.link-button {
  border: 0;
  padding: 0;
  background: transparent;
  color: var(--teal);
  font-size: 13px;
  font-weight: 700;
}

.primary {
  height: 50px;
  border-radius: var(--radius);
  border: 1px solid var(--dark);
  background: var(--dark);
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-weight: 700;
  transition:
    transform 180ms ease,
    background 180ms ease,
    box-shadow 180ms ease;
}

.primary:hover {
  transform: translateY(-1px);
  background: #050506;
  box-shadow: 0 14px 30px rgba(18, 18, 20, 0.2);
}

.primary:disabled {
  cursor: wait;
  opacity: 0.74;
  transform: none;
}

.notice {
  min-height: 42px;
  border-radius: var(--radius);
  border: 1px solid var(--line);
  background: rgba(246, 245, 242, 0.78);
  color: var(--muted);
  padding: 11px 12px;
  font-size: 13px;
  line-height: 1.45;
}

.notice.error {
  color: var(--accent);
  background: rgba(196, 63, 47, 0.08);
  border-color: rgba(196, 63, 47, 0.28);
}

.notice.success {
  color: var(--good);
  background: rgba(24, 121, 78, 0.09);
  border-color: rgba(24, 121, 78, 0.26);
}

.mini-link {
  color: var(--dark);
  text-decoration: none;
  font-weight: 700;
}

@media (max-width: 980px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .brand-panel {
    min-height: auto;
    padding: 26px;
  }

  .brand-content {
    min-height: 460px;
  }

  .login-panel {
    min-height: auto;
    padding: 28px 20px 44px;
  }
}

@media (max-width: 620px) {
  .brand-panel {
    padding: 20px;
  }

  .brand-content {
    min-height: 420px;
  }

  .hero-copy h2 {
    font-size: 40px;
  }

  .system-strip {
    grid-template-columns: 1fr;
  }

  .card-head,
  .login-form,
  .card-foot {
    padding-left: 18px;
    padding-right: 18px;
  }

  .row,
  .card-foot {
    align-items: flex-start;
    display: flex;
    flex-direction: column;
  }
}
</style>
