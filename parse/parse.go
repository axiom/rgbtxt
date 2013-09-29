package parse

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"
	"sync"
)

type ColorPair struct {
	Name  string
	Color color.RGBA
}

func (c ColorPair) Go() string {
	r, g, b, a := c.Color.RGBA()
	return fmt.Sprintf("%s = color.RGBA{%v, %v, %v, %v}\n",
		c.Name,
		r&0x00ff,
		g&0x00ff,
		b&0x00ff,
		a&0x00ff,
	)
}

func ParseLinesChan(input io.Reader, pairs chan ColorPair) error {
	lines := make(chan string, 100)
	var err error

	go func() {
		bufferedInput := bufio.NewReader(input)
		for {
			if line, err := bufferedInput.ReadString('\n'); err == io.EOF || err == nil {
				if line != "" {
					lines <- line
				}

				if err == io.EOF {
					err = nil
					break
				}
			} else {
				break
			}
		}
		close(lines)
	}()

	var wg sync.WaitGroup
	for w := 0; w < 5; w++ {
		wg.Add(1)
		go func() {
			for line := range lines {
				if colorPair, err := ParseLine(line); err == nil {
					pairs <- colorPair
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(pairs)
	return err
}

func ParseLines(input io.Reader) ([]ColorPair, error) {
	pairs := make(chan ColorPair, 100)
	collectedPairs := make([]ColorPair, 0, 1000)
	go func() {
		for colorPair := range pairs {
			collectedPairs = append(collectedPairs, colorPair)
		}
	}()

	err := ParseLinesChan(input, pairs)
	return collectedPairs, err
}

// Parse a line from rgb.txt.
func ParseLine(line string) (ColorPair, error) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return ColorPair{}, fmt.Errorf("Could not parse line, too few fields")
	}

	name := Camelize(strings.Join(parts[3:], " "))

	components := parts[0:3]
	if len(components) != 3 {
		return ColorPair{}, fmt.Errorf("Expected three color components")
	}

	red, err := strconv.Atoi(components[0])
	if err != nil {
		return ColorPair{}, fmt.Errorf("Could not parse color component")
	}
	green, err := strconv.Atoi(components[1])
	if err != nil {
		return ColorPair{}, fmt.Errorf("Could not parse color component")
	}
	blue, err := strconv.Atoi(components[2])
	if err != nil {
		return ColorPair{}, fmt.Errorf("Could not parse color component")
	}

	if red < 0 || red > 255 || blue < 0 || blue > 255 || green < 0 || green > 255 {
		return ColorPair{}, fmt.Errorf("Component out of range")
	}

	return ColorPair{
		Name:  name,
		Color: color.RGBA{uint8(red), uint8(green), uint8(blue), 255},
	}, nil
}

// Camelize a possibly space seperated string of words.
func Camelize(name string) string {
	parts := strings.Fields(name)
	for i, part := range parts {
		parts[i] = camelizeWord(part)
	}
	return strings.Join(parts, "")
}

func camelizeWord(name string) string {
	return strings.TrimSpace((strings.ToUpper(name[0:1])) + name[1:])
}
