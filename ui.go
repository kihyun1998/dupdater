package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UpdaterUI struct {
	window      fyne.Window
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
	stepLabels  []*canvas.Text
	stepIcons   []*widget.Icon
	currentStep int
	totalSteps  int
}

func newUpdaterUI(appName string) *UpdaterUI {
	myApp := app.New()
	myWindow := myApp.NewWindow(fmt.Sprintf("%s Updater", appName))

	ui := &UpdaterUI{
		window:      myWindow,
		progressBar: widget.NewProgressBar(),
		statusLabel: widget.NewLabel("Initializing..."),
		totalSteps:  7,
	}

	ui.createUI()
	return ui
}

func (ui *UpdaterUI) createUI() {
	stepsContainer := ui.createStepsUI()

	content := container.NewVBox(
		widget.NewLabel("Updating application..."),
		ui.progressBar,
		ui.statusLabel,
		layout.NewSpacer(),
		stepsContainer,
	)

	paddedContent := container.NewPadded(content)

	ui.window.SetContent(paddedContent)
	ui.window.Resize(fyne.NewSize(400, 300))
}

func (ui *UpdaterUI) createStepsUI() *fyne.Container {
	stepsContainer := container.NewVBox()

	stepNames := []string{
		"Check application status",
		"Backup files",
		"Download update",
		"Verify download",
		"Extract files",
		"Verify installation",
		"Launch application",
	}

	for _, name := range stepNames {
		icon := widget.NewIcon(theme.QuestionIcon())
		label := canvas.NewText(name, color.Black)
		label.TextStyle = fyne.TextStyle{Monospace: true}

		stepContainer := container.NewHBox(icon, label)
		stepsContainer.Add(stepContainer)

		ui.stepIcons = append(ui.stepIcons, icon)
		ui.stepLabels = append(ui.stepLabels, label)

	}
	return stepsContainer
}

func (ui *UpdaterUI) UpdateStatus(status string) {
	ui.statusLabel.SetText(status)
}

func (ui *UpdaterUI) SetCurrentStep(step int) {
	if step < 0 || step >= ui.totalSteps {
		return
	}

	if ui.currentStep < ui.totalSteps {
		ui.stepIcons[ui.currentStep].SetResource(theme.QuestionIcon())
		ui.stepLabels[ui.currentStep].Color = color.Black
	}

	ui.currentStep = step
	ui.stepIcons[ui.currentStep].SetResource(theme.NavigateNextIcon())
	ui.stepLabels[ui.currentStep].Color = theme.Color(theme.ColorNamePrimary)

	for i := 0; i < ui.currentStep; i++ {
		ui.stepIcons[i].SetResource(theme.ConfirmIcon())
	}

	ui.window.Content().Refresh()
}

func (ui *UpdaterUI) ShowError(err error) {
	errorText := canvas.NewText(fmt.Sprintf("Error: %v", err), color.NRGBA{R: 200, G: 30, B: 30, A: 255})
	errorText.Alignment = fyne.TextAlignCenter
	errorText.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewVBox(
		errorText,
		widget.NewButton("OK", func() {
			ui.window.Close()
		}),
	)

	ui.window.SetContent(content)
}

func (ui *UpdaterUI) Run() {
	ui.window.ShowAndRun()
}
