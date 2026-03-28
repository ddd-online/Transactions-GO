import axios, { type AxiosInstance, type AxiosResponse } from 'axios';
import type { Result } from "@/types/billadm";

const apiClient: AxiosInstance = axios.create({
    baseURL: 'http://127.0.0.1:31943/api',
    timeout: 10000,
    headers: { 'Content-Type': 'application/json' },
});

/**
 * Check if the response indicates an error (code !== 0).
 * Throws an Error with the message if so.
 */
function checkSuccess(result: Result, prefix?: string): void {
    if (result.code !== 0) {
        throw new Error(`${prefix || ''}响应失败: ${result.msg}`);
    }
}

export const api = {
    async get<T = any>(url: string, errorPrefix?: string): Promise<T> {
        try {
            const response: AxiosResponse<Result<T>> = await apiClient.get(url);
            checkSuccess(response.data, errorPrefix);
            return response.data.data;
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(`${errorPrefix || '请求失败'}: ${error.message}`);
            }
            throw error;
        }
    },

    async post<T = any>(url: string, data: object = {}, errorPrefix?: string): Promise<T> {
        try {
            const response: AxiosResponse<Result<T>> = await apiClient.post(url, data);
            checkSuccess(response.data, errorPrefix);
            return response.data.data;
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(`${errorPrefix || '请求失败'}: ${error.message}`);
            }
            throw error;
        }
    },

    async patch<T = any>(url: string, data: object = {}, errorPrefix?: string): Promise<T> {
        try {
            const response: AxiosResponse<Result<T>> = await apiClient.patch(url, data);
            checkSuccess(response.data, errorPrefix);
            return response.data.data;
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(`${errorPrefix || '请求失败'}: ${error.message}`);
            }
            throw error;
        }
    },

    async delete<T = any>(url: string, errorPrefix?: string): Promise<T> {
        try {
            const response: AxiosResponse<Result<T>> = await apiClient.delete(url);
            checkSuccess(response.data, errorPrefix);
            return response.data.data;
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(`${errorPrefix || '请求失败'}: ${error.message}`);
            }
            throw error;
        }
    }
};

export default api;
