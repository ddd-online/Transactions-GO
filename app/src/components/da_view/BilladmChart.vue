<template>
  <div ref="containerRef" class="billadm-chart"></div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch, nextTick } from 'vue'
import { Chart } from '@antv/g2'
import type { TimeSeriesData } from '@/backend/chart'
import { TransactionTypeToColor, TransactionTypeToLabel } from '@/backend/constant'

interface Props {
  data: TimeSeriesData[]
  xField: string
  yField: string
  seriesField: string
  title: string
}

const props = defineProps<Props>()

const containerRef = ref<HTMLDivElement | null>(null)
let chart: Chart | null = null

const initChart = () => {
  if (!containerRef.value || !props.data.length) return

  // 销毁旧图表
  if (chart) {
    chart.destroy()
    chart = null
  }

  // 获取时间轴标题
  const xAxisTitle = props.xField === 'time' ? (props.title.includes('月度') ? '月份' : '年份') : props.xField

  // 计算16:9的高度
  const width = containerRef.value.clientWidth
  const height = width * 9 / 16

  chart = new Chart({
    container: containerRef.value,
    autoFit: true,
    height: height,
    data: props.data,
  })

  // 使用TransactionTypeToColor和TransactionTypeToLabel设置颜色
  // 数据中type字段使用的是label（收入/支出/转账），需要统一颜色映射
  const colorDomain = Array.from(TransactionTypeToLabel.values()) // ['收入', '支出', '转账']
  const colorRange = Array.from(TransactionTypeToLabel.keys()).map(k => TransactionTypeToColor.get(k)!) // ['#52c41a', '#f5222d', '#1677ff']

  chart.scale('color', {
    domain: colorDomain,
    range: colorRange,
  })

  chart.axis('x', {
    title: xAxisTitle,
    labelFill: '#000000',
    labelFontSize: 15,
    titleFontSize: 16,
    line: { style: { stroke: '#000000', lineWidth: 1 } }
  })
  chart.axis('y', {
    title: '金额（元）',
    labelFill: '#000000',
    labelFontSize: 15,
    titleFontSize: 16,
    domainMin: 0,
    nice: true,
    line: { style: { stroke: '#000000', lineWidth: 1 } }
  })

  chart
    .line()
    .encode('x', props.xField)
    .encode('y', props.yField)
    .encode('color', props.seriesField)
    .style('lineWidth', 2)

  chart
    .point()
    .encode('x', props.xField)
    .encode('y', props.yField)
    .encode('color', props.seriesField)
    .style('size', 4)
    .style('stroke', '#fff')
    .style('lineWidth', 1)
    .tooltip(false)

  chart.render()
}

onMounted(async () => {
  await nextTick()
  initChart()
})

onUnmounted(() => {
  if (chart) {
    chart.destroy()
    chart = null
  }
})

watch(() => props.data, () => {
  initChart()
}, { deep: true })
</script>

<style scoped>
.billadm-chart {
  width: 100%;
  height: 100%;
}
</style>
