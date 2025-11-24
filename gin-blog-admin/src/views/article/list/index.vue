<script setup>
import { defineOptions, h, onActivated, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NImage, NInput, NPopconfirm, NSelect, NSwitch, NTabPane, NTabs, NTag, NUpload } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import { convertImgUrl, formatDate } from '@/utils'
import { useCRUD } from '@/composables'
import { articleTypeMap, articleTypeOptions } from '@/assets/config'
import api from '@/api'

// KeepAlive requires name attribute that corresponds to name in router
defineOptions({ name: 'Article List' })

const route = useRoute()
const router = useRouter()

const categoryOptions = ref([])
const tagOptions = ref([])

const $table = ref(null)

const queryItems = ref({
  title: '', // Title
  type: null, // Type
  category_id: null, // Category
  tag_id: null, // Tag
})

const extraParams = ref({
  is_delete: null, // Not deleted | Recycle bin
  status: null, // null-all, 1-public, 2-private, 3-draft
})

const { handleDelete } = useCRUD({
  name: 'Article',
  doDelete: updateOrDeleteArticles, // Soft delete
  refresh: () => $table.value?.handleSearch(),
})

onMounted(() => {
  api.getCategoryOption().then(res => (categoryOptions.value = res.data))
  api.getTagOption().then(res => (tagOptions.value = res.data))
  handleChangeTab('all') // Default view all
})

// ! When switching pages, if coming from [Write Article] page, will carry needRefresh parameter
onActivated(() => {
  const { needRefresh } = route.query
  needRefresh && ($table.value?.handleSearch())
})

const columns = [
  { type: 'selection', width: 20, fixed: 'left' },
  {
    title: 'Cover',
    key: 'img',
    width: 55,
    align: 'center',
    render(row) {
      return h(NImage, {
        imgProps: { style: { 'border-radius': '2px', 'height': '100%', 'width': '100%' } },
        src: convertImgUrl(row.img),
        fallbackSrc: 'http://dummyimage.com/400x400',
        showToolbarTooltip: true,
      })
    },
  },
  {
    title: 'Title',
    key: 'title',
    width: 120,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Category',
    key: 'category.name',
    width: 60,
    align: 'center',
    ellipsis: { tooltip: true },
    render(row) {
      return h('div', row.category.name || 'None')
    },
  },
  {
    title: 'Tags',
    key: 'tags',
    width: 100,
    align: 'center',
    render(row) {
      const tags = row.tags ?? []
      const group = []
      for (let i = 0; i < tags.length; i++) {
        group.push(
          h(NTag, { type: 'info', style: { margin: '2px 3px' } }, { default: () => tags[i].name }),
        )
      }
      return h('div', group.length ? group : 'None')
    },
  },
  {
    title: 'Views',
    key: 'view_count',
    width: 40,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Likes',
    key: 'like_count',
    width: 40,
    align: 'center',
    ellipsis: { tooltip: true },
  },
  {
    title: 'Type',
    key: 'type',
    width: 50,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: articleTypeMap[row.type]?.tag },
        { default: () => articleTypeMap[row.type]?.name },
      )
    },
  },
  {
    title: 'Published Time',
    key: 'updateDate',
    align: 'center',
    width: 80,
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
    title: 'Pin to Top',
    key: 'is_top',
    width: 50,
    align: 'center',
    fixed: 'left',
    render(row) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.is_top,
        loading: !!row.publishing,
        onUpdateValue: () => handleUpdateTop(row),
      })
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
        row.is_delete
          ? h(
            NButton,
            {
              size: 'small',
              type: 'success',
              secondary: true,
              onClick: async () => {
                await api.softDeleteArticle([row.id], false)
                await $table.value?.handleSearch()
              },
            },
            { default: () => 'Restore', icon: () => h('i', { class: 'i-majesticons:eye-line' }) },
          )
          : h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              secondary: true,
              onClick: () => router.push(`/article/write/${row.id}`), // Navigate to write article page with parameters
            },
            { default: () => 'View', icon: () => h('i', { class: 'i-majesticons:eye-line' }) },
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
            default: () => h('div', {}, 'Are you sure you want to delete this article?'),
          },
        ),
      ]
    },
  },
]

function updateOrDeleteArticles(ids) {
  extraParams.value.is_delete
    ? api.deleteArticle(ids)
    : api.softDeleteArticle(JSON.parse(ids), true)
}

