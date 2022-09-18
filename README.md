
# Automatic, Atomic, NO costs channel rebalance

Service providers and Routing node operators are frequently facing the need to rebalance their channels in order to allow directional payments and routing.

Rebalancing is done today manually, by the node owner while using tools and command line apps. The process is slow, time consuming and expensive. Rebalancing costs are significant and may eat most of the routing revenue.

Some solutions also achieve rebalancing by constantly modifying the channel routing fee based on local/remote balances. While this method may rebalance the channel, it is very slow and is suffering from a high opportunity cost as the channel can't be optimally used for routing until balanced.

The `lnagent` project allows automatic, zero costs, rule based rebalancing.

## How does it work

There are two main components involved with the solution:
1. `lnagent` - this open source component is running near the lightning node and has access to it. It collects information from the node, identifies which channels should be rebalanced and share it with the `coordinator`. Once instructed by the `coordinator` the `agent` initiates or participates in an atomic rebalance operation(see below).
2. `coordinator` - A centralised, proprietary service which collects information from multiple `agents`, identifies a `cricular route of rebalance` and orchestrates the rebalance operation

The system is combined of a single `coordinator` and multiple `agents`.

### Agent Data collection

The `agent` periodically shares with the `coordinator` full channel status and rebalancing configuration allowing the `coordinator` to decide, based on the `agent` configuration, when rebalance should be done.

### Circular rebalance route
Every few minutes, the `coordinator` analyses all the rebalance demand from the different agents. Using this information, public data (network structure) and other information provided from the `agent`, it aims to create a `circular route of rebalance`.

#### circular route of rebalance
`circular route of rebalance` is a route from a head node via 1 to N other nodes ending back at the head node. The head node is selected by the `coordinator`. The nodes are connected by channels that need rebalancing.

Example: Using the information from the `agens` the `coordinator` may create this circular route of rebalance -
`xyz->def->abc->xyz`
sending sats along this route will cause all the involved channels to improve their balance situation.
In an advanced mode the `circular route of rebalance` may include also other nodes and channels that do not need rebalancing and for which the process may need to pay a routing fee.

Once a circular route of rebalance is identified the `coordinator` instructs the involved agents to execute atomic rebalancing.

### atomic rebalancing
The `coordinator` instructs the `agents` to perform a set of separated payments which are all sharing the same payment hash.

Considering the route
`xyz->ch1->def->ch2->abc->ch3->xyz`  (circular route of 3 nodes and 3 channels)
the coordinator appoints `xyz` as the rebalance head. The head creates a random secret and hashes it to get the payment hash.
The `coordinator` instruct:
- `def` - wait for a payment from `xyz` on `ch1`. Once accepted, hold it and pay `abc` on `ch2` using the same hash
- `abc` - wait for a payment from `def` on `ch2`. Once accepted, hold it and pay `xyz` on `ch3` using the same hash
- `xyz` - create a secret and payment hash, using that hash pay to `def` on `ch1`. Wait for payment from `abc` on `ch3`. Once the payment arrives, settle it by providing the secret.

Executing 3 instructions above rebalance ch1, ch2 and ch3.

This atomic rebalance has the following properties:
1. The secret is known only to the head. The coordinator and the nodes along the route do not have access to the secret until all payments are locked.
2. all payments are using the same secret hash
3. all or nothing - the method ensures that all payments will be executed or all fail. The situation in which a participating node will pay without getting paid can't happen
4. no routing fees - These separate payments are all direct payments (`abc` pays to his peer `xyz`). Direct payments do not have any routing cost. Hence the atomic rebalance operation is executed without paying any routing fees.


## Instalation

clone this repository

you can run lnagent with
>go run cmd/lnagent/main.go

or you can generate the lanaget program

> go install  ./...

## Usage
```
NAME:
   lnagent run - start the agent

USAGE:
   lnagent run [command options] [arguments...]

OPTIONS:
   --ln-host value, --lnh value                          host name/ip of the lightning node (default: "localhost")
   --ln-port value, --lnp value                          port of the lightning node (default: 10009)
   --network value, -n value                             lightning environment (mainnet, testnet, simnet) (default: "mainnet")
   --implementation value, --impl value                  specify the lightning node implementation (lnd, c-lightning) (default: "lnd")
   --lnd-dir value                                       path to lnd directory
   --lnc-host value, --lnch value                        lightning coordinator host
   --lnc-port value, --lncp value                        lightning coordinator port (default: 2222)
   --channel-low-trigger value, --low value              when local amount (as percentage of capacity) fails below the channel-low-trigger, the channel becomes rebalance candidate. Expressed as percentage [0-100] (default: 25)
   --channel-high-trigger value, --high value            when local amount (as percentage of capacity) exceeds the channel-high-trigger, the channel becomes rebalance candidate. Expressed as percentage [0-100] (default: 75)
   --channel-target value, --target value                the target of rebalance expressed as percentage [0-100] of the channel's capacity (default: 50)
   --rebalance-budget-ppm-percent value, --budget value  rebalance budget expressed as percentage of the channel's ppm (fee_per_mil). Budget may exceeds 100% (default: 0)
   --tor.active                                          tcp connections to use Tor (default: false)
   --tor.socks value                                     The host:port that Tor's exposed SOCKS5 proxy is listening on (default: localhost:9050) (default: "localhost:9050")
   --help, -h                                            show help (default: false)
   
```

### service payment
There are currently no payment for the service.

if you are happy with the result you are welcome to tip the creator using keysend to 02ade55c81601182e0df6b59e4042869687d28661850bc18079c7d7d5fa78be39e

