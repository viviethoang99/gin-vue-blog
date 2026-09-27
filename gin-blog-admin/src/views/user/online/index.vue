<script setup>
import { NButton, NImage, NInput, NPopconfirm } from 'naive-ui'
import { h, onMounted, ref } from 'vue'

import api from '@/api'
import CommonPage from '@/components/common/CommonPage.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import QueryItem from '@/components/crud/QueryItem.vue'
import { convertImgUrl, formatDate, IMG_PLACEHOLDER } from '@/utils'

defineOptions({ name: 'Online Users' })

const $table = ref(null)
const queryItems = ref({
  keyword: '', // Username | Nickname
})

onMounted(() => {
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
      // info 来自 Redis 里反序列化的 UserAuth, 没有 Preload("UserInfo"), 可能为空;
      // 少一个可选链就会让整张表 render 抛错(user/list 那边写的是 row.info?.nickname)
      return h('span', row.info?.nickname || 'Unknown')
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
    title: 'Browser',
    key: 'browser',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.browser || 'Unknown')
    },
  },
  {
    title: 'Operating System',
    key: 'os',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.os || 'Unknown')
    },
  },
  {
    title: 'Login Time',
    key: 'last_login_time',
    align: 'center',
    width: 70,
    render(row) {
      return h('span', formatDate(row.last_login_time, 'YYYY-MM-DD HH:mm:ss'))
    },
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 60,
    align: 'center',
    fixed: 'right',
    render(row) {
      return h(
        NPopconfirm,
        { onPositiveClick: () => handleForceOffline(row) },
        {
          trigger: () =>
            h(
              NButton,
              { size: 'small', type: 'warning' },
              {
                default: () => 'Log Out',
                icon: () => h('i', { class: 'i-material-symbols:delete-outline' }),
              },
            ),
          default: () => h('div', {}, 'Are you sure you want to force this user offline?'),
        },
      )
    },
  },
]

// Force user offline
async function handleForceOffline(row) {
  try {
    await api.forceOfflineUser(row.id)
    window.$message.success('User has been forced offline!')
    $table.value?.handleSearch()
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <CommonPage title="Online Users">
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getOnlineUsers"
      :is-pagination="false"
    >
      <template #queryBar>
        <QueryItem label="Username | Nickname" :label-width="100" :content-width="200">
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
  </CommonPage>
</template>
