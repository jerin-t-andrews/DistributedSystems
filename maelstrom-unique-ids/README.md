# Generate Unique Ids

### Constraints
- 3 Nodes running
- All nodes sending "generate" requests
- Single node receiving all requests
- **You forgot this! IDs can be any type!**

### Approach
1) Define i, respond with i, and increment i
    - This approach fails because all three nodes will have overlapping intervals where their i-values are equal thus resulting in duplicate ids

2) Increment relative to node #
    - For this approach we would take the node number and increment by the total number of nodes
    - EX: For a 3 node system, node i would increment by:
    unique_id = i + 3 (everytime)
    - The problem with this solution is that it is not scalable. Once a new node is introduced, the formula has to change, otherwise there won't be unique ids

Q: Is there a way we can generate an id without needing to know any information about the system globally and/or by using individual node information that is unique implicitly?

3) Concatenate unique info from requesting nodes
    - We know that when a node sends a request we receive this information:
        EX: {c8 n0 {"type":"generate","msg_id":11812}}
        - dest: c8
        - src: n0
        - then the rest is the request body
    - We know that for an individual node, the msg_ids are unique
        - Can we use the msg_ids as our unique id? No, b/c we run into the same issue as using i and incrementing. We will have overlapping ids between the various requesting nodes
    - We can leverage a combination of the msg_id and the requesting node id/name to create a unique id
    - EX: unique_id : "n011812"
    - This solution passes!