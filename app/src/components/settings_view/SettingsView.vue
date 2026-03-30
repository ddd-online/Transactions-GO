<template>
  <div class="settings-view">
    <!-- 左侧设置列表 -->
    <div class="settings-sidebar">
      <a-menu
          v-model:selectedKeys="selectedKeys"
          mode="inline"
          theme="light"
          class="settings-menu"
      >
        <a-menu-item key="category-tag" @click="activeComponent = 'category-tag'">
          <template #icon>
            <TagOutlined/>
          </template>
          <span>分类与标签</span>
        </a-menu-item>
        <a-menu-item key="workspace" @click="activeComponent = 'workspace'">
          <template #icon>
            <FolderOpenOutlined/>
          </template>
          <span>工作空间</span>
        </a-menu-item>
      </a-menu>
    </div>

    <!-- 右侧内容显示 -->
    <div class="settings-content">
      <billadm-category-tag-setting v-if="activeComponent === 'category-tag'"/>
      <div v-else-if="activeComponent === 'workspace'" class="placeholder">
        <a-empty description="工作空间功能开发中"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';
import {FolderOpenOutlined, TagOutlined} from "@ant-design/icons-vue";

const selectedKeys = ref(['category-tag']);
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
  padding: 12px;
}

.settings-menu {
  background: transparent;
  border-inline-end: none !important;
}

.settings-content {
  flex: 1;
  min-width: 0;
  padding: 20px;
  overflow-y: auto;
}

.placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>
