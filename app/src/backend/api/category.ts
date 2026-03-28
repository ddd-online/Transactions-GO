import api from "@/backend/api/api-client";
import type { Category } from "@/types/billadm";

export async function queryCategory(trType: string): Promise<Category[]> {
    return api.get<Category[]>(`/v1/categories?type=${encodeURIComponent(trType)}`, '查询分类');
}
