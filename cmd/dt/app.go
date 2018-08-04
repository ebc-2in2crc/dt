package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	// AppName コマンドの名前
	AppName          = "dt"
	unixSeconds      = "unix"
	unixMilliSeconds = "unixm"
)

// AdjustDay 対応する月に同じ日が存在しないときの調整
type AdjustDay int

const (
	// AdjustToEndOfMonth 対応する月に同じ日が存在しないときは代わりにその月の末日が使われます
	AdjustToEndOfMonth AdjustDay = iota
	// Normalize 対応する月に同じ日が存在しないときは time.Time.	Date と同じ方法で正規化します
	Normalize
)

var version = "0.10.0"
var splitRegexp = regexp.MustCompile(`\s*=\s*`)

// Dt 日付計算とフォーマット機能をもつ
type Dt struct {
	time   time.Time
	format string
}

func (dt *Dt) get() time.Time {
	return dt.time
}

// AddYear 月を加算. 負値のときは減算.
func (dt *Dt) AddYear(year int) *Dt {
	return &Dt{
		time:   dt.time.AddDate(year, 0, 0),
		format: dt.format,
	}
}

// AddMonth 月を加算. 負値のときは減算.
func (dt *Dt) AddMonth(month int, adjust AdjustDay) *Dt {
	result := &Dt{
		time:   dt.time.AddDate(0, month, 0),
		format: dt.format,
	}
	if adjust == Normalize {
		return result
	}

	t := dt.time
	firstDayOfMonth := time.Date(t.Year(), t.Month(), 1, t.Hour(),
		t.Minute(), t.Second(), t.Nanosecond(), t.Location()).AddDate(0, month, 0)
	if result.time.Month() == firstDayOfMonth.Month() {
		return result
	}

	lastDayOfPreviousMonth := firstDayOfMonth.AddDate(0, 1, -1)
	return &Dt{
		time:   lastDayOfPreviousMonth,
		format: dt.format,
	}
}

// AddDay 日を加算. 負値のときは減算.
func (dt *Dt) AddDay(day int) *Dt {
	return &Dt{
		time:   dt.time.AddDate(0, 0, day),
		format: dt.format,
	}
}

// AddHour 時を加算. 負値のときは減算.
func (dt *Dt) AddHour(hour int) *Dt {
	return &Dt{
		time:   dt.time.Add(time.Duration(hour) * time.Hour),
		format: dt.format,
	}
}

// AddMinute 分を加算. 負値のときは減算.
func (dt *Dt) AddMinute(minute int) *Dt {
	return &Dt{
		time:   dt.time.Add(time.Duration(minute) * time.Minute),
		format: dt.format,
	}
}

// AddSecond 秒を加算. 負値のときは減算.
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
		k, v := splitFormat(scanner.Text())
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

func splitFormat(s string) (string, string) {
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
