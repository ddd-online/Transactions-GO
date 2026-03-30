<template>
  <a-config-provider :locale="locale" :theme="themeConfig">
    <router-view/>
  </a-config-provider>
</template>

<script setup lang="ts">
import zhCN from 'ant-design-vue/es/locale/zh_CN';
import {ref, computed} from 'vue';
import {useThemeStore} from '@/stores/themeStore';

const locale = ref(zhCN);
const themeStore = useThemeStore();

const themeConfig = computed(() => ({
  algorithm: themeStore.mode === 'dark' ? undefined : undefined,
  token: {
    colorPrimary: themeStore.mode === 'dark' ? '#7c9cff' : '#1677ff',
    colorBgContainer: themeStore.mode === 'dark' ? '#1a1a2e' : '#ffffff',
    colorBgElevated: themeStore.mode === 'dark' ? '#252538' : '#ffffff',
    colorBgLayout: themeStore.mode === 'dark' ? '#1a1a2e' : '#f5f5f5',
    colorBorder: themeStore.mode === 'dark' ? '#3a3a52' : '#e0e0e0',
    colorText: themeStore.mode === 'dark' ? '#e8e8f0' : '#333333',
    colorTextSecondary: themeStore.mode === 'dark' ? '#a0a0b8' : '#7b7b7b',
    borderRadius: 8,
  },
}));
</script>