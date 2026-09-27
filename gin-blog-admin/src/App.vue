<script setup>
import { darkTheme, dateZhCN, NConfigProvider, zhCN } from 'naive-ui'

import themes from '@/assets/themes'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()

// highlight.js 只有操作日志页的 NCode 用得到, 放在这里会进入首屏 chunk (约 38KB),
// 改为在 views/log/operation 里按需引入并直接传给 NCode。

// 上报用户信息, 需要时取消注释, 并补回 onMounted / api / useAuthStore 的导入
// onMounted(() => {
//   const { accessToken } = useAuthStore()
//   accessToken && api.report()
// })

// FIXME: After each Docker build, the app may inherit previous localStorage
// TODO: If no route information is found, redirect to the login page
</script>

<template>
  <NConfigProvider
    class="h-full w-full"
    :theme="themeStore.darkMode ? darkTheme : undefined"
    :theme-overrides="themes.naiveThemeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
  >
    <RouterView v-slot="{ Component }">
      <component :is="Component" />
    </RouterView>
  </NConfigProvider>
</template>
