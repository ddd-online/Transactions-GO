import api_client from "@/backend/api/api-client.ts";
import type {Result, WorkspaceStatus} from "@/types/billadm";

export async function openWorkspace(workspaceDir: string) {
    const resp: Result = await api_client.post('/v1/workspace/open', {
        'workspaceDir': workspaceDir,
    });
    api_client.isRespSuccess(resp, 'openWorkspace错误: ');
}

export async function hasOpenedWorkspace(): Promise<WorkspaceStatus> {
    const resp: Result<string> = await api_client.post('/v1/workspace/is-opened');
    api_client.isRespSuccess(resp, 'hasOpenedWorkspace错误: ');
    return {
        isOpened: resp.data.length > 0,
        workspaceDir: resp.data,
    };
}