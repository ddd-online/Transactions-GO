<template>
  <a-card
      :title="title"
      :body-style="bodyCss"
      hoverable>
    <BilladmChart :option="option"/>
  </a-card>
</template>

<script setup lang="ts">
import {computed, type CSSProperties} from 'vue';
import BilladmChart from "@/components/da_view/BilladmChart.vue";
import type {TransactionRecord} from "@/types/billadm";
import {buildLineChart, buildPieChart} from "@/backend/table.ts";

interface ChartOptions {
  granularity?: 'year' | 'month'
  lineDisplayTypes?: string[]
  includeOutlier?: boolean
  transactionType?: string
}

const bodyCss: CSSProperties = {
  aspectRatio: 3 / 2,
  minHeight: 0
}

const props = defineProps<{
  title: string
  data: TransactionRecord[]
  chartType: string
  chartOptions: ChartOptions
}>();

const option = computed(() => {
  switch (props.chartType) {
    case 'Line':
      return buildLineChart(props.data, {
        granularity: props.chartOptions.granularity || 'month',
        lineDisplayTypes: props.chartOptions.lineDisplayTypes || ['income', 'expense', 'transfer'],
        includeOutlier: props.chartOptions.includeOutlier === undefined ? true : props.chartOptions.includeOutlier,
      });
    case 'Pie':
      return buildPieChart(props.data, {transactionType: props.chartOptions.transactionType || 'expense'})
    default:
      return buildPieChart(props.data, {transactionType: props.chartOptions.transactionType || 'expense'})
  }
});
</script>