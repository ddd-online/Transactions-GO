import api from "@/backend/api/api-client";
import type { TransactionRecord, TrQueryCondition, TrQueryResult } from "@/types/billadm";

export async function queryTrOnCondition(condition: TrQueryCondition): Promise<TrQueryResult> {
    return api.post<TrQueryResult>('/v1/transactions/query', condition, '查询消费记录');
}

export async function createTrForLedger(data: TransactionRecord): Promise<string> {
    return api.post<string>('/v1/transactions', data, '创建消费记录');
}

export async function deleteTrById(id: string): Promise<void> {
    return api.delete<void>(`/v1/transactions/${id}`, '删除消费记录');
}
