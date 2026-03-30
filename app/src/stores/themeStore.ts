import {defineStore} from 'pinia';
import {ref, watch} from 'vue';

export type ThemeMode = 'light' | 'dark';

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>('light');

  const toggleTheme = () => {
    mode.value = mode.value === 'light' ? 'dark' : 'light';
  };

  const setTheme = (theme: ThemeMode) => {
    mode.value = theme;
  };

  // Apply theme to document
  watch(mode, (newTheme) => {
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('billadm-theme', newTheme);
  }, {immediate: true});

  // Load saved theme
  const loadSavedTheme = () => {
    const saved = localStorage.getItem('billadm-theme') as ThemeMode | null;
    if (saved === 'light' || saved === 'dark') {
      mode.value = saved;
    }
  };

  return {
    mode,
    toggleTheme,
    setTheme,
    loadSavedTheme,
  };
});
