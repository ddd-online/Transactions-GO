const {app, BrowserWindow, ipcMain, dialog, net} = require('electron');
const path = require('path');
const fs = require('fs');
const os = require('os');

process.noAsar = false;

const isDev = !app.isPackaged;
const appPath = isDev ? path.dirname(__dirname) : app.getAppPath();

const getServer = () => {
    return "http://127.0.0.1:31943"
};

const getUiServer = () => {
    if (isDev) {
        return "http://localhost:31945/static"
    } else {
        return "http://127.0.0.1:31943/static/index.html"
    }
}


// 应用日志
const logDir = path.join(appPath, 'logs');
const logFile = path.join(logDir, 'app.log');

if (!fs.existsSync(logDir)) {
    fs.mkdirSync(logDir, {recursive: true});
}
const logStream = fs.createWriteStream(logFile, {flags: 'a'});
const log = (message) => {
    const time = new Date().toISOString();
    const logMessage = `[${time}] ${message}\n`;
    logStream.write(logMessage);
};

let billadmCfg = {
    width: 1400, height: 1000, x: undefined, y: undefined, workspaceDir: '',
};

function readBilladmFile() {
    if (isDev) return;
    const homeDir = os.homedir();
    const filePath = path.join(homeDir, '.billadm');
    let tmpObj;
    try {
        const fileContent = fs.readFileSync(filePath, 'utf8');
        tmpObj = JSON.parse(fileContent);
        billadmCfg = {
            ...billadmCfg, ...tmpObj,
        }
    } catch (err) {
        log(`读取 .billadm 文件时发生错误:', ${err.message}`);
    }

    log(`窗口宽度 ${billadmCfg.width} 窗口高度 ${billadmCfg.height} 工作空间路径 ${billadmCfg.workspaceDir}`)
}

function saveBilladmConfig() {
    if (isDev) return;
    const homeDir = os.homedir();
    const filePath = path.join(homeDir, '.billadm');

    try {
        if (typeof billadmCfg !== 'object' || billadmCfg === null) {
            log('billadmCfg 不是一个有效的对象，无法保存');
            return false;
        }

        const data = JSON.stringify(billadmCfg, null, 2);

        fs.writeFileSync(filePath, data, 'utf8');

        log(`配置已成功保存至: ${filePath}`);
    } catch (err) {
        log(`保存 .billadm 文件时发生错误 ${err.message}`);
    }
}


// 内核
let kernelProcess = null;

const startKernel = () => {
    if (isDev) return;
    const kernelExe = path.join(appPath, 'Billadm-Kernel.exe');
    log(`Starting kernel: ${kernelExe}`);
    const cp = require("child_process");
    kernelProcess = cp.spawn(kernelExe, ['-mode', 'release', '-workspace', billadmCfg.workspaceDir], {
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
        log(`[Kernel Process] kernel [pid=${kernelProcess.pid}] closed with code ${code}`);
        kernelProcess = null;
    });

    kernelProcess.on('error', (err) => {
        log('[Kernel Process] Failed to start:', err);
    });
};

const createWindow = () => {
    const mainWindow = new BrowserWindow({
        width: billadmCfg.width,
        height: billadmCfg.height,
        x: billadmCfg.x,
        y: billadmCfg.y,
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
                    await net.fetch(getServer() + "/api/v1/app/exit", {method: "POST"});
                } catch (e) {
                    log(`请求kernel关闭失败 ${e}`)
                }
                const bounds = mainWindow.getBounds();
                billadmCfg = {...billadmCfg, ...bounds}
                mainWindow.close();
                break;
        }
    });

    ipcMain.handle('dialog:open', async (event, options) => {
        try {
            return await dialog.showOpenDialog({
                properties: ['openFile'], ...options,
            });
        } catch (err) {
            log(`Dialog error: ${err.message}`);
            return {canceled: true, filePaths: [], error: err.message};
        }
    });

    ipcMain.on('workspace:set', (event, workspaceDir) => {
        billadmCfg.workspaceDir = workspaceDir;
    });

    ipcMain.handle('app', async (event, field) => {
        switch (field) {
            case 'name':
                return app.getName();
            case 'version':
                return app.getVersion();
            default:
                return '';
        }
    });
};

app.whenReady().then(() => {
    readBilladmFile();
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
        if (kernelProcess && kernelProcess.exitCode === null) {
            kernelProcess.kill();
        }
        saveBilladmConfig();
        app.quit();
    }
});
