# 关键事件 (Key Event) 页面实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现"关键事件"功能页面，包含后端 API（Go）和前端页面（Vue），支持以日历形式查看/编辑/删除每日文本记录。

**Architecture:** 后端采用现有 DAO-Service-Controller 三层架构，新增 `key_event` 数据表；前端采用 Vue + Pinia，日历用 CSS Grid 渲染，详情弹窗复用 `BilladmModal` 组件。

**Tech Stack:** Go (Gin + GORM), Vue 3 + TypeScript, Pinia, Ant Design Vue, SCSS

---

## 文件结构一览

### 后端新建文件
- `kernel/models/key_event.go`
- `kernel/dao/key_event_dao.go`
- `kernel/service/key_event_service.go`
- `kernel/api/key_event_controller.go`

### 后端修改文件
- `kernel/api/router.go` — 添加 `/key-events` 路由组

### 前端新建文件
- `app/src/components/key_event_view/KeyEventView.vue`
- `app/src/stores/keyEventStore.ts`
- `app/src/backend/api/key-event.ts`

### 前端修改文件
- `app/src/router/router.ts` — 注册路由
- `app/src/components/AppLeftBar.vue` — 添加导航入口
- `app/src/types/billadm.d.ts` — 添加 `KeyEvent` 类型

---

## 后端实现

### Task 1: 创建 KeyEvent 模型

**Files:**
- Create: `kernel/models/key_event.go`

- [ ] **Step 1: 创建 model 文件**

```go
package models

type KeyEvent struct {
    ID        string `gorm:"primaryKey;comment:事件UUID" json:"id"`
    Date      string `gorm:"uniqueIndex;not null;comment:日期 YYYY-MM-DD" json:"date"`
    Content   string `gorm:"type:text;comment:事件内容" json:"content"`
    CreatedAt int64  `gorm:"autoCreateTime:unix;not null;comment:创建时间" json:"createdAt"`
    UpdatedAt int64  `gorm:"autoUpdateTime:unix;not null;comment:更新时间" json:"updatedAt"`
}

func (k *KeyEvent) TableName() string {
    return "tbl_billadm_key_event"
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/models/key_event.go
git commit -m "feat(kernel): add KeyEvent model"
```

---

### Task 2: 创建 KeyEvent DAO 层

**Files:**
- Create: `kernel/dao/key_event_dao.go`

- [ ] **Step 1: 创建 DAO 文件**（参考 `ledger_dao.go` 的单例模式）

```go
package dao

import (
    "sync"

    "github.com/billadm/models"
    "github.com/billadm/workspace"
)

var (
    keyEventDao     KeyEventDao
    keyEventDaoOnce sync.Once
)

func GetKeyEventDao() KeyEventDao {
    if keyEventDao != nil {
        return keyEventDao
    }
    keyEventDaoOnce.Do(func() {
        keyEventDao = &keyEventDaoImpl{}
    })
    return keyEventDao
}

type KeyEventDao interface {
    UpsertKeyEvent(ws *workspace.Workspace, event *models.KeyEvent) error
    QueryByDate(ws *workspace.Workspace, date string) (*models.KeyEvent, error)
    QueryByYear(ws *workspace.Workspace, year string) ([]models.KeyEvent, error)
    DeleteByDate(ws *workspace.Workspace, date string) error
    GetDb(ws *workspace.Workspace) KeyEventDb
}

type KeyEventDb interface {
    Create(value interface{}) *gorm.DB
    Where(query interface{}, args ...interface{}) *gorm.DB
    First(dest interface{}) *gorm.DB
    Delete(value interface{}) *gorm.DB
    Updates(values interface{}) *gorm.DB
    Find(dest interface{}) *gorm.DB
}

var _ KeyEventDao = &keyEventDaoImpl{}

type keyEventDaoImpl struct{}

func (k *keyEventDaoImpl) UpsertKeyEvent(ws *workspace.Workspace, event *models.KeyEvent) error {
    return ws.GetDb().Save(event).Error
}

func (k *keyEventDaoImpl) QueryByDate(ws *workspace.Workspace, date string) (*models.KeyEvent, error) {
    var event models.KeyEvent
    err := ws.GetDb().Where("date = ?", date).First(&event).Error
    if err != nil {
        return nil, err
    }
    return &event, nil
}

func (k *keyEventDaoImpl) QueryByYear(ws *workspace.Workspace, year string) ([]models.KeyEvent, error) {
    events := make([]models.KeyEvent, 0)
    err := ws.GetDb().Where("date LIKE ?", year+"-%").Find(&events).Error
    return events, err
}

func (k *keyEventDaoImpl) DeleteByDate(ws *workspace.Workspace, date string) error {
    return ws.GetDb().Where("date = ?", date).Delete(&models.KeyEvent{}).Error
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/dao/key_event_dao.go
git commit -m "feat(kernel): add KeyEvent DAO layer"
```

