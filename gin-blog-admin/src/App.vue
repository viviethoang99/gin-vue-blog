<script setup>
import { computed, onMounted } from 'vue'
import { NConfigProvider, darkTheme, dateViVN as dateLocale, viVN as locale } from 'naive-ui'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'

import { useAuthStore, useThemeStore } from '@/store'
import themes from '@/assets/themes'
import api from '@/api'

hljs.registerLanguage('json', json)
const themeStore = useThemeStore()
const naiveThemeOverrides = computed(() => themes.naiveThemeOverrides(themeStore.darkMode))

// onMounted(() => {
//   const { accessToken } = useAuthStore()
//   // Report user information when access token exists
//   accessToken && api.report()
// })

// FIXME: After each Docker build, the app may inherit previous localStorage
// TODO: If no route information is found, redirect to the login page
</script>

<template>
  <NConfigProvider
    class="h-full w-full"
    :theme="themeStore.darkMode ? darkTheme : undefined"
    :theme-overrides="naiveThemeOverrides"
    :locale="locale"
    :date-locale="dateLocale"
    :hljs="hljs"
  >
    <RouterView v-slot="{ Component }">
      <component :is="Component" />
    </RouterView>
  </NConfigProvider>
</template>
