import api from "@/backend/api/api-client";
import type { Tag } from "@/types/billadm";

export async function queryTags(category: string): Promise<Tag[]> {
    return api.post<Tag[]>(`/v1/tag/query/${category}`, {}, '查询标签');
}
