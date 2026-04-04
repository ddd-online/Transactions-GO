<template>
  <div class="chart-view">
    <!-- 图表头部 -->
    <div class="chart-view-header">
      <h3 class="chart-view-title">{{ title }}</h3>
    </div>

    <!-- 图表内容 -->
    <div class="chart-view-content">
      <div class="chart-wrapper">
        <div class="chart-container">
          <BilladmChart v-if="data.length > 0" :data="data" x-field="time" y-field="amount" :title="title"
            :lines="lines" />
          <a-empty v-else description="暂无数据" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import BilladmChart from '@/components/da_view/BilladmChart.vue'
import type { TimeSeriesData, ChartLine } from '@/backend/chart'

interface Props {
  title: string
  data: TimeSeriesData[]
  lines: ChartLine[]
}

defineProps<Props>()
</script>

<style scoped>
.chart-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.chart-view-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid var(--billadm-color-window-border);
  flex-shrink: 0;
}

.chart-view-title {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

.chart-view-content {
  flex: 1;
  padding: 16px;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--billadm-color-major-background);
}

.chart-wrapper {
  position: relative;
  width: 90%;
  aspect-ratio: 16 / 9;
  overflow: hidden;
}

.chart-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
</style>
