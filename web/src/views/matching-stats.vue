<template>
  <div class="stats-page">
    <div class="page-header">
      <div>
        <h1>
          匹配记录统计
          <span class="current-view-badge">{{ currentLabel }}</span>
        </h1>
        <p>展示匹配记录的成功统计，支持当日、本周、累计及指定日期查询。</p>
      </div>
    </div>

    <el-card shadow="never" class="stats-panel">
      <div class="button-group">
        <el-button
          :type="activeView === 'today' ? 'primary' : 'default'"
          @click="switchView('today')"
        >
          当日
        </el-button>
        <el-button
          :type="activeView === 'week' ? 'primary' : 'default'"
          @click="switchView('week')"
        >
          本周
        </el-button>
        <el-button
          :type="activeView === 'all' ? 'primary' : 'default'"
          @click="switchView('all')"
        >
          累计
        </el-button>
        <el-button
          :type="activeView === 'date' ? 'primary' : 'default'"
          @click="switchView('date')"
        >
          按日期
        </el-button>
        <el-button
          :type="activeView === 'range' ? 'primary' : 'default'"
          @click="switchView('range')"
        >
          按日期范围
        </el-button>
      </div>

      <div v-if="activeView === 'date'" class="range-form">
        <el-date-picker
          v-model="singleDate"
          type="date"
          placeholder="选择日期"
          value-format="YYYY-MM-DD"
          :disabled="loading"
        />
        <el-button
          type="primary"
          :loading="loading"
          :disabled="!singleDate"
          @click="loadDateStats"
        >
          查询
        </el-button>
      </div>

      <div v-if="activeView === 'range'" class="range-form">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="截至日期"
          value-format="YYYY-MM-DD"
          :disabled="loading"
        />
        <el-button
          type="primary"
          :loading="loading"
          :disabled="!dateRange || dateRange.length !== 2"
          @click="loadRangeStats"
        >
          查询
        </el-button>
      </div>

      <el-skeleton v-if="loading" :rows="8" animated />

      <template v-else>
        <el-alert
          v-if="errorMessage"
          :title="errorMessage"
          type="error"
          show-icon
          :closable="false"
          class="status-message"
        />

        <template v-else>
          <div :key="animationKey" class="content-enter">
          <!-- 第一档：成功匹配记录数（Hero 卡片，横跨整行） -->
          <div class="summary-row summary-row--primary">
            <el-card shadow="hover" class="summary-card summary-card--hero enter-hero">
              <span class="summary-label">成功匹配记录数</span>
              <strong>{{ stats.total }}</strong>
              <span class="summary-sub">所有匹配记录总量</span>
            </el-card>
          </div>

          <!-- 第二档：去重用户数 + 人均匹配次数 -->
          <div class="summary-row summary-row--secondary">
            <el-card shadow="hover" class="summary-card summary-card--users enter-card" style="--enter-delay: 0.15s">
              <span class="summary-label">去重用户数</span>
              <strong>{{ uniqueUserCount }}</strong>
              <span class="summary-sub">参与匹配的独立用户总数</span>
            </el-card>
            <el-card shadow="hover" class="summary-card summary-card--avg enter-card" style="--enter-delay: 0.3s">
              <span class="summary-label">人均匹配次数</span>
              <strong>{{ avgMatchesPerUser }}</strong>
              <span class="summary-sub">每位用户平均参与匹配次数</span>
            </el-card>
          </div>

          <el-empty
            v-if="stats.records.length === 0"
            description="当前没有匹配记录"
          />

          <template v-else>
            <el-card shadow="never" class="chart-card enter-card" style="--enter-delay: 0.45s">
              <template #header>
                <div class="section-header">
                  <div>
                    <h2>匹配成功统计图</h2>
                    <p>X 轴为 name[qq号]，Y 轴为该用户成功匹配的次数；悬停柱子可快速预览，点击柱子可在图表外固定查看详情。</p>
                  </div>
                  <el-button text :loading="loading" @click="refreshCurrent">
                    刷新数据
                  </el-button>
                </div>
              </template>

              <transition name="detail-fade">
                <div v-if="selectedBarDetail" class="selected-bar-panel">
                  <div class="selected-bar-panel__header">
                    <div>
                      <span class="detail-caption">已固定的用户详情</span>
                      <strong>{{ selectedBarDetail.name }}</strong>
                    </div>
                    <span class="selected-bar-panel__hint">再次点击同一柱子可取消固定</span>
                  </div>
                  <div class="selected-bar-panel__meta">
                    <span>QQ号：{{ selectedBarDetail.id || '-' }}</span>
                    <span>成功匹配次数：{{ selectedBarDetail.count }}</span>
                    <span>上次匹配成功时间：{{ formatTime(selectedBarDetail.lastMatchedAt) }}</span>
                  </div>
                </div>
              </transition>

              <div class="chart-scroll">
                <div class="chart-wrapper" :style="chartWrapperStyle">
                  <div class="chart-y-axis">
                    <span v-for="tick in yAxisTicks" :key="tick">{{ tick }}</span>
                  </div>

                  <div class="chart-area">
                    <div class="chart-grid">
                      <span
                        v-for="tick in yAxisTicks"
                        :key="`grid-${tick}`"
                        class="grid-line"
                      />
                    </div>

                    <div class="bars-row">
                      <div
                        v-for="(item, index) in chartData"
                        :key="item.key"
                        class="bar-column"
                      >
                        <div class="bar-track">
                          <div class="bar-interactive">
                            <transition name="detail-fade">
                              <div
                                v-if="hoveredBarKey === item.key"
                                class="bar-popover"
                              >
                                <span class="detail-caption">悬浮查看</span>
                                <strong>{{ item.name }}</strong>
                                <span>QQ号：{{ item.id || '-' }}</span>
                                <span>成功匹配次数：{{ item.count }}</span>
                                <span>上次匹配成功时间：{{ formatTime(item.lastMatchedAt) }}</span>
                              </div>
                            </transition>
                            <span class="bar-value enter-bar-value" :style="{ animationDelay: `${0.6 + index * 0.08}s` }">{{ item.count }}</span>
                            <button
                              type="button"
                              class="bar enter-bar"
                              :class="{ 'is-active': activeBarKey === item.key }"
                              :style="{ height: `${item.heightPx}px`, animationDelay: `${0.6 + index * 0.08}s` }"
                              :title="`${item.label}：${item.count}`"
                              @mouseenter="setHoveredBar(item.key)"
                              @mouseleave="clearHoveredBar(item.key)"
                              @focus="setHoveredBar(item.key)"
                              @blur="clearHoveredBar(item.key)"
                              @click="toggleSelectedBar(item.key)"
                            />
                          </div>
                        </div>
                        <span class="bar-label enter-bar-label" :style="{ animationDelay: `${0.6 + index * 0.08}s` }">{{ item.label }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>

            <el-card shadow="never" class="detail-card collapsible-card enter-card" style="--enter-delay: 0.6s">
              <template #header>
                <div
                  class="section-header collapsible-header"
                  @click="summaryCollapsed = !summaryCollapsed"
                >
                  <div>
                    <h2>统计汇总</h2>
                    <p>按用户聚合后的成功匹配次数，便于快速核对图表数据。</p>
                  </div>
                  <el-icon
                    class="collapsible-arrow"
                    :class="{ 'is-open': !summaryCollapsed }"
                  >
                    <ArrowDown />
                  </el-icon>
                </div>
              </template>

              <transition name="card-collapse">
                <div v-show="!summaryCollapsed">
                  <el-table :data="chartData" stripe>
                    <el-table-column prop="label" label="用户(name[qq号])" min-width="220" />
                    <el-table-column prop="name" label="用户名" min-width="140" />
                    <el-table-column prop="id" label="QQ号" min-width="120" />
                    <el-table-column prop="count" label="成功次数" min-width="120" sortable />
                  </el-table>
                </div>
              </transition>
            </el-card>

            <el-card shadow="never" class="detail-card collapsible-card enter-card" style="--enter-delay: 0.75s">
              <template #header>
                <div
                  class="section-header collapsible-header"
                  @click="rawDataCollapsed = !rawDataCollapsed"
                >
                  <div>
                    <h2>原始数据</h2>
                    <p>图表下方保留完整匹配记录，方便继续排查和对照。</p>
                  </div>
                  <el-icon
                    class="collapsible-arrow"
                    :class="{ 'is-open': !rawDataCollapsed }"
                  >
                    <ArrowDown />
                  </el-icon>
                </div>
              </template>

              <transition name="card-collapse">
                <div v-show="!rawDataCollapsed">
                  <el-table :data="stats.records" stripe>
                    <el-table-column prop="match_id" label="匹配 ID" min-width="180" />
                    <el-table-column label="发起用户" min-width="220">
                      <template #default="scope">
                        {{ formatUserLabel(scope.row.user_name, scope.row.user_id) }}
                      </template>
                    </el-table-column>
                    <el-table-column label="匹配对象" min-width="220">
                      <template #default="scope">
                        {{ formatUserLabel(scope.row.peer_name, scope.row.peer_id) }}
                      </template>
                    </el-table-column>
                    <el-table-column label="创建时间" min-width="180">
                      <template #default="scope">
                        {{ formatTime(scope.row.created_at) }}
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </transition>
            </el-card>
          </template>
          </div>
        </template>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'MatchingStatsPage' });

