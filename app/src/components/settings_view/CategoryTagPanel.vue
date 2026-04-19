<template>
  <div class="category-tag-panel">
    <!-- 左侧：分类列表 -->
    <div class="panel-sidebar">
      <div class="sidebar-title">
        <span>分类列表</span>
        <a-button type="link" size="small" @click="openAddCategoryModal" :disabled="!ledgerStore.currentLedgerId">
          <template #icon><PlusOutlined /></template>
          添加
        </a-button>
      </div>
      <div class="category-list">
        <div
            v-for="(category, index) in categories"
            :key="category.name"
            class="category-item"
            :class="{ 'category-item-active': selectedCategory === category.name }"
            @click="selectCategory(category.name)"
        >
          <span class="category-name">{{ category.name }}<span class="record-count" v-if="category.recordCount">({{ category.recordCount }})</span></span>
          <span class="category-actions">
            <button class="action-btn" @click.stop="moveCategory(index, -1)" :class="{ 'disabled': index === 0 }" :disabled="index === 0">
              <UpOutlined />
            </button>
            <button class="action-btn" @click.stop="moveCategory(index, 1)" :class="{ 'disabled': index === categories.length - 1 }" :disabled="index === categories.length - 1">
              <DownOutlined />
            </button>
            <a-popconfirm title="确定要删除该分类及其所有标签吗？" @confirm="handleDeleteCategory(category.name)" ok-text="确定" cancel-text="取消">
              <button class="action-btn delete-btn" @click.stop>
                <DeleteOutlined />
              </button>
            </a-popconfirm>
          </span>
        </div>
      </div>
    </div>

    <!-- 右侧：标签列表 -->
    <div class="panel-content">
      <div class="content-title">
        <span>{{ selectedCategory || '请选择分类' }}</span>
        <span class="tag-count" v-if="selectedCategory">
          ({{ selectedTags.length }} 个标签)
        </span>
        <a-button v-if="selectedCategory" type="link" size="small" @click="openAddTagModal" style="margin-left: auto;">
          <template #icon><PlusOutlined /></template>
          添加
        </a-button>
      </div>
      <div class="tag-list" v-if="selectedTags.length > 0">
        <div v-for="(tag, index) in selectedTags" :key="tag.name" class="tag-row">
          <span class="tag-name">{{ tag.name }}<span class="record-count" v-if="tag.recordCount">({{ tag.recordCount }})</span></span>
          <span class="tag-actions">
            <button class="action-btn" @click="moveTag(index, -1)" :class="{ 'disabled': index === 0 }" :disabled="index === 0">
              <UpOutlined />
            </button>
            <button class="action-btn" @click="moveTag(index, 1)" :class="{ 'disabled': index === selectedTags.length - 1 }" :disabled="index === selectedTags.length - 1">
              <DownOutlined />
            </button>
            <a-popconfirm title="确定要删除该标签吗？" @confirm="handleDeleteTag(tag.name)" ok-text="确定" cancel-text="取消">
              <button class="action-btn delete-btn" @click.stop>
                <DeleteOutlined />
              </button>
            </a-popconfirm>
          </span>
        </div>
      </div>
      <div class="empty-state" v-else>
        <a-empty description="暂无标签" />
      </div>
    </div>

    <!-- 添加分类弹窗 -->
    <a-modal v-model:open="openCategoryModal" title="添加分类" @ok="confirmAddCategory" ok-text="确定" cancel-text="取消" centered>
      <a-form :model="categoryForm" :rules="categoryRules">
        <a-form-item label="分类名称" name="name">
          <a-input v-model:value="categoryForm.name" placeholder="请输入分类名称" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 添加标签弹窗 -->
    <a-modal v-model:open="openTagModal" title="添加标签" @ok="confirmAddTag" ok-text="确定" cancel-text="取消" centered>
      <a-form :model="tagForm" :rules="tagRules">
        <a-form-item label="标签名称" name="name">
          <a-input v-model:value="tagForm.name" placeholder="请输入标签名称" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue';
import type {TransactionType, Category, Tag} from '@/types/billadm';
import {useLedgerStore} from '@/stores/ledgerStore';
import {getCategoryByType, getTagsByCategory, addCategory, removeCategory, addTag, removeTag, reorderCategory, reorderTag} from '@/backend/functions';
import {PlusOutlined, DeleteOutlined, UpOutlined, DownOutlined} from "@ant-design/icons-vue";
import {message} from "ant-design-vue";
import type {Rule} from "ant-design-vue/es/form";

