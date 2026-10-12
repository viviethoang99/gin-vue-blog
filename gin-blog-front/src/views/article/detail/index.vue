<script setup>
import { useWindowScroll, useWindowSize } from '@vueuse/core'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import { marked } from 'marked'
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'

import Comment from '@/components/comment/Comment.vue'
import AppFooter from '@/components/layout/AppFooter.vue'
import { convertImgUrl } from '@/utils'
import { addCopyButtons } from '@/utils/code-block'
import { typesetMath } from '@/utils/mathjax'
import BannerInfo from './components/BannerInfo.vue'
import Catalogue from './components/Catalogue.vue'
import Copyright from './components/Copyright.vue'
import Forward from './components/Forward.vue'
import LastNext from './components/LastNext.vue'

import LatestList from './components/LatestList.vue'
import Recommend from './components/Recommend.vue'

import Reward from './components/Reward.vue'
import 'highlight.js/styles/a11y-dark.css'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('javascript', javascript)

const route = useRoute()

const data = ref({
  id: 0,
  title: '',
  content: '',
  created_at: '',
  updated_at: '',
  like_count: 0,
  view_count: 0,
  comment_count: 0,
  img: '',
  newest_articles: [],
  tags: [],
  category: {},
  next_article: {},
  last_article: {},
  recommend_articles: [],
})

// Nội dung bài viết
const previewRef = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const resp = await api.getArticleDetail(route.params.id)
    if (!resp.data) {
      return
    }
    // Các trường này từ backend có thể là null (chưa phân loại, không có bài trước, không có đề xuất),
    // trong khi component con truy cập trực tiếp .length / .name, nên chuẩn hóa về giá trị mặc định tại đây
    // thay vì ghi đè toàn bộ bằng resp.data
    data.value = {
      ...data.value,
      ...resp.data,
      tags: resp.data.tags ?? [],
      newest_articles: resp.data.newest_articles ?? [],
      recommend_articles: resp.data.recommend_articles ?? [],
      category: resp.data.category ?? {},
      last_article: resp.data.last_article ?? {},
      next_article: resp.data.next_article ?? {},
      // marked phân tích văn bản Markdown; không truyền trực tiếp nội dung rỗng vào marked
      content: await marked.parse(resp.data.content ?? '', { async: true }),
    }
    await nextTick()
    // Tô sáng mã nguồn bằng highlight.js
    document.querySelectorAll('pre code').forEach(el => hljs.highlightElement(el))
    // Thêm nút "Sao chép" vào khối mã
    addCopyButtons(previewRef.value)
    // Chỉ tải MathJax khi nội dung có công thức
    await typesetMath(data.value.content)
  }
  catch (err) {
    console.error(err)
  }
  finally {
    loading.value = false
  }
})

// Tiến độ đọc: đường mảnh ở đầu trang, tính theo tỷ lệ vị trí cuộn trên tổng chiều cao có thể cuộn
const { y } = useWindowScroll()
const { height: windowHeight } = useWindowSize()
const readProgress = computed(() => {
  // Phải đọc y trước khi return sớm: khi nội dung ở màn hình đầu chưa render xong, trang chưa có thanh cuộn;
  // nếu return trước thì y không được ghi nhận là dependency, các lần cuộn sau sẽ không tính lại và thanh tiến độ luôn ở 0
  const scrolled = y.value
  const total = document.documentElement.scrollHeight - windowHeight.value
  if (total <= 0) {
    return 0
  }
  return Math.min(100, Math.max(0, (scrolled / total) * 100))
})

// Cảnh báo với bài viết lâu chưa cập nhật: nội dung kỹ thuật nhanh lỗi thời, tránh để độc giả làm theo thông tin cũ
const STALE_DAYS = 90
const staleDays = computed(() => {
  if (!data.value.updated_at) {
    return 0
  }
  const days = Math.floor((Date.now() - new Date(data.value.updated_at).getTime()) / 86400000)
  return days >= STALE_DAYS ? days : 0
})

const styleVal = computed(() =>
  data.value.img
    ? `background: url('${convertImgUrl(data.value.img)}') center center / cover no-repeat;`
    : 'background: rgba(0,0,0,0.1) center center / cover no-repeat;',
)
</script>

<template>
  <!-- Tiến độ đọc -->
  <div
    class="fixed inset-x-0 top-0 z-999 h-0.5 bg-#49b1f5"
    :style="{ width: `${readProgress}%`, transition: 'width .1s linear' }"
  />
  <!-- Phần đầu trang -->
  <div :style="styleVal" class="banner-fade-down absolute inset-x-0 top-0 h-[360px] f-c-c lg:h-[400px]">
    <BannerInfo v-if="!loading" :article="data" />
  </div>
  <!-- Nội dung chính -->
  <main class="flex-1">
    <div class="card-fade-up grid grid-cols-12 mx-auto mb-3 mt-[380px] gap-4 px-1 lg:mt-[440px] lg:max-w-[1200px]">
      <!-- Phần chính của bài viết -->
      <div class="card-view col-span-12 mx-2 pt-7 lg:col-span-9 lg:mx-0">
        <!-- Nội dung bài viết -->
        <!-- Cảnh báo bài viết cũ -->
        <div
          v-if="!loading && staleDays"
          class="mb-5 border-l-4 border-#f0ad4e rounded bg-#f0ad4e/10 px-4 py-2 text-sm lg:mx-10"
        >
          Bài viết này được cập nhật lần cuối cách đây {{ staleDays }} ngày. Một số nội dung có thể đã lỗi thời.
        </div>
        <article
          ref="previewRef"
          class="max-w-none prose prose-truegray lg:mx-10 dark:prose-invert"
          v-html="data.content"
        />
        <!-- Thông báo bản quyền -->
        <Copyright class="my-5 lg:mx-5" />
        <!-- Thẻ và chia sẻ -->
        <Forward :tag-list="data.tags" class="mb-12 lg:mx-5" />
        <!-- Thích và ủng hộ -->
        <Reward
          :article-id="data.id"
          :like-count="data.like_count"
          class="mb-10"
        />
        <!-- Bài trước và bài tiếp theo -->
        <LastNext
          :last-article="data.last_article"
          :next-article="data.next_article"
          class="lg:mx-5"
        />
        <!-- Bài viết đề xuất -->
        <Recommend
          :recommend-list="data.recommend_articles"
          class="mt-7 lg:mx-5"
        />
        <!-- Đường phân cách -->
        <hr class="my-10 border-2 border-color-divider border-dashed lg:mx-5">
        <!-- Bình luận bài viết -->
        <Comment :type="1" class="lg:mx-5" />
      </div>
      <!-- Thanh bên bài viết -->
      <div class="col-span-0 lg:col-span-3">
        <div class="sticky top-5 hidden lg:block space-y-4">
          <!-- Mục lục -->
          <!-- TODO: Cách dùng v-if chưa tối ưu; cần xử lý việc component cha lấy dữ liệu từ API và component con render -->
          <Catalogue v-if="!loading" :preview-ref="previewRef" />
          <!-- Bài viết mới nhất -->
          <LatestList :article-list="data.newest_articles" />
        </div>
      </div>
    </div>
  </main>
  <!-- Chân trang -->
  <footer>
    <AppFooter />
  </footer>
</template>