import axios from 'axios';
import { ArrowDown } from '@element-plus/icons-vue';
import base, { apiUrl } from '@/api/api.ts';
import { computed, onMounted, reactive, ref } from 'vue';

interface MatchingRecord {
  id: number;
  user_id: number;
  user_name: string;
  peer_id: number;
  peer_name: string;
  match_id: string;
  created_at: string;
}

interface MatchingStatsResponse {
  data?: {
    total?: number;
    records?: MatchingRecord[];
  };
}

interface UserChartItem {
  key: string;
  id: number;
  name: string;
  label: string;
  count: number;
  heightPx: number;
  lastMatchedAt: string;
}

type ViewMode = 'today' | 'week' | 'all' | 'date' | 'range';

const loading = ref(false);
const errorMessage = ref('');
const activeView = ref<ViewMode>('today');
const singleDate = ref<string>('');
const dateRange = ref<[string, string] | null>(null);
const hoveredBarKey = ref('');
const summaryCollapsed = ref(true);
const rawDataCollapsed = ref(true);
const selectedBarKey = ref('');
const animationKey = ref(0);
const maxBarHeight = 240;
const stats = reactive({
  total: 0,
  records: [] as MatchingRecord[],
});

const currentLabel = computed(() => {
  if (activeView.value === 'today') return '当日匹配统计';
  if (activeView.value === 'week') return '本周匹配统计';
  if (activeView.value === 'all') return '累计匹配统计';
  if (activeView.value === 'date' && singleDate.value) {
    return `${singleDate.value} 匹配统计`;
  }
  if (dateRange.value && dateRange.value.length === 2) {
    return `${dateRange.value[0]} 至 ${dateRange.value[1]}`;
  }
  return '指定日期范围';
});

