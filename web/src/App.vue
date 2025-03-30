<template>
  <div id="app">
    <el-config-provider :locale="zhCn" :message="config">
      <router-view></router-view>
    </el-config-provider>
  </div>
</template>

<script lang="ts" setup>
import {defineComponent, ref} from 'vue';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import {useRoute, useRouter} from "vue-router";
const router = useRouter()
const route = useRoute()

const config = ref({
  max:1,
})

</script>

<style>
/* ... existing code ... */
</style>
