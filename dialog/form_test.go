package dialog

import (
	"errors"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"
	testify "github.com/stretchr/testify/assert"
)

func TestFormDialog_Control(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	fd := controlFormDialog(ch, test.NewWindow(nil))
	fd.Show()
	go test.Tap(fd.confirm)

	select {
	case result := <-ch:
		assert.True(result, "Control should allow confirmation with no validation constraints")
	case <-time.After(500 * time.Millisecond):
		assert.Fail("Should have received a confirmation by now")
	}
}

func TestFormDialog_InvalidCannotSubmit(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	fd := validatingFormDialog(ch, test.NewWindow(nil))
	fd.Show()

	assert.False(fd.win.Hidden)
	assert.True(fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")
	go test.Tap(fd.confirm)

	select {
	case <-ch:
		assert.Fail("Callback should not have ran with an invalid form")
	case <-time.After(500 * time.Millisecond):
	}
}

func TestFormDialog_ValidCanSubmit(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	fd := validatingFormDialog(ch, test.NewWindow(nil))
	fd.Show()

	assert.False(fd.win.Hidden)
	assert.True(fd.confirm.Disabled(), "Confirm button should be disabled due to validation state")

	if validatingEntry, ok := fd.items[0].Widget.(*widget.Entry); ok {
		validatingEntry.SetText("abc")
		assert.False(fd.confirm.Disabled())
		go test.Tap(fd.confirm)

		select {
		case result := <-ch:
			assert.True(result, "Confirm should return true result")
		case <-time.After(500 * time.Millisecond):
			assert.Fail("Callback should have ran with a valid form")
		}
	} else {
		assert.Fail("First item's widget should be an Entry (check validatingFormDialog)")
	}
}

func TestFormDialog_CanCancelInvalid(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	fd := validatingFormDialog(ch, test.NewWindow(nil))
	fd.Show()
	assert.False(fd.win.Hidden)

	go test.Tap(fd.dismiss)

	select {
	case result := <-ch:
		assert.False(result, "Result should be false with cancellation")
	case <-time.After(500 * time.Millisecond):
		assert.Fail("Should have received a cancellation by now")
	}
}

func TestFormDialog_CanCancelNoValidation(t *testing.T) {
	assert := testify.New(t)
	ch := make(chan bool)
	fd := controlFormDialog(ch, test.NewWindow(nil))
	fd.Show()
	assert.False(fd.win.Hidden)

	go test.Tap(fd.dismiss)

	select {
	case result := <-ch:
		assert.False(result, "Result should be false with cancellation")
	case <-time.After(500 * time.Millisecond):
		assert.Fail("Should have received a cancellation by now")
	}
}

func validatingFormDialog(ch chan bool, parent fyne.Window) *FormDialog {
	validatingEntry := widget.NewEntry()
	validatingEntry.Validator = func(input string) error {
		if input != "abc" {
			return errors.New("only accepts 'abc'")
		}
		return nil
	}
	validatingItem := &widget.FormItem{
		Text:   "Only accepts 'abc'",
		Widget: validatingEntry,
	}
	controlEntry := widget.NewPasswordEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}

	items := []*widget.FormItem{validatingItem, controlItem}
	return NewFormDialog("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		ch <- confirm
	}, parent)
}

func controlFormDialog(ch chan bool, parent fyne.Window) *FormDialog {
	controlEntry := widget.NewEntry()
	controlItem := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry,
	}
	controlEntry2 := widget.NewPasswordEntry()
	controlItem2 := &widget.FormItem{
		Text:   "I accept anything",
		Widget: controlEntry2,
	}
	items := []*widget.FormItem{controlItem, controlItem2}
	return NewFormDialog("Validating Form Dialog", "Submit", "Cancel", items, func(confirm bool) {
		ch <- confirm
	}, parent)
}
