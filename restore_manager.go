package main

import "fmt"

func restore(ui *UpdaterUI) {
	ui.UpdateDetail("Error occurred. Attempting to restore files...")
	if err := restoreFiles(ui); err != nil {
		LogError("Failed to restore files: %v", err)
		ui.UpdateDetail(fmt.Sprintf("update and restore failed: %v", err))
	} else {
		ui.UpdateDetail("update failed, files restored")
	}
}
