import api from "@/backend/api/api-client";
import type { TransactionRecord, TrQueryCondition, TrQueryResult } from "@/types/billadm";

export async function queryTrOnCondition(condition: TrQueryCondition): Promise<TrQueryResult> {
    return api.post<TrQueryResult>('/v1/tr/query', condition, '查询消费记录');
}

export async function createTrForLedger(data: TransactionRecord): Promise<string> {
    return api.post<string>('/v1/tr/create-one', data, '创建消费记录');
}

export async function deleteTrById(id: string): Promise<void> {
    return api.post<void>('/v1/tr/delete-by-id', { trId: id }, '删除消费记录');
}
