//go:build linux

package application

import "C"

const AlertStyleWarning = C.int(0)
const AlertStyleInformational = C.int(1)
const AlertStyleCritical = C.int(2)

var alertTypeMap = map[DialogType]C.int{
	WarningDialog:  AlertStyleWarning,
	InfoDialog:     AlertStyleInformational,
	ErrorDialog:    AlertStyleCritical,
	QuestionDialog: AlertStyleInformational,
}

func (m *linuxApp) showAboutDialog(title string, message string, icon []byte) {
	//	var iconData unsafe.Pointer
	// if icon != nil {
	// 	iconData = unsafe.Pointer(&icon[0])
	// }
	//C.showAboutBox(C.CString(title), C.CString(message), iconData, C.int(len(icon)))
}

type linuxDialog struct {
	dialog *MessageDialog

	//nsDialog unsafe.Pointer
}

func (m *linuxDialog) show() {
	globalApplication.dispatchOnMainThread(func() {

		// Mac can only have 4 Buttons on a dialog
		if len(m.dialog.Buttons) > 4 {
			m.dialog.Buttons = m.dialog.Buttons[:4]
		}

		// if m.nsDialog != nil {
		// 	//C.releaseDialog(m.nsDialog)
		// }
		// var title *C.char
		// if m.dialog.Title != "" {
		// 	title = C.CString(m.dialog.Title)
		// }
		// var message *C.char
		// if m.dialog.Message != "" {
		// 	message = C.CString(m.dialog.Message)
		// }
		// var iconData unsafe.Pointer
		// var iconLength C.int
		// if m.dialog.Icon != nil {
		// 	iconData = unsafe.Pointer(&m.dialog.Icon[0])
		// 	iconLength = C.int(len(m.dialog.Icon))
		// } else {
		// 	// if it's an error, use the application Icon
		// 	if m.dialog.DialogType == ErrorDialog {
		// 		iconData = unsafe.Pointer(&globalApplication.options.Icon[0])
		// 		iconLength = C.int(len(globalApplication.options.Icon))
		// 	}
		// }

		// alertType, ok := alertTypeMap[m.dialog.DialogType]
		// if !ok {
		// 	alertType = AlertStyleInformational
		// }

		//		m.nsDialog = C.createAlert(alertType, title, message, iconData, iconLength)

		// Reverse the Buttons so that the default is on the right
		reversedButtons := make([]*Button, len(m.dialog.Buttons))
		var count = 0
		for i := len(m.dialog.Buttons) - 1; i >= 0; i-- {
			//button := m.dialog.Buttons[i]
			//C.alertAddButton(m.nsDialog, C.CString(button.Label), C.bool(button.IsDefault), C.bool(button.IsCancel))
			reversedButtons[count] = m.dialog.Buttons[i]
			count++
		}

		buttonPressed := int(0) //C.dialogRunModal(m.nsDialog))
		if len(m.dialog.Buttons) > buttonPressed {
			button := reversedButtons[buttonPressed]
			if button.callback != nil {
				button.callback()
			}
		}
	})

}

func newDialogImpl(d *MessageDialog) *linuxDialog {
	return &linuxDialog{
		dialog: d,
	}
}

type linuxOpenFileDialog struct {
	dialog *OpenFileDialog
}

func newOpenFileDialogImpl(d *OpenFileDialog) *linuxOpenFileDialog {
	return &linuxOpenFileDialog{
		dialog: d,
	}
}

