<script setup>
import { defineComponent, h, nextTick, onMounted, ref } from 'vue'
import { NInput } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import QueryItem from '@/components/crud/QueryItem.vue'
import CrudTable from '@/components/crud/CrudTable.vue'

import api from '@/api'

defineOptions({ name: 'Category Management' })

const $table = ref(null)
const queryItems = ref({
  keyword: '',
})

// Index of the currently editing row
const editIndex = ref(-1)

onMounted(() => {
  $table.value?.handleSearch()
})

const ShowOrEdit = defineComponent({
  props: {
    value: [String, Number],
    onUpdateValue: [Function, Array],
    rowIndex: [Number],
  },
  setup(props) {
    const inputRef = ref(null)
    const inputValue = ref(props.value)

    function handleOnClick() {
      editIndex.value = props.rowIndex
      nextTick(() => {
        inputRef.value.focus()
      })
    }
    function handleChange() {
      props.onUpdateValue(inputValue.value)
      // isEdit.value = false
    }
    return () =>
      h(
        'div',
        {
          style: 'min-height: 22px',
          onClick: handleOnClick,
        },
        editIndex.value === props.rowIndex
          ? h(NInput, {
            ref: inputRef,
            value: inputValue.value,
            onUpdateValue: (v) => {
              inputValue.value = v
            },
            onChange: handleChange,
            onBlur: handleChange,
          })
          : props.value,
      )
  },
})

const columns = [
  { type: 'selection', width: 15, fixed: 'left' },
  {
    title: 'Created Date',
    key: 'created_at',
    width: 80,
    align: 'center',
    render(row, index) {
      return h(
        ShowOrEdit,
        {
          value: row.created_at,
          rowIndex: index,
          onUpdateValue(v) {
            data.value[index].name = v
          },
        },
      )
    },
  },
  {
    title: 'Updated Date',
    key: 'updated_at',
    width: 80,
    align: 'center',
    render(row, index) {
      return h(
        ShowOrEdit,
        {
          value: row.updated_at,
          rowIndex: index,
          onUpdateValue(v) {
            data.value[index].name = v
          },
        },
      )
    },
  },
  {
    title: 'Actions',
    key: 'action',
    width: 80,
    align: 'center',
    render(row, index) {
      if (editIndex.value === index) {
        return h(
          'button',
          {
            onClick() {
              editIndex.value = -1
            },
          },
          'Close Edit',
        )
      }
      else {
        return h(
          'button',
          {
            onClick() {
              editIndex.value = index
            },
          },
          'Enter Edit',
        )
      }
    },
  },
]
</script>

<template>
  <CommonPage title="Test Page">
    <CrudTable
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="api.getCategorys"
    >
      <template #queryBar>
        <QueryItem label="Category" :label-width="80">
          <NInput
            v-model:value="queryItems.keyword"
            clearable
            type="text"
            placeholder="Enter category name"
            @keydown.enter="$table?.handleSearch()"
          />
        </QueryItem>
      </template>
    </CrudTable>
  </CommonPage>
</template>
