<template>
  <div class="not-found-page">
    <el-result
      icon="warning"
      title="路径输入有误"
      sub-title="你访问的页面不存在，请检查 URL 是否正确，或返回首页重新进入。"
    >
      <template #extra>
        <div class="actions">
          <el-button type="primary" @click="router.push('/')">返回首页</el-button>
          <el-button @click="router.push('/matching-stats')">前往匹配统计</el-button>
        </div>
        <p class="path-tip">当前路径：{{ route.fullPath }}</p>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'NotFoundPage' });

import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
</script>

<style scoped>
.not-found-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
  background: #f5f7fa;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
}

.path-tip {
  margin-top: 12px;
  color: #606266;
  text-align: center;
  word-break: break-all;
}
</style>
