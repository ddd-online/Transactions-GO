<template>
  <div class="settings-view">
    <!-- 左侧设置导航 -->
    <aside class="settings-sidebar">
      <nav class="settings-nav">
        <div
          class="nav-item"
          :class="{ active: activeComponent === 'category-tag' }"
          @click="activeComponent = 'category-tag'"
        >
          <TagOutlined class="nav-icon"/>
          <span class="nav-text">分类与标签</span>
        </div>
        <div
          class="nav-item"
          :class="{ active: activeComponent === 'workspace' }"
          @click="activeComponent = 'workspace'"
        >
          <FolderOpenOutlined class="nav-icon"/>
          <span class="nav-text">工作空间</span>
        </div>
        <div
          class="nav-item"
          :class="{ active: activeComponent === 'data-import-export' }"
          @click="activeComponent = 'data-import-export'"
        >
          <CloudUploadOutlined class="nav-icon"/>
          <span class="nav-text">数据导入导出</span>
        </div>
        <div
          class="nav-item"
          :class="{ active: activeComponent === 'template' }"
          @click="activeComponent = 'template'"
        >
          <FileTextOutlined class="nav-icon"/>
          <span class="nav-text">消费模板</span>
        </div>
        <div
          class="nav-item"
          :class="{ active: activeComponent === 'mcp' }"
          @click="activeComponent = 'mcp'"
        >
          <SettingOutlined class="nav-icon"/>
          <span class="nav-text">MCP</span>
        </div>
        <div
          class="nav-item"
          :class="{ active: activeComponent === 'about' }"
          @click="activeComponent = 'about'"
        >
          <InfoCircleOutlined class="nav-icon"/>
          <span class="nav-text">关于</span>
        </div>
      </nav>
    </aside>

    <!-- 右侧内容区 -->
    <main class="settings-content">
      <div class="content-inner">
        <component :is="currentComponent" />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import {
  CloudUploadOutlined,
  FolderOpenOutlined,
  TagOutlined,
  FileTextOutlined,
  SettingOutlined,
  InfoCircleOutlined
} from "@ant-design/icons-vue";
import WorkspaceSetting from './WorkspaceSetting.vue';
import BilladmTemplateSetting from './BilladmTemplateSetting.vue';
import McpSetting from './McpSetting.vue';
import AboutSetting from './AboutSetting.vue';
import CategoryTagSetting from './BilladmCategoryTagSetting.vue';
import DataImportExportSetting from './DataImportExportSetting.vue';

const activeComponent = ref('category-tag');

const componentMap = {
  'category-tag': CategoryTagSetting,
  'workspace': WorkspaceSetting,
  'data-import-export': DataImportExportSetting,
  'template': BilladmTemplateSetting,
  'mcp': McpSetting,
  'about': AboutSetting,
};

const currentComponent = computed(() => {
  return componentMap[activeComponent.value as keyof typeof componentMap] || null;
});
</script>

<style scoped>
.settings-view {
  height: 100%;
  display: flex;
  background-color: var(--billadm-color-major-warm);
}

/* Sidebar */
.settings-sidebar {
  width: 220px;
  flex-shrink: 0;
  background-color: var(--billadm-color-major-warm);
  border-right: 1px solid var(--billadm-color-divider);
  display: flex;
  flex-direction: column;
  padding: 0;
}

.settings-nav {
  display: flex;
  flex-direction: column;
  padding: var(--billadm-space-md) var(--billadm-space-sm);
  gap: 2px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-md);
  padding: var(--billadm-space-md) var(--billadm-space-md);
  border-radius: var(--billadm-radius-md);
  cursor: pointer;
  transition: all var(--billadm-transition-fast);
  color: var(--billadm-color-text-secondary);
  position: relative;
  overflow: hidden;
}

.nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%) scaleY(0);
  width: 3px;
  height: 60%;
  background-color: var(--billadm-color-primary);
  border-radius: 0 2px 2px 0;
  transition: transform var(--billadm-transition-fast);
}

.nav-item:hover {
  background-color: var(--billadm-color-hover-bg);
  color: var(--billadm-color-text-major);
}

.nav-item.active {
  background-color: var(--billadm-color-active-bg);
  color: var(--billadm-color-primary);
}

.nav-item.active::before {
  transform: translateY(-50%) scaleY(1);
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
}

.nav-text {
  font-size: var(--billadm-size-text-body);
  font-weight: 500;
  white-space: nowrap;
}

/* Content */
.settings-content {
  flex: 1;
  min-width: 0;
  height: 100%;
  overflow-y: auto;
  background-color: var(--billadm-color-major-warm);
}

.content-inner {
  min-height: 100%;
  margin: 0 auto;
  padding: var(--billadm-space-md) var(--billadm-space-lg);
}
</style>
