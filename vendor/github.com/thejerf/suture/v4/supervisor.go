package suture

// FIXMES in progress:
// 1. Ensure the supervisor actually gets to the terminated state for the
//     unstopped service report.
// 2. Save the unstopped service report in the supervisor.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	notRunning = iota
	normal
	paused
	terminated
)

type supervisorID uint32
type serviceID uint32

// ErrSupervisorNotRunning is returned by some methods if the supervisor is
// not running, either because it has not been started or because it has
// been terminated.
var ErrSupervisorNotRunning = errors.New("supervisor not running")

/*
Supervisor is the core type of the module that represents a Supervisor.

Supervisors should be constructed either by New or NewSimple.

Once constructed, a Supervisor should be started in one of three ways:

 1. Calling .Serve(ctx).
 2. Calling .ServeBackground(ctx).
 3. Adding it to an existing Supervisor.

Calling Serve will cause the supervisor to run until the passed-in
context is cancelled. Often one of the last lines of the "main" func for a
program will be to call one of the Serve methods.

Calling ServeBackground will CORRECTLY start the supervisor running in a
new goroutine. It is risky to directly run

	go supervisor.Serve()

because that will briefly create a race condition as it starts up, if you
try to .Add() services immediately afterward.
*/
type Supervisor struct {
	Name string

	spec Spec

	services             map[serviceID]serviceWithName
	cancellations        map[serviceID]context.CancelFunc
	servicesShuttingDown map[serviceID]serviceWithName
	lastFail             time.Time
	failures             float64
	restartQueue         []serviceID
	serviceCounter       serviceID
	control              chan supervisorMessage
	notifyServiceDone    chan serviceID
	resumeTimer          <-chan time.Time
	liveness             chan struct{}

	// despite the recommendation in the context package to avoid
	// holding this in a struct, I think due to the function of suture
	// and the way it works, I think it's OK in this case. This is the
	// exceptional case, basically.
	ctxMutex sync.Mutex
	ctx      context.Context
	// This function cancels this supervisor specifically.
	ctxCancel func()

	getNow       func() time.Time
	getAfterChan func(time.Duration) <-chan time.Time

	m sync.Mutex

	// The unstopped service report is generated when we finish
	// stopping.
	unstoppedServiceReport UnstoppedServiceReport

	// malign leftovers
	id    supervisorID
	state uint8
}

/*
New is the full constructor function for a supervisor.

The name is a friendly human name for the supervisor, used in logging. Suture
does not care if this is unique, but it is good for your sanity if it is.

If not set, the following values are used:

  - EventHook:         A function is created that uses log.Print.
  - FailureDecay:      30 seconds
  - FailureThreshold:  5 failures
  - FailureBackoff:    15 seconds
  - Timeout:           10 seconds
  - BackoffJitter:     DefaultJitter

The EventHook function will be called when errors occur. Suture will log the
following:

  - When a service has failed, with a descriptive message about the
    current backoff status, and whether it was immediately restarted
  - When the supervisor has gone into its backoff mode, and when it
    exits it
  - When a service fails to stop

A default hook for slog is provided with the [sutureslog
module](https://github.com/thejerf/sutureslog).

The failureRate, failureThreshold, and failureBackoff controls how failures
are handled, in order to avoid the supervisor failure case where the
program does nothing but restarting failed services. If you do not
care how failures behave, the default values should be fine for the
vast majority of services, but if you want the details:

The supervisor tracks the number of failures that have occurred, with an
exponential decay on the count. Every FailureDecay seconds, the number of
failures that have occurred is cut in half. (This is done smoothly with an
exponential function.) When a failure occurs, the number of failures
is incremented by one. When the number of failures passes the
FailureThreshold, the entire service waits for FailureBackoff seconds
before attempting any further restarts, at which point it resets its
failure count to zero.

Timeout is how long Suture will wait for a service to properly terminate.

The PassThroughPanics options can be set to let panics in services propagate
and crash the program, should this be desirable.

DontPropagateTermination indicates whether this supervisor tree will
propagate a ErrTerminateTree if a child process returns it. If false,
this supervisor will itself return an error that will terminate its
parent. If true, it will merely return ErrDoNotRestart. false by default.
*/
func New(name string, spec Spec) *Supervisor {
	spec.configureDefaults(name)

	return &Supervisor{
		name,

		spec,

		// services
		make(map[serviceID]serviceWithName),
		// cancellations
		make(map[serviceID]context.CancelFunc),
		// servicesShuttingDown
		make(map[serviceID]serviceWithName),
		// lastFail, deliberately the zero time
		time.Time{},
		// failures
		0,
		// restartQueue
		make([]serviceID, 0, 1),
		// serviceCounter
		0,
		// control
		make(chan supervisorMessage),
		// notifyServiceDone
		make(chan serviceID),
		// resumeTimer
		make(chan time.Time),

		// liveness
		make(chan struct{}),

		sync.Mutex{},
		// ctx
		nil,
		// myCancel
		nil,

		// the tests can override these for testing threshold
		// behavior
		// getNow
		time.Now,
		// getAfterChan
		time.After,

		// m
		sync.Mutex{},

		// unstoppedServiceReport
		nil,

		// id
		nextSupervisorID(),
		// state
		notRunning,
	}
}

