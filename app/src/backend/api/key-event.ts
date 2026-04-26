import api from "@/backend/api/api-client";
import type { KeyEvent } from "@/types/billadm";

export async function queryKeyEventDatesByYear(year: number): Promise<string[]> {
    return api.get<string[]>(`/v1/key-events/dates/${year}`, '查询关键事件日期列表');
}

export async function queryKeyEventByDate(date: string): Promise<KeyEvent> {
    return api.get<KeyEvent>(`/v1/key-events/${date}`, '查询关键事件详情');
}

export async function saveKeyEvent(date: string, content: string): Promise<string> {
    return api.post<string>('/v1/key-events', { date, content }, '保存关键事件');
}

export async function deleteKeyEvent(date: string): Promise<void> {
    return api.delete<void>(`/v1/key-events/${date}`, '删除关键事件');
}
