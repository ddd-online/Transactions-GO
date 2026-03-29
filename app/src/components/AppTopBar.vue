<template>
  <div class="menu-bar">
    <div class="menu-bar-left">
      <a-typography-text class="app-title">
        Transactions
      </a-typography-text>
    </div>
    <div class="menu-bar-right">
      <a-button type="text" @click="onMinimize">
        <template #icon>
          <LineOutlined/>
        </template>
      </a-button>
      <a-button type="text" @click="onMaximize">
        <template #icon>
          <BorderOutlined/>
        </template>
      </a-button>
      <a-button class="close-button" type="text" @click="onClose">
        <template #icon>
          <CloseOutlined/>
        </template>
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {BorderOutlined, CloseOutlined, LineOutlined} from "@ant-design/icons-vue";

const onMinimize = () => {
  window.electronAPI.minimizeWindow();
}

const onMaximize = () => {
  window.electronAPI.maximizeWindow();
}

const onClose = () => {
  window.electronAPI.closeWindow();
}
</script>

<style scoped>
.menu-bar {
  -webkit-app-region: drag;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  width: 100%;
  padding: 0 12px;
}

.menu-bar > * {
  -webkit-app-region: no-drag;
}

.menu-bar-left {
  display: flex;
  align-items: center;
}

.app-title {
  font-size: var(--billadm-size-text-title);
  font-weight: var(--billadm-weight-medium);
  color: var(--billadm-color-text-major);
  margin: 0;
  letter-spacing: 0.5px;
}

.menu-bar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
