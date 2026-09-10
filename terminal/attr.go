package terminal

import (
	"fmt"
	"image/color" //nolint:misspell
	"io"
	"os"
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

	sgrColourBlack = 0
	// sgrColourStdRed     = 1
	// sgrColourStdGreen   = 2
	// sgrColourStdYellow  = 3
	// sgrColourStdBlue    = 4
	// sgrColourStdMagenta = 5
	// sgrColourStdCyan    = 6
	// sgrColourPaleGrey   = 7
	// sgrColourDarkGrey   = 8
	sgrColourHIRed   = 9
	sgrColourHIGreen = 10
	// sgrColourHIYellow   = 11
	sgrColourHIBlue = 12
	// sgrColourHIMagenta  = 13
	// sgrColourHICyan     = 14
	sgrColourWhite = 15
)

// Underlining

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

// Blinking

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

// Intensity

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
func (Attr) End(w io.Writer) {
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
//
// This is an [AttrFunc] suitable for passing to [MkAttr]
func AttrOverlined(a *Attr) {
	a.overlined = true
}

// AttrCrossedOut sets the CrossedOut attribute.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrCrossedOut(a *Attr) {
	a.crossedOut = true
}

// AttrInverted sets the Inverted attribute.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrInverted(a *Attr) {
	a.inverted = true
}

// AttrItalic sets the Italic attribute.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrItalic(a *Attr) {
	a.italic = true
}

// AttrUnderline sets the underline style to single. See also
// [AttrDoubleUnderline] and [AttrNoUnderline]. Only one of these should be
// passed when creating a [Attr]. If more than one is given only the last
// takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrUnderline(a *Attr) {
	a.underline = ulSingle
}

// AttrDoubleUnderline sets the underline style to double. See also
// [AttrUnderline] and [AttrNoUnderline]. Only one of these should be passed
// when creating an [Attr]. If more than one is given only the last takes
// effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrDoubleUnderline(a *Attr) {
	a.underline = ulDouble
}

// AttrNoUnderline turns off Underlining. See also [AttrUnderline] and
// [AttrDoubleUnderline]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrNoUnderline(a *Attr) {
	a.underline = ulNone
}

// AttrBlinkFast sets the Blink attribute to fast. See also [AttrBlinkSlow]
// and [AttrNoBlink]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBlinkFast(a *Attr) {
	a.blink = blkFast
}

// AttrBlinkSlow sets the Blink attribute to slow. See also [AttrBlinkFast]
// and [AttrNoBlink]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBlinkSlow(a *Attr) {
	a.blink = blkSlow
}

// AttrNoBlink turns off blinking. See also [AttrBlinkFast] and
// [AttrBlinkSlow]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrNoBlink(a *Attr) {
	a.blink = blkNone
}

// AttrIntensityNormal sets the Intensity attribute to its normal value. See
// also [AttrBold] and [AttrFaint]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrIntensityNormal(a *Attr) {
	a.intensity = ityStd
}

// AttrBold sets the Intensity attribute to bold. See also
// [AttrIntensityNormal] and [AttrFaint]. Only one of these should be passed
// when creating an [Attr]. If more than one is given only the last takes
// effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBold(a *Attr) {
	a.intensity = ityStd
}

// AttrFaint sets the Intensity attribute to faint. See also [AttrBold] and
// [AttrIntensityNormal]. Only one of these should be passed when creating a
// [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFaint(a *Attr) {
	a.intensity = ityStd
}

// AttrFGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [AttrFGColour], [AttrFGWBlack], [AttrFGWWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGGrey(g uint8) AttrFunc {
	return func(a *Attr) {
		a.fgColour = attrColour{
			useRGB: true,
			red:    g,
			green:  g,
			blue:   g,
		}
	}
}

// AttrBGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [AttrBGColour], [AttrBGWBlack], [AttrBGWWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGGrey(g uint8) AttrFunc {
	return func(a *Attr) {
		a.bgColour = attrColour{
			useRGB: true,
			red:    g,
			green:  g,
			blue:   g,
		}
	}
}

// AttrFGShort returns an [AttrFunc] that will set the foreground colour to
// the colour represented by the supplied short code.
//
// See also [AttrFGColour], [AttrFGWBlack], [AttrFGWWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGShort(c uint8) AttrFunc {
	return func(a *Attr) {
		a.fgColour = attrColour{
			useColourIdx: true,
			colourIdx:    c,
		}
	}
}

