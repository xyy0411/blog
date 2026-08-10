<template>
  <div class="stats-page">
    <div class="page-header">
      <div>
        <h1>
          匹配申请统计
          <span class="current-view-badge">{{ currentLabel }}</span>
        </h1>
        <p>展示匹配申请的成功与失败比例。支持当日数据查询和指定日期范围查询。</p>
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
          :type="activeView === 'range' ? 'primary' : 'default'"
          @click="switchView('range')"
        >
          按日期范围
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

      <el-skeleton v-if="loading" :rows="6" animated />

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
          <!-- 第一档：申请总数（Hero 卡片，横跨整行） -->
          <div class="summary-row summary-row--primary">
            <el-card shadow="hover" class="summary-card summary-card--hero">
              <span class="summary-label">申请总数</span>
              <strong>{{ stats.total }}</strong>
              <span class="summary-sub">匹配申请记录总量</span>
            </el-card>
          </div>

          <!-- 第二档：匹配成功 / 匹配失败 -->
          <div class="summary-row summary-row--secondary">
            <el-card shadow="hover" class="summary-card summary-card--success">
              <span class="summary-label">匹配成功</span>
              <strong>{{ stats.successCount }}</strong>
              <span class="summary-sub">{{ formatRate(stats.successRate) }}</span>
            </el-card>
            <el-card shadow="hover" class="summary-card summary-card--fail">
              <span class="summary-label">匹配失败</span>
              <strong>{{ stats.failCount }}</strong>
              <span class="summary-sub">{{ formatRate(stats.failRate) }}</span>
            </el-card>
          </div>

          <el-empty
            v-if="stats.total === 0"
            description="当前条件下没有匹配申请记录"
          />

          <template v-else>
            <el-card shadow="never" class="chart-card">
              <template #header>
                <div class="section-header">
                  <div>
                    <h2>匹配成功率统计图</h2>
                    <p>扇形图展示匹配成功与失败的比例，悬停切片可查看详细数据。</p>
                  </div>
                  <el-button text :loading="loading" @click="refreshCurrent">
                    刷新数据
                  </el-button>
                </div>
              </template>

              <div class="pie-wrapper">
                <svg
                  :key="pieAnimationKey"
                  class="pie-svg pie-enter"
                  :viewBox="`0 0 ${svgSize} ${svgSize}`"
                  :style="{ width: `${svgSize}px`, height: `${svgSize}px` }"
                  role="img"
                  aria-label="匹配成功与失败比例扇形图"
                >
                  <path
                    v-for="(slice, index) in pieSlices"
                    :key="slice.key"
                    :d="slice.path"
                    :fill="slice.color"
                    :stroke="slice.key === hoveredSliceKey ? 'rgba(255,255,255,0.9)' : '#1e1e2d'"
                    stroke-width="2"
                    :transform="slice.transform"
                    class="pie-slice"
                    :style="{ animationDelay: `${index * 0.15}s` }"
                    @mouseenter="hoveredSliceKey = slice.key"
                    @mouseleave="hoveredSliceKey = ''"
                  >
                    <title>{{ slice.label }}：{{ slice.count }} 条 ({{ formatRate(slice.rate) }})</title>
                  </path>
                  <text
                    v-for="slice in pieSlices"
                    v-show="slice.count > 0"
                    :key="`label-${slice.key}`"
                    :x="slice.labelX"
                    :y="slice.labelY"
                    text-anchor="middle"
                    dominant-baseline="middle"
                    class="pie-label"
                  >
                    {{ slice.label }} {{ formatRate(slice.rate) }}
                  </text>
                </svg>

                <div class="pie-legend">
                    <div
                      v-for="slice in pieSlices"
                      :key="`legend-${slice.key}`"
                      class="legend-item"
                      :class="{ 'is-hovered': hoveredSliceKey === slice.key, 'is-filtered': (slice.key === 'success' && filterReason === 0) || (slice.key === 'fail' && filterReason === -1) }"
                      @mouseenter="hoveredSliceKey = slice.key"
                      @mouseleave="hoveredSliceKey = ''"
                      @click="onLegendClick(slice.key)"
                    >
                    <span class="legend-dot" :style="{ background: slice.color }" />
                    <span class="legend-label">{{ slice.label }}</span>
                    <span class="legend-count">{{ slice.count }} 条</span>
                    <span class="legend-rate">{{ formatRate(slice.rate) }}</span>
                  </div>
                </div>

                <transition name="detail-fade">
                  <div v-if="hoveredSlice" class="pie-popover">
                    <span class="detail-caption">悬浮详情</span>
                    <strong :style="{ color: hoveredSlice.color }">{{ hoveredSlice.label }}</strong>
                    <span>数量：{{ hoveredSlice.count }} 条</span>
                    <span>占比：{{ formatRate(hoveredSlice.rate) }}</span>
                    <span>说明：{{ hoveredSlice.description }}</span>
                  </div>
                </transition>
              </div>
            </el-card>

            <el-card shadow="never" class="chart-card">
              <template #header>
                <div class="section-header">
                  <div>
                    <h2>退出原因统计</h2>
                    <p>展示用户从匹配队列中退出的各种原因占比，包括成功匹配、主动退出、超时等。</p>
                  </div>
                  <div>
                    <el-button-group>
                      <el-button size="small" :type="filterReason === null ? 'primary' : 'default'" @click="setFilterReason(null)">全部</el-button>
                      <el-button size="small" :type="filterReason === 0 ? 'primary' : 'default'" @click="setFilterReason(0)">匹配成功</el-button>
                      <el-button size="small" :type="filterReason === 1 ? 'primary' : 'default'" @click="setFilterReason(1)">用户主动退出</el-button>
                      <el-button size="small" :type="filterReason === 2 ? 'primary' : 'default'" @click="setFilterReason(2)">匹配超时</el-button>
                      <el-button size="small" :type="filterReason === 3 ? 'primary' : 'default'" @click="setFilterReason(3)">网络错误</el-button>
                      <el-button size="small" :type="filterReason === 4 ? 'primary' : 'default'" @click="setFilterReason(4)">已过期</el-button>
                    </el-button-group>
                  </div>
                </div>
              </template>

              <div class="exit-reason-stats">
                <div class="stats-grid">
                  <!-- 匹配成功 -->
                  <div class="stat-item">
                    <div class="stat-header">
                      <span class="stat-label">匹配成功</span>
                      <span class="stat-badge success">{{ stats.exitReasonStats.success.count }}</span>
                    </div>
                    <div class="progress-bar">
                      <div
                        class="progress-fill success"
                        :style="{ width: `${stats.exitReasonStats.success.rate}%` }"
                      />
                    </div>
                    <div class="stat-footer">{{ formatRate(stats.exitReasonStats.success.rate) }}</div>
                  </div>

                  <!-- 用户主动退出 -->
                  <div class="stat-item">
                    <div class="stat-header">
                      <span class="stat-label">用户主动退出</span>
                      <span class="stat-badge info">{{ stats.exitReasonStats.user_initiated.count }}</span>
                    </div>
                    <div class="progress-bar">
                      <div
                        class="progress-fill info"
                        :style="{ width: `${stats.exitReasonStats.user_initiated.rate}%` }"
                      />
                    </div>
                    <div class="stat-footer">{{ formatRate(stats.exitReasonStats.user_initiated.rate) }}</div>
                  </div>

                  <!-- 匹配超时 -->
                  <div class="stat-item">
                    <div class="stat-header">
                      <span class="stat-label">匹配超时</span>
                      <span class="stat-badge warning">{{ stats.exitReasonStats.timeout.count }}</span>
                    </div>
                    <div class="progress-bar">
                      <div
                        class="progress-fill warning"
                        :style="{ width: `${stats.exitReasonStats.timeout.rate}%` }"
                      />
                    </div>
                    <div class="stat-footer">{{ formatRate(stats.exitReasonStats.timeout.rate) }}</div>
                  </div>

                  <!-- 错误/断连 -->
                  <div class="stat-item">
                    <div class="stat-header">
                      <span class="stat-label">网络错误</span>
                      <span class="stat-badge danger">{{ stats.exitReasonStats.error.count }}</span>
                    </div>
                    <div class="progress-bar">
                      <div
                        class="progress-fill danger"
                        :style="{ width: `${stats.exitReasonStats.error.rate}%` }"
                      />
                    </div>
                    <div class="stat-footer">{{ formatRate(stats.exitReasonStats.error.rate) }}</div>
                  </div>

                  <!-- 已过期 -->
                  <div class="stat-item">
                    <div class="stat-header">
                      <span class="stat-label">已过期</span>
                      <span class="stat-badge secondary">{{ stats.exitReasonStats.expired.count }}</span>
                    </div>
                    <div class="progress-bar">
                      <div
                        class="progress-fill secondary"
                        :style="{ width: `${stats.exitReasonStats.expired.rate}%` }"
                      />
                    </div>
                    <div class="stat-footer">{{ formatRate(stats.exitReasonStats.expired.rate) }}</div>
                  </div>
                </div>
              </div>
            </el-card>

            <el-card shadow="never" class="detail-card collapsible-card">
              <template #header>
                <div
                  class="section-header collapsible-header"
                  @click="rawDataCollapsed = !rawDataCollapsed"
                >
                  <div>
                    <h2>原始数据</h2>
                    <p>完整保留所有匹配申请记录，便于核对图表数据。</p>
                  </div>
                  <div style="display: flex; gap: 8px; align-items: center;">
                    <el-button
                      v-if="!rawDataCollapsed"
                      size="small"
                      @click.stop="showAllData"
                    >
                      显示全部数据
                    </el-button>
                    <el-icon
                      class="collapsible-arrow"
                      :class="{ 'is-open': !rawDataCollapsed }"
                    >
                      <ArrowDown />
                    </el-icon>
                  </div>
                </div>
              </template>

              <transition name="card-collapse">
                <div v-show="!rawDataCollapsed">
                  <el-table :data="displayTableData" stripe>
                    <el-table-column prop="id" label="ID" min-width="80" />
                    <el-table-column label="用户" min-width="200">
                      <template #default="scope">
                        {{ formatUserLabel(scope.row.user_name, scope.row.user_id) }}
                      </template>
                    </el-table-column>
                    <el-table-column label="匹配结果" min-width="120">
                      <template #default="scope">
                        <el-tag
                          :type="scope.row.is_matched ? 'success' : 'danger'"
                          disable-transitions
                        >
                          {{ scope.row.is_matched ? '成功' : '失败' }}
                        </el-tag>
                      </template>
                    </el-table-column>
                    <el-table-column prop="duration" label="匹配时长(秒)" min-width="120" sortable />
                    <el-table-column prop="match_id" label="匹配 ID" min-width="180" />
                    <el-table-column label="退出原因" min-width="160">
                      <template #default="scope">
                        {{ exitReasonLabel(scope.row.exit_reason) }}
                      </template>
                    </el-table-column>
                    <el-table-column label="创建时间" min-width="180">
                      <template #default="scope">
                        {{ formatTime(scope.row.created_at) }}
                      </template>
                    </el-table-column>
                  </el-table>
                  <div v-if="tablePageSize < stats.records.length" class="table-footer">
                    <span class="info-text">
                      显示 {{ displayTableData.length }} / {{ stats.records.length }} 条记录
                    </span>
                    <el-button
                      type="primary"
                      text
                      @click="showAllData"
                    >
                      加载全部
                    </el-button>
                  </div>
                  <div v-else-if="stats.records.length > 0" class="table-footer">
                    <span class="info-text">共 {{ stats.records.length }} 条记录</span>
                  </div>
                </div>
              </transition>
            </el-card>
          </template>
        </template>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'MatchingApplicationPage' });

