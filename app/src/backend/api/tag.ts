import api from "@/backend/api/api-client";
import type { Tag } from "@/types/billadm";

export async function queryTags(categoryTransactionType: string): Promise<Tag[]> {
    return api.get<Tag[]>(`/v1/tags?categoryTransactionType=${encodeURIComponent(categoryTransactionType)}`, '查询标签');
}
