<template>
  <button
    :class="[
      'billadm-btn',
      `billadm-btn-${type}`,
      `billadm-btn-${size}`,
      {
        'billadm-btn-danger': danger,
        'billadm-btn-icon-only': iconOnly,
        'billadm-btn-round': round,
      }
    ]"
    :disabled="disabled"
    @click="$emit('click', $event)"
  >
    <span v-if="$slots.icon" class="billadm-btn-icon">
      <slot name="icon" />
    </span>
    <span v-if="!iconOnly || !$slots.icon" class="billadm-btn-text">
      <slot />
    </span>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  type?: 'primary' | 'secondary' | 'text' | 'dashed'
  size?: 'small' | 'default'
  danger?: boolean
  disabled?: boolean
  iconOnly?: boolean
  round?: boolean
}>()

defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()
</script>

<style scoped>
.billadm-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border: 1px solid transparent;
  border-radius: var(--billadm-radius-sm);
  font-size: var(--billadm-size-text-body);
  cursor: pointer;
  transition: all var(--billadm-transition-normal);
  outline: none;
  white-space: nowrap;
}

.billadm-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

/* Sizes */
.billadm-btn-small {
  height: 24px;
  padding: 0 8px;
  font-size: var(--billadm-size-text-caption);
}

.billadm-btn-default {
  height: 32px;
  padding: 0 16px;
}

.billadm-btn-icon-only.billadm-btn-small {
  width: 24px;
  padding: 0;
}

.billadm-btn-icon-only.billadm-btn-default {
  width: 32px;
  padding: 0;
}

/* Primary */
.billadm-btn-primary {
  background-color: var(--billadm-color-primary);
  border-color: var(--billadm-color-primary);
  color: #fff;
}

.billadm-btn-primary:hover:not(:disabled) {
  background-color: var(--billadm-color-primary-light);
  border-color: var(--billadm-color-primary-light);
}

.billadm-btn-primary:active:not(:disabled) {
  background-color: var(--billadm-color-primary);
  border-color: var(--billadm-color-primary);
}

/* Secondary (Default) */
.billadm-btn-secondary {
  background-color: var(--billadm-color-major-background);
  border-color: var(--billadm-color-window-border);
  color: var(--billadm-color-text-major);
}

.billadm-btn-secondary:hover:not(:disabled) {
  background-color: var(--billadm-color-hover-bg);
  border-color: var(--billadm-color-primary);
  color: var(--billadm-color-primary);
}

/* Text */
.billadm-btn-text {
  background-color: transparent;
  border-color: transparent;
  color: var(--billadm-color-text-major);
}

.billadm-btn-text:hover:not(:disabled) {
  background-color: var(--billadm-color-hover-bg);
  color: var(--billadm-color-primary);
}

/* Dashed */
.billadm-btn-dashed {
  background-color: transparent;
  border-style: dashed;
  border-color: var(--billadm-color-window-border);
  color: var(--billadm-color-text-major);
}

.billadm-btn-dashed:hover:not(:disabled) {
  border-color: var(--billadm-color-primary);
  color: var(--billadm-color-primary);
}

/* Danger */
.billadm-btn-danger.billadm-btn-primary {
  background-color: var(--billadm-color-negative);
  border-color: var(--billadm-color-negative);
}

.billadm-btn-danger.billadm-btn-primary:hover:not(:disabled) {
  background-color: var(--billadm-color-negative);
  border-color: var(--billadm-color-negative);
  opacity: 0.85;
}

.billadm-btn-danger.billadm-btn-secondary,
.billadm-btn-danger.billadm-btn-text {
  color: var(--billadm-color-negative);
}

.billadm-btn-danger.billadm-btn-secondary:hover:not(:disabled),
.billadm-btn-danger.billadm-btn-text:hover:not(:disabled) {
  background-color: rgba(199, 62, 58, 0.1);
  border-color: var(--billadm-color-negative);
  color: var(--billadm-color-negative);
}

/* Round */
.billadm-btn-round {
  border-radius: 16px;
}

.billadm-btn-round.billadm-btn-small {
  border-radius: 12px;
}

/* Icon */
.billadm-btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.billadm-btn-icon-only .billadm-btn-icon {
  margin: 0;
}

.billadm-btn-text .billadm-btn-icon {
  margin-right: 0;
}
</style>