func serviceName(service Service) (serviceName string) {
	stringer, canStringer := service.(fmt.Stringer)
	if canStringer {
		serviceName = stringer.String()
	} else {
		serviceName = fmt.Sprintf("%#v", service)
	}
	return
}

// NewSimple is a convenience function to create a service with just a name
// and the sensible defaults.
func NewSimple(name string) *Supervisor {
	return New(name, Spec{})
}

// HasSupervisor is an interface that indicates the given struct contains a
// supervisor. If the struct is either already a *Supervisor, or it embeds
// a *Supervisor, this will already be implemented for you. Otherwise, a
// struct containing a supervisor will need to implement this in order to
// participate in the log function propagation and recursive
// UnstoppedService report.
//
// It is legal for GetSupervisor to return nil, in which case
// the supervisor-specific behaviors will simply be ignored.
type HasSupervisor interface {
	GetSupervisor() *Supervisor
}

func (s *Supervisor) GetSupervisor() *Supervisor {
	return s
}

/*
Add adds a service to this supervisor.

If the supervisor is currently running, the service will be started
immediately. If the supervisor has not been started yet, the service
will be started when the supervisor is. If the supervisor was already stopped,
this is a no-op returning an empty service-token.

The returned ServiceID may be passed to the Remove method of the Supervisor
to terminate the service.

As a special behavior, if the service added is itself a supervisor, the
supervisor being added will copy the EventHook function from the Supervisor it
is being added to. This allows factoring out providing a Supervisor
from its logging. This unconditionally overwrites the child Supervisor's
logging functions.
*/
func (s *Supervisor) Add(service Service) ServiceToken {
	if s == nil {
		panic("can't add service to nil *suture.Supervisor")
	}

	if hasSupervisor, isHaveSupervisor := service.(HasSupervisor); isHaveSupervisor {
		supervisor := hasSupervisor.GetSupervisor()
		if supervisor != nil {
			supervisor.spec.EventHook = s.spec.EventHook
		}
	}

	s.m.Lock()
	if s.state == notRunning {
		id := s.serviceCounter
		s.serviceCounter++

		s.services[id] = serviceWithName{service, serviceName(service)}
		s.restartQueue = append(s.restartQueue, id)

		s.m.Unlock()
		return ServiceToken{supervisor: s.id, service: id}
	}
	s.m.Unlock()

	response := make(chan serviceID)
	if s.sendControl(addService{service, serviceName(service), response}) != nil {
		return ServiceToken{}
	}
	return ServiceToken{supervisor: s.id, service: <-response}
}

// ServeBackground starts running a supervisor in its own goroutine. When
// this method returns, the supervisor is guaranteed to be in a running state.
// The returned one-buffered channel receives the error returned by .Serve.
func (s *Supervisor) ServeBackground(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.Serve(ctx)
	}()
	s.sync()
	return errChan
}

