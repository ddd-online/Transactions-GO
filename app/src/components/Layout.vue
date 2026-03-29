<template>
  <div class="app-shell">
    <!-- 工作空间选择弹窗 -->
    <billadm-file-select
        v-model="showWorkspaceSelect"
        title="新建工作目录或打开已存在的工作目录"
        @confirm="handleOpenWorkspace"
    />

    <!-- 主布局 -->
    <div class="app-shell-body">
      <!-- 顶部状态栏 -->
      <header class="app-header">
        <app-top-bar/>
      </header>

      <!-- 主体区域 -->
      <div class="app-main">
        <!-- 侧边栏 -->
        <aside class="app-sidebar">
          <app-left-bar/>
        </aside>

        <!-- 内容区域 -->
        <main class="app-content">
          <router-view/>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from "vue";
import {useLedgerStore} from "@/stores/ledgerStore.ts";
import {openWorkspace} from "@/backend/api/workspace.ts";
import NotificationUtil from "@/backend/notification.ts";

const ledgerStore = useLedgerStore();
const showWorkspaceSelect = ref(true);

const handleOpenWorkspace = async (workspaceDir: string) => {
  try {
    await openWorkspace(workspaceDir);
    await initWorkspace();
  } catch (error) {
    NotificationUtil.error('打开工作空间失败', `${error}`);
    showWorkspaceSelect.value = true;
  }
}

const initWorkspace = async () => {
  await ledgerStore.refreshWorkspaceStatus();
  if (!ledgerStore.workspaceStatus.isOpened) {
    showWorkspaceSelect.value = true;
    return;
  }
  showWorkspaceSelect.value = false;
  window.electronAPI.setWorkspace(ledgerStore.workspaceStatus.workspaceDir);
  await ledgerStore.init();
}

onMounted(initWorkspace);
</script>

<style scoped>
.app-shell {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.app-shell-body {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

/* 顶部状态栏 - 次要背景色 */
.app-header {
  height: var(--billadm-size-header-height);
  background-color: var(--billadm-color-minor-background);
  flex-shrink: 0;
  border-bottom: 1px solid var(--billadm-color-window-border);
}

/* 主体区域 */
.app-main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 侧边栏 - 次要背景色 */
.app-sidebar {
  width: var(--billadm-size-sider-width);
  min-width: var(--billadm-size-sider-width);
  height: 100%;
  background-color: var(--billadm-color-minor-background);
  flex-shrink: 0;
  border-right: 1px solid var(--billadm-color-window-border);
}

/* 内容区域 - 主背景色 */
.app-content {
  flex: 1;
  background-color: var(--billadm-color-major-background);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>
