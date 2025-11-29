<script setup>
import { NButton, NResult } from 'naive-ui'
import AppPage from '@/components/common/AppPage.vue'
</script>

<template>
  <AppPage>
    <NResult m-auto status="404" description="Sorry, the page you visited does not exist.">
      <template #icon>
        <img src="/image/404.webp" width="500">
      </template>
      <template #footer>
        <NButton @click="$router.replace('/')">
          Back Home
        </NButton>
      </template>
    </NResult>
  </AppPage>
</template>
