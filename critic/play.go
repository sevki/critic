package critic

import (
	"fmt"
	"image"
	"image/color"
	"sort"
)

func Analyze(i image.Image, plt color.Palette) []ColorArtColor {

	b := i.Bounds()

	colorMap := make(map[int]int)

	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {

			col := plt.Index(i.At(x, y))
			c, ok := colorMap[col]
			if ok {
				colorMap[col] = c + 1
			} else {
				colorMap[col] = 1
			}
		}
	}

	var colors []ColorArtColor
	for i, n := range colorMap {
		r, g, b, _ := plt[i].RGBA()
		hex := RGBToHex(r, g, b)
		c := ColorArtColor{Name: X11names[i], Frequency: n, Color: X11[i], Hex: hex}
		colors = append(colors, c)
	}
	sort.Sort(ByFrequency(colors))
	return colors
}
func AnalyzeAndConvert(i image.Image, plt color.Palette) ([]ColorArtColor, *image.RGBA) {

	ni := image.NewRGBA(i.Bounds())

	b := i.Bounds()

	colorMap := make(map[int]int)

	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {

			col := plt.Index(i.At(x, y))
			c, ok := colorMap[col]
			if ok {
				colorMap[col] = c + 1
			} else {
				colorMap[col] = 1
			}
			ni.Set(x, y, plt.Convert(i.At(x, y)))
		}
	}

	var colors []ColorArtColor
	for i, n := range colorMap {
		r, g, b, _ := plt[i].RGBA()
		hex := RGBToHex(r, g, b)
		c := ColorArtColor{Name: X11names[i], Frequency: n, Color: X11[i], Hex: hex}
		colors = append(colors, c)
	}
	sort.Sort(ByFrequency(colors))
	return colors, ni
}
func RGBToHex(r, g, b uint32) string {
	return fmt.Sprintf("#%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

type ColorArtColor struct {
	Name      string
	Color     color.Color
	Hex       string
	Frequency int
}
type ByFrequency []ColorArtColor

func (a ByFrequency) Len() int           { return len(a) }
func (a ByFrequency) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFrequency) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency }
