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
