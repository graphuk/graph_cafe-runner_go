::change dir to the BAT directory
cd /D %~dp0

@echo off
set REPO=%CD%\..
set GOPATH=%REPO%\..\..\..\..
set GOROOT=%REPO%\node_modules\go-win
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%REPO%\node_modules\packr-win
set PATH=%PATH%;%REPO%\node_modules\upx-win
@echo on

:delete old binaries
del /F /S /Q %REPO%\cmd\cafe-runner-client\cafe-runner-client*
del /F /S /Q %REPO%\cmd\cafe-runner-server\cafe-runner-server*
del /F /S /Q %REPO%\web\assets\dist >nul 2>&1
rmdir /S /Q %REPO%\web\assets\dist >nul 2>&1
rmdir /S /Q %REPO%\release >nul 2>&1

mkdir %REPO%\web\assets\dist
mkdir %REPO%\web\assets\dist\win64
mkdir %REPO%\web\assets\dist\linux64
mkdir %REPO%\release

node %REPO%\node_modules\concurrently\bin\concurrently.js "%REPO%\scripts\_buildReleaseClientLinux.cmd" "%REPO%\scripts\_buildReleaseClientWin.cmd"
node %REPO%\node_modules\concurrently\bin\concurrently.js "%REPO%\scripts\_buildReleaseServerLinux.cmd" "%REPO%\scripts\_buildReleaseServerWin.cmd"

