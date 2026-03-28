import api from "@/backend/api/api-client";
import type { Ledger } from "@/types/billadm";

export async function queryAllLedgers(): Promise<Ledger[]> {
    return api.get<Ledger[]>('/v1/ledgers?id=all', '查询账本');
}

export async function createLedgerByName(name: string): Promise<string> {
    return api.post<string>('/v1/ledgers', { name }, '创建账本');
}

export async function modifyLedgerNameById(id: string, name: string): Promise<void> {
    return api.patch<void>(`/v1/ledgers/${id}`, { name }, '修改账本名称');
}

export async function deleteLedgerById(id: string): Promise<void> {
    return api.delete<void>(`/v1/ledgers/${id}`, '删除账本');
}
