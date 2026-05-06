# 交易记录关联关键事件 — 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 支持将多条交易记录手动关联到某一天的关键事件，并在关键事件详情弹窗中集中查看和管理。

**Architecture:** 在 TransactionRecord 表新增 `key_event_date` 字段存储关联日期。后端通过三个新 API（link / unlink / listLinked）管理关联关系，link 时若目标日期无 KeyEvent 则自动创建空记录。前端在交易表格操作列增加"关联"按钮，在 KeyEvent 弹窗增加 Tabs 切换展示关联交易。

**Tech Stack:** Go (Gin + GORM), Vue 3 + Pinia + Ant Design Vue, TypeScript

---

## File Structure

```
kernel/
├── models/
│   ├── transaction_record.go          # 新增 KeyEventDate 字段
│   └── dto/
│       └── transaction_record_dto.go  # DTO 新增 KeyEventDate，双向映射
├── dao/
│   └── transaction_record_dao.go      # 新增 UpdateKeyEventDate、QueryByKeyEventDate
├── service/
│   └── transaction_record_service.go  # 新增 LinkToKeyEvent、UnlinkFromKeyEvent、QueryLinkedByDate
├── api/
│   ├── transaction_record_controller.go # 新增三个 handler
│   └── router.go                        # 注册三个新路由

app/src/
├── types/
│   └── billadm.d.ts                     # TransactionRecord 新增 key_event_date，TrForm 新增 key_event_date
├── backend/
│   ├── api/
│   │   └── tr.ts                        # 新增 linkTrToKeyEvent、unlinkTrFromKeyEvent、fetchLinkedTransactions
│   └── dto-utils.ts                    # trFormToTrDto / trDtoToTrForm 映射 key_event_date
├── stores/
│   └── keyEventStore.ts                # 新增 linkedTransactions 状态、fetchLinkedTransactions action
└── components/
    ├── tr_view/
    │   ├── TransactionRecordTable.vue   # 操作列新增"关联"/"已关联"按钮
    │   └── TransactionRecordView.vue    # 处理关联事件，弹出日期选择器弹窗
    └── key_event_view/
        └── KeyEventView.vue            # 弹窗改为 Tabs，新增关联交易 Tab
```

---

### Task 1: Go 模型 — 新增 KeyEventDate 字段

**Files:**
- Modify: `kernel/models/transaction_record.go`

- [ ] **Step 1: 在 TransactionRecord struct 中新增字段**

```go
// kernel/models/transaction_record.go，在 TransactionAt 之后、CreatedAt 之前添加：

	// 关联关键事件日期
	KeyEventDate string `gorm:"type:varchar(10);comment:关联关键事件日期" json:"key_event_date"`

	// 时间信息
	TransactionAt int64 `gorm:"not null;comment:交易时间" json:"transaction_at"`
```

完整 struct 变为：

