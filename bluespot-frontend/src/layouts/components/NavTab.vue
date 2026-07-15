<template>
  <div class="tab-nav-con">
    <div
      ref="scrollOuter"
      class="tab-nav-outter"
      @DOMMouseScroll.stop.prevent="onScroll"
      @mousewheel.stop.prevent="onScroll"
    >
      <div ref="scrollBody" class="tab-nav-list" :style="{ left: -leftValue + 'px' }">
        <t-dropdown
          v-for="(tab, index) in appStore.tabNavList"
          :key="tab._id"
          trigger="context-menu"
          :min-column-width="128"
          :popup-props="{
            onVisibleChange: (visible: boolean, ctx: Parameters<OnVisibleChangeWithId>[1]) =>
              onVisibleChange(visible, ctx, tab.path),
            visible: operationId === tab.path,
          }"
          :options="
            [
              { content: '刷新', value: 'refresh' },
              index > 0 ? { content: '关闭左侧', value: 'ahead' } : null,
              index < appStore.tabNavList.length - 1
                ? { content: '关闭右侧', value: 'behind' }
                : null,
              appStore.tabNavList.length > 1 ? { content: '关闭其它', value: 'other' } : null,
            ].filter(isNonNullable)
          "
          @click="(dropdownItem: DropdownOption) => onDropdownClick(dropdownItem, tab)"
        >
          <div
            class="tab-nav-item"
            :class="[
              {
                'is-active': activeId === tab._id,
              },
              'tab-nav-item' + index,
            ]"
            @click="onNavItemClick(tab)"
          >
            <div
              :class="{
                'text-brand': activeId === tab._id,
              }"
            >
              <House v-if="tab.isHome" />
              <template v-else>
                {{ tab.meta.title }}
              </template>
            </div>
            <div v-if="!tab.isHome" class="close-x ml-[8px]" @click.stop="onDelete(tab)">
              <X :size="14" />
            </div>
          </div>
        </t-dropdown>
      </div>
      <template v-if="contentExceed">
        <transition name="slide-left">
          <div
            v-if="leftVisible"
            class="tab-nav-operations operations-left"
            @click="handleScroll(240)"
          >
            <ChevronLeft :size="20" />
          </div>
        </transition>
        <transition name="slide-right">
          <div
            v-if="rightVisible"
            class="tab-nav-operations operations-right"
            @click="handleScroll(-240)"
          >
            <ChevronRight :size="20" />
          </div>
        </transition>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChevronLeft, ChevronRight, X, House } from 'lucide-vue-next'
import type { PopupProps, DropdownOption } from 'tdesign-vue-next'
import { homeRoute, useAppStore } from '@/stores/app'
import type { ITabNavItem } from '@/stores/app'
import { isNonNullable } from '@/utils/type'

defineOptions({
  name: 'TabNav',
})

// const navList = ref(Array(30).fill(null).map(() => ({
//   _id: Math.random().toString(36).slice(2)
// })))
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const leftValue = ref(0)
const contentExceed = ref(false)
const activeId = ref<string | number>()
const operationId = ref('')
const scrollOuter = ref()
const scrollBody = ref()

const leftVisible = computed(() => {
  return leftValue.value !== 0
})

watch(
  () => route,
  (newVal) => {
    // console.log('newVal: ', newVal)
    activeId.value = newVal.query.id ? newVal.path + '__' + newVal.query.id : newVal.path
    let title = newVal.meta?.title || ''
    if (newVal.query.id) {
      // 编辑/详情页面
      title += `-${newVal.query.label}`
    }
    appStore.addTabNav({
      path: newVal.path,
      meta: {
        title,
      },
      query: newVal.query,
      _id: activeId.value,
    } as ITabNavItem)
  },
  {
    immediate: true,
    deep: true,
  },
)

const rightVisible = computed(() => {
  const scrollOuterEL = scrollOuter.value
  const scrollBodyEL = scrollBody.value
  if (scrollOuterEL && scrollBodyEL) {
    return leftValue.value !== scrollBodyEL.offsetWidth - scrollOuterEL.offsetWidth
  }
  return false
})

// 自定义一个 IWheelEvent，默认的WheelEvent缺少属性定义
interface IWheelEvent {
  wheelDelta?: number
  detail?: number
  type: string
}
const onScroll = (event: IWheelEvent) => {
  const type = event.type
  let delta = 0
  if (type === 'DOMMouseScroll' || type === 'mousewheel') {
    delta = event.wheelDelta ? event.wheelDelta : -(event.detail || 0) * 40
  }
  handleScroll(delta)
}

const handleScroll = (offset: number) => {
  const outerWidth = scrollOuter.value.offsetWidth
  const bodyWidth = scrollBody.value.offsetWidth
  if (offset > 0) {
    // 向前滚动
    leftValue.value = Math.max(0, leftValue.value - offset)
  } else {
    // 向后滚动
    if (bodyWidth > outerWidth) {
      // 内容超出
      if (leftValue.value < bodyWidth - outerWidth) {
        leftValue.value = Math.min(leftValue.value - offset, bodyWidth - outerWidth)
      }
    } else {
      // 内容未超出
      leftValue.value = 0
    }
  }
}

