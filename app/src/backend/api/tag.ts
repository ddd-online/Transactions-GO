import api from "@/backend/api/api-client";
import type { Tag } from "@/types/billadm";

export async function queryTags(category: string): Promise<Tag[]> {
    return api.get<Tag[]>(`/v1/tags?category=${encodeURIComponent(category)}`, '查询标签');
}
