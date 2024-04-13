package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Switch struct {
	conn *websocket.Conn
}

/* Public API */

func (s Switch) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_on"

	s.conn.WriteMessage(req)
}

func (s Switch) Toggle(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "toggle"

	s.conn.WriteMessage(req)
}

func (s Switch) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "switch"
	req.Service = "turn_off"
	s.conn.WriteMessage(req)
}
