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
          <div class="sections">
            <div
              v-for="(_, idx) in sections"
              :key="idx"
              class="section-item"
            >
              <div class="section-toolbar">
                <span class="section-label">段落 {{ idx + 1 }}</span>
                <el-button
                  v-if="sections.length > 1"
                  size="small"
                  type="danger"
                  plain
                  @click="removeSection(idx)"
                >删除此段</el-button>
              </div>
              <el-input
                v-model="sections[idx]"
                type="textarea"
                :rows="6"
                placeholder="支持 markdown 语法"
              />
            </div>
          </div>
          <el-button class="add-section-btn" @click="addSection" plain>
            <el-icon><Plus /></el-icon>
            <span>添加段落</span>
          </el-button>
        </el-form-item>

        <el-form-item label="评论">
          <el-switch v-model="form.open_comment" />
          <span class="switch-hint">开启后允许登录用户评论</span>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="submitting"
            @click="submit"
          >{{ isEdit ? '保存' : '发布' }}</el-button>
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
import { Plus } from '@element-plus/icons-vue';
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const store = blogStore();

const isEdit = computed(() => !!route.params.id);
const submitting = ref(false);
const sections = ref<string[]>(['']);
const form = ref({
  title: '',
  cover: '',
  open_comment: true,
});

function addSection() {
  sections.value.push('');
}

function removeSection(idx: number) {
  sections.value.splice(idx, 1);
}

async function loadArticle() {
  if (!isEdit.value) return;
  try {
    const response = await axios.get(
      apiUrl(base.showArticle.replace(':article_id', route.params.id as string))
    );
    const data = response.data?.data;
    if (data?.article) {
      form.value.title = data.article.title || '';
      form.value.cover = data.article.cover || '';
      form.value.open_comment = data.article.open_comment ?? true;
      // 后端把 Contents 切片拼成单字符串存到 content，编辑器按双换行拆段还原
      const raw = (data.article.content || '') as string;
      const parts = raw.split(/\n{2,}/).filter((s: string) => s.trim());
      sections.value = parts.length > 0 ? parts : [''];
    }
  } catch (error) {
    console.error('获取文章失败', error);
  }
}

async function submit() {
  if (!form.value.title.trim()) return;
  const nonEmpty = sections.value.filter((s) => s.trim());
  if (nonEmpty.length === 0) return;

  submitting.value = true;
  try {
    const payload = {
      title: form.value.title,
      contents: { content: nonEmpty },
      cover: form.value.cover,
      open_comment: form.value.open_comment,
    };
    const headers = { token: store.token };
    if (isEdit.value) {
      await axios.put(
        apiUrl(base.showArticle.replace(':article_id', route.params.id as string)),
        payload,
        { headers }
      );
    } else {
      await axios.post(apiUrl(base.createArticle), payload, { headers });
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

.sections {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
}

.section-item {
  border: 1px solid #2a2a3a;
  border-radius: 6px;
  padding: 12px;
  background: #131320;
}

.section-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.section-label {
  color: #8b8b9e;
  font-size: 12px;
}

.add-section-btn {
  margin-top: 12px;
  width: 100%;
  border-style: dashed;
}

.switch-hint {
  margin-left: 8px;
  color: #8b8b9e;
  font-size: 12px;
}
</style>
