package iso8601

import (
	"math"
	"math/big"
	"time"
)

func decomposeRat(x *big.Rat) (w *big.Int, f *big.Rat) {
	n, d := x.Num(), x.Denom()

	w = (&big.Int{}).Div(n, d)
	f = (&big.Rat{}).Sub(x, (&big.Rat{}).SetInt(w))
	return
}

func addRatFrac(t time.Time, span time.Duration, x *big.Rat) time.Time {
	nanos := (&big.Rat{}).Mul(x, big.NewRat(span.Nanoseconds(), 1))
	w, _ := decomposeRat(nanos)

	return t.Add(time.Duration(w.Int64()))
}

func addRatYears(t time.Time, x *big.Rat) (time.Time, bool) {
	w, f := decomposeRat(x)
	if !w.IsInt64() {
		return t.AddDate(math.MaxInt32, 0, 0), false
	}

	v := w.Int64()
	if v > math.MaxInt32 {
		return t.AddDate(math.MaxInt32, 0, 0), false
	}

	t = t.AddDate(int(v), 0, 0)

	// Remainder is fractional years.
	span := t.AddDate(1, 0, 0).Sub(t)
	return addRatFrac(t, span, f), true
}

func addRatMonths(t time.Time, x *big.Rat) (time.Time, bool) {
	w, f := decomposeRat(x)
	years, months := w.QuoRem(w, big.NewInt(12), &big.Int{})

	t, cont := addRatYears(t, (&big.Rat{}).SetInt(years))
	if !cont {
		return t, cont
	}

	t = t.AddDate(0, int(months.Int64()), 0)

	// Remainder is fractional months.
	span := t.AddDate(0, 1, 0).Sub(t)
	return addRatFrac(t, span, f), true
}

func addRatDays(t time.Time, x *big.Rat) (time.Time, bool) {
	w, f := decomposeRat(x)
	if !w.IsInt64() {
		return t.AddDate(0, 0, math.MaxInt32), false
	}

	v := w.Int64()
	if v > math.MaxInt32 {
		return t.AddDate(0, 0, math.MaxInt32), false
	}

	t = t.AddDate(0, 0, int(v))

	// Remainder is fractional days.
	span := t.AddDate(0, 0, 1).Sub(t)
	return addRatFrac(t, span, f), true
}

func addRatDuration(t time.Time, x *big.Rat, base time.Duration) (time.Time, bool) {
	nanos := (&big.Rat{}).Mul(x, (&big.Rat{}).SetInt64(int64(base)))

	w, _ := decomposeRat(nanos)
	if !w.IsInt64() {
		return t.Add(time.Duration(math.MaxInt64)), false
	}

	return t.Add(time.Duration(w.Int64())), true
}

func addRatHours(t time.Time, x *big.Rat) (time.Time, bool) {
	return addRatDuration(t, x, time.Hour)
}

func addRatMinutes(t time.Time, x *big.Rat) (time.Time, bool) {
	return addRatDuration(t, x, time.Minute)
}

func addRatSeconds(t time.Time, x *big.Rat) (time.Time, bool) {
	return addRatDuration(t, x, time.Second)
}
