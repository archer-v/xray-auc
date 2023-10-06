# xray-auc
Simple cmd utility to control xray users from a console

Supported operations: add user, remove user, get user traffic statistics (in development)

It was written in a hurry to quickly solve the problem of managing xray users, since xray does not provide a simple api for this.

Xray has grpc api, but it's difficult to use it from linux console by automation scripts.

Some peace of code related to grpc api flow was borrowed from https://github.com/FranzKafkaYu/x-ui

```
  Usage:
    xray-auc [addUser|rmUser] [user_options..]

  Subcommands: 
    addUser   add user to an inbound proxy configuration
       user options: 
         -p --proto      proxy protocol, one of: vmess, vless, trojan, shadowsocks
         -e --email      user email (is used as a human id)
         -s --password   user secret (password for shadowsocks or id for vless/vmess/trojan)
            --flow       flow (vless proto only)
         -c --cipher     cipher (shadowsocks proto only, optional)
         -t --tag        proxy tag

    rmUser    remove user from an inbound proxy configuration
       user options: 
         -e --email      user email (is used as a human id)
         -t --tag        proxy tag

  General flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -a --addr      xray server host and port separated with a colon
    -f --file      filepath to json file with array of user records in format: [ { user_options...}, .... {}]
```


**Build**

```make distro```
