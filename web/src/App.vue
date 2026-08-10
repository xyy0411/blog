/* ... existing code ... */
import { ref } from 'vue';
      <!-- 渲染当前路由对应的组件 -->
      <router-view></router-view>
<template>
  <div id="app">
    <el-config-provider :locale="zhCn" :message="config">
      <MainLayout v-if="showLayout" />
      <router-view v-else />
    </el-config-provider>
  </div>
</template>

<script lang="ts" setup>
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import MainLayout from '@/layouts/MainLayout.vue';

const config = ref({
  max: 1,
});

const route = useRoute();
const showLayout = computed(() => !route.meta.hideLayout);
</script>

<style>
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.2s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}

/* ============== 全局深色滚动条 ============== */
html {
  scrollbar-width: thin;
  scrollbar-color: #4a4a5e #131320;
}

body {
  scrollbar-width: thin;
  scrollbar-color: #4a4a5e #131320;
}

* {
  scrollbar-width: thin;
  scrollbar-color: #4a4a5e #131320;
}

*::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

*::-webkit-scrollbar-track {
  background: #131320;
  border-radius: 4px;
}

*::-webkit-scrollbar-thumb {
  background: #4a4a5e;
  border-radius: 4px;
  border: 2px solid #131320;
  transition: background 0.2s ease;
}

*::-webkit-scrollbar-thumb:hover {
  background: #5a5a6e;
}

*::-webkit-scrollbar-thumb:active {
  background: #60a5fa;
}

/* ============== Element Plus 组件深色滚动条 ============== */
.el-table__body-wrapper,
.el-scrollbar__wrap,
.el-table__body,
.el-table__header-wrapper {
  scrollbar-width: thin;
  scrollbar-color: #4a4a5e #131320;
}

.el-table__body-wrapper::-webkit-scrollbar,
.el-scrollbar__wrap::-webkit-scrollbar,
.el-table__body::-webkit-scrollbar,
.el-table__header-wrapper::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.el-table__body-wrapper::-webkit-scrollbar-track,
.el-scrollbar__wrap::-webkit-scrollbar-track,
.el-table__body::-webkit-scrollbar-track {
  background: #131320;
}

.el-table__body-wrapper::-webkit-scrollbar-thumb,
.el-scrollbar__wrap::-webkit-scrollbar-thumb,
.el-table__body::-webkit-scrollbar-thumb {
  background: #4a4a5e;
  border-radius: 4px;
  border: 2px solid #131320;
}

.el-table__body-wrapper::-webkit-scrollbar-thumb:hover,
.el-scrollbar__wrap::-webkit-scrollbar-thumb:hover,
.el-table__body::-webkit-scrollbar-thumb:hover {
  background: #5a5a6e;
}

.el-scrollbar__wrap--hidden-default::-webkit-scrollbar {
  display: none;
}
</style>
