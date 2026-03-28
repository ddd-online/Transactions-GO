<template>
  <div class="container">
    <a-typography-text type="secondary">
      收入
    </a-typography-text>
    <a-typography-text>
      {{ centsToYuan(appDataStore.income) }}
    </a-typography-text>
    <a-typography-text type="secondary">
      支出
    </a-typography-text>
    <a-typography-text>
      {{ centsToYuan(appDataStore.expense) }}
    </a-typography-text>
    <a-typography-text type="secondary">
      转账
    </a-typography-text>
    <a-typography-text>
      {{ centsToYuan(appDataStore.transfer) }}
    </a-typography-text>
  </div>
</template>

<script setup lang="ts">
import {centsToYuan} from "@/backend/functions.ts";
import {useAppDataStore} from "@/stores/appDataStore.ts";

const appDataStore = useAppDataStore();
</script>

<style scoped>
.container {
  gap: 8px;
}
</style>