import axios from 'axios';
import { ArrowDown } from '@element-plus/icons-vue';
import base, { apiUrl } from '@/api/api.ts';
import blogStore from '@/store/arlog.ts';
import { computed, onMounted, reactive, ref } from 'vue';

interface MatchingApplication {
  id: number;
  user_id: number;
  user_name: string;
  is_matched: boolean;
  duration: number;
  match_id: string;
  created_at: string;
  exit_reason?: number;
}

interface ExitReasonItem {
  count: number;
  rate: number;
}

interface ApplicationStatsResponse {
  data?: {
    total?: number;
    success_count?: number;
    fail_count?: number;
    success_rate?: number;
    fail_rate?: number;
    exit_reason_stats?: {
      success?: ExitReasonItem;
      user_initiated?: ExitReasonItem;
      timeout?: ExitReasonItem;
      error?: ExitReasonItem;
      expired?: ExitReasonItem;
    };
    records?: MatchingApplication[];
  };
  error?: string;
}

interface PieSlice {
  key: string;
  label: string;
  description: string;
  count: number;
  rate: number;
  color: string;
  path: string;
  transform: string;
  labelX: number;
  labelY: number;
}

type ViewMode = 'today' | 'range' | 'week' | 'all';

const store = blogStore();
const loading = ref(false);
const errorMessage = ref('');
const activeView = ref<ViewMode>('today');
const dateRange = ref<[string, string] | null>(null);
const hoveredSliceKey = ref('');
const rawDataCollapsed = ref(true);
const pieAnimationKey = ref(0);
const tablePageSize = ref(10);
// filter by exit reason: null => all, otherwise numeric code
const filterReason = ref<number | null>(null);

