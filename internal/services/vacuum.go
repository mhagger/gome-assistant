package services

import (
	"context"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type Vacuum struct {
	conn *ws.WebsocketWriter
	ctx  context.Context
}

/* Public API */

// Tell the vacuum cleaner to do a spot clean-up.
// Takes an entityId.
func (v Vacuum) CleanSpot(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "clean_spot"

	v.conn.WriteMessage(v.ctx, req)
}

// Locate the vacuum cleaner robot.
// Takes an entityId.
func (v Vacuum) Locate(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "locate"

	v.conn.WriteMessage(v.ctx, req)
}

// Pause the cleaning task.
// Takes an entityId.
func (v Vacuum) Pause(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "pause"

	v.conn.WriteMessage(v.ctx, req)
}

// Tell the vacuum cleaner to return to its dock.
// Takes an entityId.
func (v Vacuum) ReturnToBase(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "return_to_base"

	v.conn.WriteMessage(v.ctx, req)
}

// Send a raw command to the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SendCommand(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "send_command"
	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	v.conn.WriteMessage(v.ctx, req)
}

// Set the fan speed of the vacuum cleaner. Takes an entityId and an optional
// map that is translated into service_data.
func (v Vacuum) SetFanSpeed(entityId string, serviceData ...map[string]any) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "set_fan_speed"

	if len(serviceData) != 0 {
		req.ServiceData = serviceData[0]
	}

	v.conn.WriteMessage(v.ctx, req)
}

// Start or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) Start(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "start"

	v.conn.WriteMessage(v.ctx, req)
}

// Start, pause, or resume the cleaning task.
// Takes an entityId.
func (v Vacuum) StartPause(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "start_pause"

	v.conn.WriteMessage(v.ctx, req)
}

// Stop the current cleaning task.
// Takes an entityId.
func (v Vacuum) Stop(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "stop"

	v.conn.WriteMessage(v.ctx, req)
}

// Stop the current cleaning task and return to home.
// Takes an entityId.
func (v Vacuum) TurnOff(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "turn_off"

	v.conn.WriteMessage(v.ctx, req)
}

// Start a new cleaning task.
// Takes an entityId.
func (v Vacuum) TurnOn(entityId string) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "vacuum"
	req.Service = "turn_on"

	v.conn.WriteMessage(v.ctx, req)
}
