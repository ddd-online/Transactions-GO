<template>
  <div class="category-tag-panel">
    <!-- 左侧：分类列表 -->
    <aside class="panel-sidebar">
      <header class="sidebar-header">
        <h2 class="sidebar-title">分类</h2>
        <a-button
          type="primary"
          size="small"
          shape="round"
          @click="openAddCategoryModal"
          :disabled="!ledgerStore.currentLedgerId"
        >
          <template #icon><PlusOutlined /></template>
          新增
        </a-button>
      </header>
      <nav class="category-list" role="list">
        <div
          v-for="(category, index) in categories"
          :key="category.name"
          class="category-item"
          :class="{ 'is-active': selectedCategory === category.name }"
          role="listitem"
          @click="selectCategory(category.name)"
        >
          <div class="category-info">
            <span class="category-name">{{ category.name }}</span>
            <span class="category-count" v-if="category.recordCount">{{ category.recordCount }}</span>
          </div>
          <div class="category-actions">
            <button
              class="icon-btn"
              @click.stop="moveCategory(index, -1)"
              :disabled="index === 0"
              aria-label="上移"
            >
              <UpOutlined />
            </button>
            <button
              class="icon-btn"
              @click.stop="moveCategory(index, 1)"
              :disabled="index === categories.length - 1"
              aria-label="下移"
            >
              <DownOutlined />
            </button>
            <a-popconfirm
              title="删除该分类及所有标签？"
              @confirm="handleDeleteCategory(category.name)"
              ok-text="删除"
              cancel-text="取消"
              placement="left"
            >
              <button class="icon-btn delete" @click.stop aria-label="删除分类">
                <DeleteOutlined />
              </button>
            </a-popconfirm>
          </div>
        </div>
        <div v-if="categories.length === 0" class="empty-hint">
          暂无分类
        </div>
      </nav>
    </aside>

    <!-- 右侧：标签列表 -->
    <main class="panel-content">
      <header class="content-header">
        <div class="content-title-wrap" v-if="selectedCategory">
          <h2 class="content-title">{{ selectedCategory }}</h2>
          <span class="tag-count">{{ selectedTags.length }} 个标签</span>
        </div>
        <div class="content-title-wrap is-placeholder" v-else>
          <span class="placeholder-text">选择一个分类</span>
        </div>
        <a-button
          v-if="selectedCategory"
          type="primary"
          size="small"
          shape="round"
          @click="openAddTagModal"
        >
          <template #icon><PlusOutlined /></template>
          新增
        </a-button>
      </header>

      <section class="tag-list" role="list" v-if="selectedTags.length > 0">
        <div
          v-for="(tag, index) in selectedTags"
          :key="tag.name"
          class="tag-item"
          role="listitem"
        >
          <span class="tag-name">{{ tag.name }}</span>
          <span class="tag-record-count" v-if="tag.recordCount">{{ tag.recordCount }}</span>
          <button
            class="icon-btn"
            @click="moveTag(index, -1)"
            :disabled="index === 0"
            aria-label="上移"
          >
            <UpOutlined />
          </button>
          <button
            class="icon-btn"
            @click="moveTag(index, 1)"
            :disabled="index === selectedTags.length - 1"
            aria-label="下移"
          >
            <DownOutlined />
          </button>
          <a-popconfirm
            title="删除该标签？"
            @confirm="handleDeleteTag(tag.name)"
            ok-text="删除"
            cancel-text="取消"
            placement="left"
          >
            <button class="icon-btn delete" @click.stop aria-label="删除标签">
              <DeleteOutlined />
            </button>
          </a-popconfirm>
        </div>
      </section>

      <div class="empty-state" v-else-if="selectedCategory">
        <EmptyStateIcon />
        <p class="empty-text">该分类下暂无标签</p>
      </div>
      <div class="empty-state is-idle" v-else>
        <EmptyStateIcon />
        <p class="empty-text">从左侧选择分类查看标签</p>
      </div>
    </main>

    <!-- 添加分类弹窗 -->
    <a-modal
      v-model:open="openCategoryModal"
      title="新增分类"
      @ok="confirmAddCategory"
      ok-text="确认"
      cancel-text="取消"
      centered
      :width="360"
    >
      <div class="modal-form">
        <label class="form-label">分类名称</label>
        <a-input
          v-model:value="categoryForm.name"
          placeholder="请输入分类名称"
          size="large"
          :maxlength="20"
        />
      </div>
    </a-modal>

    <!-- 添加标签弹窗 -->
    <a-modal
      v-model:open="openTagModal"
      title="新增标签"
      @ok="confirmAddTag"
      ok-text="确认"
      cancel-text="取消"
      centered
      :width="360"
    >
      <div class="modal-form">
        <label class="form-label">标签名称</label>
        <a-input
          v-model:value="tagForm.name"
          placeholder="请输入标签名称"
          size="large"
          :maxlength="20"
        />
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch, h} from 'vue';
import type {TransactionType, Category, Tag} from '@/types/billadm';
import {useLedgerStore} from '@/stores/ledgerStore';
import {getCategoryByType, getTagsByCategory, addCategory, removeCategory, addTag, removeTag} from '@/backend/functions';
import {PlusOutlined, DeleteOutlined, UpOutlined, DownOutlined} from "@ant-design/icons-vue";
import {message} from "ant-design-vue";
import {reorderCategory, reorderTag} from '@/backend/functions';

