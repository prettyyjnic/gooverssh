# gooverssh
gooverssh - forwards over ssh.


## Installation
```bash
go get github.com/scottkiss/gooverssh

go build gooverssh.go
```

## Usage

### edit config file config.toml

```
  [ssh]
  local_bind_address = ":3306"
  ssh_server = "222.222.222.222"
  ssh_user = "root"
  ssh_password = "passwd"
  private_key_path = ""

  [remote]
  host = "111.111.111.111"
  port = "3306"

```

### run gooverssh
```
./gooverssh
```

then you connect to localhost:3306,the gooverssh will open an ssh tunne and forwards local port 3306 to 111.111.111.111:3306.

## License
View the [LICENSE](https://github.com/scottkiss/gooverssh/blob/master/LICENSE) file


