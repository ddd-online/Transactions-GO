# 关键事件标题功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为关键事件增加标题字段，支持手动输入或自动从内容首行提取，日历格 hover 显示标题 tooltip。

**Architecture:** 后端在 KeyEvent 模型中新增 Title 列（varchar 200），API 层透传；前端在 Pinia store 缓存 title，KeyEventView 组件新增标题输入框和 hover tooltip。

**Tech Stack:** Go (GORM) + Vue 3 + TypeScript + Ant Design Vue + Pinia

---

## 文件清单

| 文件 | 改动 |
|------|------|
| `kernel/models/key_event.go` | 新增 `Title string` 字段 |
| `kernel/service/key_event_service.go` | `UpsertKeyEvent` 增加 title 参数，Update 时也更新 title |
| `kernel/api/key_event_controller.go` | 解析请求体 `title` 字段，响应包含 title |
| `app/src/types/billadm.d.ts` | `KeyEvent` 接口新增 `title: string` |
| `app/src/stores/keyEventStore.ts` | `fetchEventByDate` 返回 title；新增 `getTitle(date)` 方法 |
| `app/src/components/key_event_view/KeyEventView.vue` | hover tooltip 显示标题；模态框新增标题输入框 |

---

## Task 1: 后端 — Model + API

### 修改 `kernel/models/key_event.go`

在 `KeyEvent` 结构体中新增 `Title` 字段：

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

### 修改 `kernel/service/key_event_service.go`

**Step 1: 修改接口签名**

`UpsertKeyEvent` 方法签名增加 `title string` 参数：

```go
type KeyEventService interface {
    UpsertKeyEvent(ws *workspace.Workspace, date string, title string, content string) error
    // ... QueryByDate, QueryDatesByYear, DeleteByDate 不变
}
```

**Step 2: 修改实现**

```go
func (s *keyEventServiceImpl) UpsertKeyEvent(ws *workspace.Workspace, date string, title string, content string) error {
    existing, err := s.keyEventDao.QueryByDate(ws, date)
    if err != nil && err != gorm.ErrRecordNotFound {
        return err
    }

    if existing != nil {
        existing.Title = title
        existing.Content = content
        return s.keyEventDao.UpsertKeyEvent(ws, existing)
    }

    event := &models.KeyEvent{
        ID:      util.GetUUID(),
        Date:    date,
        Title:   title,
        Content: content,
    }
    return s.keyEventDao.UpsertKeyEvent(ws, event)
}
```

### 修改 `kernel/api/key_event_controller.go`

**Step 1: 修改 `upsertKeyEvent`**

解析请求体中新增的 `title` 字段，并传入 service：

```go
func upsertKeyEvent(c *gin.Context) {
    ret := models.NewResult()
    defer c.JSON(http.StatusOK, ret)

    ws := workspace.Manager.OpenedWorkspace()
    if ws == nil {
        ret.Code = -1
        ret.Msg = workspace.ErrOpenedWorkspaceNotFound
        return
    }

    arg, ok := JsonArg(c, ret)
    if !ok {
        return
    }

    date, ok := arg["date"].(string)
    if !ok {
        ret.Code = -1
        ret.Msg = "date is required"
        return
    }

    title, _ := arg["title"].(string)
    content, _ := arg["content"].(string)

    if err := service.GetKeyEventService().UpsertKeyEvent(ws, date, title, content); err != nil {
        ret.Code = -1
        ret.Msg = err.Error()
        return
    }

    ret.Data = date
}
```

`getKeyEvent` 和 `deleteKeyEvent` 无需改动（响应体由 GORM 自动包含 title 字段）。

---

## Task 2: 前端 — Type + Store

### 修改 `app/src/types/billadm.d.ts`

`KeyEvent` 接口新增 `title` 字段：

```ts
export interface KeyEvent {
    id: string
    date: string
    title: string       // 可为空
    content: string
    createdAt: number
    updatedAt: number
}
```

### 修改 `app/src/stores/keyEventStore.ts`

**Step 1: 新增 `titles` 缓存**

```ts
const titles = ref(new Map<string, string>());
```

**Step 2: 修改 `fetchDatesByYear`**

获取日期列表时，同时获取每个日期对应的 title（或在 fetchEventByDate 后缓存）。最简单的方式是在 `fetchDatesByYear` 中改为获取完整事件列表：

```ts
const fetchDatesByYear = async (year: number) => {
    try {
        const events = await queryKeyEventsByYear(year)  // 需要新增这个 API
        datesWithRecords.value = new Set(events.map(e => e.date))
        titles.value = new Map(events.map(e => [e.date, e.title]))
        currentYear.value = year
    } catch (error) {
        NotificationUtil.error('查询关键事件失败', `${error}`)
    }
}
```

**Step 3: 新增 `getTitle` 方法**

```ts
const getTitle = (date: string): string => {
    return titles.value.get(date) || ''
}
```

**Step 4: 修改 `fetchEventByDate`**

获取详情后同时更新 title 缓存：

```ts
const fetchEventByDate = async (date: string): Promise<KeyEvent | null> => {
    try {
        const event = await queryKeyEventByDate(date)
        if (event) {
            titles.value.set(date, event.title)
        }
        return event
    } catch (error) {
        return null
    }
}
```

**Step 5: 修改 `saveEvent`**

```ts
const saveEvent = async (date: string, title: string, content: string): Promise<void> => {
    try {
        await saveKeyEvent(date, title, content)
        datesWithRecords.value.add(date)
        titles.value.set(date, title)
        NotificationUtil.success('保存成功')
    } catch (error) {
        NotificationUtil.error('保存失败', `${error}`)
        throw error
    }
}
```

