<template>
  <div class="bottom-bar">
    <div class="bottom-bar-left">
      <span class="bottom-bar-brand">Billadm</span>
    </div>
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
  padding: 0 16px;
  background-color: var(--billadm-color-minor-background);
  border-top: 1px solid var(--billadm-color-divider);
}

.bottom-bar > * {
  -webkit-app-region: no-drag;
}

.bottom-bar-left {
  display: flex;
  align-items: center;
}

.bottom-bar-brand {
  font-family: var(--billadm-font-display);
  font-size: var(--billadm-size-text-caption);
  font-weight: 500;
  color: var(--billadm-color-text-secondary);
  letter-spacing: 0.04em;
}

.bottom-bar-right {
  display: flex;
  align-items: center;
}
</style>
