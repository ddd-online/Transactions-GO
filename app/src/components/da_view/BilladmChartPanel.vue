<template>
  <a-card :title="title" :body-style="bodyCss" hoverable>
    <template v-if="data.length > 0">
      <BilladmChart
          :data="data"
          x-field="time"
          y-field="amount"
          series-field="type"
          :title="title"
      />
    </template>
    <template v-else>
      <a-empty description="暂无数据"/>
    </template>
  </a-card>
</template>

<script setup lang="ts">
import {type CSSProperties} from 'vue';
import BilladmChart from "@/components/da_view/BilladmChart.vue";
import type {TimeSeriesData} from "@/backend/chart";

interface Props {
  title: string
  data: TimeSeriesData[]
}

defineProps<Props>();

const bodyCss: CSSProperties = {
  height: '350px',
  padding: '16px',
}
</script>
