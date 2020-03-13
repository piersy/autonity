Install Autonity
================

We'll install Autonity from a binary file.

Prerequisite
------------

1.Install Go:

  .. Code:: bash

      brew install go

  .. note:: Alternatively, install from: https://golang.org/doc/install

  run 'go version' to check that it has installed correctly.


2. Make the following directory:

  .. Code:: bash

    mkdir -p ~/go/src/github.com/clearmatics/autonity

3. Go to the directory you just created:

  .. Code:: bash

    cd ~/go/src/github.com/clearmatics/autonity



Install and build Autonity
-------------------------

1. Download the latest Autonity build:
  .. Code:: bash

      git clone https://github.com/clearmatics/autonity

2. Build Autonity:

  .. Code:: bash

      make autonity

3. Attach a javascript console:

  .. Code:: bash

      ~/go/src/github.com/clearmatics/autonity/build/bin/autonity attach http://{IP}:{rpcport}

Replace the IP and Port addresses you received when onboarding. Here is an example:

  .. Code:: bash

      attach https://participant0.magneto.network:8545

On successful installation, you see:

  .. Code:: bash

      Welcome to the Autonity JavaScript console!
      instance: Autonity/v0.3.0-e245e9cb-20191115/linux-amd64/go1.12.13
      coinbase: 0x0db5c674b49e2b1d5699ea4addcd3d154dfbe102
      at block: 5731653 (Tue, 03 Mar 2020 14:42:24 GMT)
      modules: debug:1.0 eth:1.0 net:1.0 rpc:1.0 tendermint:1.0 txpool:1.0 web3:1.0
