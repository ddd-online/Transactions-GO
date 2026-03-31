<template>
  <div class="bottom-bar">
    <div class="bottom-bar-left"></div>
    <div class="bottom-bar-right">
      <billadm-statistics-footer v-if="showStatistics" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import BilladmStatisticsFooter from '@/components/common/BilladmStatisticsFooter.vue'

const route = useRoute()

const showStatistics = computed(() => {
  return route.path === '/tr_view' || route.path === '/da_view'
})
</script>

<style scoped>
.bottom-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  width: 100%;
  padding: 0 12px;
}

.bottom-bar > * {
  -webkit-app-region: no-drag;
}

.bottom-bar-left {
  display: flex;
  align-items: center;
}

.bottom-bar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
