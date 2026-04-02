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
            v-for="category in categories"
            :key="category.name"
            class="category-item"
            :class="{ 'category-item-active': selectedCategory === category.name }"
            @click="selectCategory(category.name)"
        >
          <span class="category-name">{{ category.name }}</span>
          <a-popconfirm title="确定要删除该分类及其所有标签吗？" @confirm="handleDeleteCategory(category.name)" ok-text="确定" cancel-text="取消">
            <DeleteOutlined class="category-delete-btn" @click.stop />
          </a-popconfirm>
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
        <div class="tag-grid">
          <div v-for="tag in selectedTags" :key="tag.name" class="tag-item">
            {{ tag.name }}
            <CloseOutlined class="tag-delete-btn" @click="handleDeleteTag(tag.name)" />
          </div>
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
import type {TransactionType} from '@/types/billadm';
import {useLedgerStore} from '@/stores/ledgerStore';
import {getCategoryByType, getTagsByCategory, addCategory, removeCategory, addTag, removeTag} from '@/backend/functions';
import {PlusOutlined, DeleteOutlined, CloseOutlined} from "@ant-design/icons-vue";
import {message} from "ant-design-vue";
import type {Rule} from "ant-design-vue/es/form";

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
  display: flex;
  align-items: center;
  justify-content: space-between;
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
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.category-item:hover {
  background-color: var(--billadm-color-icon-hover-bg);
}

.category-item-active {
  background-color: var(--billadm-color-icon-hover-bg);
  color: v-bind('props.activeColor');
  border-left: 2px solid v-bind('props.activeColor');
}

.category-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-delete-btn {
  opacity: 0;
  color: #ff4d4f;
  font-size: 12px;
  transition: opacity 0.2s;
}

.category-item:hover .category-delete-btn {
  opacity: 1;
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
  display: flex;
  align-items: center;
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
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.tag-delete-btn {
  opacity: 0;
  color: #ff4d4f;
  font-size: 10px;
  transition: opacity 0.2s;
  cursor: pointer;
}

.tag-item:hover .tag-delete-btn {
  opacity: 1;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}
</style>