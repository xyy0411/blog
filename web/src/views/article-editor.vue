<template>
  <div class="article-editor-page">
    <div class="page-header">
      <h1>{{ isEdit ? '编辑文章' : '写新文章' }}</h1>
    </div>

    <el-card shadow="never" class="editor-panel">
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题" required>
          <el-input
            v-model="form.title"
            placeholder="文章标题"
            maxlength="32"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="封面">
          <el-input v-model="form.cover" placeholder="封面图片 URL（可选）" />
        </el-form-item>

        <el-form-item label="内容" required>
          <div class="content-toolbar">
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              accept=".md,.markdown,.txt"
              :on-change="handleFileUpload"
            >
              <el-button size="small">
                <el-icon><Upload /></el-icon>
                <span>上传 Markdown 文件</span>
              </el-button>
            </el-upload>
            <span class="toolbar-hint">支持 markdown 语法，详情页自动渲染</span>
          </div>
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="16"
            placeholder="写点什么... 支持 markdown 语法"
          />
        </el-form-item>

        <el-form-item label="评论">
          <el-switch v-model="form.open_comment" />
          <span class="switch-hint">开启后允许登录用户评论</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">{{ isEdit ? '保存' : '发布' }}</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios';
import base, { apiUrl } from '@/api/api.ts';
import blogStore from '@/store/arlog.ts';
import type { UploadFile } from 'element-plus';
import { Upload } from '@element-plus/icons-vue';
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const store = blogStore();

const isEdit = computed(() => !!route.params.id);
const submitting = ref(false);
const form = ref({
  title: '',
  abstract: '',
  content: '',
  cover: '',
  open_comment: true,
});

function handleFileUpload(file: UploadFile) {
  if (!file.raw) return;
  const reader = new FileReader();
  reader.onload = (e) => {
    const text = String(e.target?.result || '');
    if (text.trim()) {
      form.value.content = text;
    }
  };
  reader.onerror = () => {
    console.error('读取文件失败');
  };
  reader.readAsText(file.raw);
}

async function loadArticle() {
  if (!isEdit.value) return;
  try {
    const response = await axios.get(apiUrl(base.showArticle.replace(':article_id', route.params.id as string)));
    const data = response.data?.data;
    if (data?.article) {
      form.value.title = data.article.title || '';
      form.value.abstract = data.article.abstract || '';
      form.value.content = data.article.content || '';
      form.value.cover = data.article.cover || '';
      form.value.open_comment = data.article.open_comment ?? true;
    }
  } catch (error) {
    console.error('获取文章失败', error);
  }
}

async function submit() {
  if (!form.value.title.trim()) return;
  if (!form.value.content.trim()) return;
  submitting.value = true;
  try {
    const payload = {
      title: form.value.title,
      abstract: form.value.abstract,
      content: form.value.content,
      cover: form.value.cover,
      open_comment: form.value.open_comment,
    };
    if (isEdit.value) {
      await axios.put(
        apiUrl(base.showArticle.replace(':article_id', route.params.id as string)),
        payload,
        { headers: { Authorization: store.token } }
      );
    } else {
      await axios.post(
        apiUrl(base.createArticle),
        payload,
        { headers: { Authorization: store.token } }
      );
    }
    router.push('/articles');
  } catch (error) {
    console.error('提交文章失败', error);
  } finally {
    submitting.value = false;
  }
}

function goBack() {
  router.push('/articles');
}

onMounted(() => {
  void loadArticle();
});
</script>

<style scoped>
.article-editor-page {
  padding: 0;
}

.page-header h1 {
  margin: 0 0 20px 0;
  font-size: 22px;
  color: #e4e4e7;
}

.editor-panel {
  background: #1a1a24;
  border: 1px solid #2a2a3a;
  border-radius: 8px;
}

.content-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #131320;
  border: 1px solid #2a2a3a;
  border-radius: 6px;
}

.toolbar-hint {
  color: #8b8b9e;
  font-size: 12px;
}

.switch-hint {
  margin-left: 8px;
  color: #8b8b9e;
  font-size: 12px;
}
</style>
