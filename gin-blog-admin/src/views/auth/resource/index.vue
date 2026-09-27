<script setup>
import { NButton, NForm, NFormItem, NGradientText, NInput, NPopconfirm, NRadio, NRadioGroup, NSpace, NSwitch, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { formatDate } from '@/utils'

defineOptions({ name: 'Resource Management' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

const {
  modalVisible,
  modalAction,
  modalTitle,
  modalLoading,
  handleAdd,
  handleDelete,
  handleEdit,
  handleSave,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: 'Resource',
  doCreate: api.saveOrUpdateResource,
  doDelete: api.deleteResource,
  doUpdate: api.saveOrUpdateResource,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

// Request methods
const requestMethods = ['GET', 'POST', 'DELETE', 'PUT']

// Request methods correspond to different tag types (computed property with params)
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

const columns = [
  {
    title: 'Resource Name',
    key: 'name',
    width: 80,
    ellipsis: { tooltip: true },
  },
  {
    title: 'Resource Path',
    key: 'url',
    width: 80,
    ellipsis: { tooltip: true },
    render(row) {
      return row.children ? '-' : h('span', { class: 'color-[#1890ff]' }, row.url)
    },
  },
  {
    title: 'Request Method',
    key: 'request_method',
    width: 50,
    align: 'center',
    render(row) {
      return row.children
        ? '-'
        : h(
            NTag,
            { type: tagType(row.request_method) }, // 注意这里使用计算属性
            { default: () => row.request_method },
          )
    },
  },
  {
    title: 'Anonymous Access',
    // key 只用于导出(CrudTable 按 item[key] 取值), 以前写的是不存在的 is_hidden
    key: 'is_anonymous',
    width: 50,
    align: 'center',
    fixed: 'left',
    render(row) {
      return row.children
        ? '-'
        : h(NSwitch, {
            size: 'small',
            rubberBand: false,
            value: row.is_anonymous,
            loading: !!row.publishing, // 修改 ing 动画
            onUpdateValue: () => handleUpdateAnonymous(row),
          })
    },
  },
  {
    title: 'Created Date',
    key: 'created_at',
    width: 60,
    render(row) {
      return h('span', formatDate(row.created_at))
    },
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 115,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'primary',
            style: `display: ${row.children ? '' : 'none'};`,
            onClick: () => {
              handleAdd() // Add modal
              modalForm.value.parent_id = row.id // Parent resource id
            },
          },
          { default: () => 'Add', icon: () => h('i', { class: 'i-material-symbols:add' }) },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: () => (row.children ? handleEditModule(row) : handleEdit(row)),
          },
          { default: () => 'Edit', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => {
              handleDelete(row.id, false)
            },
          },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'tiny', quaternary: true, type: 'error' },
                { default: () => 'Delete', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
              ),
            default: () => h('div', {}, 'Are you sure you want to delete this resource?'),
          },
        ),
      ]
    },
  },
]

// Update anonymous access permission
async function handleUpdateAnonymous(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_anonymous = !row.is_anonymous
  try {
    await api.updateResourceAnonymous(row)
    $message?.success(row.is_anonymous ? 'Anonymous access allowed' : 'Anonymous access denied')
  }
  catch (err) {
    row.is_anonymous = !row.is_anonymous
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// Module related
const moduleModalVisible = ref(false)
function handleAddModule() {
  modalAction.value = 'add'
  modalForm.value = {}
  moduleModalVisible.value = true
}
function handleEditModule(row) {
  modalAction.value = 'edit'
  modalForm.value = { ...row }
  moduleModalVisible.value = true
}
// 保存模块
//
// 以前是 handleSave() 不 await 紧接着关掉弹窗: 请求还没回来弹窗就关了,
// 保存失败也一样关, 用户以为成功了
async function handleModuleSave() {
  try {
    modalLoading.value = true
    await modalFormRef.value?.validate()
    await api.saveOrUpdateResource(modalForm.value)
    $message.success(modalAction.value === 'add' ? '新增成功' : '编辑成功')
    moduleModalVisible.value = false
    $table.value?.handleSearch()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    modalLoading.value = false
  }
}
</script>

<template>
  <CommonPage title="Resource Management">
    <template #action>
      <NButton type="primary" @click="handleAddModule">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        New Module
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.getResources"
      :single-line="true"
    >
      <template #queryBar>
        <QueryItem label="Resource Name" :label-width="50">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Enter resource name"
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
        <NFormItem label="Resource Name" path="name">
          <NInput v-model:value="modalForm.name" placeholder="Enter resource name" />
        </NFormItem>
        <NFormItem label="Resource Path" path="url">
          <NInput v-model:value="modalForm.url" placeholder="Enter resource path" />
        </NFormItem>
        <NFormItem label="Request Method" path="request_method">
          <NRadioGroup v-model:value="modalForm.request_method" name="radiogroup">
            <NSpace>
              <NRadio v-for="method of requestMethods" :key="method" :value="method">
                <NGradientText :type="tagType(method)">
                  {{ method }}
                </NGradientText>
              </NRadio>
            </NSpace>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </CrudModal>

    <CrudModal
      v-model:visible="moduleModalVisible"
      :title="`${modalAction === 'add' ? 'Add' : 'Edit'} Module`"
      :loading="modalLoading"
      @save="handleModuleSave"
    >
      <NForm
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
      >
        <NFormItem label="Module Name" path="name">
          <NInput v-model:value="modalForm.name" placeholder="Enter module name" />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
