<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NDropdown, NForm, NFormItem, NImage, NInput } from 'naive-ui'

import CrudModal from '@/components/crud/CrudModal.vue'
import UploadOne from '@/components//UploadOne.vue'
import CommonPage from '@/components/common/CommonPage.vue'

import { convertImgUrl } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

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
const reloadFlag = ref(false)
const uploadOneRef = ref(null) // Image upload ref object

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  const resp = await api.getPages()
  pageList.value = resp.data
}

// Refresh preview image based on input link
function refreshImg(img) {
  reloadFlag.value = true
  uploadOneRef.value.previewImg = img
  setTimeout(() => reloadFlag.value = false, 600)
}

function handleSelect(key, page) {
  if (key === 'edit') {
    handleEdit(page)
  }
  else if (key === 'delete') {
    handleDelete([page.id])
  }
}

const options = [
  {
    label: 'Edit',
    key: 'edit',
    icon: () => h('i', { class: 'i-mingcute:edit-2-line' }),
  },
  {
    label: 'Delete',
    key: 'delete',
    icon: () => h('i', { class: 'i-mingcute:delete-back-line' }),
  },
]
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
    <div class="flex flex-wrap justify-between">
      <div
        v-for="page of pageList" :key="page.id"
        class="relative my-2 w-[300px] cursor-pointer text-center"
      >
        <div class="absolute right-2 top-1 text-white">
          <NDropdown :options="options" @select="handleSelect($event, page)">
            <span class="i-ion:ellipsis-horizontal h-5 w-5 text-white hover:text-blue" />
          </NDropdown>
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
      <div class="h-0 w-[300px]" />
      <div class="h-0 w-[300px]" />
      <div class="h-0 w-[300px]" />
    </div>

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
          label="Name"
          path="name"
          :rule="{ required: true, message: 'Please enter page name', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.name" placeholder="Page name" />
        </NFormItem>
        <NFormItem
          label="Label"
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
