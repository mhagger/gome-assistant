package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

type HomeAssistant struct {
	conn *websocket.Conn
}

func NewHomeAssistant(conn *websocket.Conn) *HomeAssistant {
	return &HomeAssistant{
		conn: conn,
	}
}

// TurnOn a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) TurnOn(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_on",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	ha.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

// Toggle a Home Assistant entity. Takes an entityID and an optional
// map that is translated into service_data.
func (ha *HomeAssistant) Toggle(entityID string, serviceData any) {
	req := CallServiceRequest{
		Domain:  "homeassistant",
		Service: "toggle",
		Target: Target{
			EntityID: entityID,
		},
		ServiceData: serviceData,
	}

	ha.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}

func (ha *HomeAssistant) TurnOff(entityID string) {
	req := CallServiceRequest{
		Domain:  "homeassistant",
		Service: "turn_off",
		Target: Target{
			EntityID: entityID,
		},
	}

	ha.conn.Send(func(lc websocket.LockedConn) error {
		req.ID = lc.NextID()
		return lc.SendMessage(req)
	})
}
