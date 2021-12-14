# GitGo

GitGo is split into three parts:

1. The API server
2. The GIT server
3. The CLI client

We need a couple of certificates before setting up the application. This README shows how to make self-signed
certificates but if you have your own way of creating certificates then those can be used as well.

Start by creating a key that's used to sign the Root CA

```bash
openssl genrsa -des3 -passout pass:changeme -out ca.pass.key 4096
openssl rsa -passin pass:changeme -in ca.pass.key -out ca.key
rm ca.pass.key
```

Then create the actual CA and sign it using the key. This is the CA we will be using when creating client certificates

```bash
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=TestRootCA" 
```

Now, let's create the server certificate that we can use for the HTTPS connection

```bash
openssl genrsa -aes256 -passout pass:apiserver -out apiserver.pass.key 4096
openssl rsa -passin pass:apiserver -in apiserver.pass.key -out apiserver.key
rm apiserver.pass.key
```

And now we sign the api server private key using a certificate sign request file

```bash
openssl req -new -key apiserver.key -out apiserver.csr \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=apiserver" 
openssl x509 -CAcreateserial -req -days 365 -in apiserver.csr -CA ca.crt -CAkey ca.key -out apiserver.crt
rm apiserver.csr
```

Ok, now do the same for the git server:

```bash
openssl genrsa -aes256 -passout pass:gitserver -out gitserver.pass.key 4096
openssl rsa -passin pass:gitserver -in gitserver.pass.key -out gitserver.key
rm gitserver.pass.key

openssl req -new -key gitserver.key -out gitserver.csr \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=gitserver" 
openssl x509 -req -days 365 -in gitserver.csr -CA ca.crt -CAkey ca.key -out gitserver.crt
rm gitserver.csr
```

And lastly we create a client-side certificate used when the GIT server communicates with the API server

```bash
openssl genrsa -aes256 -passout pass:apiserver_client -out apiserver_client.pass.key 4096
openssl rsa -passin pass:apiserver_client -in apiserver_client.pass.key -out apiserver_client.key
rm apiserver_client.pass.key

openssl req -new -key apiserver_client.key -out apiserver_client.csr \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=apiserver_client" 
openssl x509 -req -days 365 -in apiserver_client.csr -CA apiserver.crt -CAkey apiserver.key -out apiserver_client.crt
#cat apiserver_client.key apiserver_client.crt ca.crt > apiserver_client.pem
rm apiserver_client.csr
```

## API Server

The server is split into two parts. The first one is used for manipulating the server using REST requests. The second
part is used to send git data over ssh

### Installation guide

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

## GIT Server

This server is responsible for processing git requests over ssh. It validates each request by communicating with an API
server. Communication with the api server is done over https using a client-side certificate.

## client

## TODO

* SuperUser user should authenticate using client-side certificate
* SSH server should validate a fingerprint using the rest API (f.ex. HEAD /api/v1/repositories/{path}/keys)