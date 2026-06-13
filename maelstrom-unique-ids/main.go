package main

import (
	"encoding/json"
	"log"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		body["id"] = n.ID() + strconv.FormatFloat(body["msg_id"].(float64), 'f', -1, 64)
		// ^^Could use fmt.Sprintf() instead to do the conversion, but it is considered slower than strconv

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}