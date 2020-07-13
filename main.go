package main

import (
	"bufio"
	"fmt"
	"github.com/disintegration/imaging"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

func main() {
	s, err := getRandomThoughtText()
	if err != nil {
		log.Println(err)
	}

	if err := generateThoughtOfTheDay(s); err != nil {
		log.Fatal(err)
	}

	if err := generateDayInfo(); err != nil {
		log.Fatal(err)
	}
}

func getRandomThoughtText() (string, error) {
	defaultQuote := "\"Take up one idea. Make that one idea your life - think of it, dream of it, live on that idea. Let the brain, muscles, nerves, every part of your body be full of that idea, and just leave every other idea alone. This is the way to success.\" - Swami Vivekananda\n"
	quote := defaultQuote

	p := filepath.Join("assets", "thoughts.txt")

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return defaultQuote, fmt.Errorf("error while opening thoughts file %v", err)
	}

	var lines []string
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(lines)

	if len(lines) != 0 {
		quote = lines[n]
	}

	return quote, nil
}

func getBackgroundImage() (*gg.Context, error) {
	dc := gg.NewContext(1200, 628)

	baseImagePath := filepath.Join("assets", "base.png")
	backgroundImage, err := gg.LoadImage(baseImagePath)
	backgroundImage = imaging.Fill(backgroundImage, dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)
	if err != nil {
		return nil, fmt.Errorf("error while loading background image %v", err)
	}
	dc.DrawImage(backgroundImage, 0, 0)

	margin := 20.0
	x := margin
	y := margin
	w := float64(dc.Width()) - (2.0 * margin)
	h := float64(dc.Height()) - (2.0 * margin)
	dc.SetColor(color.RGBA{A: 100})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	return dc, nil
}

// https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html
func generateThoughtOfTheDay(title string) error {
	dc, err := getBackgroundImage()
	if err != nil {
		return err
	}

	fontPath := filepath.Join("assets", "fonts", "OpenSans-Bold.ttf")
	textColor := color.White
	if err := dc.LoadFontFace(fontPath, 50); err != nil {
		return fmt.Errorf("error while loading font %v", err)
	}

	textRightMargin := 60.0
	textTopMargin := 90.0
	x := textRightMargin
	y := textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin
	dc.SetColor(textColor)
	dc.DrawStringWrapped(title, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	if err := dc.SavePNG("./assets/quote.png"); err != nil {
		return fmt.Errorf("error while saving background image %v", err)
	}

	fmt.Println("generated quote ", title)

	return nil
}

func generateDayInfo() error {
	d := fmt.Sprintf(time.Now().Format("01 January 2006"))
	g := greetTime(time.Now())
	apiKey := os.Getenv("WEATHER_KEY")

	location := "Mumbai, India"

	w := weather{
		city:   "Mumbai",
		apiKey: apiKey,
	}
	wInfo, err := w.getWeatherInfo()
	if err != nil {
		return err
	}

	desc := wInfo.Current.WeatherDescriptions[0]
	temp := wInfo.Current.Temperature

	dc, err := getBackgroundImage()
	if err != nil {
		return err
	}

	fontPath := filepath.Join("assets", "fonts", "OpenSans-Bold.ttf")
	fontPathLight := filepath.Join("assets", "fonts", "OpenSans-Light.ttf")

	if err := dc.LoadFontFace(fontPath, 50); err != nil {
		return fmt.Errorf("error while loading font %v", err)
	}

	textRightMargin := 60.0
	textTopMargin := 90.0
	x := textRightMargin
	y := textTopMargin
	dc.SetColor(color.White)
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin
	dc.DrawStringWrapped(location, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	if err := dc.LoadFontFace(fontPathLight, 35); err != nil {
		return fmt.Errorf("error while loading font %v", err)
	}

	dc.DrawStringWrapped(d, x, y+80, 0, 0, maxWidth, 1.5, gg.AlignLeft)
	dc.DrawStringWrapped(desc, x, y+150, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	if err := dc.LoadFontFace(fontPath, 250); err != nil {
		return fmt.Errorf("error while loading font %v", err)
	}

	dc.DrawStringWrapped(strconv.Itoa(temp), x-70, y, 0, 0, maxWidth, 1.5, gg.AlignRight)

	if err := dc.LoadFontFace(fontPathLight, 70); err != nil {
		return fmt.Errorf("error while loading font %v", err)
	}

	dc.DrawStringWrapped("°C", x, y, 0, 0, maxWidth, 1.5, gg.AlignRight)

	if err := dc.LoadFontFace(fontPathLight, 100); err != nil {
		return fmt.Errorf("error while loading font %v", err)
	}

	dc.DrawStringWrapped(g, x, y+300, 0, 0, maxWidth, 1.5, gg.AlignCenter)

	if err := dc.SavePNG("./assets/day.png"); err != nil {
		return fmt.Errorf("error while saving background image %v", err)
	}

	fmt.Printf("got weather info %s, %s", desc, strconv.Itoa(temp))

	return err
}

func greetTime(t time.Time) string {
	var g string
	switch {
	case t.Hour() < 12:
		g = "Good morning"
	case t.Hour() < 17:
		g = "Good afternoon"
	default:
		g = "Good evening"
	}

	return g
}
