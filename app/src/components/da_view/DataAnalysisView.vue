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

    <!-- 主内容区：左侧边栏 + 右侧图表显示 -->
    <div class="da-main">
      <!-- 左侧图表列表 -->
      <div class="da-sidebar">
        <billadm-chart-list :chart-configs="KEEP_CHART_CONFIGS" @select="onChartSelect" />
      </div>

      <!-- 右侧图表显示 -->
      <div class="da-content">
        <billadm-chart-view v-if="selectedChart" :title="selectedChart.title" :data="selectedChart.data" />
        <a-empty v-else description="请选择图表" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import BilladmTimeRangePicker from '@/components/common/BilladmTimeRangePicker.vue'
import BilladmChartList from '@/components/da_view/BilladmChartList.vue'
import BilladmChartView from '@/components/da_view/BilladmChartView.vue'
import { useLedgerStore } from '@/stores/ledgerStore.ts'
import { useTrQueryConditionStore } from '@/stores/trQueryConditionStore.ts'
import { convertToUnixTimeRange } from '@/backend/timerange.ts'
import { getTrOnCondition } from '@/backend/functions.ts'
import { buildLineChartData, KEEP_CHART_CONFIGS, type ChartConfig, type TimeSeriesData } from '@/backend/chart'
import type { TransactionRecord } from '@/types/billadm'

const ledgerStore = useLedgerStore()
const trQueryConditionStore = useTrQueryConditionStore()

interface ChartInstance {
  title: string
  data: TimeSeriesData[]
}

const selectedChart = ref<ChartInstance | null>(null)

// 查询交易记录
const queryTrs = async (): Promise<TransactionRecord[]> => {
  if (!ledgerStore.currentLedgerId) return []
  const trCondition = {
    ledgerId: ledgerStore.currentLedgerId,
    tsRange: trQueryConditionStore.timeRange
      ? convertToUnixTimeRange(trQueryConditionStore.timeRange)
      : undefined,
  }
  const result = await getTrOnCondition(trCondition)
  return result.items || []
}

// 加载图表数据
const loadChartData = async (config: ChartConfig): Promise<ChartInstance> => {
  const data = await queryTrs()
  const chartData = buildLineChartData(data, {
    granularity: config.granularity,
    lineDisplayTypes: config.lineDisplayTypes,
    includeOutlier: config.includeOutlier,
  })
  return {
    title: config.title,
    data: chartData,
  }
}

// 缓存所有图表数据
const chartDataCache = ref<Map<string, ChartInstance>>(new Map())

const loadAllChartData = async () => {
  const promises = KEEP_CHART_CONFIGS.map(async (config) => {
    const chartInstance = await loadChartData(config)
    chartDataCache.value.set(config.title, chartInstance)
  })
  await Promise.all(promises)

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
  display: flex;
  gap: 12px;
  min-height: 0;
  overflow: hidden;
}

.da-sidebar {
  width: 200px;
  flex-shrink: 0;
  background-color: var(--billadm-color-minor-background, #f5f5f5);
  border: 1px solid var(--billadm-color-border, #e8e8e8);
  border-radius: 8px;
  overflow-y: auto;
}

.da-content {
  flex: 1;
  min-width: 0;
  background-color: var(--billadm-color-major-background, #fff);
  border: 1px solid var(--billadm-color-border, #e8e8e8);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.04);
}
</style>
