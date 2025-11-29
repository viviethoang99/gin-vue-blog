<script setup>
import { nextTick, reactive, ref } from 'vue'
import { NButton, NDataTable, NSpace } from 'naive-ui'
import { utils, writeFile } from 'xlsx'

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
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [5, 10, 20],
  onChange: (page) => {
    pagination.page = page
    handleQuery()
  },
  onUpdatePageSize: (pageSize) => {
    pagination.page = 1
    pagination.pageSize = pageSize
    handleQuery()
  },
  prefix({ itemCount }) {
    return `Total ${itemCount} items`
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
    tableData.value = data?.page_data || data
    pagination.itemCount = data?.total ?? data.length
  }
  catch (error) {
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

function handleExport(columns = props.columns, data = tableData.value) {
  if (!data?.length) {
    return window.$message.warning('No data available')
  }
  const columnsData = columns.filter(item => !!item.title && !item.hideInExcel)
  const thKeys = columnsData.map(item => item.key)
  const thData = columnsData.map(item => item.title)
  const trData = data.map(item => thKeys.map(key => item[key]))
  const sheet = utils.aoa_to_sheet([thData, ...trData])
  const workBook = utils.book_new()
  utils.book_append_sheet(workBook, sheet, 'Data Report')
  writeFile(workBook, 'Data Report.xlsx')
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
