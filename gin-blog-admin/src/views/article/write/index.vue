<script setup>
import { config as configureMdEditor, en_US, MdEditor } from 'md-editor-v3'
import { NButton, NDynamicTags, NForm, NFormItem, NInput, NRadio, NRadioGroup, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'
import { h, nextTick, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'

import { articleTypeOptions } from '@/assets/config'
import UploadOne from '@/components//UploadOne.vue'
import CommonPage from '@/components/common/CommonPage.vue'

import CrudModal from '@/components/crud/CrudModal.vue'
import { useTagStore } from '@/store'
import { convertImgUrl, request } from '@/utils'
import 'md-editor-v3/lib/style.css'

defineOptions({ name: 'Publish Article' })

configureMdEditor({
  editorConfig: {
    languageUserDefined: {
      'vi-VN': {
        ...en_US,
        toolbarTips: {
          bold: 'In đậm',
          underline: 'Gạch chân',
          italic: 'In nghiêng',
          strikeThrough: 'Gạch ngang',
          title: 'Tiêu đề',
          sub: 'Chỉ số dưới',
          sup: 'Chỉ số trên',
          quote: 'Trích dẫn',
          unorderedList: 'Danh sách không thứ tự',
          orderedList: 'Danh sách có thứ tự',
          task: 'Danh sách công việc',
          codeRow: 'Mã nội dòng',
          code: 'Khối mã',
          link: 'Liên kết',
          image: 'Hình ảnh',
          table: 'Bảng',
          mermaid: 'Sơ đồ Mermaid',
          katex: 'Công thức',
          revoke: 'Hoàn tác',
          next: 'Làm lại',
          save: 'Lưu',
          prettier: 'Định dạng',
          pageFullscreen: 'Toàn màn hình trong trang',
          fullscreen: 'Toàn màn hình',
          preview: 'Xem trước',
          previewOnly: 'Chỉ xem trước',
          htmlPreview: 'Xem mã HTML',
          catalog: 'Mục lục',
          github: 'Mã nguồn',
        },
        titleItem: {
          h1: 'Tiêu đề cấp 1',
          h2: 'Tiêu đề cấp 2',
          h3: 'Tiêu đề cấp 3',
          h4: 'Tiêu đề cấp 4',
          h5: 'Tiêu đề cấp 5',
          h6: 'Tiêu đề cấp 6',
        },
        imgTitleItem: {
          link: 'Thêm liên kết ảnh',
          upload: 'Tải ảnh lên',
          clip2upload: 'Cắt và tải lên',
        },
        linkModalTips: {
          linkTitle: 'Thêm liên kết',
          imageTitle: 'Thêm ảnh',
          descLabel: 'Mô tả:',
          descLabelPlaceHolder: 'Nhập mô tả...',
          urlLabel: 'Liên kết:',
          urlLabelPlaceHolder: 'Nhập liên kết...',
          buttonOK: 'Đồng ý',
        },
        clipModalTips: {
          title: 'Cắt ảnh',
          buttonUpload: 'Tải lên',
        },
        copyCode: {
          text: 'Sao chép',
          successTips: 'Đã sao chép!',
          failTips: 'Sao chép thất bại!',
        },
        footer: {
          markdownTotal: 'Số ký tự',
          scrollAuto: 'Cuộn đồng bộ',
        },
      },
    },
  },
})

const route = useRoute()
// const router = useRouter()
const tagStore = useTagStore()

const categoryOptions = ref([]) // Category options
const tagOptions = ref([]) // Tag options
let backTagOptions = [] // Backup tag options

// Fix the issue where viewing multiple articles at the same time, switching tabs doesn't refresh
// watch 的必须是 getter 而不是 route 本身: useRoute() 返回的是响应式对象的浅代理,
// 直接传进 watch 会报 "Invalid watch source" 并且完全不触发
watch(() => route.fullPath, async () => tagStore.reloadTag())

onMounted(async () => {
  fetchData()
})

onActivated(async () => {
  fetchData()
})

async function fetchData() {
  getArticleInfo()
  // 拦截器已经弹过错误提示, 这里补 catch 只是别留下 unhandled rejection
  api.getCategoryOption().then((resp) => {
    categoryOptions.value = resp.data.map(e => ({ value: e.label, label: e.label }))
  }).catch(err => console.error(err))
  api.getTagOption().then((resp) => {
    tagOptions.value = resp.data.map(e => ({ value: e.label, label: e.label }))
    backTagOptions = tagOptions.value
  }).catch(err => console.error(err))
  await nextTick()
}

const formRef = ref(null)
const formModel = ref({
  title: '',
  status: 1, // Publish status: default public
  is_top: false, // Default not pinned
  type: 1, // Default original
  tag_names: [],
  category_name: '',
})
const btnLoading = ref(false)
const modalVisible = ref(false)
const newTag = ref(null) // New tag

// Listen to selected tags, update selectable tags in real time
watch(() => formModel.value.tag_names, (newVal) => {
  tagOptions.value = backTagOptions.filter(e => !newVal.includes(e.label))
}, { deep: true })

// Get article information based on id parameter in route
async function getArticleInfo() {
  const id = route.params.id // Get parameter from route

  // No id means creating new article
  // 必须带上 tag_names 和 category_name, 否则 watch tag_names 的回调会拿到 undefined
  if (!id) {
    formModel.value = { status: 1, is_top: false, title: '', type: 1, tag_names: [], category_name: '' }
    return
  }

  // With id means editing article
  window.$loadingBar?.start()
  try {
    const resp = await api.getArticleById(id)
    const { category, tags } = resp.data
    formModel.value = resp.data
    // 导入生成的草稿没有分类和标签, category 为 null
    formModel.value.tag_names = tags?.map(e => e.name) ?? []
    formModel.value.category_name = category?.name ?? ''
    window.$loadingBar?.finish()
  }
  catch {
    window.$loadingBar?.error()
    $message?.error('Tải dữ liệu thất bại')
  }
}

// TODO: Save draft
function handleDraft() {
  $message.info('Tính năng lưu nháp đang được phát triển')
}

// Publish article
function handlePublish() {
  if (!formModel.value.title || !formModel.value.title?.trim()) {
    formModel.value.title = formModel.value.title?.trim()
    $message.info('Vui lòng nhập tiêu đề')
    return
  }
  modalVisible.value = true
}

// Save
async function handleSave() {
  formRef.value?.validate(async (err) => {
    if (!err) {
      btnLoading.value = true
      // $message.loading('Saving...')
      try {
        await api.saveOrUpdateArticle(formModel.value)
        modalVisible.value = false
        $message.success('Thao tác thành công!')
        // Close current tab and return to article list
        tagStore.removeTag(route.path)
        // await router.replace({ path: '/article/list', query: { needRefresh: true } })
      }
      catch (err) {
        console.error(err)
      }
      finally {
        btnLoading.value = false
      }
    }
  })
}

const rules = {
  category_name: {
    required: true,
    message: 'Vui lòng chọn danh mục',
    trigger: ['blur', 'change'],
  },
  tag_names: {
    required: true,
    message: 'Vui lòng chọn thẻ',
  },
}

// Render tags
function renderTag(tag, index) {
  return h(
    NTag,
    {
      type: 'info',
      disabled: index > 3,
      closable: true,
      onClose: () => formModel.value.tag_names.splice(index, 1),
    },
    { default: () => tag },
  )
}

// MdEditor only supplies the selected files. Upload them through the same Axios
// instance as the rest of the admin app so JWT and business-error handling are
// applied consistently, then give the resulting URLs back to the editor.
async function handleEditorUpload(files, callback) {
  try {
    const urls = await Promise.all(files.map(async (file) => {
      const formData = new FormData()
      formData.append('file', file)

      // Do not set Content-Type manually: Axios/browser must append the
      // multipart boundary for the Go server to parse FormFile("file").
      const resp = await request.post('/upload', formData)
      const url = convertImgUrl(resp.data)
      return url
    }))

    callback(urls)
  }
  catch {
    $message?.error('Tải ảnh thất bại')
  }
}
</script>

<template>
  <CommonPage :show-header="false" title="Viết bài">
    <div class="mb-4 flex items-center bg-white space-x-2">
      <NInput
        v-model:value="formModel.title"
        type="text"
        class="mr-5 flex-1 py-1 text-lg color-primary font-bold"
        placeholder="Nhập tiêu đề bài viết..."
      />
      <NButton ghost type="error" :loading="btnLoading" @click="handleDraft">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:uploading-loop" />
        </template>
        Lưu nháp
      </NButton>
      <NButton type="error" :loading="btnLoading" @click="handlePublish">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:confirm-circle" />
        </template>
        Xuất bản
      </NButton>
    </div>

    <MdEditor
      v-model="formModel.content"
      style="height: calc(100vh - 245px)"
      language="vi-VN"
      @on-upload-img="handleEditorUpload"
    />

    <CrudModal
      v-model:visible="modalVisible"
      title="Xuất bản bài viết"
      :loading="btnLoading"
      cancel-text="Hủy"
      ok-text="Xác nhận"
      show-footer
      @save="handleSave"
    >
      <NForm
        ref="formRef"
        label-placement="left"
        label-align="left"
        :label-width="100"
        :model="formModel"
        :rules="rules"
      >
        <NFormItem label="Danh mục" path="category_name">
          <NSelect
            v-model:value="formModel.category_name"
            style="width: 50%"
            clearable filterable tag
            placeholder="Tìm kiếm, nhấn Enter để thêm"
            :options="categoryOptions"
          />
        </NFormItem>
        <NFormItem label="Thẻ" path="tag_names">
          <NDynamicTags
            v-model:value="formModel.tag_names"
            :render-tag="renderTag"
            :max="3"
          >
            <template #input="{ submit, deactivate }">
              <NSelect
                v-model:value="newTag"
                size="small" filterable tag clearable
                :options="tagOptions"
                placeholder="Tên thẻ"
                @update:value="{
                  submit($event);
                  newTag = null;
                }"
                @blur="deactivate"
              >
                <template #action>
                  Nhập tên thẻ để tìm kiếm, nhấn Enter để thêm thẻ mới
                </template>
              </NSelect>
            </template>
          </NDynamicTags>
        </NFormItem>
        <NFormItem label="Loại bài viết" path="type">
          <NSelect
            v-model:value="formModel.type"
            style="width: 50%"
            placeholder="Chọn loại bài viết"
            :options="articleTypeOptions"
          />
        </NFormItem>
        <!-- <n-form-item label="Article Description" path="desc">
          <n-input
            v-model:value="formModel.desc"
            placeholder="Please enter article description"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
          />
        </n-form-item> -->
        <NFormItem
          v-if="(formModel.type === 2 || formModel.type === 3)"
          label="Liên kết gốc" path="original_url"
        >
          <NInput
            v-model:value="formModel.original_url"
            type="text"
            placeholder="Nhập liên kết bài viết gốc"
          />
        </NFormItem>
        <NFormItem label="Ảnh đại diện" path="img">
          <UploadOne
            v-model:preview="formModel.img"
            :width="220"
          />
        </NFormItem>
        <NFormItem label="Ghim lên đầu" path="is_top">
          <NSwitch v-model:value="formModel.is_top" />
        </NFormItem>
        <NFormItem label="Trạng thái" path="status">
          <NRadioGroup v-model:value="formModel.status" name="radiogroup">
            <NSpace>
              <NRadio :value="1">
                Công khai
              </NRadio>
              <NRadio :value="2">
                Riêng tư
              </NRadio>
            </NSpace>
          </NRadioGroup>
        </NFormItem>
      </NForm>
    </CrudModal>
  </CommonPage>
</template>

<style lang="scss" scoped>
.md-preview {
  ul,
  ol {
    list-style: revert;
  }
}
</style>
