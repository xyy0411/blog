<template>
  <div class="main-layout">
    <aside class="sidebar" :class="{ 'is-collapsed': collapsed }">
      <div class="sidebar-header">
        <span v-show="!collapsed" class="sidebar-title">博客后台</span>
        <span v-show="collapsed" class="sidebar-title sidebar-title--mini">B</span>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :collapse-transition="false"
        class="sidebar-menu"
        router
      >
        <el-menu-item index="/matching-stats">
          <el-icon><DataLine /></el-icon>
          <template #title>匹配记录统计</template>
        </el-menu-item>
        <el-menu-item index="/matching-application">
          <el-icon><PieChart /></el-icon>
          <template #title>匹配申请统计</template>
        </el-menu-item>
        <el-menu-item index="/about">
          <el-icon><InfoFilled /></el-icon>
          <template #title>关于</template>
        </el-menu-item>
      </el-menu>

      <button
        type="button"
        class="sidebar-toggle"
        :title="collapsed ? '展开导航栏' : '收起导航栏'"
        @click="toggleCollapsed"
      >
        <el-icon>
          <Fold v-if="!collapsed" />
          <Expand v-else />
        </el-icon>
      </button>
    </aside>

    <main class="content">
      <header class="content-header">
        <el-button text @click="toggleCollapsed">
          <el-icon>
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <span class="content-header__title">{{ currentRouteTitle }}</span>
        </el-button>
      </header>

      <section class="content-body">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'MainLayout' });

import {
  DataLine,
  Expand,
  Fold,
  InfoFilled,
  PieChart,
} from '@element-plus/icons-vue';
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const collapsed = ref(false);
const isSmallScreen = ref(false);

const STORAGE_KEY = 'mainLayout.collapsed';
const SMALL_SCREEN_BREAKPOINT = 768;

const activeMenu = computed(() => route.path);

const currentRouteTitle = computed(() => {
  switch (route.path) {
    case '/matching-stats':
      return '匹配记录统计';
    case '/matching-application':
      return '匹配申请统计';
    case '/about':
      return '关于';
    default:
      return '博客后台';
  }
});

function loadPersistedState() {
  if (typeof window === 'undefined') {
    return;
  }
  const saved = window.localStorage.getItem(STORAGE_KEY);
  if (saved === 'true') {
    collapsed.value = true;
  }
}

function persistState() {
  if (typeof window === 'undefined') {
    return;
  }
  window.localStorage.setItem(STORAGE_KEY, String(collapsed.value));
}

function toggleCollapsed() {
  collapsed.value = !collapsed.value;
  persistState();
}

function evaluateScreenSize() {
  if (typeof window === 'undefined') {
    return;
  }
  isSmallScreen.value = window.innerWidth < SMALL_SCREEN_BREAKPOINT;
  if (isSmallScreen.value && !collapsed.value) {
    collapsed.value = true;
    persistState();
  }
}

watch(
  () => route.path,
  () => {
    if (isSmallScreen.value && !collapsed.value) {
      collapsed.value = true;
      persistState();
    }
  },
);

onMounted(() => {
  loadPersistedState();
  evaluateScreenSize();
  window.addEventListener('resize', evaluateScreenSize);
});

onBeforeUnmount(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('resize', evaluateScreenSize);
  }
});
</script>

<style scoped>
.main-layout {
  display: flex;
  min-height: 100vh;
  background: #131320;
}

/* ============== 侧边栏：深色品牌背景（Linear/Stripe 风格） ============== */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 1001;
  display: flex;
  flex-direction: column;
  width: 240px;
  background: linear-gradient(180deg, #1e1e2d 0%, #2a2a3c 100%);
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.12);
  transition: width 0.28s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.is-collapsed {
  width: 64px;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 64px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1px;
}

.sidebar-title--mini {
  font-size: 22px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
  background: transparent;
  padding: 12px 0;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 240px;
}

/* 菜单项：深色主题适配 */
.sidebar-menu :deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
  margin: 4px 12px;
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.65);
  transition: all 0.2s ease;
}

/* 收起状态：移除水平 margin，让图标居中 */
.sidebar-menu.el-menu--collapse :deep(.el-menu-item) {
  margin: 4px 0;
  padding: 0;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(96, 165, 250, 0.2), rgba(129, 140, 248, 0.15));
  color: #60a5fa;
  box-shadow: inset 3px 0 0 #60a5fa;
}

.sidebar-menu :deep(.el-menu-item.is-active .el-icon) {
  color: #60a5fa;
}

.sidebar-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 40px;
  margin: 8px 16px 16px;
  padding: 0;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.2s ease;
}

/* 收起状态：按钮水平居中 */
.sidebar.is-collapsed .sidebar-toggle {
  margin-left: 12px;
  margin-right: 12px;
}

.sidebar-toggle:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border-color: rgba(255, 255, 255, 0.15);
}

/* ============== 内容区 ============== */
.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  margin-left: 240px;
  transition: margin-left 0.28s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar.is-collapsed ~ .content {
  margin-left: 64px;
}

.content-header {
  position: sticky;
  top: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  height: 64px;
  padding: 0 24px;
  background: rgba(30, 30, 45, 0.85);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.content-header :deep(.el-button) {
  font-size: 16px;
  color: #c0c0d0;
}

.content-header :deep(.el-button:hover) {
  color: #60a5fa;
}

.content-header__title {
  margin-left: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #e4e4e7;
}

.content-body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  scrollbar-width: thin;
  scrollbar-color: #4a4a5e #131320;
}

.content-body::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.content-body::-webkit-scrollbar-track {
  background: #131320;
}

.content-body::-webkit-scrollbar-thumb {
  background: #4a4a5e;
  border-radius: 4px;
  border: 2px solid #131320;
}

.content-body::-webkit-scrollbar-thumb:hover {
  background: #5a5a6e;
}

.content-body::-webkit-scrollbar-thumb:active {
  background: #60a5fa;
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity 0.25s ease;
}

.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}

@media (max-width: 768px) {
  .sidebar {
    max-width: 80vw;
  }

  .content-header__title {
    font-size: 14px;
  }
}
</style>
