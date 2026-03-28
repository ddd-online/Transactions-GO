<template>
  <a-table :columns="columns" :data-source="items" :pagination="false" :sticky="true" size="small">
    <template #bodyCell="{ column, record }">
      <template v-if="column.dataIndex==='transactionAt'">
        {{ formatTimestamp(record.transactionAt, 'YYYY-MM-DD') }}
      </template>

      <template v-else-if="column.dataIndex==='transactionType'">
        <a-typography-text :style="{color: TransactionTypeToColor.get(record.transactionType)}">
          {{ TransactionTypeToLabel.get(record.transactionType) || record.transactionType }}
        </a-typography-text>
      </template>

      <template v-else-if="column.dataIndex === 'tags'">
        <a-tag v-for="tag in record.tags" :key="tag" color="green">
          {{ tag }}
        </a-tag>
      </template>

      <template v-else-if="column.dataIndex === 'flags'">
        <a-tag v-if="record.outlier" key="outlier" color="orange">
          离群值
        </a-tag>
      </template>

      <template v-else-if="column.dataIndex === 'price'">
        {{ centsToYuan(record.price) }}
      </template>

      <template v-else-if="column.dataIndex === 'action'">
        <a-button type="text" @click="handleEdit(record as TransactionRecord)" :style="editButtonStyle">编辑</a-button>
        <a-popconfirm title="确认删除吗"
                      ok-text="确认"
                      @confirm="handleDelete(record as TransactionRecord)"
                      :showCancel="false">
          <a-button type="text" :style="deleteButtonStyle">删除</a-button>
        </a-popconfirm>
      </template>
    </template>
  </a-table>
</template>

<script setup lang="ts">
import type {TransactionRecord} from '@/types/billadm';
import {centsToYuan, formatTimestamp} from "@/backend/functions";
import {editButtonStyle, deleteButtonStyle} from "@/backend/styles";
import {TransactionTypeToColor, TransactionTypeToLabel} from "@/backend/constant";
import type {ColumnsType} from "ant-design-vue/es/table";

const columns: ColumnsType = [
  {
    title: '消费时间',
    dataIndex: 'transactionAt',
    width: 120,
    align: 'center'
  },
  {
    title: '交易类型',
    dataIndex: 'transactionType',
    width: 120,
    align: 'center'
  },
  {
    title: '消费类型',
    dataIndex: 'category',
    width: 120,
    align: 'center'
  },
  {
    title: '标签',
    dataIndex: 'tags'
  },
  {
    title: '标记',
    dataIndex: 'flags'
  },
  {
    title: '描述',
    dataIndex: 'description'
  },
  {
    title: '价格',
    dataIndex: 'price',
    width: 150,
    align: 'center'
  },
  {
    title: '操作',
    dataIndex: 'action',
    width: 150,
    align: 'center'
  }
];

interface Props {
  items: TransactionRecord[]
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'edit', record: TransactionRecord): void;
  (e: 'delete', record: TransactionRecord): void;
}>();

const handleEdit = (record: TransactionRecord) => {
  emit('edit', record);
};

const handleDelete = (record: TransactionRecord) => {
  emit('delete', record);
};
</script>
