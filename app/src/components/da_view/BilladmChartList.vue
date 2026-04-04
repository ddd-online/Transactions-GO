<template>
  <div class="chart-list">
    <!-- 新增图表按钮 -->
    <div class="chart-list-add">
      <a-button type="dashed" block @click="showCreateModal = true">
        <template #icon>
          <PlusOutlined />
        </template>
        新增图表
      </a-button>
    </div>

    <!-- 预设图表 -->
    <div class="chart-list-section">
      <div class="chart-list-section-title">预设图表</div>
      <div v-for="chart in presetCharts" :key="chart.title" class="chart-list-item"
        :class="{ active: selectedId === 'preset_' + chart.title }" @click="selectChart(chart, true)">
        <div class="chart-list-item-icon">
          <RiseOutlined style="font-size: 14px" />
        </div>
        <span class="chart-list-item-title">{{ chart.title }}</span>
      </div>
    </div>

    <!-- 自定义图表 -->
    <div v-if="customCharts.length > 0" class="chart-list-section">
      <div class="chart-list-section-title">自定义图表</div>
      <div v-for="chart in customCharts" :key="chart.chartId" class="chart-list-item"
        :class="{ active: selectedId === chart.chartId }" @click="selectChart(chart, false)">
        <div class="chart-list-item-icon">
          <LineChartOutlined style="font-size: 14px" />
        </div>
        <span class="chart-list-item-title">{{ chart.title }}</span>
        <div class="chart-list-item-actions" @click.stop>
          <a-button type="text" size="small" danger @click="handleDelete(chart)">
            <template #icon>
              <DeleteOutlined />
            </template>
          </a-button>
        </div>
      </div>
    </div>

    <!-- 新增图表弹窗 -->
    <a-modal v-model:open="showCreateModal" title="新增图表" @ok="handleCreate" :confirm-loading="createLoading">
      <a-form :model="createForm" layout="vertical">
        <a-form-item label="图表名称" name="title">
          <a-input v-model:value="createForm.title" placeholder="请输入图表名称" />
        </a-form-item>
        <a-form-item label="时间粒度" name="granularity">
          <a-select v-model:value="createForm.granularity" placeholder="请选择时间粒度">
            <a-select-option value="year">年度</a-select-option>
            <a-select-option value="month">月度</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RiseOutlined, LineChartOutlined, PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { ChartConfig } from '@/backend/chart'
import { KEEP_CHART_CONFIGS } from '@/backend/chart'
import type { ChartDto } from '@/backend/api/chart'
import { deleteChart as deleteChartApi } from '@/backend/api/chart'

interface Props {
  customCharts: ChartDto[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'select', config: ChartConfig | ChartDto, isPreset: boolean): void
  (e: 'create', request: { title: string; granularity: 'year' | 'month' }): void
  (e: 'delete', chartId: string): void
  (e: 'refresh'): void
}>()

const presetCharts = KEEP_CHART_CONFIGS
const selectedId = ref<string>('')
const showCreateModal = ref(false)
const createLoading = ref(false)
const createForm = ref<{ title: string; granularity: 'year' | 'month' }>({
  title: '',
  granularity: 'year'
})

const selectChart = (config: ChartConfig | ChartDto, isPreset: boolean) => {
  if (isPreset) {
    selectedId.value = 'preset_' + config.title
  } else {
    selectedId.value = (config as ChartDto).chartId
  }
  emit('select', config, isPreset)
}

const handleCreate = async () => {
  if (!createForm.value.title.trim()) {
    message.error('请输入图表名称')
    return
  }
  createLoading.value = true
  try {
    emit('create', { title: createForm.value.title, granularity: createForm.value.granularity })
    showCreateModal.value = false
    createForm.value = { title: '', granularity: 'year' }
  } finally {
    createLoading.value = false
  }
}

const handleDelete = async (chart: ChartDto) => {
  try {
    await deleteChartApi(chart.chartId)
    message.success('删除成功')
    emit('delete', chart.chartId)
  } catch (error) {
    message.error('删除失败')
  }
}
</script>

<style scoped>
.chart-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
}

.chart-list-add {
  margin-bottom: 8px;
}

.chart-list-section {
  margin-top: 8px;
}

.chart-list-section-title {
  font-size: 12px;
  color: var(--billadm-color-text-minor);
  padding: 4px 12px;
  margin-bottom: 4px;
}

.chart-list-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--billadm-color-text-minor);
}

.chart-list-item:hover {
  background-color: var(--billadm-color-hover-bg, #e4e4e4);
}

.chart-list-item.active {
  background-color: var(--billadm-color-primary);
  color: #ffffff;
}

.chart-list-item.active .chart-list-item-icon {
  color: #ffffff;
}

.chart-list-item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.chart-list-item-title {
  flex: 1;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chart-list-item-actions {
  display: none;
}

.chart-list-item:hover .chart-list-item-actions {
  display: block;
}

.chart-list-item.active .chart-list-item-actions {
  display: block;
}

.chart-list-item.active .chart-list-item-actions :deep(.ant-btn) {
  color: #ffffff;
}
</style>
