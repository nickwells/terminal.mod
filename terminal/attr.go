package terminal

import (
	"fmt"
	"image/color" //nolint:misspell
	"io"
	"os"
	"strconv"
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

// UnderlineStyle represents the allowed underline styles
type UnderlineStyle int

const (
	// UlNone is the underline style that gives no underlining
	UlNone UnderlineStyle = iota
	// UlSingle is the underline style that gives a single underline
	UlSingle
	// UlDouble is the underline style that gives a double underline
	UlDouble
)

var underlineMap = map[UnderlineStyle]string{
	UlNone:   sgrNoUnderline,
	UlSingle: sgrUnderline,
	UlDouble: sgrDoubleUnderline,
}

// Blinking

// BlinkStyle represents the allowed blink styles
type BlinkStyle int

const (
	// BlkNone is the blink style that gives no blinking
	BlkNone BlinkStyle = iota
	// BlkSlow is the blink style that gives slow blinking
	BlkSlow
	// BlkFast is the blink style that gives fast blinking
	BlkFast
)

var blinkMap = map[BlinkStyle]string{
	BlkNone: sgrNoBlinking,
	BlkSlow: sgrBlinkSlow,
	BlkFast: sgrBlinkFast,
}

// Intensity

// IntensityStyle represents the allowed intensity styles
type IntensityStyle int

const (
	// ItyNormal is the intensity style that gives standard intensity
	ItyNormal IntensityStyle = iota
	// ItyBold is the intensity style that gives bold output
	ItyBold
	// ItyFaint is the intensity style that gives faint output
	ItyFaint
)

var intensityMap = map[IntensityStyle]string{
	ItyNormal: sgrNormalIntensity,
	ItyBold:   sgrBold,
	ItyFaint:  sgrFaint,
}

// colourCoder is an interface that an attribute colour must satisfy
type colourCoder interface {
	colourCodes() []string
}

// attrColourIdx represents a colour as used in a terminal where the colour
// is given as an index into the set of Standard or High Intensity colours
type attrColourIdx uint8

// colourCodes returns the strings used to represent the index as a colour
func (ac attrColourIdx) colourCodes() []string {
	return []string{
		sgrColourIdxIntro,
		strconv.Itoa(int(ac)),
	}
}

// attrColourRGB represents a colour as used in a terminal where the colour
// is given as an RGB triplet
type attrColourRGB struct {
	red   uint8
	green uint8
	blue  uint8
}

// colourCodes returns the strings used to represent the red, green and blue
// members as a colour
func (ac attrColourRGB) colourCodes() []string {
	return []string{
		sgrColourRGBIntro,
		strconv.Itoa(int(ac.red)),
		strconv.Itoa(int(ac.green)),
		strconv.Itoa(int(ac.blue)),
	}
}

// Attr represents the attributes that text should have
type Attr struct {
	overlined  bool
	crossedOut bool
	inverted   bool
	italic     bool

	hasUnderline bool
	underline    UnderlineStyle

	hasBlink bool
	blink    BlinkStyle

	hasIntensity bool
	intensity    IntensityStyle

	hasFGColour bool
	fgColour    colourCoder

	hasBGColour bool
	bgColour    colourCoder
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

	if a.inverted {
		attrs = append(attrs, sgrInvert)
	}

	if a.italic {
		attrs = append(attrs, sgrItalic)
	}

	if a.hasUnderline {
		attrs = append(attrs, underlineMap[a.underline])
	}

	if a.hasBlink {
		attrs = append(attrs, blinkMap[a.blink])
	}

	if a.hasIntensity {
		attrs = append(attrs, intensityMap[a.intensity])
	}

	if a.hasFGColour {
		attrs = append(attrs, sgrFGColour)
		attrs = append(attrs, a.fgColour.colourCodes()...)
	}

	if a.hasBGColour {
		attrs = append(attrs, sgrBGColour)
		attrs = append(attrs, a.bgColour.colourCodes()...)
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

// AttrUnderline returns a function that sets the underline style to the
// supplied value. This should only be passed once when creating an
// [Attr]. If it is passed more than once only the last takes effect.
//
// This will panic if an unknown underline style is given; use the named
// constants.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrUnderline(ul UnderlineStyle) AttrFunc {
	if _, ok := underlineMap[ul]; !ok {
		panic(fmt.Errorf("unknown underline style: %d", ul))
	}

	return func(a *Attr) {
		a.hasUnderline = true
		a.underline = ul
	}
}

// AttrBlink returns a function that sets the blink style to the supplied
// value. This should only be passed once when creating an [Attr]. If it is
// passed more than once only the last takes effect.
//
// This will panic if an unknown blink style is given; use the named
// constants.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrBlink(b BlinkStyle) AttrFunc {
	if _, ok := blinkMap[b]; !ok {
		panic(fmt.Errorf("unknown blink style: %d", b))
	}

	return func(a *Attr) {
		a.hasBlink = true
		a.blink = b
	}
}

// AttrIntensity returns a function that sets the intensity style to the
// supplied value. This should only be passed once when creating an
// [Attr]. If it is passed more than once only the last takes effect.
//
// This will panic if an unknown intensity style is given; use the named
// constants.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrIntensity(i IntensityStyle) AttrFunc {
	if _, ok := intensityMap[i]; !ok {
		panic(fmt.Errorf("unknown intensity style: %d", i))
	}

	return func(a *Attr) {
		a.hasIntensity = true
		a.intensity = i
	}
}

// AttrFGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [AttrFGColour], [AttrFGBlack], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGGrey(g uint8) AttrFunc {
	return func(a *Attr) {
		a.hasFGColour = true
		a.fgColour = attrColourRGB{
			red:   g,
			green: g,
			blue:  g,
		}
	}
}

// AttrBGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [AttrBGColour], [AttrBGBlack], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGGrey(g uint8) AttrFunc {
	return func(a *Attr) {
		a.hasBGColour = true
		a.bgColour = attrColourRGB{
			red:   g,
			green: g,
			blue:  g,
		}
	}
}

// AttrFGShort returns an [AttrFunc] that will set the foreground colour to
// the colour represented by the supplied short code.
//
// See also [AttrFGColour], [AttrFGBlack], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGShort(c uint8) AttrFunc {
	return func(a *Attr) {
		a.hasFGColour = true
		a.fgColour = attrColourIdx(c)
	}
}

// AttrFGColour returns an [AttrFunc] that will set the foreground colour to
// the colour represented by the supplied colour.
//
// See also [AttrFGShort], [AttrFGBlack], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This returns an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGColour(c color.RGBA) AttrFunc { //nolint:misspell
	return func(a *Attr) {
		a.hasFGColour = true
		a.fgColour = attrColourRGB{
			red:   c.R,
			green: c.G,
			blue:  c.B,
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
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourBlack)
}

// AttrFGWhite sets the foreground colour to white.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGBlack], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGWhite(a *Attr) {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourWhite)
}

// AttrFGRed sets the foreground colour to red.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGBlack],
// [AttrFGGreen] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGRed(a *Attr) {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourHIRed)
}

// AttrFGGreen sets the foreground colour to green.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGBlack] and [AttrFGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGGreen(a *Attr) {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourHIGreen)
}

// AttrFGBlue sets the foreground colour to blue.
//
// See also [AttrFGShort], [AttrFGColour], [AttrFGWhite], [AttrFGRed],
// [AttrFGGreen] and [AttrFGBlack]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrFGBlue(a *Attr) {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourHIBlue)
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
		a.hasBGColour = true
		a.bgColour = attrColourIdx(c)
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
		a.hasBGColour = true
		a.bgColour = attrColourRGB{
			red:   c.R,
			green: c.G,
			blue:  c.B,
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
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourBlack)
}

// AttrBGWhite sets the background colour to white.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGBlack], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGWhite(a *Attr) {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourWhite)
}

// AttrBGRed sets the background colour to red.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGBlack],
// [AttrBGGreen] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGRed(a *Attr) {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourHIRed)
}

// AttrBGGreen sets the background colour to green.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGBlack] and [AttrBGBlue]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGGreen(a *Attr) {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourHIGreen)
}

// AttrBGBlue sets the background colour to blue.
//
// See also [AttrBGShort], [AttrBGColour], [AttrBGWhite], [AttrBGRed],
// [AttrBGGreen] and [AttrBGBlack]. Only one of these should be passed when
// creating an [Attr]. If more than one is given only the last takes effect.
//
// This is an [AttrFunc] suitable for passing to [MkAttr].
func AttrBGBlue(a *Attr) {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourHIBlue)
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

// SetUnderline sets the Underline style to the given value. It returns a
// non-nil error if the style is not recognised; use the named constants.
func (a *Attr) SetUnderline(us UnderlineStyle) error {
	if _, ok := underlineMap[us]; !ok {
		return fmt.Errorf("unknown underline style: %d", us)
	}

	a.hasUnderline = true
	a.underline = us

	return nil
}

