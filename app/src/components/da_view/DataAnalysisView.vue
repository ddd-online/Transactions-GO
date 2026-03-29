<template>
  <a-layout style="height: 100%">
    <a-layout-header class="layout-header">
      <div class="toolbar-left">
        <BilladmTimeRangePicker
            v-model:time-range="trQueryConditionStore.timeRange"
            v-model:time-range-type="trQueryConditionStore.timeRangeType"
        />
      </div>
      <div class="toolbar-center">
      </div>
      <div class="toolbar-right">
      </div>
    </a-layout-header>
    <a-layout-content :style="contentStyle">
      <billadm-chart-display/>
    </a-layout-content>
  </a-layout>
</template>

<script setup lang="ts">
import {type CSSProperties} from 'vue';
import BilladmChartDisplay from "@/components/da_view/BilladmChartDisplay.vue";
import {useTrQueryConditionStore} from "@/stores/trQueryConditionStore.ts";
import {useCssVariables} from "@/backend/css.ts";

const {majorBgColor} = useCssVariables();

const contentStyle: CSSProperties = {
  backgroundColor: majorBgColor.value,
  "overflow-y": "auto",
  "margin-bottom": "auto"
};

const trQueryConditionStore = useTrQueryConditionStore();
</script>

<style scoped>
.layout-header {
  height: auto;
  background-color: var(--billadm-color-major-background);
  padding: 0 0 16px 0;
  display: flex;
  align-items: start;
  justify-content: center;
}

.toolbar-left {
  margin-right: auto;
  display: flex;
  gap: 8px;
  align-items: center;
}

.toolbar-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  gap: 8px;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}
</style>
