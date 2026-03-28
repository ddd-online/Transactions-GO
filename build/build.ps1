# build.ps1 - 一键构建整个 Billadm 应用（Vue + Go + Electron）

# 设置输出编码为 UTF-8（防止中文乱码）
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# 获取脚本所在目录
$scriptDir = $PSScriptRoot

# 获取项目根目录（上一级）
$projectRoot = Split-Path -Parent $scriptDir

# 定义全局路径
$vueDir = Join-Path $projectRoot "app"
$kernelDir = Join-Path $projectRoot "kernel"
$electronDir = Join-Path $projectRoot "electron"

$appDistDir = Join-Path $vueDir "dist"
$kernelExe = Join-Path $kernelDir "Billadm-Kernel.exe"
$sqlFile = Join-Path $kernelDir "billadm.sql"

# 颜色辅助函数（可选，提升可读性）
function Write-Info { param($msg) Write-Host "📦 $msg" -ForegroundColor Cyan }
function Write-Step { param($msg) Write-Host "`n🛠️  $msg" -ForegroundColor Magenta }
function Write-Success { param($msg) Write-Host "✅ $msg" -ForegroundColor Green }
function Write-Warn { param($msg) Write-Host "⚠️  $msg" -ForegroundColor Yellow }
function Write-ErrorCustom { param($msg) Write-Error "❌ $msg" }

# 记录初始位置，用于最后返回
$initialLocation = Get-Location

try {
    # ==============================
    # 1. 构建 Vue 前端
    # ==============================
    Write-Step "正在构建前端 Vue 项目..."

    if (-not (Test-Path $vueDir)) {
        Write-ErrorCustom "Vue 项目目录不存在: $vueDir"
        exit 1
    }

    $packageJson = Join-Path $vueDir "package.json"
    if (-not (Test-Path $packageJson)) {
        Write-ErrorCustom "未找到 package.json，确认 '$vueDir' 是有效的 Vue 项目"
        exit 1
    }

    if (Test-Path $appDistDir) {
        Write-Warn "正在删除旧的 dist 目录..."
        Remove-Item $appDistDir -Recurse -Force -ErrorAction Stop
        Write-Success "成功删除旧的 dist 目录"
    } else {
        Write-Host "🔍 dist 目录不存在，跳过删除。" -ForegroundColor Cyan
    }

    Set-Location $vueDir
    Write-Host "   执行命令: npm run build`n" -ForegroundColor DarkGray
    & npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-ErrorCustom "Vue 构建失败，退出码: $LASTEXITCODE"
        exit $LASTEXITCODE
    }
    Write-Success "前端构建成功！输出位于: $appDistDir"


    # ==============================
    # 2. 构建 Go 后端
    # ==============================
    Write-Step "正在构建后端 Go 项目..."

    if (-not (Test-Path $kernelDir)) {
        Write-ErrorCustom "Go 项目目录不存在: $kernelDir"
        exit 1
    }

    $goMod = Join-Path $kernelDir "go.mod"
    if (-not (Test-Path $goMod)) {
        Write-ErrorCustom "未找到 go.mod，确认 '$kernelDir' 是有效的 Go 项目"
        exit 1
    }

    if (Test-Path $kernelExe) {
        Write-Warn "正在删除旧的编译文件..."
        Remove-Item $kernelExe -Force -ErrorAction Stop
        Write-Success "成功删除 $kernelExe"
    }

    Set-Location $kernelDir
    Write-Host "`n   设置 GOOS=windows, GOARCH=amd64, CGO_ENABLED=1" -ForegroundColor DarkGray
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "1"

    & go build -ldflags '-s -w -extldflags "-static"' -o $kernelExe
    if ($LASTEXITCODE -ne 0) {
        Write-ErrorCustom "Go 编译失败，退出码: $LASTEXITCODE"
        exit $LASTEXITCODE
    }

    if (-not (Test-Path $kernelExe)) {
        Write-ErrorCustom "编译完成但未生成预期文件: $kernelExe"
        exit 1
    }
    Write-Success "后端编译成功！生成文件: $kernelExe"


    # ==============================
    # 3. 准备 Electron 打包环境
    # ==============================
    Write-Step "准备 Electron 打包环境..."

    if (-not (Test-Path $electronDir)) {
        Write-ErrorCustom "Electron 项目目录不存在: $electronDir"
        exit 1
    }

    # 清理 electron 目录中的旧资源
    $targetDist = Join-Path $electronDir "dist"
    $targetKernel = Join-Path $electronDir "Billadm-Kernel.exe"
    $targetSql = Join-Path $electronDir "billadm.sql"

    foreach ($item in @($targetDist, $targetKernel, $targetSql)) {
        if (Test-Path $item) {
            Write-Warn "正在删除旧文件/目录: $(Split-Path $item -Leaf)"
            Remove-Item $item -Recurse -Force -ErrorAction Stop
        }
    }

    # 拷贝前端 dist
    if (-not (Test-Path $appDistDir)) {
        Write-ErrorCustom "前端构建产物缺失: $appDistDir"
        exit 1
    }
    Copy-Item -Path $appDistDir -Destination $electronDir -Recurse -Force -ErrorAction Stop
    Write-Success "前端 dist 已拷贝至 $electronDir\dist"

    # 拷贝后端 exe
    if (-not (Test-Path $kernelExe)) {
        Write-ErrorCustom "后端可执行文件缺失: $kernelExe"
        exit 1
    }
    Copy-Item -Path $kernelExe -Destination $electronDir -Force -ErrorAction Stop
    Write-Success "已拷贝 Billadm-Kernel.exe 到 $electronDir"

    # 拷贝 SQL 文件
    if (-not (Test-Path $sqlFile)) {
        Write-ErrorCustom "SQL 初始化文件缺失: $sqlFile"
        exit 1
    }
    Copy-Item -Path $sqlFile -Destination $electronDir -Force -ErrorAction Stop
    Write-Success "已拷贝 billadm.sql 到 $electronDir"


    # ==============================
    # 4. 执行 Electron 打包
    # ==============================
    Write-Step "正在执行 Electron 应用打包..."
    Set-Location $electronDir
    Write-Host "   执行命令: npm run package" -ForegroundColor Yellow
    & npm run package
    if ($LASTEXITCODE -ne 0) {
        Write-ErrorCustom "Electron 打包失败，退出码: $LASTEXITCODE"
        exit $LASTEXITCODE
    }
    Write-Success "Electron 应用打包成功！"


    # ==============================
    # 完成
    # ==============================
    Write-Host "`n🎉 整个构建与打包流程已完成！" -ForegroundColor Green

} finally {
    # 确保最终返回脚本目录
    Set-Location $initialLocation
    Write-Host "`n↩️  已返回脚本所在目录: $scriptDir" -ForegroundColor DarkCyan
}