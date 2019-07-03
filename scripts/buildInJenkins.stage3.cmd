k:
set PATH=%PATH%;%CD%/node_modules/hub-win
hub release create -t master -m testTitle -a release\cafe-runner-server.exe -a release\cafe-runner-server %1
