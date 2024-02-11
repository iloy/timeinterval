package timeinterval

import (
	"fmt"
	"slices"
	"time"
)

type CompareResult int

const (
	Before CompareResult = iota
	Equal
	After
	Undefined
)

// TimmePoint
type TimePoint struct {
	year       int
	month      int
	day        int
	hour       int
	minute     int
	second     int
	nanosecond int

	t time.Time
}

func (tp *TimePoint) Year() int {
	return tp.year
}

func (tp *TimePoint) Month() int {
	return tp.month
}

func (tp *TimePoint) Day() int {
	return tp.day
}

func (tp *TimePoint) Hour() int {
	return tp.hour
}

func (tp *TimePoint) Minute() int {
	return tp.minute
}

func (tp *TimePoint) Second() int {
	return tp.second
}

func (tp *TimePoint) Nanosecond() int {
	return tp.nanosecond
}

func NewTimePoint(year, month, day, hour, minute, sec, nsec int) *TimePoint {
	t := time.Date(year, time.Month(month), day, hour, minute, sec, nsec, time.UTC) // no Daylight Saving Time (DST)

	_year, _month, _day := t.Date()
	_hour, _minute, _second := t.Clock()
	_nanosecond := t.Nanosecond()

	ret := &TimePoint{
		year:       _year,
		month:      int(_month), // [1, 12]
		day:        _day,
		hour:       _hour,       // [0, 23]
		minute:     _minute,     // [0, 59]
		second:     _second,     // [0, 59]
		nanosecond: _nanosecond, // [0, 999999999]

		t: t,
	}

	return ret
}

func (tp *TimePoint) Copy() *TimePoint {
	// TimePoint is immutable
	return tp
}

func (tp *TimePoint) Compare(tp2 *TimePoint) CompareResult {
	v := tp.t.Compare(tp2.t)
	switch v {
	case -1:
		return Before
	case 0:
		return Equal
	case +1:
		return After
	default:
		panic(fmt.Sprint("invalid return value from time.Time.Compare():", v))
	}
}

func (tp *TimePoint) Equal(tp2 *TimePoint) bool {
	return tp.Compare(tp2) == Equal
}

func (tp *TimePoint) Before(tp2 *TimePoint) bool {
	return tp.Compare(tp2) == Before
}

func (tp *TimePoint) After(tp2 *TimePoint) bool {
	return tp.Compare(tp2) == After
}

// sub returns the duration tp - tp2
func (tp *TimePoint) sub(tp2 *TimePoint) time.Duration {
	return tp.t.Sub(tp2.t)
}

func (tp *TimePoint) Diff(tp2 *TimePoint) time.Duration {
	v := tp.Compare(tp2)
	switch v {
	case Before:
		return tp2.sub(tp)
	case Equal:
		return time.Duration(0)
	case After:
		return tp.sub(tp2)
	default:
		panic(fmt.Sprint("invalid return value from TimePoint.Compare():", v))
	}
}

func TimePointMax(tp ...*TimePoint) *TimePoint {
	ret := tp[0]

	for i := 1; i < len(tp); i++ {
		if tp[i].After(ret) {
			ret = tp[i]
		}
	}

	return ret
}

func TimePointMin(tp ...*TimePoint) *TimePoint {
	ret := tp[0]

	for i := 1; i < len(tp); i++ {
		if tp[i].Before(ret) {
			ret = tp[i]
		}
	}

	return ret
}

// TimeInterval
type TimeInterval struct {
	start *TimePoint
	end   *TimePoint

	duration time.Duration
}

func (ti *TimeInterval) Start() *TimePoint {
	return ti.start
}

func (ti *TimeInterval) End() *TimePoint {
	return ti.end
}

func (ti *TimeInterval) Duration() time.Duration {
	return ti.duration
}

func (ti *TimeInterval) IsZeroDuration() bool {
	return ti.Duration() == time.Duration(0)
}

func NewTimeInterval(start, end *TimePoint) *TimeInterval {
	if start == nil || end == nil {
		panic("nil argument")
	}

	if end.Before(start) {
		panic(fmt.Sprintf("end is before start: start: %v, end: %v", start, end))
	}

	ret := &TimeInterval{
		start: start,
		end:   end,

		duration: start.Diff(end),
	}

	return ret
}

func (ti *TimeInterval) Copy() *TimeInterval {
	// TimeInterval is immutable
	return ti
}

func (ti *TimeInterval) Equal(ti2 *TimeInterval) bool {
	return ti.start.Equal(ti2.start) && ti.end.Equal(ti2.end)
}

func (ti *TimeInterval) Has(tp *TimePoint) bool {
	return (!tp.Before(ti.start)) && (!tp.After(ti.end))
}

func (ti *TimeInterval) Covers(ti2 *TimeInterval) bool {
	return (!ti.start.After(ti2.start)) && (!ti.end.Before(ti2.end))
}

func (ti *TimeInterval) Intersects(ti2 *TimeInterval) bool {
	// 교차하는 부분이 time.Duration(0) 보다 큰 경우에 true 반환
	// !(
	//   (ti.start.Equal(ti2.end) || ti.start.After(ti2.end))
	//   ||
	//   (ti.end.Equal(ti2.start) || ti.end.Before(ti2.start))
	// )
	return (ti.start.Before(ti2.end)) && (ti.end.After(ti2.start))
}

