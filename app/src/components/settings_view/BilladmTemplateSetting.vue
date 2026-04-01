<template>
  <div class="template-setting">
    <div class="setting-header">
      <span class="setting-title">消费模板</span>
    </div>

    <a-table :columns="columns" :data-source="templates" :loading="loading" :pagination="false" row-key="template_id">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          <a-tag :color="getTypeColor(record.transaction_type)">
            {{ getTypeLabel(record.transaction_type) }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'tags'">
          <template v-if="record.tags && record.tags.length > 0">
            <a-tag v-for="tag in record.tags" :key="tag">{{ tag }}</a-tag>
          </template>
          <span v-else>-</span>
        </template>
        <template v-else-if="column.key === 'flags'">
          <template v-if="record.flags">
            <a-tag color="orange">离群值</a-tag>
          </template>
          <span v-else>-</span>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-popconfirm title="确定要删除该模板吗？" @confirm="handleDelete(record.template_id)" ok-text="确定" cancel-text="取消">
            <a-button type="link" danger size="small">删除</a-button>
          </a-popconfirm>
        </template>
      </template>
      <template #emptyText>
        <a-empty description="暂无模板" />
      </template>
    </a-table>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import type { TransactionTemplate } from '@/types/billadm';
import { getTemplatesByLedgerId, removeTemplate } from '@/backend/functions.ts';
import { useLedgerStore } from '@/stores/ledgerStore.ts';
import { TransactionTypeToLabel, TransactionTypeToColor } from '@/backend/constant.ts';

const ledgerStore = useLedgerStore();
const templates = ref<TransactionTemplate[]>([]);
const loading = ref(false);

const columns = [
  {
    title: '模板名称',
    dataIndex: 'template_name',
    key: 'name',
  },
  {
    title: '交易类型',
    dataIndex: 'transaction_type',
    key: 'type',
    width: 100,
  },
  {
    title: '分类',
    dataIndex: 'category',
    key: 'category',
  },
  {
    title: '标签',
    key: 'tags',
  },
  {
    title: '标记',
    dataIndex: 'flags',
    key: 'flags',
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true,
  },
  {
    title: '操作',
    key: 'action',
    width: 80,
  },
];

const loadTemplates = async () => {
  if (!ledgerStore.currentLedgerId) {
    templates.value = [];
    return;
  }
  loading.value = true;
  try {
    templates.value = await getTemplatesByLedgerId(ledgerStore.currentLedgerId);
  } finally {
    loading.value = false;
  }
};

const handleDelete = async (templateId: string) => {
  await removeTemplate(templateId);
  await loadTemplates();
};

const getTypeLabel = (type: string) => {
  return TransactionTypeToLabel.get(type) || type;
};

const getTypeColor = (type: string) => {
  return TransactionTypeToColor.get(type) || '#999';
};

onMounted(() => {
  loadTemplates();
});
</script>

<style scoped>
.template-setting {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.setting-header {
  margin-bottom: 16px;
}

.setting-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--billadm-color-text-primary);
}
</style>
