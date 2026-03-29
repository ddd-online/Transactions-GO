<template>
  <div class="da-view">
    <!-- 工具栏 -->
    <div class="da-toolbar">
      <div class="da-toolbar-left">
        <BilladmTimeRangePicker
            v-model:time-range="trQueryConditionStore.timeRange"
            v-model:time-range-type="trQueryConditionStore.timeRangeType"
        />
      </div>
      <div class="da-toolbar-right">
        <billadm-ledger-select/>
      </div>
    </div>

    <!-- 图表展示 -->
    <div class="da-content">
      <billadm-chart-display/>
    </div>
  </div>
</template>

<script setup lang="ts">
import BilladmChartDisplay from "@/components/da_view/BilladmChartDisplay.vue";
import {useTrQueryConditionStore} from "@/stores/trQueryConditionStore.ts";

const trQueryConditionStore = useTrQueryConditionStore();
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

.da-content {
  flex: 1;
  overflow-y: auto;
}
</style>
