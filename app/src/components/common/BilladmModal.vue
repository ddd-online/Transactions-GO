<template>
  <a-modal
    v-model:open="open"
    :title="title"
    :width="width"
    :centered="centered"
    :closable="closable"
    :esc-to-close="escToClose"
    :mask-closable="maskClosable"
    @cancel="handleCancel"
  >
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>

    <slot />

    <template v-if="!hideFooter" #footer>
      <billadm-button @click="handleCancel">{{ cancelText }}</billadm-button>
      <billadm-button type="primary" @click="handleOk">{{ okText }}</billadm-button>
    </template>
  </a-modal>
</template>

<script setup lang="ts">
import { watch } from 'vue'

withDefaults(defineProps<{
  title?: string
  width?: string | number
  centered?: boolean
  closable?: boolean
  escToClose?: boolean
  maskClosable?: boolean
  hideFooter?: boolean
  cancelText?: string
  okText?: string
}>(), {
  width: 800,
  centered: true,
  closable: true,
  escToClose: true,
  maskClosable: false,
  hideFooter: false,
  cancelText: '取消',
  okText: '确认',
})

const emit = defineEmits<{
  (e: 'ok'): void
  (e: 'cancel'): void
  (e: 'update:open', val: boolean): void
}>()

const open = defineModel<boolean>()

watch(open, (val) => {
  if (val !== undefined) {
    emit('update:open', val)
  }
})

const handleOk = () => {
  emit('ok')
  open.value = false
}

const handleCancel = () => {
  emit('cancel')
  open.value = false
}
</script>

<style scoped>
</style>
