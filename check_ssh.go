package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "os"
    "strings"

    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"

)

type SSHCommand struct {
    Path string
    Env  []string
    Stdin   io.Reader
    Stdout  io.Writer
    Stderr  io.Writer    
}

type SSHClient struct {
    Config *ssh.ClientConfig
    Host    string
    Port    int
}

func (client *SSHClient) RunCommand(cmd *SSHCommand) error {
    var (
        session *ssh.Session
        err error
    )

    if session, err = client.newSession(); err != nil {
        return err
    }

    if err = client.prepareCommand(session, cmd); err != nil {
        return err
    }
    err = session.Run(cmd.Path)
    return err
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
    for _, env := range cmd.Env {
        variable := strings.Split(env, "=")
        if len(variable) != 2 {
            continue
        }

        if err := session.Setenv(variable[0], variable[1]); err != nil {
            return err
        }
    }

    if cmd.Stdin != nil {
        stdin, err := session.StdinPipe()
        if err != nil {
            return fmt.Errorf("Unable to setup stdin for session: %v", err)
        }
        go io.Copy(stdin, cmd.Stdin)
    }

    if cmd.Stdout != nil {
        stdout, err := session.StdoutPipe()
        if err != nil {
            return fmt.Errorf("Unable to setup stdout for session: %v", err)
        }
        go io.Copy(cmd.Stdout, stdout)
    }

    if cmd.Stderr != nil {
        stderr, err := session.StderrPipe()
        if err != nil {
            return fmt.Errorf("Unable to setup stderr for session: %v", err)
        }
        go io.Copy(cmd.Stderr, stderr)
    }

    return nil
}


func (client *SSHClient) newSession() (*ssh.Session, error) {
    connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
    if err != nil {
        return nil, fmt.Errorf("Failed to dail: %s", err)
    }

    session, err := connection.NewSession()
    if err != nil {
        return nil, fmt.Errorf("Failed to create session: %s", err)
    }

    modes := ssh.TerminalModes{
        ssh.ECHO: 0,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }

    if err := session.RequestPty("xtrem", 80, 40, modes); err != nil {
        session.Close()
        return nil, fmt.Errorf("Request for pseudo terminal failed: %s", err)
    }

    return session, nil
}


func PublicKeyFile(file string) ssh.AuthMethod {
    buffer, err := ioutil.ReadFile(file)
    if err != nil {
        return nil
    }

    key, err := ssh.ParsePrivateKey(buffer)
    if err != nil {
        return nil
    }

    return ssh.PublicKeys(key)
}


func SSHAgent() ssh.AuthMethod {
    if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err != nil {
        return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
    }
    return nil
}


func main() {
    sshConfig := &ssh.ClientConfig{
        User: os.Getenv("SSH_USER"),
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Auth: []ssh.AuthMethod{
                ssh.Password(os.Getenv("SSH_PASSWORD")),
        },
    }

    client := &SSHClient{
        Config: sshConfig,
        Host: "127.0.0.1",
        Port: 22,
    }

    cmd := &SSHCommand{
        Path: "ls -al $LC_DIR"  ,
        Env: []string{"LC_DIR=/"},
        Stdin: os.Stdin,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
    }

    fmt.Printf("Running command: %s\n", cmd.Path)
    if err := client.RunCommand(cmd); err != nil {
        fmt.Fprintf(os.Stderr, "Command run err: %\n", err)
        os.Exit(1)
    }
}
