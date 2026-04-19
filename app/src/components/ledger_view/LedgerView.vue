<template>
  <div class="ledger-view">
    <!-- 账本卡片列表 -->
    <div class="ledger-list">
      <div
        v-for="ledger in ledgerStore.ledgers"
        :key="ledger.id"
        class="ledger-card-item"
        :class="{ active: ledger.id === ledgerStore.currentLedgerId }"
      >
        <div class="ledger-card-header">
          <div class="ledger-card-icon">
            <BookOutlined />
          </div>
          <div class="ledger-card-info">
            <h3 class="ledger-card-name">{{ ledger.name }}</h3>
            <p v-if="ledger.description" class="ledger-card-desc">{{ ledger.description }}</p>
            <p v-else class="ledger-card-desc-empty">暂无描述</p>
          </div>
        </div>

        <div class="ledger-card-meta">
          <span class="ledger-card-meta-item">
            <ClockCircleOutlined /> 创建于 {{ formatTimestamp(ledger.createdAt, 'YYYY-MM-DD') }}
          </span>
          <span class="ledger-card-meta-item">
            <EditOutlined /> {{ formatTimestamp(ledger.updatedAt, 'YYYY-MM-DD HH:mm') }}
          </span>
        </div>

        <div class="ledger-card-actions">
          <a-button type="text" class="action-btn" @click="openEditModal(ledger)">
            <EditOutlined /> 编辑
          </a-button>
          <a-popconfirm
            title="确认删除此账本吗？"
            ok-text="确认"
            :showCancel="false"
            @confirm="ledgerStore.deleteLedger(ledger.id)"
          >
            <a-button type="text" class="action-btn danger">
              <DeleteOutlined /> 删除
            </a-button>
          </a-popconfirm>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="ledgerStore.ledgers.length === 0" class="empty-state">
      <BookOutlined class="empty-state-icon" />
      <p class="empty-state-text">暂无账本，点击右下角按钮创建一个</p>
    </div>

    <!-- 悬浮按钮 -->
    <a-float-button type="primary" class="float-btn" @click="openCreateModal">
      <template #icon>
        <PlusOutlined />
      </template>
    </a-float-button>

    <!-- 新建/编辑账本弹窗 -->
    <a-modal
      :title="modalTitle"
      :open="modalVisible"
      :confirm-loading="confirmLoading"
      ok-text="确认"
      cancel-text="取消"
      centered
      @ok="handleOk"
      @cancel="modalVisible = false"
    >
      <a-form layout="vertical" class="ledger-form">
        <a-form-item label="账本名称" required>
          <a-input
            v-model:value="ledgerForm.name"
            placeholder="请输入账本名称"
            :maxlength="50"
            show-count
          />
        </a-form-item>
        <a-form-item label="账本描述">
          <a-textarea
            v-model:value="ledgerForm.description"
            placeholder="请输入账本描述（可选）"
            :rows="3"
            :maxlength="200"
            show-count
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useLedgerStore } from "@/stores/ledgerStore";
import { formatTimestamp } from "@/backend/functions";
import {
  PlusOutlined,
  BookOutlined,
  EditOutlined,
  DeleteOutlined,
  ClockCircleOutlined
} from "@ant-design/icons-vue";
import type { Ledger } from '@/types/billadm';

const ledgerStore = useLedgerStore();

const modalVisible = ref<boolean>(false);
const confirmLoading = ref<boolean>(false);
const modalTitle = ref<string>('');
const editingLedgerId = ref<string>('');

const ledgerForm = reactive({
  name: '',
  description: '',
});

const openCreateModal = () => {
  modalTitle.value = '新建账本';
  editingLedgerId.value = '';
  ledgerForm.name = '';
  ledgerForm.description = '';
  modalVisible.value = true;
};

const openEditModal = (ledger: Ledger) => {
  modalTitle.value = '编辑账本';
  editingLedgerId.value = ledger.id;
  ledgerForm.name = ledger.name;
  ledgerForm.description = ledger.description || '';
  modalVisible.value = true;
};

const handleOk = async () => {
  if (!ledgerForm.name.trim()) {
    return;
  }

  confirmLoading.value = true;
  try {
    if (editingLedgerId.value) {
      await ledgerStore.modifyLedger(editingLedgerId.value, ledgerForm.name.trim(), ledgerForm.description.trim());
    } else {
      await ledgerStore.createLedger(ledgerForm.name.trim(), ledgerForm.description.trim());
    }
    modalVisible.value = false;
  } finally {
    confirmLoading.value = false;
  }
};
</script>

<style scoped>
.ledger-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: var(--billadm-space-xl);
  gap: var(--billadm-space-xl);
  overflow-y: auto;
}

.ledger-list {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-lg);
}

.ledger-card-item {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-lg);
  padding: var(--billadm-space-xl);
  background-color: var(--billadm-color-major-background);
  border: 1px solid var(--billadm-color-window-border);
  border-radius: var(--billadm-radius-lg);
  box-shadow: var(--billadm-shadow-sm);
  transition: all var(--billadm-transition-normal);
}

.ledger-card-item:hover {
  border-color: var(--billadm-color-primary);
  box-shadow: var(--billadm-shadow-md);
  transform: translateY(-2px);
}

.ledger-card-item.active {
  border-color: var(--billadm-color-primary);
  background-color: rgba(45, 90, 39, 0.02);
}

.ledger-card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--billadm-space-lg);
}

.ledger-card-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--billadm-radius-md);
  background-color: rgba(45, 90, 39, 0.1);
  color: var(--billadm-color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.ledger-card-info {
  flex: 1;
  min-width: 0;
}

.ledger-card-name {
  font-family: var(--billadm-font-display);
  font-size: var(--billadm-size-text-title-sm);
  font-weight: 500;
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-xs) 0;
}

.ledger-card-desc {
  font-family: var(--billadm-font-body);
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin: 0;
  line-height: 1.5;
}

.ledger-card-desc-empty {
  font-family: var(--billadm-font-body);
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-disabled);
  font-style: italic;
  margin: 0;
}

.ledger-card-meta {
  display: flex;
  gap: var(--billadm-space-xl);
  padding-top: var(--billadm-space-md);
  border-top: 1px solid var(--billadm-color-divider);
}

.ledger-card-meta-item {
  font-family: var(--billadm-font-body);
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
  display: flex;
  align-items: center;
  gap: var(--billadm-space-xs);
}

.ledger-card-actions {
  display: flex;
  gap: var(--billadm-space-sm);
}

.action-btn {
  font-family: var(--billadm-font-body);
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
  border-radius: var(--billadm-radius-md);
  transition: all var(--billadm-transition-fast);
}

.action-btn:hover {
  color: var(--billadm-color-primary);
  background-color: var(--billadm-color-hover-bg);
}

.action-btn.danger:hover {
  color: var(--billadm-color-expense);
  background-color: rgba(199, 62, 58, 0.1);
}

.float-btn {
  position: fixed;
  right: 48px;
  bottom: 80px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--billadm-space-3xl);
  color: var(--billadm-color-text-secondary);
}

.empty-state-icon {
  font-size: 56px;
  color: var(--billadm-color-primary);
  opacity: 0.3;
  margin-bottom: var(--billadm-space-xl);
}

.empty-state-text {
  font-family: var(--billadm-font-body);
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  text-align: center;
  margin: 0;
}

.ledger-form :deep(.ant-form-item-label > label) {
  font-family: var(--billadm-font-body);
  font-weight: 500;
  color: var(--billadm-color-text-major);
}
</style>