/*
Serve starts the supervisor. You should call this on the top-level supervisor,
but nothing else.
*/
func (s *Supervisor) Serve(ctx context.Context) error {
	// context documentation suggests that it is legal for functions to
	// take nil contexts, it's user's responsibility to never pass them in.
	if ctx == nil {
		ctx = context.Background()
	}

	if s == nil {
		panic("Can't serve with a nil *suture.Supervisor")
	}
	// Take a separate cancellation function so this tree can be
	// indepedently cancelled.
	ctx, myCancel := context.WithCancel(ctx)
	s.ctxMutex.Lock()
	s.ctx = ctx
	s.ctxMutex.Unlock()
	s.ctxCancel = myCancel

	if s.id == 0 {
		panic("Can't call Serve on an incorrectly-constructed *suture.Supervisor")
	}

	s.m.Lock()
	if s.state == normal || s.state == paused {
		s.m.Unlock()
		panic("Called .Serve() on a supervisor that is already Serve()ing")
	}

	s.state = normal
	s.m.Unlock()

	defer func() {
		s.m.Lock()
		s.state = terminated
		s.m.Unlock()
	}()

	// for all the services I currently know about, start them
	for _, id := range s.restartQueue {
		namedService, present := s.services[id]
		if present {
			s.runService(ctx, namedService.Service, id)
		}
	}
	s.restartQueue = make([]serviceID, 0, 1)

	for {
		select {
		case <-ctx.Done():
			s.stopSupervisor()
			return ctx.Err()
		case m := <-s.control:
			switch msg := m.(type) {
			case serviceFailed:
				s.handleFailedService(ctx, msg.id, msg.panicVal, msg.stacktrace, true)
			case serviceEnded:
				_, monitored := s.services[msg.id]
				if monitored {
					cancel := s.cancellations[msg.id]

					if isErr(msg.err, ErrDoNotRestart) || ctx.Err() != nil {
						delete(s.services, msg.id)
						delete(s.cancellations, msg.id)
						go cancel()
					} else if isErr(msg.err, ErrTerminateSupervisorTree) {
						s.stopSupervisor()
						if s.spec.DontPropagateTermination {
							return ErrDoNotRestart
						} else {
							return msg.err
						}
					} else {
						err := msg.err
						if isErr(msg.err, context.DeadlineExceeded) || isErr(msg.err, context.Canceled) {
							err = fmt.Errorf("from some other context, not the service's context, so service is being restarted: %w", msg.err)
						}
						s.handleFailedService(ctx, msg.id, err, nil, false)
					}
				}
			case addService:
				id := s.serviceCounter
				s.serviceCounter++

				s.services[id] = serviceWithName{msg.service, msg.name}
				s.runService(ctx, msg.service, id)

				msg.response <- id
			case removeService:
				s.removeService(msg.id, msg.notification)
			case stopSupervisor:
				msg.done <- s.stopSupervisor()
				return nil
			case listServices:
				services := []Service{}
				for _, service := range s.services {
					services = append(services, service.Service)
				}
				msg.c <- services
			case syncSupervisor:
				// this does nothing on purpose; its sole purpose is to
				// introduce a sync point via the channel receive
			case panicSupervisor:
				// used only by tests
				panic("Panicking as requested!")
			}
		case serviceEnded := <-s.notifyServiceDone:
			delete(s.servicesShuttingDown, serviceEnded)
		case <-s.resumeTimer:
			// We're resuming normal operation after a pause due to
			// excessive thrashing
			// FIXME: Ought to permit some spacing of these functions, rather
			// than simply hammering through them
			s.m.Lock()
			s.state = normal
			s.m.Unlock()
			s.failures = 0
			s.spec.EventHook(EventResume{s, s.Name})
			for _, id := range s.restartQueue {
				namedService, present := s.services[id]
				if present {
					s.runService(ctx, namedService.Service, id)
				}
			}
			s.restartQueue = make([]serviceID, 0, 1)
		}
	}
}

