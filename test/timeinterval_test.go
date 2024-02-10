package timeinterval_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/iloy/timeinterval"
)

const (
	year  = 2024
	month = 2
	day   = 5
)

func TestTimePointNew(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 60, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 20, 0, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 0, 60, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)

	assert.Equal(t, t1.Equal(t2), true)
	assert.Equal(t, t3.Equal(t4), true)

	assert.Equal(t, t1.Year(), year)
	assert.Equal(t, t1.Month(), month)
	assert.Equal(t, t1.Day(), day)
	assert.Equal(t, t1.Hour(), 20)
	assert.Equal(t, t1.Minute(), 0)
	assert.Equal(t, t1.Second(), 0)
	assert.Equal(t, t1.Nanosecond(), 0)

	assert.Equal(t, t3.Minute(), 1)
	assert.Equal(t, t3.Second(), 0)
}

func TestTimePointDiff(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 20, 0, 0, 0)

	assert.Equal(t, t1.Diff(t2), time.Duration(0))
	assert.Equal(t, t1.Diff(t3), time.Minute)
	assert.Equal(t, t1.Diff(t4), time.Minute*2)
	assert.Equal(t, t1.Diff(t5), time.Hour*1)
	assert.Equal(t, t1.Diff(t5), time.Minute*60)
}

func TestTimePointMinMax(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)

	assert.Equal(t, timeinterval.TimePointMin(t1, t1).Equal(t1), true)
	assert.Equal(t, timeinterval.TimePointMin(t1, t2).Equal(t1), true)
	assert.Equal(t, timeinterval.TimePointMax(t1, t2).Equal(t2), true)
	assert.Equal(t, timeinterval.TimePointMax(t2, t2).Equal(t2), true)
}

func TestTimeIntervalNew(t *testing.T) {
	t0 := timeinterval.NewTimePoint(year, month, day, 18, 59, 0, 0)
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)

	assert.Panics(t, func() { _ = timeinterval.NewTimeInterval(t1, t0) })
	t11 := timeinterval.NewTimeInterval(t1, t1)
	t12 := timeinterval.NewTimeInterval(t1, t2)
	t13 := timeinterval.NewTimeInterval(t1, t3)

	assert.Equal(t, t11.Duration(), time.Duration(0))
	assert.Equal(t, t11.Duration(), time.Duration(0))
	assert.Equal(t, t11.IsZeroDuration(), true)
	assert.Equal(t, t12.Duration(), time.Duration(0))
	assert.Equal(t, t12.Duration(), time.Duration(0))
	assert.Equal(t, t12.IsZeroDuration(), true)
	assert.Equal(t, t13.Duration(), time.Minute)
	assert.Equal(t, t13.Duration(), time.Minute)
	assert.Equal(t, t13.IsZeroDuration(), false)
}

func TestTimeIntervalEqual(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)

	ti12 := timeinterval.NewTimeInterval(t1, t2)
	ti13 := timeinterval.NewTimeInterval(t1, t3)
	ti23 := timeinterval.NewTimeInterval(t2, t3)

	assert.Equal(t, ti12.Equal(ti13), false)
	assert.Equal(t, ti12.Equal(ti23), false)
	assert.Equal(t, ti13.Equal(ti23), true)
}

func TestTimeIntervalHas(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 3, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 19, 4, 0, 0)

	ti24 := timeinterval.NewTimeInterval(t2, t4)

	assert.Equal(t, ti24.Has(t1), false)
	assert.Equal(t, ti24.Has(t2), true)
	assert.Equal(t, ti24.Has(t3), true)
	assert.Equal(t, ti24.Has(t4), true)
	assert.Equal(t, ti24.Has(t5), false)
}

func TestTimeIntervalCovers(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 3, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 19, 4, 0, 0)

	ti11 := timeinterval.NewTimeInterval(t1, t1)
	ti12 := timeinterval.NewTimeInterval(t1, t2)
	ti13 := timeinterval.NewTimeInterval(t1, t3)
	ti14 := timeinterval.NewTimeInterval(t1, t4)
	ti15 := timeinterval.NewTimeInterval(t1, t5)
	ti23 := timeinterval.NewTimeInterval(t2, t3)
	ti24 := timeinterval.NewTimeInterval(t2, t4)
	ti25 := timeinterval.NewTimeInterval(t2, t5)
	ti34 := timeinterval.NewTimeInterval(t3, t4)
	ti35 := timeinterval.NewTimeInterval(t3, t5)
	ti44 := timeinterval.NewTimeInterval(t4, t4)
	ti45 := timeinterval.NewTimeInterval(t4, t5)
	ti55 := timeinterval.NewTimeInterval(t5, t5)

	assert.Equal(t, ti11.Covers(ti24), false)
	assert.Equal(t, ti12.Covers(ti24), false)
	assert.Equal(t, ti13.Covers(ti24), false)
	assert.Equal(t, ti14.Covers(ti24), true)
	assert.Equal(t, ti15.Covers(ti24), true)
	assert.Equal(t, ti23.Covers(ti24), false)
	assert.Equal(t, ti24.Covers(ti24), true)
	assert.Equal(t, ti25.Covers(ti24), true)
	assert.Equal(t, ti34.Covers(ti24), false)
	assert.Equal(t, ti35.Covers(ti24), false)
	assert.Equal(t, ti44.Covers(ti24), false)
	assert.Equal(t, ti45.Covers(ti24), false)
	assert.Equal(t, ti55.Covers(ti24), false)
}

