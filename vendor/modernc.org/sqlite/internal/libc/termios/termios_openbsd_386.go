// Code generated by 'ccgo termios/gen.c -crt-import-path "" -export-defines "" -export-enums "" -export-externs X -export-fields F -export-structs "" -export-typedefs "" -header -hide _OSSwapInt16,_OSSwapInt32,_OSSwapInt64 -ignore-unsupported-alignment -o termios/termios_openbsd_386.go -pkgname termios', DO NOT EDIT.

package termios

import (
	"math"
	"reflect"
	"sync/atomic"
	"unsafe"
)

var _ = math.Pi
var _ reflect.Kind
var _ atomic.Value
var _ unsafe.Pointer

const (
	ALTWERASE            = 0x00000200
	B0                   = 0
	B110                 = 110
	B115200              = 115200
	B1200                = 1200
	B134                 = 134
	B14400               = 14400
	B150                 = 150
	B1800                = 1800
	B19200               = 19200
	B200                 = 200
	B230400              = 230400
	B2400                = 2400
	B28800               = 28800
	B300                 = 300
	B38400               = 38400
	B4800                = 4800
	B50                  = 50
	B57600               = 57600
	B600                 = 600
	B7200                = 7200
	B75                  = 75
	B76800               = 76800
	B9600                = 9600
	BRKINT               = 0x00000002
	CCTS_OFLOW           = 65536
	CDISCARD             = 15
	CDSUSP               = 25
	CEOF                 = 4
	CEOT                 = 4
	CERASE               = 0177
	CFLUSH               = 15
	CHWFLOW              = 1114112
	CIGNORE              = 0x00000001
	CINTR                = 3
	CKILL                = 21
	CLNEXT               = 22
	CLOCAL               = 0x00008000
	CMIN                 = 1
	CQUIT                = 034
	CREAD                = 0x00000800
	CREPRINT             = 18
	CRPRNT               = 18
	CRTSCTS              = 0x00010000
	CRTS_IFLOW           = 65536
	CS5                  = 0x00000000
	CS6                  = 0x00000100
	CS7                  = 0x00000200
	CS8                  = 0x00000300
	CSIZE                = 0x00000300
	CSTART               = 17
	CSTOP                = 19
	CSTOPB               = 0x00000400
	CSUSP                = 26
	CTIME                = 0
	CWERASE              = 23
	ECHO                 = 0x00000008
	ECHOCTL              = 0x00000040
	ECHOE                = 0x00000002
	ECHOK                = 0x00000004
	ECHOKE               = 0x00000001
	ECHONL               = 0x00000010
	ECHOPRT              = 0x00000020
	ENDRUNDISC           = 9
	EXTA                 = 19200
	EXTB                 = 38400
	EXTPROC              = 0x00000800
	FLUSHO               = 0x00800000
	HUPCL                = 0x00004000
	ICANON               = 0x00000100
	ICRNL                = 0x00000100
	IEXTEN               = 0x00000400
	IGNBRK               = 0x00000001
	IGNCR                = 0x00000080
	IGNPAR               = 0x00000004
	IMAXBEL              = 0x00002000
	INLCR                = 0x00000040
	INPCK                = 0x00000010
	IOCPARM_MASK         = 0x1fff
	ISIG                 = 0x00000080
	ISTRIP               = 0x00000020
	IUCLC                = 0x00001000
	IXANY                = 0x00000800
	IXOFF                = 0x00000400
	IXON                 = 0x00000200
	MDMBUF               = 0x00100000
	MSTSDISC             = 8
	NCCS                 = 20
	NMEADISC             = 7
	NOFLSH               = 0x80000000
	NOKERNINFO           = 0x02000000
	OCRNL                = 0x00000010
	OLCUC                = 0x00000020
	ONLCR                = 0x00000002
	ONLRET               = 0x00000080
	ONOCR                = 0x00000040
	ONOEOT               = 0x00000008
	OPOST                = 0x00000001
	OXTABS               = 0x00000004
	PARENB               = 0x00001000
	PARMRK               = 0x00000008
	PARODD               = 0x00002000
	PENDIN               = 0x20000000
	PPPDISC              = 5
	SLIPDISC             = 4
	STRIPDISC            = 6
	TABLDISC             = 3
	TCIFLUSH             = 1
	TCIOFF               = 3
	TCIOFLUSH            = 3
	TCION                = 4
	TCOFLUSH             = 2
	TCOOFF               = 1
	TCOON                = 2
	TCSADRAIN            = 1
	TCSAFLUSH            = 2
	TCSANOW              = 0
	TCSASOFT             = 0x10
	TIOCFLAG_CLOCAL      = 0x02
	TIOCFLAG_CRTSCTS     = 0x04
	TIOCFLAG_MDMBUF      = 0x08
	TIOCFLAG_PPS         = 0x10
	TIOCFLAG_SOFTCAR     = 0x01
	TIOCM_CAR            = 0100
	TIOCM_CD             = 64
	TIOCM_CTS            = 0040
	TIOCM_DSR            = 0400
	TIOCM_DTR            = 0002
	TIOCM_LE             = 0001
	TIOCM_RI             = 128
	TIOCM_RNG            = 0200
	TIOCM_RTS            = 0004
	TIOCM_SR             = 0020
	TIOCM_ST             = 0010
	TIOCPKT_DATA         = 0x00
	TIOCPKT_DOSTOP       = 0x20
	TIOCPKT_FLUSHREAD    = 0x01
	TIOCPKT_FLUSHWRITE   = 0x02
	TIOCPKT_IOCTL        = 0x40
	TIOCPKT_NOSTOP       = 0x10
	TIOCPKT_START        = 0x08
	TIOCPKT_STOP         = 0x04
	TOSTOP               = 0x00400000
	TTYDEF_CFLAG         = 19200
	TTYDEF_IFLAG         = 11010
	TTYDEF_LFLAG         = 1483
	TTYDEF_OFLAG         = 3
	TTYDEF_SPEED         = 9600
	TTYDISC              = 0
	VDISCARD             = 15
	VDSUSP               = 11
	VEOF                 = 0
	VEOL                 = 1
	VEOL2                = 2
	VERASE               = 3
	VINTR                = 8
	VKILL                = 5
	VLNEXT               = 14
	VMIN                 = 16
	VQUIT                = 9
	VREPRINT             = 6
	VSTART               = 12
	VSTATUS              = 18
	VSTOP                = 13
	VSUSP                = 10
	VTIME                = 17
	VWERASE              = 4
	XCASE                = 0x01000000
	X_FILE_OFFSET_BITS   = 64
	X_ILP32              = 1
	X_MACHINE_CDEFS_H_   = 0
	X_MACHINE__TYPES_H_  = 0
	X_MAX_PAGE_SHIFT     = 12
	X_PID_T_DEFINED_     = 0
	X_POSIX_VDISABLE     = 255
	X_STACKALIGNBYTES    = 15
	X_SYS_CDEFS_H_       = 0
	X_SYS_IOCCOM_H_      = 0
	X_SYS_TERMIOS_H_     = 0
	X_SYS_TTYCOM_H_      = 0
	X_SYS_TTYDEFAULTS_H_ = 0
	X_SYS__TYPES_H_      = 0
	I386                 = 1
	Unix                 = 1
)

