<script setup>
import { NButton, NDatePicker, NForm, NFormItem, NInput, NRadio, NRadioGroup, NTabPane, NTabs } from 'naive-ui'
import { onMounted, ref } from 'vue'

import api from '@/api'
import UploadOne from '@/components//UploadOne.vue'

import CommonPage from '@/components/common/CommonPage.vue'

defineOptions({ name: 'Website Management' })

// 三个 tab 各自一个 ref: 原来四处(其中一处还是嵌套的)都叫 formRef,
// 后挂载的实例覆盖前面的, handleSave 校验到的永远是最后那个。
// 现在这些表单没有 rules 所以看不出问题, 以后任何一处加必填都会失效
const basicFormRef = ref(null)
const socialFormRef = ref(null)
const otherFormRef = ref(null)
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
  tourist_avatar: 'https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images/config/tourist_avatar.jpeg',
  user_avatar: 'https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images/config/user_avatar.jpeg',
  article_cover: 'https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images/config/default_article_cover.png',
  is_comment_review: 1,
  is_message_review: 1,
  // is_email_notice: 0,
  // social_login_list: [],
  // social_url_list: [],
  // is_reward: 0,
  // wechat_qrcode: 'https://dummyimage.com/100x100',
  // alipay_ode: 'https://dummyimage.com/100x100',
})

onMounted(async () => {
  fetchData()
})

async function fetchData() {
  // 以前是裸 await + 直接赋值: 接口失败会产生未捕获的 rejection,
  // 配置表为空时后端返回 {}, 直接赋值又会把表单里的默认值清空
  try {
    const resp = await api.getConfig()
    if (resp.data && Object.keys(resp.data).length) {
      form.value = resp.data
    }
  }
  catch (err) {
    console.error(err)
  }
}

// 参数是表单实例本身(模板里的 basicFormRef 会自动解包), 不是 ref 对象
function handleSave(formInst) {
  formInst?.validate(async (err) => {
    if (!err) {
      try {
        $loadingBar?.start()
        await api.updateConfig(form.value)
        $loadingBar?.finish()
        $message.success('Website information updated successfully')
        // fetchData()
      }
      catch {
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
          ref="basicFormRef"
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
          <NButton type="primary" @click="handleSave(basicFormRef)">
            Confirm
          </NButton>
        </NForm>
      </NTabPane>
      <NTabPane name="contact" tab="Social">
        <NForm
          ref="socialFormRef"
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
          <NButton type="primary" @click="handleSave(socialFormRef)">
            Confirm
          </NButton>
        </NForm>
      </NTabPane>
      <NTabPane name="other" tab="Others">
        <NForm
          ref="otherFormRef"
          label-placement="left"
          label-align="left"
          :label-width="120"
          :model="form"
          class="mt-4"
        >
          <!-- 这里原来又套了一个表单, 只为让两个上传并排。表单嵌表单没有意义,
               而且 ref 重名; 换成普通的 flex 容器 -->
          <div class="flex flex-wrap gap-8">
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
          </div>
          <NFormItem label-placement="top" label="Default Article Cover" path="article_cover">
            <UploadOne
              v-model:preview="form.article_cover"
              :width="300"
            />
          </NFormItem>
          <!-- 值的含义反着读: is_comment_review = true 表示免审核(直接展示) -->
          <NFormItem label="评论审核" path="is_comment_review">
            <NRadioGroup v-model:value="form.is_comment_review" name="is_comment_review">
              <NRadio value="true">
                关闭(新评论直接展示)
              </NRadio>
              <NRadio value="false">
                开启(需在评论管理里通过)
              </NRadio>
            </NRadioGroup>
          </NFormItem>
          <NFormItem label="留言审核" path="is_message_review">
            <NRadioGroup v-model:value="form.is_message_review" name="is_message_review">
              <NRadio value="true">
                关闭(新留言直接展示)
              </NRadio>
              <NRadio value="false">
                开启(需在留言管理里通过)
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
          <NButton type="primary" @click="handleSave(otherFormRef)">
            Confirm
          </NButton>
        </NForm>
      </NTabPane>
    </NTabs>
  </CommonPage>
</template>
