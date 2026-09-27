<script setup>
import { NButton, NForm, NFormItem, NInput, NTabPane, NTabs } from 'naive-ui'
import { onMounted, ref } from 'vue'

import api from '@/api'
import UploadOne from '@/components//UploadOne.vue'
import CommonPage from '@/components/common/CommonPage.vue'
import { useUserStore } from '@/store'

const userStore = useUserStore()

const infoFormRef = ref(null)
const infoForm = ref({
  // 用原始 avatar 而不是 userStore.avatar 这个 getter:
  // getter 跑过 convertImgUrl, 提交时会把展示用的地址(空头像时是占位图)写回库
  avatar: userStore.userInfo.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
})

onMounted(async () => {
  await userStore.getUserInfo()
  infoForm.value = {
    avatar: userStore.userInfo.avatar,
    nickname: userStore.nickname,
    intro: userStore.intro,
    website: userStore.website,
  }
})

// 两个提交按钮都要带 loading 并在请求期间禁用:
// 改密码接口慢时连点两次, 第二次会带着已经失效的旧密码请求
const infoLoading = ref(false)
const passwordLoading = ref(false)

async function updateProfile() {
  infoFormRef.value?.validate(async (err) => {
    if (err) {
      return
    }
    // 以前是裸 await: 更新失败会留下未捕获的 rejection
    infoLoading.value = true
    try {
      await api.updateCurrent(infoForm.value)
      $message.success('Update successful!')
      await userStore.getUserInfo()
    }
    catch (err) {
      console.error(err)
    }
    finally {
      infoLoading.value = false
    }
  })
}
const infoFormRules = {
  nickname: [
    {
      required: true,
      message: 'Please enter nickname',
      trigger: ['input', 'blur', 'change'],
    },
  ],
}

// Password change form
const passwordFormRef = ref(null)
const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

function updatePassword() {
  passwordFormRef.value?.validate(async (err) => {
    if (err) {
      return
    }
    passwordLoading.value = true
    try {
      await api.updateCurrentPassword(passwordForm.value)
      $message.success('Password updated successfully!')
      // 改成功后清空, 否则旧密码留在框里, 再点一次会用已经失效的旧密码请求
      passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
    }
    catch (err) {
      console.error(err)
    }
    finally {
      passwordLoading.value = false
    }
  })
}
const passwordFormRules = {
  old_password: [
    {
      required: true,
      message: 'Please enter old password',
      trigger: ['input', 'blur', 'change'],
    },
  ],
  new_password: [
    {
      required: true,
      message: 'Please enter new password',
      trigger: ['input', 'blur', 'change'],
    },
  ],
  confirm_password: [
    {
      required: true,
      message: 'Please enter password again',
      trigger: ['input', 'blur'],
    },
    {
      validator: validatePasswordStartWith,
      message: 'Passwords do not match',
      trigger: 'input',
    },
    {
      validator: validatePasswordSame,
      message: 'Passwords do not match',
      trigger: ['blur', 'password-input'],
    },
  ],
}
function validatePasswordStartWith(rule, value) {
  return !!passwordForm.value.new_password && passwordForm.value.new_password.startsWith(value) && passwordForm.value.new_password.length >= value.length
}
function validatePasswordSame(rule, value) {
  return value === passwordForm.value.new_password
}
</script>

<template>
  <CommonPage :show-header="false">
    <NTabs type="line" animated>
      <NTabPane name="website" tab="Edit Information">
        <div class="m-7 flex items-center">
          <div class="mr-7 w-50">
            <UploadOne
              v-model:preview="infoForm.avatar"
              :width="130"
            />
          </div>
          <NForm
            ref="infoFormRef"
            label-placement="left"
            label-align="left"
            label-width="100"
            :model="infoForm"
            :rules="infoFormRules"
            class="w-80"
          >
            <NFormItem label="Nickname" path="nickname">
              <NInput
                v-model:value="infoForm.nickname"
                type="text"
                placeholder="Please enter nickname"
              />
            </NFormItem>
            <NFormItem label="Biography" path="intro">
              <NInput
                v-model:value="infoForm.intro"
                type="text"
                placeholder="Please enter biography"
              />
            </NFormItem>
            <NFormItem label="Website" path="website">
              <NInput
                v-model:value="infoForm.website"
                type="text"
                placeholder="Please enter website"
              />
            </NFormItem>
            <NButton type="primary" :loading="infoLoading" :disabled="infoLoading" @click="updateProfile">
              Update
            </NButton>
          </NForm>
        </div>
      </NTabPane>
      <NTabPane name="contact" tab="Change Password">
        <NForm
          ref="passwordFormRef"
          label-placement="left"
          label-align="left"
          :model="passwordForm"
          label-width="100"
          :rules="passwordFormRules"
          class="m-[30px] w-[400px]"
        >
          <NFormItem label="Old Password" path="old_password">
            <NInput
              v-model:value="passwordForm.old_password"
              type="password"
              show-password-on="mousedown"
              placeholder="Please enter old password"
            />
          </NFormItem>
          <NFormItem label="New Password" path="new_password">
            <NInput
              v-model:value="passwordForm.new_password"
              :disabled="!passwordForm.old_password"
              type="password"
              show-password-on="mousedown"
              placeholder="Please enter new password"
            />
          </NFormItem>
          <NFormItem label="Confirm Password" path="confirm_password">
            <NInput
              v-model:value="passwordForm.confirm_password"
              :disabled="!passwordForm.new_password"
              type="password"
              show-password-on="mousedown"
              placeholder="Please enter new password again"
            />
          </NFormItem>
          <NButton type="primary" :loading="passwordLoading" :disabled="passwordLoading" @click="updatePassword">
            Update
          </NButton>
        </NForm>
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>
