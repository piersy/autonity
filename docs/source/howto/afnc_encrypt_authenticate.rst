Define Secure Communication Between Nodes
=================

To enable secure communications between nodes we need to:

- define DNS records in Google Cloud Platform

- include public and private keys in the genesis file

Check the current status of the network and ensure it's functioning correctly.

1. Get the information for the first validator:

      .. Code:: 

        helm status <node name>

  This command returns the following:

    .. Code:: 

        LAST DEPLOYED: Wed Nov  6 14:18:00 2019
        NAMESPACE: how-to-se
        STATUS: DEPLOYED
        ...
        ...
        ...
        ==> v1/Job
        NAME                             COMPLETIONS  DURATION  AGE
        init-job02-genesis-configurator  0/1          14m       14m
        ==> v1/Pod(related)
        NAME                                   READY  STATUS    RESTARTS  AGE
        autonity-node-0-67b5fdf8b-kqspj        0/2    Init:0/1  0         14m
        init-job02-genesis-configurator-54qp6  1/1    Running   0         14m
        ==> v1/Role
        NAME           AGE
        genesis-write  14m
        secrets-write  14m
        ==> v1/RoleBinding
        NAME           AGE
        genesis-write  14m
        secrets-write  14m
        ==> v1/Secret
        NAME             TYPE    DATA  AGE
        autonity-node-0  Opaque  2     14m
        ==> v1/Service
        NAME                 TYPE          CLUSTER-IP    EXTERNAL-IP    PORT(S)                     AGE
        autonity-node-0      ClusterIP     10.7.246.119  <none>         8545/TCP,8546/TCP,9200/TCP  14m
        p2p-autonity-node-0  LoadBalancer  10.7.251.42   35.230.150.24  30303:31449/TCP             14m
        ==> v1/ServiceAccount
        NAME                           SECRETS  AGE
        autonity-genesis-configurator  1        14m
        autonity-keys-generator        1        14m
        ...
        ...
        ...
 
    As shown in the message above under the pod information:

        .. Code:: 

            ==> v1/Pod(related)
            NAME                                   READY  STATUS    RESTARTS  AGE
            autonity-node-0-67b5fdf8b-kqspj        0/2    Init:0/1  0         14m
            init-job02-genesis-configurator-54qp6  1/1    Running   0         14m

    There are 2 running services:

        - `autonity-node-0-67b5fdf8b-kqspj` which is the Autonity node in `Init` status

        - `init-job02-genesis-configurator-54qp6` which is the peer discovery process trying to find the other nodes listed in the `gensis.yaml`. This service is currently in `Running` status

2. Check the status of the peer discovery:

    .. Code:: bash

        kubectl -n <node namespace> logs <peer discovery service name>

    The service name in this example is `init-job02-genesis-configurator-54qp6` (see the example above).

    This command returns a similar message if the peer-to-peer discovery is **not yet completed**:

        .. Code:: 

              2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
              2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-3-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
              2019-11-06 14:18:06 WARNING  A record: None of DNS query names exist: val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
              2019-11-06 14:18:06 WARNING  TXT record: None of DNS query names exist: val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com., val-4-se.4c621a00-2099-45c8-b50c-f06f95c0bcf3.com.
              2019-11-06 14:18:06 INFO     Fully resolved 0 fqdn records from 5

    The message above shows the node is trying to resolve the other network peers. The current status is listed in the following line:

        .. Code:: bash 

            2019-11-06 14:18:06 INFO     Fully resolved 0 fqdn records from 5

    The above line mentions that at that point the node cannot resolve any of the DNSs mentioned in the `genesis.yaml` file. 
    At this point, DNS entries **MUST** be created for the validator FQDN names within the domain as described in the genesis.yaml file.