// UnstoppedServiceReport will return a report of what services failed to
// stop when the supervisor was stopped. This call will return when the
// supervisor is done shutting down. It will hang on a supervisor that has
// not been stopped, because it will not be "done shutting down".
//
// Calling this on a supervisor will return a report for the whole
// supervisor tree under it.
//
// WARNING: Technically, any use of the returned data structure is a
// TOCTOU violation:
// https://en.wikipedia.org/wiki/Time-of-check_to_time-of-use
// Since the data structure was generated and returned to you, any of these
// services may have stopped since then.
//
// However, this can still be useful information at program teardown
// time. For instance, logging that a service failed to stop as expected is
// still useful, as even if it shuts down later, it was still later than
// you expected.
//
// But if you cast the Service objects back to their underlying objects and
// start trying to manipulate them ("shut down harder!"), be sure to
// account for the possibility they are in fact shut down before you get
// them.
//
// If there are no services to report, the UnstoppedServiceReport will be
// nil. A zero-length constructed slice is never returned.
func (s *Supervisor) UnstoppedServiceReport() (UnstoppedServiceReport, error) {
	// the only thing that ever happens to this channel is getting
	// closed when the supervisor terminates.
	_, _ = <-s.liveness

	// FIXME: Recurse on the supervisors
	return s.unstoppedServiceReport, nil
}

func (s *Supervisor) handleFailedService(ctx context.Context, id serviceID, err interface{}, stacktrace []byte, panic bool) {
	now := s.getNow()

	if s.lastFail.IsZero() {
		s.lastFail = now
		s.failures = 1.0
	} else {
		sinceLastFail := now.Sub(s.lastFail).Seconds()
		intervals := sinceLastFail / s.spec.FailureDecay
		s.failures = s.failures*math.Pow(.5, intervals) + 1
	}

	if s.failures > s.spec.FailureThreshold {
		s.m.Lock()
		s.state = paused
		s.m.Unlock()
		s.spec.EventHook(EventBackoff{s, s.Name})
		s.resumeTimer = s.getAfterChan(
			s.spec.BackoffJitter.Jitter(s.spec.FailureBackoff))
	}

	s.lastFail = now

	failedService, monitored := s.services[id]

	// It is possible for a service to be no longer monitored
	// by the time we get here. In that case, just ignore it.
	if monitored {
		s.m.Lock()
		curState := s.state
		s.m.Unlock()
		if curState == normal {
			s.runService(ctx, failedService.Service, id)
		} else {
			s.restartQueue = append(s.restartQueue, id)
		}
		if panic {
			s.spec.EventHook(EventServicePanic{
				Supervisor:       s,
				SupervisorName:   s.Name,
				Service:          failedService.Service,
				ServiceName:      failedService.name,
				CurrentFailures:  s.failures,
				FailureThreshold: s.spec.FailureThreshold,
				Restarting:       curState == normal,
				PanicMsg:         s.spec.Sprint(err),
				Stacktrace:       string(stacktrace),
			})
		} else {
			e := EventServiceTerminate{
				Supervisor:       s,
				SupervisorName:   s.Name,
				Service:          failedService.Service,
				ServiceName:      failedService.name,
				CurrentFailures:  s.failures,
				FailureThreshold: s.spec.FailureThreshold,
				Restarting:       curState == normal,
			}
			if err != nil {
				e.Err = err
			}
			s.spec.EventHook(e)
		}
	}
}

func (s *Supervisor) runService(ctx context.Context, service Service, id serviceID) {
	childCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	blockingCancellation := func() {
		cancel()
		<-done
	}
	s.cancellations[id] = blockingCancellation
	go func() {
		if !s.spec.PassThroughPanics {
			defer func() {
				if r := recover(); r != nil {
					buf := make([]byte, 65535)
					written := runtime.Stack(buf, false)
					buf = buf[:written]
					s.fail(id, r, buf)
				}
			}()
		}

		var err error

		defer func() {
			cancel()
			close(done)

			r := recover()
			if r == nil {
				s.serviceEnded(id, err)
			} else {
				panic(r)
			}
		}()

		err = service.Serve(childCtx)
	}()
}

func (s *Supervisor) removeService(id serviceID, notificationChan chan struct{}) {
	namedService, present := s.services[id]
	if present {
		cancel := s.cancellations[id]
		delete(s.services, id)
		delete(s.cancellations, id)

		s.servicesShuttingDown[id] = namedService
		go func() {
			successChan := make(chan struct{})
			go func() {
				cancel()
				close(successChan)
				if notificationChan != nil {
					notificationChan <- struct{}{}
				}
			}()

			select {
			case <-successChan:
				// Life is good!
			case <-s.getAfterChan(s.spec.Timeout):
				s.spec.EventHook(EventStopTimeout{
					s, s.Name,
					namedService.Service, namedService.name})
			}
			s.notifyServiceDone <- id
		}()
	} else {
		if notificationChan != nil {
			notificationChan <- struct{}{}
		}
	}
}

