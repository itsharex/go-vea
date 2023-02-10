<script lang="ts">
export default { name: 'Dashboard' }
</script>

<script setup lang="ts"></script>

<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>项目介绍</span>
        </div>
      </template>
      <p>前端技术栈: vite+vue3+ts+pinia+element-plus</p>
      <p>后端技术栈: Go+Gin+Gorm</p>
    </el-card>
  </div>
</template>

<style lang="scss" scoped></style>
