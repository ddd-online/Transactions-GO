import api from "@/backend/api/api-client";
import type { Category } from "@/types/billadm";

export async function queryCategory(trType: string): Promise<Category[]> {
    return api.post<Category[]>(`/v1/category/query/${trType}`, {}, '查询分类');
}
