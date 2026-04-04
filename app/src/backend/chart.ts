import type { TrQueryConditionItem } from '@/types/billadm'

/**
 * 按时间聚合的交易记录数据
 */
export interface TimeSeriesData {
  time: string
  type: string  // 存储transactionType: income/expense/transfer
  label: string  // 图例显示名称
  amount: number
}

/**
 * 图表曲线配置
 */
export interface ChartLine {
  label: string              // 图例显示名称
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
 * 保留的图表配置
 */
export const KEEP_CHART_CONFIGS: ChartConfig[] = [
  {
    title: '月度消费趋势',
    granularity: 'month',
    lines: [
      { label: '支出', transactionType: 'expense', includeOutlier: false, conditions: [] },
      { label: '收入', transactionType: 'income', includeOutlier: false, conditions: [] },
      { label: '转账', transactionType: 'transfer', includeOutlier: false, conditions: [] },
    ],
  },
  {
    title: '年度消费趋势',
    granularity: 'year',
    lines: [
      { label: '支出', transactionType: 'expense', includeOutlier: true, conditions: [] },
      { label: '收入', transactionType: 'income', includeOutlier: true, conditions: [] },
      { label: '转账', transactionType: 'transfer', includeOutlier: true, conditions: [] },
    ],
  },
  {
    title: '年度收入趋势',
    granularity: 'year',
    lines: [
      { label: '年度总收入', transactionType: 'income', includeOutlier: true, conditions: [] },
      {
        label: '年度工资收入', transactionType: 'income', includeOutlier: true, conditions: [
          { transactionType: 'income', category: '工资奖金', tags: ['工资'], tagPolicy: 'all', tagNot: false, description: '' },
        ]
      },
      {
        label: '年度奖金收入', transactionType: 'income', includeOutlier: true, conditions: [
          { transactionType: 'income', category: '工资奖金', tags: ['奖金'], tagPolicy: 'all', tagNot: false, description: '年奖金' },
        ]
      },
      {
        label: '年度分红收入', transactionType: 'income', includeOutlier: true, conditions: [
          { transactionType: 'income', category: '投资理财', tags: [], tagPolicy: 'all', tagNot: false, description: '年分红' },
        ]
      },
    ],
  },
]
