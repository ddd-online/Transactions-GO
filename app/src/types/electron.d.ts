export {};

declare global {
    interface Window {
        electronAPI: {
            minimizeWindow: () => void;
            maximizeWindow: () => void;
            closeWindow: () => void;
            openDialog: (options: any) => Promise<any>;
            saveDialog: (options: any) => Promise<any>;
            writeFile: (filePath: string, content: string) => Promise<{ success: boolean; error?: string }>;
            setWorkspace: (workspaceDir: string) => void;
            getAppInfo: (field: string) => Promise<any>;
        };
    }
}