func TestTimeIntervalIntersects(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 3, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 19, 4, 0, 0)

	ti11 := timeinterval.NewTimeInterval(t1, t1)
	ti12 := timeinterval.NewTimeInterval(t1, t2)
	ti13 := timeinterval.NewTimeInterval(t1, t3)
	ti14 := timeinterval.NewTimeInterval(t1, t4)
	ti15 := timeinterval.NewTimeInterval(t1, t5)
	ti23 := timeinterval.NewTimeInterval(t2, t3)
	ti24 := timeinterval.NewTimeInterval(t2, t4)
	ti25 := timeinterval.NewTimeInterval(t2, t5)
	ti34 := timeinterval.NewTimeInterval(t3, t4)
	ti35 := timeinterval.NewTimeInterval(t3, t5)
	ti44 := timeinterval.NewTimeInterval(t4, t4)
	ti45 := timeinterval.NewTimeInterval(t4, t5)
	ti55 := timeinterval.NewTimeInterval(t5, t5)

	assert.Equal(t, ti11.Intersects(ti12), false)
	assert.Equal(t, ti12.Intersects(ti23), false)
	assert.Equal(t, ti13.Intersects(ti24), true)
	assert.Equal(t, ti14.Intersects(ti24), true)
	assert.Equal(t, ti15.Intersects(ti24), true)
	assert.Equal(t, ti25.Intersects(ti24), true)
	assert.Equal(t, ti34.Intersects(ti24), true)
	assert.Equal(t, ti35.Intersects(ti24), true)
	assert.Equal(t, ti44.Intersects(ti24), false)
	assert.Equal(t, ti45.Intersects(ti24), false)
	assert.Equal(t, ti55.Intersects(ti24), false)
}

func TestTimeIntervalMerge(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 3, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 19, 4, 0, 0)

	ti11 := timeinterval.NewTimeInterval(t1, t1)
	ti12 := timeinterval.NewTimeInterval(t1, t2)
	ti13 := timeinterval.NewTimeInterval(t1, t3)
	ti14 := timeinterval.NewTimeInterval(t1, t4)
	ti15 := timeinterval.NewTimeInterval(t1, t5)
	ti23 := timeinterval.NewTimeInterval(t2, t3)
	ti24 := timeinterval.NewTimeInterval(t2, t4)
	/*
		ti25 := timeinterval.NewTimeInterval(t2, t5)
	*/
	ti34 := timeinterval.NewTimeInterval(t3, t4)
	/*
		ti35 := timeinterval.NewTimeInterval(t3, t5)
		ti44 := timeinterval.NewTimeInterval(t4, t4)
		ti45 := timeinterval.NewTimeInterval(t4, t5)
		ti55 := timeinterval.NewTimeInterval(t5, t5)
	*/

	assert.Equal(t, ti11.Merge(ti11).Equal(ti11), true)
	assert.Equal(t, ti11.Merge(ti15).Equal(ti15), true)
	assert.Equal(t, ti15.Merge(ti11).Equal(ti15), true)
	assert.Equal(t, ti12.Merge(ti23).Equal(ti13), true)
	assert.Equal(t, ti23.Merge(ti12).Equal(ti13), true)
	assert.Equal(t, ti13.Merge(ti12).Equal(ti13), true)
	assert.Equal(t, ti13.Merge(ti23).Equal(ti13), true)
	assert.Equal(t, ti23.Merge(ti13).Equal(ti13), true)
	assert.Equal(t, ti13.Merge(ti24).Equal(ti14), true)
	assert.Equal(t, ti24.Merge(ti13).Equal(ti14), true)
	assert.Panics(t, func() { ti12.Merge(ti34) })
}

