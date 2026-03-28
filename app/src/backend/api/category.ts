import api_client from "@/backend/api/api-client.ts";
import type {Category, Result} from "@/types/billadm";

export async function queryCategory(trType: string): Promise<Category[]> {
    const resp: Result<Category[]> = await api_client.post(`/v1/category/query/${trType}`);
    api_client.isRespSuccess(resp, 'queryCategory错误: ');
    return resp.data;
}