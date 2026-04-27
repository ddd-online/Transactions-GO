# 关键事件标题 + Hover 预览 设计

## 概述

为关键事件功能增加标题字段，支持用户手动输入或自动从内容首行提取，并在日历日期格上通过 Tooltip 显示标题预览。

## 1. 数据模型

### 后端 Model

文件：`kernel/models/key_event.go`

```go
type KeyEvent struct {
    ID        string `gorm:"primaryKey;comment:事件UUID" json:"id"`
    Date      string `gorm:"uniqueIndex;not null;comment:日期 YYYY-MM-DD" json:"date"`
    Title     string `gorm:"type:varchar(200);comment:标题" json:"title"`
    Content   string `gorm:"type:text;comment:事件内容" json:"content"`
    CreatedAt int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"createdAt"`
    UpdatedAt int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updatedAt"`
}
```

- `Title` 最大 200 字符，允许为空
- 提取策略：取 content 首行作为 fallback（在前端保存时处理）

### 前端 Type

文件：`app/src/types/billadm.d.ts`

```ts
interface KeyEvent {
    id: string
    date: string
    title: string       // 可为空
    content: string
    createdAt: number
    updatedAt: number
}
```

## 2. API 改动

### POST /api/v1/key-events

请求体增加 `title`（可选）字段：

```json
{
    "date": "2026-04-27",
    "title": "会议记录",
    "content": "今天讨论了项目进度..."
}
```

后端行为：
- `title` 为空字符串时存空
- `title` 超长时截断至 200 字符

### GET /api/v1/key-events/:date

响应体增加 `title` 字段：

```json
{
    "id": "uuid",
    "date": "2026-04-27",
    "title": "会议记录",
    "content": "今天讨论了项目进度...",
    "createdAt": 1714212000,
    "updatedAt": 1714212000
}
```

## 3. 前端交互

### 日历 Hover Tooltip

- 有记录的日期格 hover 时显示 title（纯文字 tooltip）
- title 为空时显示 "无标题"
- 无记录的日期格 hover 无反应
- 使用 Ant Design `a-tooltip` 组件

### 模态框

在 textarea 上方新增标题输入框：

```
[标题输入框 - maxlength=200]
[多行文本输入框 - placeholder="记录今天发生的事情..."]
```

交互：
- 新建时：标题为空
- 编辑时：回填已有标题
- 保存时：若标题为空，自动从 content 首行提取（取第一行，超长截断200字）

## 4. 文件改动

| 文件 | 改动 |
|------|------|
| `kernel/models/key_event.go` | 新增 Title 字段 |
| `kernel/api/key_event_controller.go` | 解析 title 入参，响应包含 title |
| `kernel/service/key_event_service.go` | Upsert 逻辑透传 title |
| `app/src/types/billadm.d.ts` | KeyEvent 接口新增 title |
| `app/src/stores/keyEventStore.ts` | fetchEventByDate 返回 title |
| `app/src/components/key_event_view/KeyEventView.vue` | Tooltip + 模态框标题输入 |

## 5. 实现细节

### 标题自动提取逻辑

前端保存时处理：

```ts
function extractTitle(content: string): string {
    const firstLine = content.split('\n')[0].trim()
    return firstLine.length > 200 ? firstLine.slice(0, 200) : firstLine
}
```

### Tooltip 显示逻辑

```ts
const getTooltipTitle = (dateStr: string): string | null => {
    if (!keyEventStore.hasRecord(dateStr)) return null
    const title = keyEventStore.getTitle(dateStr)
    return title || '无标题'
}
```

### 响应式布局

模态框中标题输入框宽度跟随容器，样式与设计系统一致。

## 6. 已有记录的数据迁移

已有记录 `title` 字段默认为空字符串，前端按"自动提取首行"逻辑显示。