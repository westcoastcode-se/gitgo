#!/bin/bash
set -e
cd data

export MSYS_NO_PATHCONV=1

# Create the root ca
openssl genrsa -des3 -passout pass:changeme -out ca.pass.key 4096
openssl rsa -passin pass:changeme -in ca.pass.key -out ca.key
rm ca.pass.key

openssl req -new -x509 -days 3650 -key ca.key -out ca.crt \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=TestRootCA"

# Create the certificate for the api server
openssl genrsa -aes256 -passout pass:apiserver -out apiserver.pass.key 4096
openssl rsa -passin pass:apiserver -in apiserver.pass.key -out apiserver.key
rm apiserver.pass.key

openssl req -new -key apiserver.key -out apiserver.csr \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=apiserver"
openssl x509 -CAcreateserial -req -days 365 -in apiserver.csr -CA ca.crt -CAkey ca.key -out apiserver.crt
rm apiserver.csr

# Create the certificate for the git server
openssl genrsa -aes256 -passout pass:gitserver -out gitserver.pass.key 4096
openssl rsa -passin pass:gitserver -in gitserver.pass.key -out gitserver.key
rm gitserver.pass.key

openssl req -new -key gitserver.key -out gitserver.csr \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=gitserver"
openssl x509 -req -days 365 -in gitserver.csr -CA ca.crt -CAkey ca.key -out gitserver.crt
rm gitserver.csr

# Create a certificate that the git server uses when communicating with the api server
openssl genrsa -aes256 -passout pass:apiserver_client -out apiserver_client.pass.key 4096
openssl rsa -passin pass:apiserver_client -in apiserver_client.pass.key -out apiserver_client.key
rm apiserver_client.pass.key

openssl req -new -key apiserver_client.key -out apiserver_client.csr \
  -subj "/C=SE/ST=/L=Gothenburg/O=West Coast Code AB/OU=IT Department/CN=apiserverclient"
openssl x509 -req -days 365 -in apiserver_client.csr -CA ca.crt -CAkey ca.key -out apiserver_client.crt
#cat apiserver_client.key apiserver_client.crt ca.crt > apiserver_client.pem
rm apiserver_client.csr
