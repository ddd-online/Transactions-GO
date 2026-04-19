<template>
  <div class="settings-view">
    <!-- 左侧设置列表 -->
    <div class="settings-sidebar">
      <div class="settings-list">
        <div
            class="settings-list-item"
            :class="{ active: activeComponent === 'category-tag' }"
            @click="activeComponent = 'category-tag'"
        >
          <TagOutlined class="settings-list-item-icon"/>
          <span class="settings-list-item-title">分类与标签</span>
        </div>
        <div
            class="settings-list-item"
            :class="{ active: activeComponent === 'workspace' }"
            @click="activeComponent = 'workspace'"
        >
          <FolderOpenOutlined class="settings-list-item-icon"/>
          <span class="settings-list-item-title">工作空间</span>
        </div>
        <div
            class="settings-list-item"
            :class="{ active: activeComponent === 'data-import-export' }"
            @click="activeComponent = 'data-import-export'"
        >
          <CloudUploadOutlined class="settings-list-item-icon"/>
          <span class="settings-list-item-title">数据导入导出</span>
        </div>
        <div
            class="settings-list-item"
            :class="{ active: activeComponent === 'template' }"
            @click="activeComponent = 'template'"
        >
          <FileTextOutlined class="settings-list-item-icon"/>
          <span class="settings-list-item-title">消费模板</span>
        </div>
        <div
            class="settings-list-item"
            :class="{ active: activeComponent === 'mcp' }"
            @click="activeComponent = 'mcp'"
        >
          <SettingOutlined class="settings-list-item-icon"/>
          <span class="settings-list-item-title">MCP</span>
        </div>
        <div
            class="settings-list-item"
            :class="{ active: activeComponent === 'about' }"
            @click="activeComponent = 'about'"
        >
          <InfoCircleOutlined class="settings-list-item-icon"/>
          <span class="settings-list-item-title">关于</span>
        </div>
      </div>
    </div>

    <!-- 右侧内容显示 -->
    <div class="settings-content">
      <billadm-category-tag-setting v-if="activeComponent === 'category-tag'"/>
      <workspace-setting v-else-if="activeComponent === 'workspace'"/>
      <data-import-export-setting v-else-if="activeComponent === 'data-import-export'"/>
      <billadm-template-setting v-else-if="activeComponent === 'template'"/>
      <mcp-setting v-else-if="activeComponent === 'mcp'"/>
      <about-setting v-else-if="activeComponent === 'about'"/>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';
import {CloudUploadOutlined, FolderOpenOutlined, TagOutlined, FileTextOutlined, SettingOutlined, InfoCircleOutlined} from "@ant-design/icons-vue";
import WorkspaceSetting from './WorkspaceSetting.vue';
import BilladmTemplateSetting from './BilladmTemplateSetting.vue';
import McpSetting from './McpSetting.vue';
import AboutSetting from './AboutSetting.vue';

const activeComponent = ref('category-tag');
</script>

<style scoped>
.settings-view {
  height: 100%;
  display: flex;
  padding: 0;
}

.settings-sidebar {
  flex: 0 0 200px;
  background-color: var(--billadm-color-minor-background);
  border-right: 1px solid var(--billadm-color-divider);
  overflow-y: auto;
  padding: var(--billadm-space-md) 0;
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.settings-list-item {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-md);
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  cursor: pointer;
  transition: all var(--billadm-transition-fast);
  color: var(--billadm-color-text-secondary);
  position: relative;
}

.settings-list-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  background-color: var(--billadm-color-primary);
  border-radius: 0 2px 2px 0;
  transition: height var(--billadm-transition-fast);
}

.settings-list-item:hover {
  background-color: var(--billadm-color-hover-bg);
  color: var(--billadm-color-text-major);
}

.settings-list-item.active {
  background-color: var(--billadm-color-hover-bg);
  color: var(--billadm-color-primary);
}

.settings-list-item.active::before {
  height: 20px;
}

.settings-list-item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 15px;
}

.settings-list-item-title {
  font-size: var(--billadm-size-text-body);
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.settings-content {
  flex: 1;
  min-width: 0;
  padding: var(--billadm-space-lg);
  overflow-y: auto;
  background-color: var(--billadm-color-major-background);
}
</style>
