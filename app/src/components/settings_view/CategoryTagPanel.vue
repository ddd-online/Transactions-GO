<template>
  <div class="category-tag-panel">
    <!-- 左侧：分类列表 -->
    <div class="panel-sidebar">
      <div class="sidebar-title">分类列表</div>
      <div class="category-list">
        <div
            v-for="category in categories"
            :key="category.name"
            class="category-item"
            :class="{ 'category-item-active': selectedCategory === category.name }"
            @click="selectCategory(category.name)"
        >
          {{ category.name }}
        </div>
      </div>
    </div>

    <!-- 右侧：标签列表 -->
    <div class="panel-content">
      <div class="content-title">
        {{ selectedCategory || '请选择分类' }}
        <span class="tag-count" v-if="selectedCategory">
          ({{ selectedTags.length }} 个标签)
        </span>
      </div>
      <div class="tag-list" v-if="selectedTags.length > 0">
        <div class="tag-grid">
          <div v-for="tag in selectedTags" :key="tag.name" class="tag-item">
            {{ tag.name }}
          </div>
        </div>
      </div>
      <div class="empty-state" v-else>
        <a-empty description="暂无标签"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue';
import type {TransactionType} from '@/types/billadm';
import {useLedgerStore} from '@/stores/ledgerStore';
import {getCategoryByType, getTagsByCategory} from '@/backend/functions';

interface Props {
  transactionType: TransactionType;
  activeColor?: string;
}

const props = defineProps<Props>();

const ledgerStore = useLedgerStore();

interface CategoryWithTags {
  name: string;
  tags: { name: string }[];
}

const categories = ref<CategoryWithTags[]>([]);
const selectedCategory = ref<string>('');
const selectedTags = ref<{ name: string }[]>([]);

// 加载分类数据
const loadCategories = async () => {
  const categoryList = await getCategoryByType(props.transactionType);
  categories.value = categoryList.map(c => ({name: c.name, tags: []}));

  // 加载所有分类的标签
  for (const category of categories.value) {
    const categoryTransactionType = `${category.name}:${props.transactionType}`;
    const tags = await getTagsByCategory(categoryTransactionType);
    category.tags = tags;
  }

  // 默认选中第一个分类
  const firstCategory = categories.value[0];
  if (firstCategory && !selectedCategory.value) {
    selectCategory(firstCategory.name);
  }
};

// 选择分类
const selectCategory = (categoryName: string) => {
  selectedCategory.value = categoryName;
  const category = categories.value.find(c => c.name === categoryName);
  selectedTags.value = category ? category.tags : [];
};

// 监听账本变化或交易类型变化
watch(
    () => [ledgerStore.currentLedgerId, props.transactionType],
    () => {
      selectedCategory.value = '';
      selectedTags.value = [];
      loadCategories();
    },
    {immediate: true}
);
</script>

<style scoped>
.category-tag-panel {
  height: 100%;
  display: flex;
  overflow: hidden;
}

.panel-sidebar {
  width: 200px;
  flex-shrink: 0;
  background-color: var(--billadm-color-minor-background);
  border-right: 1px solid var(--billadm-color-window-border);
  display: flex;
  flex-direction: column;
}

.sidebar-title {
  padding: 12px 16px;
  font-size: 14px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
  border-bottom: 1px solid var(--billadm-color-window-border);
}

.category-list {
  flex: 1;
  overflow-y: auto;
}

.category-item {
  padding: 12px 16px;
  cursor: pointer;
  color: var(--billadm-color-text-major);
  transition: all 0.2s;
}

.category-item:hover {
  background-color: var(--billadm-color-icon-hover-bg);
}

.category-item-active {
  background-color: var(--billadm-color-icon-hover-bg);
  color: v-bind('props.activeColor');
  border-left: 2px solid v-bind('props.activeColor');
}

.panel-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  overflow: hidden;
  background-color: var(--billadm-color-major-background);
}

.content-title {
  flex-shrink: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
  margin-bottom: 16px;
}

.tag-count {
  font-size: 13px;
  font-weight: 400;
  color: var(--billadm-color-text-minor);
  margin-left: 8px;
}

.tag-list {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.tag-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  gap: 12px;
}

.tag-item {
  padding: 8px 12px;
  text-align: center;
  background-color: var(--billadm-color-minor-background);
  border-radius: 4px;
  color: var(--billadm-color-text-major);
  font-size: 14px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}
</style>
