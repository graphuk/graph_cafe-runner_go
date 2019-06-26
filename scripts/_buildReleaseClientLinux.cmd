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
del /F /S /Q %REPO%\cmd\cafe-runner-client\cafe-runner-client
del /F /S /Q %REPO%\web\assets\dist\linux64 >nul 2>&1

mkdir %REPO%\web\assets\dist
mkdir %REPO%\web\assets\dist\linux64 

pushd %REPO%\cmd\cafe-runner-client
set GOOS=linux
go build -ldflags "-s -w"
upx --brute cafe-runner-client
popd

xcopy %REPO%\cmd\cafe-runner-client\cafe-runner-client %REPO%\web\assets\dist\linux64 /H /Y /C /R /S /I

