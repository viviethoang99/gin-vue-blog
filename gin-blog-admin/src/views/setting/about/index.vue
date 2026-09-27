<script setup>
import { MdEditor } from 'md-editor-v3'
import { NButton } from 'naive-ui'

import { onMounted, ref } from 'vue'
import api from '@/api'

import CommonPage from '@/components/common/CommonPage.vue'
import 'md-editor-v3/lib/style.css'

defineOptions({ name: 'About Me' })

const aboutContent = ref('')
const btnLoading = ref(false)

onMounted(async () => {
  // 裸 await 会在接口失败时留下未捕获的 rejection
  try {
    const resp = await api.getAbout()
    aboutContent.value = resp.data ?? ''
  }
  catch (err) {
    console.error(err)
  }
})

async function handleSave() {
  try {
    btnLoading.value = true
    await api.updateAbout({ content: aboutContent.value })
    window.$message.success('Update successful')
  }
  catch (err) {
    console.error(err)
  }
  finally {
    btnLoading.value = false
  }
}
</script>

<template>
  <CommonPage :show-header="false">
    <div class="mb-4 flex items-center justify-between">
      <div class="mx-1 text-2xl font-bold">
        About Me
      </div>
      <NButton type="primary" :loading="btnLoading" @click="handleSave">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:confirm-circle" />
        </template>
        Save
      </NButton>
    </div>
    <!-- TODO: File upload encapsulation -->
    <MdEditor v-model="aboutContent" style="height: calc(100vh - 245px)" />
  </CommonPage>
</template>

<style lang="scss">
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
