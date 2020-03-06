package structs

type Config struct {
	User      string `json:"user"`
	Password  string `json:"password"`
	StartHour string `json:"startHour"`
	EndHour   string `json:"endHour"`
	RealHours string `json:"realHours"`
}