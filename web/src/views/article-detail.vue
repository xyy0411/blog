<template>
  <div class="article-detail-page">
    <el-skeleton v-if="loading" :rows="10" animated />

    <template v-else-if="article">
      <div class="detail-toolbar">
        <el-button size="small" plain @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          <span>返回列表</span>
        </el-button>
      </div>

      <article class="article">
        <h1 class="article-title">{{ article.title }}</h1>
        <div class="article-meta">
          <span>浏览 {{ article.look_count }}</span>
          <span>评论 {{ article.comment_count }}</span>
          <span>点赞 {{ article.likes }}</span>
          <span>{{ formatTime(article.created_at) }}</span>
        </div>
        <div class="article-content" v-html="renderedContent" />
      </article>

      <section v-if="article.open_comment" class="comments-section">
        <h2>评论（{{ comments.length }}）</h2>

        <el-card v-if="store.token" shadow="never" class="comment-form">
          <el-input
            v-model="newComment"
            type="textarea"
            :rows="3"
            placeholder="说点什么..."
            maxlength="500"
            show-word-limit
          />
          <div class="comment-form-actions">
            <el-button
              type="primary"
              :loading="submitting"
              @click="submitComment"
            >发表评论</el-button>
          </div>
        </el-card>
        <el-alert
          v-else
          type="info"
          :closable="false"
          class="login-hint"
        >
          <el-link type="primary" @click="router.push('/login')">登录</el-link> 后可评论
        </el-alert>

        <div class="comment-list">
          <el-card
            v-for="c in comments"
            :key="c.id"
            shadow="never"
            class="comment-item"
          >
            <div class="comment-header">
              <span class="comment-user">用户 #{{ c.user_id }}</span>
              <span class="comment-time">{{ formatTime(c.created_at) }}</span>
            </div>
            <div class="comment-content">{{ c.content }}</div>
          </el-card>
          <el-empty
            v-if="comments.length === 0"
            description="还没有评论"
          />
        </div>
      </section>

      <el-alert
        v-else
        type="warning"
        :closable="false"
        class="comments-closed"
      >
        作者已关闭评论
      </el-alert>
    </template>

    <el-empty v-else description="文章不存在或已删除" />
  </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import base, { apiUrl } from '@/api/api.ts';
import blogStore from '@/store/arlog.ts';
import type { Article, Comment } from '@/models/article.ts';
import { ArrowLeft } from '@element-plus/icons-vue';
import MarkdownIt from 'markdown-it';
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const store = blogStore();

const md = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: true,
  typographer: true,
});

const loading = ref(false);
const submitting = ref(false);
const article = ref<Article | null>(null);
const comments = ref<Comment[]>([]);
const newComment = ref('');

const renderedContent = computed(() => {
  if (!article.value?.content) return '';
  return md.render(article.value.content);
});

const articleId = computed(() => route.params.id as string);

async function loadArticle() {
  loading.value = true;
  try {
    const response = await axios.get(
      apiUrl(base.showArticle.replace(':article_id', articleId.value))
    );
    const data = response.data?.data;
    if (data) {
      article.value = data.article ?? null;
      comments.value = Array.isArray(data.comments) ? data.comments : [];
    }
  } catch (error) {
    console.error('获取文章失败', error);
    article.value = null;
    comments.value = [];
  } finally {
    loading.value = false;
  }
}

async function submitComment() {
  if (!newComment.value.trim()) return;
  submitting.value = true;
  try {
    await axios.post(
      apiUrl(base.publishArticleComment.replace(':article_id', articleId.value)),
      { content: newComment.value },
      { headers: { Authorization: store.token } }
    );
    newComment.value = '';
    await loadArticle();
  } catch (error) {
    console.error('发表评论失败', error);
  } finally {
    submitting.value = false;
  }
}

function goBack() {
  router.push('/articles');
}

function formatTime(value?: string) {
  if (!value) return '-';
  return new Date(value).toLocaleString('zh-CN', { hour12: false });
}

watch(articleId, () => {
  void loadArticle();
});

onMounted(() => {
  void loadArticle();
});
</script>

<style scoped>
.article-detail-page {
  padding: 0;
}

.detail-toolbar {
  margin-bottom: 16px;
}

.article {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
  padding: 24px 32px;
  margin-bottom: 20px;
}

.article-title {
  margin: 0 0 12px 0;
  font-size: 26px;
  color: #e4e4e7;
}

.article-meta {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: #8b8b9e;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #2a2a3a;
}

.article-content {
  color: #d4d4d8;
  line-height: 1.7;
  font-size: 15px;
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3) {
  color: #e4e4e7;
  margin: 24px 0 12px 0;
}

.article-content :deep(p) {
  margin: 12px 0;
}

.article-content :deep(code) {
  background: #131320;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
  color: #f87171;
}

.article-content :deep(pre) {
  background: #131320;
  padding: 12px 16px;
  border-radius: 6px;
  overflow-x: auto;
}

.article-content :deep(pre code) {
  color: #d4d4d8;
  padding: 0;
}

.article-content :deep(blockquote) {
  border-left: 3px solid #60a5fa;
  padding-left: 12px;
  color: #8b8b9e;
  margin: 12px 0;
}

.article-content :deep(a) {
  color: #60a5fa;
}

.comments-section h2 {
  color: #e4e4e7;
  font-size: 20px;
  margin: 0 0 16px 0;
}

.comment-form {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
  margin-bottom: 20px;
}

.comment-form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-item {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 12px;
  color: #8b8b9e;
}

.comment-user {
  color: #60a5fa;
}

.comment-content {
  color: #d4d4d8;
  white-space: pre-wrap;
  word-break: break-word;
}

.login-hint,
.comments-closed {
  margin-bottom: 20px;
}
</style>
