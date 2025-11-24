<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NPopconfirm, NSwitch, NTag, NTree } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: 'Role Management' })

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
  name: 'Role',
  initForm: {},
  doCreate: api.saveOrUpdateRole,
  doDelete: api.deleteRole,
  doUpdate: api.saveOrUpdateRole,
  refresh: () => $table.value?.handleSearch(),
})

// Menu and resource popup menu options are different
const showMenu = ref(true)
const resourceOption = ref([]) // Resource options
const menuOption = ref([]) // Menu options

onMounted(() => {
  $table.value?.handleSearch()
  // api.getResourceOption().then(res => (resourceOption.value = res.data))
  // api.getMenuOption().then(res => (menuOption.value = res.data))
})

const columns = [
  {
    type: 'selection',
    width: 15,
    fixed: 'left',
  },
  {
    title: 'Role Name',
    key: 'name',
    width: 80,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Role Label',
    key: 'label',
    width: 80,
    align: 'center',
    render(row) {
      return h(NTag, { type: 'info' }, { default: () => row.label })
    },
  },
  {
    title: 'Created Date',
    key: 'created_at',
    width: 60,
    align: 'center',
    render(row) {
      return h('span', formatDate(row.created_at))
    },
  },
  {
    title: 'Disabled',
    key: 'is_disable',
    width: 30,
    align: 'center',
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.is_disable,
        loading: !!row.publishing, // Loading animation
        checkedValue: 1,
        uncheckedValue: 0,
        onUpdateValue: () => $message.info('This feature is not supported yet~'),
      })
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
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: async () => {
              showMenu.value = true
              await api.getMenuOption().then(resp => (menuOption.value = resp.data))
              handleEdit(row)
            },
          },
          {
            default: () => 'Menu Permissions',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline' }),
          },
        ),
        h(
          NButton,
          {
            size: 'tiny',
            quaternary: true,
            type: 'info',
            onClick: async () => {
              showMenu.value = false
              await api.getResourceOption().then(resp => (resourceOption.value = resp.data))
              handleEdit(row)
            },
          },
          {
            default: () => 'Resource Permissions',
            icon: () => h('i', { class: 'i-ic:baseline-folder-open' }),
          },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete([row.id], false),
            onNegativeClick: () => {},
          },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  style: 'margin-left: 15px;',
                },
                {
                  default: () => 'Delete',
                  icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
                },
              ),
            default: () => h('div', {}, 'Are you sure you want to delete this role?'),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <CommonPage title="Role Management">
    <template #action>
      <NButton type="primary" @click="handleAdd">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        New Role
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        Batch Delete
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getRoles"
    >
      <template #queryBar>
        <QueryItem label="Role Name" :label-width="80">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Enter role name"
            @keydown.enter=" $table?.handleSearch()"
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
        :disabled="modalAction === 'view'"
      >
        <NFormItem label="Role Name" path="name">
          <NInput v-model:value="modalForm.name" placeholder="Enter role name" />
        </NFormItem>
        <NFormItem label="Role Label" path="name">
          <NInput v-model:value="modalForm.label" placeholder="Enter role label" />
        </NFormItem>
        <!-- TODO: Can select menu and resource permissions when adding -->
        <template v-if="modalAction === 'edit'">
          <NFormItem v-if="showMenu" label="Menu Permissions" path="menu_ids">
            <NTree
              :data="menuOption"
              :checked-keys="modalForm.menu_ids"

              checkable expand-on-click block-line
              @update:checked-keys="(v) => (modalForm.menu_ids = v)"
            />
          </NFormItem>
          <NFormItem v-else label="Resource Permissions" path="resource_ids">
            <NTree
              :data="resourceOption"
              :checked-keys="modalForm.resource_ids"

              block-line checkable expand-on-click cascade accordion
              @update:checked-keys="(v) => (modalForm.resource_ids = v)"
            />
          </NFormItem>
        </template>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