func (s *Supervisor) stopSupervisor() UnstoppedServiceReport {
	notifyDone := make(chan serviceID, len(s.services))

	for id, namedService := range s.services {
		cancel := s.cancellations[id]
		delete(s.services, id)
		delete(s.cancellations, id)
		s.servicesShuttingDown[id] = namedService
		go func(sID serviceID) {
			cancel()
			notifyDone <- sID
		}(id)
	}

	timeout := s.getAfterChan(s.spec.Timeout)

SHUTTING_DOWN_SERVICES:
	for len(s.servicesShuttingDown) > 0 {
		select {
		case id := <-notifyDone:
			delete(s.servicesShuttingDown, id)
		case serviceID := <-s.notifyServiceDone:
			delete(s.servicesShuttingDown, serviceID)
		case <-timeout:
			for _, namedService := range s.servicesShuttingDown {
				s.spec.EventHook(EventStopTimeout{
					s, s.Name,
					namedService.Service, namedService.name,
				})
			}

			// failed remove statements will log the errors.
			break SHUTTING_DOWN_SERVICES
		}
	}

	// If nothing else has cancelled our context, we should now.
	s.ctxCancel()

	// Indicate that we're done shutting down
	defer close(s.liveness)

	if len(s.servicesShuttingDown) == 0 {
		return nil
	} else {
		report := UnstoppedServiceReport{}
		for serviceID, serviceWithName := range s.servicesShuttingDown {
			report = append(report, UnstoppedService{
				SupervisorPath: []*Supervisor{s},
				Service:        serviceWithName.Service,
				Name:           serviceWithName.name,
				ServiceToken:   ServiceToken{supervisor: s.id, service: serviceID},
			})
		}
		s.m.Lock()
		s.unstoppedServiceReport = report
		s.m.Unlock()
		return report
	}
}

// String implements the fmt.Stringer interface.
func (s *Supervisor) String() string {
	return s.Name
}

// sendControl abstracts checking for the supervisor to still be running
// when we send a message. This prevents blocking when sending to a
// cancelled supervisor.
func (s *Supervisor) sendControl(sm supervisorMessage) error {
	var doneChan <-chan struct{}
	s.ctxMutex.Lock()
	if s.ctx == nil {
		s.ctxMutex.Unlock()
		return ErrSupervisorNotStarted
	}
	doneChan = s.ctx.Done()
	s.ctxMutex.Unlock()

	select {
	case s.control <- sm:
		return nil
	case <-doneChan:
		return ErrSupervisorNotRunning
	}
}

/*
Remove will remove the given service from the Supervisor, and attempt to Stop() it.
The ServiceID token comes from the Add() call. This returns without waiting
for the service to stop.
*/
func (s *Supervisor) Remove(id ServiceToken) error {
	if id.supervisor != s.id {
		return ErrWrongSupervisor
	}
	err := s.sendControl(removeService{id.service, nil})
	if err == ErrSupervisorNotRunning {
		// No meaningful error handling if the supervisor is stopped.
		return nil
	}
	return err
}

/*
RemoveAndWait will remove the given service from the Supervisor and attempt
to Stop() it. It will wait up to the given timeout value for the service to
terminate. A timeout value of 0 means to wait forever.

If a nil error is returned from this function, then the service was
terminated normally. If either the supervisor terminates or the timeout
passes, ErrTimeout is returned. (If this isn't even the right supervisor
ErrWrongSupervisor is returned.)
*/
func (s *Supervisor) RemoveAndWait(id ServiceToken, timeout time.Duration) error {
	if id.supervisor != s.id {
		return ErrWrongSupervisor
	}

	var timeoutC <-chan time.Time

	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		timeoutC = timer.C
	}

	notificationC := make(chan struct{})

	sentControlErr := s.sendControl(removeService{id.service, notificationC})

	if sentControlErr != nil {
		return sentControlErr
	}

	select {
	case <-notificationC:
		// normal case; the service is terminated.
		return nil

	// This occurs if the entire supervisor ends without the service
	// having terminated, and includes the timeout the supervisor
	// itself waited before closing the liveness channel.
	case <-s.ctx.Done():
		return ErrTimeout

	// The local timeout.
	case <-timeoutC:
		return ErrTimeout
	}
}

