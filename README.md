# gitgo

## server

1. Create a git user
2. Create a gitgo user (used together with docker container)
3. Configure authorized_keys so that it redirects requests to the gitgo server

https://serverfault.com/questions/749474/ssh-authorized-keys-command-option-multiple-commands

```bash
git init --bare repository-path
```

```bash
# Create private key
openssl genrsa -des3 -out server.key 2048

openssl req \
       -newkey rsa:2048 -nodes -keyout server.key \
       -x509 -days 365 -out server.crt
# Create CA

```

```bash
ssh -p 2222 -o StrictHostKeyChecking=no gitgo@127.0.0.1 "SSH_ORIGINAL_COMMAND=\"$SSH_ORIGINAL_COMMAND\" $0 $@"
```

```bash
# SSH pubkey from git user
ssh-rsa <gitgo host key>

# other keys from users
command="docker exec gitgo --config=/gitgo/conf/app.conf serv key-1",no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty <user pubkey>
```

## client