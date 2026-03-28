import type {Dayjs} from "dayjs";

/**
 * 表示一个前端使用的消费记录
 */
export interface TrForm {
    id: string;
    price: string;
    type: string;
    category: string;
    description: string;
    tags: string[];
    flags: string[];
    time: Dayjs;
}

/**
 * 后端返回的响应的规范结构
 */
export interface Result<T = any> {
    code: number;
    msg: string;
    data: T;
}

export interface TrQueryResult {
    items: TransactionRecord[];
    total: number;
    trStatistics: TrStatistics;
}

/**
 * 账本
 */
export interface Ledger {
    id: string;           // 账本UUID
    name: string;         // 账本名称
    createdAt: number;   // 创建时间（Unix 时间戳，单位秒）
    updatedAt: number;   // 更新时间（Unix 时间戳，单位秒）
}

/**
 * 消费记录
 */
export interface TransactionRecord {
    ledgerId: string;
    transactionId: string;
    price: number;
    transactionType: string;
    category: string;
    description: string;
    tags: string[];
    transactionAt: number;
    outlier: boolean;
}

/**
 * 消费类型
 */
export interface Category {
    name: string;
    transactionType: string;
}

/**
 * 消费标签
 */
export interface Tag {
    name: string;      // 标签名称
    scope: string;     // 作用域
    category: string;  // 分类ID
}

/**
 * 消费记录统计数据
 */
export interface TrStatistics {
    income: number;    // 收入金额
    expense: number;   // 支出金额
    transfer: number;  // 转账金额
}

/**
 * 消费记录条件查询
 */
export interface TrQueryCondition {
    ledgerId: string;
    tsRange?: number[];
    items?: TrQueryConditionItem[];
    offset?: number;
    limit?: number;
}

/**
 * 消费记录条件项
 */
export interface TrQueryConditionItem {
    transactionType: string;
    category: string;
    tags: string[];
    tagPolicy: string;
    tagNot: boolean;
    description: string;
}

/**
 * 工作空间状态
 */
export interface WorkspaceStatus {
    isOpened: boolean;
    workspaceDir: string;
}

/**
 * 时间范围类型 时间范围标签类型 时间范围值类型
 */
type RangeValue = [Dayjs, Dayjs] | undefined;
type TimeRangeTypeValue = 'date' | 'month' | 'year';
type TimeRangeTypeLabel = '日' | '月' | '年';

type TransactionType = 'income' | 'expense' | 'transfer';