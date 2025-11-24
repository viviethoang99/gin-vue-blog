<script setup>
import { NResult } from 'naive-ui'
import CommonPage from '@/components/common/CommonPage.vue'
</script>

<template>
  <CommonPage>
    <NResult
      class="mt-30"
      status="404"
      title="This page is still under development..."
      description="稍安毋躁"
      size="huge"
    />
  </CommonPage>
</template>
