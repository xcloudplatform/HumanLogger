package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/kbinani/screenshot"

	"log"
	"math/rand"

	pb "github.com/ClickerAI/ClickerAI/src/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := pb.NewRobotGoClient(conn)

	a := app.New()
	w := a.NewWindow("Widgets Demo")

	moveWidget := newCoordsWidget(func(x, y int) {
		fmt.Printf("Move Widget Coordinates: x=%d, y=%d", x, y)
		client.Move(context.Background(), &pb.MoveRequest{X: int32(x), Y: int32(y)})
	})

	moveSmoothWidget := newCoordsWidget(func(x, y int) {
		fmt.Printf("Move Smooth Widget Coordinates: x=%d, y=%d\n", x, y)
		client.MoveSmooth(context.Background(), &pb.MoveSmoothRequest{X: int32(x), Y: int32(y)})

	})

	moveAndClickWidget := newCoordsWidget(func(x, y int) {
		fmt.Printf("Move and Click Widget Coordinates: x=%d, y=%d\n", x, y)
		client.MoveSmooth(context.Background(), &pb.MoveSmoothRequest{X: int32(x), Y: int32(y)})
		client.Sleep(context.Background(), &pb.SleepRequest{Tm: int32(1)})
		// client.Click(context.Background(), &pb.ClickRequest{
		// 	Args: []*anypb.Any{
		// 		any.Must(any.New("left")),
		// 	},
		// })

	})

	typeWidget := newTypeWidget(func(str string) {

	})

	grid := container.New(layout.NewFormLayout(),
		widget.NewLabel("Move"), moveWidget,
		widget.NewLabel("Move Smooth"), moveSmoothWidget,
		widget.NewLabel("Move and Click"), moveAndClickWidget,
		widget.NewLabel("Type"), typeWidget,
	)

	w.SetContent(grid)
	w.Resize(fyne.NewSize(800, 800))
	w.ShowAndRun()
}

func newTypeWidget(callback func(string)) fyne.CanvasObject {
	input := widget.NewEntry()
	randButton := widget.NewButton("RAND", func() {
		input.SetText(randString(10))
	})
	goButton := widget.NewButton("GO!", func() {
		callback(input.Text)
	})
	buttons := container.New(layout.NewGridLayout(2), randButton, goButton)

	grid := container.New(layout.NewGridLayoutWithColumns(3),
		input,
		layout.NewSpacer(),
		buttons,
	)

	return grid
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func newCoordsWidget(callback func(int, int)) fyne.CanvasObject {
	screenWidth := float64(screenSize().Width)
	screenHeight := float64(screenSize().Height)

	x := screenWidth / 2
	y := screenHeight / 2
	xData := binding.BindFloat(&x)
	yData := binding.BindFloat(&y)

	xSlider := widget.NewSliderWithData(0, screenWidth, xData)
	ySlider := widget.NewSliderWithData(0, screenHeight, yData)
	xEntry := widget.NewEntryWithData(binding.FloatToString(xData))
	yEntry := widget.NewEntryWithData(binding.FloatToString(yData))

	xWidgets := container.NewVBox(xEntry, xSlider)
	yWidgets := container.NewVBox(yEntry, ySlider)

	randButton := widget.NewButton("RAND", func() {
		xData.Set(rand.Float64() * screenWidth)
		yData.Set(rand.Float64() * screenHeight)
	})
	goButton := widget.NewButton("GO!", func() {
		callback(int(x), int(y))
	})

	buttons := container.New(layout.NewGridLayout(2), randButton, goButton)

	return container.New(layout.NewGridLayout(3), xWidgets, yWidgets, buttons)
}

func screenSize() fyne.Size {
	if screenshot.NumActiveDisplays() > 0 {
		// #0 is the main monitor
		bounds := screenshot.GetDisplayBounds(0)
		return fyne.NewSize(float32(bounds.Dx()), float32(bounds.Dy()))
	}
	return fyne.NewSize(800, 800)
}
