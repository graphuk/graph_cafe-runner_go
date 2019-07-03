::change dir to the BAT directory
cd /D %~dp0
cd ..

subst k: /d
subst k: .
k:

IF EXIST src rd /s /q src
mkdir src\github.com\graph-uk
mklink /D src\github.com\graph-uk\graph_cafe-runner_go %CD%
cd src\github.com\graph-uk\graph_cafe-runner_go


dir
npm install