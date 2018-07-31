package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	AppName          = "dt"
	unixSeconds      = "unix"
	unixMilliSeconds = "unixm"
)

var version = "0.10.0"

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

func loadConfig() {
	path, err := homedir.Expand("~/.config/dt/.dt")
	if err != nil {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	m := map[string]string{}
	for scanner.Scan() {
		k, v := SplitFormat(scanner.Text())
		if k == "" || v == "" {
			continue
		}
		m[k] = v
	}

	for k, v := range m {
		log.Printf("custom format: %s => %s\n", k, v)
		formats[k] = v
	}
}

func SplitFormat(s string) (string, string) {
	cols := splitRegexp.Split(s, 2)
	if len(cols) != 2 {
		return "", ""
	}

	if cols[0] == "" || cols[1] == "" {
		return "", ""
	}

	return cols[0], cols[1]
}

func main() {
	cli := &CLO{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
