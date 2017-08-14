#!/usr/bin/env bash

cd "$( dirname "${BASH_SOURCE[0]}" )"

# copy template to config.yaml
if [ -f config_template.yaml ]; then
    cp config_template.yaml ../config.yaml
fi

# create certcache directory
if [ ! -d ../certcache ]; then
    mkdir ../certcache
fi

cd ../certcache
find . -mindepth 1 -delete
openssl req -x509 -nodes -sha256 -days 730 -newkey rsa:2048 -keyout ca.key -out ca.crt -config ../setup/root.ini
openssl req -new -newkey rsa:2048 -nodes -sha256 -keyout server.key -out server.csr -config ../setup/req.ini
openssl x509 -req -sha256 -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.crt -extensions v3_req -extfile ../setup/req.ini
rm ca.key ca.srl server.csr
