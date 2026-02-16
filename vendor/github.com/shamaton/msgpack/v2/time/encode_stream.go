package time

import (
	"reflect"
	"time"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
)

var StreamEncoder = new(timeStreamEncoder)

type timeStreamEncoder struct{}

var _ ext.StreamEncoder = (*timeStreamEncoder)(nil)

func (timeStreamEncoder) Code() int8 {
	return def.TimeStamp
}

func (timeStreamEncoder) Type() reflect.Type {
	return typeOf
}

func (e timeStreamEncoder) Write(w ext.StreamWriter, value reflect.Value) error {
	t := value.Interface().(time.Time)

	secs := uint64(t.Unix())
	if secs>>34 == 0 {
		data := uint64(t.Nanosecond())<<34 | secs
		if data&0xffffffff00000000 == 0 {
			if err := w.WriteByte1Int(def.Fixext4); err != nil {
				return err
			}
			if err := w.WriteByte1Int(def.TimeStamp); err != nil {
				return err
			}
			if err := w.WriteByte4Uint64(data); err != nil {
				return err
			}
			return nil
		}

		if err := w.WriteByte1Int(def.Fixext8); err != nil {
			return err
		}
		if err := w.WriteByte1Int(def.TimeStamp); err != nil {
			return err
		}
		if err := w.WriteByte8Uint64(data); err != nil {
			return err
		}
		return nil
	}

	if err := w.WriteByte1Int(def.Ext8); err != nil {
		return err
	}
	if err := w.WriteByte1Int(12); err != nil {
		return err
	}
	if err := w.WriteByte1Int(def.TimeStamp); err != nil {
		return err
	}
	if err := w.WriteByte4Int(t.Nanosecond()); err != nil {
		return err
	}
	if err := w.WriteByte8Uint64(secs); err != nil {
		return err
	}
	return nil
}
