<template>
  <div id="app">
    <!-- 使用 ElementPlus 的配置提供器，设置语言为中文，并传入消息配置 -->
    <el-config-provider :locale="zhCn" :message="config">
      <!-- 渲染当前路由对应的组件 -->
      <router-view></router-view>
    </el-config-provider>
  </div>
</template>

<script lang="ts" setup>
// 导入 ElementPlus 的中文语言包
import zhCn from 'element-plus/es/locale/lang/zh-cn';
// 从 Vue 中导入 ref 用于创建响应式数据
import { ref } from 'vue';

// 定义消息配置，最大消息数量为 1
const config = ref({
  max: 1,
});
</script>

<style>
/* ... existing code ... */
</style>

