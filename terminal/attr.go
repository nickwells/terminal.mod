package terminal

import (
	"fmt"
	"image/color"
	"io"
	"strings"
)

const (
	csi    = "\033["
	sgrSep = ";"
	sgrEnd = "m"

	sgrColourIdxIntro = "5"
	sgrColourRGBIntro = "2"

	sgrReset           = "0"
	sgrBold            = "1"
	sgrFaint           = "2"
	sgrItalic          = "3"
	sgrUnderline       = "4"
	sgrBlinkSlow       = "5"
	sgrBlinkFast       = "6"
	sgrInvert          = "7"
	sgrCrossedOut      = "9"
	sgrDoubleUnderline = "21"
	sgrNormalIntensity = "22"
	sgrNoUnderline     = "24"
	sgrNoBlinking      = "25"
	sgrFGColour        = "38"
	sgrBGColour        = "48"
	sgrOverlined       = "53"

	sgrColourBlack = 1
	sgrColourWhite = 15
	sgrColourRed   = 9
	sgrColourGreen = 10
	sgrColourBlue  = 12
)

type underlineStyle int

const (
	ulNone underlineStyle = iota
	ulSingle
	ulDouble
)

var underlineMap = map[underlineStyle]string{
	ulNone:   sgrNoUnderline,
	ulSingle: sgrUnderline,
	ulDouble: sgrDoubleUnderline,
}

type blinkStyle int

const (
	blkNone blinkStyle = iota
	blkSlow
	blkFast
)

var blinkMap = map[blinkStyle]string{
	blkNone: sgrNoBlinking,
	blkSlow: sgrBlinkSlow,
	blkFast: sgrBlinkFast,
}

type intensityStyle int

const (
	ityStd intensityStyle = iota
	ityBold
	ityFaint
)

var intensityMap = map[intensityStyle]string{
	ityStd:   sgrNormalIntensity,
	ityBold:  sgrBold,
	ityFaint: sgrFaint,
}

// attrColour represents a colour as used in a terminal
type attrColour struct {
	useColourIdx bool
	useRGB       bool

	colourIdx uint8

	red   uint8
	green uint8
	blue  uint8
}

// Attr represents the attributes that text should have
type Attr struct {
	overlined  bool
	crossedOut bool
	underline  underlineStyle

	inverted bool

	italic    bool
	blink     blinkStyle
	intensity intensityStyle

	fgColour attrColour
	bgColour attrColour
}

// mkAttrStr generates the escape codes to set the text attributes on a
// terminal
func mkAttrStr(attrs ...string) string {
	if len(attrs) == 0 {
		return ""
	}

	return csi + strings.Join(attrs, sgrSep) + sgrEnd
}

// Reset sends the reset code. This should undo any previous attribute setting
func Reset(w io.Writer) {
	fmt.Fprint(w, mkAttrStr(sgrReset))
}

// Start sends the codes to set the text attributes as given.
func (a Attr) Start(w io.Writer) {
	attrs := []string{}

	if a.overlined {
		attrs = append(attrs, sgrOverlined)
	}

	if a.crossedOut {
		attrs = append(attrs, sgrCrossedOut)
	}

	attrs = append(attrs, underlineMap[a.underline])

	if a.inverted {
		attrs = append(attrs, sgrInvert)
	}

	if a.italic {
		attrs = append(attrs, sgrItalic)
	}

	attrs = append(attrs, blinkMap[a.blink])
	attrs = append(attrs, intensityMap[a.intensity])

	if a.fgColour.useColourIdx {
		attrs = append(attrs, sgrFGColour, sgrColourIdxIntro)
		attrs = append(attrs, fmt.Sprintf("%d", a.fgColour.colourIdx))
	}

	if a.fgColour.useRGB {
		attrs = append(attrs, sgrFGColour, sgrColourRGBIntro)
		attrs = append(attrs, fmt.Sprintf("%d", a.fgColour.red))
		attrs = append(attrs, fmt.Sprintf("%d", a.fgColour.green))
		attrs = append(attrs, fmt.Sprintf("%d", a.fgColour.blue))
	}

	if a.bgColour.useColourIdx {
		attrs = append(attrs, sgrBGColour, sgrColourIdxIntro)
		attrs = append(attrs, fmt.Sprintf("%d", a.bgColour.colourIdx))
	}

	if a.bgColour.useRGB {
		attrs = append(attrs, sgrBGColour, sgrColourRGBIntro)
		attrs = append(attrs, fmt.Sprintf("%d", a.bgColour.red))
		attrs = append(attrs, fmt.Sprintf("%d", a.bgColour.green))
		attrs = append(attrs, fmt.Sprintf("%d", a.bgColour.blue))
	}

	fmt.Fprint(w, mkAttrStr(attrs...))
}

// End sents the codes to reset the text attributes to their default
func (_ Attr) End(w io.Writer) {
	Reset(w)
}

// AttrFunc is the type of a function used to set a text attrribute
type AttrFunc func(*Attr)

// MkAttr generates and returns a text attributes structure
func MkAttr(af ...AttrFunc) Attr {
	var a Attr

	for _, af := range af {
		af(&a)
	}

	return a
}

// AttrOverlined sets the Overlined attribute.
func AttrOverlined(a *Attr) {
	a.overlined = true
}

