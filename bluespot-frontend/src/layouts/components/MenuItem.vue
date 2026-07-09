<template>
  <template v-for="(item, index) in menus" :key="index">
    <t-submenu v-if="item.children && item.children.length" :value="item.path">
      <template #icon>
        <ComSvgIcon v-if="item.icon" class="menu-item__icon" :name="item.icon" size="20" />
      </template>
      <template #title>
        <span class="ml-2">{{ item.title }}</span>
      </template>
      <MenuItem :menus="item.children" />
    </t-submenu>
    <t-menu-item v-else :value="item.path" :to="item.path" class="menu-item--nodropdown">
      <template #icon>
        <ComSvgIcon v-if="item.icon" class="menu-item__icon" :name="item.icon" size="20" />
      </template>
      <span class="ml-2">{{ item.title }}</span>
    </t-menu-item>
  </template>
</template>

<script setup lang="ts">
import type { IMenuItem } from './menu'

defineOptions({
  name: 'MenuItem',
})

interface Props {
  menus: IMenuItem[]
}

defineProps<Props>()
</script>

<style scoped></style>