const stats = reactive({
  total: 0,
  successCount: 0,
  failCount: 0,
  successRate: 0,
  failRate: 0,
  exitReasonStats: {
    success: { count: 0, rate: 0 },
    user_initiated: { count: 0, rate: 0 },
    timeout: { count: 0, rate: 0 },
    error: { count: 0, rate: 0 },
    expired: { count: 0, rate: 0 },
  },
  records: [] as MatchingApplication[],
});

const svgSize = 320;
const pieRadius = 130;
const pieCenter = svgSize / 2;
const successColor = '#10b981';
const failColor = '#ef4444';

const currentLabel = computed(() => {
  switch (activeView.value) {
    case 'today':
      return '当日匹配申请统计';
    case 'week':
      return '本周匹配申请统计';
    case 'all':
      return '累计匹配申请统计';
    case 'range':
    default:
      return '指定日期范围匹配申请统计';
  }
});

const authHeaders = computed(() => {
  const headers: Record<string, string> = {};
  if (store.token) {
    headers.Authorization = store.token.startsWith('Bearer ')
      ? store.token
      : `Bearer ${store.token}`;
  }
  return headers;
});

const pieSlices = computed<PieSlice[]>(() => {
  if (stats.total === 0) {
    return [];
  }

  const successAngle = (stats.successCount / stats.total) * Math.PI * 2;
  const items = [
    {
      key: 'success',
      label: '匹配成功',
      description: '该时段内匹配成功的记录。',
      count: stats.successCount,
      rate: stats.successRate,
      color: successColor,
      startAngle: 0,
      endAngle: successAngle,
    },
    {
      key: 'fail',
      label: '匹配失败',
      description: '该时段内匹配失败的记录。',
      count: stats.failCount,
      rate: stats.failRate,
      color: failColor,
      startAngle: successAngle,
      endAngle: Math.PI * 2,
    },
  ];

  return items
    .filter((item) => item.count > 0)
    .map((item) => buildSlice(item));
});

