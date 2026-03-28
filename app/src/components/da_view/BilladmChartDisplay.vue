<template>
  <a-row :gutter="16" style="width: 100%;padding: 16px">
    <a-col
        v-for="(chart, index) in chartsWithData"
        :key="index"
        :span="colSpan"
        style="margin-bottom: 16px"
    >
      <billadm-fullscreen v-model="chart.isFullscreen" :dblclick="true">
        <billadm-chart-panel
            :title="chart.title"
            :data="chart.data"
            :chart-type="chart.chartType"
            :chart-options="chart.chartOptions"
        />
      </billadm-fullscreen>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, watch} from 'vue';
import BilladmChartPanel from '@/components/da_view/BilladmChartPanel.vue';
import BilladmFullscreen from '@/components/common/BilladmFullScreen.vue';
import {useLedgerStore} from "@/stores/ledgerStore.ts";
import {useTrQueryConditionStore} from "@/stores/trQueryConditionStore.ts";
import type {TransactionRecord, TrQueryConditionItem} from "@/types/billadm";
import {convertToUnixTimeRange} from "@/backend/timerange.ts";
import {getTrOnCondition} from "@/backend/functions.ts";

interface Props {
  columns?: number;
}

const props = withDefaults(defineProps<Props>(), {
  columns: 3,
});

interface ChartOptions {
  granularity?: 'year' | 'month'
  lineDisplayTypes?: string[]
  includeOutlier?: boolean
  transactionType?: string
}

interface ChartConfig {
  title: string
  chartType: 'Line' | 'Pie'
  conditions: TrQueryConditionItem[]
  chartOptions: ChartOptions
}

interface ChartInstance {
  title: string
  chartType: 'Line' | 'Pie'
  data: TransactionRecord[]
  chartOptions: ChartOptions
  isFullscreen: boolean
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

const chartConfigs: ChartConfig[] = [
  {
    title: '月度消费趋势(不含离群值)',
    chartType: 'Line',
    conditions: [],
    chartOptions: {
      granularity: "month",
      includeOutlier: false,
    } as ChartOptions
  },
  {
    title: '年度消费趋势(含离群值)',
    chartType: 'Line',
    conditions: [],
    chartOptions: {
      granularity: "year",
      includeOutlier: true,
    } as ChartOptions
  },
  {
    title: '月度加油开销',
    chartType: 'Line',
    conditions: [{
      transactionType: 'expense',
      category: '交通出行',
      tags: ['油费'],
      tagPolicy: 'all',
      tagNot: false,
      description: ''
    }],
    chartOptions: {
      granularity: "month",
      lineDisplayTypes: ['expense'],
    } as ChartOptions
  },
  {
    title: '年度总收入',
    chartType: 'Line',
    conditions: [{
      transactionType: 'income',
      category: '工资奖金',
      tags: ['工资'],
      tagPolicy: 'all',
      tagNot: false,
      description: '-'
    }, {
      transactionType: 'income',
      category: '工资奖金',
      tags: ['奖金'],
      tagPolicy: 'all',
      tagNot: false,
      description: '年奖金'
    }, {
      transactionType: 'income',
      category: '投资理财',
      tags: [],
      tagPolicy: 'all',
      tagNot: false,
      description: '年分红'
    }],
    chartOptions: {
      granularity: "year",
      lineDisplayTypes: ['income'],
    } as ChartOptions
  },
  {
    title: '年度工资收入',
    chartType: 'Line',
    conditions: [{
      transactionType: 'income',
      category: '工资奖金',
      tags: ['工资'],
      tagPolicy: 'all',
      tagNot: false,
      description: '-'
    }],
    chartOptions: {
      granularity: "year",
      lineDisplayTypes: ['income'],
    } as ChartOptions
  },
  {
    title: '年度奖金收入',
    chartType: 'Line',
    conditions: [{
      transactionType: 'income',
      category: '工资奖金',
      tags: ['奖金'],
      tagPolicy: 'all',
      tagNot: false,
      description: '年奖金'
    }],
    chartOptions: {
      granularity: "year",
      lineDisplayTypes: ['income'],
    } as ChartOptions
  }, {
    title: '年度分红收入',
    chartType: 'Line',
    conditions: [{
      transactionType: 'income',
      category: '投资理财',
      tags: [],
      tagPolicy: 'all',
      tagNot: false,
      description: '年分红'
    }],
    chartOptions: {
      granularity: "year",
      lineDisplayTypes: ['income'],
    } as ChartOptions
  },
  {
    title: '消费分布-支出',
    chartType: 'Pie',
    conditions: [],
    chartOptions: {
      transactionType: "expense",
    } as ChartOptions
  },
  {
    title: '消费分布-收入',
    chartType: 'Pie',
    conditions: [],
    chartOptions: {
      transactionType: "income",
    } as ChartOptions
  },
];

const colSpan = computed(() => {
  const cols = props.columns;
  if (cols <= 0 || 24 % cols !== 0) {
    console.warn(`columns 应为 24 的约数（如 1,2,3,4,6,8,12,24），当前值: ${cols}，回退到 2 列`);
    return 12;
  }
  return 24 / cols;
});

// 图表实例列表（带 data）
const chartsWithData = ref<ChartInstance[]>([])

// 加载所有图表数据
const loadChartsData = async () => {
  const promises = chartConfigs.map(async (config) => {
    const data = await queryTrs(config.conditions)
    return {
      title: config.title,
      chartType: config.chartType,
      data,
      chartOptions: config.chartOptions,
      isFullscreen: false,
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