---

### Task 3: 创建 KeyEvent Service 层

**Files:**
- Create: `kernel/service/key_event_service.go`

- [ ] **Step 1: 创建 Service 文件**（参考 `ledger_service.go` 的单例模式）

```go
package service

import (
    "sync"

    "github.com/billadm/dao"
    "github.com/billadm/models"
    "github.com/billadm/util"
    "github.com/billadm/workspace"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
)

var (
    keyEventService     KeyEventService
    keyEventServiceOnce sync.Once
)

func GetKeyEventService() KeyEventService {
    if keyEventService != nil {
        return keyEventService
    }
    keyEventServiceOnce.Do(func() {
        keyEventService = &keyEventServiceImpl{
            keyEventDao: dao.GetKeyEventDao(),
        }
    })
    return keyEventService
}

type KeyEventService interface {
    UpsertKeyEvent(ws *workspace.Workspace, date string, content string) error
    QueryByDate(ws *workspace.Workspace, date string) (*models.KeyEvent, error)
    QueryDatesByYear(ws *workspace.Workspace, year string) ([]string, error)
    DeleteByDate(ws *workspace.Workspace, date string) error
}

var _ KeyEventService = &keyEventServiceImpl{}

type keyEventServiceImpl struct {
    keyEventDao dao.KeyEventDao
}

// UpsertKeyEvent 根据 date 判断是否存在：存在则更新 content，不存在则新建
func (s *keyEventServiceImpl) UpsertKeyEvent(ws *workspace.Workspace, date string, content string) error {
    existing, err := s.keyEventDao.QueryByDate(ws, date)
    if err != nil && err != gorm.ErrRecordNotFound {
        return err
    }

    if existing != nil {
        // Update
        existing.Content = content
        return s.keyEventDao.UpsertKeyEvent(ws, existing)
    }

    // Create
    event := &models.KeyEvent{
        ID:      util.GetUUID(),
        Date:    date,
        Content: content,
    }
    return s.keyEventDao.UpsertKeyEvent(ws, event)
}

func (s *keyEventServiceImpl) QueryByDate(ws *workspace.Workspace, date string) (*models.KeyEvent, error) {
    return s.keyEventDao.QueryByDate(ws, date)
}

func (s *keyEventServiceImpl) QueryDatesByYear(ws *workspace.Workspace, year string) ([]string, error) {
    events, err := s.keyEventDao.QueryByYear(ws, year)
    if err != nil {
        return nil, err
    }
    dates := make([]string, len(events))
    for i, e := range events {
        dates[i] = e.Date
    }
    return dates, nil
}

func (s *keyEventServiceImpl) DeleteByDate(ws *workspace.Workspace, date string) error {
    logrus.Infof("delete key event, date: %s", date)
    return s.keyEventDao.DeleteByDate(ws, date)
}
```

- [ ] **Step 2: 提交**

```bash
git add kernel/service/key_event_service.go
git commit -m "feat(kernel): add KeyEvent service layer with upsert logic"
```

---

### Task 4: 创建 KeyEvent Controller 并注册路由

**Files:**
- Create: `kernel/api/key_event_controller.go`
- Modify: `kernel/api/router.go`

- [ ] **Step 1: 创建 Controller 文件**

