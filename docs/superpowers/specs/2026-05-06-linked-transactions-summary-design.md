# 关联交易统计指标设计

## 概述

在 KeyEvent 弹窗的"关联交易"Tab 中，表格上方新增统计栏，按收入/支出/转账三类分别汇总金额。

## 前端改动

**仅修改** `KeyEventView.vue`。

### 1. 新增 computed — linkedSummary

```ts
const linkedSummary = computed(() => {
  let income = 0, expense = 0, transfer = 0;
  for (const t of linkedTransactions.value) {
    if (t.transactionType === 'income') income += t.price;
    else if (t.transactionType === 'expense') expense += t.price;
    else if (t.transactionType === 'transfer') transfer += t.price;
  }
  return { income, expense, transfer };
});
```

### 2. 统计栏模板

在"关联交易" Tab 中，表格上方插入：

```html
<div class="linked-summary" v-if="linkedTransactions.length > 0">
  <div class="summary-item income">
    <span class="summary-label">收入</span>
    <span class="summary-value">+{{ centsToYuan(linkedSummary.income) }}</span>
  </div>
  <div class="summary-item expense">
    <span class="summary-label">支出</span>
    <span class="summary-value">-{{ centsToYuan(linkedSummary.expense) }}</span>
  </div>
  <div class="summary-item transfer">
    <span class="summary-label">转账</span>
    <span class="summary-value">{{ centsToYuan(linkedSummary.transfer) }}</span>
  </div>
</div>
```

### 3. 样式

```css
.linked-summary {
  display: flex;
  gap: var(--billadm-space-md);
  padding: var(--billadm-space-sm) var(--billadm-space-md);
  margin-bottom: var(--billadm-space-sm);
  background: var(--billadm-color-minor-background);
  border-radius: var(--billadm-radius-md);
}

.summary-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.summary-label {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

.summary-value {
  font-family: var(--billadm-font-mono);
  font-size: var(--billadm-size-text-body);
  font-weight: var(--billadm-weight-semibold);
  font-variant-numeric: tabular-nums;
}

.summary-item.income .summary-value {
  color: var(--billadm-color-income);
}

.summary-item.expense .summary-value {
  color: var(--billadm-color-expense);
}

.summary-item.transfer .summary-value {
  color: var(--billadm-color-transfer);
}
```

### 设计说明

- 仅当有至少一条关联交易时才显示统计栏
- 即使某类型金额为 0 也照常显示（保持三列布局对称）
- 金额使用 `centsToYuan` 格式化，与表格金额列一致
- 颜色使用现有 CSS 变量，与系统其他地方一致
