::change dir to the BAT directory
cd /D %~dp0

@echo off
set REPO=%CD%
set GOPATH=%REPO%\..\..\..\..
set GOROOT=%REPO%\node_modules\go-win
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%REPO%\node_modules\packr-win
@echo on

call %REPO%\scripts\buildReleaseParallel.cmd
call %REPO%\scripts\runCombatIntegrationTests.cmd
