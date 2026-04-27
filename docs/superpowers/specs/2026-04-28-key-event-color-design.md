# 关键事件颜色标记设计

## 概述

为关键事件增加颜色标记功能，用户在编辑/创建事件时可选择预设颜色，日历格以背景色填充方式呈现。

## 1. 数据模型

### 后端 Model

文件：`kernel/models/key_event.go`

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

- `Color` 最大 20 字符，存储十六进制色值（如 `#E57373`），允许为空
- 空时前端回退至 `--billadm-color-primary`

### 前端 Type

文件：`app/src/types/billadm.d.ts`

```ts
interface KeyEvent {
    id: string
    date: string
    title: string
    content: string
    color: string      // 可为空，hex 色值
    createdAt: number
    updatedAt: number
}
```

## 2. 预设颜色

20 个区分度高的预设色：

| # | 色值 | 用途 |
|---|------|------|
| 1 | `#C73E3A` | 红色 |
| 2 | `#E57373` | 浅红 |
| 3 | `#2D7D46` | 绿色 |
| 4 | `#4CAF50` | 浅绿 |
| 5 | `#5A7FAA` | 蓝灰 |
| 6 | `#64B5F6` | 天蓝 |
| 7 | `#C9A227` | 金色 |
| 8 | `#8B7355` | 棕色 |
| 9 | `#7A5C58` | 赭石 |
| 10 | `#5C7A6A` | 灰绿 |
| 11 | `#9C27B0` | 紫色 |
| 12 | `#BA68C8` | 浅紫 |
| 13 | `#FF9800` | 橙色 |
| 14 | `#FFB74D` | 浅橙 |
| 15 | `#00BCD4` | 青色 |
| 16 | `#4DD0E1` | 浅青 |
| 17 | `#795548` | 深棕 |
| 18 | `#A1887F` | 浅棕 |
| 19 | `#607D8B` | 蓝灰 |
| 20 | `#90A4AE` | 浅蓝灰 |

前端定义常量：

```ts
const EVENT_COLORS = [
  '#C73E3A', '#E57373', '#2D7D46', '#4CAF50',
  '#5A7FAA', '#64B5F6', '#C9A227', '#8B7355',
  '#7A5C58', '#5C7A6A',
  '#9C27B0', '#BA68C8', '#FF9800', '#FFB74D',
  '#00BCD4', '#4DD0E1', '#795548', '#A1887F',
  '#607D8B', '#90A4AE'
]
const DEFAULT_COLOR = '#2D5A27'  // --billadm-color-primary
```

## 3. API 改动

### POST /api/v1/key-events

请求体增加 `color`（可选）字段：

```json
{
    "date": "2026-04-28",
    "title": "会议记录",
    "content": "今天讨论了项目进度...",
    "color": "#C73E3A"
}
```

后端行为：`color` 为空字符串时存空。

### GET /api/v1/key-events/:date

响应体增加 `color` 字段（由 GORM 自动序列化）。

### GET /api/v1/key-events/year/:year

响应体 `KeyEvent[]` 中每个对象包含 `color`（由 GORM 自动序列化）。

## 4. 前端交互

### 日历显示

有记录的日期格背景色使用记录的颜色：

```css
.day-cell--has-record {
    /* 优先使用记录颜色，回退到主色 */
    background-color: var(--event-color, var(--billadm-color-primary));
}
```

文字颜色固定为白色（`var(--billadm-color-text-inverse)`），因为预设颜色中最浅的 `#E57373` 配合白色文字也有足够对比度。若背景色为空则回退到主色。

### 模态框颜色选择器

在标题输入框下方、内容文本框上方，新增颜色选择条：

```
[标题输入框]
[○ ○ ○ ○ ○ ○ ○ ○ ○ ○]  ← 颜色选择条
[多行文本输入框]
```

交互：
- 10 个圆形色块（24px 直径），hover 时边框高亮
- 选中时显示勾选标记（白色 SVG check）
- 点击切换选中状态
- 默认颜色：`#2D5A27`（主色）

实现示例：

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

### 样式

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

## 5. 文件改动

| 文件 | 改动 |
|------|------|
| `kernel/models/key_event.go` | 新增 `Color` 字段 |
| `kernel/api/key_event_controller.go` | 解析 `color` 入参 |
| `kernel/service/key_event_service.go` | `UpsertKeyEvent` 增加 `color` 参数 |
| `app/src/types/billadm.d.ts` | `KeyEvent` 接口新增 `color` |
| `app/src/stores/keyEventStore.ts` | 无需改动（已通过 `queryKeyEventsByYear` 返回完整对象） |
| `app/src/components/key_event_view/KeyEventView.vue` | 颜色选择器 + 日历格背景色 |

## 6. 数据迁移

已有记录的 `color` 字段默认为空，前端按默认主色显示，无缝衔接。