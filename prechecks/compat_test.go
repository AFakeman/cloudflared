package prechecks

import "github.com/cloudflare/cloudflared/connection"

func protocolPointer(protocol connection.Protocol) *connection.Protocol {
	return &protocol
}
