<template>
  <v-chart
      :option="option"
      :autoresize="true"
      class="billadm-chart-instance"
  />
</template>

<script setup lang="ts">
import VChart from 'vue-echarts';
import type {EChartsOption} from "echarts";

interface Props {
  option: EChartsOption
}

defineProps<Props>();
</script>

<style scoped>
.billadm-chart-instance {
  width: 100%;
  height: 100%;
}
</style>