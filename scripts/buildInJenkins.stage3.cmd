k:
set PATH=%PATH%;%CD%/node_modules/hub-win
set
echo %1
echo %2
set GITHUB_TOKEN=%2
set

hub release create -t master -m cafe-runner -a release\cafe-runner-server.exe -a release\cafe-runner-server %1
