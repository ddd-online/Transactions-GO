# 关联交易统计指标 — 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 KeyEvent 弹窗的"关联交易"Tab 表格上方，显示收入/支出/转账三类金额的汇总统计栏。

**Architecture:** 前端纯计算 — 从已有的 `linkedTransactions` 数组中按 `transactionType` 分组求和，通过 computed 属性实时计算，无需额外 API。

**Tech Stack:** Vue 3 (Composition API), TypeScript

---

### Task 1: 关联交易统计栏

**Files:**
- Modify: `app/src/components/key_event_view/KeyEventView.vue`

- [ ] **Step 1: 新增 computed — linkedSummary**

在 script 中，`linkedCount` 定义之后添加：

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

- [ ] **Step 2: 新增统计栏模板**

在"关联交易" Tab (`a-tab-pane key="linked"`) 中，`<a-spin />` 之前插入：

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

- [ ] **Step 3: 新增样式**

在 `<style scoped>` 末尾添加：

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

- [ ] **Step 4: 类型检查**

```bash
cd app && npx vue-tsc --noEmit
```

Expected: 无类型错误。

- [ ] **Step 5: 提交**

```bash
git add app/src/components/key_event_view/KeyEventView.vue
git commit -m "feat: add income/expense/transfer summary bar to linked transactions tab"
```
