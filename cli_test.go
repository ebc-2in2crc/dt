package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

type MyTime struct{}

func (n *MyTime) Now() time.Time {
	return time.Date(2018, 5, 12, 17, 30, 0, 0, time.Local)
}

func TestSplitFormat(t *testing.T) {
	params := []struct {
		input string
		key   string
		value string
	}{
		{input: "a=b", key: "a", value: "b"},
		{input: "a = b", key: "a", value: "b"},
		{input: "a = b = c", key: "a", value: "b = c"},
	}

	for _, p := range params {
		actualKey, actualValue := SplitFormat(p.input)
		if actualKey != p.key || actualValue != p.value {
			t.Errorf("SplitFormat() = %s => %s, want %s => %s", actualKey, actualValue, p.key, p.value)
		}
	}
}

func TestNow(t *testing.T) {
	myTime := &MyTime{}
	nowInterface = myTime

	actual := now()
	expect := myTime.Now()
	if actual != expect {
		t.Errorf("now() = %v, want %v", actual, expect)
	}
}

func TestRun_versionFlag(t *testing.T) {
	params := []struct {
		argstr string
	}{
		{argstr: AppName + " --version"},
		{argstr: AppName + " -v"},
	}

	for _, p := range params {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		clo := &CLO{outStream: outStream, errStream: errStream}

		args := strings.Split(p.argstr, " ")
		status := clo.Run(args)
		if status != ExitCodeOK {
			t.Errorf("Run(%s): ExitStatus = %d; want %d", p.argstr, status, ExitCodeOK)
		}

		actual := outStream.String()
		expect := fmt.Sprintf(AppName+" version %s", version)
		if strings.Contains(actual, expect) == false {
			t.Errorf("Run(%s): Output = %q; want %q", p.argstr, actual, expect)
		}
	}
}

func TestRun_outputEqualsInput(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	clo := &CLO{outStream: outStream, errStream: errStream}

	args := []string{AppName, "2018/05/12 17:30:00"}
	status := clo.Run(args)
	if status != ExitCodeOK {
		t.Errorf("Run(%s): ExitStatus = %d; want %d", args, status, ExitCodeOK)
	}

	actual := outStream.String()
	expect := fmt.Sprintf("2018/05/12 17:30:00")
	if strings.Contains(actual, expect) == false {
		t.Errorf("Run(%s): Output = %q; want %q", args, actual, expect)
	}
}