interface Props {
  transactionType: TransactionType;
  activeColor?: string;
}

const props = defineProps<Props>();

const ledgerStore = useLedgerStore();

interface CategoryWithTags extends Category {
  tags: Tag[];
}

const categories = ref<CategoryWithTags[]>([]);
const selectedCategory = ref<string>('');
const selectedTags = ref<Tag[]>([]);

// 添加分类弹窗
const openCategoryModal = ref(false);
const categoryForm = ref({ name: '' });
const categoryRules: Record<string, Rule[]> = {
  name: [{
    validator: (_: any, value: string) => {
      if (!value || !value.trim()) return Promise.reject(new Error('请输入分类名称'));
      // 检查是否已存在
      if (categories.value.some(c => c.name === value.trim())) {
        return Promise.reject(new Error('该分类已存在'));
      }
      return Promise.resolve();
    },
    trigger: 'blur',
  }],
};

// 添加标签弹窗
const openTagModal = ref(false);
const tagForm = ref({ name: '' });
const tagRules: Record<string, Rule[]> = {
  name: [{
    validator: (_: any, value: string) => {
      if (!value || !value.trim()) return Promise.reject(new Error('请输入标签名称'));
      // 检查是否已存在
      if (selectedTags.value.some(t => t.name === value.trim())) {
        return Promise.reject(new Error('该标签已存在'));
      }
      return Promise.resolve();
    },
    trigger: 'blur',
  }],
};

// 打开添加分类弹窗
const openAddCategoryModal = () => {
  categoryForm.value.name = '';
  openCategoryModal.value = true;
};

// 确认添加分类
const confirmAddCategory = async () => {
  const name = categoryForm.value.name.trim();
  if (!name) return;
  if (categories.value.some(c => c.name === name)) {
    message.error('该分类已存在');
    return;
  }
  try {
    await addCategory(ledgerStore.currentLedgerId!, name, props.transactionType);
    message.success('添加分类成功');
    openCategoryModal.value = false;
    await loadCategories();
    selectCategory(name);
  } catch {
    // error already shown in addCategory
  }
};

// 删除分类
const handleDeleteCategory = async (name: string) => {
  try {
    await removeCategory(name, props.transactionType, ledgerStore.currentLedgerId!);
    message.success('删除分类成功');
    if (selectedCategory.value === name) {
      selectedCategory.value = '';
      selectedTags.value = [];
    }
    await loadCategories();
  } catch {
    // error already shown in removeCategory
  }
};

// 移动分类
const moveCategory = async (index: number, direction: number) => {
  const newIndex = index + direction;
  if (newIndex < 0 || newIndex >= categories.value.length) return;

  const category = categories.value[index];
  const targetCategory = categories.value[newIndex];
  if (!category || !targetCategory) return;

  const categorySortOrder = category.sortOrder || 0;
  const targetSortOrder = targetCategory.sortOrder || 0;

  try {
    await reorderCategory(category.name, props.transactionType, targetSortOrder);
    await reorderCategory(targetCategory.name, props.transactionType, categorySortOrder);
    await loadCategories();
  } catch {
    // error already shown in reorderCategory
  }
};

// 打开添加标签弹窗
const openAddTagModal = () => {
  tagForm.value.name = '';
  openTagModal.value = true;
};

// 确认添加标签
const confirmAddTag = async () => {
  const name = tagForm.value.name.trim();
  if (!name) return;
  if (selectedTags.value.some(t => t.name === name)) {
    message.error('该标签已存在');
    return;
  }
  const categoryTransactionType = `${selectedCategory.value}:${props.transactionType}`;
  try {
    await addTag(name, categoryTransactionType);
    message.success('添加标签成功');
    openTagModal.value = false;
    await loadCategories();
    await selectCategory(selectedCategory.value);
  } catch {
    // error already shown in addTag
  }
};

// 删除标签
const handleDeleteTag = async (name: string) => {
  const categoryTransactionType = `${selectedCategory.value}:${props.transactionType}`;
  try {
    await removeTag(name, categoryTransactionType, ledgerStore.currentLedgerId!);
    message.success('删除标签成功');
    await loadCategories();
    await selectCategory(selectedCategory.value);
  } catch {
    // error already shown in removeTag
  }
};

