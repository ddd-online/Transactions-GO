<template>
  <div class="menu-bar">
    <div class="menu-bar-left">
      <div class="menu-bar-avatar">
        <a-avatar shape="square" :src="IconBilladm"/>
      </div>
      <billadm-ledger-select v-if="route.path!='/ledger_view'"/>
    </div>
    <div class="menu-bar-center">
      <a-typography-text class="typography-section">
        Billadm-{{ route.name }}
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
import {useRoute} from 'vue-router'
import {BorderOutlined, CloseOutlined, LineOutlined} from "@ant-design/icons-vue";
import IconBilladm from '@/assets/icons/billadm.svg';

const route = useRoute()

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
  padding: 0 8px;
}

.menu-bar > * {
  -webkit-app-region: no-drag;
}

.menu-bar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-bar-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.menu-bar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.menu-bar-avatar {
  width: var(--billadm-size-sider-width);
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
