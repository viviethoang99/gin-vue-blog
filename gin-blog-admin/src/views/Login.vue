<script setup>
import { useStorage } from '@vueuse/core'
import { NButton, NCheckbox, NInput } from 'naive-ui'
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import api from '@/api'

import AppPage from '@/components/common/AppPage.vue'
import { addDynamicRoutes } from '@/router'
import { useAuthStore, useUserStore } from '@/store'
import { getLocal, removeLocal, setLocal } from '@/utils'

const title = import.meta.env.VITE_TITLE // Read from environment variables

const userStore = useUserStore()
const authStore = useAuthStore()

const router = useRouter()
const { query } = useRoute()

/*
  Mock 模式(GitHub Pages 演示站)预填演示账号, 真实部署留空

  演示站没有后端, 随便什么账号都能登进去, 预填是为了让访客少输两下;
  但连了真后端时预填账号/密码等于把凭据写在页面上, 所以只在 mock 下给。
*/
const isMock = import.meta.env.VITE_USE_MOCK === 'true'

const loginForm = reactive({
  username: isMock ? 'guest' : '',
  password: isMock ? '123456' : '',
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

      // isRemember 是 useStorage 返回的 Ref, 必须取 .value:
      // 直接判断 Ref 恒为真, 取消勾选也会把账号密码存进 localStorage
      isRemember.value ? setLocal('loginInfo', { username, password }) : removeLocal('loginInfo')
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

  // doLogin 里只有 try/finally, 失败时的 rejection 需要在这里兜住
  doLogin(username, password).catch(err => console.error(err))

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
    <div style="transform: translateY(25px)" class="m-auto max-w-[700px] min-w-[345px] flex items-center justify-center rounded-2 bg-white bg-opacity-60 p-4 shadow">
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
          placeholder="用户名"
          :maxlength="20"
        />
        <NInput
          v-model:value="loginForm.password"
          class="h-[50px] items-center pl-2"
          type="password"
          show-password-on="mousedown"
          placeholder="密码"
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
