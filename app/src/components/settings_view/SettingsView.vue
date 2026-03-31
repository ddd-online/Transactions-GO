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
      </div>
    </div>

    <!-- 右侧内容显示 -->
    <div class="settings-content">
      <billadm-category-tag-setting v-if="activeComponent === 'category-tag'"/>
      <workspace-setting v-else-if="activeComponent === 'workspace'"/>
      <data-import-export-setting v-else-if="activeComponent === 'data-import-export'"/>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';
import {CloudUploadOutlined, FolderOpenOutlined, TagOutlined} from "@ant-design/icons-vue";
import WorkspaceSetting from './WorkspaceSetting.vue';

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
  border-right: 1px solid var(--billadm-color-window-border);
  overflow-y: auto;
  padding: 12px 8px;
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.settings-list-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--billadm-color-text-secondary);
}

.settings-list-item:hover {
  background-color: var(--billadm-color-icon-hover-bg);
}

.settings-list-item.active {
  background-color: #ffffff;
  color: var(--billadm-color-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.settings-list-item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 14px;
}

.settings-list-item-title {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.settings-content {
  flex: 1;
  min-width: 0;
  padding: 20px;
  overflow-y: auto;
}
</style>
