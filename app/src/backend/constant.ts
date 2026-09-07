import type { StockTradeTag } from '@/types/transactions';

export const TransactionTypeToLabel = new Map([
    ['income', '收入'],
    ['expense', '支出'],
    ['transfer', '转账']
]);

export const TransactionTypeToColor = new Map([
    ['income', '#16a34a'],
    ['expense', '#dc2626'],
    ['transfer', '#3b82f6']
]);

export const TimeRangeValueToLabel = {
    'date': '日',
    'month': '月',
    'year': '年'
} as const;

export const TimeRangeLabelToValue = {
    '日': 'date',
    '月': 'month',
    '年': 'year'
} as const;

/**
 * 股票交易轮次标签（策略分类）：展示值即存储值，默认「分析」。
 */
export const STOCK_TRADE_TAG_OPTIONS: { value: StockTradeTag; label: StockTradeTag }[] = [
    { value: '分析', label: '分析' },
    { value: '打板', label: '打板' },
    { value: '尾盘', label: '尾盘' },
    { value: '追涨', label: '追涨' },
];

/** 股票交易轮次标签默认值 */
export const STOCK_TRADE_TAG_DEFAULT: StockTradeTag = '分析';
