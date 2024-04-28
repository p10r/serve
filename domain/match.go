package domain

type UntrackedMatch struct {
	HomeName  string
	AwayName  string
	StartTime int64
	Country   string
	League    string
}

type Match struct {
	ID        int64
	HomeName  string
	AwayName  string
	StartTime int64
	Country   string
	League    string
}

type Matches []Match

func NewMatch(id int64, match UntrackedMatch) Match {
	return Match{
		ID:        id,
		HomeName:  match.HomeName,
		AwayName:  match.AwayName,
		StartTime: match.StartTime,
		Country:   match.Country,
		League:    match.League,
	}
}
