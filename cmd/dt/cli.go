package main

import (
	"errors"
	"fmt"
	"io"

	"log"
	"regexp"
	"strconv"
	"time"

	"io/ioutil"

	"sort"

	"github.com/codegangsta/cli"
)

const (
	// ExitCodeOK コマンドが成功
	ExitCodeOK = iota

	// ExitCodeError コマンドが失敗
	ExitCodeError

	def = "def"
)

// CLO コマンドのメインの構造体
type CLO struct {
	outStream, errStream io.Writer
}

const defaultFormat = "2006/01/02 15:04:05"

var clo *CLO
var cliContext *cli.Context

var formats = map[string]string{
	"def":      defaultFormat,
	"YMDhms/":  defaultFormat,
	"YMDhms-":  "2006-01-02 15:04:05",
	"YMDhm/":   "2006/01/02 15:04",
	"YMDhm-":   "2006-01-02 15:04",
	"YMD/":     "2006/01/02",
	"YMD-":     "2006-01-02",
	"ANSIC":    time.ANSIC,
	"UnixDate": time.UnixDate,
	"RubyDate": time.RubyDate,
	"RFC822":   time.RFC822,
	"RFC822Z":  time.RFC822Z,
	"RFC850":   time.RFC850,
	"RFC1123":  time.RFC1123,
	"RFC1123Z": time.RFC1123Z,
	"RFC3339":  time.RFC3339,
}

// Run CLO のエントリーポイント
func (c *CLO) Run(args []string) int {
	clo = c

	app := cli.NewApp()
	app.Name = AppName
	app.Usage = "日付の計算や書式の変換"
	app.Version = version
	app.HideHelp = true
	app.HideVersion = true
	app.Description = description()
	app.Flags = flags()
	cli.AppHelpTemplate = appHelpTemplate()
	cli.HelpPrinter = helpPrinter(cli.HelpPrinter)

	app.Action = action()
	app.Writer = c.outStream
	app.ErrWriter = c.errStream

	err := app.Run(args)
	if err == nil {
		return ExitCodeOK
	}

	fmt.Fprintf(clo.errStream, "%v\n", err)
	return ExitCodeError
}

func description() string {
	return `日付を, 年や月などの単位で計算します.
  日付は, 年月日時分秒のいずれかの単位ごとに加算したり減算できます. 
  単位は, 年月日時分秒それぞれを YMDhms で指定します.

  たとえば, 以下のコマンドでシステム時刻の1年3ヶ月20秒前を調べられます.

  $ dt now +1Y +3M +20s
  # ...
  
  計算元の日付を指定することもできます.

  $ dt "2018/05/12 17:30:00" +1Y +3M +20s
  2019/08/12 17:30:20

  日付のフォーマットは, 入力から自動で判断されます. 利用できるフォー
  マットについては DATE FORMATS を参照してください. 計算元の日付が
  数字のみで構成される場合は, 自動的に unix 秒と判断されます.

  $ dt -o def 1526113800 +1Y +3M +20s
  2019/08/12 17:30:20

  --input-format, -i オプションにより, unix ミリ秒も指定できます.

  $ dt -i unixm -o def 1526113800000 +1Y +3M +20s
  2019/08/12 17:30:20

  デフォルトでは出力フォーマットは入力フォーマットと同じですが,
  --output-format, -o オプションで出力フォーマットを指定できます.

  $ dt 1526113800 +1Y +3M +20s
  1565598620

  $ dt -o "Mon Jan _2 15:04:05 2006" 1526113800 +1Y +3M +20s
  Mon Aug 12 17:30:20 2019

  $ dt -o ANSIC 1526113800 +1Y +3M +20s
  Mon Aug 12 17:30:20 2019
`
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "adjust-day, a",
			Usage: "結果日付が無効なときその月末日に調整します",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "デバッグログを出力します",
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "このヘルプを表示します",
		},
		cli.StringFlag{
			Name:  "input-format, i",
			Usage: "入力フォーマットを指定します",
		},
		cli.StringFlag{
			Name:  "output-format, o",
			Usage: "出力フォーマットを指定します",
		},
		cli.BoolFlag{
			Name:  "version, v",
			Usage: "バージョンを表示します",
		},
	}
}

func appHelpTemplate() string {
	return `NAME:
  {{.Name}} - {{.Usage}}
	
USAGE:
  {{.Name}} [options] date [expr [expr ...]]
	
DESCRIPTION:
  {{.Description}}
	
OPTIONS:
  {{range .Flags}}{{.}}
  {{end}}{{if .DateFormats}}
DATE FORMATS:
  {{range .DateFormats}}{{.}}
  {{end}}{{end}}
`
}

func helpPrinter(printer func(w io.Writer, templ string, d interface{})) func(w io.Writer, templ string, d interface{}) {
	return func(w io.Writer, templ string, d interface{}) {
		app := d.(*cli.App)
		data := newHelpData(app)
		printer(w, templ, data)
	}
}

func newHelpData(app *cli.App) interface{} {
	slice := make([]string, len(formats))
	i := 0
	for k, v := range formats {
		slice[i] = k + ": " + v
		i++
	}
	sort.Strings(slice)
	return &customParameter{
		App:         app,
		DateFormats: slice,
	}
}

type customParameter struct {
	*cli.App
	DateFormats []string
}

func action() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cliContext = c

		if c.Bool("h") == true {
			loadConfig()
			cli.ShowAppHelp(c)
			return nil
		}
		if c.Bool("v") == true {
			cli.ShowVersion(c)
			return nil
		}
		if c.Bool("d") == false {
			log.SetOutput(ioutil.Discard)
		}

		log.Printf("args: %s", c.Args())

		loadConfig()

		var dt = &Dt{time: now(), format: defaultFormat}
		for i, arg := range c.Args() {
			newDt, err := processArg(i, arg, dt)
			if err != nil {
				return err
			}
			dt = newDt
		}

		output(dt)
		return nil
	}
}

