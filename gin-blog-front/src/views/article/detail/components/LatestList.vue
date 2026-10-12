<script setup>
import dayjs from 'dayjs'
import { convertImgUrl } from '@/utils'

const { articleList } = defineProps({
  articleList: Array,
})
</script>

<template>
  <Transition name="slide-fade" appear>
    <div class="card-view space-y-3">
      <div class="flex items-center">
        <span class="i-akar-icons:clock" />
        <span class="ml-2"> Bài viết mới nhất </span>
      </div>
      <ul class="space-y-1">
        <!-- Đường phân cách phải đặt trên li: trước đây được đặt trên RouterLink, là phần tử con duy nhất của li,
             khiến last:border-none luôn khớp và đường nét đứt không bao giờ hiển thị (border thực tế là none/0px) -->
        <li
          v-for="item of articleList" :key="item.id"
          class="border-b-1 border-color-divider border-dashed p-1 last:border-none hover:bg-surface-soft"
        >
          <RouterLink :to="`/article/${item.id}`">
            <div class="flex items-center">
              <!-- object-cover: trước đây dùng fill mặc định khiến ảnh bìa không vuông bị kéo méo -->
              <img
                :src="convertImgUrl(item.img)"
                class="h-15 w-15 shrink-0 rounded bg-surface-soft object-cover"
                loading="lazy" :alt="item.title"
              >
              <div class="min-w-0 flex-1 pl-2">
                <p class="line-clamp-2">
                  {{ item.title }}
                </p>
                <p class="text-sm text-muted">
                  {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
                </p>
              </div>
            </div>
          </RouterLink>
        </li>
      </ul>
    </div>
  </Transition>
</template>
