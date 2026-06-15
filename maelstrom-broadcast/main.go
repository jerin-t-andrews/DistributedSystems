package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	var topology map[string]any
	messages := make(map[int]struct{})

	// Topology handler
	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		json.Unmarshal(msg.Body, &topology)
		
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
		if _, exists := messages[int(body["message"].(float64))]; !exists {
			messages[int(body["message"].(float64))] = struct{}{}
		}

		// Propagate broadcast to neighbors


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
		list := make([]int, 0, len(messages))
		for key := range messages {
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