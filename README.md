# dt

dt は日付の計算や書式を変換するプログラムです.

## Description

dt は, 日付を年や月などの単位で計算したり, 書式を変換します.

日付は, 年月日時分秒のいずれかの単位ごとに加算したり減算できます. たとえば, 以下のコマンドでシステム時刻の1年3ヶ月20秒前を調べられます.

```
$ dt now +1Y +3M +20s
```

計算元の日付を指定することもできます. 

```
$ dt "2018/05/12 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

日付のフォーマットは, 入力から自動で判断されます. 利用できるフォーマットについては --help の DATE FORMATS を参照してください.  また, 計算元の日付が数字のみで構成される場合は, 自動的に unix 秒と判断されます.

```
$ dt -o def 1526113800 +1Y +3M +20s
2019/08/12 17:30:20
```

--input-format, -i オプションにより, unix ミリ秒も指定できます.

```
$ dt -i unixm -o def 1526113800000 +1Y +3M +20s
2019/08/12 17:30:20
```

デフォルトでは出力フォーマットは入力フォーマットと同じですが, --output-format, -o オプションで出力フォーマットを指定できます.

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

### 計算元の日付を指定

```
$ dt "2018/05/12 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

### 日付の加算

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

# '+' は省略できます
$ dt "2018/05/12 17:30:00" 1Y 3M 20s
2019/08/12 17:30:20
```

### 日付の減算

```
$ dt "2018/05/12 17:30:00" -1Y -3M -20s
2017/02/12 17:29:40
```

### 入力フォーマット

#### dt 標準

```
$ dt "2018/05/12 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

#### RFC822

```
$ dt -o def "12 May 18 17:30 MST" +1Y +3M +20s
2019/08/12 17:30:20
```

#### unix 秒

```
$ dt -o def 1526113800 +1Y +3M +20s
2019/08/12 17:30:20
```

#### 自動判断されるフォーマット

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

### 入力フォーマットを指定

#### unix ミリ秒

```
$ dt -i unixm -o def 1526113800000 +1Y +3M +20s
2019/08/12 17:30:20
``` 

#### カスタムフォーマット

```
$ dt -i "02-Jan-06 15:04:05" -o def "12-May-18 17:30:00" +1Y +3M +20s
2019/08/12 17:30:20
```

### 出力フォーマット

デフォルトでは入力フォーマットと同じ

```
$ dt "2018-05-12 17:30:00" +1Y +3M +20s
2019-08-12 17:30:20

$ dt 1526113800 +1Y +3M +20s
15300622820
```

### 出力フォーマットを指定

```
$ dt -o def 1526113800 +1Y +3M +20s
2019/08/12 17:30:20

$ dt -o "02-Jan-06 15:04:05" 1526113800 +1Y +3M +20s
12-Aug-19 17:30:20
```

### ヘルプ

```
$ dt --help
# ...

```

## Installation

### Developer

```
$ go get -u github.com/ebc-2in2crc/dt/...
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
