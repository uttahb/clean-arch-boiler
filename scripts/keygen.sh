#!/bin/bash

if ! command -v openssl &> /dev/null
then
    echo "OpenSSL is not installed. Installing..."
    sudo apt-get update && sudo apt-get install -y openssl
else
    echo "OpenSSL is already installed."
fi
[ -d "./keys" ] && rm -rf "./keys"

mkdir ./keys

openssl genrsa -out ./keys/auth-private.pem 2048

openssl rsa -in ./keys/auth-private.pem -outform PEM -pubout -out ./keys/auth-public.pem

openssl genrsa -out ./keys/auth-refresh-private.pem 2048

openssl rsa -in ./keys/auth-refresh-private.pem -outform PEM -pubout -out ./keys/auth-refresh-public.pem