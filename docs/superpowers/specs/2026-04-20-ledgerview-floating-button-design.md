# LedgerView 悬浮按钮设计

## 概述

将 LedgerView.vue 的页面标题区改为悬浮按钮，与 TransactionRecordView.vue 保持一致的交互方式。

## 变更内容

### 1. 删除页面标题区

删除 `.view-header` 区域（包含标题"账本"、账本数量、"新建账本"按钮）。

### 2. 添加悬浮按钮

在右下角添加新建账本悬浮按钮：
- 使用 Ant Design `a-float-button` 组件
- 主按钮使用 `PlusOutlined` 图标
- 位置：`right: 48px; bottom: 80px;`
- 绑定 `openCreateModal` 方法

### 3. 更新空状态文案

将空状态描述从"点击右上角按钮"改为"点击下方按钮"。

## 样式对照

| 元素 | TransactionRecordView.vue | LedgerView.vue (变更后) |
|------|-------------------------|------------------------|
| 主悬浮按钮类名 | `.float-primary` | `.float-primary` |
| 主悬浮按钮位置 | `right: 48px; bottom: 80px;` | `right: 48px; bottom: 80px;` |
| 按钮图标 | `PlusOutlined` | `PlusOutlined` |
| 点击事件 | `createTr` | `openCreateModal` |

## 涉及文件

- `app/src/components/ledger_view/LedgerView.vue`
