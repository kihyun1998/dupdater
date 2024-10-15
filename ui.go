package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UpdaterUI struct {
	window      fyne.Window
	stepLabel   *widget.Label
	detailLabel *widget.Label
	spinnerIcon *canvas.Image
	currentStep int
	totalSteps  int
}

func newUpdaterUI(appName string) *UpdaterUI {
	myApp := app.New()
	myWindow := myApp.NewWindow(fmt.Sprintf("%s Updater", appName))

	ui := &UpdaterUI{
		window:     myWindow,
		totalSteps: 7,
	}

	ui.createUI()
	return ui
}

func (ui *UpdaterUI) createUI() {
	icon := fyne.NewStaticResource("icon", resourceIconPng.StaticContent)

	ui.spinnerIcon = canvas.NewImageFromResource(icon)
	ui.spinnerIcon.FillMode = canvas.ImageFillContain
	// image.Resize(fyne.NewSize(200, 200))

	// ui.spinnerIcon = widget.NewIcon(icon)
	// ui.spinnerIcon.Resize(fyne.NewSize(50, 50))

	ui.stepLabel = widget.NewLabel("")
	ui.stepLabel.TextStyle = fyne.TextStyle{Bold: true}

	ui.detailLabel = widget.NewLabel("")
	ui.detailLabel.TextStyle = fyne.TextStyle{Italic: true}
	ui.detailLabel.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		ui.spinnerIcon,
		ui.stepLabel,
		ui.detailLabel,
	)

	paddedContent := container.NewPadded(content)

	ui.window.SetContent(paddedContent)
	ui.window.Resize(fyne.NewSize(300, 150))
}

func (ui *UpdaterUI) UpdateStep(status string) {
	ui.stepLabel.SetText(status)
	ui.window.Content().Refresh()
}

func (ui *UpdaterUI) UpdateDetail(status string) {
	ui.detailLabel.SetText(status)
	ui.window.Content().Refresh()
}

func (ui *UpdaterUI) SetCurrentStep(step int) {
	if step < 0 || step >= ui.totalSteps {
		return
	}

	ui.currentStep = step
	ui.UpdateStep(fmt.Sprintf("Step %d of %d", ui.currentStep+1, ui.totalSteps))
}

func (ui *UpdaterUI) ShowError(err error) {
	errorText := widget.NewLabel(fmt.Sprintf("Error: %v", err))
	restoreButton := widget.NewButton("Restore", func() {
		LogInfo("Restore button clicked")
		if err := restoreFiles(ui); err != nil {
			LogError("Restore failed: %v", err)
			ui.ShowError(fmt.Errorf("restore failed: %v", err))
		} else {
			ui.UpdateDetail("Restore completed successfully")
		}
	})

	content := container.NewVBox(
		errorText,
		restoreButton,
	)

	ui.window.SetContent(content)
}

func (ui *UpdaterUI) Run() {
	ui.window.ShowAndRun()
}