const hoveredSlice = computed(() =>
  pieSlices.value.find((slice) => slice.key === hoveredSliceKey.value) ?? null,
);

// filter helpers
function setFilterReason(reason: number | null) {
  filterReason.value = reason;
}

function onLegendClick(key: string) {
  // legend only has 'success' and 'fail'
  if (key === 'success') {
    setFilterReason(0);
    return;
  }
  if (key === 'fail') {
    // -1 means filter by is_matched === false
    setFilterReason(-1);
    return;
  }
}

function exitReasonLabel(code?: number) {
  switch (code) {
    case 0:
      return '匹配成功';
    case 1:
      return '用户主动退出';
    case 2:
      return '匹配超时';
    case 3:
      return '网络错误';
    case 4:
      return '已过期';
    default:
      return '-';
  }
}

const displayTableData = computed(() => {
  // apply exit reason filter first
  let filtered = stats.records;
  if (filterReason.value != null) {
    if (filterReason.value === -1) {
      filtered = stats.records.filter((r) => !r.is_matched);
    } else {
      filtered = stats.records.filter((r) => r.exit_reason === filterReason.value);
    }
  }
  const start = 0;
  const end = tablePageSize.value;
  return filtered.slice(start, end);
});

interface SliceInput {
  key: string;
  label: string;
  description: string;
  count: number;
  rate: number;
  color: string;
  startAngle: number;
  endAngle: number;
}

