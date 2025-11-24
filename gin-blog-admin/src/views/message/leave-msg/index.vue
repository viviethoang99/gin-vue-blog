<script setup>
import { h, onMounted, ref } from 'vue'
import { NButton, NImage, NInput, NPopconfirm, NTabPane, NTabs, NTag } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import api from '@/api'

defineOptions({ name: 'Message Management' })

onMounted(() => {
  handleChangeTab('all') // Default to view all
})

const $table = ref(null)
const queryItems = ref({
  nickname: '',
})
const extraParams = ref({
  is_review: null, // Message status: Under review | Approved
})

const { handleDelete } = useCRUD({
  name: 'Message',
  doDelete: api.deleteMessages,
  refresh: () => $table.value?.handleSearch(),
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
        'src': convertImgUrl(row.avatar),
        'fallback-src': 'http://dummyimage.com/400x400', // Load failed
        'show-toolbar-tooltip': true,
      })
    },
  },
  {
    title: 'Messenger',
    key: 'nickname',
    width: 60,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Message Content',
    key: 'content',
    width: 120,
    align: 'center',
  },
  {
    title: 'IP Address',
    key: 'ip_address',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'IP Source',
    key: 'ip_source',
    width: 70,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', row.ip_source || 'Unknown')
    },
  },
  {
    title: 'Message Time',
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
    title: 'Status',
    key: 'is_review',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.is_review ? 'success' : 'error' },
        { default: () => (row.is_review ? 'Approved' : 'Under Review') },
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
        row.is_review
          ? h(
            NButton,
            {
              size: 'small',
              type: 'warning',
              onClick: () => handleUpdateReview([row.id], false),
            },
            {
              default: () => 'Revoke',
              icon: () => h('i', { class: 'i-mi:circle-error' }),
            },
          )
          : h(
            NButton,
            {
              size: 'small',
              type: 'success',
              style: 'margin-left: 15px;',
              onClick: () => handleUpdateReview([row.id], true),
            },
            {
              default: () => 'Approve',
              icon: () => h('i', { class: 'i-mi:circle-check' }),
            },
          ),
        h(
          NPopconfirm,
          { onPositiveClick: () => handleDelete([row.id], false) },
          {
            trigger: () =>
              h(
                NButton,
                { size: 'small', type: 'error', style: 'margin-left: 15px;' },
                { default: () => 'Delete', icon: () => h('i', { class: 'i-material-symbols:delete-outline' }) },
              ),
            default: () => h('div', {}, 'Are you sure you want to delete this message?'),
          },
        ),
      ]
    },
  },
]

// Update message review
async function handleUpdateReview(ids, is_review) {
  if (!ids.length) {
    $message.info('Please select data to review')
    return
  }

  await api.updateMessageReview(ids, is_review)
  $message?.success(is_review ? 'Review successful' : 'Revoke successful')
  $table.value?.handleSearch()
}

// Switch tab: [All, Approved, Under Review]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.value.is_review = null
      break
    case 'has_review': // Approved
      extraParams.value.is_review = 1
      break
    case 'not_review': // Under Review
      extraParams.value.is_review = 0
      break
  }
  $table.value?.handleSearch()
}
</script>

<template>
  <CommonPage title="Message Management">
    <template #action>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <span class="i-material-symbols:recycling-rounded" />
        </template>
        Batch Delete
      </NButton>
      <NButton
        type="success"
        :disabled="!$table?.selections.length"
        @click="handleUpdateReview($table.selections, true)"
      >
        <template #icon>
          <span class="i-ic:outline-approval" />
        </template>
        Batch Approve
      </NButton>
    </template>
    <NTabs
      type="line"
      animated
      @update:value="handleChangeTab"
    >
      <template #prefix>
        Status
      </template>
      <NTabPane name="all" tab="All" />
      <NTabPane name="has_review" tab="Approved" />
      <NTabPane name="not_review" tab="Under Review" />
    </NTabs>
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getMessages"
    >
      <template #queryBar>
        <QueryItem label="User" :label-width="40" :content-width="180">
          <NInput
            v-model:value="queryItems.nickname"
            clearable
            type="text"
            placeholder="Enter user nickname"
            @keydown.enter=" $table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
