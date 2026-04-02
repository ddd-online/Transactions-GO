import api from "@/backend/api/api-client";
import type { Tag } from "@/types/billadm";

export async function queryTags(categoryTransactionType: string): Promise<Tag[]> {
    return api.get<Tag[]>(`/v1/tags?categoryTransactionType=${encodeURIComponent(categoryTransactionType)}`, '查询标签');
}

export async function createTag(name: string, categoryTransactionType: string): Promise<void> {
    await api.post<void>('/v1/tags', { name, categoryTransactionType }, '创建标签');
}

export async function deleteTag(name: string, categoryTransactionType: string, ledgerId: string): Promise<void> {
    await api.delete<void>(`/v1/tags/${encodeURIComponent(name)}?categoryTransactionType=${encodeURIComponent(categoryTransactionType)}&ledgerId=${encodeURIComponent(ledgerId)}`, '删除标签');
}
