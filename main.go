package main

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomTitleBar struct {
	widget.BaseWidget
	Title       string
	CloseButton *widget.Button
	OnClickHold func()

	holdMutex    sync.Mutex
	cancelHoldCh chan struct{}
}

func NewCustomTitleBar() *CustomTitleBar {
	titleBar := &CustomTitleBar{
		cancelHoldCh: make(chan struct{}),
	}
	titleBar.ExtendBaseWidget(titleBar)
	return titleBar
}

func (t *CustomTitleBar) CreateRenderer() fyne.WidgetRenderer {
	title := widget.NewLabel(t.Title)
	bar := container.NewHBox(
		title,
		layout.NewSpacer(),
		t.CloseButton,
	)
	return widget.NewSimpleRenderer(bar)
}

func (t *CustomTitleBar) MouseDown(ev *desktop.MouseEvent) {
	t.holdMutex.Lock()
	defer t.holdMutex.Unlock()

	// Reset previous hold detection
	t.cancelHoldCh = make(chan struct{})

	go func(cancelCh chan struct{}) {
		select {
		case <-time.After(1 * time.Second):
			if t.OnClickHold != nil {
				t.OnClickHold()
			}
		case <-cancelCh:
			// Canceled hold detection
		}
	}(t.cancelHoldCh)
}

func (t *CustomTitleBar) MouseUp(ev *desktop.MouseEvent) {
	t.holdMutex.Lock()
	defer t.holdMutex.Unlock()
	close(t.cancelHoldCh) // Cancel any ongoing hold detection
}

func (t *CustomTitleBar) MouseMoved(ev *desktop.MouseEvent) {
	// Optional: Handle mouse moved events if needed
}

func (t *CustomTitleBar) MouseOut() {
	// Optional: Handle mouse out events if needed
}

func (t *CustomTitleBar) MouseIn(ev *desktop.MouseEvent) {
	// Optional: Handle mouse in events if needed
}

var _ desktop.Hoverable = (*CustomTitleBar)(nil) // Asserting that CustomTitleBar implements Hoverable
var _ desktop.Mouseable = (*CustomTitleBar)(nil) // Asserting that CustomTitleBar implements Mouseable

func main() {
	myApp := app.New()

	if drv, ok := myApp.Driver().(desktop.Driver); ok {
		splashWindow := drv.CreateSplashWindow()
		splashWindow.SetMaster()

		closeButton := widget.NewButton("", func() {

			splashWindow.Close()
		})
		closeButton.Importance = widget.LowImportance
		closeButton.SetIcon(theme.CancelIcon())

		titleBar := NewCustomTitleBar()
		titleBar.Title = "My App"
		titleBar.OnClickHold = func() {
			fmt.Println("click hold")
		}
		titleBar.CloseButton = closeButton

		content := widget.NewLabel("Contenido principal")

		splashWindow.SetContent(container.NewVBox(
			titleBar,
			content,
		))

		splashWindow.ShowAndRun()
	} else {
		mainWindow := myApp.NewWindow("Ventana con bordes")
		mainWindow.SetContent(widget.NewLabel("El driver no es compatible con desktop."))
		mainWindow.ShowAndRun()
	}
}
