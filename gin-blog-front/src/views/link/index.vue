<script setup>
import { onMounted, ref } from 'vue'

import api from '@/api'
import BannerPage from '@/components/BannerPage.vue'
import Comment from '@/components/comment/Comment.vue'
import AddLink from './components/AddLink.vue'
import LinkList from './components/LinkList.vue'

const loading = ref(true)
const linkList = ref([])

onMounted(() => {
  // 补 catch: 以前失败会留下未捕获的 rejection
  // data 为空也不能给 LinkList 传 null, 它里面直接取 linkList.length
  api.getLinks().then((res) => {
    linkList.value = res.data ?? []
  }).catch((err) => {
    console.error(err)
  }).finally(() => {
    loading.value = false
  })
})
</script>

<template>
  <BannerPage label="link" title="Friendly Links" card :loading="loading">
    <div class="space-y-5">
      <!-- Link list -->
      <LinkList :link-list="linkList" />
      <!-- Add link -->
      <AddLink />
      <!-- Comments -->
      <Comment class="mt-30" :type="2" />
    </div>
  </BannerPage>
</template>