// SetBlink sets the Blink style to the given value. It returns a
// non-nil error if the style is not recognised; use the named constants.
func (a *Attr) SetBlink(bs BlinkStyle) error {
	if _, ok := blinkMap[bs]; !ok {
		return fmt.Errorf("unknown blink style: %d", bs)
	}

	a.hasBlink = true
	a.blink = bs

	return nil
}

// SetIntensity sets the Intensity style to the given value. It returns a
// non-nil error if the style is not recognised; use the named constants.
func (a *Attr) SetIntensity(is IntensityStyle) error {
	if _, ok := intensityMap[is]; !ok {
		return fmt.Errorf("unknown intensity style: %d", is)
	}

	a.hasIntensity = true
	a.intensity = is

	return nil
}

// SetFGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [SetFGColour], [SetFGBlack], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGGrey(g uint8) {
	a.hasFGColour = true
	a.fgColour = attrColourRGB{
		red:   g,
		green: g,
		blue:  g,
	}
}

// SetBGGrey returns an [AttrFunc] that will set the foreground
// colour to the colour on a grey scale where 0
// represents black and 255 represents white.
//
// See also [SetBGColour], [SetBGBlack], [SetBGWWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGGrey(g uint8) {
	a.hasBGColour = true
	a.bgColour = attrColourRGB{
		red:   g,
		green: g,
		blue:  g,
	}
}

// SetFGShort sets the foreground colour to the colour represented by the
// supplied short code.
//
// See also [SetFGColour], [SetFGBlack], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGShort(c uint8) {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(c)
}

// SetFGColour sets the foreground colour to the colour represented by the
// supplied colour.
//
// See also [SetFGShort], [SetFGBlack], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGColour(c color.RGBA) { //nolint:misspell
	a.hasFGColour = true
	a.fgColour = attrColourRGB{
		red:   c.R,
		green: c.G,
		blue:  c.B,
	}
}

// SetFGBlack sets the foreground colour to black.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGBlack() {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourBlack)
}

// SetFGWhite sets the foreground colour to white.
//
// See also [SetFGShort], [SetFGColour], [SetFGBlack], [SetFGRed],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGWhite() {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourWhite)
}

// SetFGRed sets the foreground colour to red.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGBlack],
// [SetFGGreen] and [SetFGBlue].
func (a *Attr) SetFGRed() {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourHIRed)
}

// SetFGGreen sets the foreground colour to green.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGRed],
// [SetFGBlack] and [SetFGBlue].
func (a *Attr) SetFGGreen() {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourHIGreen)
}

// SetFGBlue sets the foreground colour to blue.
//
// See also [SetFGShort], [SetFGColour], [SetFGWhite], [SetFGRed],
// [SetFGGreen] and [SetFGBlack].
func (a *Attr) SetFGBlue() {
	a.hasFGColour = true
	a.fgColour = attrColourIdx(sgrColourHIBlue)
}

// SetBGShort sets the background colour to the colour represented by the
// supplied short code.
//
// See also [SetBGColour], [SetBGBlack], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGShort(c uint8) {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(c)
}

// SetBGColour sets the background colour to the colour represented by the
// supplied colour.
//
// See also [SetBGShort], [SetBGBlack], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGColour(c color.RGBA) { //nolint:misspell
	a.hasBGColour = true
	a.bgColour = attrColourRGB{
		red:   c.R,
		green: c.G,
		blue:  c.B,
	}
}

// SetBGBlack sets the background colour to black.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGBlack() {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourBlack)
}

// SetBGWhite sets the background colour to white.
//
// See also [SetBGShort], [SetBGColour], [SetBGBlack], [SetBGRed],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGWhite() {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourWhite)
}

// SetBGRed sets the background colour to red.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGBlack],
// [SetBGGreen] and [SetBGBlue].
func (a *Attr) SetBGRed() {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourHIRed)
}

// SetBGGreen sets the background colour to green.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGRed],
// [SetBGBlack] and [SetBGBlue].
func (a *Attr) SetBGGreen() {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourHIGreen)
}

// SetBGBlue sets the background colour to blue.
//
// See also [SetBGShort], [SetBGColour], [SetBGWhite], [SetBGRed],
// [SetBGGreen] and [SetBGBlack].
func (a *Attr) SetBGBlue() {
	a.hasBGColour = true
	a.bgColour = attrColourIdx(sgrColourHIBlue)
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
