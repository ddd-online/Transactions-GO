<template>
  <div class="chart-view">
    <!-- 图表头部 -->
    <div class="chart-view-header">
      <h3 class="chart-view-title">{{ title }}</h3>
      <a-button type="text" size="small" @click="drawerVisible = true">
        <template #icon>
          <SettingOutlined />
        </template>
        配置
      </a-button>
    </div>

    <!-- 配置抽屉 -->
    <a-drawer v-model:open="drawerVisible" title="图表配置详情" placement="right" width="600">
      <a-descriptions :column="1" size="small" bordered>
        <a-descriptions-item label="图表名称">
          <template v-if="!isPreset">
            <a-input v-model:value="editTitle" style="width: 200px;" />
          </template>
          <template v-else>{{ title }}</template>
        </a-descriptions-item>
        <a-descriptions-item label="时间粒度">
          <template v-if="!isPreset">
            <a-select v-model:value="editGranularity" style="width: 120px;">
              <a-select-option value="year">年度</a-select-option>
              <a-select-option value="month">月度</a-select-option>
            </a-select>
          </template>
          <template v-else>
            <a-tag :color="granularity === 'year' ? 'blue' : 'green'">
              {{ granularity === 'year' ? '年度' : '月度' }}
            </a-tag>
          </template>
        </a-descriptions-item>
        <a-descriptions-item label="曲线数量">{{ lines.length }} 条</a-descriptions-item>
      </a-descriptions>

      <div v-if="!isPreset" style="margin-bottom: 16px;">
        <a-button type="primary" @click="showAddLineModal = true">
          <template #icon>
            <PlusOutlined />
          </template>
          添加曲线
        </a-button>
        <a-button type="primary" style="margin-left: 8px;" @click="handleSave">
          保存修改
        </a-button>
      </div>

      <a-divider orientation="left">曲线详情</a-divider>

      <a-table :data-source="localLines" :pagination="false" size="small">
        <a-table-column title="曲线名称" data-index="label" />
        <a-table-column title="交易类型" data-index="transactionType">
          <template #default="{ text }">
            <a-tag :color="getTypeColor(text)">{{ getTypeLabel(text) }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="包含离群值">
          <template #default="{ record: r }">
            <a-tag :color="r.includeOutlier ? 'orange' : 'green'">
              {{ r.includeOutlier ? '是' : '否' }}
            </a-tag>
          </template>
        </a-table-column>
        <a-table-column title="筛选条件">
          <template #default="{ record }">
            <template v-if="record.conditions && record.conditions.length > 0">
              <div style="display: flex; flex-wrap: wrap; gap: 4px;">
                <a-tag v-for="cond in record.conditions" :key="cond.description" color="purple">
                  {{ cond.category }}
                  <template v-if="cond.tags && cond.tags.length > 0">
                    / {{ cond.tags.join(', ') }}
                  </template>
                  <template v-if="cond.description">
                    / {{ cond.description }}
                  </template>
                </a-tag>
              </div>
            </template>
            <span v-else style="color: #999;">无</span>
          </template>
        </a-table-column>
        <a-table-column v-if="!isPreset" title="操作" width="60">
          <template #default="{ index }">
            <a-button type="text" size="small" danger @click="handleDeleteLine(index)">
              <template #icon>
                <DeleteOutlined />
              </template>
            </a-button>
          </template>
        </a-table-column>
      </a-table>
    </a-drawer>

    <!-- 添加曲线弹窗 -->
    <a-modal v-model:open="showAddLineModal" title="添加曲线" @ok="handleAddLine" :confirm-loading="addLineLoading">
      <a-form :model="newLineForm" layout="vertical">
        <a-form-item label="曲线名称" name="label">
          <a-input v-model:value="newLineForm.label" placeholder="请输入曲线名称" />
        </a-form-item>
        <a-form-item label="交易类型" name="transactionType">
          <a-select v-model:value="newLineForm.transactionType" placeholder="请选择交易类型">
            <a-select-option value="income">收入</a-select-option>
            <a-select-option value="expense">支出</a-select-option>
            <a-select-option value="transfer">转账</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="包含离群值" name="includeOutlier">
          <a-switch v-model:checked="newLineForm.includeOutlier" />
        </a-form-item>
      </a-form>
    </a-modal>

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
import { ref, watch } from 'vue'
import { SettingOutlined, PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import BilladmChart from '@/components/da_view/BilladmChart.vue'
import type { TimeSeriesData, ChartLine } from '@/backend/chart'
import { TransactionTypeToColor, TransactionTypeToLabel } from '@/backend/constant'

interface Props {
  title: string
  data: TimeSeriesData[]
  lines: ChartLine[]
  granularity?: 'year' | 'month'
  isPreset?: boolean
  chartId?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  granularity: 'year',
  isPreset: false,
  chartId: null
})

const emit = defineEmits<{
  (e: 'update', chartId: string, request: { title?: string; granularity?: 'year' | 'month'; lines?: ChartLine[] }): void
  (e: 'addLine', chartId: string, line: ChartLine): void
}>()

const drawerVisible = ref(false)
const editTitle = ref(props.title)
const editGranularity = ref(props.granularity)
const localLines = ref<ChartLine[]>([...props.lines])
const showAddLineModal = ref(false)
const addLineLoading = ref(false)
const newLineForm = ref({
  label: '',
  transactionType: 'income' as string,
  includeOutlier: true,
})

watch(() => props.title, (v) => { editTitle.value = v })
watch(() => props.granularity, (v) => { editGranularity.value = v })
watch(() => props.lines, (v) => { localLines.value = [...v] }, { deep: true })

const getTypeColor = (type: string) => {
  return TransactionTypeToColor.get(type) || '#999'
}

const getTypeLabel = (type: string) => {
  return TransactionTypeToLabel.get(type) || type
}

const handleSave = () => {
  if (!props.chartId) return
  emit('update', props.chartId, {
    title: editTitle.value,
    granularity: editGranularity.value,
    lines: localLines.value,
  })
  drawerVisible.value = false
}

const handleAddLine = () => {
  if (!newLineForm.value.label.trim()) {
    message.error('请输入曲线名称')
    return
  }
  if (!props.chartId) return
  addLineLoading.value = true
  const line: ChartLine = {
    label: newLineForm.value.label,
    transactionType: newLineForm.value.transactionType,
    includeOutlier: newLineForm.value.includeOutlier,
    conditions: [],
  }
  emit('addLine', props.chartId, line)
  showAddLineModal.value = false
  newLineForm.value = { label: '', transactionType: 'income', includeOutlier: true }
  addLineLoading.value = false
}

const handleDeleteLine = (index: number) => {
  localLines.value.splice(index, 1)
}
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