// 空白状态图标组件
const EmptyStateIcon = () => h(
  'svg',
  {
    width: '48',
    height: '48',
    viewBox: '0 0 48 48',
    fill: 'none',
    xmlns: 'http://www.w3.org/2000/svg',
  },
  [
    h('rect', {
      x: '8',
      y: '12',
      width: '32',
      height: '24',
      rx: '4',
      stroke: 'currentColor',
      'stroke-width': '1.5',
      'stroke-dasharray': '4 2',
    }),
    h('line', {
      x1: '8',
      y1: '20',
      x2: '40',
      y2: '20',
      stroke: 'currentColor',
      'stroke-width': '1.5',
    }),
    h('line', {
      x1: '16',
      y1: '12',
      x2: '16',
      y2: '36',
      stroke: 'currentColor',
      'stroke-width': '1.5',
    }),
  ]
);

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
const categoryForm = ref({name: ''});

// 添加标签弹窗
const openTagModal = ref(false);
const tagForm = ref({name: ''});

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
    message.success('分类已添加');
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
    message.success('分类已删除');
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
    message.success('标签已添加');
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
    message.success('标签已删除');
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
  categories.value = categoryList.map(c => ({
    name: c.name,
    transactionType: c.transactionType,
    sortOrder: c.sortOrder,
    recordCount: c.recordCount,
    tags: []
  }));

  // 加载所有分类的标签
  for (const category of categories.value) {
    const categoryTransactionType = `${category.name}:${props.transactionType}`;
    const tags = await getTagsByCategory(categoryTransactionType, ledgerStore.currentLedgerId!);
    category.tags = tags.map(t => ({
      name: t.name,
      categoryTransactionType: t.categoryTransactionType,
      sortOrder: t.sortOrder,
      recordCount: t.recordCount
    }));
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
  background-color: var(--billadm-color-major-background);
}

/* ========== Sidebar ========== */
.panel-sidebar {
  width: 240px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background-color: var(--billadm-color-minor-background);
  border-right: 1px solid var(--billadm-color-divider);
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--billadm-space-xl) var(--billadm-space-xl) var(--billadm-space-md);
  /* Tighter bottom — hierarchy should come from vertical gap in the main layout */
  border-bottom: 1px solid var(--billadm-color-divider);
}

.sidebar-title {
  font-size: var(--billadm-size-text-body);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  margin: 0;
}

