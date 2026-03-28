<!-- components/BilladmFullscreen.vue -->
<template>
  <!-- 非全屏状态：正常渲染插槽内容 -->
  <div v-if="!isFullscreen" class="billadm-fullscreen-wrapper" @dblclick.stop="handleDblClick">
    <slot/>
  </div>

  <!-- 全屏状态：使用 teleport 渲染到 body -->
  <teleport to="body" :disabled="!isFullscreen">
    <div v-if="isFullscreen" class="billadm-fullscreen-wrapper fullscreen">
      <div class="fullscreen-mask" @click.self="toggleFullscreen"></div>
      <div class="fullscreen-panel">
        <slot/>
      </div>
    </div>
  </teleport>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  dblclick: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue']);

const isFullscreen = ref(props.modelValue);

watch(
    () => props.modelValue,
    (newVal) => {
      isFullscreen.value = newVal;
    }
);

watch(isFullscreen, (val) => {
  emit('update:modelValue', val);
});

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value;
};

const handleDblClick = () => {
  if (props.dblclick) {
    toggleFullscreen();
  }
};
</script>

<style scoped>
.billadm-fullscreen-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.billadm-fullscreen-wrapper.fullscreen {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}

.fullscreen-panel {
  width: 80%;
  height: auto;
  max-width: 1200px;
  background: var(--billadm-color-major-background-color);
  border: 1px solid var(--billadm-color-window-border-color);
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.fullscreen-mask {
  position: absolute;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: -1;
}
</style>