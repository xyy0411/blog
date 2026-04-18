<template>
  <div class="home-page">
    <div class="hero">
      <h1>欢迎来到我的博客</h1>
      <p>搜索文章，或者进入匹配统计页面查看当日与累计数据。</p>
      <div class="hero-actions">
        <el-button type="primary" @click="router.push('/matching-stats')">
          查看匹配统计
        </el-button>
      </div>
    </div>

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

    <div class="articles-grid">
      <div class="article-card" v-for="(article, index) in articles" :key="index">
        <div class="article-title">{{ article.title }}</div>
        <div class="article-content">{{ article.content }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'HomePage' });
import axios from 'axios';
import base, { apiUrl } from '@/api/api.ts';
import type { Article } from '@/models/article.ts';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const searchQuery = ref('');
const articles = ref<Article[]>([]);
const searchUrl = apiUrl('/api/article/search');

const onSearch = async () => {
  try {
    const response = await axios.get<{ data: Article[] }>(searchUrl, {
      params: {
        query: searchQuery.value,
      },
    });

    if (response.status === 200) {
      articles.value = response.data.data;
    }
  } catch (error) {
    console.error('搜索失败', error);
  }
};

const fetchAllArticles = async () => {
  try {
    const allArticlesUrl = apiUrl(base.articleList);
    const response = await axios.get<{ data: Article[] }>(allArticlesUrl);

    if (response.status === 200) {
      articles.value = response.data.data;
    }
  } catch (error) {
    console.error('获取文章列表失败', error);
  }
};

onMounted(() => {
  void fetchAllArticles();
});
</script>

<style scoped>
.home-page {
  padding: 32px;
}

.hero {
  margin-bottom: 24px;
  text-align: center;
}

.hero h1 {
  margin-bottom: 12px;
}

.hero p {
  margin: 0 0 16px;
  color: #606266;
}

.hero-actions {
  display: flex;
  justify-content: center;
}

.search-container {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 24px;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.article-card {
  min-height: 200px;
  border: 1px solid #ccc;
  border-radius: 10px;
  padding: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
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