// Update article pin status
async function handleUpdateTop(row) {
  if (!row.id) {
    return
  }
  row.publishing = true
  row.is_top = !row.is_top
  try {
    await api.updateArticleTop(row.id, row.is_top)
    $message?.success(row.is_top ? 'Article pinned successfully' : 'Article unpinned successfully')
    $table.value?.handleSearch()
  }
  catch (err) {
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// Export articles
async function exportArticles(ids) {
  // Method 1: Export based on article content and title from frontend
  const list = $table.value?.tableData.filter(e => ids.includes(e.id))
  for (const item of list)
    downloadFile(item.content, `${item.title}.md`)

  // Method 2: Backend export returns links, frontend downloads based on links
  // const res = await api.exportArticles(ids)
  // for (const url of res.data)
  // downloadFile(url)
}

// Switch tabs: [All, Public, Private, Draft, Trash]
function handleChangeTab(value) {
  switch (value) {
    case 'all':
      extraParams.value.is_delete = 0
      extraParams.value.status = null
      break
    case 'public':
      extraParams.value.is_delete = 0
      extraParams.value.status = 1
      break
    case 'secret':
      extraParams.value.is_delete = 0
      extraParams.value.status = 2
      break
    case 'draft':
      extraParams.value.is_delete = 0
      extraParams.value.status = 3
      break
    case 'delete':
      extraParams.value.is_delete = 1
      extraParams.value.status = null
      break
  }
  $table.value?.handleSearch()
}

// Check file type before upload
function beforeUpload(data) {
  if (!data.file.name.endsWith('.md')) {
    $message.error('Only .md format files can be uploaded, please re-upload')
    return false
  }
  return true
}

// Operations after file upload
function afterUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  if (res.code === 0) {
    $table.value?.handleSearch()
    $message.success('Article imported successfully!')
  }
  else {
    $message.error('Article import failed!')
  }
}

function downloadFile(content, fileName) {
  const aEle = document.createElement('a') // Create download link
  aEle.download = fileName // Set download filename
  aEle.style.display = 'none'// Hidden downloadable link
  // Convert string content to blob address
  const blob = new Blob([content])
  aEle.href = URL.createObjectURL(blob)
  // Bind click event
  document.body.appendChild(aEle)
  aEle.click()
  // Then remove
  document.body.removeChild(aEle)
}
</script>

<template>
  <CommonPage title="Article List">
    <template #action>
      <NButton type="primary" @click="$router.replace('/article/write')">
        <template #icon>
          <i class="i-material-symbols:add" />
        </template>
        New Article
      </NButton>
      <NButton
        type="error"
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <i class="i-material-symbols:recycling-rounded" />
        </template>
        Batch Delete
      </NButton>
      <NButton
        type="info"
        :disabled="!$table?.selections.length"
        @click="exportArticles($table?.selections)"
      >
        <template #icon>
          <i class="i-mdi:export" />
        </template>
        Batch Export
      </NButton>
      <div class="inline-block">
        <NUpload
          action="/api/article/import"
          :show-file-list="false"
          multiple
          @before-upload="beforeUpload"
          @finish="afterUpload"
        >
          <NButton type="success">
            <template #icon>
              <i class="i-mdi:import" />
            </template>
            Batch Import
          </NButton>
        </NUpload>
      </div>
    </template>

    <NTabs type="line" animated @update:value="handleChangeTab">
      <template #prefix>
        Status
      </template>
      <NTabPane name="all" tab="All" />
      <NTabPane name="public" tab="Public" />
      <NTabPane name="secret" tab="Private" />
      <NTabPane name="draft" tab="Draft" />
      <NTabPane name="delete" tab="Trash" />
    </NTabs>

    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :extra-params="extraParams"
      :columns="columns"
      :get-data="api.getArticles"
    >
      <template #queryBar>
        <QueryItem label="Title" :label-width="40" :content-width="180">
          <NInput
            v-model:value="queryItems.title"
            clearable
            type="text"
            placeholder="Enter title"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="Type" :label-width="40" :content-width="160">
          <NSelect
            v-model:value="queryItems.type"
            clearable
            placeholder="Select article type"
            :options="articleTypeOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="Category" :label-width="60" :content-width="160">
          <NSelect
            v-model:value="queryItems.category_id"
            clearable
            filterable
            placeholder="Select article category"
            :options="categoryOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
        <QueryItem label="Tags" :label-width="40" :content-width="160">
          <NSelect
            v-model:value="queryItems.tag_id"
            clearable
            filterable
            placeholder="Select article tags"
            :options="tagOptions"
            @update:value="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
