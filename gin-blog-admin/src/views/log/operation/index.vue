<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NCode, NForm, NFormItem, NInput, NPopconfirm, NTag } from 'naive-ui'
import { useClipboard } from '@vueuse/core'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: 'Operation Log' })

// Request method corresponds to different types of tags (computed property with parameters)
function tagType(type) {
  switch (type) {
    case 'GET':
      return 'info'
    case 'POST':
      return 'success'
    case 'PUT':
      return 'warning'
    case 'DELETE':
      return 'error'
    default:
      return 'info'
  }
}

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

const {
  modalVisible,
  modalLoading,
  handleDelete,
  modalForm,
  modalFormRef,
  handleView,
} = useCRUD({
  name: 'Log',
  doDelete: api.deleteOperationLogs,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  { title: 'System Module', key: 'opt_module', width: 70, align: 'center', ellipsis: { tooltip: true } },
  { title: 'Operation Type', key: 'opt_type', width: 70, align: 'center', ellipsis: { tooltip: true } },
  // { title: 'Operation Description', key: 'opt_desc', width: 80, align: 'center', ellipsis: { tooltip: true } },
  {
    title: 'Request Method',
    key: 'request_method',
    width: 80,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h(
        NTag,
        { type: tagType(row.request_method) }, // Note: using computed property here
        { default: () => row.request_method },
      )
    },
  },
  { title: 'Operator', key: 'nickname', width: 80, align: 'center', ellipsis: { tooltip: true } },
  { title: 'IP Address', key: 'ip_address', width: 80, align: 'center', ellipsis: { tooltip: true } },
  { title: 'Location', key: 'ip_source', width: 80, align: 'center', ellipsis: { tooltip: true } },
  {
    title: 'Created Time',
    key: 'created_at',
    align: 'center',
    width: 80,
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.created_at),
          icon: () => h('i', { class: 'i-mdi:update' }),
        },
      )
    },
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 120,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            quaternary: true,
            type: 'info',
            onClick: () => handleView(row),
          },
          {
            default: () => 'View',
            icon: () => h('i', { class: 'i-ic:outline-remove-red-eye' }),
          },
        ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  quaternary: true,
                  type: 'error',
                  style: 'margin-left: 15px;',
                },
                {
                  default: () => 'Delete',
                  icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
                },
              ),
            default: () => h('div', {}, 'Are you sure you want to delete this log?'),
          },
        ),
      ]
    },
  },
]

function copyFormatCode(code) {
  const { copy } = useClipboard()
  copy(JSON.stringify(JSON.parse(code), null, 2))
  window.$message.success('Content copied to clipboard!')
}
</script>

<template>
  <CommonPage title="Operation Log">
    <template #action>
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
      :get-data="api.getOperationLogs"
    >
      <template #queryBar>
        <QueryItem label="Module" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Please enter module name or description"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>

    <CrudModal
      v-model:visible="modalVisible"
      title="Log Details"
      :show-footer="false"
      :loading="modalLoading"
      width="full"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="90"
        :model="modalForm"
      >
        <NFormItem label="Module: " path="opt_module">
          {{ modalForm.opt_module }}
        </NFormItem>
        <NFormItem label="Request URL: " path="opt_url">
          {{ modalForm.opt_url }}
        </NFormItem>
        <NFormItem label="Request Method: " path="request_method">
          <NTag :type="tagType(modalForm.request_method)">
            {{ modalForm.request_method }}
          </NTag>
        </NFormItem>
        <NFormItem label="Operation Type: " path="opt_type">
          {{ modalForm.opt_type }}
        </NFormItem>
        <NFormItem label="Operation Method: " path="opt_method">
          <NCode
            :code="modalForm.opt_method"
            code-wrap
            language="json"
          />
        </NFormItem>
        <NFormItem label="Operator: " path="nickname">
          {{ modalForm.nickname }}
        </NFormItem>
        <NFormItem label="Request Parameters: " path="request_param">
          <NCode
            class="word-wrap cursor-pointer p-7"
            :code="JSON.stringify(JSON.parse(modalForm.request_param), null, 2)"
            language="json"
            @click="copyFormatCode(modalForm.request_param)"
          />
        </NFormItem>
        <NFormItem label="Response Data: " path="response_data">
          <NCode
            class="cursor-pointer p-7"
            :code="JSON.stringify(JSON.parse(modalForm.response_data), null, 2)"
            language="json"
            @click="copyFormatCode(modalForm.response_data)"
          />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
