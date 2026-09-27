<script setup>
import { NIcon, NText, NUpload, NUploadDragger } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useAuthStore } from '@/store'
import { convertImgUrl, parseJson } from '@/utils'

const props = defineProps({
  preview: {
    type: String,
    default: '',
  },
  width: {
    type: Number,
    default: 120,
  },
})

const emit = defineEmits(['update:preview'])

// 不解构 token: 解构后失去响应性, 重新登录仍会用旧值
const authStore = useAuthStore()
const previewImg = ref(props.preview)

watch(() => props.preview, val => previewImg.value = val)

// Upload image
function handleImgUpload({ event }) {
  // 网关拦截或鉴权失败时响应不是 JSON, 不能直接 JSON.parse
  const res = parseJson(event?.target?.response)
  if (res?.code !== 0) {
    $message?.error(res?.message || '图片上传失败')
    return
  }
  previewImg.value = res.data
  emit('update:preview', previewImg.value)
}

// Determine whether it's a locally uploaded image or a network resource
// In development you can use local file upload; in production, cloud storage is recommended
const imgUrl = computed(() => convertImgUrl(previewImg.value))

defineExpose({ previewImg })
</script>

<template>
  <div>
    <NUpload
      action="/api/upload"
      :headers="{ Authorization: `Bearer ${authStore.token}` }"
      accept="image/jpeg,image/png,image/gif,image/webp,image/bmp"
      :show-file-list="false"
      @finish="handleImgUpload"
    >
      <template v-if="previewImg">
        <img
          border-color="#d9d9d9"
          class="cursor-pointer border-2 rounded-lg border-dashed hover:border-color-lightblue"
          :style="{ width: `${props.width}px` }"
          :src="imgUrl"
          alt="Article Cover"
        >
      </template>
      <template v-else>
        <NUploadDragger>
          <div class="mb-3">
            <NIcon size="50" :depth="3">
              <span class="i-mdi:upload" />
            </NIcon>
          </div>
          <NText>
            Click or drag file to this area to upload
          </NText>
        </NUploadDragger>
      </template>
    </NUpload>
  </div>
</template>