```go
type TransactionRecord struct {
	TransactionID string `gorm:"primaryKey;comment:交易UUID" json:"transaction_id"`
	LedgerID      string `gorm:"not null;comment:关联账本ID" json:"ledger_id"`

	Price           int64  `gorm:"not null;comment:交易金额" json:"price"`
	TransactionType string `gorm:"not null;comment:交易类型" json:"transaction_type"`

	Category    string `gorm:"not null;comment:分类ID" json:"category"`
	Description string `gorm:"comment:交易描述" json:"description"`

	Flags string `gorm:"comment:标记集" json:"flags"`

	KeyEventDate string `gorm:"type:varchar(10);comment:关联关键事件日期" json:"key_event_date"`

	TransactionAt int64 `gorm:"not null;comment:交易时间" json:"transaction_at"`
	CreatedAt     int64 `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"created_at"`
	UpdatedAt     int64 `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updated_at"`
}
```

- [ ] **Step 2: 构建验证编译通过**

```bash
cd kernel && go build ./...
```

Expected: 编译成功

- [ ] **Step 3: 提交**

```bash
git add kernel/models/transaction_record.go
git commit -m "feat: add KeyEventDate field to TransactionRecord model"
```

---

### Task 2: Go DTO — 新增 KeyEventDate 字段及映射

**Files:**
- Modify: `kernel/models/dto/transaction_record_dto.go`

- [ ] **Step 1: DTO struct 新增字段**

在 `TransactionRecordDto` struct 中，`Outlier` 字段之后添加：

```go
KeyEventDate string `json:"key_event_date"`
```

修改后的 struct：

```go
type TransactionRecordDto struct {
	LedgerID        string   `json:"ledgerId"`
	TransactionID   string   `json:"transactionId"`
	Price           int64    `json:"price"`
	TransactionType string   `json:"transactionType"`
	Category        string   `json:"category"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	TransactionAt   int64    `json:"transactionAt"`
	Outlier         bool     `json:"outlier"`
	KeyEventDate    string   `json:"key_event_date"`
}
```

- [ ] **Step 2: ToTransactionRecord 映射（不含 KeyEventDate，由专用 link API 管理）**

不修改 `ToTransactionRecord` — `KeyEventDate` 通过专用 link/unlink API 设置，不走创建/更新的 DTO 转换。在方法中添加注释说明。

- [ ] **Step 3: FromTransactionRecord 映射**

在 `FromTransactionRecord` 方法末尾添加：

```go
func (dto *TransactionRecordDto) FromTransactionRecord(tr *models.TransactionRecord) {
	dto.LedgerID = tr.LedgerID
	dto.TransactionID = tr.TransactionID
	dto.Price = tr.Price
	dto.TransactionType = tr.TransactionType
	dto.Category = tr.Category
	dto.Description = tr.Description
	dto.Tags = make([]string, 0)
	dto.TransactionAt = tr.TransactionAt
	dto.KeyEventDate = tr.KeyEventDate
	flags := models.TransactionRecordFlags{}
	if err := json.Unmarshal([]byte(tr.Flags), &flags); err == nil {
		dto.Outlier = flags.Outlier
	}
}
```

- [ ] **Step 4: 构建验证**

```bash
cd kernel && go build ./...
```

Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add kernel/models/dto/transaction_record_dto.go
git commit -m "feat: add KeyEventDate to TransactionRecordDto"
```

---

### Task 3: Go DAO — 新增 UpdateKeyEventDate 和 QueryByKeyEventDate

**Files:**
- Modify: `kernel/dao/transaction_record_dao.go`

- [ ] **Step 1: 接口新增两个方法签名**

在 `TransactionRecordDao` interface 末尾添加：

```go
UpdateKeyEventDate(ws *workspace.Workspace, trId string, date string) error
QueryByKeyEventDate(ws *workspace.Workspace, date string) ([]*models.TransactionRecord, error)
```

- [ ] **Step 2: 实现 UpdateKeyEventDate**

在 `transactionRecordDaoImpl` 上添加：

```go
func (t *transactionRecordDaoImpl) UpdateKeyEventDate(ws *workspace.Workspace, trId string, date string) error {
	return ws.GetDb().
		Model(&models.TransactionRecord{}).
		Where("transaction_id = ?", trId).
		Update("key_event_date", date).Error
}
```

- [ ] **Step 3: 实现 QueryByKeyEventDate**

```go
func (t *transactionRecordDaoImpl) QueryByKeyEventDate(ws *workspace.Workspace, date string) ([]*models.TransactionRecord, error) {
	trs := make([]*models.TransactionRecord, 0)
	err := ws.GetDb().
		Where("key_event_date = ?", date).
		Order("transaction_at desc").
		Find(&trs).Error
	return trs, err
}
```

- [ ] **Step 4: 构建验证**

```bash
cd kernel && go build ./...
```

Expected: 编译成功

- [ ] **Step 5: 提交**

```bash
git add kernel/dao/transaction_record_dao.go
git commit -m "feat: add UpdateKeyEventDate and QueryByKeyEventDate to TransactionRecordDao"
```

---

### Task 4: Go Service — 新增 LinkToKeyEvent、UnlinkFromKeyEvent、QueryLinkedByDate

**Files:**
- Modify: `kernel/service/transaction_record_service.go`

- [ ] **Step 1: 新增 imports**

在 import 块中添加：

```go
import (
	"errors"
	"fmt"
	"sync"

	"github.com/billadm/dao"
	"github.com/billadm/models"
	"github.com/billadm/models/dto"
	"github.com/billadm/pkg/operator"
	"github.com/billadm/util"
	"github.com/billadm/workspace"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)
```

（仅新增 `"errors"` 和 `"gorm.io/gorm"`，其他已存在）

- [ ] **Step 2: 接口新增三个方法签名**

在 `TransactionRecordService` interface 末尾添加：

```go
LinkToKeyEvent(ws *workspace.Workspace, trId string, date string) error
UnlinkFromKeyEvent(ws *workspace.Workspace, trId string) error
QueryLinkedByDate(ws *workspace.Workspace, date string) ([]*dto.TransactionRecordDto, error)
```

- [ ] **Step 3: 实现 LinkToKeyEvent**

在文件末尾添加：

```go
func (t *transactionRecordServiceImpl) LinkToKeyEvent(ws *workspace.Workspace, trId string, date string) error {
	logrus.Infof("link transaction %s to key event date %s", trId, date)

	// Update the transaction record's key_event_date
	if err := t.trDao.UpdateKeyEventDate(ws, trId, date); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("transaction not found: %s", trId)
		}
		return fmt.Errorf("update key event date: %w", err)
	}

	// Ensure KeyEvent exists for this date; auto-create if not
	keyEventSvc := GetKeyEventService()
	_, err := keyEventSvc.QueryByDate(ws, date)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check key event: %w", err)
	}
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		if err := keyEventSvc.UpsertKeyEvent(ws, date, "", "", ""); err != nil {
			return fmt.Errorf("auto-create key event: %w", err)
		}
		logrus.Infof("auto-created empty key event for date %s", date)
	}

	logrus.Infof("linked transaction %s to key event date %s", trId, date)
	return nil
}
```

- [ ] **Step 4: 实现 UnlinkFromKeyEvent**

```go
func (t *transactionRecordServiceImpl) UnlinkFromKeyEvent(ws *workspace.Workspace, trId string) error {
	logrus.Infof("unlink transaction %s from key event", trId)

	if err := t.trDao.UpdateKeyEventDate(ws, trId, ""); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("transaction not found: %s", trId)
		}
		return fmt.Errorf("unlink key event date: %w", err)
	}

	logrus.Infof("unlinked transaction %s from key event", trId)
	return nil
}
```

- [ ] **Step 5: 实现 QueryLinkedByDate**

```go
func (t *transactionRecordServiceImpl) QueryLinkedByDate(ws *workspace.Workspace, date string) ([]*dto.TransactionRecordDto, error) {
	logrus.Infof("query linked transactions for date %s", date)

	trs, err := t.trDao.QueryByKeyEventDate(ws, date)
	if err != nil {
		return nil, fmt.Errorf("query by key event date: %w", err)
	}

	// Batch query tags
	trIds := make([]string, len(trs))
	for i, tr := range trs {
		trIds[i] = tr.TransactionID
	}
	tagMap, err := t.trTagDao.QueryTrTagsByTrIds(ws, trIds)
	if err != nil {
		return nil, err
	}

	// Assemble DTOs
	dtos := make([]*dto.TransactionRecordDto, 0, len(trs))
	for _, tr := range trs {
		trDto := &dto.TransactionRecordDto{}
		trDto.FromTransactionRecord(tr)
		if tags, ok := tagMap[tr.TransactionID]; ok {
			for _, tag := range tags {
				trDto.Tags = append(trDto.Tags, tag.Tag)
			}
		}
		dtos = append(dtos, trDto)
	}

	logrus.Infof("query linked transactions for date %s, count: %d", date, len(dtos))
	return dtos, nil
}
```

- [ ] **Step 6: 构建验证**

```bash
cd kernel && go build ./...
```

Expected: 编译成功

- [ ] **Step 7: 提交**

```bash
git add kernel/service/transaction_record_service.go
git commit -m "feat: add LinkToKeyEvent, UnlinkFromKeyEvent, QueryLinkedByDate to TrService"
```

---

### Task 5: Go API — 新增三个 handler + 注册路由

**Files:**
- Modify: `kernel/api/transaction_record_controller.go`
- Modify: `kernel/api/router.go`

- [ ] **Step 1: 新增 handler — linkTransactionToKeyEvent**

在 `kernel/api/transaction_record_controller.go` 末尾添加：

```go
// POST /transactions/link
func linkTransactionToKeyEvent(c *gin.Context) {
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

	trId, _ := arg["transaction_id"].(string)
	date, _ := arg["date"].(string)

	if trId == "" || date == "" {
		ret.Code = -1
		ret.Msg = "transaction_id and date are required"
		return
	}

	if err := service.GetTrService().LinkToKeyEvent(ws, trId, date); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = date
}

