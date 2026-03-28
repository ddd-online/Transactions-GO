import api from "@/backend/api/api-client";
import type { Ledger } from "@/types/billadm";

export async function queryAllLedgers(): Promise<Ledger[]> {
    return api.post<Ledger[]>('/v1/ledger/query-all', { id: 'all' }, '查询账本');
}

export async function createLedgerByName(name: string): Promise<string> {
    return api.post<string>('/v1/ledger/create-one', { name }, '创建账本');
}

export async function modifyLedgerNameById(id: string, name: string): Promise<void> {
    return api.post<void>('/v1/ledger/modify-name', { id, name }, '修改账本名称');
}

export async function deleteLedgerById(id: string): Promise<void> {
    return api.post<void>('/v1/ledger/delete-one', { id }, '删除账本');
}
