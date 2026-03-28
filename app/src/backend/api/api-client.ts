import axios, {type AxiosInstance, type AxiosResponse} from 'axios';
import type {ApiClient, Result} from "@/types/billadm";

const apiClient: AxiosInstance = axios.create({
    baseURL: 'http://127.0.0.1:31943/api',
    timeout: 10000, // 请求超时时间 (毫秒)
    headers: {
        'Content-Type': 'application/json',
    },
});

const api: ApiClient = {
    async post(url: string, data: object = {}): Promise<Result> {
        try {
            const response: AxiosResponse<Result> = await apiClient.post(url, data);
            return response.data;
        } catch (error) {
            console.error('API POST Error:', error);
            throw error;
        }
    },
    async get(url: string): Promise<Result> {
        try {
            const response: AxiosResponse = await apiClient.get(url);
            return response.data;
        } catch (error) {
            console.error('API GET Error:', error);
            throw error;
        }
    },
    isRespSuccess(result: Result, prefix?: string): void {
        if (result.code !== 0) {
            throw `${prefix}响应code不为0, 响应msg: ${result.msg}`;
        }
    }
};

export default api;