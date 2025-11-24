<script setup>
import { h, nextTick, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NDynamicTags, NForm, NFormItem, NInput, NRadio, NRadioGroup, NSelect, NSpace, NSwitch, NTag } from 'naive-ui'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import CommonPage from '@/components/common/CommonPage.vue'
import CrudModal from '@/components/crud/CrudModal.vue'
import UploadOne from '@/components//UploadOne.vue'

import { articleTypeOptions } from '@/assets/config'
import { useTagStore } from '@/store'
import api from '@/api'

defineOptions({ name: 'Publish Article' })

const route = useRoute()
// const router = useRouter()
const tagStore = useTagStore()

const categoryOptions = ref([]) // Category options
const tagOptions = ref([]) // Tag options
let backTagOptions = [] // Backup tag options

// Fix the issue where viewing multiple articles at the same time, switching tabs doesn't refresh
watch(route, async () => tagStore.reloadTag())

onMounted(async () => {
  fetchData()
})

onActivated(async () => {
  fetchData()
})

async function fetchData() {
  getArticleInfo()
  api.getCategoryOption().then((resp) => {
    categoryOptions.value = resp.data.map(e => ({ value: e.label, label: e.label }))
  })
  api.getTagOption().then((resp) => {
    tagOptions.value = resp.data.map(e => ({ value: e.label, label: e.label }))
    backTagOptions = tagOptions.value
  })
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
  if (!id) {
    formModel.value = { status: 1, is_top: false, title: '', type: 1 }
    return
  }

  // With id means editing article
  window.$loadingBar?.start()
  try {
    const resp = await api.getArticleById(id)
    const { category, tags } = resp.data
    formModel.value = resp.data
    formModel.value.tag_names = tags.map(e => e.name)
    formModel.value.category_name = category.name
    window.$loadingBar?.finish()
  }
  catch (err) {
    window.$loadingBar?.error()
    $message?.error('Loading failed')
  }
}

// TODO: Save draft
function handleDraft() {
  $message.info('Save draft feature is under development')
}

// Publish article
function handlePublish() {
  if (!formModel.value.title || !formModel.value.title?.trim()) {
    formModel.value.title = formModel.value.title?.trim()
    $message.info('Please enter title')
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
        $message.success('Operation successful!')
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
    message: 'Please select article category',
    trigger: ['blur', 'change'],
  },
  tag_names: {
    required: true,
    message: 'Please select article tags',
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
</script>

<template>
  <CommonPage :show-header="false" title="Write Article">
    <div class="mb-4 flex items-center bg-white space-x-2">
      <NInput
        v-model:value="formModel.title"
        type="text"
        class="mr-5 flex-1 py-1 text-lg color-primary font-bold"
        placeholder="Enter article title..."
      />
      <NButton ghost type="error" :loading="btnLoading" @click="handleDraft">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:uploading-loop" />
        </template>
        Save Draft
      </NButton>
      <NButton type="error" :loading="btnLoading" @click="handlePublish">
        <template #icon>
          <span v-if="!btnLoading" class="i-line-md:confirm-circle" />
        </template>
        Publish Article
      </NButton>
    </div>

    <!-- TODO: File upload -->
    <MdEditor v-model="formModel.content" style="height: calc(100vh - 245px)" />

    <CrudModal
      v-model:visible="modalVisible"
      title="Publish Article"
      :loading="btnLoading"
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
        <NFormItem label="Category" path="category_name">
          <NSelect
            v-model:value="formModel.category_name"
            style="width: 50%"
            clearable filterable tag
            placeholder="Search keywords, press enter to add"
            :options="categoryOptions"
          />
        </NFormItem>
        <NFormItem label="Tags" path="tag_names">
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
                placeholder="Tag name"
                @update:value="{
                  submit($event);
                  newTag = null;
                }"
                @blur="deactivate"
              >
                <template #action>
                  Enter tag name to search, press enter to add custom tag
                </template>
              </NSelect>
            </template>
          </NDynamicTags>
        </NFormItem>
        <NFormItem label="Type" path="type">
          <NSelect
            v-model:value="formModel.type"
            style="width: 50%"
            placeholder="Please select article type"
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
          label="Original URL" path="original_url"
        >
          <NInput
            v-model:value="formModel.original_url"
            type="text"
            placeholder="Please enter original article link"
          />
        </NFormItem>
        <NFormItem label="Thumbnail" path="img">
          <UploadOne
            v-model:preview="formModel.img"
            :width="220"
          />
        </NFormItem>
        <NFormItem label="Pin to Top" path="is_top">
          <NSwitch v-model:value="formModel.is_top" />
        </NFormItem>
        <NFormItem label="Status" path="status">
          <NRadioGroup v-model:value="formModel.status" name="radiogroup">
            <NSpace>
              <NRadio :value="1">
                Public
              </NRadio>
              <NRadio :value="2">
                Private
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
