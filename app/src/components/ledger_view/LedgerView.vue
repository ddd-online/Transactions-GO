<template>
  <div class="ledger-view">
    <!-- 悬浮按钮 -->
    <a-float-button type="primary" class="float-btn" @click="openCreateModal">
      <template #icon>
        <PlusOutlined />
      </template>
    </a-float-button>

    <!-- 账本卡片列表 -->
    <div class="ledger-list">
      <a-card v-for="ledger in ledgerStore.ledgers" :key="ledger.id" class="ledger-card">
        <template #title>
          <div class="ledger-card-header">
            <BookOutlined class="ledger-icon" />
            <span class="ledger-name">{{ ledger.name }}</span>
          </div>
        </template>
        <template #extra>
          <a-space>
            <a-button type="text" @click="openEditModal(ledger)">
              <template #icon>
                <EditOutlined />
              </template>
              编辑
            </a-button>
            <a-popconfirm title="确认删除吗" ok-text="确认" :showCancel="false" @confirm="ledgerStore.deleteLedger(ledger.id)">
              <a-button type="text" danger>
                <template #icon>
                  <DeleteOutlined />
                </template>
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>

        <div class="ledger-content">
          <p v-if="ledger.description" class="ledger-description">{{ ledger.description }}</p>
          <p v-else class="ledger-description-empty">暂无描述</p>

          <a-descriptions :column="2" size="small" class="ledger-meta">
            <a-descriptions-item label="创建时间">
              {{ formatTimestamp(ledger.createdAt, 'YYYY-MM-DD HH:mm') }}
            </a-descriptions-item>
            <a-descriptions-item label="更新时间">
              {{ formatTimestamp(ledger.updatedAt, 'YYYY-MM-DD HH:mm') }}
            </a-descriptions-item>
          </a-descriptions>
        </div>
      </a-card>
    </div>

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
          <a-input v-model:value="ledgerForm.name" placeholder="请输入账本名称" :maxlength="50" show-count />
        </a-form-item>
        <a-form-item label="账本描述">
          <a-textarea v-model:value="ledgerForm.description" placeholder="请输入账本描述（可选）" :rows="3" :maxlength="200" show-count />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useLedgerStore } from "@/stores/ledgerStore";
import { formatTimestamp } from "@/backend/functions";
import { PlusOutlined, BookOutlined, EditOutlined, DeleteOutlined } from "@ant-design/icons-vue";
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
  padding: 16px;
  gap: 16px;
  overflow-y: auto;
}

.float-btn {
  position: fixed;
  right: 50px;
  bottom: 80px;
}

.ledger-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ledger-card {
  border-radius: 8px;
}

.ledger-card :deep(.ant-card-head) {
  min-height: 48px;
  background-color: var(--billadm-color-minor-background);
  border-bottom: 1px solid var(--billadm-color-window-border);
}

.ledger-card :deep(.ant-card-body) {
  padding: 16px;
}

.ledger-card-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ledger-icon {
  font-size: 20px;
  color: var(--billadm-color-primary);
}

.ledger-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

.ledger-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ledger-description {
  margin: 0;
  font-size: 14px;
  color: var(--billadm-color-text-minor);
  line-height: 1.5;
}

.ledger-description-empty {
  margin: 0;
  font-size: 14px;
  color: var(--billadm-color-text-disabled);
  font-style: italic;
}

.ledger-meta {
  margin-top: 8px;
}

.ledger-form :deep(.ant-form-item-label > label) {
  font-weight: 500;
}
</style>
