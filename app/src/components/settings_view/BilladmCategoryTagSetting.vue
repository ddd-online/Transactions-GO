<template>
  <div class="category-tag-setting">
    <!-- Tab切换：收入、支出、转账 -->
    <a-tabs v-model:activeKey="activeTab" class="setting-tabs">
      <a-tab-pane key="expense">
        <template #tab>
          <span :style="{ color: expenseColor }">支出</span>
        </template>
        <div class="tab-content">
          <category-tag-panel transaction-type="expense" :active-color="expenseColor"/>
        </div>
      </a-tab-pane>
      <a-tab-pane key="income">
        <template #tab>
          <span :style="{ color: incomeColor }">收入</span>
        </template>
        <div class="tab-content">
          <category-tag-panel transaction-type="income" :active-color="incomeColor"/>
        </div>
      </a-tab-pane>
      <a-tab-pane key="transfer">
        <template #tab>
          <span :style="{ color: transferColor }">转账</span>
        </template>
        <div class="tab-content">
          <category-tag-panel transaction-type="transfer" :active-color="transferColor"/>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';
import {TransactionTypeToColor} from "@/backend/constant.ts";
import CategoryTagPanel from './CategoryTagPanel.vue';

const activeTab = ref('expense');
const incomeColor = TransactionTypeToColor.get('income') || '#52c41a';
const expenseColor = TransactionTypeToColor.get('expense') || '#f5222d';
const transferColor = TransactionTypeToColor.get('transfer') || '#faad14';
</script>

<style scoped>
.category-tag-setting {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.setting-tabs {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.setting-tabs :deep(.ant-tabs-content-holder) {
  flex: 1;
  overflow: hidden;
}

.setting-tabs :deep(.ant-tabs-content) {
  height: 100%;
}

.setting-tabs :deep(.ant-tabs-tabpane) {
  height: 100%;
}

.tab-content {
  height: 100%;
}
</style>
