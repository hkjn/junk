package main

type Data []int

func (d1 Data) Equal(d2 Data) bool {
	if len(d1) != len(d2) {
		return false
	}
	for i, x := range d1 {
		if x != d2[i] {
			return false
		}
	}
	return true
}

// Insertion sorts the data using in-place insertion sort.
//
// Insertion sort works by iterating through the data, growing the
// sorted sublist. In each iteration a new out-of-place item with
// index j is found, then moved left until it's in the proper place:
//   Sorted       j  Unsorted
// |-------------|x|--------|
func (d Data) Insertion() {
	for i := 1; i < len(d); i++ {
		for j := i; j > 0 && d[j-1] > d[j]; j-- {
			d[j], d[j-1] = d[j-1], d[j]
		}
	}
}

// Merge sorts the data between low and high indices using in-place merge sort.
//func (d Data) Merge() {
//	if len(d) <= 1 {
//		return
//	}
//
//	left := make(Data, middle-1)
//	right := make(Data, len(d)-middle)
//	middle := len(d) / 2
//	for i := 0; i < len(d); i++ {
//		if i < middle-1 {
//			left[i] = d[i]
//		} else {
//			right[i-middle] = d[i]
//		}
//	}
//
//	left.Merge()
//	right.Merge()
//	d = merge(left, right)
//}

//func merge(left, right Data) Data {
//	r := make(Data, len(left)+len(right))
//	i := 0
//	j := 0
//	k := 0
//	for k <
//	if left[i] > right[j] {
//		r[k] = left[i]
//		i++
//	} else {
//		r[k] = right[j]
//		j++
//	}
//	k++
//}
