package alfred

type Result struct {
	Items []Item `json:"items"`
}

type Item struct {
	UID      string  `json:"uid,omitempty"`
	Title    string  `json:"title"`
	Subtitle string  `json:"subtitle,omitempty"`
	Arg      string  `json:"arg,omitempty"`
	Action   *Action `json:"action,omitempty"`
	Text     string  `json:"text,omitempty"`

	Icon *Icon `json:"icon,omitempty"`

	Valid *bool  `json:"valid,omitempty"`
	Type  string `json:"type,omitempty"`

	Autocomplete string `json:"autocomplete,omitempty"`

	Mods map[string]Mod `json:"mods,omitempty"`
}

type Mod struct {
	Valid    *bool  `json:"valid,omitempty"`
	Arg      string `json:"arg"`
	Subtitle string `json:"subtitle"`
}

type Action struct {
	Text []string `json:"text,omitempty"`
	URL  string   `json:"url,omitempty"`
	File string   `json:"file,omitempty"`
	Auto string   `json:"auto,omitempty"`
}

type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path"`
}
