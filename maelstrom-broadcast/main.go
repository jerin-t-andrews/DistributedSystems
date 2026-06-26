package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Mutex for map
type map_access struct {
	mu sync.RWMutex
	msgs map[int]struct{}
}

func main() {
	n := maelstrom.NewNode()
	var typed struct {
		Topology map[string][]string `json:"topology"`
	}
	messages := &map_access{
		msgs: make(map[int]struct{}),
	}

	// Topology handler
	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		json.Unmarshal(msg.Body, &typed)

		body["type"] = "topology_ok"
		delete(body, "topology")
		return n.Reply(msg, body)
	})

	// Broadcast handler
	n.Handle("broadcast", func (msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update messages received
		messages.mu.Lock()
		defer messages.mu.Unlock()
		if _, exists := messages.msgs[int(body["message"].(float64))]; !exists {
			messages.msgs[int(body["message"].(float64))] = struct{}{}

			// Propagate broadcast to neighbors
			topology := typed.Topology
			for i := range topology[n.ID()] {
				new_body := map[string]any {
					"type": "broadcast",
					"message": int(body["message"].(float64)),
				} 
				n.RPC(topology[n.ID()][i], new_body, func(msg maelstrom.Message) error {
					return nil
				})
			}
		}

		body["type"] = "broadcast_ok"
		delete(body, "message")

		return n.Reply(msg, body)
	})

	// Read/Received handler
	n.Handle("read", func (msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Convert message hashset to list
		messages.mu.RLock()
		defer messages.mu.RUnlock()
		list := make([]int, 0, len(messages.msgs))
		for key := range messages.msgs {
			list = append(list, key)
		}

		body["type"] = "read_ok"
		body["messages"] = list


		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}