```go
package api

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/billadm/models"
    "github.com/billadm/service"
    "github.com/billadm/workspace"
)

// GET /api/v1/key-events/dates/:year
func listKeyEventDates(c *gin.Context) {
    ret := models.NewResult()
    defer c.JSON(http.StatusOK, ret)

    ws := workspace.Manager.OpenedWorkspace()
    if ws == nil {
        ret.Code = -1
        ret.Msg = workspace.ErrOpenedWorkspaceNotFound
        return
    }

    year := c.Param("year")
    if year == "" {
        ret.Code = -1
        ret.Msg = "missing year parameter"
        return
    }

    dates, err := service.GetKeyEventService().QueryDatesByYear(ws, year)
    if err != nil {
        ret.Code = -1
        ret.Msg = err.Error()
        return
    }

    ret.Data = dates
}

// GET /api/v1/key-events/:date
func getKeyEvent(c *gin.Context) {
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
        ret.Msg = "missing date parameter"
        return
    }

    event, err := service.GetKeyEventService().QueryByDate(ws, date)
    if err != nil {
        ret.Code = -1
        ret.Msg = err.Error()
        return
    }

    ret.Data = event
}

// POST /api/v1/key-events  body: { date, content }
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

    content, _ := arg["content"].(string)

    if err := service.GetKeyEventService().UpsertKeyEvent(ws, date, content); err != nil {
        ret.Code = -1
        ret.Msg = err.Error()
        return
    }

    ret.Data = date
}

// DELETE /api/v1/key-events/:date
func deleteKeyEvent(c *gin.Context) {
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
        ret.Msg = "missing date parameter"
        return
    }

    if err := service.GetKeyEventService().DeleteByDate(ws, date); err != nil {
        ret.Code = -1
        ret.Msg = err.Error()
        return
    }
}
```

- [ ] **Step 2: 修改 router.go，在 v1 路由组添加 `/key-events` 路由**

在 `ServeAPI` 函数中，`v1 := ginServer.Group("/api/v1")` 的内部，在 MCP 路由组之前添加：

```go
// Key Events
keyEvents := v1.Group("/key-events")
{
    keyEvents.GET("/dates/:year", listKeyEventDates)
    keyEvents.GET("/:date", getKeyEvent)
    keyEvents.POST("", upsertKeyEvent)
    keyEvents.DELETE("/:date", deleteKeyEvent)
}
```

- [ ] **Step 3: 提交**

```bash
git add kernel/api/key_event_controller.go kernel/api/router.go
git commit -m "feat(kernel): add KeyEvent API endpoints and router registration"
```

---

## 前端实现

### Task 5: 添加 KeyEvent 类型定义

**Files:**
- Modify: `app/src/types/billadm.d.ts`

- [ ] **Step 1: 在 `billadm.d.ts` 末尾添加 `KeyEvent` 类型**

```ts
/**
 * 关键事件
 */
export interface KeyEvent {
    id: string;           // 事件UUID
    date: string;        // 日期 YYYY-MM-DD
    content: string;     // 事件内容
    createdAt: number;   // 创建时间戳
    updatedAt: number;   // 更新时间戳
}
```

- [ ] **Step 2: 提交**

```bash
git add app/src/types/billadm.d.ts
git commit -m "feat(types): add KeyEvent interface"
```

---

### Task 6: 创建 KeyEvent API 客户端

**Files:**
- Create: `app/src/backend/api/key-event.ts`

- [ ] **Step 1: 创建 API 封装文件**（参考 `ledger.ts` 模式）

```ts
import api from "@/backend/api/api-client";
import type { KeyEvent } from "@/types/billadm";

export async function queryKeyEventDatesByYear(year: number): Promise<string[]> {
    return api.get<string[]>(`/v1/key-events/dates/${year}`, '查询关键事件日期列表');
}

export async function queryKeyEventByDate(date: string): Promise<KeyEvent> {
    return api.get<KeyEvent>(`/v1/key-events/${date}`, '查询关键事件详情');
}

export async function saveKeyEvent(date: string, content: string): Promise<string> {
    return api.post<string>('/v1/key-events', { date, content }, '保存关键事件');
}

export async function deleteKeyEvent(date: string): Promise<void> {
    return api.delete<void>(`/v1/key-events/${date}`, '删除关键事件');
}
```