const uniqueUserCount = computed(() => {
  const userSet = new Set<string>();
  stats.records.forEach((record) => {
    if (record.user_id) {
      userSet.add(`${record.user_id}`);
    }
    if (record.peer_id) {
      userSet.add(`${record.peer_id}`);
    }
  });
  return userSet.size;
});

const avgMatchesPerUser = computed(() => {
  if (uniqueUserCount.value === 0) {
    return '0.0';
  }
  return (stats.total / uniqueUserCount.value).toFixed(1);
});

const chartData = computed<UserChartItem[]>(() => {
  const counter = new Map<string, Omit<UserChartItem, 'count' | 'heightPx' | 'lastMatchedAt'>>();
  const countMap = new Map<string, number>();
  const lastMatchedMap = new Map<string, string>();

  const upsertUser = (id: number, name: string, createdAt: string) => {
    if (!id && !name) {
      return;
    }

    const normalizedName = name?.trim() || '未知用户';
    const key = `${normalizedName}-${id}`;

    if (!counter.has(key)) {
      counter.set(key, {
        key,
        id,
        name: normalizedName,
        label: formatUserLabel(normalizedName, id),
      });
    }

    countMap.set(key, (countMap.get(key) ?? 0) + 1);

    const currentLastMatched = lastMatchedMap.get(key);
    if (!currentLastMatched || new Date(createdAt).getTime() > new Date(currentLastMatched).getTime()) {
      lastMatchedMap.set(key, createdAt);
    }
  };

  stats.records.forEach((record) => {
    upsertUser(record.user_id, record.user_name, record.created_at);
    upsertUser(record.peer_id, record.peer_name, record.created_at);
  });

  const maxCount = Math.max(...countMap.values(), 0);

  return Array.from(countMap.entries())
    .map(([key, count]) => {
      const user = counter.get(key)!;
      return {
        ...user,
        count,
        heightPx: maxCount > 0 ? Math.max((count / maxCount) * maxBarHeight, 12) : 0,
        lastMatchedAt: lastMatchedMap.get(key) ?? '',
      };
    })
    .sort((left, right) => right.count - left.count || left.id - right.id);
});

