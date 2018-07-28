package main

import (
	"testing"
	"time"
)

func TestDt_Initialized(t *testing.T) {
	expect := time.Now()
	dt := &Dt{time: expect}

	actual := dt.Get()
	if actual != expect {
		t.Errorf("Dt.Get() = %v, want %v", actual, expect)
	}
}

func TestDt_AddYear(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddYear(1).Get()
	expect := now.AddDate(1, 0, 0)
	if actual != expect {
		t.Errorf("Dt.AddYear() = %v, want %v", actual, expect)
	}
}

func TestDt_AddMonth(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddMonth(1).Get()
	expect := now.AddDate(0, 1, 0)
	if actual != expect {
		t.Errorf("Dt.AddMonth() = %v, want %v", actual, expect)
	}
}

func TestDt_AddDay(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddDay(1).Get()
	expect := now.AddDate(0, 0, 1)
	if actual != expect {
		t.Errorf("Dt.AddDay() = %v, want %v", actual, expect)
	}
}

func TestDt_AddHour(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddHour(1).Get()
	expect := now.Add(1 * time.Hour)
	if actual != expect {
		t.Errorf("Dt.AddHour() = %v, want %v", actual, expect)
	}
}

func TestDt_AddMinute(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddMinute(1).Get()
	expect := now.Add(1 * time.Minute)
	if actual != expect {
		t.Errorf("Dt.AddMinute() = %v, want %v", actual, expect)
	}
}

func TestDt_AddSecond(t *testing.T) {
	now := time.Now()
	dt := &Dt{time: now}

	actual := dt.AddSecond(1).Get()
	expect := now.Add(1 * time.Second)
	if actual != expect {
		t.Errorf("Dt.AddSecond() = %v, want %v", actual, expect)
	}
}
