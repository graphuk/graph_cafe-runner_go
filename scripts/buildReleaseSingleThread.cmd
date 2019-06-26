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

pushd %REPO%\cmd\cafe-runner-client
set GOOS=windows
go build -ldflags "-s -w"
popd

pushd %REPO%\cmd\cafe-runner-client
set GOOS=linux
go build -ldflags "-s -w"
popd

xcopy %REPO%\cmd\cafe-runner-client\cafe-runner-client.exe %REPO%\web\assets\dist\win64 /H /Y /C /R /S /I
xcopy %REPO%\cmd\cafe-runner-client\cafe-runner-client %REPO%\web\assets\dist\linux64 /H /Y /C /R /S /I

pushd %REPO%\cmd\cafe-runner-server
set GOOS=windows
packr build -ldflags "-s -w"
upx --brute cafe-runner-server.exe
popd

pushd %REPO%\cmd\cafe-runner-server
set GOOS=linux
packr build -ldflags "-s -w"
upx --brute cafe-runner-server
popd

xcopy %REPO%\cmd\cafe-runner-server\cafe-runner-server.exe %REPO%\release /H /Y /C /R /S /I
xcopy %REPO%\cmd\cafe-runner-server\cafe-runner-server %REPO%\release /H /Y /C /R /S /I

