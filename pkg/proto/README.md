



# Install depends
```bash
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.0/protoc-3.9.0-linux-x86_64.zip -O protoc.zip
unzip protoc.zip
sudo apt install npm

sudo su -
PROTOV=3.17.3
WEBGENV=1.3.1
wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOV}/protoc-${PROTOV}-linux-x86_64.zip
unzip -o protoc.zip -d /usr/local bin/protoc
chmod +x /usr/local/bin/protoc
unzip -o protoc.zip -d /usr/local include/*
chmod 755 /usr/local/include/ -R
rm -f protoc.zip
wget -O /usr/local/bin/protoc-gen-grpc-web https://github.com/grpc/grpc-web/releases/download/${WEBGENV}/protoc-gen-grpc-web-${WEBGENV}-linux-x86_64
chmod +x /usr/local/bin/protoc-gen-grpc-web
exit

go get -u github.com/golang/protobuf/protoc-gen-go
npm install

```


# GENERATE 
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I=./pkg/proto/ \
--go_out="./pkg/proto/" \
--go-grpc_out="./pkg/proto/" \
./pkg/proto/gophkeeper.proto


```