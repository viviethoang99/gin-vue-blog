<script setup>
import { useClipboard } from '@vueuse/core'
import { NButton, NForm, NFormItem, NImage, NInput, NPopconfirm } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { convertImgUrl, formatDate, IMG_PLACEHOLDER } from '@/utils'

defineOptions({ name: 'Friend Links' })

// 在 setup 里创建一次: 原来写在列的 render 里, 每次点击都新建一个实例
const { copy } = useClipboard()

const $table = ref(null)
const queryItems = ref({
  keyword: '', // Friend link name | address | description
})

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
  name: 'Friend Link',
  initForm: {},
  doCreate: api.saveOrUpdateLink,
  doDelete: api.deleteLinks,
  doUpdate: api.saveOrUpdateLink,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: 'Avatar',
    key: 'avatar',
    width: 40,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 40,
        'imgProps': { style: { 'border-radius': '3px' } },
        // 本地上传的头像存的是相对路径, 不转换直接给 img 会裂图
        'src': convertImgUrl(row.avatar),
        'fallback-src': IMG_PLACEHOLDER, // 加载失败时用内联占位图, 不再请求外网
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: 'Link Name',
    key: 'name',
    width: 100,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Address',
    key: 'address',
    width: 120,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h(
        'a',
        {
          class: 'hover:underline hover:underline-blue-500 hover:underline-2 hover:underline-solid hover:underline-offset-4 cursor-pointer',
          // href: row.address,
          // target: '_blank',
          onClick: () => {
            copy(row.address)
            $message.info('Link copied to clipboard!')
          },
        },
        row.address,
      )
    },
  },
  {
    title: 'Description',
    key: 'intro',
    width: 120,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Created Date',
    key: 'created_at',
    width: 80,
    align: 'center',
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.created_at),
          icon: () => h('i', { class: 'i-mdi:clock-time-three-outline' }),
        },
      )
    },
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 100,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            onClick: () => handleEdit(row),
          },
          { default: () => 'Edit', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () => h(
              NButton,
              { size: 'small', type: 'error', style: 'margin-left: 15px;' },
              { default: () => 'Delete', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
            ),
            default: () => h('div', {}, '确定删除该友链吗?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <CommonPage title="Friend Links">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        New Friend Link
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <span class="i-material-symbols:playlist-remove" />
        </template>
        Batch Delete
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getLinks"
    >
      <template #queryBar>
        <QueryItem label="Name | Address | Description" :label-width="150">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Search keywords"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>

    <CrudModal
      v-model:visible="modalVisible"
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
          label="Link Name"
          path="name"
          :rule="{ required: true, message: 'Please enter link name', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.name" placeholder="Please enter link name" />
        </NFormItem>
        <NFormItem
          label="Avatar"
          path="avatar"
          :rule="{ required: true, message: 'Please enter avatar URL', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.avatar" placeholder="Please enter avatar URL" />
        </NFormItem>
        <NFormItem
          label="Address"
          path="address"
          :rule="{ required: true, message: 'Please enter link address', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.address" placeholder="Please enter link address" />
        </NFormItem>
        <NFormItem
          label="Description"
          path="intro"
          :rule="{ required: true, message: 'Please enter description', trigger: ['input', 'blur'] }"
        >
          <NInput v-model:value="modalForm.intro" placeholder="Please enter description" />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