const yAxisTicks = computed(() => {
  const maxCount = Math.max(...chartData.value.map((item) => item.count), 0);
  const safeMax = Math.max(maxCount, 1);
  const step = Math.max(1, Math.ceil(safeMax / 4));
  const ticks: number[] = [];

  for (let tick = step * 4; tick >= 0; tick -= step) {
    ticks.push(tick);
  }

  if (!ticks.includes(safeMax)) {
    ticks.splice(1, 0, safeMax);
  }

  return Array.from(new Set(ticks)).sort((left, right) => right - left);
});

const chartWrapperStyle = computed(() => {
  const minWidth = Math.max(chartData.value.length * 92, 720);
  return {
    minWidth: `${minWidth}px`,
  };
});

const activeBarKey = computed(() => hoveredBarKey.value || selectedBarKey.value);

const selectedBarDetail = computed(() =>
  chartData.value.find((item) => item.key === selectedBarKey.value) ?? null,
);

const endpointMap: Record<Exclude<ViewMode, 'date' | 'range'>, string> = {
  today: apiUrl(base.matchingToday),
  week: apiUrl(base.matchingWeek),
  all: apiUrl(base.matchingAll),
};

function resetStats() {
  stats.total = 0;
  stats.records = [];
  hoveredBarKey.value = '';
  selectedBarKey.value = '';
}

async function loadPreset(view: 'today' | 'week' | 'all') {
  loading.value = true;
  errorMessage.value = '';

  try {
    const response = await axios.get<MatchingStatsResponse>(endpointMap[view]);
    stats.total = response.data.data?.total ?? 0;
    stats.records = response.data.data?.records ?? [];
    hoveredBarKey.value = '';
    selectedBarKey.value = '';
    animationKey.value++;
  } catch (error) {
    console.error('获取匹配统计失败', error);
    resetStats();
    errorMessage.value = '获取匹配统计失败，请稍后重试。';
  } finally {
    loading.value = false;
  }
}

async function loadDateStats() {
  if (!singleDate.value) {
    errorMessage.value = '请选择日期';
    return;
  }

  loading.value = true;
  errorMessage.value = '';

  try {
    const response = await axios.get<MatchingStatsResponse>(
      apiUrl(base.matchingRecordRange),
      {
        params: { start_date: singleDate.value, end_date: singleDate.value },
      },
    );
    stats.total = response.data.data?.total ?? 0;
    stats.records = response.data.data?.records ?? [];
    hoveredBarKey.value = '';
    selectedBarKey.value = '';
    animationKey.value++;
  } catch (error) {
    console.error('获取指定日期匹配统计失败', error);
    resetStats();
    errorMessage.value = '获取指定日期数据失败，请检查日期参数后稍后重试。';
  } finally {
    loading.value = false;
  }
}

async function loadRangeStats() {
  if (!dateRange.value || dateRange.value.length !== 2) {
    errorMessage.value = '请选择开始日期和截至日期';
    return;
  }

  loading.value = true;
  errorMessage.value = '';

  try {
    const [startDate, endDate] = dateRange.value;
    const response = await axios.get<MatchingStatsResponse>(
      apiUrl(base.matchingRecordRange),
      {
        params: { start_date: startDate, end_date: endDate },
      },
    );
    stats.total = response.data.data?.total ?? 0;
    stats.records = response.data.data?.records ?? [];
    hoveredBarKey.value = '';
    selectedBarKey.value = '';
    animationKey.value++;
  } catch (error) {
    console.error('获取日期范围匹配统计失败', error);
    resetStats();
    errorMessage.value = '获取指定日期范围数据失败，请检查日期参数后稍后重试。';
  } finally {
    loading.value = false;
  }
}

