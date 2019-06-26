set GOPATH=%cd%
set GOROOT=%cd%\..\..\..\node_modules\go-win
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%GOPATH%\bin

start %cd%\..\..\..\node_modules\liteide\bin\liteide.exe %cd%\src\Tests\session_creation\main.go