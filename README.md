# Billadm-GO

一款桌面端记账工具

# 简介

账本管理页面
<div style="text-align: center;">
  <img alt="账本管理页面" style="width: 80%;" src="./doc/images/账本管理页面.png"/>
</div>

消费记录页面
<div style="text-align: center;">
  <img alt="消费记录页面" style="width: 80%;" src="./doc/images/消费记录页面.png"/>
</div>

数据分析页面
<div style="text-align: center;">
  <img alt="数据分析页面" style="width: 80%;" src="./doc/images/数据分析页面.png"/>
</div>

设置页面
<div style="text-align: center;">
  <img alt="设置页面" style="width: 80%;" src="./doc/images/设置页面.png"/>
</div>

# 安装

在release页面中下载安装包进行安装

> 仅支持windows-x86版本。其他架构需要自行编译

# 调试

### 自行打包

下载项目后，进入`build`目录，执行ps1脚本`build.ps1`

在`electron`目录下会生成`out`目录，其中存在打包后的安装程序

### 热更新调试

使用vue的热更新能力，需要以下三个步骤(建议使用goland)

1. goland打开`kernel`目录下的`main.go`文件，点击启动键，启动`go`服务
2. goland打开`app`目录下的`package.json`文件，点击启动键，启动`dev`命令
3. goland打开`electron`目录下的`package.json`文件，点击启动键，启动`start`命令