function switchView(view: ViewMode) {
  if (activeView.value === view) {
    return;
  }
  activeView.value = view;

  if (view === 'date') {
    if (singleDate.value) {
      void loadDateStats();
    } else {
      resetStats();
    }
  } else if (view === 'range') {
    if (dateRange.value && dateRange.value.length === 2) {
      void loadRangeStats();
    } else {
      resetStats();
    }
  } else {
    void loadPreset(view);
  }
}

function refreshCurrent() {
  if (activeView.value === 'date') {
    void loadDateStats();
  } else if (activeView.value === 'range') {
    void loadRangeStats();
  } else {
    void loadPreset(activeView.value);
  }
}

const setHoveredBar = (key: string) => {
  hoveredBarKey.value = key;
};

const clearHoveredBar = (key: string) => {
  if (hoveredBarKey.value === key) {
    hoveredBarKey.value = '';
  }
};

const toggleSelectedBar = (key: string) => {
  selectedBarKey.value = selectedBarKey.value === key ? '' : key;
};

const formatUserLabel = (name: string, id: number) => `${name || '未知用户'}[${id || '-'}]`;

const formatTime = (value: string) => {
  if (!value) {
    return '-';
  }

  return new Date(value).toLocaleString('zh-CN', {
    hour12: false,
  });
};

onMounted(() => {
  void loadPreset('today');
});
</script>

<style scoped>
.stats-page {
  min-height: 100vh;
  padding: 28px 32px;
  background: #131320;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.page-header h1 {
  margin: 0 0 6px;
  font-size: 24px;
  font-weight: 700;
  color: #e4e4e7;
}

.page-header p {
  margin: 0;
  color: #8b8b9e;
  font-size: 14px;
}

/* 主面板卡片 */
.stats-panel {
  max-width: 1200px;
  margin: 0 auto;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: #1e1e2d;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.stats-panel :deep(.el-card__body) {
  background: #1e1e2d;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
}

.button-group :deep(.el-button) {
  border-radius: 10px;
  font-weight: 500;
  background: #2a2a3c;
  border-color: rgba(255, 255, 255, 0.08);
  color: #c0c0d0;
  transition: all 0.2s ease;
}

.button-group :deep(.el-button:hover) {
  background: #353548;
  border-color: rgba(96, 165, 250, 0.3);
  color: #60a5fa;
}

.button-group :deep(.el-button.el-button--primary) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  color: #fff;
}

/* 日期选择器深色 */
.range-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.range-form :deep(.el-date-editor) {
  border-radius: 10px;
}

.range-form :deep(.el-range-editor),
.range-form :deep(.el-input__wrapper) {
  background: #2a2a3c;
  border-color: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.08) inset;
}

.range-form :deep(.el-range-editor .el-range-input),
.range-form :deep(.el-input__inner) {
  color: #c0c0d0;
}

.range-form :deep(.el-range-editor .el-range-separator) {
  color: #8b8b9e;
}

.range-form :deep(.el-input__inner::placeholder) {
  color: #6b6b7e;
}

/* Empty 深色 */
.stats-panel :deep(.el-empty__description) {
  color: #8b8b9e;
}

.summary-row {
  display: grid;
  gap: 16px;
  margin-bottom: 16px;
}

/* 第一档：Hero 卡片，横跨整行 */
.summary-row--primary {
  grid-template-columns: 1fr;
}

.summary-row--primary .summary-card :deep(.el-card__body) {
  padding: 28px 32px;
}

.summary-row--primary .summary-card strong {
  font-size: 42px;
}

/* 第二档：次级卡片，两栏并列 */
.summary-row--secondary {
  grid-template-columns: 1fr 1fr;
}

/* 当前视图徽章 */
.current-view-badge {
  display: inline-block;
  margin-left: 12px;
  padding: 4px 12px;
  font-size: 13px;
  font-weight: 600;
  color: #60a5fa;
  background: rgba(96, 165, 250, 0.12);
  border: 1px solid rgba(96, 165, 250, 0.25);
  border-radius: 999px;
  vertical-align: middle;
  letter-spacing: 0.3px;
}

.summary-card {
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: #2a2a3c;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
  overflow: hidden;
}

