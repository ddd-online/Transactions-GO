<template>
  <div class="da-view">
    <!-- 工具栏 -->
    <div class="da-toolbar">
      <div class="da-toolbar-left">
        <BilladmTimeRangePicker v-model:time-range="trQueryConditionStore.timeRange"
          v-model:time-range-type="trQueryConditionStore.timeRangeType" />
      </div>
      <div class="da-toolbar-right">
        <billadm-ledger-select />
      </div>
    </div>

    <!-- 主内容区 -->
    <a-card class="da-main" :body-style="{ padding: '0', display: 'flex', height: '100%' }">
      <!-- 左侧图表列表 -->
      <div class="da-sidebar">
        <billadm-chart-list :chart-configs="KEEP_CHART_CONFIGS" @select="onChartSelect" />
      </div>

      <!-- 右侧图表显示 -->
      <div class="da-content">
        <billadm-chart-view v-if="selectedChart" :title="selectedChart.title" :data="selectedChart.data" />
        <a-empty v-else description="请选择图表" />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import BilladmTimeRangePicker from '@/components/common/BilladmTimeRangePicker.vue'
import BilladmChartList from '@/components/da_view/BilladmChartList.vue'
import BilladmChartView from '@/components/da_view/BilladmChartView.vue'
import { useLedgerStore } from '@/stores/ledgerStore.ts'
import { useTrQueryConditionStore } from '@/stores/trQueryConditionStore.ts'
import { useAppDataStore } from '@/stores/appDataStore.ts'
import { convertToUnixTimeRange } from '@/backend/timerange.ts'
import { getTrOnCondition } from '@/backend/functions.ts'
import { buildLineChartData, KEEP_CHART_CONFIGS, type ChartConfig, type TimeSeriesData } from '@/backend/chart'
import type { TransactionRecord, TrStatistics } from '@/types/billadm'

const ledgerStore = useLedgerStore()
const trQueryConditionStore = useTrQueryConditionStore()
const appDataStore = useAppDataStore()

interface ChartInstance {
  title: string
  data: TimeSeriesData[]
}

const selectedChart = ref<ChartInstance | null>(null)

// 查询交易记录
const queryTrs = async (conditions: import('@/types/billadm').TrQueryConditionItem[] = []): Promise<{ items: TransactionRecord[], trStatistics: TrStatistics | null }> => {
  if (!ledgerStore.currentLedgerId) return { items: [], trStatistics: null }
  const trCondition = {
    ledgerId: ledgerStore.currentLedgerId,
    tsRange: trQueryConditionStore.timeRange
      ? convertToUnixTimeRange(trQueryConditionStore.timeRange)
      : undefined,
    items: conditions,
  }
  const result = await getTrOnCondition(trCondition)
  return { items: result.items || [], trStatistics: result.trStatistics || null }
}

// 加载图表数据
const loadChartData = async (config: ChartConfig): Promise<{ chartInstance: ChartInstance, trStatistics: TrStatistics | null }> => {
  // 查询所有符合条件的交易记录（不过滤具体条件，由buildLineChartData根据lines中的条件分别过滤）
  const { items, trStatistics } = await queryTrs([])
  const chartData = buildLineChartData(items, {
    granularity: config.granularity,
    lines: config.lines,
  })
  return {
    chartInstance: {
      title: config.title,
      data: chartData,
    },
    trStatistics,
  }
}

// 缓存所有图表数据
const chartDataCache = ref<Map<string, ChartInstance>>(new Map())

const loadAllChartData = async () => {
  let statistics: TrStatistics | null = null
  const promises = KEEP_CHART_CONFIGS.map(async (config) => {
    const { chartInstance, trStatistics } = await loadChartData(config)
    chartDataCache.value.set(config.title, chartInstance)
    if (trStatistics) statistics = trStatistics
  })
  await Promise.all(promises)

  // 更新底部统计信息
  if (statistics) {
    appDataStore.setStatistics(statistics)
  }

  // 初始化选中第一个图表，或更新当前选中图表的数据
  if (!selectedChart.value && KEEP_CHART_CONFIGS.length > 0) {
    const firstConfig = KEEP_CHART_CONFIGS[0]!
    selectedChart.value = chartDataCache.value.get(firstConfig.title) || null
  } else if (selectedChart.value) {
    // 更新已选中图表的数据
    const updatedData = chartDataCache.value.get(selectedChart.value.title)
    if (updatedData) {
      selectedChart.value = updatedData
    }
  }
}

// 图表选择
const onChartSelect = (config: ChartConfig) => {
  selectedChart.value = chartDataCache.value.get(config.title) || null
}

onMounted(() => {
  loadAllChartData()
})

// 监听查询条件或账本变化，重新加载
watch(
  () => [ledgerStore.currentLedgerId, trQueryConditionStore.timeRange],
  () => loadAllChartData(),
  { deep: true }
)
</script>

<style scoped>
.da-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 16px;
  gap: 16px;
}

.da-toolbar {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  flex-shrink: 0;
}

.da-toolbar-left {
  display: flex;
  gap: 8px;
}

.da-toolbar-right {
  display: flex;
  gap: 8px;
}

.da-main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.da-sidebar {
  flex: 0 0 200px;
  background-color: var(--billadm-color-minor-background);
  border-right: 1px solid var(--billadm-color-window-border);
  overflow-y: auto;
}

.da-content {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  background-color: var(--billadm-color-major-background);
}
</style>
