# LedgerView 悬浮按钮实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 LedgerView.vue 的页面标题区改为悬浮按钮，与 TransactionRecordView.vue 保持一致的交互方式。

**Architecture:** 删除 `.view-header` 区域，添加 `a-float-button` 悬浮按钮，移除页面标题但保留空状态和弹窗功能。

**Tech Stack:** Vue 3, Ant Design Vue (a-float-button, PlusOutlined), TypeScript

---

## 文件清单

**修改文件:**
- `app/src/components/ledger_view/LedgerView.vue`

---

## 任务 1: 删除页面标题区

**文件:**
- 修改: `app/src/components/ledger_view/LedgerView.vue:3-15`

- [ ] **Step 1: 删除 .view-header 区域**

删除以下代码（第 3-15 行）:
```html
<!-- 页面标题区 -->
<header class="view-header">
  <div class="view-header-left">
    <h1 class="view-title">账本</h1>
    <span class="view-count">{{ ledgerStore.ledgers.length }} 个账本</span>
  </div>
  <button class="create-btn" @click="openCreateModal">
    <svg class="create-btn-icon" viewBox="0 0 20 20" fill="none">
      <path d="M10 4v12M4 10h12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
    </svg>
    <span>新建账本</span>
  </button>
</header>
```

---

## 任务 2: 添加悬浮按钮

**文件:**
- 修改: `app/src/components/ledger_view/LedgerView.vue` (在空状态之后、弹窗之前添加)

- [ ] **Step 1: 添加悬浮按钮**

在 `</div>` (空状态结束) 之后、`<!-- 新建/编辑账本弹窗 -->` 之前添加:
```html
<!-- 悬浮按钮 -->
<a-float-button type="primary" class="float-primary" @click="openCreateModal">
  <template #icon>
    <PlusOutlined />
  </template>
</a-float-button>
```

- [ ] **Step 2: 导入 PlusOutlined 图标**

在 script setup 的 import 语句中添加:
```typescript
import { PlusOutlined } from "@ant-design/icons-vue";
```

---

## 任务 3: 更新空状态文案

**文件:**
- 修改: `app/src/components/ledger_view/LedgerView.vue:91`

- [ ] **Step 1: 更新空状态描述文字**

将:
```html
<p class="empty-state-desc">点击右上角按钮创建你的第一个账本</p>
```
改为:
```html
<p class="empty-state-desc">点击下方按钮创建你的第一个账本</p>
```

---

## 任务 4: 添加悬浮按钮样式

**文件:**
- 修改: `app/src/components/ledger_view/LedgerView.vue` (在 `<style scoped>` 中添加)

- [ ] **Step 1: 添加 .float-primary 样式**

在 `.ledger-form :deep()` 样式之后添加:
```css
/* ========== Floating Button ========== */
.float-primary {
  right: 48px;
  bottom: 80px;
}
```

- [ ] **Step 2: 删除旧的 .create-btn 相关样式**

删除以下样式定义:
- `.create-btn` (第 213-242 行)
- `.create-btn-icon` (第 228-231 行)
- `.create-btn:hover` (第 233-237 行)
- `.create-btn:active` (第 239-242 行)

- [ ] **Step 3: 删除 .view-header 相关样式**

删除以下样式定义:
- `.view-header` (第 185-192 行)
- `.view-header-left` (第 194-198 行)
- `.view-title` (第 200-205 行)
- `.view-count` (第 207-210 行)

---

## 任务 5: 验证

- [ ] **Step 1: 检查代码变更完整性**

确认以下内容:
1. `.view-header` 区域已删除
2. 悬浮按钮已添加在空状态之后
3. `PlusOutlined` 已导入
4. 空状态文案已更新为"点击下方按钮"
5. `.float-primary` 样式已添加
6. 旧的 `.create-btn` 和 `.view-header` 相关样式已删除

---

## 提交

```bash
git add app/src/components/ledger_view/LedgerView.vue
git commit -m "refactor(ui): replace header button with floating button in LedgerView"
```