const onNavItemClick = (row: ITabNavItem) => {
  if (row._id === activeId.value) {
    return
  }
  activeId.value = row._id
  router.push({
    path: row.path,
    query: row.query,
  })
}

const onDelete = (row: ITabNavItem) => {
  appStore.deleteTabNav(row)
}

type OnVisibleChangeWithId = (
  ...args: [...Parameters<NonNullable<PopupProps['onVisibleChange']>>, id: string]
) => void

const onVisibleChange: OnVisibleChangeWithId = (visible, ctx, id) => {
  if (ctx?.trigger === 'document') {
    operationId.value = ''
  }
  if (visible) {
    operationId.value = id
  }
}

const onDropdownClick = (dropdownItem: DropdownOption, row: ITabNavItem) => {
  console.log('dropdownItem, ctx: ', dropdownItem, row)
  if (dropdownItem.value === 'refresh') {
    if (row.path !== route.path) {
      router.push({
        path: row.path,
        query: row.query,
      })
    }
    appStore.toggleRefreshFlag()
  } else {
    const clickIndex = appStore.tabNavList.findIndex((item) => item.path === row.path)
    if (dropdownItem.value === 'ahead') {
      const retList = appStore.tabNavList.slice(clickIndex)
      appStore.setTabNavList([homeRoute].concat(retList))
    } else if (dropdownItem.value === 'behind') {
      const retList = appStore.tabNavList.slice(0, clickIndex + 1)
      appStore.setTabNavList(retList)
    } else {
      const retList = appStore.tabNavList.filter((item) => item.path === row.path)
      appStore.setTabNavList([homeRoute].concat(retList))
    }
    handleOperationEffect(dropdownItem.value as string, clickIndex)
  }
  operationId.value = ''
}

const handleOperationEffect = (type: string, clickIndex: number) => {
  const currentId = route.path
  const currentIndex = appStore.tabNavList.findIndex((item) => item.path === currentId)
  // 存在三种情况需要刷新当前路由
  // 点击非当前路由的关闭其他、点击非当前路由的关闭左侧且当前路由小于触发路由、点击非当前路由的关闭右侧且当前路由大于触发路由
  const needRefreshRouter =
    (type === 'other' && currentIndex !== clickIndex) ||
    (type === 'ahead' && currentIndex < clickIndex) ||
    (type === 'behind' && currentIndex === -1)
  if (needRefreshRouter) {
    const nextRouteIndex = type === 'behind' ? appStore.tabNavList.length - 1 : 0
    const nextRouter = appStore.tabNavList[nextRouteIndex]
    if (nextRouter) {
      router.push({
        path: nextRouter.path,
        query: nextRouter.query,
      })
    }
  }
}

const onResize = () => {
  const scrollOuterEL = scrollOuter.value
  const scrollBodyEL = scrollBody.value
  if (scrollOuterEL && scrollBodyEL) {
    contentExceed.value = scrollOuterEL.offsetWidth < scrollBodyEL.offsetWidth
    leftValue.value = 0
  }
}

onMounted(() => {
  onResize()
  window.addEventListener('resize', onResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
})
</script>

<style scoped>
.tab-nav-con {
  z-index: 100;
  position: fixed;
  inset-block-start: 60px;
  width: 100%;
  height: 50px;
  background-color: #fff;
}
.tab-nav-outter {
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  border: 1px solid #f3f3f3;
  overflow: hidden;
}
.tab-nav-list {
  position: absolute;
  display: inline-flex;
  white-space: nowrap;
  transition: left 0.2s;
}
.tab-nav-item {
  position: relative;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  padding: 13px 16px;
  line-height: 24px;
  background-color: #fff;
  cursor: pointer;
  color: rgba(0, 0, 0, 0.6);
  & + & {
    border-left: 1px solid #eee;
  }
  &:hover {
    color: var(--td-brand-color, #346fff);
  }
  &.is-active {
    /* background-color: #eee; */
    color: var(--td-brand-color, #346fff);
    background-color: var(--td-brand-color-light, #f2f3ff);
  }
}
.close-x:hover {
  color: rgba(0, 0, 0, 0.8);
}
.tab-nav-operations {
  opacity: 1;
  z-index: 1;
  position: absolute;
  top: 0;
  width: 40px;
  height: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
  /* background-color: #f3f3f3; */
  background-color: #fff;
  border-left: 1px solid #eee;
  cursor: pointer;
  color: rgb(0 0 0 / 60%);
  cursor: pointer;
  transition: all 0.3s ease-in-out;
  &:hover {
    background-color: #eee;
  }
}
.operations-left {
  left: 0;
  border-right: 1px solid #eee;
  box-shadow: 10px 0 20px 5px rgb(0 0 0 / 5%);
}
.operations-right {
  right: 0;
  border-left: 1px solid #eee;
  box-shadow: -10px 0 20px 5px rgb(0 0 0 / 5%);
}
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.3s;
}
.slide-left-enter-from,
.slide-left-leave-to {
  transform: translateX(-40px);
  opacity: 0;
}
.slide-right-enter-from,
.slide-right-leave-to {
  transform: translateX(40px);
  opacity: 0;
}
</style>
