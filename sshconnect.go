package sshconnect

// provide a way to run shell commands from go programs over a shell connection
// based off code found at https://zaiste.net/posts/executing_commands_via_ssh_using_go/

import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "os"
    "strings"
    "bytes"
    "bufio"
    "path/filepath"
    "io/ioutil"
    "errors"
)

// wrapper; get an ssh client
func getClient(hostname string, port string, config *ssh.ClientConfig) (*ssh.Client, error) {
    host_string := fmt.Sprintf("%s:%s", hostname, port)
    return ssh.Dial("tcp", host_string, config)
}

// parse public key
func GetPublicKey(keypath string) (ssh.PublicKey, error) {
    file, err := ioutil.ReadFile(keypath)
    if err != nil {
        panic(err)
    }
    publickey, _, _, _, err := ssh.ParseAuthorizedKey(file)
    return publickey, err
}

// get the remote host key
func GetRemoteHostKey(hostname string) (ssh.PublicKey, error) {

    known_hosts, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh/known_hosts"))
    if err != nil {
        panic(err)
    }
    defer known_hosts.Close()

    scanner := bufio.NewScanner(known_hosts)
    var hostkey ssh.PublicKey
    for scanner.Scan(){
        fields := strings.Split(scanner.Text(), " ")
        if len(fields) != 3 {
            continue
        }
        if strings.Contains(fields[0], hostname) {
            hostkey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
            if err != nil {
                return nil, errors.New(fmt.Sprintf("error parsing %q: %v", fields[2], err))
            }
            break
        }
    }
    if hostkey == nil {
        return nil, errors.New(fmt.Sprintf("no hostkey for %s", hostname))
    }
    return hostkey, nil
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

    // get stdout
    var stdoutBuf bytes.Buffer
    session.Stdout = &stdoutBuf

    err := session.Run(command)
    if err != nil {
        return fmt.Sprintf("Got error: %s\n", err), nil
    }
    return fmt.Sprintf(stdoutBuf.String()), nil
}
