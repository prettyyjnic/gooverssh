# gooverssh
gooverssh - forwards over ssh.


## Installation
```bash
go get github.com/prettyyjnic/gooverssh

go build gooverssh.go

mv config.example config.toml

vim config.toml // just modify the conf to what you want
```

## Usage

### edit config file config.toml

```
[ssh]
[ssh.innser_jump]
remoteAddress = "aaaaa.com:1111"
localPort = "11111"
sshServerAddress = "bbbbb.com:2222"
sshUser = "yourname"
sshPassword = "yourpassword"
sshPriveteKeyPath = "yourpath"
sshPrivateKeyPassword = "yourprivatekeypassword"

[ssh.useInnerJump]
require = "innser_jump"
remoteAddress = "cccc.com:3333"
localPort = "22222"
sshServerAddress = "127.0.0.1:11111"
sshUser = "yourname"
sshPassword = "yourpassword"
sshPriveteKeyPath = "yourpath"
sshPrivateKeyPassword = "yourprivatekeypassword"
```

### run gooverssh
```
./gooverssh
```
then you connect to localhost:11111,the gooverssh will open an ssh tunne and forwards local port 11111 to aaaaa.com:1111.
and connect to localhost:22222,the gooverssh will open an ssh tunne and forwards local port 22222 to cccc.com:3333.



