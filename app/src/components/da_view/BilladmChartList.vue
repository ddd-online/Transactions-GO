<template>
  <div class="chart-list">
    <div v-for="chart in chartConfigs" :key="chart.title" class="chart-list-item"
      :class="{ active: selectedTitle === chart.title }" @click="selectChart(chart)">
      <div class="chart-list-item-icon">
        <RiseOutlined style="font-size: 14px" />
      </div>
      <span class="chart-list-item-title">{{ chart.title }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RiseOutlined } from '@ant-design/icons-vue'
import type { ChartConfig } from '@/backend/chart'

interface Props {
  chartConfigs: ChartConfig[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'select', config: ChartConfig): void
}>()

const selectedTitle = ref<string>(props.chartConfigs[0]?.title || '')

const selectChart = (config: ChartConfig) => {
  selectedTitle.value = config.title
  emit('select', config)
}
</script>

<style scoped>
.chart-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
}

.chart-list-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--billadm-color-text-secondary, #666);
}

.chart-list-item:hover {
  background-color: var(--billadm-color-hover-background, #f5f5f5);
}

.chart-list-item.active {
  background-color: #ffffff;
  color: var(--billadm-color-primary, #1677ff);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chart-list-item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.chart-list-item-title {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
