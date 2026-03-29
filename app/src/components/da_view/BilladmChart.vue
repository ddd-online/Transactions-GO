<template>
  <div ref="containerRef" class="billadm-chart"></div>
</template>

<script setup lang="ts">
import {onMounted, onUnmounted, ref, watch} from 'vue'
import {Chart} from '@antv/g2'
import type {TimeSeriesData} from '@/backend/chart'
import {TransactionTypeToColor} from '@/backend/constant'

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

  chart = new Chart({
    container: containerRef.value,
    autoFit: true,
    height: 300,
    data: props.data,
  })

  // 使用TransactionTypeToColor设置颜色
  const colorDomain = Array.from(TransactionTypeToColor.keys())
  const colorRange = colorDomain.map(k => TransactionTypeToColor.get(k)!)

  chart.scale(props.seriesField, {
    domain: colorDomain,
    range: colorRange,
  })

  chart.axis(props.xField, {
    title: {text: xAxisTitle, fontSize: 12},
    labelFontSize: 12,
  })
  chart.axis(props.yField, {
    title: {text: '金额（元）', fontSize: 12},
    labelFontSize: 12,
    labelFormatter: (value: string) => `¥${parseFloat(value).toFixed(0)}`,
  })

  chart
    .line()
    .encode('x', props.xField)
    .encode('y', props.yField)
    .encode('color', props.seriesField)
    .encode('shape', 'smooth')
    .style('lineWidth', 2)

  chart
    .point()
    .encode('x', props.xField)
    .encode('y', props.yField)
    .encode('color', props.seriesField)
    .style('size', 4)
    .style('stroke', '#fff')
    .style('lineWidth', 1)

  chart.render()
}

onMounted(() => {
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
}, {deep: true})
</script>

<style scoped>
.billadm-chart {
  width: 100%;
  height: 100%;
  min-height: 300px;
}
</style>
