package shared

import (
	"hash/fnv"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	log "github.com/sirupsen/logrus"
)

func HashOfString(text string) uint32 {
	algorithm := fnv.New32a()
	if _, err := algorithm.Write([]byte(text)); err != nil {
		log.Error(err)
		return 0
	}
	return algorithm.Sum32()
}

func CenteredText(text string) fyne.CanvasObject {
	color := color.NRGBA{R: 0x10, G: 0xff, B: 0xff, A: 0xaa}
	title := canvas.NewText(text, color)
	title.TextSize = 14
	title.TextStyle = fyne.TextStyle{Bold: true}
	return container.New(layout.NewHBoxLayout(),
		layout.NewSpacer(),
		title,
		layout.NewSpacer(),
	)
}

func GetDate(text string) time.Time {
	var err error

	if buffer := strings.TrimSpace(text); buffer != "" {
		if items := strings.Split(buffer, "-"); len(items) == 3 {
			if len(items[0]) == 4 && len(items[1]) == 2 && len(items[2]) == 2 {
				if year, err := strconv.Atoi(items[0]); err == nil {
					if month, err := strconv.Atoi(items[1]); err == nil {
						if day, err := strconv.Atoi(items[2]); err == nil {
							return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
						}
					}
				}
			}
		}
	}
	log.Error(err)
	return time.Time{}
}
