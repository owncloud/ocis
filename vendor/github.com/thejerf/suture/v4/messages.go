package suture

// sum type pattern for type-safe message passing; see
// http://www.jerf.org/iri/post/2917

type supervisorMessage interface {
	isSupervisorMessage()
}

type listServices struct {
	c chan []Service
}

func (ls listServices) isSupervisorMessage() {}

type removeService struct {
	id           serviceID
	notification chan struct{}
}

func (rs removeService) isSupervisorMessage() {}

func (s *Supervisor) sync() {
	select {
	case s.control <- syncSupervisor{}:
	case <-s.liveness:
	}
}

type syncSupervisor struct {
}

func (ss syncSupervisor) isSupervisorMessage() {}

func (s *Supervisor) fail(id serviceID, panicVal interface{}, stacktrace []byte) {
	select {
	case s.control <- serviceFailed{id, panicVal, stacktrace}:
	case <-s.liveness:
	}
}

type serviceFailed struct {
	id         serviceID
	panicVal   interface{}
	stacktrace []byte
}

func (sf serviceFailed) isSupervisorMessage() {}

func (s *Supervisor) serviceEnded(id serviceID, err error) {
	s.sendControl(serviceEnded{id, err})
}

type serviceEnded struct {
	id  serviceID
	err error
}

func (s serviceEnded) isSupervisorMessage() {}

// added by the Add() method
type addService struct {
	service  Service
	name     string
	response chan serviceID
}

func (as addService) isSupervisorMessage() {}

type stopSupervisor struct {
	done chan UnstoppedServiceReport
}

func (ss stopSupervisor) isSupervisorMessage() {}

func (s *Supervisor) panic() {
	s.control <- panicSupervisor{}
}

type panicSupervisor struct {
}

func (ps panicSupervisor) isSupervisorMessage() {}
