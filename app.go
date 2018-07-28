package main

import (
	"fmt"
	"os"
	"time"
)

const (
	AppName          = "dt"
	Version          = "0.10.0"
	unixSeconds      = "unix"
	unixMilliSeconds = "unixm"
)

type Dt struct {
	time   time.Time
	format string
}

func (dt *Dt) Get() time.Time {
	return dt.time
}

func (dt *Dt) AddYear(year int) *Dt {
	return &Dt{
		time:   dt.time.AddDate(year, 0, 0),
		format: dt.format,
	}
}

func (dt *Dt) AddMonth(month int) *Dt {
	return &Dt{
		time:   dt.time.AddDate(0, month, 0),
		format: dt.format,
	}
}

func (dt *Dt) AddDay(day int) *Dt {
	return &Dt{
		time:   dt.time.AddDate(0, 0, day),
		format: dt.format,
	}
}

func (dt *Dt) AddHour(hour int) *Dt {
	return &Dt{
		time:   dt.time.Add(time.Duration(hour) * time.Hour),
		format: dt.format,
	}
}

func (dt *Dt) AddMinute(minute int) *Dt {
	return &Dt{
		time:   dt.time.Add(time.Duration(minute) * time.Minute),
		format: dt.format,
	}
}

func (dt *Dt) AddSecond(second int) *Dt {
	return &Dt{
		time:   dt.time.Add(time.Duration(second) * time.Second),
		format: dt.format,
	}
}

func (dt *Dt) String() string {
	t := dt.time
	f := dt.format
	switch f {
	case unixSeconds:
		return fmt.Sprintf("%d", t.Unix())
	case unixMilliSeconds:
		return fmt.Sprintf("%d", t.UnixNano()/int64(time.Millisecond))
	default:
		return t.Format(f)
	}
}

func main() {
	cli := &CLO{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
