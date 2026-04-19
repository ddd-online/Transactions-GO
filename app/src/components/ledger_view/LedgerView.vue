<template>
  <div class="ledger-view">
    <!-- 页面标题区 -->
    <header class="view-header">
      <div class="view-header-left">
        <h1 class="view-title">账本</h1>
        <span class="view-count">{{ ledgerStore.ledgers.length }} 个账本</span>
      </div>
      <button class="create-btn" @click="openCreateModal">
        <svg class="create-btn-icon" viewBox="0 0 20 20" fill="none">
          <path d="M10 4v12M4 10h12" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
        </svg>
        <span>新建账本</span>
      </button>
    </header>

    <!-- 账本网格 -->
    <div class="ledger-grid" v-if="ledgerStore.ledgers.length > 0">
      <article v-for="(ledger, index) in ledgerStore.ledgers" :key="ledger.id" class="ledger-card">
        <div class="ledger-card-inner">
          <!-- 左侧装饰条 -->
          <div class="ledger-card-accent" :style="{ backgroundColor: ledgerColors[index % ledgerColors.length] }"></div>

          <!-- 主要内容 -->
          <div class="ledger-card-body">
            <div class="ledger-card-header">
              <div class="ledger-icon"
                :style="{ backgroundColor: ledgerColors[index % ledgerColors.length] + '15', color: ledgerColors[index % ledgerColors.length] }">
                <svg viewBox="0 0 24 24" fill="none">
                  <path d="M4 19.5A2.5 2.5 0 016.5 17H20" stroke="currentColor" stroke-width="1.5"
                    stroke-linecap="round" stroke-linejoin="round" />
                  <path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" stroke="currentColor"
                    stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </div>
              <div class="ledger-info">
                <h3 class="ledger-name">{{ ledger.name }}</h3>
                <p class="ledger-desc" :class="{ 'is-empty': !ledger.description }">
                  {{ ledger.description || '暂无描述' }}
                </p>
              </div>
            </div>

            <div class="ledger-card-footer">
              <div class="ledger-meta">
                <span class="ledger-meta-item">
                  <svg class="meta-icon" viewBox="0 0 16 16" fill="none">
                    <rect x="2" y="3" width="12" height="11" rx="1.5" stroke="currentColor" stroke-width="1.2" />
                    <path d="M5 1.5v3M11 1.5v3M2 7h12" stroke="currentColor" stroke-width="1.2"
                      stroke-linecap="round" />
                  </svg>
                  {{ formatTimestamp(ledger.createdAt, 'YYYY/MM/dd') }}
                </span>
              </div>

              <div class="ledger-actions">
                <button class="action-btn" @click="openEditModal(ledger)" title="编辑">
                  <svg viewBox="0 0 16 16" fill="none">
                    <path d="M11.5 2.5l2 2-8 8H3.5v-2l8-8z" stroke="currentColor" stroke-width="1.2"
                      stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </button>
                <a-popconfirm title="确认删除此账本吗？" ok-text="确认" :showCancel="false"
                  @confirm="ledgerStore.deleteLedger(ledger.id)">
                  <button class="action-btn action-btn--danger" title="删除">
                    <svg viewBox="0 0 16 16" fill="none">
                      <path
                        d="M3 5h10M6 5V4a1 1 0 011-1h2a1 1 0 011 1v1M11 5v7a1.5 1.5 0 01-1.5 1.5H6.5A1.5 1.5 0 015 12V5"
                        stroke="currentColor" stroke-width="1.2" stroke-linecap="round" />
                    </svg>
                  </button>
                </a-popconfirm>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-state-visual">
        <svg class="empty-state-icon" viewBox="0 0 64 64" fill="none">
          <rect x="8" y="12" width="48" height="40" rx="4" stroke="currentColor" stroke-width="2" />
          <path d="M8 22h48" stroke="currentColor" stroke-width="2" />
          <path d="M20 8v8M44 8v8" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          <path d="M24 36h16M28 44h8" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
        </svg>
      </div>
      <h3 class="empty-state-title">暂无账本</h3>
      <p class="empty-state-desc">点击右上角按钮创建你的第一个账本</p>
    </div>

    <!-- 新建/编辑账本弹窗 -->
    <a-modal :title="modalTitle" :open="modalVisible" :confirm-loading="confirmLoading" ok-text="确认" cancel-text="取消"
      centered :width="400" @ok="handleOk" @cancel="modalVisible = false">
      <a-form layout="vertical" class="ledger-form">
        <a-form-item label="账本名称" required>
          <a-input v-model:value="ledgerForm.name" placeholder="请输入账本名称" :maxlength="50" show-count />
        </a-form-item>
        <a-form-item label="账本描述">
          <a-textarea v-model:value="ledgerForm.description" placeholder="请输入账本描述（可选）" :rows="3" :maxlength="200"
            show-count />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useLedgerStore } from "@/stores/ledgerStore";