func toCString(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

func (m *linuxOpenFileDialog) show() ([]string, error) {
	openFileResponses[m.dialog.id] = make(chan string)
	//	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		//nsWindow = m.dialog.window.impl.(*macosWebviewWindow).nsWindow
	}

	// Massage filter patterns into macOS format
	// We iterate all filter patterns, tidy them up and then join them with a semicolon
	// This should produce a single string of extensions like "png;jpg;gif"
	// 	var filterPatterns string
	// if len(m.dialog.filters) > 0 {
	// 	var allPatterns []string
	// 	for _, filter := range m.dialog.filters {
	// 		patternComponents := strings.Split(filter.Pattern, ";")
	// 		for i, component := range patternComponents {
	// 			filterPattern := strings.TrimSpace(component)
	// 			filterPattern = strings.TrimPrefix(filterPattern, "*.")
	// 			patternComponents[i] = filterPattern
	// 		}
	// 		allPatterns = append(allPatterns, strings.Join(patternComponents, ";"))
	// 	}
	// 	filterPatterns = strings.Join(allPatterns, ";")
	// }

	// C.showOpenFileDialog(C.uint(m.dialog.id),
	// 	C.bool(m.dialog.canChooseFiles),
	// 	C.bool(m.dialog.canChooseDirectories),
	// 	C.bool(m.dialog.canCreateDirectories),
	// 	C.bool(m.dialog.showHiddenFiles),
	// 	C.bool(m.dialog.allowsMultipleSelection),
	// 	C.bool(m.dialog.resolvesAliases),
	// 	C.bool(m.dialog.hideExtension),
	// 	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	// 	C.bool(m.dialog.allowsOtherFileTypes),
	// 	toCString(filterPatterns),
	// 	C.uint(len(filterPatterns)),
	// 	toCString(m.dialog.message),
	// 	toCString(m.dialog.directory),
	// 	toCString(m.dialog.buttonText),
	// 	nsWindow)
	var result []string
	for filename := range openFileResponses[m.dialog.id] {
		result = append(result, filename)
	}
	return result, nil
}

//export openFileDialogCallback
func openFileDialogCallback(cid C.uint, cpath *C.char) {
	path := C.GoString(cpath)
	id := uint(cid)
	channel, ok := openFileResponses[id]
	if ok {
		channel <- path
	} else {
		panic("No channel found for open file dialog")
	}
}

//export openFileDialogCallbackEnd
func openFileDialogCallbackEnd(cid C.uint) {
	id := uint(cid)
	channel, ok := openFileResponses[id]
	if ok {
		close(channel)
		delete(openFileResponses, id)
		freeDialogID(id)
	} else {
		panic("No channel found for open file dialog")
	}
}

type linuxSaveFileDialog struct {
	dialog *SaveFileDialog
}

func newSaveFileDialogImpl(d *SaveFileDialog) *linuxSaveFileDialog {
	return &linuxSaveFileDialog{
		dialog: d,
	}
}

func (m *linuxSaveFileDialog) show() (string, error) {
	saveFileResponses[m.dialog.id] = make(chan string)
	//	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		//		nsWindow = m.dialog.window.impl.(*linuxWebviewWindow).nsWindow
	}

	// C.showSaveFileDialog(C.uint(m.dialog.id),
	// 	C.bool(m.dialog.canCreateDirectories),
	// 	C.bool(m.dialog.showHiddenFiles),
	// 	C.bool(m.dialog.canSelectHiddenExtension),
	// 	C.bool(m.dialog.hideExtension),
	// 	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	// 	C.bool(m.dialog.allowOtherFileTypes),
	// 	toCString(m.dialog.message),
	// 	toCString(m.dialog.directory),
	// 	toCString(m.dialog.buttonText),
	// 	toCString(m.dialog.filename),
	// 	nsWindow)
	return <-saveFileResponses[m.dialog.id], nil
}

//export saveFileDialogCallback
func saveFileDialogCallback(cid C.uint, cpath *C.char) {
	// Covert the path to a string
	path := C.GoString(cpath)
	id := uint(cid)
	// put response on channel
	channel, ok := saveFileResponses[id]
	if ok {
		channel <- path
		close(channel)
		delete(saveFileResponses, id)
		freeDialogID(id)

	} else {
		panic("No channel found for save file dialog")
	}
}
