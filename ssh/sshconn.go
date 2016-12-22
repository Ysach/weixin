package main

import (
	"golang.org/x/crypto/ssh"
	"time"
	"fmt"
	"log"
	//"bytes"
	"os"
	"golang.org/x/crypto/ssh/terminal"
)

func connect(user, password, host string, port int) (*ssh.Session, error) {
	//定义变量
	var (
		auth         []ssh.AuthMethod
		address string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err error
	)
	// 获取认证方式
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	// 设置config配置
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}

	// 连接
	address = fmt.Sprintf("%s:%d", host, port)
	fmt.Println("连接的地址为：", address)

	if client, err = ssh.Dial("tcp", address, clientConfig); err != nil {
		return nil, err
	}

	// 建立session连接
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func main() {
	session, err := connect("lottery", "zzc!qazxsw2", "192.168.1.81", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	/*
	//打印出来
	//session.Stdout = os.Stdout
	//session.Stderr = os.Stderr
	var b bytes.Buffer
	session.Stdout = &b
	CMD := "ls /tmp"
	fmt.Println("执行的命令为：", CMD)
	session.Run(CMD)
	fmt.Println(b.String())
	*/

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// 执行命令
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	// 设置终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		log.Fatal(err)
	}

	session.Run("ls /tmp")
}
