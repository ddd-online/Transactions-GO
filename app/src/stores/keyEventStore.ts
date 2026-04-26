import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
    queryKeyEventDatesByYear,
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

    // 获取某年有记录的日期列表
    const fetchDatesByYear = async (year: number) => {
        try {
            const dates = await queryKeyEventDatesByYear(year);
            datesWithRecords.value = new Set(dates);
            currentYear.value = year;
        } catch (error) {
            NotificationUtil.error('查询关键事件失败', `${error}`);
        }
    };

    // 获取单日详情
    const fetchEventByDate = async (date: string): Promise<KeyEvent | null> => {
        try {
            return await queryKeyEventByDate(date);
        } catch (error) {
            // 404 表示当天没有记录，返回 null
            return null;
        }
    };

    // 保存事件（新建或更新）
    const saveEvent = async (date: string, content: string): Promise<void> => {
        try {
            await saveKeyEvent(date, content);
            datesWithRecords.value.add(date);
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

    return {
        datesWithRecords,
        currentYear,
        fetchDatesByYear,
        fetchEventByDate,
        saveEvent,
        deleteEvent,
        hasRecord,
    };
});
