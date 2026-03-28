<template>
  <a-tabs v-model:activeKey="activeKey" style="height: 100%">
    <a-tab-pane key="category&tag">
      <template #tab>
        <span>
          <TagOutlined style="font-size: medium"/>
          分类与标签
        </span>
      </template>
      <billadm-category-tag-setting/>
    </a-tab-pane>
    <a-tab-pane key="workspace">
      <template #tab>
        <span>
          <FolderOpenOutlined style="font-size: medium"/>
          工作空间
        </span>
      </template>
    </a-tab-pane>
  </a-tabs>
</template>
<script lang="ts" setup>
import {ref} from 'vue';
import {FolderOpenOutlined, TagOutlined} from "@ant-design/icons-vue";

const activeKey = ref('category&tag');
</script>

<style scoped>
:deep(.ant-tabs-content-holder) {
  overflow-y: auto;
}
</style>