function buildSlice(input: SliceInput): PieSlice {
  const { path, labelX, labelY } = describeArc(
    pieCenter,
    pieCenter,
    pieRadius,
    input.startAngle,
    input.endAngle,
  );

  const midAngle = (input.startAngle + input.endAngle) / 2;
  const offsetX = Math.cos(midAngle) * 8;
  const offsetY = Math.sin(midAngle) * 8;
  const transform = hoveredSliceKey.value === input.key
    ? `translate(${offsetX} ${offsetY})`
    : '';

  return {
    key: input.key,
    label: input.label,
    description: input.description,
    count: input.count,
    rate: input.rate,
    color: input.color,
    path,
    transform,
    labelX,
    labelY,
  };
}

function polarToCartesian(cx: number, cy: number, r: number, angle: number) {
  return {
    x: cx + r * Math.cos(angle),
    y: cy + r * Math.sin(angle),
  };
}

function describeArc(cx: number, cy: number, r: number, startAngle: number, endAngle: number) {
  // 当切片占满整圆时，endAngle 与 startAngle 会重合，导致 path 不可见；
  // 这里将 endAngle 略微回退以避免该退化情况。
  const safeEnd = endAngle - startAngle >= Math.PI * 2 ? endAngle - 0.0001 : endAngle;
  const start = polarToCartesian(cx, cy, r, startAngle);
  const end = polarToCartesian(cx, cy, r, safeEnd);
  const largeArc = safeEnd - startAngle > Math.PI ? 1 : 0;

  const path = `M ${cx} ${cy} L ${start.x} ${start.y} A ${r} ${r} 0 ${largeArc} 1 ${end.x} ${end.y} Z`;
  const labelAngle = (startAngle + safeEnd) / 2;
  const labelPos = polarToCartesian(cx, cy, r * 0.62, labelAngle);

  return { path, labelX: labelPos.x, labelY: labelPos.y };
}

function applyResponse(payload: ApplicationStatsResponse) {
  const data = payload.data ?? {};
  stats.total = data.total ?? 0;
  stats.successCount = data.success_count ?? 0;
  stats.failCount = data.fail_count ?? 0;
  stats.successRate = data.success_rate ?? 0;
  stats.failRate = data.fail_rate ?? 0;
  stats.exitReasonStats = data.exit_reason_stats ?? {
    success: { count: 0, rate: 0 },
    user_initiated: { count: 0, rate: 0 },
    timeout: { count: 0, rate: 0 },
    error: { count: 0, rate: 0 },
    expired: { count: 0, rate: 0 },
  };
  stats.records = data.records ?? [];
  tablePageSize.value = 10;
  pieAnimationKey.value++;
}

