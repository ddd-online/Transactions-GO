export {};

declare global {
    interface Window {
        electronAPI: {
            minimizeWindow: () => void;
            maximizeWindow: () => void;
            closeWindow: () => void;
            openDialog: (options: any) => Promise<any>;
            setWorkspace: (workspaceDir: string) => void;
            getAppInfo: (field: string) => Promise<any>;
        };
    }
}