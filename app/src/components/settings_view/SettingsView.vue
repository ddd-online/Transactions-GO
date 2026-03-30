<template>
  <div class="settings-view">
    <!-- 主内容区：左侧边栏 + 右侧内容显示 -->
    <div class="settings-main">
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
  padding: 16px;
}

.settings-main {
  flex: 1;
  display: flex;
  gap: 12px;
  min-height: 0;
  overflow: hidden;
}

.settings-sidebar {
  width: 200px;
  flex-shrink: 0;
  background-color: var(--billadm-color-minor-background, #f5f5f5);
  border: 1px solid var(--billadm-color-border, #e8e8e8);
  border-radius: 8px;
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
  background-color: var(--billadm-color-major-background, #fff);
  border: 1px solid var(--billadm-color-border, #e8e8e8);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.04);
}

.placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>