对应的 `saveKeyEvent` API 函数需要修改为传入 title 参数。

**Step 6: 修改 `deleteEvent`**

```ts
const deleteEvent = async (date: string): Promise<void> => {
    try {
        await deleteKeyEvent(date)
        datesWithRecords.value.delete(date)
        titles.value.delete(date)
        NotificationUtil.success('删除成功')
    } catch (error) {
        NotificationUtil.error('删除失败', `${error}`)
        throw error
    }
}
```

### 修改 `app/src/backend/api/key-event.ts`

**Step 1: 新增获取年度事件列表 API**

```ts
export async function queryKeyEventsByYear(year: number): Promise<KeyEvent[]> {
    return api.get<KeyEvent[]>(`/v1/key-events/year/${year}`, '查询关键事件列表')
}
```

**Step 2: 修改 `saveKeyEvent` 签名**

```ts
export async function saveKeyEvent(date: string, title: string, content: string): Promise<string> {
    return api.post<string>('/v1/key-events', { date, title, content }, '保存关键事件')
}
```

需要新增后端 `GET /api/v1/key-events/year/:year` 端点（见 Task 1）。

---

## Task 3: 前端 — UI（KeyEventView）

### 修改 `app/src/components/key_event_view/KeyEventView.vue`

#### 3.1 新增 hover Tooltip

在日期格上包裹 `a-tooltip`，hover 时显示标题：

```html
<a-tooltip :title="getTooltipTitle(formatDate(selectedYear, month, day))" :hidden="!hasRecord(selectedYear, month, day)">
  <div
    class="day-cell"
    :class="{ 'day-cell--has-record': hasRecord(selectedYear, month, day), 'day-cell--today': isToday(selectedYear, month, day) }"
    role="button"
    tabindex="0"
    :aria-label="`${month}月${day}日`"
    @click="onDayClick(selectedYear, month, day)"
    @keydown.enter.prevent="onDayClick(selectedYear, month, day)"
    @keydown.space.prevent="onDayClick(selectedYear, month, day)"
  >
    {{ day }}
  </div>
</a-tooltip>
```

#### 3.2 新增 Tooltip 逻辑函数

```ts
const getTooltipTitle = (dateStr: string): string => {
    if (!keyEventStore.hasRecord(dateStr)) return ''
    const title = keyEventStore.getTitle(dateStr)
    return title || '无标题'
}
```

#### 3.3 模态框新增标题输入框

在 textarea 上方新增标题输入框：

```html
<div class="event-modal-content">
  <a-input
    v-model:value="eventTitle"
    placeholder="标题（可选）"
    :maxlength="200"
    class="event-title-input"
  />
  <a-textarea
    v-model:value="eventContent"
    placeholder="记录今天发生的事情..."
    :rows="5"
    :maxlength="5000"
    show-count
  />
</div>
```

#### 3.4 新增响应式变量

```ts
const eventTitle = ref('')
```

#### 3.5 修改 `onDayClick`

加载事件详情时同时设置 title：

```ts
const onDayClick = async (year: number, month: number, day: number) => {
    const dateStr = formatDate(year, month, day)
    selectedDate.value = dateStr

    if (keyEventStore.hasRecord(dateStr)) {
        confirmLoading.value = true
        try {
            const event = await keyEventStore.fetchEventByDate(dateStr)
            if (event) {
                eventTitle.value = event.title
                eventContent.value = event.content
                isEditMode.value = true
            } else {
                eventTitle.value = ''
                eventContent.value = ''
                isEditMode.value = false
            }
        } finally {
            confirmLoading.value = false
        }
    } else {
        eventTitle.value = ''
        eventContent.value = ''
        isEditMode.value = false
    }

    modalVisible.value = true
}
```

#### 3.6 修改 `handleSave`

保存时传入 title，自动从 content 首行提取 fallback：

```ts
const extractTitle = (content: string): string => {
    const firstLine = content.split('\n')[0].trim()
    return firstLine.length > 200 ? firstLine.slice(0, 200) : firstLine
}

const handleSave = async () => {
    if (!eventContent.value.trim()) {
        NotificationUtil.warning('内容不能为空', '请输入事件内容')
        return
    }
    const title = eventTitle.value.trim() || extractTitle(eventContent.value)
    confirmLoading.value = true
    try {
        await keyEventStore.saveEvent(selectedDate.value, title, eventContent.value.trim())
        modalVisible.value = false
    } finally {
        confirmLoading.value = false
    }
}
```

#### 3.7 样式新增

```css
.event-title-input {
  margin-bottom: var(--billadm-space-sm);
}
```

---

## 自检清单

| Spec 要求 | 实现位置 |
|-----------|---------|
| 后端 Title 字段 | Task 1: `kernel/models/key_event.go` |
| API 解析 title | Task 1: `kernel/api/key_event_controller.go` |
| Service 透传 title | Task 1: `kernel/service/key_event_service.go` |
| 前端 KeyEvent 接口 | Task 2: `app/src/types/billadm.d.ts` |
| Store 缓存 title | Task 2: `app/src/stores/keyEventStore.ts` |
| 日历 hover tooltip | Task 3: `KeyEventView.vue` 3.1-3.2 |
| 模态框标题输入框 | Task 3: `KeyEventView.vue` 3.3-3.7 |
| 自动提取首行作 fallback | Task 3: `KeyEventView.vue` 3.6 |