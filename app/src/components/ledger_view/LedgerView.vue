<template>
  <div class="ledger-view">
    <!-- 账本网格 -->
    <div class="ledger-grid" v-if="ledgerStore.ledgers.length > 0">
      <article v-for="(ledger, index) in ledgerStore.ledgers" :key="ledger.id"
        class="ledger-card"
        :style="{
          '--ledger-accent': `var(${ledgerColorVars[index % ledgerColorVars.length]})`,
          animationDelay: `${index * 60}ms`
        }">
        <div class="ledger-card-inner">
          <!-- 主要内容 -->
          <div class="ledger-card-body">
            <div class="ledger-card-header">
              <div class="ledger-icon"
                :style="{
                  backgroundColor: `color-mix(in srgb, var(${ledgerColorVars[index % ledgerColorVars.length]}) 15%, transparent)`,
                  color: `var(${ledgerColorVars[index % ledgerColorVars.length]})`,
                  '--card-index': index
                }">
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
                <a-popconfirm :title="`确认删除账本「${ledger.name}」？`" ok-text="确认" :showCancel="false"
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
      <h3 class="empty-state-title">还没有账本</h3>
      <p class="empty-state-desc">点击下方按钮创建你的第一个账本</p>
    </div>

    <!-- 悬浮按钮 -->
    <a-float-button type="primary" class="float-primary" @click="openCreateModal">
      <template #icon>
        <PlusOutlined />
      </template>
    </a-float-button>

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
import { PlusOutlined } from "@ant-design/icons-vue";

const ledgerStore = useLedgerStore();

// Ledger accent colors - use CSS custom properties for token system + dark mode support
const ledgerColorVars = [
  '--billadm-ledger-forest',
  '--billadm-ledger-warm-brown',
  '--billadm-ledger-amber',
  '--billadm-ledger-vermillion',
  '--billadm-ledger-slate-blue',
  '--billadm-ledger-ochre',
  '--billadm-ledger-moss',
  '--billadm-ledger-camel',
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
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  gap: var(--billadm-space-md);
}

/* ========== Ledger Grid ========== */
.ledger-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--billadm-space-md);
  flex: 1;
  overflow-y: auto;
  align-content: start;
}

/* ========== Ledger Card ========== */
.ledger-card {
  border-radius: var(--billadm-radius-lg);
  background-color: var(--billadm-color-major-background);
  border: 1px solid var(--billadm-color-window-border);
  transition: box-shadow var(--billadm-transition-normal), transform var(--billadm-transition-fast);
  position: relative;
  overflow: hidden;
}

.ledger-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background-color: var(--ledger-accent, var(--billadm-ledger-forest));
  opacity: 0.6;
  transition: opacity var(--billadm-transition-fast);
}

.ledger-card:hover {
  box-shadow: var(--billadm-shadow-lg);
  transform: translateY(-2px);
}

.ledger-card:hover::before {
  opacity: 1;
}

.ledger-card-inner {
  display: flex;
  overflow: hidden;
  border-radius: var(--billadm-radius-lg);
}

.ledger-card-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: var(--billadm-space-md);
  min-width: 0;
}

.ledger-card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--billadm-space-sm);
  margin-bottom: var(--billadm-space-md);
}

.ledger-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--billadm-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.ledger-icon svg {
  width: 20px;
  height: 20px;
}

.ledger-info {
  flex: 1;
  min-width: 0;
}

.ledger-name {
  font-size: var(--billadm-size-text-section);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-xs) 0;
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
  padding-top: var(--billadm-space-sm);
  border-top: 1px solid var(--billadm-color-divider);
}

.ledger-meta {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-sm);
}

.ledger-meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
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
  width: 28px;
  height: 28px;
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
  width: 14px;
  height: 14px;
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
  position: relative;
}

.empty-state::before {
  content: '';
  position: absolute;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--billadm-color-primary) 0%, transparent 60%);
  opacity: 0.08;
  pointer-events: none;
}

.empty-state-visual {
  width: 96px;
  height: 96px;
  margin-bottom: var(--billadm-space-xl);
  color: var(--billadm-color-primary);
  opacity: 0.25;
  transform: scale(1);
  transition: transform var(--billadm-transition-slow);
}

.empty-state:hover .empty-state-visual {
  transform: scale(1.05);
  opacity: 0.35;
}

.empty-state-icon {
  width: 100%;
  height: 100%;
}

.empty-state-title {
  font-size: var(--billadm-size-text-display-sm);
  font-weight: var(--billadm-weight-bold);
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-sm) 0;
  letter-spacing: -0.01em;
}

.empty-state-desc {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin: 0;
  max-width: 280px;
}

/* ========== Form ========== */
.ledger-form :deep(.ant-form-item-label > label) {
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

/* ========== Floating Button ========== */
.float-primary {
  right: 40px;
  bottom: 72px;
  box-shadow: var(--billadm-shadow-md);
  transition: box-shadow var(--billadm-transition-normal), transform var(--billadm-transition-fast);
}

.float-primary:hover {
  box-shadow: var(--billadm-shadow-xl);
  transform: scale(1.08);
}
</style>
