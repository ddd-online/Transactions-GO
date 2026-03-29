<template>
  <div class="app-left-bar">
    <div class="app-left-bar-top">
      <a-button :type="route.path==='/ledger_view' ? 'primary': 'text'" @click="navigate('ledger_view')">
        <template #icon>
          <BookOutlined style="font-size: 18px"/>
        </template>
      </a-button>
      <a-button :type="route.path==='/tr_view' ? 'primary': 'text'" @click="navigate('tr_view')">
        <template #icon>
          <TransactionOutlined style="font-size: 18px"/>
        </template>
      </a-button>
      <a-button :type="route.path==='/da_view' ? 'primary': 'text'" @click="navigate('da_view')">
        <template #icon>
          <LineChartOutlined style="font-size: 18px"/>
        </template>
      </a-button>
    </div>
    <div class="app-left-bar-bottom">
      <a-button :type="route.path==='/settings_view' ? 'primary': 'text'" @click="navigate('settings_view')">
        <template #icon>
          <SettingOutlined style="font-size: 18px"/>
        </template>
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useRoute, useRouter} from 'vue-router';
import {BookOutlined, LineChartOutlined, SettingOutlined, TransactionOutlined} from "@ant-design/icons-vue";

const router = useRouter()
const route = useRoute()

const navigate = (path: string) => {
  router.push(path)
}
</script>

<style scoped>
.app-left-bar {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  width: 100%;
  padding: 8px 4px;
  background-color: var(--billadm-color-minor-background);
}

.app-left-bar-top {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.app-left-bar-bottom {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
