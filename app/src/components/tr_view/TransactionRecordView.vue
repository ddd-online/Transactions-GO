<template>
  <a-float-button
      type="primary"
      style="right: 50px;bottom: 80px"
      @click="createTr">
    <template #icon>
      <PlusOutlined/>
    </template>
  </a-float-button>
  <a-float-button
      style="right: 110px;bottom: 80px"
      @click="openTrFilterModal=true"
      :badge="{ count: trQueryConditionStore.conditionLen, color: 'blue' }"
  >
    <template #icon>
      <FilterOutlined/>
    </template>
  </a-float-button>
  <TransactionRecordFilter v-model="openTrFilterModal"/>
  <a-layout style="height: 100%">
    <a-layout>
      <a-layout-header class="headerStyle">
        <div class="left-groups">
          <BilladmTimeRangePicker
              v-model:time-range="trQueryConditionStore.timeRange"
              v-model:time-range-type="trQueryConditionStore.timeRangeType"
          />
        </div>
        <div class="center-groups">
        </div>
        <div class="right-groups">
        </div>
      </a-layout-header>
      <a-layout-content :style="contentStyle">
        <transaction-record-table :items="tableData" @edit="updateTr" @delete="deleteTr"/>
      </a-layout-content>
      <a-layout-footer class="footerStyle">
        <a-pagination
            v-model:current="currentPage"
            v-model:pageSize="pageSize"
            :total="trTotal"
            :show-total="total => `共${total}记录`"
            :pageSizeOptions="['15','30','50','100']"
            show-size-changer
        />
      </a-layout-footer>
      <a-modal
          :title="trModalTitle"
          :open="openTrModal"
          width="800px"
          @ok="confirmTrModal"
          ok-text="确认"
          @cancel="closeTrModal"
          cancel-text="取消"
          centered
      >
        <a-form :model="trForm" :rules="rules">
          <a-form-item label="时间" name="time">
            <a-date-picker v-model:value="trForm.time" style="width: 100%"/>
          </a-form-item>

          <a-form-item label="类型" name="type">
            <a-radio-group v-model:value="trForm.type" button-style="solid">
              <a-radio-button value="income">收入</a-radio-button>
              <a-radio-button value="expense">支出</a-radio-button>
              <a-radio-button value="transfer">转账</a-radio-button>
            </a-radio-group>
          </a-form-item>

          <a-form-item label="分类" name="category">
            <a-select v-model:value="trForm.category" :options="categories"/>
          </a-form-item>

          <a-form-item label="标签" name="tags">
            <a-select v-model:value="trForm.tags" :options="tags" mode="multiple" placeholder="选择一个或多个标签"/>
          </a-form-item>

          <a-form-item label="标记" name="flags">
            <a-checkbox-group v-model:value="trForm.flags" :options="flagOptions"/>
          </a-form-item>

          <a-form-item label="描述" name="description">
            <a-input v-model:value="trForm.description" placeholder="描述消费内容" allowClear/>
          </a-form-item>

          <a-form-item label="金额" name="price">
            <a-input v-model:value="trForm.price" prefix="￥" style="width: 100%"/>
          </a-form-item>
        </a-form>
      </a-modal>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import {type CSSProperties, ref, watch} from 'vue';
import TransactionRecordTable from '@/components/tr_view/TransactionRecordTable.vue';
import type {TransactionRecord, TrForm, TrQueryCondition} from "@/types/billadm";
import {useCssVariables} from "@/backend/css.ts";
import {convertToUnixTimeRange} from "@/backend/timerange.ts";
import {
  createTransactionRecord,
  deleteTransactionRecord,
  getCategoryByType,
  getTagsByCategory,
  getTrOnCondition,
  updateTransactionRecord
} from "@/backend/functions.ts";
import {useLedgerStore} from "@/stores/ledgerStore.ts";
import {useTrQueryConditionStore} from "@/stores/trQueryConditionStore.ts";
import dayjs from "dayjs";
import {trDtoToTrForm, trFormToTrDto} from "@/backend/dto-utils.ts";
import type {DefaultOptionType} from "ant-design-vue/es/vc-cascader";
import {useAppDataStore} from "@/stores/appDataStore.ts";
import type {Rule} from "ant-design-vue/es/form";
import {FilterOutlined, PlusOutlined} from "@ant-design/icons-vue";

const {majorBgColor} = useCssVariables();

const contentStyle: CSSProperties = {
  backgroundColor: majorBgColor.value,
  "overflow-y": "auto",
  "margin-bottom": "auto"
};

const ledgerStore = useLedgerStore();
const trQueryConditionStore = useTrQueryConditionStore();
const appDataStore = useAppDataStore();

// 表单校验规则
const rules: Record<string, Rule[]> = {
  price: [
    {
      trigger: 'blur',
    },
    {
      validator: (_: any, value: string) => {
        if (!value) return Promise.reject(new Error('请输入价格'));

        // 必须是非负数，且最多两位小数
        const regex = /^(0|[1-9]\d*)(\.\d{1,2})?$/;
        // 允许 "0", "0.00", "123", "123.4", "123.45"
        // 不允许 ".5", "01.23", "123.456"

        if (!regex.test(value)) {
          return Promise.reject(new Error('请输入 ≥0 的有效金额，最多两位小数'));
        }

        return Promise.resolve();
      },
      trigger: 'blur',
    },
  ],
};

