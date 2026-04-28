<template>
  <div class="workspace-setting">
    <!-- 主要工作空间卡片 -->
    <div class="workspace-hero">
      <div class="hero-content">
        <div class="hero-icon">
          <FolderOpenOutlined />
        </div>
        <div class="hero-text">
          <h2 class="hero-title">当前工作空间</h2>
          <p class="hero-path" :class="{ empty: !workspaceDir }">
            <span v-if="workspaceDir">{{ workspaceDir }}</span>
            <span v-else class="path-placeholder">未设置工作空间</span>
          </p>
        </div>
      </div>
      <div class="hero-action">
        <a-button type="primary" size="large" @click="showFileSelect = true">
          <template #icon>
            <SwapOutlined />
          </template>
          切换工作空间
        </a-button>
      </div>
    </div>

    <!-- 信息卡片区域 -->
    <div class="info-section">
      <div class="section-header">
        <span class="section-label">存储信息</span>
      </div>
      <div class="info-cards">
        <div class="info-card">
          <div class="info-card-icon storage">
            <DatabaseOutlined />
          </div>
          <div class="info-card-body">
            <div class="info-card-title">本地 SQLite</div>
            <div class="info-card-desc">每个账本独立数据库文件</div>
          </div>
        </div>
        <div class="info-card">
          <div class="info-card-icon security">
            <SafetyOutlined />
          </div>
          <div class="info-card-body">
            <div class="info-card-title">本地安全存储</div>
            <div class="info-card-desc">数据始终保存在本设备</div>
          </div>
        </div>
        <div class="info-card">
          <div class="info-card-icon folder">
            <FileProtectOutlined />
          </div>
          <div class="info-card-body">
            <div class="info-card-title">文件夹管理</div>
            <div class="info-card-desc">支持任意本地文件夹</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 工作空间选择弹窗 -->
    <billadm-file-select
      v-model="showFileSelect"
      title="选择工作目录"
      placeholder="请输入或选择工作目录路径"
      @confirm="handleSwitchWorkspace"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  FolderOpenOutlined,
  SwapOutlined,
  DatabaseOutlined,
  SafetyOutlined,
  FileProtectOutlined
} from "@ant-design/icons-vue";
import { useLedgerStore } from '@/stores/ledgerStore';
import { openWorkspace } from '@/backend/api/workspace';
import NotificationUtil from '@/backend/notification';

const ledgerStore = useLedgerStore();
const showFileSelect = ref(false);
const workspaceDir = ref('');

onMounted(async () => {
  workspaceDir.value = await window.electronAPI.getWorkspace() || '';
});

const handleSwitchWorkspace = async (newWorkspaceDir: string) => {
  try {
    await openWorkspace(newWorkspaceDir);
    window.electronAPI.setWorkspace(newWorkspaceDir);
    workspaceDir.value = newWorkspaceDir;
    await ledgerStore.init();
    NotificationUtil.success('切换工作空间成功');
  } catch (error) {
    NotificationUtil.error('切换工作空间失败', `${error}`);
  }
};
</script>

<style scoped>
.workspace-setting {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-lg);
}

/* Hero Section */
.workspace-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--billadm-space-lg);
  padding: var(--billadm-space-lg);
  background-color: var(--billadm-color-major-background);
  border-radius: var(--billadm-radius-lg);
  border: 1px solid var(--billadm-color-divider);
}

.hero-content {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-md);
}

.hero-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background-color: var(--billadm-color-primary);
  border-radius: var(--billadm-radius-md);
  color: var(--billadm-color-text-inverse);
  font-size: 20px;
  flex-shrink: 0;
}

.hero-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.hero-title {
  font-size: var(--billadm-size-text-body-sm);
  font-weight: 500;
  color: var(--billadm-color-text-secondary);
  margin: 0;
}

.hero-path {
  font-size: var(--billadm-size-text-body);
  font-family: var(--billadm-font-mono);
  color: var(--billadm-color-text-major);
  margin: 0;
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.hero-path.empty {
  color: var(--billadm-color-text-disabled);
  font-style: italic;
}

/* Info Section */
.info-section {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-sm);
}

.section-label {
  font-size: var(--billadm-size-text-body-sm);
  font-weight: 600;
  color: var(--billadm-color-text-secondary);
  padding-left: var(--billadm-space-xs);
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--billadm-space-sm);
}

.info-card {
  display: flex;
  align-items: flex-start;
  gap: var(--billadm-space-sm);
  padding: var(--billadm-space-md);
  background-color: var(--billadm-color-major-background);
  border-radius: var(--billadm-radius-lg);
  border: 1px solid var(--billadm-color-window-border);
}

.info-card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--billadm-radius-md);
  font-size: 16px;
  flex-shrink: 0;
}

.info-card-icon.storage {
  background-color: rgba(45, 125, 70, 0.12);
  color: var(--billadm-color-positive);
}

.info-card-icon.security {
  background-color: rgba(90, 127, 170, 0.12);
  color: var(--billadm-color-transfer);
}

.info-card-icon.folder {
  background-color: rgba(201, 162, 39, 0.12);
  color: var(--billadm-color-accent);
}

.info-card-body {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-card-title {
  font-size: var(--billadm-size-text-body-sm);
  font-weight: 600;
  color: var(--billadm-color-text-major);
}

.info-card-desc {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
  line-height: 1.4;
}
</style>
