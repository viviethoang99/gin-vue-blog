<script setup>
import { NDropdown } from 'naive-ui'
import { computed, h } from 'vue'
import { useRoute } from 'vue-router'

import { useTagStore } from '@/store'

const props = defineProps({
  show: { type: Boolean, default: false },
  currentPath: { type: String, default: '' },
  x: { type: Number, default: 0 },
  y: { type: Number, default: 0 },
})

const emit = defineEmits(['update:show'])

const route = useRoute()
const tagStore = useTagStore()

const options = computed(() => [
  {
    label: 'Reload',
    key: 'reload',
    disabled: props.currentPath !== tagStore.activeTag, // Can only reload current tag
    icon: () => h('i', { class: 'i-mdi:refresh' }),
  },
  {
    label: 'Close',
    key: 'close',
    disabled: tagStore.tags.length <= 1, // Cannot close when only one tag exists
    icon: () => h('i', { class: 'i-mdi:close' }),
  },
  {
    label: 'Close Others',
    key: 'close-other',
    disabled: tagStore.tags.length <= 1, // Cannot close others when only one tag exists
    icon: () => h('i', { class: 'i-mdi:arrow-expand-horizontal' }),
  },
  {
    label: 'Close Left',
    key: 'close-left',
    // Cannot close left when only one tag or current selected is the first tag
    disabled: tagStore.tags.length <= 1 || props.currentPath === tagStore.tags[0].path,
    icon: () => h('i', { class: 'i-mdi:arrow-expand-left' }),
  },
  {
    label: 'Close Right',
    key: 'close-right',
    // Cannot close right when only one tag or current selected is the last tag
    disabled: tagStore.tags.length <= 1 || props.currentPath === tagStore.tags[tagStore.tags.length - 1].path,
    icon: () => h('i', { class: 'i-mdi:arrow-expand-right' }),
  },
])

const actionMap = new Map([
  [
    'reload',
    () => {
      // Reload, regardless of keepAlive, need to re-fetch data
      tagStore.updateAliveKey(route.name)
      tagStore.reloadTag()
    },
  ],
  [
    'close',
    () => {
      // resetKeepAlive()
      tagStore.removeTag(props.currentPath)
    },
  ],
  [
    'close-other',
    () => tagStore.removeOther(props.currentPath),
  ],
  [
    'close-left',
    () => tagStore.removeLeft(props.currentPath),
  ],
  [
    'close-right',
    () => tagStore.removeRight(props.currentPath),
  ],
])

function handleHideDropdown() {
  emit('update:show', false)
}

function handleSelect(key) {
  const actionFn = actionMap.get(key)
  actionFn && actionFn()
  handleHideDropdown()
}
</script>

<template>
  <NDropdown
    :show="show"
    :options="options"
    :x="x"
    :y="y"
    placement="bottom-start"
    @clickoutside="handleHideDropdown"
    @select="handleSelect"
  />
</template>
