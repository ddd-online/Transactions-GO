# 关键事件页面 UI 优化实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 对关键事件页面进行 UI 优化：工具栏结构调整（左右中三栏布局）、年份选择器替换为年选择器+左右箭头按钮、当天日期突出显示。

**Architecture:** 仅修改 `KeyEventView.vue` 一个文件，参照交易记录页面 `BilladmTimeRangePicker` 的左右箭头按钮样式。

**Tech Stack:** Vue 3 + TypeScript + Ant Design Vue + SCSS

---

## 修改文件

- `app/src/components/key_event_view/KeyEventView.vue`

---

## Task 1: 重构工具栏布局

### 修改 `.key-event-toolbar` 样式

将现有 `.key-event-toolbar` 替换为左-中-右三栏布局：

```css
.key-event-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
  padding-bottom: var(--billadm-space-md);
  border-bottom: 1px solid var(--billadm-color-divider);
}

.key-event-toolbar-left,
.key-event-toolbar-right {
  flex: 1;
}

.key-event-toolbar-center {
  display: flex;
  align-items: center;
  gap: 8px;
}
```

### 新增模板结构

```html
<div class="key-event-toolbar">
  <div class="key-event-toolbar-left"></div>
  <div class="key-event-toolbar-center">
    <!-- 年份选择器 + 左右箭头按钮 -->
  </div>
  <div class="key-event-toolbar-right"></div>
</div>
```

---

## Task 2: 替换年份选择器为年选择器+左右箭头

### 新增 Data 和方法

```ts
import { LeftOutlined, RightOutlined } from '@ant-design/icons-vue';
import type { Dayjs } from 'dayjs';

// 将 selectedYear 从 number ref 改为 Dayjs ref
const selectedYearDayjs = ref<Dayjs>(dayjs());

// 年份值（用于 a-date-picker）
const selectedYear = computed(() => selectedYearDayjs.value.year());

// 上一年
const goToPrevYear = () => {
  selectedYearDayjs.value = dayjs().year(selectedYearDayjs.value.year() - 1);
  keyEventStore.fetchDatesByYear(selectedYearDayjs.value.year());
};

// 下一年
const goToNextYear = () => {
  selectedYearDayjs.value = dayjs().year(selectedYearDayjs.value.year() + 1);
  keyEventStore.fetchDatesByYear(selectedYearDayjs.value.year());
};

// 年选择器变化
const onYearChange = (val: Dayjs | null) => {
  if (val) {
    selectedYearDayjs.value = val;
    keyEventStore.fetchDatesByYear(val.year());
  }
};
```

### 替换模板中年份选择器

将 Task 1 中的 `key-event-toolbar-center` 内容替换为：

```html
<div class="key-event-toolbar-center">
  <a-button type="text" @click="goToPrevYear">
    <template #icon><LeftOutlined /></template>
  </a-button>
  <a-date-picker
    v-model:value="selectedYearDayjs"
    picker="year"
    :style="{ width: '120px' }"
    @change="onYearChange"
  />
  <a-button type="text" @click="goToNextYear">
    <template #icon><RightOutlined /></template>
  </a-button>
</div>
```

同时删除原有的 `yearOptions` computed 和 `onYearChange` 旧的 `(val: any) => onYearChange(val as number)` 处理。

---

## Task 3: 当天日期突出显示

### 新增 `isToday` 判断函数

```ts
const isToday = (year: number, month: number, day: number) => {
  const today = dayjs();
  return year === today.year() && month === today.month() + 1 && day === today.date();
};
```

### 修改日历格子模板

在 `day-cell` 的 `:class` 绑定中追加当天判断：

```html
:class="{ 'day-cell--has-record': hasRecord(selectedYear, month, day), 'day-cell--today': isToday(selectedYear, month, day) }"
```

### 新增当天日期样式

```css
.day-cell--today {
  border: 2px solid var(--billadm-color-primary);
  font-weight: var(--billadm-weight-bold);
}
```

---

## 自检清单

| Spec 要求 | 实现位置 |
|-----------|---------|
| 工具栏左-中-右布局 | Task 1 |
| 年选择器+左右箭头按钮 | Task 2 |
| 当天日期突出显示 | Task 3 |
