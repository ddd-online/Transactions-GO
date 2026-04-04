import type { TransactionRecord } from '@/types/billadm'
import type { TrQueryConditionItem } from '@/types/billadm'
import dayjs from 'dayjs'

/**
 * 按时间聚合的交易记录数据
 */
export interface TimeSeriesData {
  time: string
  type: string  // 存储transactionType: income/expense/transfer
  name: string  // 曲线名称，用于图例显示
  amount: number
}

/**
 * 图表曲线配置
 */
export interface ChartLine {
  name: string              // 曲线名称，用于图例显示
  transactionType: string   // 交易类型：income/expense/transfer
  includeOutlier: boolean  // 是否包含离群值
  conditions: TrQueryConditionItem[]  // 查询条件
}

/**
 * 图表配置项
 */
export interface ChartConfig {
  title: string
  granularity: 'year' | 'month'
  lines: ChartLine[]
}

/**
 * 根据条件过滤交易记录
 */
function filterByConditions(
  item: TransactionRecord,
  conditions: TrQueryConditionItem[]
): boolean {
  if (conditions.length === 0) {
    return true
  }
  // 所有条件必须同时满足
  return conditions.every((cond) => {
    if (cond.transactionType && item.transactionType !== cond.transactionType) {
      return false
    }
    if (cond.category && item.category !== cond.category) {
      return false
    }
    if (cond.tags && cond.tags.length > 0) {
      const hasAllTags = cond.tags.every((tag) => item.tags.includes(tag))
      if (cond.tagNot ? hasAllTags : !hasAllTags) {
        return false
      }
    }
    if (cond.description && !item.description.includes(cond.description)) {
      return false
    }
    return true
  })
}

/**
 * 构建G2折线图数据（支持多曲线）
 */
export function buildLineChartData(
  trList: TransactionRecord[],
  options: {
    granularity: 'year' | 'month'
    lines: ChartLine[]
  }
): TimeSeriesData[] {
  const { granularity, lines } = options

  if (lines.length === 0) {
    return []
  }

  // 收集所有曲线的时间范围
  let minYear = Infinity,
    minMonth = Infinity
  let maxYear = -Infinity,
    maxMonth = -Infinity

  lines.forEach((line) => {
    const filteredData = trList.filter((item) => {
      if (item.transactionType !== line.transactionType) return false
      if (!line.includeOutlier && item.outlier) return false
      return filterByConditions(item, line.conditions)
    })

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
  })

  // 如果没有数据，使用当前时间范围
  if (minYear === Infinity) {
    const now = dayjs()
    minYear = now.year()
    minMonth = now.month() + 1
    maxYear = minYear
    maxMonth = minMonth
  }

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

  // 初始化时间数据结构
  const timeDataMap = new Map<string, Record<string, number>>()
  timeLabels.forEach((label) => {
    const initObj: Record<string, number> = {}
    lines.forEach((line) => {
      initObj[line.name] = 0
    })
    timeDataMap.set(label, initObj)
  })

  // 聚合每条曲线的数据
  lines.forEach((line) => {
    const filteredData = trList.filter((item) => {
      if (item.transactionType !== line.transactionType) return false
      if (!line.includeOutlier && item.outlier) return false
      return filterByConditions(item, line.conditions)
    })

    filteredData.forEach((item) => {
      const date = dayjs(item.transactionAt * 1000)
      const label =
        granularity === 'year'
          ? String(date.year())
          : `${date.year()}-${String(date.month() + 1).padStart(2, '0')}`

      if (timeDataMap.has(label)) {
        const record = timeDataMap.get(label)!
        record[line.name] = (record[line.name] ?? 0) + item.price
      }
    })
  })

  // 转换为G2数据格式
  const result: TimeSeriesData[] = []
  timeLabels.forEach((label) => {
    const record = timeDataMap.get(label)!
    lines.forEach((line) => {
      result.push({
        time: label,
        type: line.transactionType,  // 使用transactionType而非label
        name: line.name,
        amount: (record[line.name] ?? 0) / 100, // 转换为元
      })
    })
  })

  return result
}

/**
 * 保留的图表配置
 */
export const KEEP_CHART_CONFIGS: ChartConfig[] = [
  {
    title: '月度消费趋势',
    granularity: 'month',
    lines: [
      { name: '支出', transactionType: 'expense', includeOutlier: false, conditions: [] },
      { name: '收入', transactionType: 'income', includeOutlier: false, conditions: [] },
      { name: '转账', transactionType: 'transfer', includeOutlier: false, conditions: [] },
    ],
  },
  {
    title: '年度消费趋势',
    granularity: 'year',
    lines: [
      { name: '支出', transactionType: 'expense', includeOutlier: true, conditions: [] },
      { name: '收入', transactionType: 'income', includeOutlier: true, conditions: [] },
      { name: '转账', transactionType: 'transfer', includeOutlier: true, conditions: [] },
    ],
  },
  {
    title: '年度总收入',
    granularity: 'year',
    lines: [
      { name: '年度总收入', transactionType: 'income', includeOutlier: true, conditions: [] },
    ],
  },
  {
    title: '年度工资收入',
    granularity: 'year',
    lines: [
      {
        name: '年度工资收入', transactionType: 'income', includeOutlier: true, conditions: [
          { transactionType: 'income', category: '工资奖金', tags: ['工资'], tagPolicy: 'all', tagNot: false, description: '' },
        ]
      },
    ],
  },
  {
    title: '年度奖金收入',
    granularity: 'year',
    lines: [
      {
        name: '年度奖金收入', transactionType: 'income', includeOutlier: true, conditions: [
          { transactionType: 'income', category: '工资奖金', tags: ['奖金'], tagPolicy: 'all', tagNot: false, description: '年奖金' },
        ]
      },
    ],
  },
  {
    title: '年度分红收入',
    granularity: 'year',
    lines: [
      {
        name: '年度分红收入', transactionType: 'income', includeOutlier: true, conditions: [
          { transactionType: 'income', category: '投资理财', tags: [], tagPolicy: 'all', tagNot: false, description: '年分红' },
        ]
      },
    ],
  },
]
