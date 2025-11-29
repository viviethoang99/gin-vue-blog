<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NPopconfirm } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: 'Category Management' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

onMounted(() => {
  $table.value?.handleSearch()
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
  name: 'Category',
  initForm: {},
  doCreate: api.saveOrUpdateCategory,
  doDelete: api.deleteCategory,
  doUpdate: api.saveOrUpdateCategory,
  refresh: () => $table.value?.handleSearch(),
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  { title: 'Category Name', key: 'name', width: 100, align: 'center', ellipsis: { tooltip: true } },
  { title: 'Article Count', key: 'article_count', width: 30, align: 'center' },
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
    title: 'Updated Date',
    key: 'updated_at',
    width: 80,
    align: 'center',
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.updated_at),
          icon: () => h('i', { class: 'i-mdi:update' }),
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
          { size: 'small', type: 'primary', onClick: () => handleEdit(row) },
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
            default: () => h('div', {}, 'Are you sure you want to delete this category?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <CommonPage title="Category Management">
    <template #action>
      <NButton type="primary" secondary @click="$table?.handleExport()">
        <template #icon>
          <i class="i-mdi:download" />
        </template>
        Export
      </NButton>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        New Category
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
      :get-data="api.getCategorys"
    >
      <template #queryBar>
        <QueryItem label="Search" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Please enter category name"
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
          label="Article Category"
          path="name"
          :rule="{ required: true, message: 'Please enter category name', trigger: ['input', 'blur'] }"
        >
          <NInput
            v-model:value="modalForm.name"
            placeholder="Please enter category name"
            clearable
          />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
