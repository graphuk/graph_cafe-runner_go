::change dir to the BAT directory
cd /D %~dp0

@echo off
set REPO=%CD%\..
set GOPATH=%REPO%\..\..\..\..\
set GOROOT=%REPO%\node_modules\go-win

set COMBATPATH=%REPO%\node_modules\combat-win
set CURLPATH=%REPO%\node_modules\curl-win
set CHROMEPATH=%REPO%\node_modules\chrome-portable-win\App\Chrome-bin

set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%COMBATPATH%
set PATH=%PATH%;%CURLPATH%
set PATH=%PATH%;%CHROMEPATH%

@echo on

:delete old binaries
del /F /S /Q %REPO%\tests\integration\combat\src\Tests_shared\cafe-runner-server*
xcopy %REPO%\cmd\cafe-runner-server\cafe-runner-server.exe %REPO%\tests\integration\combat\src\Tests_shared /H /Y /C /R /S /I

cd %REPO%\tests\integration\combat\src\tests
combat run -HostName=localhost:3133

IF "%USERPROFILE%" == "C:\Users\jenkins" GOTO EXIT
pause

EXIT:
