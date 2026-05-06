# 交易记录关联关键事件设计

## 概述

支持将多个消费记录（TransactionRecord）手动关联到某一天的关键事件（KeyEvent），在关键事件详情中集中查看。

## 1. 数据模型

### TransactionRecord 新增字段

`kernel/models/transaction_record.go`：

```go
KeyEventDate string `gorm:"type:varchar(10);comment:关联关键事件日期" json:"key_event_date"`
```

- 类型 VARCHAR(10)，格式 YYYY-MM-DD
- 可为空（空字符串表示未关联）
- GORM AutoMigrate 启动时自动添加列

### 前端类型同步

`app/src/types/billadm.d.ts`：

```ts
interface TransactionRecord {
    // ...existing fields
    key_event_date: string  // 可为空，关联的关键事件日期
}
```

## 2. 关联逻辑

### Link（关联）

1. 根据 `transactionId` 查到交易记录
2. 设置 `key_event_date = date`，保存
3. 检查该 date 是否存在 KeyEvent，不存在则创建空记录（仅 date，title/content/color 为空）

### Unlink（解除关联）

1. 根据 `transactionId` 查到交易记录
2. 设置 `key_event_date = ""`，保存
3. 不删除 KeyEvent

## 3. 后端 API

### POST /api/v1/tr/link

```json
// Request
{ "transaction_id": "uuid-xxx", "date": "2026-05-06" }

// Response
{ "code": 0, "msg": "", "data": "2026-05-06" }
```

自动创建：若该日期 KeyEvent 不存在，service 层自动创建空记录。

### POST /api/v1/tr/unlink

```json
// Request
{ "transaction_id": "uuid-xxx" }

// Response
{ "code": 0, "msg": "", "data": "uuid-xxx" }
```

仅清空 key_event_date，不删除 KeyEvent。

### GET /api/v1/tr/linked/:date

```json
// Response
{ "code": 0, "msg": "", "data": [ /* TransactionRecord[] */ ] }
```

返回完整的 TransactionRecord 列表，跨所有账本。

### 后端文件改动

| 文件 | 改动 |
|------|------|
| `kernel/models/transaction_record.go` | 新增 KeyEventDate 字段 |
| `kernel/dao/transaction_record_dao.go` | 新增 UpdateKeyEventDate、QueryByKeyEventDate 方法 |
| `kernel/service/transaction_record_service.go` | 新增 LinkToKeyEvent、UnlinkFromKeyEvent、QueryLinkedByDate |
| `kernel/api/transaction_record_controller.go` | 新增 link、unlink、listLinked 三个 handler |
| `kernel/api/router.go` | 注册三个新路由 |

## 4. 前端改动

### 4.1 交易记录表格 — 操作列

`TransactionRecordTable.vue`：

- 操作列新增"关联"按钮
- 已关联的记录显示"已关联"按钮（hover tooltip 显示关联日期如"已关联至 2026-05-06"）
- 点击"关联"或"已关联"均 emit `link` 事件
- emit 定义新增: `(e: 'link', record: TransactionRecord): void`

### 4.2 关联日期选择弹窗

`TransactionRecordView.vue`：

- 点击"关联"后弹出日期选择器弹窗（a-date-picker）
- 任意日期可选，未选日期时确认按钮置灰
- 确认后调用 `/api/v1/tr/link`，toast 提示成功
- 已关联的记录弹出时日期选择器预填当前关联日期，支持更改或解除（底部额外"解除关联"按钮）

### 4.3 关键事件弹窗 — Tabs + 关联交易表格

`KeyEventView.vue`：

- 弹窗内新增 Tabs 切换：**详情** | **关联交易 (N)**
- "关联交易" Tab 内嵌简化交易表格，列：

| 账本 | 分类 | 标签 | 描述 | 金额 | 操作 |
|------|------|------|------|------|------|

- 金额列通过颜色区分收入/支出/转账
- 操作列仅"删除"（解除关联），调用 `/api/v1/tr/unlink`
- 无关联时显示空状态："暂无关联交易记录"
- 不分页（关联记录数通常有限），内容区可滚动

### 4.4 Store 和 API 封装

`keyEventStore.ts`：新增 `fetchLinkedTransactions(date: string)` action
`app/src/backend/api/`：新增 link/unlink/fetchLinked API 封装函数

### 前端文件改动

| 文件 | 改动 |
|------|------|
| `TransactionRecordTable.vue` | 操作列新增关联/已关联按钮，emit link 事件 |
| `TransactionRecordView.vue` | 关联日期选择弹窗 |
| `KeyEventView.vue` | 弹窗改为 Tabs，新增关联交易 Tab |
| `keyEventStore.ts` | 新增 fetchLinkedTransactions |
| `functions.ts` / `api/*` | 新增 linkTrToKeyEvent、unlinkTrFromKeyEvent、fetchLinkedTrs |

## 5. 边界情况

| 场景 | 行为 |
|------|------|
| 交易记录不存在 | API 返回错误，前端 toast 提示 |
| 关联时日期为空 | 前端校验，确认按钮置灰 |
| 该日期无 KeyEvent | 自动新建空 KeyEvent，然后关联 |
| 重复关联同一交易到同一日期 | 覆盖，幂等 |
| KeyEvent 被删除 | 关联记录保留旧 key_event_date，查询时 KeyEvent 不存在则返回空 |
| 解除关联后该日期无其他关联 | 不删除 KeyEvent |
| 关联交易 Tab 数据为空 | 显示空状态提示 |
| 跨账本关联 | 支持，Tab 表格显示账本列 |
