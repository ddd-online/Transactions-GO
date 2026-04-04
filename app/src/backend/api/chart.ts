import api from "@/backend/api/api-client";
import type { ChartLine } from '@/backend/chart';

export interface ChartDto {
    chartId: string
    ledgerId: string
    title: string
    granularity: 'year' | 'month'
    lines: ChartLine[]
    chartType: string
    isPreset: boolean
    sortOrder: number
}

export interface CreateChartRequest {
    ledgerId: string
    title: string
    granularity: 'year' | 'month'
    lines: ChartLine[]
    chartType: string
}

export interface UpdateChartRequest {
    chartId: string
    title: string
    granularity: 'year' | 'month'
    lines: ChartLine[]
    chartType: string
    sortOrder: number
}

export async function queryCharts(ledgerId: string): Promise<ChartDto[]> {
    return api.get<ChartDto[]>(`/v1/charts?ledgerId=${ledgerId}`, '查询图表列表');
}

export async function createChart(request: CreateChartRequest): Promise<ChartDto> {
    return api.post<ChartDto>('/v1/charts', request, '创建图表');
}

export async function updateChart(request: UpdateChartRequest): Promise<ChartDto> {
    return api.patch<ChartDto>('/v1/charts', request, '更新图表');
}

export async function deleteChart(chartId: string): Promise<void> {
    return api.delete<void>(`/v1/charts/${chartId}`, '删除图表');
}
