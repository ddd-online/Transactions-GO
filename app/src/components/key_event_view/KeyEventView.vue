<template>
  <div class="key-event-view">
    <!-- 顶部工具栏 -->
    <div class="key-event-toolbar">
      <div class="key-event-toolbar-left"></div>
      <div class="key-event-toolbar-center">
        <a-button type="text" @click="goToPrevYear">
          <template #icon><LeftOutlined /></template>
        </a-button>
        <a-date-picker
          v-model:value="selectedYearDayjs"
          picker="year"
          :style="{ width: '120px' }"
          @change="(val) => onYearChange(val)"
        />
        <a-button type="text" @click="goToNextYear">
          <template #icon><RightOutlined /></template>
        </a-button>
      </div>
      <div class="key-event-toolbar-right"></div>
    </div>

    <!-- 全年日历 -->
    <div class="calendar-container" :class="{ 'is-loading': isLoading }">
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
            :class="{ 'day-cell--has-record': hasRecord(selectedYear, month, day), 'day-cell--today': isToday(selectedYear, month, day) }"
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
      :width="modalWidth"
      :body-style="{ height: modalBodyHeight }"
      @ok="handleSave"
      @cancel="modalVisible = false"
    >
      <div class="event-modal-content">
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
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useKeyEventStore } from "@/stores/keyEventStore";
import dayjs, { type Dayjs } from "dayjs";
import NotificationUtil from "@/backend/notification";
import { LeftOutlined, RightOutlined } from "@ant-design/icons-vue";

const keyEventStore = useKeyEventStore();
const isLoading = ref(false);

// ========== 弹窗尺寸 ==========
const windowWidth = ref(window.innerWidth);
const windowHeight = ref(window.innerHeight);
const modalWidth = computed(() => Math.floor(windowWidth.value * (2 / 3)));
const modalBodyHeight = computed(() => Math.floor(modalWidth.value * (3 / 4)) + 'px');

const handleResize = () => {
  windowWidth.value = window.innerWidth;
  windowHeight.value = window.innerHeight;
};

// ========== 年份选择 ==========
const selectedYearDayjs = ref<Dayjs>(dayjs());
const selectedYear = computed(() => selectedYearDayjs.value.year());

const goToPrevYear = () => {
  selectedYearDayjs.value = selectedYearDayjs.value.year(selectedYearDayjs.value.year() - 1);
  keyEventStore.fetchDatesByYear(selectedYearDayjs.value.year());
};

const goToNextYear = () => {
  selectedYearDayjs.value = selectedYearDayjs.value.year(selectedYearDayjs.value.year() + 1);
  keyEventStore.fetchDatesByYear(selectedYearDayjs.value.year());
};

const onYearChange = (val: Dayjs | string | null) => {
  if (val) {
    const d = typeof val === 'string' ? dayjs(val) : val;
    selectedYearDayjs.value = d;
    keyEventStore.fetchDatesByYear(d.year());
  }
};

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

const isToday = (year: number, month: number, day: number) => {
  const today = dayjs();
  return year === today.year() && month === today.month() + 1 && day === today.date();
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
    NotificationUtil.warning('内容不能为空', '请输入事件内容');
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

// ========== 初始化 ==========
onMounted(() => {
  keyEventStore.fetchDatesByYear(selectedYear.value);
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
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

/* ========== 日历容器 ========== */
.calendar-container {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--billadm-space-md);
  overflow-y: auto;
  align-content: start;
}

.calendar-container.is-loading {
  opacity: 0.6;
  pointer-events: none;
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

.day-cell--today::after {
  content: '';
  position: absolute;
  top: 3px;
  right: 3px;
  width: 6px;
  height: 6px;
  background-color: var(--billadm-color-error, #ff4d4f);
  border-radius: 50%;
}

.day-cell {
  position: relative;
}

/* ========== 弹窗内容 ========== */
.event-modal-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.event-modal-content :deep(.ant-input-textarea) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.event-modal-content :deep(.ant-input-textarea textarea) {
  flex: 1;
  resize: none;
}

.event-date-label {
  font-size: var(--billadm-size-text-body);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
}
</style>