.summary-card:hover {
  transform: translateY(-2px);
  border-color: rgba(96, 165, 250, 0.2);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.summary-card :deep(.el-card__body) {
  padding: 20px 24px;
  background: transparent;
}

.summary-card--hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
}

.summary-card--hero:hover {
  border: none;
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.4);
}

.summary-card--hero :deep(.el-card__body) {
  padding: 24px 28px;
}

.summary-card--hero .summary-label {
  color: rgba(255, 255, 255, 0.8);
}

.summary-card--hero strong {
  color: #fff;
  font-size: 32px;
}

.summary-card--hero .summary-sub {
  color: rgba(255, 255, 255, 0.7);
}

.summary-label {
  display: block;
  margin-bottom: 8px;
  color: #8b8b9e;
  font-size: 13px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.summary-card strong {
  font-size: 28px;
  font-weight: 700;
  color: #e4e4e7;
  line-height: 1.2;
}

.summary-sub {
  display: block;
  margin-top: 6px;
  font-size: 13px;
  color: #8b8b9e;
  font-weight: 500;
}

.summary-card--users {
  border-left: 4px solid #60a5fa;
}

.summary-card--users strong {
  color: #93c5fd;
}

.summary-card--avg {
  border-left: 4px solid #a78bfa;
}

.summary-card--avg strong {
  color: #c4b5fd;
}

.chart-card,
.detail-card {
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  background: #2a2a3c;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.chart-card :deep(.el-card__header),
.detail-card :deep(.el-card__header),
.collapsible-card :deep(.el-card__header) {
  background: #2a2a3c;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.chart-card :deep(.el-card__body),
.detail-card :deep(.el-card__body) {
  background: #2a2a3c;
}

.status-message,
.chart-card,
.detail-card {
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-header h2 {
  margin: 0 0 4px;
  font-size: 18px;
  font-weight: 700;
  color: #e4e4e7;
}

.section-header p {
  margin: 0;
  color: #8b8b9e;
  font-size: 13px;
}

/* 刷新按钮深色 */
.section-header :deep(.el-button) {
  color: #60a5fa;
}

.section-header :deep(.el-button:hover) {
  background: rgba(96, 165, 250, 0.1);
}

.chart-scroll {
  overflow-x: auto;
  overflow-y: visible;
  padding: 8px 4px 12px;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  scrollbar-color: #4a4a5e #131320;
}

.chart-scroll::-webkit-scrollbar {
  height: 8px;
}

.chart-scroll::-webkit-scrollbar-track {
  background: #131320;
}

.chart-scroll::-webkit-scrollbar-thumb {
  background: #4a4a5e;
  border-radius: 4px;
  border: 2px solid #131320;
}

.chart-scroll::-webkit-scrollbar-thumb:hover {
  background: #5a5a6e;
}

/* 选中柱子详情面板：深色 */
.selected-bar-panel {
  display: grid;
  gap: 14px;
  margin-bottom: 16px;
  padding: 16px 18px;
  border: 1px solid rgba(96, 165, 250, 0.25);
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.12), rgba(118, 75, 162, 0.08));
  box-shadow: 0 8px 24px rgba(96, 165, 250, 0.08);
}

.selected-bar-panel__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.selected-bar-panel__header strong {
  display: block;
  margin-top: 4px;
  font-size: 16px;
  color: #e4e4e7;
}

.selected-bar-panel__hint {
  color: #8b8b9e;
  font-size: 12px;
  text-align: right;
}

.selected-bar-panel__meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px 14px;
  color: #c0c0d0;
  font-size: 14px;
}

.chart-wrapper {
  display: grid;
  grid-template-columns: 56px minmax(0, 1fr);
  gap: 16px;
  align-items: stretch;
  min-height: 460px;
}

.chart-y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-end;
  padding: 120px 0 72px;
  color: #8b8b9e;
  font-size: 12px;
}

.chart-area {
  position: relative;
  display: flex;
  min-height: 460px;
  padding-top: 120px;
}

.chart-grid {
  position: absolute;
  inset: 120px 0 72px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  pointer-events: none;
}

.grid-line {
  border-top: 1px dashed rgba(255, 255, 255, 0.06);
}

.bars-row {
  position: relative;
  z-index: 1;
  display: grid;
  grid-auto-flow: column;
  grid-auto-columns: minmax(72px, 1fr);
  gap: 12px;
  align-items: end;
  width: 100%;
  min-height: 100%;
  padding: 0 8px;
}

