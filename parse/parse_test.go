package parse_test

import (
	"bytes"
	"github.com/axiom/rgbtxt/parse"
	"image/color"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

func (s *S) TestParseLine(c *C) {
	testData := []struct {
		line     string
		expected parse.ColorPair
	}{
		{"0	0 0  		black", parse.ColorPair{"Black", color.RGBA{0, 0, 0, 255}}},
		{"1 2	 3		purple haze", parse.ColorPair{"PurpleHaze", color.RGBA{1, 2, 3, 255}}},
		{"255 255 255				black hazel nut", parse.ColorPair{"BlackHazelNut", color.RGBA{255, 255, 255, 255}}},
	}
	for _, entry := range testData {
		colorPair, err := parse.ParseLine(entry.line)
		c.Assert(err, IsNil)
		c.Assert(colorPair.Name, Equals, entry.expected.Name)
		c.Assert(colorPair.Color, Equals, entry.expected.Color)
	}
}

func (s *S) TestParseLineBad(c *C) {
	_, err := parse.ParseLine("255 255 522		nut balls")
	c.Assert(err, NotNil)

	_, err = parse.ParseLine("255 255 252	")
	c.Assert(err, NotNil)

	_, err = parse.ParseLine("-1 255 522		nut balls")
	c.Assert(err, NotNil)

	_, err = parse.ParseLine("1 255 522		nut balls")
	c.Assert(err, NotNil)
}

func (s *S) TestParseLines(c *C) {
	lines := bytes.NewBufferString("0 0 0		black\n255 255 255		pure white")

	colorPairs, err := parse.ParseLines(lines)
	c.Assert(err, IsNil)
	c.Assert(len(colorPairs), Equals, 2)

	first, second := colorPairs[0], colorPairs[1]
	var black parse.ColorPair
	var white parse.ColorPair

	w := color.RGBA{0, 0, 0, 255}
	if first.Color == w {
		black = first
		white = second
	} else {
		black = second
		white = first
	}

	c.Assert(white.Name, Equals, "PureWhite")
	c.Assert(white.Color, Equals, color.RGBA{255, 255, 255, 255})
	c.Assert(black.Name, Equals, "Black")
	c.Assert(black.Color, Equals, color.RGBA{0, 0, 0, 255})
}
