import type { TransactionRecord } from '@/types/billadm'
import type { TrQueryConditionItem } from '@/types/billadm'
import { TransactionTypeToLabel } from '@/backend/constant'
import dayjs from 'dayjs'

/**
 * 按时间聚合的交易记录数据
 */
export interface TimeSeriesData {
  time: string
  type: string
  amount: number
}

/**
 * 构建G2折线图数据
 */
export function buildLineChartData(
  trList: TransactionRecord[],
  options: {
    granularity: 'year' | 'month'
    lineDisplayTypes: string[]
    includeOutlier: boolean
  }
): TimeSeriesData[] {
  const { granularity, lineDisplayTypes, includeOutlier } = options

  // 过滤非异常且类型匹配的数据
  let filteredData = trList
    .filter((item) => lineDisplayTypes.includes(item.transactionType))

  if (!includeOutlier) {
    filteredData = filteredData.filter((item) => !item.outlier)
  }

  if (filteredData.length === 0) {
    return []
  }

  // 找出时间范围
  let minYear = Infinity,
    minMonth = Infinity
  let maxYear = -Infinity,
    maxMonth = -Infinity

  filteredData.forEach((item) => {
    const date = dayjs(item.transactionAt * 1000)
    const year = date.year()
    const month = date.month() + 1

    if (year < minYear || (year === minYear && month < minMonth)) {
      minYear = year
      minMonth = month
    }
    if (year > maxYear || (year === maxYear && month > maxMonth)) {
      maxYear = year
      maxMonth = month
    }
  })

  // 生成完整的时间轴（年 or 月）
  const timeLabels: string[] = []
  if (granularity === 'year') {
    for (let y = minYear; y <= maxYear; y++) {
      timeLabels.push(String(y))
    }
  } else {
    let currentYear = minYear
    let currentMonth = minMonth
    while (currentYear < maxYear || (currentYear === maxYear && currentMonth <= maxMonth)) {
      timeLabels.push(`${currentYear}-${String(currentMonth).padStart(2, '0')}`)
      currentMonth++
      if (currentMonth > 12) {
        currentMonth = 1
        currentYear++
      }
    }
  }

  // 初始化每月/每年数据结构
  const initData = () =>
    lineDisplayTypes.reduce(
      (acc, type) => {
        acc[type] = 0
        return acc
      },
      {} as Record<string, number>
    )

  const timeDataMap = new Map<string, Record<string, number>>()
  timeLabels.forEach((label) => {
    timeDataMap.set(label, initData())
  })

  // 聚合数据
  filteredData.forEach((item) => {
    const date = dayjs(item.transactionAt * 1000)
    const label =
      granularity === 'year'
        ? String(date.year())
        : `${date.year()}-${String(date.month() + 1).padStart(2, '0')}`

    const type = item.transactionType
    const amount = item.price

    if (timeDataMap.has(label) && lineDisplayTypes.includes(type)) {
      const record = timeDataMap.get(label)!
      record[type]! += amount
    }
  })

  // 转换为G2数据格式
  const result: TimeSeriesData[] = []
  timeLabels.forEach((label) => {
    const record = timeDataMap.get(label)!
    lineDisplayTypes.forEach((type) => {
      result.push({
        time: label,
        type: TransactionTypeToLabel.get(type) || type,
        amount: (record[type] ?? 0) / 100, // 转换为元
      })
    })
  })

  return result
}

/**
 * 图表配置项
 */
export interface ChartConfig {
  title: string
  granularity: 'year' | 'month'
  lineDisplayTypes: string[]
  includeOutlier: boolean
  conditions: TrQueryConditionItem[]
}

/**
 * 保留的图表配置
 */
export const KEEP_CHART_CONFIGS: ChartConfig[] = [
  {
    title: '月度消费趋势(不含离群值)',
    granularity: 'month',
    lineDisplayTypes: ['expense', 'income', 'transfer'],
    includeOutlier: false,
    conditions: [],
  },
  {
    title: '年度消费趋势(含离群值)',
    granularity: 'year',
    lineDisplayTypes: ['expense', 'income', 'transfer'],
    includeOutlier: true,
    conditions: [],
  },
  {
    title: '年度总收入',
    granularity: 'year',
    lineDisplayTypes: ['income'],
    includeOutlier: true,
    conditions: [],
  },
  {
    title: '年度工资收入',
    granularity: 'year',
    lineDisplayTypes: ['income'],
    includeOutlier: true,
    conditions: [{
      transactionType: 'income',
      category: '工资奖金',
      tags: ['工资'],
      tagPolicy: 'all',
      tagNot: false,
      description: '',
    }],
  },
  {
    title: '年度奖金收入',
    granularity: 'year',
    lineDisplayTypes: ['income'],
    includeOutlier: true,
    conditions: [{
      transactionType: 'income',
      category: '工资奖金',
      tags: ['奖金'],
      tagPolicy: 'all',
      tagNot: false,
      description: '年奖金',
    }],
  },
  {
    title: '年度分红收入',
    granularity: 'year',
    lineDisplayTypes: ['income'],
    includeOutlier: true,
    conditions: [{
      transactionType: 'income',
      category: '投资理财',
      tags: [],
      tagPolicy: 'all',
      tagNot: false,
      description: '年分红',
    }],
  },
]
