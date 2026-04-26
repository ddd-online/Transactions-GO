# 关键事件 (Key Event) 页面设计

## 概述

新增"关键事件"功能页面，用于记录某一天发生的事情（文本形式）。用户可以通过日历视图快速浏览和编辑每日事件。

## 功能需求

- 顶部工具栏包含年份选择器（仅可选择年份）
- 内容区域显示全年日历（12 个月 × 每周 7 天网格）
- 有记录的日期格子以主题绿色背景突出显示
- 点击日期弹出详情弹窗，支持查看、编辑、保存、删除

## 数据模型

### 新增数据表 `tbl_billadm_key_event`

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | VARCHAR(36) PK | UUID，主键 |
| `date` | VARCHAR(10) NOT NULL UNIQUE | 日期，格式 `YYYY-MM-DD`，唯一索引 |
| `content` | TEXT | 事件内容文本 |
| `created_at` | BIGINT | 创建时间戳 |
| `updated_at` | BIGINT | 更新时间戳 |

## API 接口

### 路由注册位置
`kernel/api/router.go` — 在 `ServeAPI` 函数的 `v1` 路由组下添加 `/key-events` 路由组。

### 接口清单

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/v1/key-events/dates/{year}` | 获取指定年份有记录的所有日期列表 |
| `GET` | `/api/v1/key-events/{date}` | 获取单日详情 |
| `POST` | `/api/v1/key-events` | Upsert 创建或更新（body: `{ date, content }`） |
| `DELETE` | `/api/v1/key-events/{date}` | 删除指定日期的记录 |

### 响应格式

```json
// GET /dates/{year}
{ "code": 0, "msg": "", "data": ["2024-01-15", "2024-03-20", ...] }

// GET /{date}
{ "code": 0, "msg": "", "data": { "id": "...", "date": "2024-01-15", "content": "...", "createdAt": 1234567890, "updatedAt": 1234567890 } }

// POST /key-events (Upsert)
{ "code": 0, "msg": "", "data": "2024-01-15" }

// DELETE /key-events/{date}
{ "code": 0, "msg": "", "data": null }
```

## 后端实现

### 文件结构

```
kernel/
├── models/
│   └── key_event.go           # 模型定义
├── dao/
│   └── key_event_dao.go       # DAO 层
├── service/
│   └── key_event_service.go   # 业务逻辑（含 Upsert）
└── api/
    ├── key_event_controller.go # HTTP 控制器
    └── router.go              # 注册路由
```

### 核心逻辑

**Upsert 逻辑（service 层）**
- 根据 `date` 查询记录
- 若存在则更新 `content` 和 `updated_at`
- 若不存在则创建新记录

## 前端实现

### 文件结构

```
app/src/
├── components/
│   └── key_event_view/
│       └── KeyEventView.vue   # 页面组件
├── stores/
│   └── keyEventStore.ts       # Pinia store
├── backend/
│   └── api/
│       └── key-event.ts       # API 调用封装
└── router/
    └── router.ts             # 注册路由
```

### 组件设计

**KeyEventView.vue**
- 顶部：年份选择器（a-select），默认当前年份
- 内容区：12 个月份的日历网格，使用 CSS Grid 布局
- 弹窗：详情弹窗，含日期标题、文本框、保存/删除按钮

### 日历渲染逻辑

- 每月按 7 列（周日~周六）渲染网格
- 有记录的日期格子背景色为 `var(--billadm-color-primary)`（主题绿），文字色为白色
- 无记录的日期格子保持默认背景
- 点击日期时：
  - 若有记录：先 GET 详情再弹窗
  - 若无记录：直接弹窗（内容为空）

### 状态管理（keyEventStore）

```ts
// state
records: Record<string, boolean>  // { "2024-01-15": true, ... }  用于快速判断哪天有记录

// actions
fetchDatesByYear(year: number): Promise<string[]>
fetchEventByDate(date: string): Promise<KeyEvent>
saveEvent(date: string, content: string): Promise<void>
deleteEvent(date: string): Promise<void>
```

## UI 风格

遵循现有项目设计系统：
- 使用 CSS 变量定义的主题色 `var(--billadm-color-primary)`
- 使用现有 spacing/typography token
- 暗模式支持（CSS 变量自动切换）
- 卡片/弹窗使用 `var(--billadm-radius-lg)` 圆角

## 路由注册

前端路由：`/key_event_view` → `KeyEventView.vue`

侧边栏导航需同步添加入口（参考 `AppLeftBar.vue`）。