3. Create the required DNSs through the GCP for all 5 nodes by setting the following information for each DNS:

    a. Create an `A` record which includes the node IP address

    b. Create a `TXT` record which includes the node port number and the node's public key value. The public key will be used for authentication between the nodes and enables peer-to-peer discovery

    c. Run the following command **per node** to get the required values: 

         .. Code::

              `helm status <node-name>`

        In the resulting output, scroll to the second-to-last section `### Get enode` and copy the string resembling:

            .. Code:: 

                ...
                ...
                ...
                ### Get enode:
                # It can take a time to wait until Public IP will allocated
                  IP=$(kubectl -n how-to-se get svc p2p-autonity-node-0 -o jsonpath="{.status.loadBalancer.ingress[*].ip}"); \
                    PUB_KEY=$(kubectl -n how-to-se get configmap autonity-node-0 -o jsonpath="{.data.pub_key}"); \
                    PORT=$(kubectl -n how-to-se get svc p2p-autonity-node-0 -o jsonpath="{.spec.ports[0].port}"); \
                    echo "enode://"${PUB_KEY}\@${IP}\:${PORT}
                    ...
                    ...
                    ...
                  
        You see the following message:

            .. Code:: 

              enode://fcd5c05d98846325f5578f825ed05fbd96ef073b8d45e88eb3e9cc298b92326d5c2c4d3b0492862ce9f142ba04a109ee726449f97f009e24be4e898b000dad62@35.230.150.24:30303

  The above message includes all the required information for the DNS creation as follows:

    .. list-table:: 
       :widths: 5 30 30 10
       :header-rows: 1

       * - Type
         - Name
         - Value
         - TTL
       * - A
         - [node name].[cloud DNS domain]
         - [node ip address]
         - 1 min
       * - TXT
         - [node name].[cloud DNS domain]
         - "p=[port number]; k=[node public key]"
         - 1 min 



d. Repeat step 3c for all the 5 nodes 

    .. Note:: Setting up DNS records needs to be undertaken on the GCP platform under 'Network Services > Cloud DNS. Please contact your system/GCP administrator to setup the domains and provide the above information

    Your nodes can now communicate. If the process is successful, the command:

        .. Code::  

            kubectl -n <node namespace> logs <peer discovery service name>

    returns the list of peers (or `users`):

        .. Code:: 

            2019-12-04 12:01:33 INFO     Fully resolved 5 fqdn records from 5
            2019-12-04 12:01:33 INFO     All fqdn records was resolved successfully
            2019-12-04 12:01:33 INFO     Generated genesis was written successfully to ConfigMap genesis
               ...
               ...
                  "users": [
                    {
                      "enode": "enode://5dc5a89ac3fddac21f06392c44fb7fcf39e5319aafd70677b57f769fc16bd9e5dc5d2906fe3bc501319e8b9d58e4cc309acb0a03613c2bb91ef8bcc74752cd59@35.189.95.220:30303",
                      "stake": 50000,
                      "type": "validator"
                    },
                    {
                      "enode": "enode://44df1498fad8d9065fb571b0acfa719e616614d63a47da232d680f5e3714638bb094a0c40807cc849cdfd5ad44306037c47d8c150b34aec1ffd1859117cde04d@35.234.155.208:30303",
                      "stake": 50000,
                      "type": "validator"
                    },
                    {
                      "enode": "enode://28da4f0ddd440866a006c9d5c389bccdb6ddc8669564262c1a32be52e996935d8accd11048ce425fc3035c53bb1981b52d5854127aaa4be7430f2a7d4b1aa255@35.197.209.67:30303",
                      "stake": 50000,
                      "type": "validator"
                    },
                    {
                      "enode": "enode://e741dd1ca01764c0f935b81c588db140fdb4441eb218a68147b3deff900da5baf75c0ef71e4a6a6053187a928ea23e8f07ffb3dec801d123e2844ac8ed0d302a@35.246.69.161:30303",
                      "stake": 50000,
                      "type": "validator"
                    },
                    {
                      "enode": "enode://f0276cc49732c10d3a74f65e1260354690a90df134d8bcbed1f6bf5bb913db7c465bfd12bc0e292a91836ef8a6ca0feeb68e9e64f475f4d9542ba0aa93279878@35.197.239.151:30303",
                      "type": "participant"
                    }
                  ]
                },
                ...
                ...


  The message states that all nodes has been resolved and returns the updated network genesis file that can be used later to deploy any new nodes.

5. You can verify the network's active status by running the following commands. For any node:

    a. Run `helm status <node name>`. From the resultant output under the `NOTES` section:

    b. Copy and run the command to `# Forward rpcapi autonity-node-0 to localhost`

    c. Copy the command under `### HTTP(s)-RPC ###` to `# Get last block number` and open a new terminal.

    d. On a new terminal, paste and run the copied command: `curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545` to get the last block number

    e. Repeat step d and compare the returned block number. The block number should be different to verify that the network is running and blocks are being mined
