<script setup lang="ts">
import {reactive, ref} from "vue";
import type {Article} from "@/models/article.ts";
import axios from "axios";

let searchQuery = ref(``)
  let onSearch = () => {
  }
  function btn() {
    alert(`按钮`)
  }
interface ResponseType {
  data: Article; // 用实际类型替换 any
  title: string;
}
let articles = reactive<Article[]>([])
  let url = `http://127.0.0.1:3000/api/article/1`
  try {
    const response = await axios.get<ResponseType>(url)
    // 检查响应状态码
    if (response.status === 200) {
      // 从response.data中解构data和title
      const {data} = response.data;
      articles.push(data)
    }
  }
</script>
<template>
  <h1>欢迎来到我的博客</h1>
  <div class="search-container">
      <el-input
        v-model="searchQuery"
        placeholder="请输入搜索内容"
        clearable
        @input="onSearch"
        suffix-icon="el-icon-search"
      />
      <el-button type="primary" @click="btn">搜索</el-button>
  </div>
  <div class="article-card">
    <div class="article-card" v-for="(article, index) in articles" :key="index">
      <div class="article-title">{{ article.title }}</div>
      <div class="article-content">{{ article.content }}</div>
    </div>
  </div>
</template>

<style scoped>
h1 {
  display: flex;
  align-items: center;
  justify-content: center;
}
.search-container {
  display: flex;
  align-items: center; /* 垂直居中对齐 */
  gap: 10px; /* 元素之间的间距 */
}
.article-card {
  width: 250px; /* 方格的宽度 */
  height: 200px; /* 方格的高度 */
  border: 1px solid #ccc; /* 边框 */
  border-radius: 10px; /* 圆角 */
  padding: 10px; /* 内边距 */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1); /* 阴影效果 */
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.article-title {
  font-weight: bold;
  margin-bottom: 10px;
}
.article-content {
  font-size: 14px;
  color: #555;
}
</style>
