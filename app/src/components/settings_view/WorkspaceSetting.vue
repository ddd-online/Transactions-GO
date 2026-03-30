<template>
  <div class="workspace-setting">
    <div class="workspace-info">
      <div class="info-label">当前工作空间</div>
      <div class="info-value">{{ workspaceDir || '未打开工作空间' }}</div>
    </div>
    <div class="workspace-actions">
      <a-button type="primary" @click="showFileSelect = true">
        <template #icon>
          <FolderOpenOutlined/>
        </template>
        切换工作空间
      </a-button>
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
import {ref} from 'vue';
import {FolderOpenOutlined} from "@ant-design/icons-vue";
import {useLedgerStore} from '@/stores/ledgerStore';
import {openWorkspace} from '@/backend/api/workspace';
import NotificationUtil from '@/backend/notification';

const ledgerStore = useLedgerStore();
const showFileSelect = ref(false);

const workspaceDir = ref(ledgerStore.workspaceStatus.workspaceDir || '');

// 切换工作空间
const handleSwitchWorkspace = async (newWorkspaceDir: string) => {
  try {
    await openWorkspace(newWorkspaceDir);
    await ledgerStore.refreshWorkspaceStatus();
    workspaceDir.value = ledgerStore.workspaceStatus.workspaceDir || '';
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
  gap: 24px;
}

.workspace-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-label {
  font-size: 14px;
  color: var(--billadm-color-text-minor);
}

.info-value {
  font-size: 16px;
  color: var(--billadm-color-text-major);
  word-break: break-all;
  padding: 12px;
  background-color: var(--billadm-color-minor-background);
  border-radius: 6px;
  font-family: monospace;
}

.workspace-actions {
  display: flex;
  gap: 12px;
}
</style>
