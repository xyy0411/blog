<template>
  <div class="article-list-page">
    <div class="page-header">
      <div>
        <h1>文章列表</h1>
        <p>所有已发布的文章。点击标题进入详情。</p>
      </div>
      <el-button v-if="store.token" type="primary" @click="goNew">
        <el-icon><EditPen /></el-icon>
        <span>写新文章</span>
      </el-button>
    </div>

    <el-card shadow="never" class="list-panel">
      <el-skeleton v-if="loading" :rows="6" animated />
      <el-empty
        v-else-if="articles.length === 0"
        description="还没有文章，去写第一篇吧"
      />
      <el-table v-else :data="articles" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="标题" min-width="240">
          <template #default="{ row }">
            <el-link type="primary" @click="goDetail(row.id)">{{ row.title }}</el-link>
            <div v-if="row.abstract" class="article-abstract">{{ row.abstract }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="look_count" label="浏览" width="100" sortable />
        <el-table-column prop="comment_count" label="评论" width="100" sortable />
        <el-table-column prop="likes" label="点赞" width="100" sortable />
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="goDetail(row.id)">查看</el-button>
            <el-button
              v-if="store.token"
              size="small"
              type="primary"
              @click="goEdit(row.id)"
            >编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
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

function goDetail(id: number) {
  router.push(`/article/${id}`);
}

function goEdit(id: number) {
  router.push(`/article/${id}/edit`);
}

function goNew() {
  router.push('/article/new');
}

function formatTime(value?: string) {
  if (!value) return '-';
  return new Date(value).toLocaleString('zh-CN', { hour12: false });
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

.list-panel {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
}

.article-abstract {
  margin-top: 4px;
  font-size: 12px;
  color: #8b8b9e;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
