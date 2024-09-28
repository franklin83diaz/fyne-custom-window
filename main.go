package main

import (
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
	OnClickHold func()
}

func NewCustomTitleBar() *CustomTitleBar {
	titleBar := &CustomTitleBar{}
	titleBar.ExtendBaseWidget(titleBar)
	return titleBar
}

func (t *CustomTitleBar) CreateRenderer() fyne.WidgetRenderer {
	// Crear un botón de cierre personalizado (sin bordes)
	closeButton := widget.NewButton("", func() {

	})
	closeButton.Importance = widget.LowImportance // Estilo flat, sin resaltado
	closeButton.SetIcon(theme.CancelIcon())       // Usar un ícono "X"

	// Crear la barra de título personalizada
	title := widget.NewLabel("Mi Aplicación")

	// Crear la caja contenedora
	bar := container.NewHBox(
		title,
		layout.NewSpacer(),
		closeButton,
	)

	return widget.NewSimpleRenderer(bar)
}

// MouseDown se llama cuando se presiona el mouse sobre la barra de título
func (t *CustomTitleBar) MouseDown(ev *desktop.MouseEvent) {
	go func() {
		time.Sleep(1 * time.Second) // Esperar un segundo para detectar si es un clic sostenido
		if t.OnClickHold != nil {
			t.OnClickHold()
		}
	}()
}

func main() {
	myApp := app.New()

	if drv, ok := myApp.Driver().(desktop.Driver); ok {
		splashWindow := drv.CreateSplashWindow()

		// Crear un botón de cierre personalizado (sin bordes)
		closeButton := widget.NewButton("", func() {
			splashWindow.Close()
		})
		closeButton.Importance = widget.LowImportance // Estilo flat, sin resaltado
		closeButton.SetIcon(theme.CancelIcon())       // Usar un ícono "X"

		// Crear la barra de título personalizada
		titleBar := container.NewHBox(
			widget.NewLabel("Mi Aplicación"),
			layout.NewSpacer(),
			closeButton,
		)

		// Contenido principal de la ventana
		content := widget.NewLabel("Contenido principal")

		// Agregar la barra de título y contenido a la ventana
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