function resetStats() {
  stats.total = 0;
  stats.successCount = 0;
  stats.failCount = 0;
  stats.successRate = 0;
  stats.failRate = 0;
  stats.exitReasonStats = {
    success: { count: 0, rate: 0 },
    user_initiated: { count: 0, rate: 0 },
    timeout: { count: 0, rate: 0 },
    error: { count: 0, rate: 0 },
    expired: { count: 0, rate: 0 },
  };
  stats.records = [];
  tablePageSize.value = 10;
}

function formatDateYYYYMMDD(d: Date) {
  const yyyy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, '0');
  const dd = String(d.getDate()).padStart(2, '0');
  return `${yyyy}-${mm}-${dd}`;
}

async function loadWeekStats() {
  loading.value = true;
  errorMessage.value = '';
  try {
    const now = new Date();
    const weekday = now.getDay(); // 0 Sun .. 6 Sat
    // compute Monday as start of week (China commonly uses Monday start)
    const daysToMonday = weekday === 0 ? 6 : weekday - 1;
    const start = new Date(now.getFullYear(), now.getMonth(), now.getDate() - daysToMonday);
    const end = new Date(start.getFullYear(), start.getMonth(), start.getDate() + 6);
    const startStr = formatDateYYYYMMDD(start);
    const endStr = formatDateYYYYMMDD(end);
    const response = await axios.get<ApplicationStatsResponse>(
      apiUrl(base.matchingApplicationRange),
      { params: { start_date: startStr, end_date: endStr }, headers: authHeaders.value },
    );
    applyResponse(response.data);
  } catch (error) {
    console.error('获取本周匹配申请数据失败', error);
    resetStats();
    errorMessage.value = '获取本周匹配申请数据失败，请稍后重试。';
  } finally {
    loading.value = false;
  }
}

async function loadAllStats() {
  loading.value = true;
  errorMessage.value = '';
  try {
    const startStr = '1970-01-01';
    const endStr = formatDateYYYYMMDD(new Date());
    const response = await axios.get<ApplicationStatsResponse>(
      apiUrl(base.matchingApplicationRange),
      { params: { start_date: startStr, end_date: endStr }, headers: authHeaders.value },
    );
    applyResponse(response.data);
  } catch (error) {
    console.error('获取累计匹配申请数据失败', error);
    resetStats();
    errorMessage.value = '获取累计匹配申请数据失败，请稍后重试。';
  } finally {
    loading.value = false;
  }
}
async function loadTodayStats() {
  loading.value = true;
  errorMessage.value = '';

  try {
    const response = await axios.get<ApplicationStatsResponse>(
      apiUrl(base.matchingApplicationToday),
      { headers: authHeaders.value },
    );
    applyResponse(response.data);
  } catch (error) {
    console.error('获取当日匹配申请数据失败', error);
    resetStats();
    errorMessage.value = '获取当日匹配申请数据失败，请确认登录状态后稍后重试。';
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
    const response = await axios.get<ApplicationStatsResponse>(
      apiUrl(base.matchingApplicationRange),
      {
        params: { start_date: startDate, end_date: endDate },
        headers: authHeaders.value,
      },
    );
    applyResponse(response.data);
  } catch (error) {
    console.error('获取日期范围匹配申请数据失败', error);
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
  hoveredSliceKey.value = '';

  if (view === 'today') {
    void loadTodayStats();
  } else if (view === 'week') {
    void loadWeekStats();
  } else if (view === 'all') {
    void loadAllStats();
  } else if (dateRange.value && dateRange.value.length === 2) {
    void loadRangeStats();
  } else {
    resetStats();
  }
}

function refreshCurrent() {
  if (activeView.value === 'today') {
    void loadTodayStats();
  } else if (activeView.value === 'week') {
    void loadWeekStats();
  } else if (activeView.value === 'all') {
    void loadAllStats();
  } else {
    void loadRangeStats();
  }
}

function formatRate(value: number) {
  if (!Number.isFinite(value)) {
    return '0.00%';
  }
  return `${value.toFixed(2)}%`;
}

function formatUserLabel(name: string, id: number) {
  return `${name || '未知用户'}[${id || '-'}]`;
}

function formatTime(value: string) {
  if (!value) {
    return '-';
  }
  return new Date(value).toLocaleString('zh-CN', { hour12: false });
}

function showAllData() {
  tablePageSize.value = stats.records.length;
}

onMounted(() => {
  void loadTodayStats();
});
</script>

<style scoped>
.stats-page {
  padding: 28px 32px;
  background: #131320;
  min-height: 100%;
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

/* 主面板卡片：深色 */
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

/* 按钮组 */
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

.range-form :deep(.el-range-editor) {
  background: #2a2a3c;
  border-color: rgba(255, 255, 255, 0.08);
}

.range-form :deep(.el-range-editor .el-range-input) {
  color: #c0c0d0;
}

.range-form :deep(.el-range-editor .el-range-separator) {
  color: #8b8b9e;
}

/* 骨架屏深色 */
.stats-panel :deep(.el-skeleton__item) {
  background: #2a2a3c;
}

/* Alert 深色 */
.stats-panel :deep(.el-alert) {
  background: #2a2a3c;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.stats-panel :deep(.el-alert--error .el-alert__title) {
  color: #f87171;
}

/* Empty 深色 */
.stats-panel :deep(.el-empty__description) {
  color: #8b8b9e;
}

/* ============== KPI 卡片行 ============== */
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
  font-size: 36px;
}

/* 第二档：成功/失败，两栏并列 */
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

/* Hero 卡片：渐变背景 */
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
  padding: 20px 24px;
}

