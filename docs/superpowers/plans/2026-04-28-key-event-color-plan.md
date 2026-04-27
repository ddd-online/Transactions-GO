# 关键事件颜色标记实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为关键事件增加颜色标记功能，用户可选择 10 种预设颜色，日历格以背景色填充显示。

**Architecture:** 后端在 KeyEvent 模型中新增 Color 列（varchar 20），API 层透传；前端在 KeyEventView 组件中新增颜色选择器，日历格通过 CSS var 实现动态背景色。

**Tech Stack:** Go (GORM) + Vue 3 + TypeScript + Ant Design Vue

---

## 文件清单

| 文件 | 改动 |
|------|------|
| `kernel/models/key_event.go` | 新增 `Color string` 字段 |
| `kernel/service/key_event_service.go` | `UpsertKeyEvent` 增加 `color` 参数 |
| `kernel/api/key_event_controller.go` | 解析 `color` 入参 |
| `app/src/types/billadm.d.ts` | `KeyEvent` 接口新增 `color: string` |
| `app/src/components/key_event_view/KeyEventView.vue` | 颜色选择器 + 日历格背景色 |

---

## Task 1: 后端 — Model + API

### 修改 `kernel/models/key_event.go`

在 `KeyEvent` 结构体中新增 `Color` 字段：

```go
type KeyEvent struct {
    ID        string `gorm:"primaryKey;comment:事件UUID" json:"id"`
    Date      string `gorm:"uniqueIndex;not null;comment:日期 YYYY-MM-DD" json:"date"`
    Title     string `gorm:"type:varchar(200);comment:标题" json:"title"`
    Content   string `gorm:"type:text;comment:事件内容" json:"content"`
    Color     string `gorm:"type:varchar(20);comment:颜色标记 hex" json:"color"`
    CreatedAt int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"createdAt"`
    UpdatedAt int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updatedAt"`
}
```

GORM AutoMigrate 会自动创建 `color` 列。

### 修改 `kernel/service/key_event_service.go`

**Step 1: 修改接口签名**

`UpsertKeyEvent` 方法签名增加 `color string` 参数：

```go
type KeyEventService interface {
    UpsertKeyEvent(ws *workspace.Workspace, date string, title string, content string, color string) error
    // ... QueryByDate, QueryDatesByYear, DeleteByDate 不变
}
```

**Step 2: 修改实现**

```go
func (s *keyEventServiceImpl) UpsertKeyEvent(ws *workspace.Workspace, date string, title string, content string, color string) error {
    if len(title) > 200 {
        title = title[:200]
    }

    existing, err := s.keyEventDao.QueryByDate(ws, date)
    if err != nil && err != gorm.ErrRecordNotFound {
        return err
    }

    if existing != nil {
        existing.Title = title
        existing.Content = content
        existing.Color = color
        return s.keyEventDao.UpsertKeyEvent(ws, existing)
    }

    event := &models.KeyEvent{
        ID:      util.GetUUID(),
        Date:    date,
        Title:   title,
        Content: content,
        Color:   color,
    }
    return s.keyEventDao.UpsertKeyEvent(ws, event)
}
```

### 修改 `kernel/api/key_event_controller.go`

修改 `upsertKeyEvent` 函数，解析 `color` 字段：

```go
title, _ := arg["title"].(string)
content, _ := arg["content"].(string)
color, _ := arg["color"].(string)

