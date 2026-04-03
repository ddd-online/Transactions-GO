const { app, BrowserWindow, ipcMain, dialog, net } = require('electron');
const path = require('path');
const fs = require('fs');
const os = require('os');

process.noAsar = false;

const isDev = !app.isPackaged;
const appPath = isDev ? path.dirname(__dirname) : app.getAppPath();

const API_PORT = isDev ? '28080' : '31943';
const API_SERVER = `http://127.0.0.1:${API_PORT}`;

const getUiServer = () => {
    if (isDev) {
        return 'http://localhost:31945/static';
    } else {
        return `${API_SERVER}/static/index.html`;
    }
};

// 应用日志
const logDir = path.join(appPath, 'logs');
const logFile = path.join(logDir, 'app.log');

if (!fs.existsSync(logDir)) {
    fs.mkdirSync(logDir, { recursive: true });
}
const logStream = fs.createWriteStream(logFile, { flags: 'a' });
const log = (message) => {
    const time = new Date().toISOString();
    logStream.write(`[${time}] ${message}\n`);
};

let transactionsCfg = {
    width: 1400, height: 1000, x: undefined, y: undefined, workspaceDir: '',
};

function readTransactionsCfg() {
    if (isDev) return;
    const homeDir = os.homedir();
    const filePath = path.join(homeDir, '.transactions.json');
    try {
        const fileContent = fs.readFileSync(filePath, 'utf8');
        const tmpObj = JSON.parse(fileContent);
        transactionsCfg = { ...transactionsCfg, ...tmpObj };
    } catch (err) {
        log(`读取 .transactions.json 文件失败: ${err.message}`);
    }
    log(`窗口 ${transactionsCfg.width}x${transactionsCfg.height} workspace ${transactionsCfg.workspaceDir}`);
}

function saveTransactionsCfg() {
    if (isDev) return;
    const homeDir = os.homedir();
    const filePath = path.join(homeDir, '.transactions.json');
    try {
        if (typeof transactionsCfg !== 'object' || transactionsCfg === null) {
            log('transactionsCfg 无效，无法保存');
            return;
        }
        fs.writeFileSync(filePath, JSON.stringify(transactionsCfg, null, 2), 'utf8');
        log(`配置已保存至 ${filePath}`);
    } catch (err) {
        log(`保存配置失败: ${err.message}`);
    }
}


// 内核
let kernelProcess = null;

const startKernel = () => {
    if (isDev) return;
    const kernelExe = path.join(appPath, 'Billadm-Kernel.exe');
    log(`Starting kernel: ${kernelExe}`);
    const cp = require("child_process");
    kernelProcess = cp.spawn(kernelExe, ['-mode', 'release', '-port', API_PORT, '-workspace', transactionsCfg.workspaceDir], {
        detached: false,
    });

    // 捕获标准输出
    kernelProcess.stdout.on('data', (data) => {
        log(`[Kernel STDOUT]: ${data.toString()}`);
    });

    // 捕获错误输出
    kernelProcess.stderr.on('data', (data) => {
        log(`[Kernel STDERR]: ${data.toString()}`);
    });

    // 进程关闭
    kernelProcess.on('close', (code) => {
        if (kernelProcess) {
            log(`[Kernel Process] kernel [pid=${kernelProcess.pid}] closed with code ${code}`);
        } else {
            log(`[Kernel Process] kernel closed with code ${code}`);
        }
        kernelProcess = null;
    });

    // 进程异常退出
    kernelProcess.on('exit', (code) => {
        const pid = kernelProcess ? kernelProcess.pid : 'unknown';
        log(`[Kernel Process] kernel [pid=${pid}] exited with code ${code}`);
        if (code !== 0 && code !== null) {
            dialog.showMessageBox({
                type: 'error',
                title: '后台服务异常退出',
                message: `后台服务异常退出，退出码: ${code}\n请重启应用`,
            });
        }
        kernelProcess = null;
    });

    kernelProcess.on('error', (err) => {
        log('[Kernel Process] Failed to start:', err);
    });
};

const createWindow = () => {
    const mainWindow = new BrowserWindow({
        width: transactionsCfg.width,
        height: transactionsCfg.height,
        x: transactionsCfg.x,
        y: transactionsCfg.y,
        frame: false,
        webPreferences: {
            nodeIntegration: false, contextIsolation: true, preload: path.join(__dirname, 'preload.js'),
        },
    });

    mainWindow.loadURL(getUiServer());

    if (isDev) {
        mainWindow.webContents.openDevTools();
    }

    ipcMain.on('window-control', async (event, command) => {
        switch (command) {
            case 'minimize':
                mainWindow.minimize();
                break;
            case 'maximize':
                mainWindow.isMaximized() ? mainWindow.unmaximize() : mainWindow.maximize();
                break;
            case 'close':
                try {
                    await net.fetch(API_SERVER + "/api/v1/app/exit", { method: "POST" });
                } catch (e) {
                    log(`请求kernel关闭失败 ${e}`);
                }
                const bounds = mainWindow.getBounds();
                transactionsCfg = { ...transactionsCfg, ...bounds }
                mainWindow.close();
                break;
        }
    });

    ipcMain.handle('dialog:open', async (event, options) => {
        try {
            return await dialog.showOpenDialog({
                properties: ['openDirectory'], ...options,
            });
        } catch (err) {
            log(`Dialog error: ${err.message}`);
            return { canceled: true, filePaths: [], error: err.message };
        }
    });

    ipcMain.handle('dialog:save', async (event, options) => {
        try {
            return await dialog.showSaveDialog({
                properties: ['showOverwriteConfirmation'], ...options,
            });
        } catch (err) {
            log(`Save Dialog error: ${err.message}`);
            return { canceled: true, filePath: '', error: err.message };
        }
    });

    ipcMain.handle('file:write', async (event, { filePath, content }) => {
        try {
            fs.writeFileSync(filePath, content, 'utf8');
            log(`File written: ${filePath}`);
            return { success: true };
        } catch (err) {
            log(`File write error: ${err.message}`);
            return { success: false, error: err.message };
        }
    });

    ipcMain.on('workspace:set', (event, workspaceDir) => {
        transactionsCfg.workspaceDir = workspaceDir;
        saveTransactionsCfg();
    });

    ipcMain.handle('app', async (event, field) => {
        switch (field) {
            case 'name':
                return app.getName();
            case 'version':
                return app.getVersion();
            case 'apiServer':
                return API_SERVER;
            default:
                return '';
        }
    });
};

app.whenReady().then(() => {
    readTransactionsCfg();
    startKernel();
    createWindow();

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow();
        }
    });
});

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        if (kernelProcess) {
            kernelProcess.kill();
        }
        saveTransactionsCfg();
        app.quit();
    }
});
