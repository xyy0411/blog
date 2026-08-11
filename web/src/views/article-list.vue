<template>
  <div class="article-list-page">
    <div class="page-header">
      <div>
        <h1>文章列表</h1>
        <p>所有已发布的文章。点击卡片进入详情。</p>
      </div>
      <el-button v-if="store.token" type="primary" @click="goNew">
        <el-icon><EditPen /></el-icon>
        <span>写新文章</span>
      </el-button>
    </div>

    <el-skeleton v-if="loading" :rows="8" animated />

    <el-empty
      v-else-if="articles.length === 0"
      description="还没有文章，去写第一篇吧"
    />

    <div v-else class="article-cards">
      <article
        v-for="a in articles"
        :key="a.id"
        class="article-card"
        @click="goDetail(a.id!)"
      >
        <div v-if="a.cover" class="card-cover">
          <img :src="a.cover" :alt="a.title" loading="lazy" />
        </div>
        <div class="card-body">
          <h3 class="card-title">{{ a.title }}</h3>
          <p class="card-excerpt">{{ getExcerpt(a) }}</p>
          <div class="card-meta">
            <span class="meta-time">{{ formatTime(a.created_at) }}</span>
            <span class="meta-item">浏览 {{ a.look_count }}</span>
            <span class="meta-item">评论 {{ a.comment_count }}</span>
            <span class="meta-item">点赞 {{ a.likes }}</span>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import base, { apiUrl } from '@/api/api.ts';
import blogStore from '@/store/arlog.ts';
import type { Article } from '@/models/article.ts';
import { EditPen } from '@element-plus/icons-vue';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const store = blogStore();
const router = useRouter();
const loading = ref(false);
const articles = ref<Article[]>([]);

async function loadArticles() {
  loading.value = true;
  try {
    const response = await axios.get(apiUrl(base.articleList));
    // 后端 bug：空列表时 data 字段缺失，前端兜底 []
    const data = response.data?.data;
    articles.value = Array.isArray(data) ? data : [];
  } catch (error) {
    console.error('获取文章列表失败', error);
    articles.value = [];
  } finally {
    loading.value = false;
  }
}

function getExcerpt(a: Article): string {
  const source = a.abstract?.trim() || a.content || '';
  if (!source) return '';
  // 简单去掉 markdown 符号，避免预览里全是 # * ` 这种
  const text = source
    .replace(/^#+\s*/gm, '')
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/!\[.*?\]\(.*?\)/g, '')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/^>\s*/gm, '')
    .replace(/^\s*[-+*]\s+/gm, '')
    .replace(/\n{2,}/g, '\n')
    .trim();
  return text.length > 120 ? text.slice(0, 120) + '…' : text;
}

function goDetail(id: number) {
  router.push(`/article/${id}`);
}

function goNew() {
  router.push('/article/new');
}

function formatTime(value?: string) {
  if (!value) return '-';
  return new Date(value).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

onMounted(() => {
  void loadArticles();
});
</script>

<style scoped>
.article-list-page {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  gap: 16px;
}

.page-header h1 {
  margin: 0 0 6px 0;
  font-size: 22px;
  color: #e4e4e7;
}

.page-header p {
  margin: 0;
  color: #8b8b9e;
  font-size: 13px;
}

.article-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.article-card {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: border-color 0.2s ease, transform 0.2s ease;
  display: flex;
  flex-direction: column;
}

.article-card:hover {
  border-color: #60a5fa;
  transform: translateY(-2px);
}

.card-cover {
  aspect-ratio: 16 / 9;
  overflow: hidden;
  background: #131320;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.card-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #e4e4e7;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-excerpt {
  margin: 0;
  color: #8b8b9e;
  font-size: 13px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.card-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 12px;
  color: #8b8b9e;
  padding-top: 8px;
  border-top: 1px solid #2a2a3a;
}

.meta-time {
  color: #60a5fa;
}

@media (max-width: 640px) {
  .article-cards {
    grid-template-columns: 1fr;
  }
}
</style>
