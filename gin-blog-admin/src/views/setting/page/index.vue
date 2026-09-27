<script setup>
import { NButton, NEmpty, NForm, NFormItem, NImage, NInput } from 'naive-ui'
import { onMounted, ref } from 'vue'

import api from '@/api'
import UploadOne from '@/components//UploadOne.vue'
import CommonPage from '@/components/common/CommonPage.vue'

import CrudModal from '@/components/crud/CrudModal.vue'
import { useCRUD } from '@/composables'
import { convertImgUrl } from '@/utils'

// FIXME: Why doesn't KeepAlive work for this page only?

const {
  modalVisible,
  modalTitle,
  modalLoading,
  handleAdd,
  handleDelete,
  handleEdit,
  handleSave,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: 'Page',
  initForm: {},
  doCreate: api.saveOrUpdatePage,
  doDelete: api.deletePage,
  doUpdate: api.saveOrUpdatePage,
  refresh: fetchData,
})

const pageList = ref([])
// 首屏拉取期间不显示空状态, 否则会闪一下「还没有页面」
const loading = ref(true)
const reloadFlag = ref(false)
const uploadOneRef = ref(null) // Image upload ref object

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  // 裸 await 会在接口失败时产生未捕获的 rejection; data 为空时也不能让 pageList 变成 undefined
  loading.value = true
  try {
    const resp = await api.getPages()
    pageList.value = resp.data ?? []
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
}

// Refresh preview image based on input link
function refreshImg(img) {
  // 弹窗未打开时 UploadOne 还没挂载, 直接取 .previewImg 会抛异常
  if (!uploadOneRef.value) {
    return
  }
  reloadFlag.value = true
  uploadOneRef.value.previewImg = img
  setTimeout(() => reloadFlag.value = false, 600)
}
</script>

<template>
  <CommonPage title="Page Management">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        New Page
      </NButton>
    </template>
    <!-- 原来是 flex + justify-between, 末尾还要垫三个空 div 才对得齐;
         换成 grid, 列数按容器宽度自适应, 卡片数变化也不会跳 -->
    <div class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(300px, 1fr))">
      <div
        v-for="page of pageList" :key="page.id"
        class="relative my-2 cursor-pointer text-center"
      >
        <div class="absolute right-2 top-1 flex gap-1">
          <!-- 直接放两个按钮, 比藏在下拉里的「···」显眼 -->
          <button
            class="h-7 w-7 flex cursor-pointer items-center justify-center rounded-full bg-black/50 text-white transition-300 hover:bg-blue"
            title="Edit"
            @click="handleEdit(page)"
          >
            <i class="i-mingcute:edit-2-line text-base" />
          </button>
          <button
            class="h-7 w-7 flex cursor-pointer items-center justify-center rounded-full bg-black/50 text-white transition-300 hover:bg-red"
            title="Delete"
            @click="handleDelete([page.id])"
          >
            <i class="i-mingcute:delete-back-line text-base" />
          </button>
        </div>
        <NImage
          :src="convertImgUrl(page.cover)"
          height="170" width="300"
          :img-props="{ style: { 'border-radius': '5px' } }"
        />
        <p class="text-base">
          {{ page.name }}
        </p>
      </div>
    </div>

    <NEmpty v-if="!loading && !pageList.length" class="py-16" description="还没有页面, 点右上角新建" />

    <CrudModal
      v-model:visible="modalVisible"
      width="550px"
      :title="modalTitle"
      :loading="modalLoading"
      @save="handleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <NFormItem
          label="Page name"
          path="name"
          :rule="{ required: true, message: 'Please enter page name', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.name" placeholder="Page name" />
        </NFormItem>
        <NFormItem
          label="Page label"
          path="label"
          :rule="{ required: true, message: 'Please enter page label', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.label" placeholder="Page label" />
        </NFormItem>
        <NFormItem
          label="Cover"
          path="cover"
          :rule="{ required: true, message: 'Please upload cover image', trigger: ['input', 'blur'] }"
        >
          <div class="w-full flex items-center justify-between">
            <UploadOne
              ref="uploadOneRef"
              v-model:preview="modalForm.cover"
              :width="300"
              @finish="val => (modalForm.cover = val)"
            />

            <span
              class="i-uiw:reload h-5 w-5 cursor-pointer"
              :class="reloadFlag ? 'animate-spin' : ''"
              @click="refreshImg(modalForm.cover)"
            />
          </div>
        </NFormItem>
        <NFormItem
          label="URL"
          path="cover"
          :rule="{ required: true, message: 'Please enter cover URL', trigger: ['input', 'blur'] }"
        >
          <NInput
            v-model:value="modalForm.cover"
            type="textarea"
            placeholder="Auto-generated after successful image upload, or directly paste external link"
          />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
