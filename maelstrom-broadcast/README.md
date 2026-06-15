# Broadcast

### Constraints/Notes
- The broadcast has 3 RPC request bodies:
    - **topology**:
        ```
        {
        "type": "topology",
        "topology": {
            "n1": ["n2", "n3"],
            "n2": ["n1"],
            "n3": ["n1"]
        }
        }
        ```
        - The topology describes which nodes are connected
        - Technically, all nodes can communicate with each other regardless of the network topology
    - **broadcast**:
        ```
        {
            "type": "broadcast",
            "message": 1000
        }
        ```
        - the "message" containst the vaue that the node should broadcast to all nodes in the cluster
        - Every message is unique (i.e. the whole system will only have unique messages)
    - **read**:
        ```
        {
            "type": "broadcast_ok"
        }
        ```
        - This request is for the node to return/reply with all values it has seen so far (from broadcasts and requests for broadcast)
        - response:
            ```
            {
                "type": "read_ok",
                "messages": [1, 8, 72, 25]
            }
            ```
            * Order of returned values doesn't matter

### Approach(s)
1) Store topology, send broadcast to neighbors as broadcast, neigbors propogate said broadcast, thus leading to everyone receiving the broadcast
    - This idea assumes that the network is fully connected and that every node is reachable.
    - Each node will need to store a list/hashset of the messages it has received either from a maelstrom broadcast request or from a broadcast request from another node
    - When receiving a broadcast request, we take the message and update our current list only if it isn't a message we already have received (e.g. check hashset for existence). If we have not received this message yet, we update our list and then send a broadcast to all our neighboring nodes
    - When sending a broadcast, we iterate through all neigboring nodes and send the message we are broadcasting
    - **Q: How can we guarantee that no node receives the same broadcast request twice?**
        - In our custom broadcast message should we add an extra key that keeps track of the nodes that have already received a message already? This way the number of messages sent is just O(n), the number of nodes.
        - Otherwise we will send messages to every node multiple times and it double checks its received messages to see if it needs to send the message. This is the naive and easier approach, but it is more inefficient