func (ti *TimeInterval) Merge(ti2 *TimeInterval) *TimeInterval {
	if !ti.Mergeable(ti2) {
		panic(fmt.Sprintf("not mergeable: %v, %v", ti, ti2))
	}

	ret := NewTimeInterval(
		TimePointMin(ti.start, ti2.start),
		TimePointMax(ti.end, ti2.end),
	)

	return ret
}

func (ti *TimeInterval) Mergeable(ti2 *TimeInterval) bool {
	// Intersects() 의 조건과 다른 점에 주의
	//
	// 두 개의 TimeInterval 이 바로 붙어 있는 경우에는
	// Intersects() 는 false 이지만
	// Mergeable() 은 true
	if (ti.start.After(ti2.end)) || (ti.end.Before(ti2.start)) {
		return false
	}

	return true
}

func (ti *TimeInterval) Subtract(ti2 *TimeInterval) *TimeIntervalSet {
	ret := NewTimeIntervalSet()

	if ti2.IsZeroDuration() || !ti.Intersects(ti2) {
		ret.Add(ti)
		return ret
	}

	if ti.Start().Before(ti2.Start()) {
		ret.Add(NewTimeInterval(ti.Start(), ti2.Start()))
	}

	if ti.End().After(ti2.End()) {
		ret.Add(NewTimeInterval(ti2.End(), ti.End()))
	}

	return ret
}

// TimeIntervalSet
type TimeIntervalSet struct {
	elements []*TimeInterval
}

func (tis *TimeIntervalSet) Elements() []*TimeInterval {
	return tis.elements
}

func (tis *TimeIntervalSet) Duration() time.Duration {
	ret := time.Duration(0)

	for _, ti := range tis.elements {
		ret += ti.Duration()
	}

	return ret
}

func NewTimeIntervalSet() *TimeIntervalSet {
	ret := &TimeIntervalSet{
		elements: []*TimeInterval{},
	}

	return ret
}

func (tis *TimeIntervalSet) Copy() *TimeIntervalSet {
	// TimeIntervalSet is not immutable

	ret := NewTimeIntervalSet()

	ret.Add(tis.elements...)

	return ret
}

func (tis *TimeIntervalSet) Clear() {
	tis.elements = []*TimeInterval{}
}

func (tis *TimeIntervalSet) Add(ti ...*TimeInterval) {
	tis.elements = append(tis.elements, ti...)
}

func (tis *TimeIntervalSet) Merge(tis2 ...*TimeIntervalSet) {
	for _, v := range tis2 {
		tis.elements = append(tis.elements, v.elements...)
	}
}

func (tis *TimeIntervalSet) Cleanup(removeZeroDuration bool) {
	if removeZeroDuration {
	cleanupStart0:
		l := len(tis.elements)

		for i := 0; i < l; i++ {
			elem := tis.elements[i]
			if elem.IsZeroDuration() {
				tis.elements = append(tis.elements[:i], tis.elements[i+1:]...)
				goto cleanupStart0
			}
		}
	}

	// merge two TimeIntervals if mergeable
	{
	cleanupStart1:
		l := len(tis.elements)

		for i := 0; i < l; i++ {
			for j := i + 1; j < l; j++ {
				first := tis.elements[i]
				second := tis.elements[j]

				if first.Mergeable(second) {
					tis.elements[i] = first.Merge(second)
					tis.elements = append(tis.elements[:j], tis.elements[j+1:]...)

					goto cleanupStart1
				}
			}
		}
	}

	tis.Sort()
}

func (tis *TimeIntervalSet) Sort() {
	slices.SortFunc(tis.elements, func(a, b *TimeInterval) int {
		if a.Start().Before(b.Start()) {
			return -1
		}
		if a.Start().After(b.Start()) {
			return 1
		}
		if a.End().Before(b.End()) {
			return -1
		}
		if a.End().After(b.End()) {
			return 1
		}
		return 0
	})
}

/*
// TimeIntervalMap
type TimeIntervalMap struct {
	m map[string]*TimeIntervalSet
}

func NewTimeIntervalMap(keys ...string) *TimeIntervalMap {
	ret := &TimeIntervalMap{
		m: map[string]*TimeIntervalSet{},
	}

	for _, key := range keys {
		ret.m[key] = NewTimeIntervalSet()
	}

	return ret
}

func (tim *TimeIntervalMap) Copy() *TimeIntervalMap {
	// TimeIntervalMap is not immmutable

	ret := NewTimeIntervalMap()

	for k, v := range tim.GetAll() {
		tim.m[k] = v.Copy()
	}

	return ret
}

func (tim *TimeIntervalMap) Clear() {
	tim.m = map[string]*TimeIntervalSet{}
}

func (tim *TimeIntervalMap) Add(key string, ti *TimeInterval) {
	if val, ok := tim.m[key]; ok {
		val.Add(ti)
	} else {
		panic("invalid key name: " + key)
	}
}


//func (tim *TimeIntervalMap) Delete(key string) {
//	delete(tim.m, key)
//}

func (tim *TimeIntervalMap) Get(key string) *TimeIntervalSet {
	if val, ok := tim.m[key]; ok {
		return val
	} else {
		panic("invalid key name: " + key)
	}
}

func (tim *TimeIntervalMap) GetAll() map[string]*TimeIntervalSet {
	return tim.m
}
*/
