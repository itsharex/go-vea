<template>
  <el-upload drag :action="uploadFileUrl" :headers="headers">
    <el-icon class="el-icon--upload"><upload-filled /></el-icon>
    <div class="el-upload__text">Drop file here or <em>click to upload</em></div>
    <template #tip>
      <div class="el-upload__tip">jpg/png files with a size less than 500kb</div>
    </template>
  </el-upload>
</template>

<script setup lang="ts">
import { getToken } from '@/utils/auth'
import { UploadFilled } from '@element-plus/icons-vue'

const uploadFileUrl = ref(import.meta.env.VITE_APP_BASE_API + '/system/file/upload') // 上传文件服务器地址
const data = reactive({
  headers: {
    Authorization: 'Bearer ' + getToken()
  }
})

const { headers } = toRefs(data)
</script>
