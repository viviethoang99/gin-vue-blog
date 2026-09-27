<script setup>
import { NButton, NCheckbox, NCheckboxGroup, NForm, NFormItem, NImage, NInput, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import { loginTypeMap, loginTypeOptions } from '@/assets/config'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'

import CrudTable from '@/components/crud/CrudTable.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { convertImgUrl, formatDate, IMG_PLACEHOLDER } from '@/utils'

defineOptions({ name: 'User List' })

const $table = ref(null)
const queryItems = ref({
  username: '',
  nickname: '',
  login_type: null,
})

const {
  modalVisible,
  modalLoading,
  handleSave,
  handleEdit,
  modalForm,
  modalFormRef,
} = useCRUD({
  name: 'User',
  doUpdate: api.updateUser,
  refresh: () => $table.value?.handleSearch(),
})

const roleOptions = ref([])

onMounted(() => {
  // 拦截器已经弹过错误提示, 这里补 catch 只是别留下 unhandled rejection
  api.getRoleOption().then(resp => roleOptions.value = resp.data).catch(err => console.error(err))
  $table.value?.handleSearch()
})

const columns = [
  {
    title: 'Avatar',
    key: 'avatar',
    width: 30,
    align: 'center',
    render(row) {
      return h(NImage, {
        'height': 30,
        'imgProps': { style: { 'border-radius': '3px' } },
        'src': convertImgUrl(row.info?.avatar),
        'fallback-src': IMG_PLACEHOLDER, // 加载失败时用内联占位图, 不再请求外网
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: 'Nickname',
    key: 'nickname',
    width: 60,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.info?.nickname)
    },
  },
  {
    title: 'Login Type',
    key: 'login_type',
    width: 40,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: loginTypeMap[row.login_type]?.tag },
        { default: () => loginTypeMap[row.login_type]?.name || 'Unknown' },
      )
    },
  },
  {
    title: 'User Role',
    key: 'role',
    width: 80,
    align: 'center',
    render(row) {
      if (row.is_super) {
        return h(NTag, { type: 'error' }, { default: () => 'Super Admin' })
      }
      const roles = row.roles ?? []
      const groups = []
      for (let i = 0; i < roles.length; i++) {
        groups.push(h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => roles[i].name }))
      }
      return h('span', groups.length ? groups : 'None')
    },
  },
  {
    title: 'Login IP',
    key: 'ip_address',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.ip_address || 'Unknown')
    },
  },
  {
    title: 'Login Location',
    key: 'ip_source',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.ip_source || 'Unknown')
    },
  },
  {
    title: 'Created Time',
    key: 'created_at',
    align: 'center',
    width: 70,
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
    title: 'Last Login Time',
    key: 'last_login_time',
    align: 'center',
    width: 70,
    render(row) {
      return h(
        NButton,
        { size: 'small', type: 'text', ghost: true },
        {
          default: () => formatDate(row.last_login_time),
          icon: () => h('i', { class: 'i-mdi:update' }),
        },
      )
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
        loading: !!row.publishing,
        onUpdateValue: () => handleUpdateDisable(row),
      })
    },
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 60,
    align: 'center',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            onClick: () => {
              row.nickname = row.info?.nickname
              // roles => role_ids, 没有任何角色的用户 roles 可能是 null
              row.role_ids = row.roles?.map(e => e.id) ?? []
              handleEdit(row)
            },
          },
          {
            default: () => 'Edit',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline' }),
          },
        ),
      ]
    },
  },
]

// Update user disabled status
async function handleUpdateDisable(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_disable = !row.is_disable
  try {
    await api.updateUserDisable(row.id, row.is_disable)
    $message?.success(row.is_disable ? 'User disabled' : 'User enabled')
    $table.value?.handleSearch()
  }
  catch (err) {
    row.is_disable = !row.is_disable
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}
</script>

<template>
  <CommonPage title="User List">
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getUsers"
    >
      <template #queryBar>
        <QueryItem label="Nickname" :label-width="40" :content-width="160">
          <NInput
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="Enter nickname"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="Username" :label-width="60" :content-width="160">
          <NInput
            v-model:value="queryItems.username"
            clearable
            type="text"
            placeholder="Enter username"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="Login Type" :label-width="70" :content-width="160">
          <NSelect
            v-model:value="queryItems.login_type"
            clearable
            filterable
            placeholder="Select login type"
            :options="loginTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>

    <CrudModal
      v-model:visible="modalVisible"
      title="Edit User"
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
        <NFormItem label="User Nickname" path="name">
          <NInput
            v-model:value="modalForm.nickname"
            clearable
            placeholder="Enter user nickname"
          />
        </NFormItem>
        <NFormItem label="Role" path="role_ids">
          <NCheckboxGroup v-model:value="modalForm.role_ids">
            <NSpace item-style="display: flex;">
              <NCheckbox
                v-for="item in roleOptions"
                :key="item.value"
                :value="item.value"
                :label="item.label"
              />
            </NSpace>
          </NCheckboxGroup>
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>