- [ ] **Step 2: 提交**

```bash
git add app/src/backend/api/key-event.ts
git commit -m "feat(api): add key event API client"
```

---

### Task 7: 创建 KeyEvent Store

**Files:**
- Create: `app/src/stores/keyEventStore.ts`

- [ ] **Step 1: 创建 Store 文件**（参考 `ledgerStore.ts` 模式）

```ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
    queryKeyEventDatesByYear,
    queryKeyEventByDate,
    saveKeyEvent,
    deleteKeyEvent
} from "@/backend/api/key-event";
import NotificationUtil from "@/backend/notification";
import type { KeyEvent } from "@/types/billadm";

export const useKeyEventStore = defineStore('keyEvent', () => {
    // 某年有记录的日期集合，用于日历高亮
    const datesWithRecords = ref(new Set<string>());
    const currentYear = ref(new Date().getFullYear());

    // 获取某年有记录的日期列表
    const fetchDatesByYear = async (year: number) => {
        try {
            const dates = await queryKeyEventDatesByYear(year);
            datesWithRecords.value = new Set(dates);
            currentYear.value = year;
        } catch (error) {
            NotificationUtil.error('查询关键事件失败', `${error}`);
        }
    };

    // 获取单日详情
    const fetchEventByDate = async (date: string): Promise<KeyEvent | null> => {
        try {
            return await queryKeyEventByDate(date);
        } catch (error) {
            // 404 表示当天没有记录，返回 null
            return null;
        }
    };

    // 保存事件（新建或更新）
    const saveEvent = async (date: string, content: string): Promise<void> => {
        try {
            await saveKeyEvent(date, content);
            datesWithRecords.value.add(date);
            NotificationUtil.success('保存成功');
        } catch (error) {
            NotificationUtil.error('保存失败', `${error}`);
            throw error;
        }
    };

    // 删除事件
    const deleteEvent = async (date: string): Promise<void> => {
        try {
            await deleteKeyEvent(date);
            datesWithRecords.value.delete(date);
            NotificationUtil.success('删除成功');
        } catch (error) {
            NotificationUtil.error('删除失败', `${error}`);
            throw error;
        }
    };

    // 某天是否有记录
    const hasRecord = (date: string): boolean => {
        return datesWithRecords.value.has(date);
    };

    return {
        datesWithRecords,
        currentYear,
        fetchDatesByYear,
        fetchEventByDate,
        saveEvent,
        deleteEvent,
        hasRecord,
    };
});
```

- [ ] **Step 2: 提交**

```bash
git add app/src/stores/keyEventStore.ts
git commit -m "feat(store): add keyEventStore"
```

---

### Task 8: 创建 KeyEventView 页面组件

**Files:**
- Create: `app/src/components/key_event_view/KeyEventView.vue`

- [ ] **Step 1: 创建 Vue 组件**

