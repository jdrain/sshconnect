package sshconnect

// provide a way to run shell commands from go programs over a shell connection
// based off code found at https://zaiste.net/posts/executing_commands_via_ssh_using_go/

import (
    "fmt"
    "golang.org/x/crypto/ssh"
)

// wrapper get an ssh client
func getClient(hostname string, port string, config *ssh.ClientConfig) (*ssh.Client, error) {
    host_string := fmt.Sprintf("%s:%s", hostname, port)
    return ssh.Dial("tcp", host_string, config)
}

// execute a plurality of shell commands
func ExecuteCmds(commands []string, hostname string, port string, config *ssh.ClientConfig) (map[string]string, error) {

    conn, err := getClient(hostname, port, config)
    if err != nil {
        panic(err)
    }

    results := make(map[string]string)
    for _, command := range commands {
        res, err := ExecuteCmd(command, hostname, port, conn)
        if err != nil {
            panic(err)
        }
        results[command] = res
    }
    return results, nil
}

// execute a single shell command. Designed only to be called by `ExecuteCmds`
func ExecuteCmd(command string, hostname string, port string, client *ssh.Client) (string, error) {

    session, _ := client.NewSession()
    defer session.Close()

    res, err := session.Output(command)
    if err != nil {
        panic(err)
    }
    return fmt.Sprintf(string(res)), nil
}
