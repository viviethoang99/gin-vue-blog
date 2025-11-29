<script setup>
import { onMounted, ref } from 'vue'

import LinkList from './components/LinkList.vue'
import AddLink from './components/AddLink.vue'
import Comment from '@/components/comment/Comment.vue'
import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const loading = ref(true)
const linkList = ref([])

onMounted(() => {
  api.getLinks().then((res) => {
    linkList.value = res.data
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

