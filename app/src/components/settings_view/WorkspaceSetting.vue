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
import {ref, onMounted} from 'vue';
import {FolderOpenOutlined} from "@ant-design/icons-vue";
import {useLedgerStore} from '@/stores/ledgerStore';
import {openWorkspace} from '@/backend/api/workspace';
import NotificationUtil from '@/backend/notification';

const ledgerStore = useLedgerStore();
const showFileSelect = ref(false);

const workspaceDir = ref('');

onMounted(async () => {
  workspaceDir.value = await window.electronAPI.getWorkspace() || '';
});

// 切换工作空间
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

.workspace-info {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-sm);
}

.info-label {
  font-size: var(--billadm-size-text-caption);
  font-weight: 500;
  color: var(--billadm-color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.info-value {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-major);
  word-break: break-all;
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  background-color: var(--billadm-color-minor-background);
  border-radius: var(--billadm-radius-md);
  font-family: monospace;
}

.workspace-actions {
  display: flex;
  gap: var(--billadm-space-md);
}
</style>
