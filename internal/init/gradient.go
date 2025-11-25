// Package init provides initialization functionality for Spectr,
// including an interactive wizard for project setup and gradient
// rendering utilities for enhanced visual presentation.
package init

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// ANSI 256 color code constants
const (
	ansiMaxColorCode     = 255
	ansiStandardMax      = 16
	ansiCubeStart        = 16
	ansiCubeEnd          = 231
	ansiGrayscaleStart   = 232
	ansiGrayscaleEnd     = 255
	ansiCubeSize         = 6
	ansiCubePlaneSize    = 36 // 6 * 6
	ansiGrayscaleSteps   = 23.0
	ansiColorSteps       = 5.0
	standardColorDim     = 0.5
	standardColorBright  = 0.75
	fullBrightness       = 1.0
	zeroBrightness       = 0.0
	singleCharacterRatio = 0.0
)

// applyGradient applies a color gradient from colorA to colorB
// across text. It processes multi-line ASCII art character-by-
// character. Supports hex colors (#RRGGBB) and ANSI 256 codes.
func applyGradient(
	text string,
	colorA, colorB lipgloss.Color,
) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	startColor, endColor, err := parseColorPair(colorA, colorB)
	if err != nil {
		return text // Fallback to unstyled text
	}

	totalChars := countTotalChars(lines)
	if totalChars == 0 {
		return text
	}

	return renderGradientText(lines, startColor, endColor, totalChars)
}

// parseColorPair parses both gradient colors or returns an error
func parseColorPair(
	colorA, colorB lipgloss.Color,
) (start, end colorful.Color, err error) {
	start, err = parseColor(string(colorA))
	if err != nil {
		return colorful.Color{}, colorful.Color{}, err
	}

	end, err = parseColor(string(colorB))
	if err != nil {
		return colorful.Color{}, colorful.Color{}, err
	}

	return start, end, nil
}

// countTotalChars counts characters across all lines
func countTotalChars(lines []string) int {
	total := 0
	for _, line := range lines {
		total += len(line)
	}

	return total
}

// renderGradientText builds the final gradient-styled string
func renderGradientText(
	lines []string,
	startColor, endColor colorful.Color,
	totalChars int,
) string {
	var result strings.Builder
	charIndex := 0

	for lineIdx, line := range lines {
		if lineIdx > 0 {
			result.WriteString("\n")
		}

		for _, char := range line {
			ratio := calculateColorRatio(charIndex, totalChars)
			styledChar := styleCharacter(
				char,
				startColor,
				endColor,
				ratio,
			)
			result.WriteString(styledChar)
			charIndex++
		}
	}

	return result.String()
}

// calculateColorRatio determines interpolation ratio for a position
func calculateColorRatio(charIndex, totalChars int) float64 {
	if totalChars == 1 {
		return singleCharacterRatio
	}

	return float64(charIndex) / float64(totalChars-1)
}

// styleCharacter applies gradient color to a single character
func styleCharacter(
	char rune,
	startColor, endColor colorful.Color,
	ratio float64,
) string {
	interpolated := startColor.BlendLab(endColor, ratio)
	hexColor := interpolated.Hex()

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(hexColor)).
		Render(string(char))
}

// parseColor converts a lipgloss color to a colorful.Color.
// Supports hex format (#RRGGBB) and ANSI 256 color codes.
func parseColor(color string) (colorful.Color, error) {
	if strings.HasPrefix(color, "#") {
		return colorful.Hex(color)
	}

	colorCode, err := strconv.Atoi(color)
	if err == nil && colorCode >= 0 && colorCode <= ansiMaxColorCode {
		return ansi256ToRGB(colorCode), nil
	}

	return colorful.Color{}, fmt.Errorf(
		"invalid color format: %s",
		color,
	)
}

// ansi256ToRGB converts an ANSI 256 color code to RGB values
func ansi256ToRGB(code int) colorful.Color {
	if code < ansiStandardMax {
		return getStandardColor(code)
	}

	if code >= ansiCubeStart && code <= ansiCubeEnd {
		return getColorCubeColor(code)
	}

	if code >= ansiGrayscaleStart && code <= ansiGrayscaleEnd {
		return getGrayscaleColor(code)
	}

	// Default to white if out of range
	return colorful.Color{
		R: fullBrightness,
		G: fullBrightness,
		B: fullBrightness,
	}
}

// getStandardColor returns one of the 16 standard ANSI colors
func getStandardColor(code int) colorful.Color {
	standardColors := [ansiStandardMax]colorful.Color{
		{R: zeroBrightness, G: zeroBrightness, B: zeroBrightness},
		{R: standardColorDim, G: zeroBrightness, B: zeroBrightness},
		{R: zeroBrightness, G: standardColorDim, B: zeroBrightness},
		{R: standardColorDim, G: standardColorDim, B: zeroBrightness},
		{R: zeroBrightness, G: zeroBrightness, B: standardColorDim},
		{R: standardColorDim, G: zeroBrightness, B: standardColorDim},
		{R: zeroBrightness, G: standardColorDim, B: standardColorDim},
		{
			R: standardColorBright,
			G: standardColorBright,
			B: standardColorBright,
		},
		{R: standardColorDim, G: standardColorDim, B: standardColorDim},
		{R: fullBrightness, G: zeroBrightness, B: zeroBrightness},
		{R: zeroBrightness, G: fullBrightness, B: zeroBrightness},
		{R: fullBrightness, G: fullBrightness, B: zeroBrightness},
		{R: zeroBrightness, G: zeroBrightness, B: fullBrightness},
		{R: fullBrightness, G: zeroBrightness, B: fullBrightness},
		{R: zeroBrightness, G: fullBrightness, B: fullBrightness},
		{R: fullBrightness, G: fullBrightness, B: fullBrightness},
	}

	return standardColors[code]
}

// getColorCubeColor returns a color from the 216-color cube
func getColorCubeColor(code int) colorful.Color {
	index := code - ansiCubeStart
	r := index / ansiCubePlaneSize
	g := (index % ansiCubePlaneSize) / ansiCubeSize
	b := index % ansiCubeSize

	return colorful.Color{
		R: float64(r) / ansiColorSteps,
		G: float64(g) / ansiColorSteps,
		B: float64(b) / ansiColorSteps,
	}
}

// getGrayscaleColor returns a color from the 24-step grayscale ramp
func getGrayscaleColor(code int) colorful.Color {
	gray := float64(code-ansiGrayscaleStart) / ansiGrayscaleSteps

	return colorful.Color{R: gray, G: gray, B: gray}
}
