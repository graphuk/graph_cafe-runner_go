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
del /F /S /Q %REPO%\cmd\cafe-runner-client\cafe-runner-client.exe
del /F /S /Q %REPO%\web\assets\dist\win64 >nul 2>&1

mkdir %REPO%\web\assets\dist
mkdir %REPO%\web\assets\dist\win64

pushd %REPO%\cmd\cafe-runner-client
set GOOS=windows
go build -ldflags "-s -w"
upx --brute cafe-runner-client.exe
popd

xcopy %REPO%\cmd\cafe-runner-client\cafe-runner-client.exe %REPO%\web\assets\dist\win64 /H /Y /C /R /S /I

