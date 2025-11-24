<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NForm, NFormItem, NInput, NTabPane, NTabs } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import UploadOne from '@/components//UploadOne.vue'
import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()

const infoFormRef = ref(null)
const infoForm = ref({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
})

onMounted(async () => {
  await userStore.getUserInfo()
  infoForm.value = {
    avatar: userStore.avatar,
    nickname: userStore.nickname,
    intro: userStore.intro,
    website: userStore.website,
  }
})

async function updateProfile() {
  infoFormRef.value?.validate(async (err) => {
    if (!err) {
      await api.updateCurrent(infoForm.value)
      $message.success('Update successful!')
      userStore.getUserInfo()
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
    if (!err) {
      await api.updateCurrentPassword(passwordForm.value)
      $message.success('Password updated successfully!')
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
            <NButton type="primary" @click="updateProfile">
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
          <NButton type="primary" @click="updatePassword">
            Update
          </NButton>
        </NForm>
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>
