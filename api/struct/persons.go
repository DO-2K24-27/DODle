package persons

type Person struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
	Type      string `json:"type"`
	Workplace string `json:"workplace"`
	Image     string `json:"image"`
	Hint      string `json:"hint"`
}

type Persons struct {
	Persons []Person `json:"persons"`
}
