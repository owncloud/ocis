package memlimit

import (
	"errors"
	"runtime/debug"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/c2h5oh/datasize"
)

var (
	setEnabled *bool
	setRatio   *float64
	setAmount  *int64
)

func Set(enabled bool, ratio float64, amount string) error {

	if setEnabled != nil {
		if *setEnabled != enabled {
			return errors.New("")
		}
	}

	if amount == "" {
		if setAmount != nil {
			return errors.New("memory is already limited by amount, cannot be limited by ratio")
		}

		if setRatio != nil && *setRatio != ratio {
			return errors.New("memory is already limited by different ratio than given")
		}

		if ratio <= 0 || ratio > 1 {
			return errors.New("ratio must be greater 0 and not greater 1")
		}

	}

	var bytesAmount int64 = 0
	if amount != "" {

		if setRatio != nil {
			return errors.New("memory is already limited by ratio, cannot be limited by amount")
		}

		var v datasize.ByteSize
		err := v.UnmarshalText([]byte(amount))
		if err != nil {
			return errors.New("invalid memory amount given")
		}

		bytesAmount = int64(v.Bytes())

		if bytesAmount == 0 {
			return errors.New("zero bytes is no valid memory limit")
		}

	}

	if setAmount != nil {
		if *setAmount != bytesAmount {
			return errors.New("memory is already limited by different amount than given")
		}
	}

	// actually apply settings
	if !enabled {
		t := false
		setEnabled = &t
		return nil
	}
	t := true
	setEnabled = &t

	if bytesAmount == 0 {
		r := ratio
		setRatio = &r
		memlimit.SetGoMemLimit(ratio)
		return nil
	}

	a := bytesAmount
	setAmount = &a
	debug.SetMemoryLimit(bytesAmount)
	return nil
}
