## Bullions Blockchain

> An [Bullions Blockchain]downstream effort to make the Bullion Protocol accessible and extensible for a diverse ecosystem.

Priority is given to reducing opinions around chain configuration, IP-based feature implementations, and API predictability.
Upstream development from is merged to this repository regularly,
 usually at every upstream tagged release. Every effort is made to maintain seamless compatibility with upstream source, including compatible RPC, JS, and CLI
 APIs, data storage locations and schemas, and, of course, interoperable node protocols. Applicable bug reports, bug fixes, features, and proposals should be
 made upstream whenever possible.

[![OpenRPC](https://img.shields.io/static/v1.svg?label=OpenRPC&message=1.0.10&color=blue)](#openrpc-discovery)
[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://godoc.org/github.com/etclabscore/core-geth)
[![Go Report Card](https://goreportcard.com/badge/github.com/etclabscore/core-geth)](https://goreportcard.com/report/github.com/etclabscore/core-geth)
[![Travis](https://travis-ci.org/etclabscore/core-geth.svg?branch=master)](https://travis-ci.org/etclabscore/core-geth)
[![Gitter](https://badges.gitter.im/core-geth/community.svg)](https://gitter.im/core-geth/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)


## Install

### Pre-built executable

If you just want to download and run `geth` or any of the other tools here, this is the quickest and simplest way.

Binary archives are published at releases. Find the latest one for your OS, download it, (check the SHA sum), unarchive it, and run!

### With Docker

All runnable examples below are for images limited to `geth`. For images including the full suite of
tools available from this source, use the Docker Hub tag prefix `alltools.`, like `etclabscore/core-geth:alltools.latest`, or the associated Docker file directly `./Dockerfile.alltools`.

#### `docker run`

One of the quickest ways to get Ethereum up and running on your machine is by using
Docker:

```shell
$ docker run -d \
    --name core-geth \
    -v $LOCAL_DATADIR:/root \
    -p 30303:30303 \
    -p 8545:8545 \
    etclabscore/core-geth \
    --classic \
    --rpc --rpcport 8545
```

This will start `geth` in fast-sync mode with a DB memory allowance of 1GB just as the
above command does.  It will also create a persistent volume in your `$LOCAL_DATADIR` for
saving your blockchain, as well as map the default devp2p and JSON-RPC API ports.

Do not forget `--http.addr 0.0.0.0`, if you want to access RPC from other containers
and/or hosts. By default, `geth` binds to the local interface and RPC endpoints is not
accessible from the outside.


#### `docker pull`

Docker images are automatically [published on Docker Hub](https://hub.docker.com/r/etclabscore/core-geth/tags).

##### Image: `latest`

Image `latest` is built automatically from the `master` branch whenever it's updated.

```shell
$ docker pull bullion/core-geth:latest
```

##### Image: `<tag>`

Repository tags like `v1.2.3` correspond to Docker tags like __`version-1.2.3`__.

An example:
```shell
$ docker pull bullion/core-geth:version-1.11.1
```

#### `docker build`

You can build a local docker image directly from the source:

```shell
$ git clone https://github.com/bullion/core-geth.git
$ cd core-geth
$ docker build -t=core-geth .
```

Or with all tools:

```shell
$ docker build -t core-geth-alltools -f Dockerfile.alltools .
```

### Build from source

#### Dependencies

- Make sure your system has __Go__ installed. Version 1.13+ is recommended. https://golang.org/doc/install
- Make sure your system has a C compiler installed. For example, with Linux Ubuntu:

```shell
$ sudo apt-get install -y build-essential
```

#### Source

Once the dependencies have been installed, it's time to clone and build the source:

```shell
$ git clone https://github.com/bullion/core-geth.git
$ cd core-geth
$ make all
$ ./build/bin/geth --help
```


