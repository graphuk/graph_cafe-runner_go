package cmdutils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type CmdProcess struct {
	Cmd       *exec.Cmd
	Command   string
	StdOut    io.ReadCloser
	StdErr    io.ReadCloser
	StdErrBuf []byte
	StdOutBuf []byte
}

// this loop function - for separate concurrency go-routine.
// it is get text from console pipe.
// if command's buffer will overflow - command was paused untill we get this bytes
func (t *CmdProcess) refreshErrBufLoop() {
	buf := make([]byte, 512)
	for {
		len, err := t.StdErr.Read(buf)
		if err != nil {
			if err.Error() == `EOF` { // if the pipe closed (app is finished) - stop watching
				break
			} else {
				panic(err)
			}
		}
		if len > 0 {
			t.StdErrBuf = append(t.StdErrBuf, buf[:len]...)
		}
		if len == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// this function returns cut filepath on t.Command, and return short command
//D:\combat_server_current\src\github.com\graph-uk\combat-server\integration-tests\client\combat-client.exe
//combat-client.exe
func (t *CmdProcess) GetShortCommand() string {
	arr := strings.Split(t.Command, `\`) // split by '/' or '\'
	if len(arr) > 0 {
		return arr[len(arr)-1]
	} else {
		return `Cannot extract short command`
	}
}

// this loop function - for separate concurrency go-routine.
// it is get text from console pipe.
// if command's buffer will overflow - command was paused untill we get this bytes
func (t *CmdProcess) refreshOutBufLoop() {
	buf := make([]byte, 512)
	for {
		len, err := t.StdOut.Read(buf)
		if err != nil {
			if err.Error() == `EOF` { // if the pipe closed (app is finished) - stop watching
				break
			} else {
				panic(err)
			}
		}
		if len > 0 {
			t.StdOutBuf = append(t.StdOutBuf, buf[:len]...)
		}
		if len == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *CmdProcess) WaitingForStdErrContains(textPart string, timeout time.Duration) {
	startMoment := time.Now()
	log.Println(`AwaitOut - ` + t.GetShortCommand() + `: ` + textPart)
	for {
		if strings.Contains(string(t.StdErrBuf), textPart) {
			break
		}
		if startMoment.Add(timeout).Before(time.Now()) { // if timed out
			panic(`TimeoutErr - ` + t.GetShortCommand() + `: ` + textPart)
		}
		time.Sleep(time.Second)
	}
}

func (t *CmdProcess) WaitingForStdOutContains(textPart string, timeout time.Duration) {
	startMoment := time.Now()
	log.Println(`AwaitOut - ` + t.GetShortCommand() + `: ` + textPart)
	for {
		if strings.Contains(string(t.StdOutBuf), textPart) {
			break
		}
		if startMoment.Add(timeout).Before(time.Now()) { // if timed out
			panic(`TimeoutOut - ` + t.GetShortCommand() + `: ` + textPart)
		}
		time.Sleep(time.Second)
	}
}

func (t *CmdProcess) WaitingForExitWithCode(expectedExitCode int, timeout time.Duration) {
	log.Println(`AwaitExitWithExitCode ` + strconv.Itoa(expectedExitCode) + ` ` + t.GetShortCommand())

	ch := make(chan bool, 1)
	defer close(ch)

	go func() {
		t.Cmd.Wait()
		ch <- true
	}()

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-ch:
	case <-timer.C:
		panic(`TimeoutOut - Wait for exit with code ` + strconv.Itoa(expectedExitCode) + ` ` + t.GetShortCommand())
	}

	ws := t.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	exitCode := ws.ExitStatus()
	if exitCode != expectedExitCode {
		panic(strconv.Itoa(expectedExitCode) + ` exit code expected, but the process is finished, with '` + strconv.Itoa(exitCode) + `' exit code. ` + t.GetShortCommand())
	}
}

//func (t *CmdProcess) WaitingForExitCode(timeout time.Duration) (int, error) {
//	done := make(chan error, 1)
//	go func() {
//		done <- t.Cmd.Wait()
//	}()
//	select {
//	case <-time.After(timeout):
//		return 1, errors.New(`Timeoit was reached.`)
//	case err := <-done:
//		if err != nil {
//			//log.Fatalf("process finished with error = %v", err)
//		}

//	}
//	ws := t.Cmd.ProcessState.Sys().(syscall.WaitStatus)
//	return ws.ExitStatus(), nil
//}

func (t *CmdProcess) KillWithChilds() {
	exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(t.Cmd.Process.Pid)).Run()
}

func MustStartCmd(command string, arguments []string, dir string, env []string) *CmdProcess {
	var res CmdProcess
	res.Command = command

	res.Cmd = exec.Command(command, arguments...)
	if env != nil {
		res.Cmd.Env = env
	}
	if dir != `` {
		res.Cmd.Dir = dir
	}

	var err error
	res.StdErr, err = res.Cmd.StderrPipe()
	check(err)
	res.StdOut, err = res.Cmd.StdoutPipe()
	check(err)

	log.Println(`------------------------------------------------------------------`)
	log.Println(dir)
	log.Println(command)
	log.Println(arguments)

	check(res.Cmd.Start())

	go res.refreshOutBufLoop() // stdout/stderr pipe-readers routines
	go res.refreshErrBufLoop()

	return &res
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) {

	// Get the data
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	check(err)
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	check(err)
}

func MustPost(url, body string) {
	//url := "http://localhost:3133/api/v1/sessions"

	payload := strings.NewReader(body)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	ioutil.ReadAll(res.Body)
}