if err := service.GetKeyEventService().UpsertKeyEvent(ws, date, title, content, color); err != nil {
```

`getKeyEvent` 和 `listKeyEventsByYear` 无需改动（GORM 自动序列化 `color` 字段）。

---

## Task 2: 前端 — Type

### 修改 `app/src/types/billadm.d.ts`

在 `KeyEvent` 接口中新增 `color` 字段：

```ts
export interface KeyEvent {
    id: string
    date: string
    title: string
    content: string
    color: string      // 可为空，hex 色值
    createdAt: number
    updatedAt: number
}
```

---

## Task 3: 前端 — UI（KeyEventView.vue）

### 修改 `app/src/components/key_event_view/KeyEventView.vue`

#### 3.1 定义预设颜色常量

在 `<script setup>` 中添加：

```ts
const EVENT_COLORS = [
  '#C73E3A', '#E57373', '#2D7D46', '#4CAF50',
  '#5A7FAA', '#64B5F6', '#C9A227', '#8B7355',
  '#7A5C58', '#5C7A6A'
]

const DEFAULT_COLOR = '#2D5A27'  // --billadm-color-primary
```

#### 3.2 新增 `eventColor` ref

```ts
const eventColor = ref(DEFAULT_COLOR)
```

#### 3.3 修改 `onDayClick` 加载颜色

当获取已有事件时，同步设置 `eventColor`：

```ts
if (event) {
    eventTitle.value = event.title
    eventContent.value = event.content
    eventColor.value = event.color || DEFAULT_COLOR
    isEditMode.value = true
} else {
    eventTitle.value = ''
    eventContent.value = ''
    eventColor.value = DEFAULT_COLOR
    isEditMode.value = false
}
```

新建时也重置为默认色：

```ts
eventColor.value = DEFAULT_COLOR
```

#### 3.4 修改 `handleSave` 传入颜色

```ts
await keyEventStore.saveEvent(selectedDate.value, title, eventContent.value.trim(), eventColor.value)
```

注意：`saveEvent` 目前签名是 `saveEvent(date, title, content)` — 需要新增 `color` 参数。查看 `keyEventStore.ts` 中的 `saveEvent` 签名并更新。

#### 3.5 在模态框中添加颜色选择器

在标题输入框和文本框之间添加颜色选择条：

```html
<div class="color-picker">
  <div
    v-for="c in EVENT_COLORS"
    :key="c"
    class="color-swatch"
    :class="{ 'is-selected': eventColor === c }"
    :style="{ backgroundColor: c }"
    @click="eventColor = c"
  >
    <CheckOutlined v-if="eventColor === c" class="check-icon" />
  </div>
</div>
```

#### 3.6 修改日历格样式使用动态颜色

修改 `.day-cell--has-record` 样式，使用 CSS 变量：

```css
.day-cell--has-record {
    background-color: var(--event-color, var(--billadm-color-primary));
    color: var(--billadm-color-text-inverse);
    font-weight: var(--billadm-weight-semibold);
}
```

在日期格上通过 `style` 绑定设置 `--event-color`：

```html
<div
  class="day-cell"
  :class="{ 'day-cell--has-record': hasRecord(selectedYear, month, day), 'day-cell--today': isToday(selectedYear, month, day) }"
  :style="getDayCellStyle(selectedYear, month, day)"
  ...
>
```

`getDayCellStyle` 函数实现：

```ts
const getDayCellStyle = (year: number, month: number, day: number) => {
    const dateStr = formatDate(year, month, day)
    if (!keyEventStore.hasRecord(dateStr)) return {}
    const color = keyEventStore.getColor(dateStr) || DEFAULT_COLOR
    return { '--event-color': color } as Record<string, string>
}
```

同时 `keyEventStore` 需要新增 `getColor(date)` 方法，与 `getTitle` 类似。

#### 3.7 添加颜色选择器和样式

```css
.color-picker {
  display: flex;
  gap: 6px;
  margin-bottom: var(--billadm-space-sm);
  flex-wrap: wrap;
}

.color-swatch {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color var(--billadm-transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
}

.color-swatch:hover {
  border-color: rgba(0, 0, 0, 0.3);
}

.color-swatch.is-selected {
  border-color: #000;
}

.check-icon {
  color: #fff;
  font-size: 12px;
}
```

#### 3.8 更新 store 的 saveEvent 签名

修改 `app/src/stores/keyEventStore.ts` 的 `saveEvent` 签名，增加 `color` 参数并传递给 API：

```ts
const saveEvent = async (date: string, title: string, content: string, color: string): Promise<void> => {
    try {
        await saveKeyEvent(date, title, content, color)
        datesWithRecords.value.add(date)
        titles.value.set(date, title)
        // 颜色也存入缓存
        colors.value.set(date, color)
        NotificationUtil.success('保存成功')
    } catch (error) {
        NotificationUtil.error('保存失败', `${error}`)
        throw error
    }
}
```

同样 `queryKeyEventsByYear` 返回的事件对象已包含 `color`，需要在 `fetchDatesByYear` 中同时缓存颜色：

```ts
const colors = ref(new Map<string, string>())

// fetchDatesByYear 中：
colors.value = new Map(events.map(e => [e.date, e.color]))

// 新增 getColor
const getColor = (date: string): string => {
    return colors.value.get(date) || ''
}
```

更新 `deleteEvent` 也清理颜色缓存：`colors.value.delete(date)`

#### 3.9 更新 API saveKeyEvent 签名

`app/src/backend/api/key-event.ts` 中 `saveKeyEvent` 需要增加 `color` 参数：

```ts
export async function saveKeyEvent(date: string, title: string, content: string, color: string): Promise<string> {
    return api.post<string>('/v1/key-events', { date, title, content, color }, '保存关键事件')
}
```

---

## 自检清单

| Spec 要求 | 实现位置 |
|-----------|---------|
| 后端 Color 字段 | Task 1: `kernel/models/key_event.go` |
| Service 透传 color | Task 1: `kernel/service/key_event_service.go` |
| API 解析 color | Task 1: `kernel/api/key_event_controller.go` |
| 前端 KeyEvent 接口 | Task 2: `app/src/types/billadm.d.ts` |
| 颜色选择器 UI | Task 3: KeyEventView.vue 3.5, 3.7 |
| 日历格背景色 | Task 3: KeyEventView.vue 3.6 |
| store 缓存 color | Task 3: keyEventStore.ts 3.8 |
| API 传递 color | Task 3: key-event.ts 3.9 |