.bar-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: end;
  gap: 8px;
  min-height: 100%;
}

.bar-track {
  display: flex;
  flex: 1;
  width: 100%;
  min-height: 268px;
  flex-direction: column;
  align-items: center;
  justify-content: end;
}

.bar-interactive {
  position: relative;
  display: flex;
  width: 100%;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.bar-value {
  font-size: 13px;
  color: #60a5fa;
  font-weight: 600;
}

.bar {
  width: 100%;
  min-height: 12px;
  border: none;
  border-radius: 10px 10px 0 0;
  background: linear-gradient(180deg, #93c5fd 0%, #60a5fa 100%);
  box-shadow: 0 8px 16px rgba(96, 165, 250, 0.25);
  cursor: pointer;
  transition: height 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
}

.bar:hover,
.bar:focus-visible,
.bar.is-active {
  filter: saturate(1.15);
  transform: translateY(-4px);
  box-shadow: 0 14px 24px rgba(96, 165, 250, 0.4);
  outline: none;
}

.bar-label {
  min-height: 48px;
  color: #c0c0d0;
  font-size: 12px;
  line-height: 1.4;
  text-align: center;
  word-break: break-word;
}

/* 柱状图气泡：深色玻璃 */
.bar-popover {
  position: absolute;
  inset: auto auto calc(100% + 12px) 50%;
  z-index: 3;
  display: grid;
  min-width: 196px;
  max-width: min(260px, calc(100vw - 56px));
  gap: 6px;
  padding: 14px 16px;
  border: 1px solid rgba(96, 165, 250, 0.3);
  border-radius: 14px;
  background: rgba(30, 30, 45, 0.96);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.5);
  color: #e4e4e7;
  backdrop-filter: blur(16px);
  transform: translateX(-50%);
}

.bar-popover::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: -8px;
  width: 14px;
  height: 14px;
  border-right: 1px solid rgba(96, 165, 250, 0.3);
  border-bottom: 1px solid rgba(96, 165, 250, 0.3);
  background: rgba(30, 30, 45, 0.96);
  transform: translateX(-50%) rotate(45deg);
}

.bar-popover strong {
  font-size: 15px;
}

.detail-caption {
  color: #8b8b9e;
  font-size: 12px;
}

.detail-fade-enter-active,
.detail-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.detail-fade-enter-from,
.detail-fade-leave-to {
  opacity: 0;
  transform: translateY(6px);
}

/* ============== 启动动画 ============== */
.content-enter {
  animation: contentFadeIn 0.4s ease-out both;
}

