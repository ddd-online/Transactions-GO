<template>
  <div class="statistics-footer">
    <div class="statistics-footer-item">
      <a-typography-text class="statistics-footer-item-label">
        收入
      </a-typography-text>
      <a-typography-text class="statistics-footer-item-value income">
        {{ centsToYuan(appDataStore.income) }}
      </a-typography-text>
    </div>
    <div class="statistics-footer-item">
      <a-typography-text class="statistics-footer-item-label">
        支出
      </a-typography-text>
      <a-typography-text class="statistics-footer-item-value expense">
        {{ centsToYuan(appDataStore.expense) }}
      </a-typography-text>
    </div>
    <div class="statistics-footer-item">
      <a-typography-text class="statistics-footer-item-label">
        转账
      </a-typography-text>
      <a-typography-text class="statistics-footer-item-value transfer">
        {{ centsToYuan(appDataStore.transfer) }}
      </a-typography-text>
    </div>
  </div>
</template>

<script setup lang="ts">
import {centsToYuan} from "@/backend/functions.ts";
import {useAppDataStore} from "@/stores/appDataStore.ts";

const appDataStore = useAppDataStore();
</script>

<style scoped>
.statistics-footer {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 24px;
}

.statistics-footer-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.statistics-footer-item-label {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-minor);
  margin: 0;
}

.statistics-footer-item-value {
  font-size: var(--billadm-size-text-body);
  font-weight: 600;
  font-family: 'PingFang SC', 'Microsoft YaHei', monospace;
  margin: 0;
}

.statistics-footer-item-value.income {
  color: #52c41a;
}

.statistics-footer-item-value.expense {
  color: #f5222d;
}

.statistics-footer-item-value.transfer {
  color: #1677ff;
}
</style>
