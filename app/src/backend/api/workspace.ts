import api from "@/backend/api/api-client";

export async function openWorkspace(workspaceDir: string): Promise<void> {
    return api.post<void>('/v1/workspace', { workspaceDir }, '打开工作空间');
}
