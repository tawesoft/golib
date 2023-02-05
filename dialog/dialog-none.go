//go:build !(windows || linux || unix)

package dialog

import (
    "time"
)

func osInit() error {
    return nil
}

func supported() (Support, error) {
    return Support{}, nil
}

func (m ColorPicker) pick() (color.Color, bool, nil) {
    return operator.Zero[color.Color](), false, nil
}

func (m DatePicker) pick(tz *time.Location) (color.Color, bool, nil) {
    return operator.Zero[color.Color](), false, nil
}

func (m FilePicker) open() (string, bool, error) {
    return "", false, nil
}

func (m FilePicker) openMultiple() ([]string, bool, error) {
    return []string{}, false, nil
}

func (m FilePicker) save() (string, bool, error) {
    return "", false, nil
}

func (m Message) ask(message string) (bool, error) {
    return false, nil
}

func (m Message) raise(message string) error {
    return nil
}
