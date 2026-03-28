<template>
  <div class="menu-bar">
    <div class="top-groups">
      <a-button :type="route.path==='/ledger_view' ? 'primary': 'text'" @click="navigate('ledger_view')">
        <template #icon>
          <BookOutlined style="display: flex;justify-content: center;align-items: center;font-size: medium"/>
        </template>
      </a-button>
      <a-button :type="route.path==='/tr_view' ? 'primary': 'text'" @click="navigate('tr_view')">
        <template #icon>
          <TransactionOutlined style="display: flex;justify-content: center;align-items: center;font-size: medium"/>
        </template>
      </a-button>
      <a-button :type="route.path==='/da_view' ? 'primary': 'text'" @click="navigate('da_view')">
        <template #icon>
          <LineChartOutlined style="display: flex;justify-content: center;align-items: center;font-size: medium"/>
        </template>
      </a-button>
    </div>
    <div class="bottom-groups">
      <a-button :type="route.path==='/settings_view' ? 'primary': 'text'" @click="navigate('settings_view')">
        <template #icon>
          <SettingOutlined style="display: flex;justify-content: center;align-items: center;font-size: medium"/>
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
.menu-bar {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
}

.top-groups {
  margin-bottom: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bottom-groups {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>