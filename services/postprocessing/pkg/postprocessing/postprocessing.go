package postprocessing

import (
	"time"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
)

// Postprocessing handles postprocessing of a file
type Postprocessing struct {
	id       string
	url      string
	u        *user.User
	m        map[events.Postprocessingstep]interface{}
	filename string
	c        config.Postprocessing
	steps    []events.Postprocessingstep
}

// New returns a new postprocessing instance
func New(uploadID string, uploadURL string, user *user.User, filename string, c config.Postprocessing) *Postprocessing {
	return &Postprocessing{
		id:       uploadID,
		url:      uploadURL,
		u:        user,
		m:        make(map[events.Postprocessingstep]interface{}),
		c:        c,
		filename: filename,
		steps:    getSteps(c),
	}
}

// Init is the first step of the postprocessing
func (pp *Postprocessing) Init(ev events.BytesReceived) interface{} {
	pp.m["init"] = ev

	if len(pp.steps) == 0 {
		return pp.finished(events.PPOutcomeContinue)
	}

	return pp.nextStep(pp.steps[0])
}

// Virusscan is the virusscanning step of the postprocessing
func (pp *Postprocessing) Virusscan(ev events.VirusscanFinished) interface{} {
	pp.m[events.PPStepAntivirus] = ev

	switch ev.Outcome {
	case events.PPOutcomeContinue:
		return pp.next(events.PPStepAntivirus)
	default:
		return pp.finished(ev.Outcome)

	}
}

// Delay will sleep the configured time then continue
func (pp *Postprocessing) Delay(ev events.StartPostprocessingStep) interface{} {
	pp.m[events.PPStepDelay] = ev
	time.Sleep(pp.c.Delayprocessing)
	return pp.next(events.PPStepDelay)
}

func (pp *Postprocessing) next(current events.Postprocessingstep) interface{} {
	l := len(pp.steps)
	for i, s := range pp.steps {
		if s == current && i+1 < l {
			return pp.next(pp.steps[i+1])
		}
	}
	return pp.finished(events.PPOutcomeContinue)
}

func (pp *Postprocessing) nextStep(next events.Postprocessingstep) events.StartPostprocessingStep {
	return events.StartPostprocessingStep{
		UploadID:      pp.id,
		URL:           pp.url,
		ExecutingUser: pp.u,
		Filename:      pp.filename,
		StepToStart:   next,
	}
}

func (pp *Postprocessing) finished(outcome events.PostprocessingOutcome) events.PostprocessingFinished {
	return events.PostprocessingFinished{
		UploadID:      pp.id,
		Result:        pp.m,
		ExecutingUser: pp.u,
		Filename:      pp.filename,
		Outcome:       outcome,
	}
}

func getSteps(c config.Postprocessing) []events.Postprocessingstep {
	// NOTE: first version only contains very basic configuration options
	// But we aim for a system where postprocessing steps and their order can be configured per space
	// ideally by the spaceadmin itself
	// We need to iterate over configuring PP service when we see fit
	var steps []events.Postprocessingstep
	if c.Virusscan {
		steps = append(steps, events.PPStepAntivirus)
	}

	if c.FTSIndex {
		steps = append(steps, events.PPStepFTS)
	}

	if c.Delayprocessing != 0 {
		steps = append(steps, events.PPStepDelay)
	}
	return steps
}
