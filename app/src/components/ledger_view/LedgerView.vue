<template>
  <div class="ledger-view">
    <!-- 悬浮按钮 -->
    <a-float-button type="primary" class="float-primary" @click="openLedgerModal">
      <template #icon>
        <PlusOutlined />
      </template>
    </a-float-button>

    <!-- 账本卡片网格 -->
    <div class="ledger-grid">
      <a-card v-for="ledger in ledgerStore.ledgers" :key="ledger.id" hoverable>
        <a-descriptions :title="ledger.name" layout="vertical">
          <template #extra>
            <a-button type="text" class="btn-primary" @click="modifyLedgerName(ledger.id, ledger.name)">
              编辑
            </a-button>
            <a-popconfirm title="确认删除吗" ok-text="确认" :showCancel="false" @confirm="ledgerStore.deleteLedger(ledger.id)">
              <a-button type="text" class="btn-danger">删除</a-button>
            </a-popconfirm>
          </template>
          <a-descriptions-item :label="createTimeLabel">
            {{ formatTimestamp(ledger.createdAt, 'YYYY-MM-DD HH:mm:ss') }}
          </a-descriptions-item>
          <a-descriptions-item :label="updateTimeLabel">
            {{ formatTimestamp(ledger.updatedAt, 'YYYY-MM-DD HH:mm:ss') }}
          </a-descriptions-item>
        </a-descriptions>
      </a-card>
    </div>

    <!-- 编辑/新建弹窗 -->
    <a-modal :title="modalTitle" :open="ledgerModal" width="800px" @ok="confirmLedgerModal" ok-text="确认"
      @cancel="ledgerModal = false" cancel-text="取消" centered>
      <a-input v-model:value.lazy="ledgerName" placeholder="输入账本名称" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useLedgerStore } from "@/stores/ledgerStore";
import { formatTimestamp } from "@/backend/functions";
import { PlusOutlined } from "@ant-design/icons-vue";

const ledgerStore = useLedgerStore();

const ledgerModal = ref<boolean>(false);
const modalTitle = ref<string>("");
const ledgerId = ref<string>("");
const ledgerName = ref<string>("");
const createTimeLabel = '创建时间';
const updateTimeLabel = '更新时间';

const openLedgerModal = () => {
  modalTitle.value = "创建账本";
  ledgerName.value = "";
  ledgerModal.value = true;
};

const modifyLedgerName = (id: string, name: string) => {
  modalTitle.value = "修改账本名称";
  ledgerId.value = id;
  ledgerName.value = name;
  ledgerModal.value = true;
};

const confirmLedgerModal = async () => {
  if (!ledgerName) return;
  if (modalTitle.value === "创建账本") {
    await ledgerStore.createLedger(ledgerName.value);
  } else if (modalTitle.value === "修改账本名称") {
    await ledgerStore.modifyLedgerName(ledgerId.value, ledgerName.value);
  }
  ledgerModal.value = false;
};
</script>

<style scoped>
.ledger-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 16px;
  gap: 16px;
}

.ledger-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
  overflow-y: auto;
}

.float-primary {
  right: 50px;
  bottom: 80px;
}
</style>
