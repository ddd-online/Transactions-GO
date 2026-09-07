<template>
  <a-select
    :value="value"
    class="stats-tag-filter-select"
    @change="emit('change', $event)"
    aria-label="按交易标签筛选"
  >
    <a-select-option value="">全部标签</a-select-option>
    <a-select-option v-for="opt in STOCK_TRADE_TAG_OPTIONS" :key="opt.value" :value="opt.value">
      {{ opt.label }}
    </a-select-option>
  </a-select>
</template>

<script setup lang="ts">
import { STOCK_TRADE_TAG_OPTIONS } from '@/backend/constant'
import type { StockTradeTag } from '@/types/transactions'

defineProps<{
  value: StockTradeTag | ''
}>()

const emit = defineEmits<{
  (e: 'change', value: unknown): void
}>()
</script>

<style scoped lang="scss">
.stats-tag-filter-select {
  width: 128px;
}
</style>