.category-list {
  flex: 1;
  overflow-y: auto;
  padding: var(--billadm-space-sm) 0;
}

.category-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  cursor: pointer;
  transition: background-color var(--billadm-transition-fast);
  position: relative;
}

.category-item:hover {
  background-color: var(--billadm-color-hover-bg);
}

.category-item.is-active {
  background-color: var(--billadm-color-active-bg);
}

.category-item.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background-color: var(--billadm-color-primary);
}

.category-info {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-sm);
  min-width: 0;
}

.category-name {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-major);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.category-count {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
  background-color: var(--billadm-color-divider);
  padding: 2px 6px;
  border-radius: var(--billadm-radius-full);
  flex-shrink: 0;
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

/* ========== Content ========== */
.panel-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  min-width: 0;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--billadm-space-xl) var(--billadm-space-xl) var(--billadm-space-md);
  border-bottom: 1px solid var(--billadm-color-divider);
  flex-shrink: 0;
  /* Tighter bottom — consistent with sidebar header pattern */
}

.content-title-wrap {
  display: flex;
  align-items: baseline;
  gap: var(--billadm-space-md);
}

.content-title {
  font-size: var(--billadm-size-text-title-sm);
  font-weight: var(--billadm-weight-semibold);
  color: var(--billadm-color-text-major);
  margin: 0;
}

.tag-count {
  font-size: var(--billadm-size-text-body-sm);
  color: var(--billadm-color-text-secondary);
  /* Body-sm matches sidebar title weight — clean reading level */
}

.content-title-wrap.is-placeholder .placeholder-text {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-disabled);
}

.tag-list {
  flex: 1;
  overflow-y: auto;
  padding: var(--billadm-space-md) var(--billadm-space-xl);
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-xs);
}

.tag-item {
  display: flex;
  align-items: center;
  gap: var(--billadm-space-md);
  padding: var(--billadm-space-md) var(--billadm-space-lg);
  border-radius: var(--billadm-radius-md);
  transition: background-color var(--billadm-transition-fast);
}

.tag-item:hover {
  background-color: var(--billadm-color-hover-bg);
}

.tag-name {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-major);
  flex: 1;
}

.tag-record-count {
  font-size: var(--billadm-size-text-caption);
  color: var(--billadm-color-text-secondary);
}

/* ========== Empty State ========== */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  gap: var(--billadm-space-md);
  color: var(--billadm-color-text-disabled);
}

.empty-state.is-idle {
  opacity: 0.6;
}

.empty-text {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin: 0;
}

.empty-hint {
  padding: var(--billadm-space-xl) var(--billadm-space-lg);
  text-align: center;
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-disabled);
}

/* ========== Icon Button ========== */
.icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  font-size: 14px;
  cursor: pointer;
  color: var(--billadm-color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--billadm-radius-md);
  transition: all var(--billadm-transition-fast);
}

.icon-btn:hover {
  color: var(--billadm-color-negative);
  background-color: rgba(199, 62, 58, 0.08);
}

.icon-btn.delete:hover {
  color: var(--billadm-color-negative);
  background-color: rgba(199, 62, 58, 0.08);
}

/* ========== Modal Form ========== */
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--billadm-space-sm);
}

.form-label {
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin-bottom: var(--billadm-space-xs);
}

/* ========== Scrollbar ========== */
.category-list::-webkit-scrollbar,
.tag-list::-webkit-scrollbar {
  width: 6px;
}

.category-list::-webkit-scrollbar-track,
.tag-list::-webkit-scrollbar-track {
  background: transparent;
}

.category-list::-webkit-scrollbar-thumb,
.tag-list::-webkit-scrollbar-thumb {
  background-color: var(--billadm-color-divider);
  border-radius: 3px;
}

.category-list::-webkit-scrollbar-thumb:hover,
.tag-list::-webkit-scrollbar-thumb:hover {
  background-color: var(--billadm-color-text-disabled);
}
</style>
