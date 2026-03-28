<template>
  <a-row :gutter="16" style="width: 100%">
    <a-col :span="12">
      <a-space direction="vertical" style="width: 100%">
        <a-card title="收入" :head-style="incomeCss" :body-style="{padding:'0px'}">
          <a-collapse :bordered="false" ghost>
            <a-collapse-panel v-for="category in incomeCategory" :key="category.name" :header="category.name">
              <a-card>
                <a-card-grid v-for="tag in category.tags" :key="tag.name" style="width: 25%; text-align: center">
                  {{ tag.name }}
                </a-card-grid>
              </a-card>
            </a-collapse-panel>
          </a-collapse>
        </a-card>
        <a-card title="转账" :head-style="transferCss" :body-style="{padding:'0px'}">
          <a-collapse :bordered="false" ghost>
            <a-collapse-panel v-for="category in transferCategory" :key="category.name" :header="category.name">
              <a-card>
                <a-card-grid v-for="tag in category.tags" :key="tag.name" style="width: 25%; text-align: center">
                  {{ tag.name }}
                </a-card-grid>
              </a-card>
            </a-collapse-panel>
          </a-collapse>
        </a-card>
      </a-space>
    </a-col>
    <a-col :span="12">
      <a-card title="支出" :head-style="expenseCss" :body-style="{padding:'0px'}">
        <a-collapse :bordered="false" ghost>
          <a-collapse-panel v-for="category in expenseCategory" :key="category.name" :header="category.name">
            <a-card>
              <a-card-grid v-for="tag in category.tags" :key="tag.name" style="width: 25%; text-align: center">
                {{ tag.name }}
              </a-card-grid>
            </a-card>
          </a-collapse-panel>
        </a-collapse>
      </a-card>
    </a-col>
  </a-row>
</template>

<script lang="ts" setup>
import {type CSSProperties, ref, watch} from 'vue';
import type {Category, Tag, TransactionType} from '@/types/billadm';
import {useLedgerStore} from '@/stores/ledgerStore.ts';
import {getCategoryByType, getTagsByCategory} from '@/backend/functions.ts';
import {TransactionTypeToColor} from "@/backend/constant.ts";


const incomeCss: CSSProperties = {
  color: TransactionTypeToColor.get('income')
}

const expenseCss: CSSProperties = {
  color: TransactionTypeToColor.get('expense')
}

const transferCss: CSSProperties = {
  color: TransactionTypeToColor.get('transfer')
}

const ledgerStore = useLedgerStore();

// 扩展 Category 类型，加入 tags 字段（或使用接口继承）
interface CategoryWithTags extends Category {
  tags: Tag[];
}

const incomeCategory = ref<CategoryWithTags[]>([]);
const expenseCategory = ref<CategoryWithTags[]>([]);
const transferCategory = ref<CategoryWithTags[]>([]);

const refreshDataWithType = async (trType: TransactionType) => {
  const categories = await getCategoryByType(trType);
  // 为每个 category 并行加载其 tags
  return await Promise.all(
      categories.map(async (category) => {
        const tags = await getTagsByCategory(category.name);
        return {...category, tags};
      })
  );
}

const refreshData = async () => {
  incomeCategory.value = await refreshDataWithType('income');
  expenseCategory.value = await refreshDataWithType('expense');
  transferCategory.value = await refreshDataWithType('transfer');
};

// 监听 currentLedgerId 变化，并立即执行一次
watch(
    () => ledgerStore.currentLedgerId,
    async () => {
      await refreshData();
    },
    {immediate: true}
);
</script>