func TestTimeIntervalSubtract(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 3, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 19, 4, 0, 0)
	t6 := timeinterval.NewTimePoint(year, month, day, 19, 5, 0, 0)
	t7 := timeinterval.NewTimePoint(year, month, day, 19, 6, 0, 0)
	t8 := timeinterval.NewTimePoint(year, month, day, 19, 7, 0, 0)

	ti11 := timeinterval.NewTimeInterval(t1, t1)
	ti12 := timeinterval.NewTimeInterval(t1, t2)
	ti13 := timeinterval.NewTimeInterval(t1, t3)
	ti24 := timeinterval.NewTimeInterval(t2, t4)
	ti26 := timeinterval.NewTimeInterval(t2, t6)
	ti27 := timeinterval.NewTimeInterval(t2, t7)
	ti33 := timeinterval.NewTimeInterval(t3, t3)
	ti34 := timeinterval.NewTimeInterval(t3, t4)
	ti35 := timeinterval.NewTimeInterval(t3, t5)

	ti36 := timeinterval.NewTimeInterval(t3, t6)

	ti37 := timeinterval.NewTimeInterval(t3, t7)
	ti44 := timeinterval.NewTimeInterval(t4, t4)
	ti45 := timeinterval.NewTimeInterval(t4, t5)
	ti46 := timeinterval.NewTimeInterval(t4, t6)
	ti47 := timeinterval.NewTimeInterval(t4, t7)
	ti56 := timeinterval.NewTimeInterval(t5, t6)
	ti66 := timeinterval.NewTimeInterval(t6, t6)
	ti67 := timeinterval.NewTimeInterval(t6, t7)
	ti78 := timeinterval.NewTimeInterval(t7, t8)
	ti88 := timeinterval.NewTimeInterval(t8, t8)

	{
		tis := ti36.Subtract(ti11)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti12)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti13)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti24)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti46), true)
	}
	{
		tis := ti36.Subtract(ti26)
		assert.Equal(t, len(tis.Elements()), 0)
		assert.Equal(t, tis.Duration(), time.Duration(0))
	}
	{
		tis := ti36.Subtract(ti27)
		assert.Equal(t, len(tis.Elements()), 0)
		assert.Equal(t, tis.Duration(), time.Duration(0))
	}

	{
		tis := ti36.Subtract(ti33)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti35)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti56), true)
	}
	{
		tis := ti36.Subtract(ti36)
		assert.Equal(t, len(tis.Elements()), 0)
		assert.Equal(t, tis.Duration(), time.Duration(0))
	}
	{
		tis := ti36.Subtract(ti37)
		assert.Equal(t, len(tis.Elements()), 0)
		assert.Equal(t, tis.Duration(), time.Duration(0))
	}

	{
		// TODO
		// is this right?
		tis := ti36.Subtract(ti44)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti45)
		assert.Equal(t, len(tis.Elements()), 2)
		assert.Equal(t, tis.Elements()[0].Equal(ti34), true)
		assert.Equal(t, tis.Elements()[1].Equal(ti56), true)
	}

	{
		tis := ti36.Subtract(ti46)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti34), true)
	}
	{
		tis := ti36.Subtract(ti47)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti34), true)
	}
	{
		tis := ti36.Subtract(ti66)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti67)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti78)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
	{
		tis := ti36.Subtract(ti88)
		assert.Equal(t, len(tis.Elements()), 1)
		assert.Equal(t, tis.Elements()[0].Equal(ti36), true)
	}
}

func TestTimeIntervalSetNewAddCleanup(t *testing.T) {
	t1 := timeinterval.NewTimePoint(year, month, day, 19, 0, 0, 0)
	t2 := timeinterval.NewTimePoint(year, month, day, 19, 1, 0, 0)
	t3 := timeinterval.NewTimePoint(year, month, day, 19, 2, 0, 0)
	t4 := timeinterval.NewTimePoint(year, month, day, 19, 3, 0, 0)
	t5 := timeinterval.NewTimePoint(year, month, day, 19, 4, 0, 0)

	ti11 := timeinterval.NewTimeInterval(t1, t1)
	ti12 := timeinterval.NewTimeInterval(t1, t2)
	ti13 := timeinterval.NewTimeInterval(t1, t3)
	ti14 := timeinterval.NewTimeInterval(t1, t4)
	ti15 := timeinterval.NewTimeInterval(t1, t5)
	ti23 := timeinterval.NewTimeInterval(t2, t3)
	ti24 := timeinterval.NewTimeInterval(t2, t4)

	tis := timeinterval.NewTimeIntervalSet()

	tis.Add(ti11, ti12, ti13, ti14, ti15, ti23, ti24)
	tis.Cleanup(true)

	assert.Equal(t, len(tis.Elements()), 1)
	assert.Equal(t, tis.Elements()[0].Equal(ti15), true)
}
