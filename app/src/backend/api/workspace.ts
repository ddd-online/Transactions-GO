import api from "@/backend/api/api-client";
import type { WorkspaceStatus } from "@/types/billadm";

export async function openWorkspace(workspaceDir: string): Promise<void> {
    return api.post<void>('/v1/workspace/open', { workspaceDir }, '打开工作空间');
}

export async function hasOpenedWorkspace(): Promise<WorkspaceStatus> {
    const data = await api.post<string>('/v1/workspace/is-opened', {}, '检查工作空间');
    return {
        isOpened: data.length > 0,
        workspaceDir: data,
    };
}
