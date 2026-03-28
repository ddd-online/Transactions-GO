import dayjs from "dayjs";
import type {Category, Tag, TransactionRecord, TrQueryCondition, TrQueryResult} from "@/types/billadm";
import {createTrForLedger, deleteTrById, queryTrOnCondition} from "@/backend/api/tr.ts";
import NotificationUtil from "@/backend/notification.ts";
import {queryCategory} from "@/backend/api/category.ts";
import {queryTags} from "@/backend/api/tag.ts";

export function centsToYuan(cents: number): string {
    // 确保是整数
    if (!Number.isInteger(cents)) {
        console.warn('传入的不是整数分值:', cents);
    }
    // 转为带两位小数的字符串
    return (cents / 100).toLocaleString('zh-CN', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
        useGrouping: false
    });
}

export function yuanToCents(yuanStr: string): number {
    // 去除空格
    yuanStr = yuanStr.trim();

    // 支持负号
    const isNegative = yuanStr.startsWith('-');
    if (isNegative) yuanStr = yuanStr.slice(1);

    // 拆分整数和小数部分
    let [integerPart = '0', decimalPart = '00'] = yuanStr.split('.');

    // 小数部分最多取两位，不足补零，超过截断
    decimalPart = (decimalPart + '00').substring(0, 2);

    // 防止非数字字符
    if (!/^\d+$/.test(integerPart) || !/^\d{2}$/.test(decimalPart)) {
        throw new Error('无效的金额格式');
    }

    const totalCents = parseInt(integerPart, 10) * 100 + parseInt(decimalPart, 10);
    return isNegative ? -totalCents : totalCents;
}

/**
 * 将秒级时间戳转换为格式化时间字符串
 * @param timestamp 秒级时间戳
 * @param format 格式，默认为 'YYYY-MM-DD HH:mm:ss'
 * @returns 格式化后的时间字符串
 */
export function formatTimestamp(timestamp: number, format: string = 'YYYY-MM-DD'): string {
    return dayjs(timestamp * 1000).format(format);
}

/**
 * 消费记录
 */
export async function getTrOnCondition(condition: TrQueryCondition): Promise<TrQueryResult> {
    try {
        return await queryTrOnCondition(condition);
    } catch (error) {
        NotificationUtil.error('查询消费记录失败', `${error}`);
        return {
            items: [],
            total: 0,
            trStatistics: {
                income: 0,
                expense: 0,
                transfer: 0,
            }
        };
    }
}

export async function createTransactionRecord(tr: TransactionRecord) {
    try {
        await createTrForLedger(tr);
    } catch (error) {
        NotificationUtil.error('创建消费记录失败', `${error}`);
    }
}

export async function deleteTransactionRecord(trId: string) {
    try {
        await deleteTrById(trId);
    } catch (error) {
        NotificationUtil.error('删除消费记录失败', `${error}`);
    }
}

export async function updateTransactionRecord(tr: TransactionRecord) {
    try {
        await deleteTrById(tr.transactionId);
        await createTrForLedger(tr);
    } catch (error) {
        NotificationUtil.error('更新消费记录失败', `${error}`);
    }
}

/**
 * 分类与标签
 */
export async function getCategoryByType(trType: string): Promise<Category[]> {
    try {
        return await queryCategory(trType);
    } catch (error) {
        NotificationUtil.error(`查询 ${trType} 消费类型失败`, `${error}`);
        return [];
    }
}

export async function getTagsByCategory(category: string): Promise<Tag[]> {
    try {
        return await queryTags(category);
    } catch (error) {
        NotificationUtil.error(`查询 ${category} 消费标签失败`, `${error}`);
        return [];
    }
}