// AttrFGColour returns an [AttrFunc] that will set the foreground colour to
// the colour represented by the supplied colour.
//
// See also [AttrFGShort], [AttrFGWBlack], [AttrFGWWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGColour(c color.RGBA) AttrFunc { //nolint:misspell
	return func(a *Attr) {
		a.fgColour = attrColour{
			useRGB: true,
			red:    c.R,
			green:  c.G,
			blue:   c.B,
		}
	}
}

// AttrFGBlack sets the foreground colour to black.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGBlack(a *Attr) {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlack,
	}
}

// AttrFGWhite sets the foreground colour to white.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWBlack], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGWhite(a *Attr) {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourWhite,
	}
}

// AttrFGRed sets the foreground colour to red.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGBlack],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGRed(a *Attr) {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIRed,
	}
}

// AttrFGGreen sets the foreground colour to green.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGBlack] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGGreen(a *Attr) {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIGreen,
	}
}

// AttrFGBlue sets the foreground colour to blue.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlack]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGBlue(a *Attr) {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIBlue,
	}
}

// AttrBGShort returns an [AttrFunc] that will set the background colour to
// the colour represented by the supplied short code.
//
// See also [AttrBGColour], [AttrBGBlack], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGShort(c uint8) AttrFunc {
	return func(a *Attr) {
		a.bgColour = attrColour{
			useColourIdx: true,
			colourIdx:    c,
		}
	}
}

// AttrBGColour returns an [AttrFunc] that will set the background colour to
// the colour represented by the supplied colour.
//
// See also [AttrBGShort], [AttrBGBlack], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGColour(c color.RGBA) AttrFunc { //nolint:misspell
	return func(a *Attr) {
		a.bgColour = attrColour{
			useRGB: true,
			red:    c.R,
			green:  c.G,
			blue:   c.B,
		}
	}
}

// AttrBGBlack sets the background colour to black.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGBlack(a *Attr) {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlack,
	}
}

// AttrBGWhite sets the background colour to white.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWBlack], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGWhite(a *Attr) {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourWhite,
	}
}

// AttrBGRed sets the background colour to red.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGBlack],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGRed(a *Attr) {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIRed,
	}
}

// AttrBGGreen sets the background colour to green.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGBlack] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGGreen(a *Attr) {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIGreen,
	}
}

// AttrBGBlue sets the background colour to blue.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlack]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGBlue(a *Attr) {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIBlue,
	}
}

// SetOverlined sets the Overlined attribute.
func (a *Attr) SetOverlined() {
	a.overlined = true
}

// SetCrossedOut sets the CrossedOut attribute.
func (a *Attr) SetCrossedOut() {
	a.crossedOut = true
}

// SetInverted sets the Inverted attribute.
func (a *Attr) SetInverted() {
	a.inverted = true
}

// SetItalic sets the Italic attribute.
func (a *Attr) SetItalic() {
	a.italic = true
}

// SetUnderline sets the underline style to single. See also
// [SetDoubleUnderline] and [SetNoUnderline].
func (a *Attr) SetUnderline() {
	a.underline = ulSingle
}

// SetDoubleUnderline sets the underline style to double. See also
// [SetUnderline] and [SetNoUnderline].
func (a *Attr) SetDoubleUnderline() {
	a.underline = ulDouble
}

// SetNoUnderline turns off Underlining. See also [SetUnderline] and
// [SetDoubleUnderline].
func (a *Attr) SetNoUnderline() {
	a.underline = ulNone
}

// SetBlinkFast sets the Blink attribute to fast. See also [SetBlinkSlow]
// and [SetNoBlink].
func (a *Attr) SetBlinkFast() {
	a.blink = blkFast
}

// SetBlinkSlow sets the Blink attribute to slow. See also [SetBlinkFast]
// and [SetNoBlink].
func (a *Attr) SetBlinkSlow() {
	a.blink = blkSlow
}

// SetNoBlink turns off blinking. See also [SetBlinkFast] and
// [SetBlinkSlow].
func (a *Attr) SetNoBlink() {
	a.blink = blkNone
}

// SetIntensityNormal sets the Intensity attribute to its normal value. See
// also [SetBold] and [SetFaint].
func (a *Attr) SetIntensityNormal() {
	a.intensity = ityStd
}

// SetBold sets the Intensity attribute to bold. See also
// [SetIntensityNormal] and [SetFaint].
func (a *Attr) SetBold() {
	a.intensity = ityStd
}

