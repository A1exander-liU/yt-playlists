package ui

type Keybindings struct {
	global    map[string]string
	main      map[string]string
	modal     map[string]string
	dialog    map[string]string
	playlists map[string]string
	videos    map[string]string
	form      map[string]string
}

func initKeys() Keybindings {
	return Keybindings{
		global:    globalKeys(),
		main:      mainKeys(),
		modal:     modalKeys(),
		dialog:    dialogKeys(),
		playlists: playlistsKeys(),
		videos:    videosKeys(),
		form:      formKeys(),
	}
}

func globalKeys() map[string]string {
	keys := make(map[string]string, 2)

	keys["Quit"] = "q"
	keys["Logout"] = "Q"
	keys["Help"] = "?"

	return keys
}

func mainKeys() map[string]string {
	keys := make(map[string]string, 1)

	keys["Switch"] = "Tab"

	return keys
}

func modalKeys() map[string]string {
	keys := make(map[string]string, 1)

	keys["Close"] = "Esc"

	return keys
}

func dialogKeys() map[string]string {
	keys := make(map[string]string, 2)

	keys["Switch"] = "Tab"
	keys["Select"] = "Space"

	return keys
}

func playlistsKeys() map[string]string {
	keys := make(map[string]string, 3)

	keys["Select"] = "Space"
	keys["New"] = "a"
	keys["Delete"] = "d"

	return keys
}

func videosKeys() map[string]string {
	keys := make(map[string]string, 4)

	keys["Select"] = "Space"
	keys["Add"] = "a"
	keys["Move"] = "m"
	keys["Delete"] = "d"
	keys["Open"] = "o"

	return keys
}

func formKeys() map[string]string {
	keys := make(map[string]string, 3)

	keys["Next"] = "Tab"
	keys["Prev"] = "Shift-Tab"
	keys["Select"] = "Space"

	return keys
}
