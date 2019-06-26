::change dir to the BAT directory
cd /D %~dp0

@echo off
set REPO=%CD%\..
set GOPATH=%REPO%\..\..\..\..
set GOROOT=%REPO%\node_modules\go-win
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%REPO%\node_modules\packr-win
@echo on

:delete old binaries
del /F /S /Q %REPO%\cmd\cafe-runner-client\cafe-runner-client*
del /F /S /Q %REPO%\cmd\cafe-runner-server\cafe-runner-server*
del /F /S /Q %REPO%\web\assets\dist >nul 2>&1
rmdir /S /Q %REPO%\web\assets\dist >nul 2>&1

mkdir %REPO%\web\assets\dist
mkdir %REPO%\web\assets\dist\win64

pushd %REPO%\cmd\cafe-runner-client
go build
popd

xcopy %REPO%\cmd\cafe-runner-client\cafe-runner-client.exe %REPO%\web\assets\dist\win64 /H /Y /C /R /S /I

pushd %REPO%\cmd\cafe-runner-server
packr build
popd