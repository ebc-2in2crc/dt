# dt

dt is a simple and powerful command that calculate date or convert date format.

## Description

dt calculates the date in units of year, month or etc.. , and converts the format.

The date can be added or subtracted for each unit of year, month, day, hour, minute, and second.
For example, you can check date and time that 1 year 3 months 20 seconds before the system time with the following command.

```
$ dt now +1Y +3M +20s
```

You can also specify base date and time.

```
$ dt "2018/05/12 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

The date and time format are automatically determined from the input. For the available formats, see `DATE FORMATS` in `--help` option. Also, if the base date and time consists only of numbers, it is automatically determined as unix seconds.

```
$ dt -o def 1526113800 +1Y +3M +20s
2019/08/12 17:30:20
```

With the `--input-format` or `-i` option, you can use unix milliseconds.

```
$ dt -i unixm -o def 1526113800000 +1Y +3M +20s
2019/08/12 17:30:20
```

By default the output format is the same as the input format. You can specify the output format, if you use `--output-format` or `-o` option.

```
$ dt 1526113800 +1Y +3M +20s
15300622820

$ dt -o "02 Jan 06 15:04 MST" 1526113800 +1Y +3M +20s
12 Aug 19 17:00 JST
```

## Usage

```sh
$ date "+%Y/%m/%d %H:%M:%S"
2018/05/12 17:30:00

$ dt now +1Y +3M +20s
2019/08/12 17:30:20
```

### Specify base date and time

```
$ dt "2018/05/12 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

### Date and Time addition

```
$ dt "2018/05/12 17:30:00" +1Y
2019/05/12 17:30:00

$ dt "2018/05/12 17:30:00" +1M
2018/06/12 17:30:00

$ dt "2018/05/12 17:30:00" +1D
2018/05/13 17:30:00

$ dt "2018/05/12 17:30:00" +1h
2018/05/12 18:30:00

$ dt "2018/05/12 17:30:00" +1m
2018/05/12 17:31:00

$ dt "2018/05/12 17:30:00" +1s
2018/05/12 17:30:01

# '+' can omit
$ dt "2018/05/12 17:30:00" 1Y 3M 20s
2019/08/12 17:30:20
```

### Date and Time subtraction

```
$ dt "2018/05/12 17:30:00" -1Y -3M -20s
2017/02/12 17:29:40
```

### input format

#### default format

```
$ dt "2018/05/12 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

#### RFC822

```
$ dt -o def "12 May 18 17:30 MST" +1Y +3M +20s
2019/08/12 17:30:20
```

#### unix seconds

```
$ dt -o def 1526113800 +1Y +3M +20s
2019/08/12 17:30:20
```

#### automatically determined format

Following formats, and configuration.
For configuration, see `Configuration` section.

- 2006/01/02 15:04:05
- 2006-01-02 15:04:05
- 2006/01/02 15:04
- 2006-01-02 15:04
- 2006/01/02
- 2006-01-02
- Mon Jan _2 15:04:05 2006
- Mon Jan _2 15:04:05 MST 2006
- Mon Jan 02 15:04:05 -0700 2006
- 02 Jan 06 15:04 MST
- 02 Jan 06 15:04 -0700
- Monday, 02-Jan-06 15:04:05 MST
- Mon, 02 Jan 2006 15:04:05 MST
- Mon, 02 Jan 2006 15:04:05 -0700
- 2006-01-02T15:04:05Z07:00

Please see [https://golang.org/src/time/format.go](https://golang.org/src/time/format.go)

### Configuration

dt reads the configuration file (`~/.config/dt/.dt`) at startup.

#### Example

```
myformat = 02-Jan-06 15:04:05
yearonly = 2006
```

### Specify input format

#### unix seconds

```
$ dt -i unixm -o def 1526113800000 +1Y +3M +20s
2019/08/12 17:30:20
``` 

#### Specify custom format directly

```
$ dt -i "02-Jan-06 15:04:05" -o def "12-May-18 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

#### Specify custom format that defined in configuration

```
$ cat ~/.config/dt/.dt
myformat = 02-Jan-06 15:04:05

$ dt -i myformat -o def "12-May-18 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20

# You can omit the custom format that defined in the configuration file
$ dt -o def "12-May-18 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

### output format

The default output format is the same as the input format.

```
$ dt "2018-05-12 17:30:00" +1Y +3M +20s
2019-08-12 17:30:20

$ dt 1526113800 +1Y +3M +20s
15300622820
```

### Specify output format directly

```
$ dt -o def 1526113800 +1Y +3M +20s
2019/08/12 17:30:20

$ dt -o "02-Jan-06 15:04:05" 1526113800 +1Y +3M +20s
12-Aug-19 17:30:20
```

### Specify output format that defined in configuration file

```
$ cat ~/.config/dt/.dt
myformat = 02-Jan-06 15:04:05

$ dt -o myformat 1526113800 +1Y +3M +20s
12-Aug-19 17:30:20
```

### help option

```
$ dt --help
# ...

```

## Installation

### Developer

```
$ go get github.com/ebc-2in2crc/dt
$ cd $GOPATH/src/github.com/ebc-2in2crc/dt
$ make deps
$ make install
```

### User

Download from the following url.

- [https://github.com/ebc-2in2crc/dt/releases](https://github.com/ebc-2in2crc/dt/releases)

Or, you can use Homebrew (Only macOS).

```sh
$ brew tap ebc-2in2crc/dt
$ brew install dt
```

## Contribution

1. Fork this repository
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Rebase your local changes against the master branch
5. Run test suite with the go test ./... command and confirm that it passes
6. Run gofmt -s
7. Create new Pull Request

## Licence

[MIT](https://github.com/ebc-2in2crc/dt/blob/master/LICENSE)

## Author

[ebc-2in2crc](https://github.com/ebc-2in2crc)