import { formatTimestamp } from "@/backend/functions";
import type { Ledger } from '@/types/billadm';

const ledgerStore = useLedgerStore();

// 账本卡片颜色池 - 暖色调配色方案
const ledgerColors = [
  '#2D5A27', // 森林绿
  '#8B7355', // 暖棕
  '#C9A227', // 琥珀金
  '#C73E3A', // 朱红
  '#5A7FAA', // 灰蓝
  '#7A5C58', // 赭石
  '#5C7A6A', // 苔绿
  '#8A6B5C', // 驼色
];

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
  padding: var(--billadm-space-xl) var(--billadm-space-2xl);
  gap: var(--billadm-space-xl);
}

/* ========== View Header ========== */
.view-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.view-header-left {
  display: flex;
  align-items: baseline;
  gap: var(--billadm-space-md);
}

.view-title {
  font-size: var(--billadm-size-text-display-sm);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  margin: 0;
}

.view-count {
  font-size: var(--billadm-size-text-body-sm);
  color: var(--billadm-color-text-secondary);
}

/* ========== Create Button ========== */
.create-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 10px 18px;
  font-size: var(--billadm-size-text-body-sm);
  font-weight: var(--billadm-weight-medium);
  color: var(--billadm-color-text-inverse);
  background-color: var(--billadm-color-primary);
  border: none;
  border-radius: var(--billadm-radius-md);
  cursor: pointer;
  transition: all var(--billadm-transition-fast);
}

.create-btn-icon {
  width: 16px;
  height: 16px;
}

.create-btn:hover {
  background-color: var(--billadm-color-primary-light);
  transform: translateY(-1px);
  box-shadow: var(--billadm-shadow-sm);
}

.create-btn:active {
  transform: translateY(0);
  box-shadow: none;
}

/* ========== Ledger Grid ========== */
.ledger-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--billadm-space-lg);
  flex: 1;
  overflow-y: auto;
  align-content: start;
}

/* ========== Ledger Card ========== */
.ledger-card {
  position: relative;
  border-radius: var(--billadm-radius-lg);
  background-color: var(--billadm-color-major-background);
  transition: all var(--billadm-transition-normal);
}

.ledger-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--billadm-shadow-md);
}

.ledger-card-inner {
  display: flex;
  overflow: hidden;
  border-radius: var(--billadm-radius-lg);
  border: 1px solid var(--billadm-color-window-border);
}

/* 左侧装饰条 */
.ledger-card-accent {
  width: 4px;
  flex-shrink: 0;
}

.ledger-card-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: var(--billadm-space-lg);
  min-width: 0;
}

.ledger-card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--billadm-space-md);
  margin-bottom: var(--billadm-space-lg);
}

.ledger-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--billadm-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ledger-icon svg {
  width: 22px;
  height: 22px;
}

.ledger-info {
  flex: 1;
  min-width: 0;
}

.ledger-name {
  font-size: var(--billadm-size-text-section);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-2xs) 0;
  line-height: var(--billadm-height-tight);
}

.ledger-desc {
  font-size: var(--billadm-size-text-body-sm);
  color: var(--billadm-color-text-secondary);
  margin: 0;
  line-height: var(--billadm-height-snug);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.ledger-desc.is-empty {
  color: var(--billadm-color-text-disabled);
  font-style: italic;
}

.ledger-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: var(--billadm-space-md);
  border-top: 1px solid var(--billadm-color-divider);
}

.ledger-meta {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-md);
}

.ledger-meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

.meta-icon {
  width: 14px;
  height: 14px;
}

.ledger-actions {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-xs);
}

.action-btn {
  width: 30px;
  height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--billadm-color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--billadm-radius-sm);
  cursor: pointer;
  transition: all var(--billadm-transition-fast);
}

.action-btn svg {
  width: 15px;
  height: 15px;
}

.action-btn:hover {
  color: var(--billadm-color-primary);
  background-color: var(--billadm-color-hover-bg);
}

.action-btn--danger:hover {
  color: var(--billadm-color-expense);
  background-color: rgba(199, 62, 58, 0.08);
}

/* ========== Empty State ========== */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--billadm-space-3xl);
  text-align: center;
}

.empty-state-visual {
  width: 120px;
  height: 120px;
  margin-bottom: var(--billadm-space-xl);
  color: var(--billadm-color-primary);
  opacity: 0.25;
}

.empty-state-icon {
  width: 100%;
  height: 100%;
}

.empty-state-title {
  font-size: var(--billadm-size-text-title);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-sm) 0;
}

.empty-state-desc {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin: 0;
}

/* ========== Form ========== */
.ledger-form :deep(.ant-form-item-label > label) {
  font-weight: 500;
  color: var(--billadm-color-text-major);
}
</style>