// AttrCrossedOut sets the CrossedOut attribute.
func AttrCrossedOut(a *Attr) {
	a.crossedOut = true
}

// AttrInverted sets the Inverted attribute.
func AttrInverted(a *Attr) {
	a.inverted = true
}

// AttrItalic sets the Italic attribute.
func AttrItalic(a *Attr) {
	a.italic = true
}

// AttrUnderline sets the Underline attribute. See also [AttrDoubleUnderline]
// and [AttrNoUnderline]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
func AttrUnderline(a *Attr) {
	a.underline = ulSingle
}

// AttrDoubleUnderline sets the DoubleUnderline attribute. See also
// [AttrUnderline] and [AttrNoUnderline]. Only one of these should be passed
// when creating a [Attr]. If more than one is given only the last takes
// effect.
func AttrDoubleUnderline(a *Attr) {
	a.underline = ulDouble
}

// AttrNoUnderline turns off Underlining. See also [AttrUnderline] and
// [AttrDoubleUnderline]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
func AttrNoUnderline(a *Attr) {
	a.underline = ulNone
}

// AttrBlinkFast sets the Blink attribute to fast. See also [AttrBlinkSlow]
// and [AttrNoBlink]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
func AttrBlinkFast(a *Attr) {
	a.blink = blkFast
}

// AttrBlinkSlow sets the Blink attribute to slow. See also [AttrBlinkFast]
// and [AttrNoBlink]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
func AttrBlinkSlow(a *Attr) {
	a.blink = blkSlow
}

// AttrNoBlink turns off blinking. See also [AttrBlinkFast] and
// [AttrBlinkSlow]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
func AttrNoBlink(a *Attr) {
	a.blink = blkNone
}

// AttrIntensityNormal sets the Intensity attribute to its normal value. See
// also [AttrBold] and [AttrFaint]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrIntensityNormal(a *Attr) {
	a.intensity = ityStd
}

// AttrBold sets the Intensity attribute to bold. See also
// [AttrIntensityNormal] and [AttrFaint]. Only one of these should be passed
// when creating a [Attr]. If more than one is given only the last takes
// effect.
func AttrBold(a *Attr) {
	a.intensity = ityStd
}

// AttrFaint sets the Intensity attribute to faint. See also [AttrBold] and
// [AttrIntensityNormal]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
func AttrFaint(a *Attr) {
	a.intensity = ityStd
}

// AttrFGShort returns an AttrFunc that will set the foreground colour to the
// colour represented by the supplied short code.
//
// See also [AttrFGColour], [AttrFGWBlack], [AttrFGWWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGShort(c uint8) AttrFunc {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    c,
	}
	return func(a *Attr) {
		a.fgColour = aC
	}
}

// AttrFGColour returns an AttrFunc that will set the foreground colour to the
// colour represented by the supplied colour.
//
// See also [AttrFGShort], [AttrFGWBlack], [AttrFGWWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGColour(c color.RGBA) AttrFunc {
	aC := attrColour{
		useRGB: true,
		red:    c.R,
		green:  c.G,
		blue:   c.B,
	}
	return func(a *Attr) {
		a.fgColour = aC
	}
}

// AttrFGBlack sets the foreground colour to black.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGBlack(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlack,
	}
	a.fgColour = aC
}

// AttrFGWhite sets the foreground colour to white.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWBlack], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGWhite(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourWhite,
	}
	a.fgColour = aC
}

// AttrFGRed sets the foreground colour to red.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGBlack],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGRed(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourRed,
	}
	a.fgColour = aC
}

// AttrFGGreen sets the foreground colour to green.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGBlack] and [AttrFGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGGreen(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourGreen,
	}
	a.fgColour = aC
}

// AttrFGBlue sets the foreground colour to blue.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlack]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrFGBlue(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlue,
	}
	a.fgColour = aC
}

// AttrBGShort returns an AttrFunc that will set the background colour to the
// colour represented by the supplied short code.
//
// See also [AttrBGColour], [AttrBGBlack], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGShort(c uint8) AttrFunc {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    c,
	}
	return func(a *Attr) {
		a.bgColour = aC
	}
}

// AttrBGColour returns an AttrFunc that will set the background colour to the
// colour represented by the supplied colour.
//
// See also [AttrBGShort], [AttrBGBlack], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGColour(c color.RGBA) AttrFunc {
	aC := attrColour{
		useRGB: true,
		red:    c.R,
		green:  c.G,
		blue:   c.B,
	}
	return func(a *Attr) {
		a.bgColour = aC
	}
}

// AttrBGBlack sets the background colour to black.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGBlack(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlack,
	}
	a.bgColour = aC
}

// AttrBGWhite sets the background colour to white.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWBlack], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGWhite(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourWhite,
	}
	a.bgColour = aC
}

// AttrBGRed sets the background colour to red.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGBlack],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGRed(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourRed,
	}
	a.bgColour = aC
}

// AttrBGGreen sets the background colour to green.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGBlack] and [AttrBGBlue]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGGreen(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourGreen,
	}
	a.bgColour = aC
}

// AttrBGBlue sets the background colour to blue.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlack]. Only one of these should be passed when
// creating a [Attr]. If more than one is given only the last takes effect.
func AttrBGBlue(a *Attr) {
	aC := attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlue,
	}
	a.bgColour = aC
}
