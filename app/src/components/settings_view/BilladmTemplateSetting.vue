<template>
  <div class="template-setting">
    <div class="setting-header">
      <span class="setting-title">消费模板</span>
    </div>

    <a-table :columns="columns" :data-source="templates" :loading="loading" :pagination="false" row-key="template_id">
      <template #bodyCell="{ column, record, index }">
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
          <span class="action-buttons">
            <button class="action-btn" @click="moveTemplate(index, -1)" :disabled="index === 0">
              <UpOutlined />
            </button>
            <button class="action-btn" @click="moveTemplate(index, 1)" :disabled="index === templates.length - 1">
              <DownOutlined />
            </button>
            <a-popconfirm title="确定要删除该模板吗？" @confirm="handleDelete(record.template_id)" ok-text="确定" cancel-text="取消">
              <button class="action-btn delete-btn" @click.stop>
                <DeleteOutlined />
              </button>
            </a-popconfirm>
          </span>
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
import { message } from 'ant-design-vue';
import { UpOutlined, DownOutlined, DeleteOutlined } from "@ant-design/icons-vue";
import type { TransactionTemplate } from '@/types/billadm';
import { getTemplatesByLedgerId, removeTemplate, reorderTemplate } from '@/backend/functions.ts';
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
    width: 120,
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
  message.success('删除模板成功');
  await loadTemplates();
};

const moveTemplate = async (index: number, direction: number) => {
  const newIndex = index + direction;
  if (newIndex < 0 || newIndex >= templates.value.length) return;

  const template = templates.value[index];
  const targetTemplate = templates.value[newIndex];
  if (!template || !targetTemplate) return;

  const templateSortOrder = template.sort_order || 0;
  const targetSortOrder = targetTemplate.sort_order || 0;

  try {
    await reorderTemplate(template.template_id!, ledgerStore.currentLedgerId!, targetSortOrder);
    await reorderTemplate(targetTemplate.template_id!, ledgerStore.currentLedgerId!, templateSortOrder);
    await loadTemplates();
  } catch {
    // error already shown in reorderTemplate
  }
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

.action-buttons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  font-size: 12px;
  cursor: pointer;
  color: var(--billadm-color-text-secondary);
  background: transparent;
  border: 1px solid transparent;
  border-radius: 4px;
  transition: all 0.2s;
}

.action-btn:hover:not(.disabled) {
  color: var(--billadm-color-primary, #1890ff);
  background-color: var(--billadm-color-icon-hover-bg);
  border-color: var(--billadm-color-primary, #1890ff);
}

.action-btn.disabled {
  color: var(--billadm-color-text-disabled, #ccc);
  cursor: not-allowed;
  border-color: transparent;
  background: transparent;
}

.action-btn.delete-btn {
  color: #ff4d4f;
  border-color: transparent;
}

.action-btn.delete-btn:hover {
  color: #fff;
  background-color: #ff4d4f;
  border-color: #ff4d4f;
}
</style>