// POST /transactions/unlink
func unlinkTransactionFromKeyEvent(c *gin.Context) {
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

	trId, _ := arg["transaction_id"].(string)
	if trId == "" {
		ret.Code = -1
		ret.Msg = "transaction_id is required"
		return
	}

	if err := service.GetTrService().UnlinkFromKeyEvent(ws, trId); err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = trId
}

// GET /transactions/linked/:date
func listLinkedTransactions(c *gin.Context) {
	ret := models.NewResult()
	defer c.JSON(http.StatusOK, ret)

	ws := workspace.Manager.OpenedWorkspace()
	if ws == nil {
		ret.Code = -1
		ret.Msg = workspace.ErrOpenedWorkspaceNotFound
		return
	}

	date := c.Param("date")
	if date == "" {
		ret.Code = -1
		ret.Msg = "date is required"
		return
	}

	dtos, err := service.GetTrService().QueryLinkedByDate(ws, date)
	if err != nil {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	ret.Data = dtos
}
```

- [ ] **Step 2: 注册路由**

在 `kernel/api/router.go` 的 `transactions` group 中添加：

```go
transactions.POST("/link", linkTransactionToKeyEvent)
transactions.POST("/unlink", unlinkTransactionFromKeyEvent)
transactions.GET("/linked/:date", listLinkedTransactions)
```

修改后的 transactions group：

```go
transactions := v1.Group("/transactions")
{
	transactions.POST("/query", queryTransactions)
	transactions.POST("/query-chart-data", queryChartData)
	transactions.POST("/batch", batchCreateTransactions)
	transactions.POST("", createTransaction)
	transactions.DELETE("/:id", deleteTransaction)
	transactions.POST("/link", linkTransactionToKeyEvent)
	transactions.POST("/unlink", unlinkTransactionFromKeyEvent)
	transactions.GET("/linked/:date", listLinkedTransactions)
}
```

- [ ] **Step 3: 构建验证**

```bash
cd kernel && go build ./...
```

Expected: 编译成功

- [ ] **Step 4: 提交**

```bash
git add kernel/api/transaction_record_controller.go kernel/api/router.go
git commit -m "feat: add link/unlink/listLinked API handlers for transaction-keyevent association"
```

---

### Task 6: 前端 TypeScript 类型 — 新增 key_event_date

**Files:**
- Modify: `app/src/types/billadm.d.ts`

- [ ] **Step 1: TransactionRecord 新增字段**

```ts
export interface TransactionRecord {
    ledgerId: string;
    transactionId: string;
    price: number;
    transactionType: string;
    category: string;
    description: string;
    tags: string[];
    transactionAt: number;
    outlier: boolean;
    key_event_date: string;  // 关联的关键事件日期，可为空
}
```

- [ ] **Step 2: TrForm 新增字段**

```ts
export interface TrForm {
    id: string;
    price: string;
    type: string;
    category: string;
    description: string;
    tags: string[];
    flags: string[];
    time: Dayjs;
    key_event_date: string;  // 关联的关键事件日期，可为空
}
```

- [ ] **Step 3: 提交**

```bash
git add app/src/types/billadm.d.ts
git commit -m "feat: add key_event_date to TransactionRecord and TrForm types"
```

---

### Task 7: 前端 DTO 工具 — key_event_date 映射

**Files:**
- Modify: `app/src/backend/dto-utils.ts`

- [ ] **Step 1: trFormToTrDto 添加 key_event_date**

```ts
export function trFormToTrDto(data: TrForm, ledgerId: string = ''): TransactionRecord {
    let tr: TransactionRecord = {
        ledgerId: ledgerId,
        transactionId: data.id,
        price: yuanToCents(data.price),
        transactionType: data.type,
        category: data.category,
        description: data.description,
        tags: data.tags,
        transactionAt: data.time.unix(),
        outlier: false,
        key_event_date: data.key_event_date || '',
    };

    if (data.flags.includes('outlier')) {
        tr.outlier = true;
    }

    return tr;
}
```

- [ ] **Step 2: trDtoToTrForm 添加 key_event_date**

```ts
export function trDtoToTrForm(dto: TransactionRecord): TrForm {
    let trForm: TrForm = {
        id: dto.transactionId,
        price: centsToYuan(dto.price),
        type: dto.transactionType,
        category: dto.category,
        description: dto.description,
        tags: dto.tags,
        flags: [],
        time: dayjs(dto.transactionAt * 1000),
        key_event_date: dto.key_event_date || '',
    };

    if (dto.outlier) {
        trForm.flags.push('outlier');
    }

    return trForm;
}
```

- [ ] **Step 3: 类型检查**

```bash
cd app && npx vue-tsc --noEmit
```

- [ ] **Step 4: 提交**

```bash
git add app/src/backend/dto-utils.ts
git commit -m "feat: map key_event_date in frontend DTO conversion"
```

---

### Task 8: 前端 API 封装 — 新增三个 API 调用

**Files:**
- Modify: `app/src/backend/api/tr.ts`

- [ ] **Step 1: 新增三个导出函数**

在文件末尾添加：

```ts
export async function linkTrToKeyEvent(transactionId: string, date: string): Promise<string> {
    return api.post<string>('/v1/transactions/link', { transaction_id: transactionId, date }, '关联关键事件');
}

export async function unlinkTrFromKeyEvent(transactionId: string): Promise<string> {
    return api.post<string>('/v1/transactions/unlink', { transaction_id: transactionId }, '解除关联');
}

export async function fetchLinkedTransactions(date: string): Promise<TransactionRecord[]> {
    return api.get<TransactionRecord[]>(`/v1/transactions/linked/${date}`, '查询关联交易记录');
}
```

- [ ] **Step 2: 提交**

```bash
git add app/src/backend/api/tr.ts
git commit -m "feat: add linkTrToKeyEvent, unlinkTrFromKeyEvent, fetchLinkedTransactions API wrappers"
```

---

### Task 9: 前端 functions.ts — 新增关联操作封装

**Files:**
- Modify: `app/src/backend/functions.ts`

- [ ] **Step 1: 新增 import**

在文件顶部 import 中添加：

```ts
import { linkTrToKeyEvent, unlinkTrFromKeyEvent, fetchLinkedTransactions } from "@/backend/api/tr.ts";
```

- [ ] **Step 2: 新增三个函数**

在文件末尾添加：

```ts
/**
 * 关联交易记录到关键事件
 */
export async function linkTransactionToKeyEvent(transactionId: string, date: string): Promise<boolean> {
    try {
        await linkTrToKeyEvent(transactionId, date);
        NotificationUtil.success('关联成功');
        return true;
    } catch (error) {
        NotificationUtil.error('关联失败', `${error}`);
        return false;
    }
}

/**
 * 解除交易记录与关键事件的关联
 */
export async function unlinkTransactionFromKeyEvent(transactionId: string): Promise<boolean> {
    try {
        await unlinkTrFromKeyEvent(transactionId);
        NotificationUtil.success('已解除关联');
        return true;
    } catch (error) {
        NotificationUtil.error('解除关联失败', `${error}`);
        return false;
    }
}

/**
 * 获取某日期关联的所有交易记录（跨账本）
 */
export async function getLinkedTransactions(date: string): Promise<TransactionRecord[]> {
    try {
        return await fetchLinkedTransactions(date);
    } catch (error) {
        NotificationUtil.error('查询关联交易失败', `${error}`);
        return [];
    }
}
```

- [ ] **Step 3: 提交**

```bash
git add app/src/backend/functions.ts
git commit -m "feat: add link/unlink/fetchLinked helper functions with error handling"
```

---

### Task 10: 前端 TransactionRecordTable — 操作列新增关联按钮

**Files:**
- Modify: `app/src/components/tr_view/TransactionRecordTable.vue`

- [ ] **Step 1: 操作列模板新增"关联"/"已关联"按钮**

在模板的操作列 slot 中，在"删除"按钮之后添加：

```html
<template v-else-if="column.dataIndex === 'action'">
  <div class="cell-actions">
    <a-button type="text" class="action-btn" @click="handleEdit(record as TransactionRecord)">
      <EditOutlined /> 编辑
    </a-button>
    <a-tooltip v-if="(record as TransactionRecord).key_event_date" :title="'已关联至 ' + (record as TransactionRecord).key_event_date">
      <a-button type="text" class="action-btn" @click="handleLink(record as TransactionRecord)">
        <LinkOutlined /> 已关联
      </a-button>
    </a-tooltip>
    <a-button v-else type="text" class="action-btn" @click="handleLink(record as TransactionRecord)">
      <LinkOutlined /> 关联
    </a-button>
    <a-popconfirm
      title="确认删除此条记录？"
      ok-text="确认"
      @confirm="handleDelete(record as TransactionRecord)"
      :showCancel="false"
    >
      <a-button type="text" class="action-btn danger">
        <DeleteOutlined /> 删除
      </a-button>
    </a-popconfirm>
  </div>
</template>
```

- [ ] **Step 2: 更新 import 和 emit**

在 script 中，将 `EditOutlined, DeleteOutlined` 改为：

```ts
import {EditOutlined, DeleteOutlined, LinkOutlined} from "@ant-design/icons-vue";
```

修改 emit 定义，新增 `link` 事件：

```ts
const emit = defineEmits<{
  (e: 'edit', record: TransactionRecord): void;
  (e: 'delete', record: TransactionRecord): void;
  (e: 'link', record: TransactionRecord): void;
}>();
```

- [ ] **Step 3: 新增 handleLink 方法**

```ts
const handleLink = (record: TransactionRecord) => {
  emit('link', record);
};
```

- [ ] **Step 4: 调整操作列宽度**

```ts
{
  title: '操作',
  dataIndex: 'action',
  width: 200,  // 从 160 增加到 200，容纳新按钮
  align: 'center'
}
```

- [ ] **Step 5: 提交**

```bash
git add app/src/components/tr_view/TransactionRecordTable.vue
git commit -m "feat: add link button to transaction record table action column"
```

---

### Task 11: 前端 TransactionRecordView — 关联日期选择弹窗

**Files:**
- Modify: `app/src/components/tr_view/TransactionRecordView.vue`

- [ ] **Step 1: 新增 import**

在 script 中新增：

```ts
import { linkTransactionToKeyEvent, unlinkTransactionFromKeyEvent } from "@/backend/functions.ts";
```

- [ ] **Step 2: 新增关联弹窗状态**

在 script 的状态区添加：

```ts
// 关联关键事件弹窗
const openLinkModal = ref(false);
const linkingRecord = ref<TransactionRecord | null>(null);
const linkDate = ref<Dayjs>(dayjs());
```

- [ ] **Step 3: 新增关联日期选择弹窗模板**

在模板末尾（`</a-modal>` 筛选弹窗之后、`</div>` 关闭标签之前）添加：

```html
<!-- 关联关键事件弹窗 -->
<a-modal
  v-model:open="openLinkModal"
  title="关联关键事件"
  ok-text="确认关联"
  cancel-text="取消"
  centered
  @ok="confirmLink"
  @cancel="openLinkModal = false"
>
  <a-form>
    <a-form-item label="选择日期">
      <a-date-picker
        v-model:value="linkDate"
        style="width: 100%"
        placeholder="选择要关联的日期"
      />
    </a-form-item>
  </a-form>
  <template v-if="linkingRecord?.key_event_date" #footer>
    <a-button danger @click="handleUnlink">解除关联</a-button>
    <a-button @click="openLinkModal = false">取消</a-button>
    <a-button type="primary" @click="confirmLink">确认关联</a-button>
  </template>
</a-modal>
```

- [ ] **Step 4: 新增 handleLink、confirmLink、handleUnlink 方法**

在 script 中添加：

```ts
const handleLink = (record: TransactionRecord) => {
  linkingRecord.value = record;
  linkDate.value = record.key_event_date ? dayjs(record.key_event_date) : dayjs();
  openLinkModal.value = true;
};

const confirmLink = async () => {
  if (!linkingRecord.value || !linkDate.value) return;
  const date = linkDate.value.format('YYYY-MM-DD');
  const ok = await linkTransactionToKeyEvent(linkingRecord.value.transactionId, date);
  if (ok) {
    openLinkModal.value = false;
    linkingRecord.value = null;
    await refreshTable();
  }
};

const handleUnlink = async () => {
  if (!linkingRecord.value) return;
  const ok = await unlinkTransactionFromKeyEvent(linkingRecord.value.transactionId);
  if (ok) {
    openLinkModal.value = false;
    linkingRecord.value = null;
    await refreshTable();
  }
};
```

- [ ] **Step 5: 在表格上绑定 @link 事件**

在模板中的 `<transaction-record-table>` 上新增 `@link`：

```html
<transaction-record-table :items="tableData" @edit="updateTr" @delete="deleteTr" @link="handleLink" />
```

- [ ] **Step 6: 提交**

```bash
git add app/src/components/tr_view/TransactionRecordView.vue
git commit -m "feat: add link date picker modal to TransactionRecordView"
```

---

### Task 12: 前端 KeyEventView — 弹窗改为 Tabs + 关联交易表格

**Files:**
- Modify: `app/src/components/key_event_view/KeyEventView.vue`

- [ ] **Step 1: 新增 import**

在 script 顶部新增：

```ts
import { getLinkedTransactions, unlinkTransactionFromKeyEvent } from "@/backend/functions.ts";
import { centsToYuan } from "@/backend/functions";
import { useLedgerStore } from "@/stores/ledgerStore";
import type { TransactionRecord } from "@/types/billadm";
```

- [ ] **Step 2: 新增状态变量**

在 script 的状态区添加：

```ts
// 关联交易 Tab
const activeTab = ref('detail');
const linkedTransactions = ref<TransactionRecord[]>([]);
const linkedLoading = ref(false);
const ledgerStore = useLedgerStore();

const getLedgerName = (ledgerId: string): string => {
  const ledger = ledgerStore.ledgers.find(l => l.id === ledgerId);
  return ledger?.name || ledgerId;
};

const linkedCount = computed(() => linkedTransactions.value.length);
```

- [ ] **Step 3: 修改弹窗模板 — 替换 event-modal-content 为 Tabs 结构**

将现有的 `<a-modal>` 内容替换为 Tabs 结构。修改 `event-modal-content` div 内的内容（保留 color-picker 样式不变）：

```html
<div class="event-modal-content">
  <a-tabs v-model:activeKey="activeTab">
    <a-tab-pane key="detail" tab="详情">
      <a-input
        v-model:value="eventTitle"
        placeholder="标题（可选）"
        :maxlength="200"
        class="event-title-input"
      />
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
      <a-textarea
        v-model:value="eventContent"
        placeholder="记录今天发生的事情..."
        :rows="5"
        :maxlength="5000"
        show-count
      />
    </a-tab-pane>

    <a-tab-pane key="linked" :tab="`关联交易 (${linkedCount})`">
      <div v-if="linkedLoading" style="text-align:center;padding:24px">
        <a-spin />
      </div>
      <div v-else-if="linkedTransactions.length === 0" style="text-align:center;padding:24px;color:var(--billadm-color-text-secondary)">
        暂无关联交易记录
      </div>
      <a-table
        v-else
        :columns="linkedColumns"
        :data-source="linkedTransactions"
        :pagination="false"
        size="small"
        :scroll="{ y: 300 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'ledgerName'">
            <span style="font-size:12px;color:var(--billadm-color-text-secondary)">{{ getLedgerName(record.ledgerId) }}</span>
          </template>
          <template v-else-if="column.dataIndex === 'category'">
            <span>{{ record.category }}</span>
          </template>
          <template v-else-if="column.dataIndex === 'tags'">
            <div style="display:flex;flex-wrap:wrap;gap:4px">
              <a-tag v-for="tag in record.tags" :key="tag" style="font-size:11px">{{ tag }}</a-tag>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'price'">
            <span :style="{ color: record.transactionType === 'expense' ? 'var(--billadm-color-expense)' : record.transactionType === 'income' ? 'var(--billadm-color-income)' : 'var(--billadm-color-transfer)', fontFamily: 'var(--billadm-font-mono)' }">
              <template v-if="record.transactionType === 'expense'">-</template>
              <template v-else-if="record.transactionType === 'income'">+</template>
              {{ centsToYuan(record.price) }}
            </span>
          </template>
          <template v-else-if="column.dataIndex === 'action'">
            <a-button type="text" danger size="small" @click="handleUnlinkTr(record)">删除</a-button>
          </template>
        </template>
      </a-table>
    </a-tab-pane>
  </a-tabs>
</div>
```

- [ ] **Step 4: 新增 linkedColumns 定义**

在 script 中添加：

```ts
const linkedColumns = [
  { title: '账本', dataIndex: 'ledgerName', width: 100 },
  { title: '分类', dataIndex: 'category', width: 100 },
  { title: '标签', dataIndex: 'tags', width: 160 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '金额', dataIndex: 'price', width: 110, align: 'right' as const },
  { title: '操作', dataIndex: 'action', width: 70, align: 'center' as const },
];
```

- [ ] **Step 5: 新增 loadLinkedTransactions 和 handleUnlinkTr 方法**

```ts
const loadLinkedTransactions = async (date: string) => {
  linkedLoading.value = true;
  try {
    linkedTransactions.value = await getLinkedTransactions(date);
  } finally {
    linkedLoading.value = false;
  }
};

const handleUnlinkTr = async (record: TransactionRecord) => {
  const ok = await unlinkTransactionFromKeyEvent(record.transactionId);
  if (ok) {
    linkedTransactions.value = linkedTransactions.value.filter(
      t => t.transactionId !== record.transactionId
    );
  }
};
```

- [ ] **Step 6: 修改 onDayClick — 弹窗打开时加载关联交易**

在 `onDayClick` 方法中，`modalVisible.value = true;` 之前添加：

```ts
// 重置并加载关联交易
activeTab.value = 'detail';
loadLinkedTransactions(dateStr);
```

- [ ] **Step 7: 修改 modalTitle — 显示关联数量**

不需要修改，保持原样即可（modalTitle 已经按现有逻辑工作）。

- [ ] **Step 8: 保留 footer 模板不变**

原 `#footer` slot 的删除/取消/保存按钮逻辑不变（仅在 detail tab 下有效）。

- [ ] **Step 9: 类型检查**

```bash
cd app && npx vue-tsc --noEmit
```

- [ ] **Step 10: 提交**

```bash
git add app/src/components/key_event_view/KeyEventView.vue
git commit -m "feat: add tabs with linked transactions table to KeyEventView modal"
```

---

### Task 13: 前端 keyEventStore — 关联交易状态暂存

**Files:**
- Modify: `app/src/stores/keyEventStore.ts`

> 注：关联交易数据目前在 KeyEventView 组件内管理（local ref）。如果后续需要跨组件共享，再迁移到 store。当前无需修改 store，此 Task 可直接跳过或标记为 no-op。

- [ ] **Step 1: 不需要修改 — 跳过**

关联交易状态由 KeyEventView 本地管理即可，无需入 store。

---

### Task 14: 最终验证与端到端测试

- [ ] **Step 1: Go 完整构建**

```bash
cd kernel && go build -ldflags '-s -w -extldflags "-static"' -o Billadm-Kernel.exe
```

Expected: 构建成功

- [ ] **Step 2: 前端类型检查**

```bash
cd app && npx vue-tsc --noEmit
```

Expected: 无类型错误

- [ ] **Step 3: 前端构建**

```bash
cd app && npm run build
```

Expected: 构建成功

- [ ] **Step 4: Go 测试**

```bash
cd kernel && go test ./...
```

Expected: 所有测试通过

- [ ] **Step 5: 提交**

```bash
git add -A
git commit -m "chore: final verification - all builds and tests pass"
```