.summary-card--hero .summary-label {
  color: rgba(255, 255, 255, 0.8);
}

.summary-card--hero strong {
  color: #fff;
  font-size: 28px;
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

.summary-card--success {
  border-left: 4px solid #10b981;
}

.summary-card--success strong {
  color: #34d399;
}

.summary-card--fail {
  border-left: 4px solid #ef4444;
}

.summary-card--fail strong {
  color: #f87171;
}

/* ============== 图表/详情卡片 ============== */
.status-message,
.chart-card,
.detail-card {
  margin-bottom: 20px;
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

/* ============== 扇形图区域 ============== */
.pie-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
  align-items: center;
  justify-content: center;
  padding: 24px 8px;
  position: relative;
}

.pie-svg {
  display: block;
  max-width: 100%;
  height: auto;
  filter: drop-shadow(0 16px 32px rgba(0, 0, 0, 0.4));
}

/* ============== 扇形图启动动画 ============== */
.pie-enter {
  animation: pieContainerEnter 0.8s cubic-bezier(0.34, 1.2, 0.64, 1) both;
  transform-origin: center;
}

.pie-slice {
  cursor: pointer;
  animation: pieSliceFadeIn 0.5s ease-out both;
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1), stroke 0.2s ease, filter 0.2s ease;
}

.pie-slice:hover {
  filter: brightness(1.12);
}

@keyframes pieContainerEnter {
  0% {
    opacity: 0;
    transform: scale(0.4) rotate(-90deg);
  }
  60% {
    opacity: 1;
  }
  100% {
    opacity: 1;
    transform: scale(1) rotate(0deg);
  }
}