type Ptrdiff_t = int32

type Size_t = uint32

type Wchar_t = int32

type X__builtin_va_list = uintptr
type X__float128 = float64

type Tcflag_t = uint32
type Cc_t = uint8
type Speed_t = uint32

type Termios = struct {
	Fc_iflag  Tcflag_t
	Fc_oflag  Tcflag_t
	Fc_cflag  Tcflag_t
	Fc_lflag  Tcflag_t
	Fc_cc     [20]Cc_t
	Fc_ispeed int32
	Fc_ospeed int32
}

type X__int8_t = int8
type X__uint8_t = uint8
type X__int16_t = int16
type X__uint16_t = uint16
type X__int32_t = int32
type X__uint32_t = uint32
type X__int64_t = int64
type X__uint64_t = uint64

type X__int_least8_t = X__int8_t
type X__uint_least8_t = X__uint8_t
type X__int_least16_t = X__int16_t
type X__uint_least16_t = X__uint16_t
type X__int_least32_t = X__int32_t
type X__uint_least32_t = X__uint32_t
type X__int_least64_t = X__int64_t
type X__uint_least64_t = X__uint64_t

type X__int_fast8_t = X__int32_t
type X__uint_fast8_t = X__uint32_t
type X__int_fast16_t = X__int32_t
type X__uint_fast16_t = X__uint32_t
type X__int_fast32_t = X__int32_t
type X__uint_fast32_t = X__uint32_t
type X__int_fast64_t = X__int64_t
type X__uint_fast64_t = X__uint64_t

type X__intptr_t = int32
type X__uintptr_t = uint32

type X__intmax_t = X__int64_t
type X__uintmax_t = X__uint64_t

type X__register_t = int32

type X__vaddr_t = uint32
type X__paddr_t = uint32
type X__vsize_t = uint32
type X__psize_t = uint32

type X__double_t = float64
type X__float_t = float64
type X__ptrdiff_t = int32
type X__size_t = uint32
type X__ssize_t = int32
type X__va_list = X__builtin_va_list

type X__wchar_t = int32
type X__wint_t = int32
type X__rune_t = int32
type X__wctrans_t = uintptr
type X__wctype_t = uintptr

type X__blkcnt_t = X__int64_t
type X__blksize_t = X__int32_t
type X__clock_t = X__int64_t
type X__clockid_t = X__int32_t
type X__cpuid_t = uint32
type X__dev_t = X__int32_t
type X__fixpt_t = X__uint32_t
type X__fsblkcnt_t = X__uint64_t
type X__fsfilcnt_t = X__uint64_t
type X__gid_t = X__uint32_t
type X__id_t = X__uint32_t
type X__in_addr_t = X__uint32_t
type X__in_port_t = X__uint16_t
type X__ino_t = X__uint64_t
type X__key_t = int32
type X__mode_t = X__uint32_t
type X__nlink_t = X__uint32_t
type X__off_t = X__int64_t
type X__pid_t = X__int32_t
type X__rlim_t = X__uint64_t
type X__sa_family_t = X__uint8_t
type X__segsz_t = X__int32_t
type X__socklen_t = X__uint32_t
type X__suseconds_t = int32
type X__time_t = X__int64_t
type X__timer_t = X__int32_t
type X__uid_t = X__uint32_t
type X__useconds_t = X__uint32_t

type X__mbstate_t = struct {
	F__ccgo_pad1 [0]uint32
	F__mbstate8  [128]int8
}

type Pid_t = X__pid_t

type Winsize = struct {
	Fws_row    uint16
	Fws_col    uint16
	Fws_xpixel uint16
	Fws_ypixel uint16
}

type Tstamps = struct {
	Fts_set int32
	Fts_clr int32
}

var _ int8