<template>
  <div class="basic-layout flex">
    <div class="aside-placeholder" :class="{ collapsed }" />
    <div class="basic-aside">
      <t-menu
        :value="currentMenu"
        :expanded="expandedMenu"
        theme="light"
        default-value="3-2"
        :collapsed="collapsed"
        @expand="onExpanded"
      >
        <MenuItem :menus="menus" />
        <template #operations>
          <t-button variant="text" shape="square" @click="onChangeCollapsed">
            <template #icon>
              <t-icon name="view-list" />
            </template>
          </t-button>
        </template>
      </t-menu>
    </div>
    <div class="basic-main flex-auto min-w-0">
      <div class="header-placeholder" />
      <div class="basic-header">
        <div class="header-content">
          <router-link class="flex items-center gap-2" to="/" title="bluespot">
            <img src="@/assets/logo.svg" alt="logo" width="48" />
            <h1 class="font-bold text-2xl">bluespot</h1>
          </router-link>
          <div class="header-user">
            <t-dropdown>
              <div class="user-info">
                <t-icon name="user-circle" size="20px" />
                <span class="username">{{ userInfo?.username || '用户' }}</span>
                <t-icon name="chevron-down" size="16px" />
              </div>
              <t-dropdown-menu>
                <t-dropdown-item @click="handleLogout">
                  <t-icon name="logout" />
                  退出登录
                </t-dropdown-item>
              </t-dropdown-menu>
            </t-dropdown>
          </div>
        </div>
      </div>
      <!-- <NavTab></NavTab> -->
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component
            :is="Component"
            :key="route.query.id ? route.path + route.query.id : route.path"
          />
        </keep-alive>
      </router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { difference, remove, union } from 'lodash-es'
import { MessagePlugin } from 'tdesign-vue-next'
import type { MenuValue } from 'tdesign-vue-next'
import MenuItem from './components/MenuItem.vue'
import { menuList } from './components/menu'
import { useUserStore } from '@/stores/user'
import { getActive } from '@/router'
// import NavTab from './components/NavTab.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const collapsed = ref(false)
const menus = ref(menuList)

// 获取用户信息
const userInfo = computed(() => userStore.userInfo)

defineOptions({
  name: 'BasicLayout',
})

const currentMenu = computed(() => getActive())

const expandedMenu = ref<MenuValue[]>([])

watch(
  () => currentMenu.value,
  () => {
    const path = getActive()
    const parentPath = path.substring(0, path.lastIndexOf('/'))
    expandedMenu.value = union([parentPath], expandedMenu.value)
  },
)

const onExpanded = (value: MenuValue[]) => {
  const currentOperationMenu = difference(expandedMenu.value, value)
  const allExpanded = union(value, expandedMenu.value)
  remove(allExpanded, (item) => currentOperationMenu.includes(item))
  expandedMenu.value = allExpanded
}

const onChangeCollapsed = () => {
  collapsed.value = !collapsed.value
}

// 处理登出
const handleLogout = () => {
  userStore.logout()
  MessagePlugin.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.basic-aside {
  /* width: fit-content;
  position: relative;
  transition: all .2s; */
  z-index: 100;
  position: fixed;
  height: calc(100vh - 60px);
  inset-block-start: 60px;
  border-inline-end: 1px solid rgba(5, 5, 5, 0.06);
  :deep(.t-default-menu__inner .t-menu) {
    --td-comp-size-l: 40px;
    padding: 8px;
  }
}
.aside-placeholder {
  width: 232px;
  flex: 0 0 232px;
  transition: all 0.28s cubic-bezier(0.645, 0.045, 0.355, 1);
  &.collapsed {
    width: 64px;
    flex: 0 0 64px;
  }
}
.basic-header {
  z-index: 100;
  position: fixed;
  inset-block-start: 0;
  inset-inline-end: 0;
  width: 100%;
  height: 60px;
  background-color: rgba(255, 255, 255, 0.6);
  border-block-end: 1px solid rgba(5, 5, 5, 0.06);
  backdrop-filter: blur(8px);
  transition: background-color 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  display: flex;
  align-items: center;
  padding: 0 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.header-user {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.user-info:hover {
  background-color: rgba(0, 0, 0, 0.04);
}

.username {
  font-size: 14px;
  color: #333;
}
.header-placeholder {
  height: 110px;
  background-color: transparent;
}
</style>
