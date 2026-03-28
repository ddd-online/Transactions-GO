<template>
  <a-modal
      title="筛选消费记录"
      v-model:open="open"
      width="800px"
      @cancel="confirmFilterModal"
      centered
      :closable="false"
  >
    <template #footer>
      <a-button key="clear" @click="clearAllConditions">清除条件</a-button>
      <a-button key="confirm" type="primary" @click="confirmFilterModal">确认</a-button>
    </template>

    <!-- 输入区域 -->
    <a-form layout="vertical">
      <a-row :gutter="16">
        <!-- 交易类型 -->
        <a-col :span="4">
          <a-form-item label="交易类型">
            <a-select
                v-model:value="tempTransactionType"
                placeholder="请选择交易类型"
                allow-clear
            >
              <a-select-option
                  v-for="opt in transactionTypeOptions"
                  :key="opt.value"
                  :value="opt.value"
              >
                {{ opt.label }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>

        <!-- 分类 -->
        <a-col :span="4">
          <a-form-item label="分类">
            <a-select
                v-model:value="tempCategory"
                placeholder="请选择分类"
                :options="categories"
                allow-clear
                @change="onCategoryChange"
            />
          </a-form-item>
        </a-col>

        <!-- 标签 -->
        <a-col :span="8">
          <a-form-item label="标签">
            <a-select
                v-model:value="tempTags"
                mode="multiple"
                placeholder="请选择标签"
                :options="tags"
                allow-clear
            />
          </a-form-item>
        </a-col>

        <!-- 标签策略 -->
        <a-col :span="4">
          <a-form-item label="标签匹配策略">
            <a-select v-model:value="tempTagPolicy" placeholder="策略">
              <a-select-option value="any">任意</a-select-option>
              <a-select-option value="all">全部</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>

        <!-- 标签取反 -->
        <a-col :span="4">
          <a-form-item label="是否取反">
            <a-select v-model:value="tempTagNot" placeholder="策略">
              <a-select-option value="yes">是</a-select-option>
              <a-select-option value="no">否</a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>

      <!-- 描述关键词 -->
      <a-form-item label="描述包含">
        <a-input v-model:value="tempDescription" placeholder="输入关键词"/>
      </a-form-item>

      <!-- 添加按钮 -->
      <a-form-item>
        <a-button type="dashed" @click="addCondition" block>
          + 添加筛选条件
        </a-button>
      </a-form-item>

      <!-- 已添加条件列表 -->
      <a-divider orientation="left">已添加条件 ({{ trQueryConditionItems.length }})</a-divider>
      <div v-if="trQueryConditionItems.length === 0" style="color: #999; padding: 12px 0;">
        暂无筛选条件
      </div>
      <a-space direction="vertical" style="width: 100%">
        <a-card
            v-for="(item, index) in trQueryConditionItems"
            :key="index"
            size="small"
            style="border: 1px solid #ddd; border-radius: 4px"
        >
          <template #title>
            <span style="font-weight: normal">
              {{ getTransactionTypeLabel(item.transactionType) }}
              {{ item.category ? ` · ${item.category}` : '' }}
              {{ item.tags && item.tags.length ? ` · 标签(${item.tags.join(', ')})` : '' }}
              {{ item.description ? ` · "${item.description}"` : '' }}
            </span>
          </template>
          <template #extra>
            <a-button type="text" danger @click="deleteCondition(index)">删除</a-button>
          </template>
          <div>
            <a-tag color="blue" v-if="item.tagPolicy === 'all'">全部标签匹配</a-tag>
            <a-tag color="green" v-else>任意标签匹配</a-tag>
            <a-tag color="red" v-if="item.tagNot">标签取反</a-tag>
          </div>
        </a-card>
      </a-space>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue';
import type {Category, TrQueryConditionItem} from '@/types/billadm';
import type {DefaultOptionType} from 'ant-design-vue/es/vc-cascader';
import {getCategoryByType, getTagsByCategory} from '@/backend/functions.ts';
import {useTrQueryConditionStore} from '@/stores/trQueryConditionStore.ts';

// 双向绑定 modal 开关
const open = defineModel<boolean>();

const trQueryConditionStore = useTrQueryConditionStore();
const trQueryConditionItems = ref<TrQueryConditionItem[]>([]);

// 临时输入状态
const tempTransactionType = ref<string | undefined>(undefined);
const tempCategory = ref<string | undefined>(undefined);
const tempTags = ref<string[]>([]);
const tempTagPolicy = ref<'any' | 'all'>('any');
const tempTagNot = ref<'yes' | 'no'>("no");
const tempDescription = ref<string>('');

// 交易类型选项
const transactionTypeOptions = [
  {label: '收入', value: 'income'},
  {label: '支出', value: 'expense'},
  {label: '转账', value: 'transfer'},
];

// 分类 & 标签选项
const categories = ref<DefaultOptionType[]>([]);
const tags = ref<DefaultOptionType[]>([]);

// 交易类型变化 → 刷新分类
watch(() => tempTransactionType.value, async (newVal) => {
  if (!newVal) {
    categories.value = [];
    tempCategory.value = undefined;
    return;
  }
  const categoryList: Category[] = await getCategoryByType(newVal);
  categories.value = categoryList.map((c) => ({value: c.name}));
});

// 分类变化 → 刷新标签
watch(() => tempCategory.value, async (newVal) => {
  if (!newVal) {
    tags.value = [];
    tempTags.value = [];
    return;
  }
  const tagList = await getTagsByCategory(newVal);
  tags.value = tagList.map((t) => ({value: t.name}));
});

// 打开时加载 store 中的条件
watch(open, (newVal) => {
  if (newVal) {
    trQueryConditionItems.value = [...(trQueryConditionStore.trQueryConditionItems || [])];
    // 重置临时输入
    resetTempInputs();
  }
});

function resetTempInputs() {
  tempTransactionType.value = undefined;
  tempCategory.value = undefined;
  tempTags.value = [];
  tempTagPolicy.value = 'any';
  tempTagNot.value = "no";
  tempDescription.value = '';
}

function getTransactionTypeLabel(type: string): string {
  const map: Record<string, string> = {
    income: '收入',
    expense: '支出',
    transfer: '转账',
  };
  return map[type] || type;
}

function addCondition() {
  // 至少需要一个有效字段
  if (
      !tempTransactionType.value &&
      !tempCategory.value &&
      tempTags.value.length === 0 &&
      !tempDescription.value.trim()
  ) {
    // 可选：提示用户
    return;
  }

  const newItem: TrQueryConditionItem = {
    transactionType: tempTransactionType.value || '',
    category: tempCategory.value || '',
    tags: [...tempTags.value],
    tagPolicy: tempTagPolicy.value,
    tagNot: tempTagNot.value === "yes",
    description: tempDescription.value.trim(),
  };

  trQueryConditionItems.value.push(newItem);
  resetTempInputs();
}

function deleteCondition(index: number) {
  trQueryConditionItems.value.splice(index, 1);
}

function clearAllConditions() {
  trQueryConditionItems.value = [];
}

function confirmFilterModal() {
  trQueryConditionStore.trQueryConditionItems = trQueryConditionItems.value;
  open.value = false;
}

function onCategoryChange() {
  tempTags.value = [];
}
</script>