```vue
<template>
  <div class="key-event-view">
    <!-- 顶部工具栏 -->
    <div class="key-event-toolbar">
      <a-select
        v-model:value="selectedYear"
        :style="{ width: '120px' }"
        @change="onYearChange"
      >
        <a-select-option v-for="year in yearOptions" :key="year" :value="year">
          {{ year }}
        </a-select-option>
      </a-select>
    </div>

    <!-- 全年日历 -->
    <div class="calendar-container">
      <div v-for="month in 12" :key="month" class="month-grid">
        <div class="month-header">{{ month }}月</div>
        <div class="weekday-header">
          <span v-for="day in weekdays" :key="day">{{ day }}</span>
        </div>
        <div class="days-grid">
          <!-- 空白占位格 -->
          <div v-for="n in getMonthFirstWeekday(month)" :key="'blank-' + n" class="day-cell day-cell--blank" />
          <!-- 日期格 -->
          <div
            v-for="day in getMonthDays(month)"
            :key="day"
            class="day-cell"
            :class="{ 'day-cell--has-record': hasRecord(selectedYear, month, day) }"
            @click="onDayClick(selectedYear, month, day)"
          >
            {{ day }}
          </div>
        </div>
      </div>
    </div>

    <!-- 详情弹窗 -->
    <a-modal
      :open="modalVisible"
      :title="modalTitle"
      :confirm-loading="confirmLoading"
      ok-text="保存"
      cancel-text="取消"
      centered
      :width="480"
      @ok="handleSave"
      @cancel="modalVisible = false"
    >
      <div class="event-modal-content">
        <div class="event-date-label">{{ selectedDate }}</div>
        <a-textarea
          v-model:value="eventContent"
          placeholder="记录今天发生的事情..."
          :rows="5"
          :maxlength="5000"
          show-count
        />
      </div>
      <template #footer v-if="isEditMode">
        <a-popconfirm
          title="确认删除此条记录？"
          ok-text="确认"
          :showCancel="false"
          @confirm="handleDelete"
        >
          <a-button danger>删除</a-button>
        </a-popconfirm>
        <a-button @click="modalVisible = false">取消</a-button>
        <a-button type="primary" :loading="confirmLoading" @click="handleSave">保存</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useKeyEventStore } from "@/stores/keyEventStore";
import dayjs from "dayjs";

const keyEventStore = useKeyEventStore();

// ========== 年份选择 ==========
const currentYear = new Date().getFullYear();
const selectedYear = ref(currentYear);
const yearOptions = computed(() => {
  const years = [];
  for (let y = currentYear - 10; y <= currentYear + 5; y++) {
    years.push(y);
  }
  return years;
});

// ========== 工具函数 ==========
const weekdays = ['日', '一', '二', '三', '四', '五', '六'];

const getMonthDays = (month: number) => {
  return dayjs(`${selectedYear.value}-${String(month).padStart(2, '0')}-01`).daysInMonth();
};

const getMonthFirstWeekday = (month: number) => {
  return dayjs(`${selectedYear.value}-${String(month).padStart(2, '0')}-01`).day();
};

const formatDate = (year: number, month: number, day: number) => {
  return `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
};

const hasRecord = (year: number, month: number, day: number) => {
  const dateStr = formatDate(year, month, day);
  return keyEventStore.hasRecord(dateStr);
};

// ========== 弹窗状态 ==========
const modalVisible = ref(false);
const confirmLoading = ref(false);
const selectedDate = ref('');
const eventContent = ref('');
const isEditMode = ref(false);

const modalTitle = computed(() => {
  return isEditMode.value ? selectedDate.value : `添加事件 — ${selectedDate.value}`;
});

const onDayClick = async (year: number, month: number, day: number) => {
  const dateStr = formatDate(year, month, day);
  selectedDate.value = dateStr;

  if (keyEventStore.hasRecord(dateStr)) {
    // 有记录，加载详情
    confirmLoading.value = true;
    try {
      const event = await keyEventStore.fetchEventByDate(dateStr);
      if (event) {
        eventContent.value = event.content;
        isEditMode.value = true;
      } else {
        eventContent.value = '';
        isEditMode.value = false;
      }
    } finally {
      confirmLoading.value = false;
    }
  } else {
    // 无记录，直接新建
    eventContent.value = '';
    isEditMode.value = false;
  }

  modalVisible.value = true;
};

const handleSave = async () => {
  if (!eventContent.value.trim()) {
    return;
  }
  confirmLoading.value = true;
  try {
    await keyEventStore.saveEvent(selectedDate.value, eventContent.value.trim());
    modalVisible.value = false;
  } finally {
    confirmLoading.value = false;
  }
};

const handleDelete = async () => {
  confirmLoading.value = true;
  try {
    await keyEventStore.deleteEvent(selectedDate.value);
    modalVisible.value = false;
  } finally {
    confirmLoading.value = false;
  }
};

const onYearChange = (year: number) => {
  keyEventStore.fetchDatesByYear(year);
};

// ========== 初始化 ==========
onMounted(() => {
  keyEventStore.fetchDatesByYear(selectedYear.value);
});
</script>

<style scoped>
.key-event-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  gap: var(--billadm-space-md);
}

/* ========== 工具栏 ========== */
.key-event-toolbar {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-md);
}

/* ========== 日历容器 ========== */
.calendar-container {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--billadm-space-md);
  overflow-y: auto;
  align-content: start;
}

