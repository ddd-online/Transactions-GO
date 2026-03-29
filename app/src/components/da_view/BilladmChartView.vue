<template>
  <div class="chart-view">
    <!-- 图表头部 -->
    <div class="chart-view-header">
      <h3 class="chart-view-title">{{ title }}</h3>
      <a-button type="text" @click="showEnlarged = true">
        <template #icon>
          <ArrowsAltOutlined/>
        </template>
      </a-button>
    </div>

    <!-- 图表内容 -->
    <div class="chart-view-content">
      <div class="chart-wrapper">
        <div class="chart-container">
          <BilladmChart
              v-if="data.length > 0"
              :data="data"
              x-field="time"
              y-field="amount"
              series-field="type"
              :title="title"
          />
          <a-empty v-else description="暂无数据" />
        </div>
      </div>
    </div>

    <!-- 放大图表弹窗 -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showEnlarged" class="chart-enlarge-overlay" @click.self="showEnlarged = false">
          <div class="chart-enlarge-modal">
            <div class="chart-enlarge-header">
              <h3 class="chart-enlarge-title">{{ title }}</h3>
              <a-button type="text" @click="showEnlarged = false">
                <template #icon>
                  <CloseOutlined/>
                </template>
              </a-button>
            </div>
            <div class="chart-enlarge-content">
              <div class="chart-wrapper">
                <div class="chart-container">
                  <BilladmChart
                      v-if="data.length > 0"
                      :data="data"
                      x-field="time"
                      y-field="amount"
                      series-field="type"
                      :title="title"
                  />
                  <a-empty v-else description="暂无数据" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {ArrowsAltOutlined, CloseOutlined} from '@ant-design/icons-vue'
import BilladmChart from '@/components/da_view/BilladmChart.vue'
import type {TimeSeriesData} from '@/backend/chart'

interface Props {
  title: string
  data: TimeSeriesData[]
}

defineProps<Props>()

const showEnlarged = ref(false)
</script>

<style scoped>
.chart-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--billadm-color-major-background, #fff);
  border-radius: 8px;
  overflow: hidden;
}

.chart-view-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--billadm-color-border, #f0f0f0);
  flex-shrink: 0;
}

.chart-view-title {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--billadm-color-text, #333);
}

.chart-view-content {
  flex: 1;
  padding: 16px;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-wrapper {
  position: relative;
  width: 80%;
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

/* 放大弹窗样式 */
.chart-enlarge-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.chart-enlarge-modal {
  width: 66vw;
  height: 80vh;
  background-color: #fff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
}

.chart-enlarge-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.chart-enlarge-title {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.chart-enlarge-content {
  flex: 1;
  padding: 20px;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 过渡动画 */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .chart-enlarge-modal,
.modal-leave-active .chart-enlarge-modal {
  transition: transform 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .chart-enlarge-modal,
.modal-leave-to .chart-enlarge-modal {
  transform: scale(0.95);
}
</style>
