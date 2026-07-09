<template>
  <svg aria-hidden="true" class="svg-icon" :class="className" :style="sizeStyle">
    <use :xlink:href="symbolId" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({
  name: 'ComSvgIcon',
})

interface Props {
  name: string
  prefix?: string
  className?: string
  size?: number | string
}
const props = withDefaults(defineProps<Props>(), {
  prefix: 'icon',
  className: '',
  size: '',
})

const symbolId = computed(() => `#${props.prefix}-${props.name}`)

const sizeStyle = computed(() => {
  if (!props.size) return {}

  const size = typeof props.size === 'number' ? `${props.size}px` : props.size
  return {
    width: size,
    height: size,
  }
})
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  fill: currentColor;
  vertical-align: middle;
}
</style>
