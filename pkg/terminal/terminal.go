package terminal

import (
	"fmt"
	"io"
	"strings"
)

const (
	ResetCode     = "\x1b[0m"
	NormalCode    = "\x1b[0m"
	BoldCode      = "\x1b[1m"
	DimCode       = "\x1b[2m"
	UnderlineCode = "\x1b[4m"
	BlackCode     = "\x1b[30m"
	RedCode       = "\x1b[31m"
	GreenCode     = "\x1b[32m"
	YellowCode    = "\x1b[33m"
	BlueCode      = "\x1b[34m"
	MagentaCode   = "\x1b[35m"
	CyanCode      = "\x1b[36m"
	WhiteCode     = "\x1b[37m"
	BgBlackCode   = "\x1b[40m"
	BgRedCode     = "\x1b[41m"
	BgGreenCode   = "\x1b[42m"
	BgYellowCode  = "\x1b[43m"
	BgBlueCode    = "\x1b[44m"
	BgMagentaCode = "\x1b[45m"
	BgCyanCode    = "\x1b[46m"
	BgWhiteCode   = "\x1b[47m"
)

var (
	None = New()

	Normal        = New(NormalCode).Sprint
	Bold          = New(BoldCode).Sprint
	Dim           = New(DimCode).Sprint
	Underline     = New(UnderlineCode).Sprint
	Black         = New(BlackCode).Sprint
	Red           = New(RedCode).Sprint
	Green         = New(GreenCode).Sprint
	Yellow        = New(YellowCode).Sprint
	Blue          = New(BlueCode).Sprint
	Magenta       = New(MagentaCode).Sprint
	Cyan          = New(CyanCode).Sprint
	White         = New(WhiteCode).Sprint
	BgBlack       = New(BgBlackCode).Sprint
	BgRed         = New(BgRedCode).Sprint
	BgGreen       = New(BgGreenCode).Sprint
	BgYellow      = New(BgYellowCode).Sprint
	BgBlue        = New(BgBlueCode).Sprint
	BgMagenta     = New(BgMagentaCode).Sprint
	BgCyan        = New(BgCyanCode).Sprint
	BgWhite       = New(BgWhiteCode).Sprint
	DimUnderline  = New(DimCode, UnderlineCode).Sprint
	DimBlack      = New(DimCode, BlackCode).Sprint
	DimRed        = New(DimCode, RedCode).Sprint
	DimGreen      = New(DimCode, GreenCode).Sprint
	DimYellow     = New(DimCode, YellowCode).Sprint
	DimBlue       = New(DimCode, BlueCode).Sprint
	DimMagenta    = New(DimCode, MagentaCode).Sprint
	DimCyan       = New(DimCode, CyanCode).Sprint
	DimWhite      = New(DimCode, WhiteCode).Sprint
	DimBgBlack    = New(DimCode, BgBlackCode).Sprint
	DimBgRed      = New(DimCode, BgRedCode).Sprint
	DimBgGreen    = New(DimCode, BgGreenCode).Sprint
	DimBgYellow   = New(DimCode, BgYellowCode).Sprint
	DimBgBlue     = New(DimCode, BgBlueCode).Sprint
	DimBgMagenta  = New(DimCode, BgMagentaCode).Sprint
	DimBgCyan     = New(DimCode, BgCyanCode).Sprint
	DimBgWhite    = New(DimCode, BgWhiteCode).Sprint
	BoldUnderline = New(BoldCode, UnderlineCode).Sprint
	BoldBlack     = New(BoldCode, BlackCode).Sprint
	BoldRed       = New(BoldCode, RedCode).Sprint
	BoldGreen     = New(BoldCode, GreenCode).Sprint
	BoldYellow    = New(BoldCode, YellowCode).Sprint
	BoldBlue      = New(BoldCode, BlueCode).Sprint
	BoldMagenta   = New(BoldCode, MagentaCode).Sprint
	BoldCyan      = New(BoldCode, CyanCode).Sprint
	BoldWhite     = New(BoldCode, WhiteCode).Sprint
	BoldBgBlack   = New(BoldCode, BgBlackCode).Sprint
	BoldBgRed     = New(BoldCode, BgRedCode).Sprint
	BoldBgGreen   = New(BoldCode, BgGreenCode).Sprint
	BoldBgYellow  = New(BoldCode, BgYellowCode).Sprint
	BoldBgBlue    = New(BoldCode, BgBlueCode).Sprint
	BoldBgMagenta = New(BoldCode, BgMagentaCode).Sprint
	BoldBgCyan    = New(BoldCode, BgCyanCode).Sprint
	BoldBgWhite   = New(BoldCode, BgWhiteCode).Sprint

	Normalf        = New(NormalCode).Sprintf
	Boldf          = New(BoldCode).Sprintf
	Dimf           = New(DimCode).Sprintf
	Underlinef     = New(UnderlineCode).Sprintf
	Blackf         = New(BlackCode).Sprintf
	Redf           = New(RedCode).Sprintf
	Greenf         = New(GreenCode).Sprintf
	Yellowf        = New(YellowCode).Sprintf
	Bluef          = New(BlueCode).Sprintf
	Magentaf       = New(MagentaCode).Sprintf
	Cyanf          = New(CyanCode).Sprintf
	Whitef         = New(WhiteCode).Sprintf
	BgBlackf       = New(BgBlackCode).Sprintf
	BgRedf         = New(BgRedCode).Sprintf
	BgGreenf       = New(BgGreenCode).Sprintf
	BgYellowf      = New(BgYellowCode).Sprintf
	BgBluef        = New(BgBlueCode).Sprintf
	BgMagentaf     = New(BgMagentaCode).Sprintf
	BgCyanf        = New(BgCyanCode).Sprintf
	BgWhitef       = New(BgWhiteCode).Sprintf
	DimUnderlinef  = New(DimCode, UnderlineCode).Sprintf
	DimBlackf      = New(DimCode, BlackCode).Sprintf
	DimRedf        = New(DimCode, RedCode).Sprintf
	DimGreenf      = New(DimCode, GreenCode).Sprintf
	DimYellowf     = New(DimCode, YellowCode).Sprintf
	DimBluef       = New(DimCode, BlueCode).Sprintf
	DimMagentaf    = New(DimCode, MagentaCode).Sprintf
	DimCyanf       = New(DimCode, CyanCode).Sprintf
	DimWhitef      = New(DimCode, WhiteCode).Sprintf
	DimBgBlackf    = New(DimCode, BgBlackCode).Sprintf
	DimBgRedf      = New(DimCode, BgRedCode).Sprintf
	DimBgGreenf    = New(DimCode, BgGreenCode).Sprintf
	DimBgYellowf   = New(DimCode, BgYellowCode).Sprintf
	DimBgBluef     = New(DimCode, BgBlueCode).Sprintf
	DimBgMagentaf  = New(DimCode, BgMagentaCode).Sprintf
	DimBgCyanf     = New(DimCode, BgCyanCode).Sprintf
	DimBgWhitef    = New(DimCode, BgWhiteCode).Sprintf
	BoldUnderlinef = New(BoldCode, UnderlineCode).Sprintf
	BoldBlackf     = New(BoldCode, BlackCode).Sprintf
	BoldRedf       = New(BoldCode, RedCode).Sprintf
	BoldGreenf     = New(BoldCode, GreenCode).Sprintf
	BoldYellowf    = New(BoldCode, YellowCode).Sprintf
	BoldBluef      = New(BoldCode, BlueCode).Sprintf
	BoldMagentaf   = New(BoldCode, MagentaCode).Sprintf
	BoldCyanf      = New(BoldCode, CyanCode).Sprintf
	BoldWhitef     = New(BoldCode, WhiteCode).Sprintf
	BoldBgBlackf   = New(BoldCode, BgBlackCode).Sprintf
	BoldBgRedf     = New(BoldCode, BgRedCode).Sprintf
	BoldBgGreenf   = New(BoldCode, BgGreenCode).Sprintf
	BoldBgYellowf  = New(BoldCode, BgYellowCode).Sprintf
	BoldBgBluef    = New(BoldCode, BgBlueCode).Sprintf
	BoldBgMagentaf = New(BoldCode, BgMagentaCode).Sprintf
	BoldBgCyanf    = New(BoldCode, BgCyanCode).Sprintf
	BoldBgWhitef   = New(BoldCode, BgWhiteCode).Sprintf
)

type Formatter struct {
	code string
}

func New(codes ...string) Formatter {
	code := strings.Join(codes, "")
	formatter := Formatter{code}
	return formatter
}

func (f Formatter) Sprintf(format string, args ...interface{}) string {
	var str string
	str = fmt.Sprintf(format, args...)
	if f.code == "" {
		return str
	}
	str = f.code + strings.ReplaceAll(str, ResetCode, ResetCode+f.code) + ResetCode
	return str
}

func (f Formatter) Sprint(args ...interface{}) string {
	var str string
	str = fmt.Sprint(args...)
	if f.code == "" {
		return str
	}
	str = f.code + strings.ReplaceAll(str, ResetCode, ResetCode+f.code) + ResetCode
	return str
}

func Revert(w io.Writer) (n int, err error) {
	return fmt.Fprint(w, "\x1b[0m")
}
