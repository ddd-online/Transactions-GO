export const TransactionTypeToLabel = new Map([
    ['income', '收入'],
    ['expense', '支出'],
    ['transfer', '转账']
]);

export const TransactionTypeToColor = new Map([
    ['income', '#52c41a'],
    ['expense', '#f5222d'],
    ['transfer', '#1677ff']
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