func TestRun_add(t *testing.T) {
	nowInterface = &MyTime{}
	params := []struct {
		args   []string
		expect string
	}{
		// 日付の計算
		{args: []string{AppName, "2018/05/12 17:30:00", "+1Y"}, expect: "2019/05/12 17:30:00"},
		{args: []string{AppName, "2018/05/12 17:30:00", "-1Y"}, expect: "2017/05/12 17:30:00"},

		{args: []string{AppName, "2018/05/12 17:30:00", "+1M"}, expect: "2018/06/12 17:30:00"},
		{args: []string{AppName, "2018/05/12 17:30:00", "-1M"}, expect: "2018/04/12 17:30:00"},

		{args: []string{AppName, "2018/05/12 17:30:00", "+1D"}, expect: "2018/05/13 17:30:00"},
		{args: []string{AppName, "2018/05/12 17:30:00", "-1D"}, expect: "2018/05/11 17:30:00"},

		{args: []string{AppName, "2018/05/12 17:30:00", "+1h"}, expect: "2018/05/12 18:30:00"},
		{args: []string{AppName, "2018/05/12 17:30:00", "-1h"}, expect: "2018/05/12 16:30:00"},

		{args: []string{AppName, "2018/05/12 17:30:00", "+1m"}, expect: "2018/05/12 17:31:00"},
		{args: []string{AppName, "2018/05/12 17:30:00", "-1m"}, expect: "2018/05/12 17:29:00"},

		{args: []string{AppName, "2018/05/12 17:30:00", "+1s"}, expect: "2018/05/12 17:30:01"},
		{args: []string{AppName, "2018/05/12 17:30:00", "-1s"}, expect: "2018/05/12 17:29:59"},

		// 入力フォーマット
		{args: []string{AppName, "now", "+1Y"}, expect: "2019/05/12 17:30:00"},
		{args: []string{AppName, "1526113800", "+1Y"}, expect: "1557649800"},
		{args: []string{AppName, "2018/05/12 17:30:00", "+1Y"}, expect: "2019/05/12 17:30:00"},
		{args: []string{AppName, "2018-05-12 17:30:00", "+1Y"}, expect: "2019-05-12 17:30:00"},
		{args: []string{AppName, "2018/05/12 17:30", "+1Y"}, expect: "2019/05/12 17:30"},
		{args: []string{AppName, "2018-05-12 17:30", "+1Y"}, expect: "2019-05-12 17:30"},
		{args: []string{AppName, "2018/05/12", "+1Y"}, expect: "2019/05/12"},
		{args: []string{AppName, "2018-05-12", "+1Y"}, expect: "2019-05-12"},
		{args: []string{AppName, "Sat May 12 17:30:00 2018", "+1Y"}, expect: "Sun May 12 17:30:00 2019"},               // ANSI
		{args: []string{AppName, "Sat May 12 17:30:00 UTC 2018", "+1Y"}, expect: "Sun May 12 17:30:00 UTC 2019"},       // UnixDate
		{args: []string{AppName, "Sat May 12 17:30:00 +0900 2018", "+1Y"}, expect: "Sun May 12 17:30:00 +0900 2019"},   // RubyDate
		{args: []string{AppName, "12 May 18 17:30 UTC", "+1Y"}, expect: "12 May 19 17:30 UTC"},                         // RFC822
		{args: []string{AppName, "12 May 18 17:30 +0900", "+1Y"}, expect: "12 May 19 17:30 +0900"},                     // RFC822Z
		{args: []string{AppName, "Saturday, 12-May-18 17:30:00 UTC", "+1Y"}, expect: "Sunday, 12-May-19 17:30:00 UTC"}, // RFC850
		{args: []string{AppName, "Sat, 12 May 2018 17:30:00 UTC", "+1Y"}, expect: "Sun, 12 May 2019 17:30:00 UTC"},     // RFC1123
		{args: []string{AppName, "Sat, 12 May 2018 17:30:00 +0900", "+1Y"}, expect: "Sun, 12 May 2019 17:30:00 +0900"}, // RFC1123Z
		{args: []string{AppName, "2018-05-12T17:30:00+09:00", "+1Y"}, expect: "2019-05-12T17:30:00+09:00"},             // RFC3339

		{args: []string{AppName, "now", "+1Y", "-2M", "3D"}, expect: "2019/03/15 17:30:00"},
		{args: []string{AppName, "--input-format", "unixm", "--output-format", "def", "1526113800000"}, expect: "2018/05/12 17:30:00"},
		{args: []string{AppName, "-i", "unixm", "-o", "def", "1526113800000"}, expect: "2018/05/12 17:30:00"},

		{args: []string{AppName, "now"}, expect: "2018/05/12 17:30:00"},
		{args: []string{AppName, "--output-format", "def", "1526113800"}, expect: "2018/05/12 17:30:00"},
		{args: []string{AppName, "-o", "def", "1526113800"}, expect: "2018/05/12 17:30:00"},
		{args: []string{AppName, "--output-format", "unix", "now"}, expect: "1526113800"},
		{args: []string{AppName, "--output-format", "2006-01-02 15:04:05", "1526113800"}, expect: "2018-05-12 17:30:00"},
		{args: []string{AppName, "--output-format", "ANSIC", "1526113800"}, expect: "Sat May 12 17:30:00 2018"},
		{args: []string{AppName, "--input-format", "unixm", "--output-format", "unixm", "1526113800000"}, expect: "1526113800000"},

		// 引数がないときはシステム日付を出力
		{args: []string{AppName}, expect: "2018/05/12 17:30:00"},
	}

	for _, p := range params {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		clo := &CLO{outStream: outStream, errStream: errStream}

		args := p.args
		status := clo.Run(args)
		if status != ExitCodeOK {
			t.Errorf("Run(%s): ExitStatus = %d; want %d", args, status, ExitCodeOK)
		}

		actual := outStream.String()
		expect := p.expect
		if strings.Contains(actual, expect) == false {
			t.Errorf("Run(%s): Output = %v; want %v", args, actual, expect)
		}
	}
}

func TestRun_error(t *testing.T) {
	nowInterface = &MyTime{}
	params := []struct {
		args   []string
		expect string
	}{
		// 指定ミス: "Y" とすべきところを "y"
		{args: []string{AppName, "now", "+1y"}, expect: "'+1y' is invalid format."},
	}

	for _, p := range params {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		clo := &CLO{outStream: outStream, errStream: errStream}

		args := p.args
		status := clo.Run(args)
		if status != ExitCodeError {
			t.Errorf("Run(%s): ExitStatus = %d; want %d", args, status, ExitCodeError)
		}

		actual := errStream.String()
		expect := p.expect
		if strings.Contains(actual, expect) == false {
			t.Errorf("Run(%s): Output = %v; want %v", args, actual, expect)
		}
	}
}
