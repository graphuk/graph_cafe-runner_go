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
del /F /S /Q %REPO%\cmd\cafe-runner-server\cafe-runner-server
del /F /S /Q %REPO%\release\cafe-runner-server

:mkdir %REPO%\web\assets\dist
:mkdir %REPO%\web\assets\dist\win64
:mkdir %REPO%\web\assets\dist\linux64
:mkdir %REPO%\release

pushd %REPO%\cmd\cafe-runner-server
set GOOS=linux
packr build -ldflags "-s -w"
upx --brute cafe-runner-server
popd

xcopy %REPO%\cmd\cafe-runner-server\cafe-runner-server %REPO%\release /H /Y /C /R /S /I