import { defineStore } from 'pinia'
import { cloneDeep } from 'lodash-es'
import type { RouteLocationNormalizedLoadedGeneric } from 'vue-router'

export interface ITabNavItem {
  path: string
  meta: {
    title: string
  }
  isHome?: boolean
  query?: RouteLocationNormalizedLoadedGeneric['query']
  _id?: string | number
}

export const homeRoute: ITabNavItem = {
  path: '/home',
  meta: {
    title: '首页',
  },
  isHome: true,
  _id: '/home',
}

export const useAppStore = defineStore('app', {
  state: () => ({
    collapse: false,
    tabNavList: [homeRoute],
    refreshFlag: true,
  }),
  getters: {
    expand: (state) => !state.collapse,
  },
  actions: {
    toggleCollapse() {
      this.collapse = !this.collapse
    },
    toggleRefreshFlag() {
      this.refreshFlag = false
      setTimeout(() => {
        this.refreshFlag = true
      }, 20)
    },
    addTabNav(row: ITabNavItem) {
      if (!this.tabNavList.some((item) => item._id === row._id) && row.path !== '/home') {
        this.tabNavList.push(cloneDeep(row))
      }
    },
    deleteTabNav(row: ITabNavItem) {
      this.tabNavList = this.tabNavList.filter((item) => item._id !== row._id)
    },
    setTabNavList(list: ITabNavItem[]) {
      this.tabNavList = list
    },
  },
})
