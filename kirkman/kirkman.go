// Rev. Thomas Penyngton Kirkman proposed this problem in 1850 in  The Lady's and Gentleman's Diary, page 48.
//
//	Fifteen young ladies in a school walk out three abreast for seven days in succession:
// it is required to arrange them daily so that no two shall walk twice abreast.
//
//	You can get started pretty easily. Name the 15 girls by the first 15 letters of the alphabet.
// Then on the first day, just choose the first thing that comes to mind, namely,
//
// ABC DEF GHI JKL MNO.
//
// A schedule for the second day is easy, too, but you have to be at least a little careful about it.
// If you start off with ADG BEH CFI, then you can't complete the second day's schedule.  But it's
// still pretty easy to do the second day.  One solution is
//  ADG BEJ CFM HKN ILO.
package main

const N = 15

type (
	Girl     uint
	Grouping []Girl
	Girls    []Grouping
)

// Valid returns true if it's a valid arrangement.
func (g Girls) Valid() bool {
	for i, group := range g {
		if !group.Valid() {
			return false
		}
	}
	return true
}

// Valid returns true if it's a valid arrangement.
func (g Grouping) Valid() bool {
	for i, girl := range g {
		if !group.Valid() {
			return false
		}
	}
	return true
}
