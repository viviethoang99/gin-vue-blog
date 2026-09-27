<script setup>
import { NButton, NImage, NInput, NPopconfirm, NSelect, NSwitch, NTabPane, NTabs, NTag, NUpload } from 'naive-ui'
import { defineOptions, h, onActivated, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import api from '@/api'
import { articleTypeMap, articleTypeOptions } from '@/assets/config'
import CommonPage from '@/components/common/CommonPage.vue'

import CrudTable from '@/components/crud/CrudTable.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import { useCRUD } from '@/composables'
import { useAuthStore } from '@/store'
import { convertImgUrl, downloadFile, formatDate, IMG_PLACEHOLDER, parseJson } from '@/utils'

// KeepAlive requires name attribute that corresponds to name in router
defineOptions({ name: 'Article List' })

const route = useRoute()
const router = useRouter()
// 不解构, 否则重新登录后拿到的还是旧 token
const authStore = useAuthStore()

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
  // 拦截器已经弹过错误提示, 这里补 catch 只是别留下 unhandled rejection
  api.getCategoryOption().then(res => (categoryOptions.value = res.data)).catch(err => console.error(err))
  api.getTagOption().then(res => (tagOptions.value = res.data)).catch(err => console.error(err))
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
      // 尺寸必须写死: 原来是 height/width 都给 100%, 而父容器高度又来自图片本身,
      // 循环依赖, 实测这一列渲染出来是 0x0, 封面永远看不见。
      // 同项目其他表格(留言/友链/用户)都是写死 40 / 30 的
      return h(NImage, {
        width: 40,
        height: 40,
        // object-fit 要走 NImage 自己的 prop: 写在 imgProps.style 里会被它
        // 追加在后面的 object-fit(默认 fill) 覆盖掉 —— 实测过
        objectFit: 'cover',
        imgProps: {
          alt: row.title,
          style: { 'border-radius': '2px', 'width': '40px', 'height': '40px' },
        },
        src: convertImgUrl(row.img),
        fallbackSrc: IMG_PLACEHOLDER,
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
      // 导入的文章是草稿且不带分类, category 为 null, 不能直接取 name
      return h('div', row.category?.name || 'None')
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
                onClick: () => handleRestore(row),
              },
              { default: () => 'Restore', icon: () => h('i', { class: 'i-majesticons:eye-line' }) },
            )
          : h(
              NButton,
              {
                size: 'small',
                type: 'primary',
                secondary: true,
                onClick: () => router.push(`/article/write/${row.id}`), // 携带参数前往 写文章 页面
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
  // 必须 return, 否则 useCRUD 里 await 到 undefined: 删除成功不提示, 失败也无法捕获
  return extraParams.value.is_delete
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
    // 乐观更新后失败必须回滚, 否则开关显示已开、后端其实没变
    // (菜单/接口/用户三个页面的同类开关都是这么写的)
    row.is_top = !row.is_top
    console.error(err)
  }
  finally {
    row.publishing = false
  }
}

// 从回收站恢复: 相邻的置顶/审核/删除都有提示, 只有这里原来既不 catch 也不提示,
// 失败时是未捕获 rejection, 成功时只能靠列表刷新猜
async function handleRestore(row) {
  if (!row.id) {
    return
  }
  try {
    await api.softDeleteArticle([row.id], false)
    $message?.success('已恢复该文章')
    await $table.value?.handleSearch()
  }
  catch (err) {
    console.error(err)
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
  // 后端 F9 之后同时接受 .md 与 .markdown, 两边保持一致
  const name = data.file.name.toLowerCase()
  if (!name.endsWith('.md') && !name.endsWith('.markdown')) {
    $message.error('只能上传 .md / .markdown 格式的文件，请重新上传')
    return false
  }
  return true
}

// Operations after file upload
function afterUpload({ event }) {
  // 网关拦截或鉴权失败时响应不是 JSON, 不能直接 JSON.parse
  const res = parseJson(event?.target?.response)
  if (res?.code === 0) {
    $table.value?.handleSearch()
    $message.success('Article imported successfully!')
  }
  else {
    $message.error(res?.message || 'Article import failed!')
  }
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
      <!-- 只有「新建文章」是主操作(填充), 其余三个用 secondary:
           原来四个按钮四种填充色平铺, 看不出主次 -->
      <NButton
        type="error"
        secondary
        :disabled="!$table?.selections.length"
        @click="handleDelete($table?.selections)"
      >
        <template #icon>
          <i class="i-material-symbols:recycling-rounded" />
        </template>
        Batch Delete
      </NButton>
      <NButton
        secondary
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
          :headers="{ Authorization: `Bearer ${authStore.token}` }"
          :show-file-list="false"
          multiple
          @before-upload="beforeUpload"
          @finish="afterUpload"
        >
          <NButton secondary>
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
        <QueryItem label="Category" :label-width="40" :content-width="160">
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
