const {contextBridge, ipcRenderer} = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
    minimizeWindow: () => {
        ipcRenderer.send('window-control', 'minimize');
    },
    maximizeWindow: () => {
        ipcRenderer.send('window-control', 'maximize');
    },
    closeWindow: () => {
        ipcRenderer.send('window-control', 'close');
    },
    openDialog: async (options) => {
        return await ipcRenderer.invoke('dialog:open', options);
    },
    saveDialog: async (options) => {
        return await ipcRenderer.invoke('dialog:save', options);
    },
    writeFile: async (filePath, content) => {
        return await ipcRenderer.invoke('file:write', { filePath, content });
    },
    setWorkspace: (workspaceDir) => {
        ipcRenderer.send('workspace:set', workspaceDir);
    },
    getWorkspace: async () => {
        return await ipcRenderer.invoke('workspace:get');
    },
    getAppInfo: async (field) => {
        return await ipcRenderer.invoke('app', field);
    },
    getApiServer: async () => {
        return await ipcRenderer.invoke('app', 'apiServer');
    },
});