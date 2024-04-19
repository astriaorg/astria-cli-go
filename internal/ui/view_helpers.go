package ui

// Append the settings status to the end of the input string
func appendStatus(text string, status bool) string {
	if status {
		return text + ": [black:white]ON [-:-]"
	} else {
		return text + ": [white:darkslategray]off[-:-]"
	}
}
