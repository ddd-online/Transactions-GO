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
  token: {
    colorPrimary: '#2D5A27',
    colorBgContainer: themeStore.mode === 'dark' ? '#1A1A18' : '#FFFFFF',
    colorBgElevated: themeStore.mode === 'dark' ? '#252522' : '#FFFFFF',
    colorBgLayout: themeStore.mode === 'dark' ? '#1A1A18' : '#FAFAF8',
    colorBorder: themeStore.mode === 'dark' ? '#3A3A38' : '#E5E5E0',
    colorText: themeStore.mode === 'dark' ? '#F5F5F3' : '#1A1A18',
    colorTextSecondary: themeStore.mode === 'dark' ? '#A0A09A' : '#6B6B66',
    borderRadius: 6,
    fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
  },
}));
</script>
