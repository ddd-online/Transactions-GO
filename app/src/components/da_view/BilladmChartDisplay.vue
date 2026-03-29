<template>
  <a-row :gutter="16" style="width: 100%; padding: 16px">
    <a-col
        v-for="(chart, index) in chartsWithData"
        :key="index"
        :span="8"
        style="margin-bottom: 16px"
    >
      <billadm-chart-panel
          :title="chart.title"
          :data="chart.data"
      />
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import {onMounted, ref, watch} from 'vue';
import BilladmChartPanel from '@/components/da_view/BilladmChartPanel.vue';
import {useLedgerStore} from "@/stores/ledgerStore.ts";
import {useTrQueryConditionStore} from "@/stores/trQueryConditionStore.ts";
import type {TransactionRecord, TrQueryConditionItem} from "@/types/billadm";
import {convertToUnixTimeRange} from "@/backend/timerange.ts";
import {getTrOnCondition} from "@/backend/functions.ts";
import {buildLineChartData, KEEP_CHART_CONFIGS, type TimeSeriesData} from "@/backend/chart";

interface ChartInstance {
  title: string
  data: TimeSeriesData[]
}

const ledgerStore = useLedgerStore();
const trQueryConditionStore = useTrQueryConditionStore();

const queryTrs = async (conditions: TrQueryConditionItem[]): Promise<TransactionRecord[]> => {
  if (!ledgerStore.currentLedgerId) return []
  const trCondition = {
    ledgerId: ledgerStore.currentLedgerId,
    tsRange: trQueryConditionStore.timeRange
        ? convertToUnixTimeRange(trQueryConditionStore.timeRange)
        : undefined,
    items: conditions.length > 0 ? conditions : undefined,
  }
  const result = await getTrOnCondition(trCondition)
  return result.items || []
};

// 图表实例列表（带 data）
const chartsWithData = ref<ChartInstance[]>([])

// 加载所有图表数据
const loadChartsData = async () => {
  const promises = KEEP_CHART_CONFIGS.map(async (config) => {
    const data = await queryTrs([]) // 暂时用空条件，实际会根据配置调整
    const chartData = buildLineChartData(data, {
      granularity: config.granularity,
      lineDisplayTypes: config.lineDisplayTypes,
      includeOutlier: config.includeOutlier,
    })
    return {
      title: config.title,
      data: chartData,
    }
  })

  chartsWithData.value = await Promise.all(promises)
}

onMounted(() => {
  loadChartsData()
})

// 监听查询条件或账本变化，重新加载
watch(
    () => [ledgerStore.currentLedgerId, trQueryConditionStore.timeRange],
    () => {
      loadChartsData()
    },
    {deep: true}
)
</script>
