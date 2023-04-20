// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package events

import (
	"log"
	"reflect"

	"github.com/google/uuid"
	"go-micro.dev/v4/events"
)

var (
	// MainQueueName is the name of the main queue
	// All events will go through here as they are forwarded to the consumer via the
	// group name
	// TODO: "fan-out" so not all events go through the same queue? requires investigation
	MainQueueName = "main-queue"

	// MetadatakeyEventType is the key used for the eventtype in the metadata map of the event
	MetadatakeyEventType = "eventtype"

	// MetadatakeyEventID is the key used for the eventID in the metadata map of the event
	MetadatakeyEventID = "eventid"
)

type (
	// Unmarshaller is the interface events need to fulfill
	Unmarshaller interface {
		Unmarshal([]byte) (interface{}, error)
	}

	// Publisher is the interface publishers need to fulfill
	Publisher interface {
		Publish(string, interface{}, ...events.PublishOption) error
	}

	// Consumer is the interface consumer need to fulfill
	Consumer interface {
		Consume(string, ...events.ConsumeOption) (<-chan events.Event, error)
	}

	// Stream is the interface common to Publisher and Consumer
	Stream interface {
		Publish(string, interface{}, ...events.PublishOption) error
		Consume(string, ...events.ConsumeOption) (<-chan events.Event, error)
	}

	// Event is the envelope for events
	Event struct {
		Type  string
		ID    string
		Event interface{}
	}
)

// Consume returns a channel that will get all events that match the given evs
// group defines the service type: One group will get exactly one copy of a event that is emitted
// NOTE: uses reflect on initialization
func Consume(s Consumer, group string, evs ...Unmarshaller) (<-chan Event, error) {
	c, err := s.Consume(MainQueueName, events.WithGroup(group))
	if err != nil {
		return nil, err
	}

	registeredEvents := map[string]Unmarshaller{}
	for _, e := range evs {
		typ := reflect.TypeOf(e)
		registeredEvents[typ.String()] = e
	}

	outchan := make(chan Event)
	go func() {
		for {
			e := <-c
			et := e.Metadata[MetadatakeyEventType]
			ev, ok := registeredEvents[et]
			if !ok {
				continue
			}

			event, err := ev.Unmarshal(e.Payload)
			if err != nil {
				log.Printf("can't unmarshal event %v", err)
				continue
			}

			outchan <- Event{
				Type:  et,
				ID:    e.Metadata[MetadatakeyEventID],
				Event: event,
			}
		}
	}()
	return outchan, nil
}

// ConsumeAll allows consuming all events. Note that unmarshalling must be done manually in this case, therefore Event.Event will always be of type []byte
func ConsumeAll(s Consumer, group string) (<-chan Event, error) {
	c, err := s.Consume(MainQueueName, events.WithGroup(group))
	if err != nil {
		return nil, err
	}

	outchan := make(chan Event)
	go func() {
		for {
			e := <-c
			outchan <- Event{
				Type:  e.Metadata[MetadatakeyEventType],
				ID:    e.Metadata[MetadatakeyEventID],
				Event: e.Payload,
			}
		}
	}()
	return outchan, nil
}

// Publish publishes the ev to the MainQueue from where it is distributed to all subscribers
// NOTE: needs to use reflect on runtime
func Publish(s Publisher, ev interface{}) error {
	evName := reflect.TypeOf(ev).String()
	return s.Publish(MainQueueName, ev, events.WithMetadata(map[string]string{
		MetadatakeyEventType: evName,
		MetadatakeyEventID:   uuid.New().String(),
	}))
}
