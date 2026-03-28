<template>
  <div class="menu-bar">
    <div class="avatar">
      <a-avatar shape="square" :src="IconBilladm"/>
    </div>
    <div class="left-groups">
      <billadm-ledger-select v-if="route.path!='/ledger_view'"/>
    </div>
    <div class="center-groups">
      Billadm-{{ route.name }}
    </div>
    <div class="right-groups">
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
      <a-button class="closeButton" type="text" @click="onClose">
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

// 当前视图
const route = useRoute()

// 窗口控制
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
  justify-content: center;
  height: 100%;
  position: relative;
}

.menu-bar > * {
  -webkit-app-region: no-drag;
}

.avatar {
  width: var(--billadm-size-sider-width);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 左边按钮 将它与后面的元素隔开 */
.left-groups {
  margin-right: auto;
  display: flex;
  align-items: center;
  justify-content: center;
}

.center-groups {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

/* 右边按钮组 */
.right-groups {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-right: 8px;
}

.closeButton.ant-btn:hover {
  background-color: #f5222d;
  color: white;
}
</style>