<script setup>
import { onMounted, ref } from 'vue'
import { NButton, NDatePicker, NForm, NFormItem, NInput, NRadio, NRadioGroup, NTabPane, NTabs } from 'naive-ui'

import CommonPage from '@/components/common/CommonPage.vue'
import UploadOne from '@/components//UploadOne.vue'

import api from '@/api'

defineOptions({ name: 'Website Management' })

const formRef = ref(null)
const form = ref({
  website_avatar: '',
  website_name: 'Personal Blog',
  website_author: 'Admin',
  website_intro: 'Let the past go with the wind',
  website_notice: 'Blog backend based on Gin, GORM development\nBlog frontend based on Vue3, TS, NaiveUI development\nUnder active development... Keep going!',
  website_createtime: '2023-12-27 22:40:22',
  website_record: 'ICP Record Number',
  qq: '123456789',
  github: 'https://github.com/szluyu99',
  gitee: 'https://gitee.com/szluyu99',
  tourist_avatar: 'https://cdn.hahacode.cn/16815451239215dc82548dcadcd578a5bbc8d5deaa.jpg',
  user_avatar: 'https://cdn.hahacode.cn/2299fc4d14c94e6183b082973b35855d.png',
  article_cover: 'https://cdn.hahacode.cn/1679461519cc592408198d67faf1290ff8969dc614.png',
  is_comment_review: 1,
  is_message_review: 1,
  // is_email_notice: 0,
  // social_login_list: [],
  // social_url_list: [],
  // is_reward: 0,
  // wechat_qrcode: 'http://dummyimage.com/100x100',
  // alipay_ode: 'http://dummyimage.com/100x100',
})

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  const resp = await api.getConfig()
  form.value = resp.data
}

function handleSave() {
  formRef.value?.validate(async (err) => {
    if (!err) {
      try {
        $loadingBar?.start()
        await api.updateConfig(form.value)
        $loadingBar?.finish()
        $message.success('Website information updated successfully')
        // fetchData()
      }
      catch (err) {
        $loadingBar?.error()
      }
    }
  })
}
</script>

<template>
  <CommonPage :show-header="false" show-footer>
    <NTabs type="line" animated>
      <NTabPane name="website" tab="Website">
        <NForm
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4 w-[500px]"
        >
          <NFormItem label="Website Avatar" path="website_avatar">
            <UploadOne
              v-model:preview="form.website_avatar"
              :width="120"
            />
          </NFormItem>
          <NFormItem label="Website Name" path="website_name">
            <NInput v-model:value="form.website_name" placeholder="Please enter website name" />
          </NFormItem>
          <NFormItem label="Website Author" path="website_author">
            <NInput v-model:value="form.website_author" placeholder="Please enter website author" />
          </NFormItem>
          <NFormItem label="Website Description" path="website_intro">
            <NInput v-model:value="form.website_intro" placeholder="Please enter website description" />
          </NFormItem>
          <NFormItem label="Website Creation Date" path="website_createtime">
            <NDatePicker
              v-model:formatted-value="form.website_createtime"
              value-format="yyyy-MM-dd HH:mm:ss"
              type="datetime"
            />
          </NFormItem>
          <NFormItem label="Website Notice" path="website_notice">
            <NInput
              v-model:value="form.website_notice"
              type="textarea"
              placeholder="Please enter website notice"
              :autosize="{ minRows: 4, maxRows: 6 }"
            />
          </NFormItem>
          <NFormItem label="Website Record Number" path="website_record">
            <NInput v-model:value="form.website_record" placeholder="Please enter website record number" />
          </NFormItem>
          <!-- TODO: Third-party login -->
          <!-- <n-form-item label="Third-party Login" path="social_login_list">
            <n-checkbox-group v-model:value="cities">
              <n-space item-style="display: flex;">
                <n-checkbox value="QQ" label="QQ" />
                <n-checkbox value="WeiBo" label="Weibo" />
                <n-checkbox value="WeChat" label="WeChat" />
              </n-space>
            </n-checkbox-group>
          </n-form-item> -->
          <NButton type="primary" @click="handleSave">
            Confirm
          </NButton>
        </NForm>
      </NTabPane>
      <NTabPane name="contact" tab="Social">
        <NForm
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4 w-[500px]"
        >
          <NFormItem label="QQ" path="qq">
            <NInput v-model:value="form.qq" placeholder="Please enter QQ" />
          </NFormItem>
          <NFormItem label="Github" path="github">
            <NInput v-model:value="form.github" placeholder="Please enter Github" />
          </NFormItem>
          <NFormItem label="Gitee" path="gitee">
            <NInput v-model:value="form.gitee" placeholder="Please enter Gitee" />
          </NFormItem>
          <NButton type="primary" @click="handleSave">
            Confirm
          </NButton>
        </NForm>
      </NTabPane>
      <NTabPane name="other" tab="Others">
        <NForm
          ref="formRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4"
        >
          <NForm ref="formRef" label-align="left" :label-width="120" :model="form" inline>
            <NFormItem label="User Avatar" path="user_avatar">
              <UploadOne
                v-model:preview="form.user_avatar"
                :width="120"
              />
            </NFormItem>
            <NFormItem label="Guest Avatar" path="tourist_avatar">
              <UploadOne
                v-model:preview="form.tourist_avatar"
                :width="120"
              />
            </NFormItem>
            <!-- <n-form-item label="WeChat Payment QR" path="tourist_avatar">
              <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
            </n-form-item>
            <n-form-item label="Alipay Payment QR" path="tourist_avatar">
              <n-image border-dashed border-1 text-gray width="120" :src="form.tourist_avatar" />
            </n-form-item> -->
          </NForm>
          <NFormItem label-placement="top" label="Default Article Cover" path="article_cover">
            <UploadOne
              v-model:preview="form.article_cover"
              :width="300"
            />
          </NFormItem>
          <NFormItem label="Comment Default Review" path="is_comment_review">
            <NRadioGroup v-model:value="form.is_comment_review" name="is_comment_review">
              <NRadio value="true">
                Disabled
              </NRadio>
              <NRadio value="false">
                Enabled
              </NRadio>
            </NRadioGroup>
          </NFormItem>
          <NFormItem label="Message Default Review" path="is_message_review">
            <NRadioGroup v-model:value="form.is_message_review" name="is_message_review">
              <NRadio value="true">
                Disabled
              </NRadio>
              <NRadio value="false">
                Enabled
              </NRadio>
            </NRadioGroup>
          </NFormItem>
          <!-- <NFormItem label="Email Notification" path="is_email_notice">
            <NRadioGroup v-model:value="form.is_email_notice" name="is_email_notice">
              <NRadio :value="0">
                Disabled
              </NRadio>
              <NRadio :value="1">
                Enabled
              </NRadio>
            </NRadioGroup>
          </NFormItem> -->
          <NButton type="primary" @click="handleSave">
            Confirm
          </NButton>
        </NForm>
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>
