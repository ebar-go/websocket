package utils

func Assert(guard bool, text string) {
	if !guard {
		panic(text)
	}
}
