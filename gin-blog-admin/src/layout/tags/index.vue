<script setup>
import { nextTick, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NTag } from 'naive-ui'

import ContextMenu from './ContextMenu.vue'
import TheIcon from '@/components/icon/TheIcon.vue'
import ScrollX from '@/components/common/ScrollX.vue'
import { useTagStore } from '@/store'

const route = useRoute()
const router = useRouter()
const tagStore = useTagStore()

const scrollXRef = ref(null)
const tabRefs = ref([])

const contextMenuOption = reactive({
  show: false,
  x: 0,
  y: 0,
  currentPath: '',
})

// Watch current route path, add to tab bar when changed
watch(
  () => route.path,
  () => {
    const { name, fullPath: path } = route
    const title = route.meta?.title
    const icon = route.meta?.icon
    tagStore.addTag({ name, path, title, icon })
  },
  { immediate: true },
)

// Watch current active tag, scroll to make it visible
watch(
  () => tagStore.activeIndex,
  async (activeIndex) => {
    await nextTick()
    const activeTabElement = tabRefs.value[activeIndex]?.$el
    if (activeTabElement) {
      const { offsetLeft: x, offsetWidth: width } = activeTabElement
      scrollXRef.value?.handleScroll(x + width, width)
    }
  },
  { immediate: true },
)

function handleTagClick(path) {
  tagStore.setActiveTag(path) // Activate current clicked tag
  router.push(path)
}

// Show or hide right-click menu
function setContextMenuShow(flag) {
  contextMenuOption.show = flag
}

function setContextMenu(x, y, currentPath) {
  // Object.assign(a, b) copies properties of b to a (overwrites if same), shallow copy
  Object.assign(contextMenuOption, { x, y, currentPath })
}

// Right-click menu
async function handleContextMenu(e, tagItem) {
  const { clientX, clientY } = e
  setContextMenuShow(false)
  setContextMenu(clientX, clientY, tagItem.path)
  await nextTick()
  setContextMenuShow(true)
}

function handleRefresh(tag) {
  // Only current tag will refresh
  if (route.name === tag.name) {
    tagStore.updateAliveKey(route.name)
    tagStore.reloadTag()
  }
}
</script>

<template>
  <ScrollX ref="scrollXRef" class="bg-white dark:bg-dark!">
    <NTag
      v-for="tag in tagStore.tags" :key="tag.path"
      ref="tabRefs"
      class="mx-1 hover:border-blue hover:border-red hover:text-primary"
      :type="tagStore.activeTag === tag.path ? 'primary' : 'default'"
      :closable="tagStore.tags.length > 1"
      @click="handleTagClick(tag.path)"
      @close.stop="tagStore.removeTag(tag.path)"
      @contextmenu.prevent="handleContextMenu($event, tag)"
    >
      <template #icon>
        <div :class="{ 'cursor-pointer': $route.name === tag.name }" @click="handleRefresh(tag)">
          <TheIcon v-if="tag.icon" :icon="tag.icon" :size="16" />
          <i v-else class="i-mdi:refresh" />
        </div>
      </template>
      <div class="px-0.5">
        {{ tag.title }}
      </div>
    </NTag>
    <ContextMenu
      v-if="contextMenuOption.show"
      v-model:show="contextMenuOption.show"
      :current-path="contextMenuOption.currentPath"
      :x="contextMenuOption.x"
      :y="contextMenuOption.y"
    />
  </ScrollX>
</template>
