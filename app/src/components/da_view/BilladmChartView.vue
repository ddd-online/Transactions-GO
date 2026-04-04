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
        <a-descriptions-item label="图表名称">{{ title }}</a-descriptions-item>
        <a-descriptions-item label="时间粒度">
          <a-tag :color="granularity === 'year' ? 'blue' : 'green'">
            {{ granularity === 'year' ? '年度' : '月度' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="曲线数量">{{ lines.length }} 条</a-descriptions-item>
      </a-descriptions>

      <a-divider orientation="left">曲线详情</a-divider>

      <a-table :data-source="lines" :pagination="false" size="small">
        <a-table-column title="曲线名称" data-index="label" />
        <a-table-column title="交易类型" data-index="transactionType">
          <template #default="{ text }">
            <a-tag :color="getTypeColor(text)">{{ getTypeLabel(text) }}</a-tag>
          </template>
        </a-table-column>
        <a-table-column title="包含离群值">
          <template #default="{ record }">
            <a-tag :color="record.includeOutlier ? 'orange' : 'green'">
              {{ record.includeOutlier ? '是' : '否' }}
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
      </a-table>
    </a-drawer>

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
import { ref } from 'vue'
import { SettingOutlined } from '@ant-design/icons-vue'
import BilladmChart from '@/components/da_view/BilladmChart.vue'
import type { TimeSeriesData, ChartLine } from '@/backend/chart'
import { TransactionTypeToColor, TransactionTypeToLabel } from '@/backend/constant'

interface Props {
  title: string
  data: TimeSeriesData[]
  lines: ChartLine[]
  granularity?: 'year' | 'month'
}

const props = withDefaults(defineProps<Props>(), {
  granularity: 'year'
})

const drawerVisible = ref(false)

const getTypeColor = (type: string) => {
  return TransactionTypeToColor.get(type) || '#999'
}

const getTypeLabel = (type: string) => {
  return TransactionTypeToLabel.get(type) || type
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
