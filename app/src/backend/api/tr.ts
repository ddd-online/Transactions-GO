import api_client from "@/backend/api/api-client.ts";
import type {Result, TransactionRecord, TrQueryCondition, TrQueryResult} from "@/types/billadm";

export async function queryTrOnCondition(condition: TrQueryCondition): Promise<TrQueryResult> {
    const resp: Result<TrQueryResult> = await api_client.post('/v1/tr/query', condition);
    api_client.isRespSuccess(resp, 'queryTrsOnCondition错误: ');
    return resp.data;
}

export async function createTrForLedger(data: TransactionRecord): Promise<string> {
    const resp: Result<string> = await api_client.post('/v1/tr/create-one', data);
    api_client.isRespSuccess(resp, 'createTrForLedger错误: ');
    return resp.data;
}

export async function deleteTrById(id: string) {
    const resp: Result = await api_client.post('/v1/tr/delete-by-id', {
        'trId': id,
    });
    api_client.isRespSuccess(resp, 'deleteTrById错误: ');
}