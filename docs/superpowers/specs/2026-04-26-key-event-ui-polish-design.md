# 关键事件页面 UI 优化设计

## 概述

对关键事件页面进行 UI 优化：工具栏结构调整、年份选择器组件替换、当天日期突出显示。

## 1. 工具栏结构调整

参考交易记录页面 `.tr-toolbar` 布局，采用左-中-右三栏分区，内容居中对齐。

### 布局结构

```html
<div class="key-event-toolbar">
  <div class="key-event-toolbar-left"></div>
  <div class="key-event-toolbar-center">
    <!-- 年份选择器 + 左右箭头按钮 -->
  </div>
  <div class="key-event-toolbar-right"></div>
</div>
```

### CSS 样式

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

## 2. 年份选择器

将 `<a-select>` 替换为 `<a-date-picker picker="year" />`。

### 组件结构

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

### 交互逻辑

- `goToPrevYear()`: `selectedYear.value - 1`，触发 `fetchDatesByYear`
- `goToNextYear()`: `selectedYear.value + 1`，触发 `fetchDatesByYear`
- `a-date-picker` 绑定 `dayjs` 类型，通过 `dayjs().year()` 和 `dayjs().year(val)` 读写
- `selectedYearDayjs` 为 `ref<Dayjs | null>(dayjs())`

## 3. 当天日期突出显示

### 判断逻辑

```ts
const isToday = (year: number, month: number, day: number) => {
  const today = dayjs();
  return year === today.year() && month === today.month() + 1 && day === today.date();
};
```

### 样式

```css
.day-cell--today {
  border: 2px solid var(--billadm-color-primary);
  font-weight: var(--billadm-weight-bold);
}
```

有记录优先显示记录样式（绿色背景），当天标记在有记录的基础上额外叠加边框。

## 修改文件

- `app/src/components/key_event_view/KeyEventView.vue`