@keyframes pieSliceFadeIn {
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

.pie-label {
  font-size: 14px;
  font-weight: 600;
  fill: #fff;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.5);
  pointer-events: none;
}

.pie-legend {
  display: grid;
  gap: 12px;
  min-width: 240px;
}

.legend-item {
  display: grid;
  grid-template-columns: 16px 1fr auto auto;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  background: #1e1e2d;
  cursor: default;
  transition: all 0.2s ease;
}

.legend-item.is-hovered {
  border-color: rgba(96, 165, 250, 0.4);
  box-shadow: 0 4px 12px rgba(96, 165, 250, 0.15);
  transform: translateX(4px);
}

.legend-item.is-filtered {
  border-color: rgba(96, 165, 250, 0.45);
  box-shadow: 0 6px 18px rgba(96, 165, 250, 0.18);
}

.legend-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.legend-label {
  font-weight: 600;
  color: #e4e4e7;
  font-size: 14px;
}

.legend-count {
  color: #8b8b9e;
  font-size: 13px;
}

.legend-rate {
  color: #a0a0b0;
  font-size: 13px;
  min-width: 60px;
  text-align: right;
  font-weight: 600;
}

.pie-popover {
  position: absolute;
  top: 16px;
  right: 16px;
  display: grid;
  gap: 6px;
  min-width: 200px;
  max-width: 260px;
  padding: 16px 18px;
  border: 1px solid rgba(96, 165, 250, 0.25);
  border-radius: 16px;
  background: rgba(30, 30, 45, 0.95);
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(16px);
}

.pie-popover strong {
  font-size: 16px;
  font-weight: 700;
}

.detail-caption {
  color: #8b8b9e;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
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

.detail-card :deep(.el-table--enable-row-transition .el-table__body td.el-table__cell) {
  transition: background-color 0.2s ease;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  margin-top: 12px;
}

.info-text {
  color: #8b8b9e;
  font-size: 13px;
}

@media (max-width: 768px) {
  .stats-page {
    padding: 16px;
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
    margin-bottom: 12px;
  }

  .summary-row--secondary {
    grid-template-columns: 1fr;
  }

  .summary-row--primary .summary-card :deep(.el-card__body) {
    padding: 20px 24px;
  }

  .summary-row--primary .summary-card strong {
    font-size: 28px;
  }

  .current-view-badge {
    display: block;
    margin: 8px 0 0;
  }

  .pie-wrapper {
    gap: 16px;
    padding: 12px 0;
  }

  .pie-legend {
    width: 100%;
  }

  .pie-popover {
    position: static;
    width: 100%;
    max-width: none;
  }

  .section-header h2 {
    font-size: 16px;
  }
}

/* ============== 可折叠卷帘（原始数据卡片） ============== */
.collapsible-card :deep(.el-card__header) {
  padding: 14px 20px;
}

.collapsible-header {
  cursor: pointer;
  user-select: none;
}

.collapsible-header :deep(.el-button) {
  color: #60a5fa;
  border-color: rgba(96, 165, 250, 0.3);
}

.collapsible-header :deep(.el-button:hover) {
  background: rgba(96, 165, 250, 0.1);
  border-color: rgba(96, 165, 250, 0.5);
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

/* ============== 退出原因统计卡片 ============== */
.exit-reason-stats {
  padding: 8px 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 14px;
  font-weight: 600;
  color: #e4e4e7;
}

.stat-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  min-width: 30px;
  text-align: center;
  color: #fff;
}

.stat-badge.success {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}

.stat-badge.info {
  background: rgba(96, 165, 250, 0.2);
  color: #60a5fa;
}

.stat-badge.warning {
  background: rgba(251, 191, 36, 0.2);
  color: #fbbf24;
}

.stat-badge.danger {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
}

.stat-badge.secondary {
  background: rgba(139, 139, 158, 0.2);
  color: #a0a0b0;
}

.progress-bar {
  height: 8px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.06);
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.progress-fill.success {
  background: linear-gradient(90deg, #10b981, #34d399);
}

.progress-fill.info {
  background: linear-gradient(90deg, #60a5fa, #93c5fd);
}

.progress-fill.warning {
  background: linear-gradient(90deg, #fbbf24, #fcd34d);
}

.progress-fill.danger {
  background: linear-gradient(90deg, #ef4444, #f87171);
}

.progress-fill.secondary {
  background: linear-gradient(90deg, #8b8b9e, #a0a0b0);
}

.stat-footer {
  font-size: 12px;
  color: #8b8b9e;
  text-align: right;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
