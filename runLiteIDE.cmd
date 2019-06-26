set GOPATH=%cd%\..\..\..\..\
set GOROOT=%cd%\node_modules\go-win
set PATH=%PATH%;%GOPATH%\bin
set PATH=%PATH%;%GOROOT%\bin
set PATH=%PATH%;%cd%\node_modules\packr-win

start %cd%\node_modules\liteide\bin\liteide.exe %cd%\cmd\cafe-runner-server\main.go

