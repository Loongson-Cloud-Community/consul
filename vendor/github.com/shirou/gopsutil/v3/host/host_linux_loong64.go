// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_linux.go

package host

const (
	sizeofPtr      = 0x8
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x8
	sizeofLongLong = 0x8
	sizeOfUtmp     = 0x190
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

type (
	utmp struct {
		Type              int16
		Pid               int32
		Line              [32]int8
		Id                [4]int8
		User              [32]int8
		Host              [256]int8
		Exit              exit_status
		Session           int64
		Tv                timeval
		Addr_v6           [4]int32
		X__glibc_reserved [20]int8
		Pad_cgo_0         [4]byte
	}
	exit_status struct {
		Termination int16
		Exit        int16
	}
	timeval struct {
		Sec  int64
		Usec int64
	}
)
