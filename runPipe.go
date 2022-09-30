package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

var escapePrompt = []byte{'$', ' '}

// type Context struct {
// 	SshBuffer   *SshBuffer
// 	SshTerminal *SshTerminal
// 	Client      *ssh.Client
// 	Session     *ssh.Session
// 	Start       bool
// 	User        string
// }
// type SshBuffer struct {
// 	outBuf   *bytes.Buffer
// 	stdinBuf io.WriteCloser
// }

// type SshTerminal struct {
// 	in  chan string
// 	out chan string
// }

// func NewContext() *Context {
// 	var stdinBuf io.WriteCloser
// 	return &Context{
// 		SshBuffer: &SshBuffer{
// 			bytes.NewBuffer(make([]byte, 0)), stdinBuf},
// 		SshTerminal: &SshTerminal{
// 			make(chan string, 1),
// 			make(chan string, 1)}}
// }

func InitSeesion(c *Context) *ssh.Session {
	user := "root"
	pass := "NetEase123"
	host := "59.111.59.43:22"
	fmt.Println("runPipe...")
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
		// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// 	return nil
		// },
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	sshConfig.SetDefaults()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		// errors.Wrap(err, err.Error())
		fmt.Println(err.Error())
	}

	var session *ssh.Session

	session, err = client.NewSession()
	if err != nil {
		// errors.Wrap(err, err.Error())
		fmt.Println(err.Error())
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		// log.Fatal(err)
		panic(err)
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}

	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// e, err := session.StderrPipe()
	// if err != nil {
	// 	panic(err)
	// }

	cmdstring := "docker exec -it 2f2 ./system.exe pipe_detailed /akochan-reviewer/tactics.json 3"
	// cmdstring := "telnet 127.0.0.1 6666"
	// cmdstring := "docker exec -it 2f2 telnet 127.0.0.1 6666"

	if err := session.Start(cmdstring); err != nil {
		log.Fatal(err)
	}

	// write(w, "{\"type\":\"dahai\",\"actor\":0,\"pai\":\"F\",\"tsumogiri\":false}")
	write(w, "{\"type\":\"start_game\",\"kyoku_first\":0,\"aka_flag\":true,\"names\":[\"オオバつよし\",\"USO八百\",\"icepeach\",\"被秒杀的杂兵\"]}")

	out, err := readUntil(r, escapePrompt)
	if err != nil {
		fmt.Printf("error: %s\n", 84)
	}

	fmt.Printf("akochan: %s\n", *out)

	write(w, "exit")

	session.Wait()

	return session

}

func runShellPipe(session *ssh.Session) {

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		// log.Fatal(err)
		panic(err)
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}

	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// e, err := session.StderrPipe()
	// if err != nil {
	// 	panic(err)
	// }

	cmdstring := "docker exec -it 2f2 ./system.exe pipe_detailed /akochan-reviewer/tactics.json 3"
	// cmdstring := "telnet 127.0.0.1 6666"
	// cmdstring := "docker exec -it 2f2 telnet 127.0.0.1 6666"

	if err := session.Start(cmdstring); err != nil {
		log.Fatal(err)
	}

	// write(w, "{\"type\":\"dahai\",\"actor\":0,\"pai\":\"F\",\"tsumogiri\":false}")
	write(w, "{\"type\":\"start_game\",\"kyoku_first\":0,\"aka_flag\":true,\"names\":[\"オオバつよし\",\"USO八百\",\"icepeach\",\"被秒杀的杂兵\"]}")

	out, err := readUntil(r, escapePrompt)
	if err != nil {
		fmt.Printf("error: %s\n", 84)
	}

	fmt.Printf("akochan: %s\n", *out)

	write(w, "exit")

	session.Wait()

}

func write(w io.WriteCloser, command string) error {
	_, err := w.Write([]byte(command + "\n"))
	return err
}

func readUntil(r io.Reader, matchingByte []byte) (*string, error) {
	var buf [64 * 1024]byte
	var t int
	for {
		n, err := r.Read(buf[t:])
		if err != nil {
			return nil, err
		}
		t += n
		// if isMatch(buf[:t], t, matchingByte) {
		stringResult := string(buf[:t])
		return &stringResult, nil
		// }
	}
}

func isMatch(bytes []byte, t int, matchingBytes []byte) bool {
	if t >= len(matchingBytes) {
		for i := 0; i < len(matchingBytes); i++ {
			if bytes[t-len(matchingBytes)+i] != matchingBytes[i] {
				return false
			}
		}
		return true
	}
	return false
}

// func MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
// 	in := make(chan string, 5)
// 	out := make(chan string, 5)
// 	var wg sync.WaitGroup
// 	wg.Add(1) //for the shell itself
// 	go func() {
// 		for cmd := range in {
// 			wg.Add(1)
// 			w.Write([]byte(cmd + "\n"))
// 			wg.Wait()
// 		}
// 	}()

// 	go func() {
// 		var (
// 			buf [1024 * 1024]byte
// 			t   int
// 		)
// 		for {
// 			n, err := r.Read(buf[t:])
// 			if err != nil {
// 				fmt.Println(err.Error())
// 				close(in)
// 				close(out)
// 				return
// 			}
// 			t += n
// 			// result := string(buf[:t])
// 			// if strings.Contains(string(buf[t-n:t]), "More") {
// 			// 	w.Write([]byte("\n"))
// 			// }
// 			// if strings.Contains(result, "username:") ||
// 			// 	strings.Contains(result, "password:") ||
// 			// 	strings.Contains(result, ">") {
// 			out <- string(buf[:t])
// 			t = 0
// 			wg.Done()
// 			// }
// 		}
// 	}()
// 	return in, out
// }
