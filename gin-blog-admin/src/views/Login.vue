<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStorage } from '@vueuse/core'
import { NButton, NCheckbox, NInput } from 'naive-ui'

import AppPage from '@/components/common/AppPage.vue'

import { addDynamicRoutes } from '@/router'
import { getLocal, removeLocal, setLocal } from '@/utils'
import { useAuthStore, useUserStore } from '@/store'
import api from '@/api'

const title = import.meta.env.VITE_TITLE // Read from environment variables

const userStore = useUserStore()
const authStore = useAuthStore()

const router = useRouter()
const { query } = useRoute()

const loginForm = reactive({
  username: 'guest',
  password: '123456',
})

initLoginInfo()

// Get remembered username and password from localStorage
function initLoginInfo() {
  const localLoginInfo = getLocal('loginInfo')
  if (localLoginInfo) {
    loginForm.username = localLoginInfo.username
    loginForm.password = localLoginInfo.password
  }
}

// Reactive LocalStorage/SessionStorage - vueuse
const isRemember = useStorage('isRemember', false)
const loading = ref(false)

async function handleLogin() {
  const { username, password } = loginForm
  if (!username || !password) {
    $message.warning('Please enter username and password')
    return
  }

  const doLogin = async (username, password) => {
    loading.value = true

    // Login API
    try {
      const resp = await api.login({ username, password })
      authStore.setToken(resp.data.token)

      await userStore.getUserInfo()
      await addDynamicRoutes()

      isRemember ? setLocal('loginInfo', { username, password }) : removeLocal('loginInfo')
      $message.success('Login successful')

      // Page navigation: Navigate based on redirect in URL
      if (query.redirect) {
        const path = query.redirect
        Reflect.deleteProperty(query, 'redirect') // Delete property from object
        router.push({ path, query })
      }
      else {
        router.push('/')
      }
    }
    finally {
      loading.value = false
    }
  }

  doLogin(username, password)

  // Check if verification code is needed
  // if (JSON.parse(import.meta.env.VITE_USE_CAPTCHA)) {
  //   // Tencent slide verification code (import js file in index.html)
  //   const captcha = new TencentCaptcha(config.TENCENT_CAPTCHA, async res => res.ret === 0 && doLogin(username, password))
  //   captcha.show()
  // }
  // else {
  // doLogin(username, password)
  // }
}
</script>

<template>
  <!-- FIXME: Using style="background-image: url(/image/login_bg.webp);" doesn't work; set it in a CSS class instead -->
  <AppPage class="backgroundImg bg-cover">
    <div style="transform: translateY(25px)" class="m-auto max-w-[700px] min-w-[345px] flex items-center justify-center rounded-2 bg-white bg-opacity-60 p-4 shadow dark:bg-dark dark:bg-opacity-80">
      <div class="hidden w-[380px] px-5 py-9 md:block">
        <img src="/image/login_banner.webp" class="w-full" alt="login_banner">
      </div>

      <div class="w-[320px] flex flex-col px-4 py-9 space-y-5.5">
        <h5 class="flex items-center justify-center text-2xl text-gray font-normal">
          <img src="/image/logo.svg" alt="logo" class="mr-2 h-[50px] w-[50px]">
          <span> {{ title }} </span>
        </h5>
        <NInput
          v-model:value="loginForm.username"
          class="h-[50px] items-center pl-2"
          autofocus
          placeholder="test@qq.com"
          :maxlength="20"
        />
        <NInput
          v-model:value="loginForm.password"
          class="h-[50px] items-center pl-2"
          type="password"
          show-password-on="mousedown"
          placeholder="11111"
          :maxlength="20"
          @keydown.enter="handleLogin"
        />
        <NCheckbox
          :checked="isRemember"
          label="Remember me"
          :on-update:checked="(val) => (isRemember = val)"
        />
        <NButton
          class="h-[50px] w-full rounded-5"
          type="primary"
          :loading="loading"
          @click="handleLogin"
        >
          Login
        </NButton>
      </div>
    </div>
  </AppPage>
</template>

<style scoped>
.backgroundImg{
  background-image: url(/image/login_bg.webp);
}
</style>
