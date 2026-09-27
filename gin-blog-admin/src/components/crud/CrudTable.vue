<script setup>
import { NButton, NDataTable, NSpace } from 'naive-ui'
import { nextTick, reactive, ref } from 'vue'

const props = defineProps({
  /** Whether to not set column dividers */
  singleLine: { type: Boolean, default: false },
  /** true: backend pagination false: frontend pagination */
  remote: { type: Boolean, default: true },
  /** Whether to enable pagination */
  isPagination: { type: Boolean, default: true },
  /** Horizontal width of table content */
  scrollX: { type: Number, default: 1200 },
  /** Primary key name */
  rowKey: { type: String, default: 'id' },
  /** Columns to display */
  columns: { type: Array, required: true },
  /** Parameters in queryBar */
  queryItems: {
    type: Object,
    default() { return {} },
  },
  /** Additional parameters (optional) */
  extraParams: {
    type: Object,
    default() { return {} },
  },
  /**
   * TODO: How to handle if you want both url and body at the same time
   * Request API to get data
   */
  getData: {
    type: Function,
    required: true,
  },
})

const emit = defineEmits(['update:queryItems', 'checked', 'dataChange', 'sorterChange'])

const loading = ref(false) // Loading
const selections = ref([]) // Multiple selected rowKeys
const tableData = ref([]) // Table data
const initQuery = { ...props.queryItems }

// Pagination configuration
// ! 翻页只能有一个入口: naive-ui 内部的 mergedOnUpdatePage 会先调 pagination.onChange
// ! 再触发组件的 onUpdate:page, 两处都发请求的话每次翻页会请求两次, 且慢的响应会覆盖新的
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [5, 10, 20],
  onUpdatePageSize: (pageSize) => {
    pagination.page = 1
    pagination.pageSize = pageSize
    handleQuery()
  },
  prefix({ itemCount }) {
    return `共 ${itemCount} 条`
  },
})

async function handleQuery() {
  selections.value = [] // Reset selection

  try {
    loading.value = true
    let paginationParams = {}
    // If not pagination mode or using frontend pagination, no need to pass pagination parameters
    if (props.isPagination && props.remote) {
      paginationParams = {
        page_num: pagination.page,
        page_size: pagination.pageSize,
      }
    }
    const { data } = await props.getData({
      ...props.queryItems,
      ...props.extraParams,
      ...paginationParams,
    })
    // data 为 null 时原来会在 data.length 上抛 TypeError, 又被下面的空 catch 吞掉,
    // 表现是「点了搜索什么都没发生」。这里显式当成空列表处理
    const list = data?.page_data ?? data ?? []
    tableData.value = Array.isArray(list) ? list : []
    pagination.itemCount = data?.total ?? tableData.value.length
  }
  catch (err) {
    console.error(err)
    window.$message?.error('数据加载失败, 请重试')
    tableData.value = []
    pagination.itemCount = 0
  }
  finally {
    emit('dataChange', tableData.value)
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1 // Go back to page 1
  handleQuery()
}

async function handleReset() {
  const queryItems = { ...props.queryItems } // Reset search parameters
  for (const key in queryItems) {
    queryItems[key] = null // Note the type
  }
  emit('update:queryItems', { ...queryItems, ...initQuery })
  await nextTick()
  pagination.page = 1 // Go back to page 1
  handleQuery()
}

function onPageChange(currentPage) {
  pagination.page = currentPage
  props.remote && handleQuery()
}

function onChecked(rowKeys) {
  selections.value = rowKeys
  // Contains selection
  if (props.columns.some(item => item.type === 'selection')) {
    emit('checked', rowKeys)
  }
}

function onSorterChange(sorter) {
  emit('sorterChange', sorter)
}

async function handleExport(columns = props.columns, data = tableData.value) {
  if (!data?.length) {
    return window.$message.warning('No data available')
  }
  const columnsData = columns.filter(item => !!item.title && !item.hideInExcel)
  const thKeys = columnsData.map(item => item.key)
  const thData = columnsData.map(item => item.title)
  const trData = data.map(item => thKeys.map(key => item[key]))

  // exceljs 不提供浏览器端下载, 需自行生成 Blob 触发
  // 本函数为 async, 失败时的 rejection 不会被 Vue 的错误处理捕获, 故在此兜住
  try {
    // exceljs 体积较大(gzip 后约 300KB), 点击导出时才加载, 不进入首屏
    const { default: ExcelJS } = await import('exceljs')
    const workbook = new ExcelJS.Workbook()
    const sheet = workbook.addWorksheet('Data Report')
    sheet.addRows([thData, ...trData])

    const buffer = await workbook.xlsx.writeBuffer()
    const blob = new Blob([buffer], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'Data Report.xlsx'
    link.click()
    URL.revokeObjectURL(url)
  }
  catch {
    window.$message.error('导出失败')
  }
}

defineExpose({
  handleQuery,
  handleSearch,
  handleReset,
  handleExport,
  selections,
  tableData,
})
</script>

<template>
  <div
    v-if="$slots.queryBar"
    class="mb-7 min-h-[60px] flex items-start justify-between border border-gray-200 border-gray-400 rounded-2 border-solid bg-gray-50 p-3.5 dark:bg-black dark:bg-opacity-5"
  >
    <NSpace wrap :size="[35, 15]">
      <slot name="queryBar" />
    </NSpace>
    <div class="flex-shrink-0 space-x-4">
      <NButton ghost type="primary" @click="handleReset">
        <template #icon>
          <i class="i-lucide:rotate-ccw" />
        </template>
        Reset
      </NButton>
      <NButton type="primary" @click="handleSearch">
        <template #icon>
          <i class="i-fe:search" />
        </template>
        Search
      </NButton>
      <!-- TODO: Add extra slots to let users customize buttons -->
    </div>
  </div>
  <NDataTable
    :remote="remote"
    :loading="loading"
    :scroll-x="scrollX"
    :columns="columns"
    :data="tableData"
    :row-key="(row) => row[rowKey]"
    :single-line="singleLine"
    :pagination="isPagination ? pagination : false"
    :checked-row-keys="selections"
    @update:checked-row-keys="onChecked"
    @update:page="onPageChange"
    @update:sorter="onSorterChange"
  />
</template>
