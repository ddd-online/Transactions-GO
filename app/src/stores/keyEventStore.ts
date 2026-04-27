import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
    queryKeyEventsByYear,
    queryKeyEventByDate,
    saveKeyEvent,
    deleteKeyEvent
} from "@/backend/api/key-event";
import NotificationUtil from "@/backend/notification";
import type { KeyEvent } from "@/types/billadm";

export const useKeyEventStore = defineStore('keyEvent', () => {
    // 某年有记录的日期集合，用于日历高亮
    const datesWithRecords = ref(new Set<string>());
    const currentYear = ref(new Date().getFullYear());
    // 日期 -> 标题 的缓存
    const titles = ref(new Map<string, string>());

    // 获取某年有记录的日期列表
    const fetchDatesByYear = async (year: number) => {
        try {
            const events = await queryKeyEventsByYear(year);
            datesWithRecords.value = new Set(events.map(e => e.date));
            titles.value = new Map(events.map(e => [e.date, e.title]));
            currentYear.value = year;
        } catch (error) {
            NotificationUtil.error('查询关键事件失败', `${error}`);
        }
    };

    // 获取单日详情
    const fetchEventByDate = async (date: string): Promise<KeyEvent | null> => {
        try {
            const event = await queryKeyEventByDate(date);
            if (event) {
                titles.value.set(date, event.title);
            }
            return event;
        } catch (error) {
            // 404 表示当天没有记录，返回 null
            return null;
        }
    };

    // 保存事件（新建或更新）
    const saveEvent = async (date: string, title: string, content: string): Promise<void> => {
        try {
            await saveKeyEvent(date, title, content);
            datesWithRecords.value.add(date);
            titles.value.set(date, title);
            NotificationUtil.success('保存成功');
        } catch (error) {
            NotificationUtil.error('保存失败', `${error}`);
            throw error;
        }
    };

    // 删除事件
    const deleteEvent = async (date: string): Promise<void> => {
        try {
            await deleteKeyEvent(date);
            datesWithRecords.value.delete(date);
            titles.value.delete(date);
            NotificationUtil.success('删除成功');
        } catch (error) {
            NotificationUtil.error('删除失败', `${error}`);
            throw error;
        }
    };

    // 某天是否有记录
    const hasRecord = (date: string): boolean => {
        return datesWithRecords.value.has(date);
    };

    // 获取某天的标题
    const getTitle = (date: string): string => {
        return titles.value.get(date) || '';
    };

    return {
        datesWithRecords,
        currentYear,
        fetchDatesByYear,
        fetchEventByDate,
        saveEvent,
        deleteEvent,
        hasRecord,
        getTitle,
    };
});
