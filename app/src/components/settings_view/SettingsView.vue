<template>
  <div class="settings-view">
    <a-tabs v-model:activeKey="activeKey" class="settings-tabs">
      <a-tab-pane key="category&tag">
        <template #tab>
          <span>
            <TagOutlined style="font-size: 14px"/>
            分类与标签
          </span>
        </template>
        <billadm-category-tag-setting/>
      </a-tab-pane>
      <a-tab-pane key="workspace">
        <template #tab>
          <span>
            <FolderOpenOutlined style="font-size: 14px"/>
            工作空间
          </span>
        </template>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue';
import {FolderOpenOutlined, TagOutlined} from "@ant-design/icons-vue";

const activeKey = ref('category&tag');
</script>

<style scoped>
.settings-view {
  height: 100%;
  padding: 16px;
}

.settings-tabs {
  height: 100%;
}

.settings-tabs :deep(.ant-tabs-content-holder) {
  overflow-y: auto;
}
</style>