/* ========== 月份网格 ========== */
.month-grid {
  background-color: var(--billadm-color-major-background);
  border: 1px solid var(--billadm-color-window-border);
  border-radius: var(--billadm-radius-lg);
  padding: var(--billadm-space-sm);
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-xs);
}

.month-header {
  font-size: var(--billadm-size-text-body);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  text-align: center;
  padding-bottom: var(--billadm-space-xs);
}

.weekday-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  text-align: center;
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
  font-weight: var(--billadm-weight-medium);
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.day-cell {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--billadm-size-text-body-sm);
  color: var(--billadm-color-text-major);
  border-radius: var(--billadm-radius-sm);
  cursor: pointer;
  transition: background-color var(--billadm-transition-fast);
}

.day-cell:hover {
  background-color: var(--billadm-color-hover-bg);
}

.day-cell--blank {
  pointer-events: none;
}

.day-cell--has-record {
  background-color: var(--billadm-color-primary);
  color: var(--billadm-color-text-inverse);
  font-weight: var(--billadm-weight-semibold);
}

.day-cell--has-record:hover {
  background-color: var(--billadm-color-primary-light);
}

/* ========== 弹窗内容 ========== */
.event-modal-content {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-sm);
}

.event-date-label {
  font-size: var(--billadm-size-text-body);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
}
</style>
```

- [ ] **Step 2: 提交**

```bash
git add app/src/components/key_event_view/KeyEventView.vue
git commit -m "feat(ui): add KeyEventView calendar page component"
```

---

### Task 9: 注册路由

**Files:**
- Modify: `app/src/router/router.ts`

- [ ] **Step 1: 在 children 数组中添加关键事件路由**

在 `账本管理`、`消费记录`、`数据分析` 路由之后、`应用设置` 之前添加：

```ts
{
  name: '关键事件',
  path: 'key_event_view',
  component: () => import('@/components/key_event_view/KeyEventView.vue')
},
```

- [ ] **Step 2: 提交**

```bash
git add app/src/router/router.ts
git commit -m "feat(router): add key_event_view route"
```

---

### Task 10: 添加侧边栏导航入口

**Files:**
- Modify: `app/src/components/AppLeftBar.vue`

- [ ] **Step 1: 在导航按钮列表中添加关键事件入口**

在 `数据分析` 按钮之后、`settings_view` 分隔之前添加：

```vue
<button
  class="nav-btn"
  :class="{ active: route.path === '/key_event_view' }"
  @click="navigate('key_event_view')"
  aria-label="关键事件"
  title="关键事件"
>
  <StarOutlined style="font-size: 20px"/>
</button>
```

并导入图标：

```ts
import {StarOutlined, ...} from "@ant-design/icons-vue";
```

- [ ] **Step 2: 提交**

```bash
git add app/src/components/AppLeftBar.vue
git commit -m "feat(ui): add key event navigation button"
```

---

## 收尾

- [ ] **最终验证：** 在工作区运行 `cd kernel && go build` 确认后端编译通过
- [ ] **最终验证：** `cd app && npm run dev` 确认前端编译通过

---

## 自检清单

**Spec 覆盖检查：**
- [x] 数据表 `tbl_billadm_key_event` — Task 1
- [x] GET `/dates/{year}` 接口 — Task 4
- [x] GET `/{date}` 接口 — Task 4
- [x] POST `/key-events` (Upsert) — Task 4
- [x] DELETE `/{date}` — Task 4
- [x] KeyEvent 模型 — Task 1
- [x] DAO 层 — Task 2
- [x] Service 层（含 Upsert 逻辑）— Task 3
- [x] 前端日历组件 — Task 8
- [x] keyEventStore — Task 7
- [x] API 客户端 — Task 6
- [x] 路由注册 — Task 9
- [x] 侧边栏入口 — Task 10
- [x] 类型定义 — Task 5

**占位符检查：** 无 TBD/TODO，无模糊描述，所有步骤含实际代码。

**类型一致性：** `KeyEvent` 接口字段（id, date, content, createdAt, updatedAt）在 Task 1 定义后，Task 6 API 客户端、Task 7 Store 中保持一致。
