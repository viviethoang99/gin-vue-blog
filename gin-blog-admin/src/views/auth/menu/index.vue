<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NInputNumber, NPopconfirm, NRadio, NRadioGroup, NSpace, NSwitch, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import CrudTable from '@/components/crud/CrudTable.vue'
import IconPicker from '@/components/icon/IconPicker.vue'
import TheIcon from '@/components/icon/TheIcon.vue'

import { formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: 'Menu Management' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

const initForm = {
  order_num: 1,
  is_hidden: false, // Whether hidden
  is_catalogue: false, // Whether directory
  is_external: false, // Whether external link
  keep_alive: false,
  icon: 'mdi-account',
  order_num: 1,
  name: '',
  path: '',
  redirect: '',
  component: '',
  parent_id: 0,
}

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
  name: 'Menu',
  initForm,
  doCreate: api.saveOrUpdateMenu,
  doDelete: api.deleteMenu,
  doUpdate: api.saveOrUpdateMenu,
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  {
    title: 'Menu Name',
    key: 'name',
    width: 100,
    render: (row) => {
      const groups = []
      groups.push(h('span', row.name))

      if (row.parent_id === 0) {
        groups.push(
          h(
            NTag,
            { type: row.is_catalogue ? 'info' : 'success', class: 'ml-1.5' },
            { default: () => row.is_catalogue ? 'Directory' : 'Top Menu' },
          ),
        )
      }
      else {
        groups.push(
          h(
            NTag,
            { type: 'default', class: 'ml-1.5' },
            { default: () => 'Submenu' },
          ),
        )
      }

      if (row.is_external) {
        groups.push(
          h(
            NTag,
            { type: 'warning', class: 'ml-1.5' },
            { default: () => 'External' },
          ),
        )
      }

      return groups
    },
  },
  {
    title: 'Icon',
    key: 'icon',
    width: 30,
    render(row) {
      return h(TheIcon, { icon: row.icon, size: 20 })
    },
  },
  { title: 'Order', key: 'order_num', width: 30, ellipsis: { tooltip: true } },
  { title: 'Path', key: 'path', width: 60, ellipsis: { tooltip: true } },
  {
    title: 'Redirect',
    key: 'redirect',
    width: 80,
    render(row) {
      if (row.parent_id === 0 && !row.is_catalogue) {
        return h('span', row.redirect)
      }
      return h('span', '-')
    },
  },
  {
    title: 'Component',
    key: 'component',
    width: 80,
    render(row) {
      if (!row.is_catalogue) {
        return h('span', row.component)
      }
      return h('span', '-')
    },
  },
  {
    title: 'Keep Alive',
    key: 'keep_alive',
    width: 30,
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.keep_alive,
        loading: !!row.publishing,
        onUpdateValue: () => handleUpdateKeepAlive(row),
      })
    },
  },
  {
    title: 'Hidden',
    key: 'is_hidden',
    width: 30,
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.is_hidden,
        loading: !!row.publishing,
        onUpdateValue: () => handleUpdateHidden(row),
      })
    },
  },
  {
    title: 'Updated Date',
    key: 'updated_at',
    width: 70,
    render(row) {
      return h('span', formatDate(row.updated_at))
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
            style: `display: ${!row.is_catalogue && row.parent_id === 0 ? '' : 'none'};`,
            onClick: () => {
              initForm.component = '' // Manually clear component path
              initForm.parent_id = row.id // Set parent menu id
              initForm.is_catalogue = false
              handleAdd()
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
            onClick: () => {
              handleEdit(row)
            },
          },
          { default: () => 'Edit', icon: () => h('i', { class: 'i-material-symbols:edit-outline' }) },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row.id, false),
          },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'tiny', quaternary: true, type: 'error' },
                {
                  default: () => 'Delete',
                  icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
                },
              ),
            default: () => h('div', {}, 'Are you sure you want to delete this menu?'),
          },
        ),
      ]
    },
  },
]

async function handleUpdateKeepAlive(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.keep_alive = !row.keep_alive
  try {
    await api.saveOrUpdateMenu(row)
    $message?.success(row.keep_alive ? 'Keep alive enabled' : 'Keep alive disabled')
  }
  catch (err) {
    row.keep_alive = !row.keep_alive
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

async function handleUpdateHidden(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_hidden = !row.is_hidden
  try {
    await api.saveOrUpdateMenu(row)
    $message?.success(row.is_hidden ? 'Hidden' : 'Visible')
  }
  catch (err) {
    row.is_hidden = !row.is_hidden
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// Add menu (optional directory)
function handleClickAdd() {
  initForm.is_catalogue = true // Default select "Directory"
  initForm.component = 'Layout' // Directory must be "Layout", top menu can be "Layout"
  initForm.parent_id = 0 // Directory and top menu parent id is 0
  handleAdd()
}
</script>

<template>
  <CommonPage title="Menu Management">
    <template #action>
      <NButton type="primary" @click="handleClickAdd">
        <template #icon>
          <span class="i-material-symbols:add" />
        </template>
        New Menu
      </NButton>
    </template>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.getMenus"
      :single-line="true"
    >
      <template #queryBar>
        <QueryItem label="Menu Name" :label-width="80">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Enter menu name"
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
        <NFormItem v-if="modalForm.parent_id === 0" label="Menu Type" path="type">
          <NRadioGroup v-model:value="modalForm.is_catalogue" name="radiogroup">
            <NSpace>
              <NRadio :value="true">
                Directory
              </NRadio>
              <NRadio :value="false">
                Top Menu
              </NRadio>
            </NSpace>
          </NRadioGroup>
        </NFormItem>
        <NFormItem label="Menu Name" path="name">
          <NInput v-model:value="modalForm.name" placeholder="Enter menu name" />
        </NFormItem>
        <NFormItem label="Menu Icon" path="icon">
          <IconPicker v-model:value="modalForm.icon" />
        </NFormItem>
        <NFormItem v-if="!modalForm.is_catalogue" label="Component Path" path="component">
          <NInput v-model:value="modalForm.component" placeholder="Enter component path" />
        </NFormItem>
        <NFormItem label="Access Path" path="path">
          <NInput v-model:value="modalForm.path" placeholder="Enter access path" />
        </NFormItem>
        <NFormItem v-if="!modalForm.is_catalogue" label="Redirect Path" path="redirect">
          <NInput
            v-model:value="modalForm.redirect"
            :disabled="modalForm.parent_id !== 0"
            placeholder="Only top-level menus can set redirect path"
          />
        </NFormItem>
        <NFormItem label="Display Order" path="order_num">
          <NInputNumber v-model:value="modalForm.order_num" />
        </NFormItem>
        <NFormItem label="Hidden" path="is_hidden">
          <NSwitch v-model:value="modalForm.is_hidden" />
        </NFormItem>
        <NFormItem label="External Link" path="is_external">
          <NSwitch v-model:value="modalForm.is_external" />
        </NFormItem>
        <NFormItem label="KeepAlive" path="keep_alive">
          <NSwitch v-model:value="modalForm.keep_alive" />
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