@keyframes contentFadeIn {
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

/* Hero 卡片：从下方滑入 + 缩放 */
.enter-hero {
  animation: heroEnter 0.6s cubic-bezier(0.34, 1.2, 0.64, 1) both;
}

@keyframes heroEnter {
  0% {
    opacity: 0;
    transform: translateY(30px) scale(0.96);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* 次级卡片 / 图表卡片 / 详情卡片：依次淡入上移 */
.enter-card {
  animation: cardEnter 0.5s cubic-bezier(0.34, 1.1, 0.64, 1) both;
  animation-delay: var(--enter-delay, 0s);
}

@keyframes cardEnter {
  0% {
    opacity: 0;
    transform: translateY(20px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 柱子：从底部生长 */
.enter-bar {
  animation: barGrow 0.6s cubic-bezier(0.34, 1.15, 0.64, 1) both;
  transform-origin: bottom;
}

@keyframes barGrow {
  0% {
    opacity: 0;
    transform: scaleY(0);
  }
  60% {
    opacity: 1;
  }
  100% {
    opacity: 1;
    transform: scaleY(1);
  }
}

/* 柱子数值标签：淡入 */
.enter-bar-value {
  animation: barLabelFade 0.4s ease-out both;
}

/* 柱子底部标签：淡入 */
.enter-bar-label {
  animation: barLabelFade 0.4s ease-out both;
}

@keyframes barLabelFade {
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

/* ============== 表格深色主题 ============== */
.detail-card :deep(.el-table) {
  background: transparent;
  color: #c0c0d0;
  --el-table-bg-color: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-header-bg-color: #1e1e2d;
  --el-table-header-text-color: #8b8b9e;
  --el-table-text-color: #c0c0d0;
  --el-table-border-color: rgba(255, 255, 255, 0.06);
  --el-table-row-hover-bg-color: rgba(96, 165, 250, 0.08);
}

.detail-card :deep(.el-table th.el-table__cell) {
  background: #1e1e2d;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.detail-card :deep(.el-table tr) {
  background: transparent;
}

.detail-card :deep(.el-table td.el-table__cell),
.detail-card :deep(.el-table .el-table__cell) {
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.detail-card :deep(.el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell) {
  background: rgba(255, 255, 255, 0.02);
}

.detail-card :deep(.el-table__body tr:hover > td.el-table__cell) {
  background: rgba(96, 165, 250, 0.08) !important;
}

@media (max-width: 768px) {
  .stats-page {
    padding: 14px;
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    margin-bottom: 16px;
  }

  .button-group {
    width: 100%;
    gap: 10px;
    margin-bottom: 14px;
  }

  .button-group :deep(.el-button) {
    flex: 1 1 calc(50% - 5px);
    min-height: 40px;
    margin-left: 0;
    white-space: normal;
  }

  .range-form {
    width: 100%;
  }

  .range-form :deep(.el-date-editor) {
    width: 100%;
  }

  .summary-row {
    gap: 12px;
    margin-bottom: 14px;
  }

  .summary-row--primary .summary-card :deep(.el-card__body) {
    padding: 20px 24px;
  }

  .summary-row--primary .summary-card strong {
    font-size: 28px;
  }

  .summary-row--secondary {
    grid-template-columns: 1fr;
  }

  .current-view-badge {
    display: block;
    margin: 8px 0 0;
  }

  .section-header h2 {
    font-size: 16px;
  }

  .chart-scroll {
    margin: 0 -8px;
    padding: 4px 8px 12px;
  }

  .selected-bar-panel {
    gap: 12px;
    margin-bottom: 14px;
    padding: 14px;
  }

  .selected-bar-panel__header {
    flex-direction: column;
  }

  .selected-bar-panel__hint,
  .selected-bar-panel__meta {
    text-align: left;
  }

  .chart-wrapper {
    grid-template-columns: 40px minmax(0, 1fr);
    gap: 10px;
    min-height: 400px;
  }

  .chart-y-axis {
    padding: 108px 0 64px;
    font-size: 11px;
  }

  .chart-area {
    min-height: 400px;
    padding-top: 108px;
  }

  .chart-grid {
    inset: 108px 0 64px;
  }

  .bars-row {
    grid-auto-columns: minmax(60px, 1fr);
    gap: 10px;
    padding: 0 4px;
  }

  .bar-track {
    min-height: 228px;
  }

  .bar-popover {
    min-width: 168px;
    max-width: min(220px, calc(100vw - 40px));
    padding: 12px 14px;
    font-size: 12px;
  }

  .bar-label {
    min-height: 40px;
    font-size: 11px;
  }

  .detail-card :deep(.el-card__body),
  .chart-card :deep(.el-card__body) {
    padding: 14px;
  }

  .detail-card :deep(.el-table) {
    font-size: 12px;
  }
}

/* ============== 可折叠卷帘（统计汇总 / 原始数据卡片） ============== */
.collapsible-card :deep(.el-card__header) {
  padding: 14px 20px;
}

.collapsible-header {
  cursor: pointer;
  user-select: none;
}

.collapsible-arrow {
  font-size: 16px;
  color: #8b8b9e;
  transform: rotate(-90deg);
  transition: transform 0.25s ease, color 0.25s ease;
  flex-shrink: 0;
  margin-left: 16px;
}

.collapsible-arrow.is-open {
  transform: rotate(0deg);
  color: #60a5fa;
}

.card-collapse-enter-active,
.card-collapse-leave-active {
  overflow: hidden;
  transition: max-height 0.3s ease, opacity 0.25s ease,
    margin-top 0.28s ease, margin-bottom 0.28s ease;
}

.card-collapse-enter-from,
.card-collapse-leave-to {
  max-height: 0 !important;
  opacity: 0;
}

.card-collapse-enter-to,
.card-collapse-leave-from {
  max-height: 2500px;
  opacity: 1;
}
</style>