/*
Services returns a []Service containing a snapshot of the services this
Supervisor is managing.
*/
func (s *Supervisor) Services() []Service {
	ls := listServices{make(chan []Service)}

	if s.sendControl(ls) == nil {
		return <-ls.c
	}
	return nil
}

var currentSupervisorID uint32

func nextSupervisorID() supervisorID {
	return supervisorID(atomic.AddUint32(&currentSupervisorID, 1))
}

// ServiceToken is an opaque identifier that can be used to terminate a service that
// has been Add()ed to a Supervisor.
type ServiceToken struct {
	supervisor supervisorID
	service    serviceID
}

// An UnstoppedService is the component member of an
// UnstoppedServiceReport.
//
// The SupervisorPath is the path down the supervisor tree to the given
// service.
type UnstoppedService struct {
	SupervisorPath []*Supervisor
	Service        Service
	Name           string
	ServiceToken   ServiceToken
}

// An UnstoppedServiceReport will be returned by StopWithReport, reporting
// which services failed to stop.
type UnstoppedServiceReport []UnstoppedService

type serviceWithName struct {
	Service Service
	name    string
}

// Jitter returns the sum of the input duration and a random jitter.  It is
// compatible with the jitter functions in github.com/lthibault/jitterbug.
type Jitter interface {
	Jitter(time.Duration) time.Duration
}

// NoJitter does not apply any jitter to the input duration
type NoJitter struct{}

// Jitter leaves the input duration d unchanged.
func (NoJitter) Jitter(d time.Duration) time.Duration { return d }

// DefaultJitter is the jitter function that is applied when spec.BackoffJitter
// is set to nil.
type DefaultJitter struct {
	rand *rand.Rand
}

// Jitter will jitter the backoff time by uniformly distributing it into
// the range [FailureBackoff, 1.5 * FailureBackoff).
func (dj *DefaultJitter) Jitter(d time.Duration) time.Duration {
	// this is only called by the core supervisor loop, so it is
	// single-thread safe.
	if dj.rand == nil {
		dj.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	jitter := dj.rand.Float64() / 2
	return d + time.Duration(float64(d)*jitter)
}

// ErrWrongSupervisor is returned by the (*Supervisor).Remove method
// if you pass a ServiceToken from the wrong Supervisor.
var ErrWrongSupervisor = errors.New("wrong supervisor for this service token, no service removed")

// ErrTimeout is returned when an attempt to RemoveAndWait for a service to
// stop has timed out.
var ErrTimeout = errors.New("waiting for service to stop has timed out")

// ErrSupervisorNotTerminated is returned when asking for a stopped service
// report before the supervisor has been terminated.
var ErrSupervisorNotTerminated = errors.New("supervisor not terminated")

// ErrSupervisorNotStarted is returned if you try to send control messages
// to a supervisor that has not started yet. See note on Supervisor struct
// about the legal ways to start a supervisor.
var ErrSupervisorNotStarted = errors.New("supervisor not started yet")

// Spec is used to pass arguments to the New function to create a
// supervisor. See the New function for full documentation.
type Spec struct {
	EventHook                EventHook
	Sprint                   SprintFunc
	FailureDecay             float64
	FailureThreshold         float64
	FailureBackoff           time.Duration
	BackoffJitter            Jitter
	Timeout                  time.Duration
	PassThroughPanics        bool
	DontPropagateTermination bool
}

func (s *Spec) configureDefaults(supervisorName string) {
	if s.FailureDecay == 0 {
		s.FailureDecay = 30
	}
	if s.FailureThreshold == 0 {
		s.FailureThreshold = 5
	}
	if s.FailureBackoff == 0 {
		s.FailureBackoff = time.Second * 15
	}
	if s.BackoffJitter == nil {
		s.BackoffJitter = &DefaultJitter{}
	}
	if s.Timeout == 0 {
		s.Timeout = time.Second * 10
	}

	// set up the default logging handlers
	if s.EventHook == nil {
		s.EventHook = func(e Event) {
			log.Print(e)
		}
	}

	if s.Sprint == nil {
		s.Sprint = func(v interface{}) string {
			return fmt.Sprintf("%v", v)
		}
	}
}
