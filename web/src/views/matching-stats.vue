<template>
  <div class="stats-page">
    <div class="page-header">
      <div>
        <h1>匹配统计</h1>
        <p>点击按钮切换当日或累计统计，图表会按“用户名[QQ号]”展示成功匹配次数。</p>
      </div>
      <el-button plain @click="router.push('/')">返回首页</el-button>
    </div>

    <el-card shadow="never" class="stats-panel">
      <div class="button-group">
        <el-button
          :type="activeView === 'today' ? 'primary' : 'default'"
          @click="loadStats('today')"
        >
          显示当日的匹配统计
        </el-button>
        <el-button
          :type="activeView === 'all' ? 'primary' : 'default'"
          @click="loadStats('all')"
        >
          显示累计匹配统计
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
          <div class="summary-row">
            <el-card shadow="hover" class="summary-card">
              <span class="summary-label">当前视图</span>
              <strong>{{ currentLabel }}</strong>
            </el-card>
            <el-card shadow="hover" class="summary-card">
              <span class="summary-label">成功匹配记录数</span>
              <strong>{{ stats.total }}</strong>
            </el-card>
            <el-card shadow="hover" class="summary-card">
              <span class="summary-label">去重后的用户数</span>
              <strong>{{ chartData.length }}</strong>
            </el-card>
          </div>

          <el-empty
            v-if="stats.records.length === 0"
            description="当前没有匹配记录"
          />

          <template v-else>
            <el-card shadow="never" class="chart-card">
              <template #header>
                <div class="section-header">
                  <div>
                    <h2>匹配成功统计图</h2>
                    <p>X 轴为 name[qq号]，Y 轴为该用户成功匹配的次数。</p>
                  </div>
                </div>
              </template>

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
                        v-for="item in chartData"
                        :key="item.key"
                        class="bar-column"
                      >
                        <span class="bar-value">{{ item.count }}</span>
                        <div
                          class="bar"
                          :style="{ height: `${item.heightPercent}%` }"
                          :title="`${item.label}：${item.count}`"
                        />
                        <span class="bar-label">{{ item.label }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>

            <el-card shadow="never" class="detail-card">
              <template #header>
                <div class="section-header">
                  <div>
                    <h2>统计汇总</h2>
                    <p>按用户聚合后的成功匹配次数，便于快速核对图表数据。</p>
                  </div>
                </div>
              </template>

              <el-table :data="chartData" stripe>
                <el-table-column prop="label" label="用户(name[qq号])" min-width="220" />
                <el-table-column prop="name" label="用户名" min-width="140" />
                <el-table-column prop="id" label="QQ号" min-width="120" />
                <el-table-column prop="count" label="成功次数" min-width="120" sortable />
              </el-table>
            </el-card>

            <el-card shadow="never" class="detail-card">
              <template #header>
                <div class="section-header">
                  <div>
                    <h2>原始数据</h2>
                    <p>图表下方保留完整匹配记录，方便继续排查和对照。</p>
                  </div>
                </div>
              </template>

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
            </el-card>
          </template>
        </template>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'MatchingStatsPage' });

import axios from 'axios';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';

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
  heightPercent: number;
}

const router = useRouter();
const loading = ref(false);
const errorMessage = ref('');
const activeView = ref<'today' | 'all'>('today');
const stats = reactive({
  total: 0,
  records: [] as MatchingRecord[],
});

const endpointMap = {
  today: 'http://127.0.0.1:3000/api/matching/record/today',
  all: 'http://127.0.0.1:3000/api/matching/record/all',
} as const;

const currentLabel = computed(() =>
  activeView.value === 'today' ? '当日匹配统计' : '累计匹配统计',
);

const chartData = computed<UserChartItem[]>(() => {
  const counter = new Map<string, Omit<UserChartItem, 'count' | 'heightPercent'>>();
  const countMap = new Map<string, number>();

  const upsertUser = (id: number, name: string) => {
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
  };

  stats.records.forEach((record) => {
    upsertUser(record.user_id, record.user_name);
    upsertUser(record.peer_id, record.peer_name);
  });

  const maxCount = Math.max(...countMap.values(), 0);

  return Array.from(countMap.entries())
    .map(([key, count]) => {
      const user = counter.get(key)!;
      return {
        ...user,
        count,
        heightPercent: maxCount > 0 ? Math.max((count / maxCount) * 100, 8) : 0,
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

const loadStats = async (view: 'today' | 'all') => {
  loading.value = true;
  errorMessage.value = '';
  activeView.value = view;

  try {
    const response = await axios.get<MatchingStatsResponse>(endpointMap[view]);
    stats.total = response.data.data?.total ?? 0;
    stats.records = response.data.data?.records ?? [];
  } catch (error) {
    console.error('获取匹配统计失败', error);
    stats.total = 0;
    stats.records = [];
    errorMessage.value = '获取匹配统计失败，请稍后重试。';
  } finally {
    loading.value = false;
  }
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
  void loadStats('today');
});
</script>

<style scoped>
.stats-page {
  min-height: 100vh;
  padding: 32px;
  background: #f5f7fa;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0 0 8px;
}

.page-header p {
  margin: 0;
  color: #606266;
}

.stats-panel {
  max-width: 1200px;
  margin: 0 auto;
}

.button-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 24px;
}

.summary-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.summary-card,
.chart-card,
.detail-card {
  border-radius: 12px;
}

.summary-label {
  display: block;
  margin-bottom: 8px;
  color: #909399;
  font-size: 14px;
}

.status-message,
.chart-card,
.detail-card {
  margin-bottom: 24px;
}

.section-header h2 {
  margin: 0 0 6px;
  font-size: 18px;
}

.section-header p {
  margin: 0;
  color: #606266;
}

.chart-scroll {
  overflow-x: auto;
  padding-bottom: 8px;
}

.chart-wrapper {
  display: grid;
  grid-template-columns: 56px minmax(0, 1fr);
  gap: 16px;
  align-items: stretch;
  min-height: 380px;
}

.chart-y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-end;
  padding: 8px 0 72px;
  color: #909399;
  font-size: 12px;
}

.chart-area {
  position: relative;
  display: flex;
  min-height: 380px;
  padding-top: 8px;
}

.chart-grid {
  position: absolute;
  inset: 8px 0 72px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  pointer-events: none;
}

.grid-line {
  border-top: 1px dashed #dcdfe6;
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

.bar-value {
  font-size: 13px;
  color: #409eff;
  font-weight: 600;
}

.bar {
  width: 100%;
  min-height: 8px;
  border-radius: 10px 10px 0 0;
  background: linear-gradient(180deg, #79bbff 0%, #409eff 100%);
  box-shadow: 0 8px 16px rgb(64 158 255 / 18%);
  transition: height 0.2s ease;
}

.bar-label {
  min-height: 48px;
  color: #303133;
  font-size: 12px;
  line-height: 1.4;
  text-align: center;
  word-break: break-word;
}

@media (max-width: 768px) {
  .stats-page {
    padding: 20px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .button-group {
    width: 100%;
  }

  .button-group :deep(.el-button) {
    flex: 1;
    margin-left: 0;
  }
}
</style>
