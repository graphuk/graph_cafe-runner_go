# graph_testcafe-runner_go
A server that provides [TestCafe](https://github.com/DevExpress/testcafe "TestCafe") tests via link and collection of results

### Features
- Run TestCafe tests on all supported devices (phones, desktops, tablets), by link-clicking. Results available in web-interface.<br/>
- Client app for upload tests from repo to testcafe-runner server.
- Actual client binary embedded to server.

### How to contribute
- Requirements: windows, git, nodejs, chrome.
- Go to the folder you want to host the project.
- Check the folder's path has no Cyrillic chars, or spaces
- Run `mkdir src\github.com\graph-uk && cd src\github.com\graph-uk && git clone https://github.com/graphuk/graph_cafe-runner_go.git`
- Run `cd graph_cafe-runner_go && npm install`
- Run `buildDevAndTestIntegration.cmd` and then `buildReleaseAndTestIntegration.cmd` to make sure build works
- Make feature-branch
- Make your improves, add tests, and run buildDevAndTestIntegration.cmd. If tests passed - push your feature-branch, and merge to master.
