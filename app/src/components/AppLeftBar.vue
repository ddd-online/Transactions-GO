<template>
  <div class="app-left-bar">
    <div class="app-left-bar-top">
      <button
        class="nav-btn"
        :class="{ active: route.path === '/ledger_view' }"
        @click="navigate('ledger_view')"
        aria-label="账本"
        title="账本"
      >
        <BookOutlined style="font-size: 20px"/>
      </button>
      <button
        class="nav-btn"
        :class="{ active: route.path === '/tr_view' }"
        @click="navigate('tr_view')"
        aria-label="交易记录"
        title="交易记录"
      >
        <TransactionOutlined style="font-size: 20px"/>
      </button>
      <button
        class="nav-btn"
        :class="{ active: route.path === '/da_view' }"
        @click="navigate('da_view')"
        aria-label="数据分析"
        title="数据分析"
      >
        <LineChartOutlined style="font-size: 20px"/>
      </button>
    </div>
    <div class="app-left-bar-bottom">
      <button
        class="nav-btn"
        :class="{ active: route.path === '/settings_view' }"
        @click="navigate('settings_view')"
        aria-label="设置"
        title="设置"
      >
        <SettingOutlined style="font-size: 20px"/>
      </button>
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
  padding: 12px 8px;
}

.app-left-bar-top {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.app-left-bar-bottom {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.nav-btn {
  position: relative;
  width: 40px;
  height: 40px;
  border: none;
  background: none;
  border-radius: var(--billadm-radius-md);
  color: var(--billadm-color-icon);
  cursor: pointer;
  transition: background-color var(--billadm-transition-fast),
              color var(--billadm-transition-fast);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-btn:hover {
  background-color: var(--billadm-color-hover-bg);
  color: var(--billadm-color-primary);
}

.nav-btn:focus-visible {
  outline: 2px solid var(--billadm-color-primary);
  outline-offset: 2px;
}

.nav-btn.active {
  background-color: var(--billadm-color-active-bg);
  color: var(--billadm-color-primary);
}
</style>
