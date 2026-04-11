import api from "@/backend/api/api-client";

export interface McpStatus {
    running: boolean;
    port: number;
}

export async function startMcpServer(): Promise<McpStatus> {
    return api.post<McpStatus>('/v1/mcp/start', {}, '启动MCP服务');
}

export async function stopMcpServer(): Promise<McpStatus> {
    return api.post<McpStatus>('/v1/mcp/stop', {}, '停止MCP服务');
}

export async function getMcpStatus(): Promise<McpStatus> {
    return api.get<McpStatus>('/v1/mcp/status', '获取MCP状态');
}
