package main

import (
	"testing"
	"time"
)

func TestDt_Initialized(t *testing.T) {
	expect := time.Now()
	dt := &Dt{time: expect}

	actual := dt.get()
	if actual != expect {
		t.Errorf("Dt.get() = %v, want %v", actual, expect)
	}
}

func TestDt_AddYear(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddYear(1).get()
	expect := now.AddDate(1, 0, 0)
	if actual != expect {
		t.Errorf("Dt.AddYear() = %v, want %v", actual, expect)
	}
}

func TestDt_AddMonth(t *testing.T) {
	params := []struct {
		initial  time.Time
		addition int
		adjust   AdjustDay
		expect   time.Time
	}{
		{initial: createTime(2018, 1, 1), addition: 1, adjust: Normalize, expect: createTime(2018, 2, 1)},
		{initial: createTime(2018, 1, 31), addition: 3, adjust: Normalize, expect: createTime(2018, 5, 1)},
		{initial: createTime(2018, 1, 31), addition: 3, adjust: AdjustToEndOfMonth, expect: createTime(2018, 4, 30)},
	}

	for _, p := range params {
		dt := &Dt{time: p.initial}

		actual := dt.AddMonth(p.addition, p.adjust).get()
		expect := p.expect
		if actual != expect {
			t.Errorf("Dt.AddMonth() = %v, want %v", actual, expect)
		}
	}
}

func createTime(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func TestDt_AddDay(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddDay(1).get()
	expect := now.AddDate(0, 0, 1)
	if actual != expect {
		t.Errorf("Dt.AddDay() = %v, want %v", actual, expect)
	}
}

func TestDt_AddHour(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddHour(1).get()
	expect := now.Add(1 * time.Hour)
	if actual != expect {
		t.Errorf("Dt.AddHour() = %v, want %v", actual, expect)
	}
}

func TestDt_AddMinute(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddMinute(1).get()
	expect := now.Add(1 * time.Minute)
	if actual != expect {
		t.Errorf("Dt.AddMinute() = %v, want %v", actual, expect)
	}
}

func TestDt_AddSecond(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddSecond(1).get()
	expect := now.Add(1 * time.Second)
	if actual != expect {
		t.Errorf("Dt.AddSecond() = %v, want %v", actual, expect)
	}
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
		actualKey, actualValue := splitFormat(p.input)
		if actualKey != p.key || actualValue != p.value {
			t.Errorf("splitFormat() = %s => %s, want %s => %s", actualKey, actualValue, p.key, p.value)
		}
	}
}