func processArg(i int, arg string, dt *Dt) (*Dt, error) {
	log.Printf("arg[%d]: %s, time: %v", i, arg, dt.time)

	if i == 0 {
		return processFirst(arg)
	}
	return processRest(arg, dt)
}

func processFirst(arg string) (*Dt, error) {
	functions := []func(s string) *Dt{
		func(s string) *Dt {
			// 入力フォーマット指定
			f := cliContext.String("i")
			switch f {
			case def:
				return nil
			case unixMilliSeconds:
				match, _ := regexp.MatchString(`^\d+$`, arg)
				if match == false {
					return nil
				}
				milliSec, _ := strconv.Atoi(arg)
				return &Dt{time: time.Unix(0, int64(milliSec)*int64(time.Millisecond)), format: unixMilliSeconds}
			default:
				t, err := time.Parse(f, arg)
				if err == nil {
					return &Dt{time: t, format: f}
				}
				return nil
			}
		},
		func(s string) *Dt {
			// 現在時刻
			if arg == "now" {
				return &Dt{time: now(), format: defaultFormat}
			}
			return nil
		},
		func(s string) *Dt {
			// unix 秒として解釈
			match, _ := regexp.MatchString(`^\d+$`, arg)
			if match == true {
				unixSec, _ := strconv.Atoi(arg)
				return &Dt{time: time.Unix(int64(unixSec), 0), format: unixSeconds}
			}
			return nil
		},
		func(s string) *Dt {
			// 所定のフォーマットとして解釈
			for _, f := range formats {
				t, err := time.Parse(f, arg)
				if err == nil {
					return &Dt{time: t, format: f}
				}
			}
			return nil
		},
	}

	for _, f := range functions {
		dt := f(arg)
		if dt != nil {
			return dt, nil
		}
	}

	text := fmt.Sprintf("'%s' is invalid format.", arg)
	return nil, errors.New(text)
}

func processRest(arg string, dt *Dt) (*Dt, error) {
	match, _ := regexp.MatchString(`^[-+]?\d+[YMDhms]$`, arg)
	if match == false {
		text := fmt.Sprintf("'%s' is invalid format.", arg)
		return dt, errors.New(text)
	}

	current := dt
	month := getMonthFunc(cliContext.Bool("a"))
	functions := []func(*Dt, string) (*Dt, error){year, month, day, hour, minute, second}
	for i, f := range functions {
		newDt, err := f(current, arg)
		if err != nil || newDt.time != current.time {
			log.Printf("func: %d, time: %v -> %v", i, current.time, newDt.time)
			return newDt, err
		}
		current = newDt
	}

	return current, nil
}

func getMonthFunc(adjust bool) func(*Dt, string) (*Dt, error) {
	adjustDay := Normalize
	if adjust {
		adjustDay = AdjustToEndOfMonth
	}

	return func(dt *Dt, s string) (*Dt, error) {
		return month(dt, s, adjustDay)
	}
}

// NowInterface テスト用のインタフェース
type NowInterface interface {
	Now() time.Time
}

var nowInterface NowInterface

func now() time.Time {
	if nowInterface == nil {
		t := time.Now()
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
	}
	return nowInterface.Now()
}

func year(dt *Dt, s string) (*Dt, error) {
	return addIfMatch(s, "^[-+]?\\d+Y$", func(i int) *Dt {
		return dt.AddYear(i)
	})
}

func addIfMatch(s, pattern string, f func(i int) *Dt) (*Dt, error) {
	match, _ := regexp.MatchString(pattern, s)
	if match == false {
		return f(0), nil
	}

	duration, err := strconv.Atoi(s[0 : len(s)-1])
	if err != nil {
		return f(0), err
	}
	return f(duration), nil
}

func month(dt *Dt, s string, adjustDay AdjustDay) (*Dt, error) {
	return addIfMatch(s, "^[-+]?\\d+M$", func(i int) *Dt {
		return dt.AddMonth(i, adjustDay)
	})
}

func day(dt *Dt, s string) (*Dt, error) {
	return addIfMatch(s, "^[-+]?\\d+D$", func(i int) *Dt {
		return dt.AddDay(i)
	})
}

func hour(dt *Dt, s string) (*Dt, error) {
	return addIfMatch(s, "^[-+]?\\d+h$", func(i int) *Dt {
		return dt.AddHour(i)
	})
}

func minute(dt *Dt, s string) (*Dt, error) {
	return addIfMatch(s, "^[-+]?\\d+m$", func(i int) *Dt {
		return dt.AddMinute(i)
	})
}

func second(dt *Dt, s string) (*Dt, error) {
	return addIfMatch(s, "^[-+]?\\d+s$", func(i int) *Dt {
		return dt.AddSecond(i)
	})
}

func output(dt *Dt) {
	outputFormat := cliContext.String("o")
	switch outputFormat {
	case "":
		fmt.Fprintf(clo.outStream, "%v\n", dt)
	case "def":
		fmt.Fprintf(clo.outStream, "%s\n", &Dt{time: dt.time, format: defaultFormat})
	default:
		if v, ok := formats[outputFormat]; ok {
			fmt.Fprintf(clo.outStream, "%s\n", &Dt{time: dt.time, format: v})
		} else {
			fmt.Fprintf(clo.outStream, "%s\n", &Dt{time: dt.time, format: outputFormat})
		}
	}
}
