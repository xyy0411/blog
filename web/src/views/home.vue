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
    <el-button type="primary" @click="onSearch">搜索</el-button>
  </div>
  <div class="article-card">
    <div class="article-card" v-for="(article, index) in articles" :key="index">
      <div class="article-title">{{ article.title }}</div>
      <div class="article-content">{{ article.content }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import type { Article } from "@/models/article.ts";
import axios from "axios";

// 定义搜索查询的响应式变量
let searchQuery = ref('');
// 定义文章列表的响应式变量
let articles = reactive<Article[]>([]);

// 定义搜索接口的URL，这里假设搜索接口为 /api/article/search
let searchUrl = `http://127.0.0.1:3000/api/article/search`;

// 处理搜索逻辑
const onSearch = async () => {
  try {
    // 发起搜索请求，将搜索查询作为参数传递
    const response = await axios.get<{ data: Article[] }>(searchUrl, {
      params: {
        query: searchQuery.value
      }
    });
    // 检查响应状态码
    if (response.status === 200) {
      // 更新文章列表
      articles = response.data.data;
    }
  } catch (error) {
    console.error('搜索失败', error);
  }
};

// 页面加载时获取所有文章
const fetchAllArticles = async () => {
  try {
    // 假设获取所有文章的接口为 /api/article/all
    const allArticlesUrl = `http://127.0.0.1:3000/api/article/all`;
    const response = await axios.get<{ data: Article[] }>(allArticlesUrl);
    if (response.status === 200) {
      articles = response.data.data;
    }
  } catch (error) {
    console.error('获取文章列表失败', error);
  }
};

// 页面加载时调用获取所有文章的函数
fetchAllArticles();

</script>

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