// SetFaint sets the Intensity attribute to faint. See also [SetBold] and
// [SetIntensityNormal].
func (a *Attr) SetFaint() {
	a.intensity = ityStd
}

// SetFGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [SetFGColour], [SetFGWBlack], [SetFGWWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGGrey(g uint8) {
	a.fgColour = attrColour{
		useRGB: true,
		red:    g,
		green:  g,
		blue:   g,
	}
}

// SetBGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [SetBGColour], [SetBGWBlack], [SetBGWWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGGrey(g uint8) {
	a.bgColour = attrColour{
		useRGB: true,
		red:    g,
		green:  g,
		blue:   g,
	}
}

// SetFGShort sets the foreground colour to the colour represented by the
// supplied short code.
//
// See also [SetFGColour], [SetFGWBlack], [SetFGWWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGShort(c uint8) {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    c,
	}
}

// SetFGColour sets the foreground colour to the colour represented by the
// supplied colour.
//
// See also [SetFGShort], [SetFGWBlack], [SetFGWWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGColour(c color.RGBA) { //nolint:misspell
	a.fgColour = attrColour{
		useRGB: true,
		red:    c.R,
		green:  c.G,
		blue:   c.B,
	}
}

// SetFGBlack sets the foreground colour to black.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGBlack() {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlack,
	}
}

// SetFGWhite sets the foreground colour to white.
//
// See also [SetFGShort], [SetFGColour], [SetFGWBlack], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGWhite() {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourWhite,
	}
}

// SetFGRed sets the foreground colour to red.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGBlack],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGRed() {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIRed,
	}
}

// SetFGGreen sets the foreground colour to green.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGRed],
// [SetFGBlack] and [SetFGBlue].
func (a *Attr) SetFGGreen() {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIGreen,
	}
}

// SetFGBlue sets the foreground colour to blue.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlack].
func (a *Attr) SetFGBlue() {
	a.fgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIBlue,
	}
}

// SetBGShort sets the background colour to the colour represented by the
// supplied short code.
//
// See also [SetBGColour], [SetBGBlack], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGShort(c uint8) {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    c,
	}
}

// SetBGColour sets the background colour to the colour represented by the
// supplied colour.
//
// See also [SetBGShort], [SetBGBlack], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGColour(c color.RGBA) { //nolint:misspell
	a.bgColour = attrColour{
		useRGB: true,
		red:    c.R,
		green:  c.G,
		blue:   c.B,
	}
}

// SetBGBlack sets the background colour to black.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGBlack() {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourBlack,
	}
}

// SetBGWhite sets the background colour to white.
//
// See also [SetBGShort], [SetBGColour], [SetBGWBlack], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGWhite() {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourWhite,
	}
}

// SetBGRed sets the background colour to red.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGBlack],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGRed() {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIRed,
	}
}

// SetBGGreen sets the background colour to green.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGRed],
// [SetBGBlack] and [SetBGBlue].
func (a *Attr) SetBGGreen() {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIGreen,
	}
}

// SetBGBlue sets the background colour to blue.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlack].
func (a *Attr) SetBGBlue() {
	a.bgColour = attrColour{
		useColourIdx: true,
		colourIdx:    sgrColourHIBlue,
	}
}

// Fprint will print the arguments to the supplied io.Writer with the
// associated attributes set
func (a Attr) Fprint(w io.Writer, args ...any) {
	a.Start(w)
	fmt.Fprint(w, args...)
	a.End(w)
}

// Fprintln will print the arguments and a newline to the supplied io.Writer
// with the associated attributes set
func (a Attr) Fprintln(w io.Writer, args ...any) {
	a.Start(w)
	fmt.Fprintln(w, args...)
	a.End(w)
}

// Fprintf will format and print the arguments to the supplied io.Writer with
// the associated attributes set
func (a Attr) Fprintf(w io.Writer, f string, args ...any) {
	a.Start(w)
	fmt.Fprintf(w, f, args...)
	a.End(w)
}

// Print will print the arguments with the associated attributes set
func (a Attr) Print(args ...any) {
	a.Fprint(os.Stdout, args...)
}

// Println will print the arguments and a newline with the associated
// attributes set
func (a Attr) Println(args ...any) {
	a.Fprintln(os.Stdout, args...)
}

// Printf will format and print the arguments with the associated attributes
// set
func (a Attr) Printf(f string, args ...any) {
	a.Fprintf(os.Stdout, f, args...)
}