// 移动标签
const moveTag = async (index: number, direction: number) => {
  const newIndex = index + direction;
  if (newIndex < 0 || newIndex >= selectedTags.value.length) return;

  const tag = selectedTags.value[index];
  const targetTag = selectedTags.value[newIndex];
  if (!tag || !targetTag) return;

  const categoryTransactionType = `${selectedCategory.value}:${props.transactionType}`;

  const tagSortOrder = tag.sortOrder || 0;
  const targetSortOrder = targetTag.sortOrder || 0;

  try {
    await reorderTag(tag.name, categoryTransactionType, targetSortOrder);
    await reorderTag(targetTag.name, categoryTransactionType, tagSortOrder);
    await loadCategories();
    selectCategory(selectedCategory.value);
  } catch {
    // error already shown in reorderTag
  }
};

// 加载分类数据
const loadCategories = async () => {
  const categoryList = await getCategoryByType(props.transactionType, ledgerStore.currentLedgerId!);
  categories.value = categoryList.map(c => ({name: c.name, transactionType: c.transactionType, sortOrder: c.sortOrder, recordCount: c.recordCount, tags: []}));

  // 加载所有分类的标签
  for (const category of categories.value) {
    const categoryTransactionType = `${category.name}:${props.transactionType}`;
    const tags = await getTagsByCategory(categoryTransactionType, ledgerStore.currentLedgerId!);
    category.tags = tags.map(t => ({name: t.name, categoryTransactionType: t.categoryTransactionType, sortOrder: t.sortOrder, recordCount: t.recordCount}));
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
  width: 260px;
  flex-shrink: 0;
  background-color: var(--billadm-color-minor-background);
  border-right: 1px solid var(--billadm-color-divider);
  display: flex;
  flex-direction: column;
}

.sidebar-title {
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  font-size: var(--billadm-size-text-body);
  font-weight: 600;
  color: var(--billadm-color-text-major);
  border-bottom: 1px solid var(--billadm-color-divider);
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.category-list {
  flex: 1;
  overflow-y: auto;
}

.category-item {
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  cursor: pointer;
  color: var(--billadm-color-text-major);
  transition: all var(--billadm-transition-fast);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
}

.category-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 0;
  background-color: v-bind('props.activeColor');
  border-radius: 0 2px 2px 0;
  transition: height var(--billadm-transition-fast);
}

.category-item:hover {
  background-color: var(--billadm-color-hover-bg);
}

.category-item-active {
  background-color: var(--billadm-color-hover-bg);
  color: v-bind('props.activeColor');
}

.category-item-active::before {
  height: 24px;
}

.category-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: var(--billadm-space-sm);
}

.record-count {
  color: var(--billadm-color-text-secondary);
  font-size: var(--billadm-size-text-caption);
  margin-left: var(--billadm-space-xs);
}

.category-actions {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-xs);
  opacity: 0;
  transition: opacity var(--billadm-transition-fast);
}

.category-item:hover .category-actions {
  opacity: 1;
}

.panel-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: var(--billadm-color-major-background);
}

.content-title {
  flex-shrink: 0;
  font-size: var(--billadm-size-text-title-sm);
  font-weight: 600;
  color: var(--billadm-color-text-major);
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  border-bottom: 1px solid var(--billadm-color-divider);
  display: flex;
  align-items: center;
}

.tag-count {
  font-size: var(--billadm-size-text-body);
  font-weight: 400;
  color: var(--billadm-color-text-secondary);
  margin-left: var(--billadm-space-sm);
}

.tag-list {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
  padding: var(--billadm-space-md);
}

.tag-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  transition: all var(--billadm-transition-fast);
  border-radius: var(--billadm-radius-md);
}

.tag-row:hover {
  background-color: var(--billadm-color-hover-bg);
}

.tag-name {
  flex: 1;
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-major);
}

.tag-actions {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-xs);
  opacity: 0;
  transition: opacity var(--billadm-transition-fast);
}

.tag-row:hover .tag-actions {
  opacity: 1;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  font-size: 14px;
  cursor: pointer;
  color: var(--billadm-color-text-secondary);
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--billadm-radius-md);
  transition: all var(--billadm-transition-fast);
}

.action-btn:hover:not(.disabled) {
  color: v-bind('props.activeColor');
  background-color: var(--billadm-color-hover-bg);
}

.action-btn.disabled {
  color: var(--billadm-color-text-disabled);
  cursor: not-allowed;
  border-color: transparent;
  background: transparent;
}

.action-btn.delete-btn:hover {
  color: var(--billadm-color-expense);
  background-color: rgba(199, 62, 58, 0.08);
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  flex: 1;
  min-height: 0;
}
</style>