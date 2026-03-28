<template>
  <a-select @change="handleChange" style="width: 150px" :value="ledgerStore.currentLedgerName">
    <a-select-option
        v-for="option in options"
        :key="option.value"
        :value="option.value">
      {{ option.label }}
    </a-select-option>
  </a-select>
</template>

<script setup lang="ts">
import {computed} from 'vue';
import {useLedgerStore} from '@/stores/ledgerStore.ts';
import type {Ledger} from '@/types/billadm';
import type {SelectValue} from "ant-design-vue/es/select";
import type {DefaultOptionType} from "ant-design-vue/es/vc-cascader";

const ledgerStore = useLedgerStore();

const options = computed(() => {
  if (!Array.isArray(ledgerStore.ledgers)) return [];
  return ledgerStore.ledgers.map((ledger: Ledger) => ({
    label: ledger.name,
    value: ledger.id,
  }));
});

const handleChange = (value: SelectValue, _: DefaultOptionType | DefaultOptionType[]) => {
  ledgerStore.setCurrentLedger(value as string);
};
</script>
