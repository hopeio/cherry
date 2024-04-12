package win

import (
	"github.com/gonutz/w32/v2"
	"time"
)

func TapKey(keys ...uint16) {
	var onInput []w32.INPUT

	for i := 0; i < len(keys); i++ {
		input := w32.KeyboardInput(w32.KEYBDINPUT{
			Vk: keys[i],
		})
		onInput = append(onInput, input)
	}

	w32.SendInput(onInput...)

	time.Sleep(time.Millisecond * 100)

	var offInput []w32.INPUT

	for i := len(keys) - 1; i >= 0; i-- {
		input := w32.KeyboardInput(w32.KEYBDINPUT{
			Vk:    keys[i],
			Flags: w32.KEYEVENTF_KEYUP,
		})
		offInput = append(offInput, input)
	}

	w32.SendInput(offInput...)
}
