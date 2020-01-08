package models

type DrawedAction struct {
	ID          int
	AwardID     int
	AwardName   string
	PeopleCount int
	LeaderCount int
	BackMoney   int
	Memo        string
	//ToDo Or Done
	Status string
}
