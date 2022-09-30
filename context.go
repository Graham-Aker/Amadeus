package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Ako struct {
	aKotype    string     `json:"type"`
	kyokuFirst int        `json:"kyoku_first,omitempty"`
	akaFlag    bool       `json:"aka_flag,omitempty"`
	names      []string   `json:"names,omitempty"`
	bakaze     string     `json:"bakaze,omitempty"`
	doraMarker string     `json:"dora_marker,omitempty"`
	kyoku      int        `json:"kyoku,omitempty"`
	honba      int        `json:"honba,omitempty"`
	oya        string     `json:"oya,omitempty"`
	scores     []int      `json:"scores,omitempty"`
	tehais     [][]string `json:"tehais,omitempty"`
	kyotaku    string     `json:"kyotaku,omitempty"`
	actor      string     `json:"actor,omitempty"`
	pai        string     `json:"pai,omitempty"`
	tsumogiri  string     `json:"tsumogiri,omitempty"`
	target     string     `json:"target,omitempty"`
	consumed   []string   `json:"consumed,omitempty"`
}

type Context struct {
	SshBuffer   *SshBuffer
	SshTerminal *SshTerminal
	Client      *ssh.Client
	Session     *ssh.Session
	Start       bool
	User        string
}

type SshBuffer struct {
	outBuf   *bytes.Buffer
	stdinBuf io.WriteCloser
}

type SshTerminal struct {
	in  chan string
	out chan string
}

//var terminatorMap = map[string]byte{"root":'#',"common":'$'}

func NewContext() *Context {
	var stdinBuf io.WriteCloser
	return &Context{
		SshBuffer: &SshBuffer{
			bytes.NewBuffer(make([]byte, 0)),
			stdinBuf,
		},
		SshTerminal: &SshTerminal{
			make(chan string, 1),
			make(chan string, 1),
		},
	}
}

func (c *Context) InitCommonTerminal() error {
	if c.Start {
		return fmt.Errorf("session is start terminal")
	}
	err := c.InitCommonSession()
	if err != nil {
		return err
	}
	session := c.Session
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Printf("get pty error:%v\n", err)
		return err
	}
	stdinBuf, err := session.StdinPipe()
	if err != nil {
		log.Printf("get stdin pipe error%v\n", err)
		return err
	}
	c.SshBuffer.stdinBuf = stdinBuf
	session.Stdout = c.SshBuffer.outBuf

	err = session.Shell()
	if err != nil {
		fmt.Printf("shell session error%v", err)
		return err
	}
	// ch := make(chan struct{})
	// go resetOutBuf(c.SshBuffer.outBuf, ch, '\\')
	// <-ch
	// fmt.Println("after reset buf", c.SshBuffer.outBuf.String())
	c.Start = true
	fmt.Println("start wait session")
	go session.Wait()
	return err
}

func resetOutBuf(outBuf *bytes.Buffer, ch chan struct{}, terminator byte) {
	buf := make([]byte, 8192)
	var t int
	for {
		if outBuf != nil {
			_, err := outBuf.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("read out buffer err:%v", err)
				break
			}
			t = bytes.LastIndexByte(buf, terminator)
			if t > 0 {
				ch <- struct{}{}
				break
			}
		}
	}
}

func (c *Context) InitCommonSession() error {
	user := "root"
	pass := "NetEase123"
	host := "59.111.59.43:22"
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
		fmt.Printf("dial ssh error :%v\n", err)
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("get session error %v\n", err)
		return err
	}
	fmt.Printf("init ssh session successed :%s\n", host)
	c.Session = session
	c.User = sshConfig.User
	return nil
}

func (c *Context) InitClient() error {
	user := "root"
	pass := "NetEase123"
	host := "59.111.59.43:22"
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
		fmt.Printf("dial ssh error :%v\n", err)
		return err
	}
	fmt.Printf("init ssh clent successed :%s\n", host)
	c.Client = client
	c.User = sshConfig.User
	return nil
}

func (c *Context) initAkoShell() {

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	session := c.Session

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

func (c *Context) EnableSudo() error {
	err := c.SendCmd("su")
	err = c.SendCmd("abcdefg")
	c.User = "root"
	return err
}

func (c *Context) SendCmd(cmd string) error {
	c.SshTerminal.in <- cmd
	err := c.listenMessages()
	<-c.SshTerminal.out
	return err
}

func (c *Context) SendCmdWithOut(cmd string) (string, error) {
	c.SshTerminal.in <- cmd
	err := c.listenMessages()
	out := <-c.SshTerminal.out
	// c.SshTerminal.out = make(chan string, 1)
	// resetOutBuf(c.SshBuffer.outBuf, ch, '\\')
	return strings.TrimSpace(strings.Split(out, "["+c.User)[0]), err
}

func (c *Context) listenMessages() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		cmd := <-c.SshTerminal.in
		_, _ = c.SshBuffer.stdinBuf.Write([]byte(fmt.Sprintf("%v\n", cmd)))
		fmt.Printf("send cmd %v\n", cmd)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		buf := make([]byte, 8192)
		var t int
		terminator := '~'
		if c.User == "root" {
			terminator = '#'
		}
		for {
			time.Sleep(time.Millisecond * 200)
			n, err := c.SshBuffer.outBuf.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Printf("read out buffer err:%v", err)
				break
			}
			if n > 0 {
				t = bytes.LastIndexByte(buf, byte(terminator))
				if t > 0 {
					c.SshTerminal.out <- string(buf[:t])
					break
				}
			} else {
				c.SshTerminal.out <- string(buf)
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}
