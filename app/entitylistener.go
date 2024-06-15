package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-module/carbon"
	"saml.dev/gome-assistant/internal"
	"saml.dev/gome-assistant/websocket"
)

type EntityListener struct {
	entityIDs []string
	callback  EntityListenerCallback
	fromState string
	toState   string
	throttle  time.Duration
	lastRan   carbon.Carbon

	betweenStart string
	betweenEnd   string

	delay      time.Duration
	delayTimer *time.Timer

	exceptionDates  []time.Time
	exceptionRanges []timeRange

	runOnStartup          bool
	runOnStartupCompleted bool

	enabledEntities  []internal.EnabledDisabledInfo
	disabledEntities []internal.EnabledDisabledInfo
}

type EntityListenerCallback func(EntityData)

type EntityData struct {
	TriggerEntityID string
	FromState       string
	FromAttributes  map[string]any
	ToState         string
	ToAttributes    map[string]any
	LastChanged     websocket.TimeStamp
}

type stateChangedMsg struct {
	websocket.BaseMessage
	Event struct {
		Data struct {
			EntityID string   `json:"entity_id"`
			NewState msgState `json:"new_state"`
			OldState msgState `json:"old_state"`
		} `json:"data"`
		EventType string `json:"event_type"`
		Origin    string `json:"origin"`
	} `json:"event"`
}

type msgState struct {
	EntityID    string              `json:"entity_id"`
	LastChanged websocket.TimeStamp `json:"last_changed"`
	State       string              `json:"state"`
	Attributes  map[string]any      `json:"attributes"`
}

/* Methods */

func NewEntityListener() elBuilder1 {
	return elBuilder1{EntityListener{
		lastRan: carbon.Now().StartOfCentury(),
	}}
}

type elBuilder1 struct {
	entityListener EntityListener
}

func (b elBuilder1) EntityIDs(entityIDs ...string) elBuilder2 {
	if len(entityIDs) == 0 {
		panic("must pass at least one entityID to EntityIDs()")
	} else {
		b.entityListener.entityIDs = entityIDs
	}
	return elBuilder2(b)
}

type elBuilder2 struct {
	entityListener EntityListener
}

func (b elBuilder2) Call(callback EntityListenerCallback) elBuilder3 {
	b.entityListener.callback = callback
	return elBuilder3(b)
}

type elBuilder3 struct {
	entityListener EntityListener
}

func (b elBuilder3) OnlyBetween(start string, end string) elBuilder3 {
	b.entityListener.betweenStart = start
	b.entityListener.betweenEnd = end
	return b
}

func (b elBuilder3) OnlyAfter(start string) elBuilder3 {
	b.entityListener.betweenStart = start
	return b
}

func (b elBuilder3) OnlyBefore(end string) elBuilder3 {
	b.entityListener.betweenEnd = end
	return b
}

func (b elBuilder3) FromState(s string) elBuilder3 {
	b.entityListener.fromState = s
	return b
}

func (b elBuilder3) ToState(s string) elBuilder3 {
	b.entityListener.toState = s
	return b
}

func (b elBuilder3) Duration(s DurationString) elBuilder3 {
	d := internal.ParseDuration(string(s))
	b.entityListener.delay = d
	return b
}

func (b elBuilder3) Throttle(s DurationString) elBuilder3 {
	d := internal.ParseDuration(string(s))
	b.entityListener.throttle = d
	return b
}

func (b elBuilder3) ExceptionDates(t time.Time, tl ...time.Time) elBuilder3 {
	b.entityListener.exceptionDates = append(tl, t)
	return b
}

func (b elBuilder3) ExceptionRange(start, end time.Time) elBuilder3 {
	b.entityListener.exceptionRanges = append(
		b.entityListener.exceptionRanges, timeRange{start, end},
	)
	return b
}

func (b elBuilder3) RunOnStartup() elBuilder3 {
	b.entityListener.runOnStartup = true
	return b
}

// Enable this listener only when the current state of {entityID}
// matches {state}. If there is a network error while retrieving
// state, the listener runs if {runOnNetworkError} is true.
func (b elBuilder3) EnabledWhen(entityID, state string, runOnNetworkError bool) elBuilder3 {
	if entityID == "" {
		panic(
			fmt.Sprintf(
				"entityID is empty in EnabledWhen entityID='%s' state='%s'",
				entityID, state,
			),
		)
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.entityListener.enabledEntities = append(b.entityListener.enabledEntities, i)
	return b
}

// Disable this listener when the current state of {entityID} matches
// {state}. If there is a network error while retrieving state, the
// listener runs if {runOnNetworkError} is true.
func (b elBuilder3) DisabledWhen(entityID, state string, runOnNetworkError bool) elBuilder3 {
	if entityID == "" {
		panic(
			fmt.Sprintf(
				"entityID is empty in EnabledWhen entityID='%s' state='%s'",
				entityID, state,
			),
		)
	}
	i := internal.EnabledDisabledInfo{
		Entity:     entityID,
		State:      state,
		RunOnError: runOnNetworkError,
	}
	b.entityListener.disabledEntities = append(b.entityListener.disabledEntities, i)
	return b
}

func (b elBuilder3) Build() EntityListener {
	return b.entityListener
}

/* Functions */
func (app *App) callEntityListeners(chanMsg websocket.Message) {
	msgBytes := chanMsg.Raw
	msg := stateChangedMsg{}
	json.Unmarshal(msgBytes, &msg)
	data := msg.Event.Data
	eid := data.EntityID
	listeners, ok := app.entityListeners[eid]
	if !ok {
		// no listeners registered for this id
		return
	}

	// if new state is same as old state, don't call
	// event listener. I noticed this with iOS app location,
	// every time I refresh the app it triggers a device_tracker
	// entity listener.
	if msg.Event.Data.NewState.State == msg.Event.Data.OldState.State {
		return
	}

	for _, l := range listeners {
		// Check conditions
		if c := checkWithinTimeRange(l.betweenStart, l.betweenEnd); c.fail {
			continue
		}
		if c := checkStatesMatch(l.fromState, data.OldState.State); c.fail {
			continue
		}
		if c := checkStatesMatch(l.toState, data.NewState.State); c.fail {
			if l.delayTimer != nil {
				l.delayTimer.Stop()
			}
			continue
		}
		if c := checkThrottle(l.throttle, l.lastRan); c.fail {
			continue
		}
		if c := checkExceptionDates(l.exceptionDates); c.fail {
			continue
		}
		if c := checkExceptionRanges(l.exceptionRanges); c.fail {
			continue
		}
		if c := checkEnabledEntity(app.State, l.enabledEntities); c.fail {
			continue
		}
		if c := checkDisabledEntity(app.State, l.disabledEntities); c.fail {
			continue
		}

		entityData := EntityData{
			TriggerEntityID: eid,
			FromState:       data.OldState.State,
			FromAttributes:  data.OldState.Attributes,
			ToState:         data.NewState.State,
			ToAttributes:    data.NewState.Attributes,
			LastChanged:     data.OldState.LastChanged,
		}

		if l.delay != 0 {
			l := l
			l.delayTimer = time.AfterFunc(l.delay, func() {
				go l.callback(entityData)
				l.lastRan = carbon.Now()
			})
			continue
		}

		// run now if no delay set
		go l.callback(entityData)
		l.lastRan = carbon.Now()
	}
}