// modal
const openTrFilterModal = ref<boolean>()
// 消费记录
const tableData = ref<TransactionRecord[]>([]);
// 分页
const currentPage = ref<number>(1);
const pageSize = ref<number>(15);
const trTotal = ref<number>(0);

const openTrModal = ref(false);
const trModalTitle = ref('');
const trForm = ref<TrForm>({
  id: '',
  price: '',
  type: '',
  category: '',
  description: '',
  tags: [],
  flags: [],
  time: dayjs()
});
const categories = ref<DefaultOptionType[]>([]);
const tags = ref<DefaultOptionType[]>([]);
const flagOptions = [
  {label: '离群值', value: 'outlier'}
]

const createTr = () => {
  trForm.value.type = 'expense';
  if (trQueryConditionStore.timeRange) {
    trForm.value.time = trQueryConditionStore.timeRange[1];
  }
  trModalTitle.value = '新增消费记录';
  openTrModal.value = true;
}

const updateTr = (tr: TransactionRecord) => {
  trModalTitle.value = '编辑消费记录';
  trForm.value = trDtoToTrForm(tr);
  openTrModal.value = true;
}

const deleteTr = async (tr: TransactionRecord) => {
  await deleteTransactionRecord(tr.transactionId);
  await refreshTable();
}

const closeTrModal = () => {
  trForm.value = {
    id: '',
    price: '',
    type: '',
    category: '',
    description: '',
    tags: [],
    flags: [],
    time: dayjs()
  };
  openTrModal.value = false;
}

const confirmTrModal = async () => {
  trForm.value.time = trForm.value.time.hour(12).minute(0).second(0);
  const tr = trFormToTrDto(trForm.value, ledgerStore.currentLedgerId);
  if (tr.transactionId === '') {
    // 新建
    if (!tr.description) tr.description = '-';
    await createTransactionRecord(tr);
  } else {
    // 更新
    await updateTransactionRecord(tr);
  }
  await refreshTable();
  closeTrModal();
}

const refreshTable = async () => {
  if (!ledgerStore.currentLedgerId) return;

  const trCondition: TrQueryCondition = {
    ledgerId: ledgerStore.currentLedgerId,
    offset: pageSize.value * (currentPage.value - 1),
    limit: pageSize.value
  };
  if (trQueryConditionStore.timeRange) {
    trCondition.tsRange = convertToUnixTimeRange(trQueryConditionStore.timeRange);
  }
  if (trQueryConditionStore.trQueryConditionItems) {
    trCondition.items = trQueryConditionStore.trQueryConditionItems
  }
  let trQueryResult = await getTrOnCondition(trCondition);
  tableData.value = trQueryResult.items
  trTotal.value = trQueryResult.total
  appDataStore.setStatistics(trQueryResult.trStatistics);
}

// 查询条件变化 → 重置分页 + 刷新
watch(() => [
      ledgerStore.currentLedgerId,
      trQueryConditionStore.timeRange,
      trQueryConditionStore.trQueryConditionItems
    ], async () => {
      if (currentPage.value !== 1) {
        currentPage.value = 1;
        return;
      }
      await refreshTable();
    },
    {immediate: true}
);

// 分页变化 → 仅刷新
watch(() => [currentPage.value, pageSize.value], async () => {
  await refreshTable();
});
/**
 * 交易类型变化时要重新刷新分类列表并如果当前分类为空则选择第一个分类作为分类
 * 如果当前分类不为空则选择查看分类是否在列表中，不在列表中则需要选择第一个分类作为分类
 */
watch(() => trForm.value.type, async () => {
      if (trForm.value.type === '') return;
      const categoryList = await getCategoryByType(trForm.value.type);
      categories.value = categoryList.map(category => {
        return {
          value: category.name,
        };
      });
      const categoryNames = categoryList.map(category => category.name);
      if (categoryNames.length > 0) {
        if (!trForm.value.category || !categoryNames.includes(trForm.value.category)) {
          trForm.value.category = categoryNames[0] as string;
        }
      } else {
        trForm.value.category = '';
      }
    }
);
/**
 * 分类变化时要重新刷新标签列表清除不在候选中的标签
 */
watch(() => trForm.value.category, async () => {
      if (trForm.value.category === '') return;
      const tagList = await getTagsByCategory(trForm.value.category);
      tags.value = tagList.map(tag => {
        return {
          value: tag.name,
        };
      });
      const tagNames = tagList.map(tag => tag.name);
      if (tagNames.length > 0 && trForm.value.tags) {
        let newTags: string[] = [];
        trForm.value.tags.forEach(tag => {
          if (tag && tagNames.includes(tag)) {
            newTags.push(tag);
          }
        });
        trForm.value.tags = newTags;
      } else {
        trForm.value.tags = [];
      }
    }
);
</script>

<style scoped>
.headerStyle {
  height: auto;
  background-color: var(--billadm-color-major-background);
  padding: 0 0 16px 0;
  display: flex;
  align-items: start;
  justify-content: center;
}

.footerStyle {
  height: auto;
  background-color: var(--billadm-color-major-background);
  padding: 16px 0 0 0;
  display: flex;
  align-items: end;
  justify-content: center;
}

/* 左边按钮 将它与后面的元素隔开 */
.left-groups {
  margin-right: auto;
  display: flex;
  gap: 8px;
  align-items: center;
}

/* 中间按钮 */
.center-groups {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  gap: 8px;
}

/* 右边按钮组 */
.right-groups {
  display: flex;
  gap: 8px;
}
</style>