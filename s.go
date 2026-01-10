package main

// import "fmt"
//
//	func main() {
//		asd := []byte{65, 65, 65, 65, 65, 65, 65, 65,
//			65, 65, 10, 65, 65, 65, 83, 65,
//			65, 65, 65, 65, 65, 65, 65, 65,
//			65, 10, 65, 65, 10, 65, 10, 65,
//			65, 65, 65, 10, 65, 65, 65, 10,
//			65, 65, 65, 10, 65, 10, 10, 10}
//		fmt.Println(splitSliceBy11(asd))
//	}
//
// // split byte slice into sub slices by 10
func splitSliceBy10(s []byte) [][]byte {
	var bslices [][]byte
	var bslice []byte
	for i := 0; i < len(s); i++ {
		if s[i] == 10 && len(bslice) > 0 {
			bslices = append(bslices, bslice)
			bslice = []byte{}
			continue
		}
		if s[i] != 10 {
			bslice = append(bslice, s[i])
		}
	}
	if len(bslice) > 0 {
		bslices = append(bslices, bslice)
	}
	return bslices
}
