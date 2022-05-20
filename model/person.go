package model

import (
	"encoding/json"

	"calendar/sqlite"
	log "github.com/sirupsen/logrus"
)

var createTableCmd = `
CREATE TABLE person (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT NUL NULL COLLATE NOCASE,
	second_name TEXT NULL COLLATE NOCASE,
	last_name TEXT NUL NULL COLLATE NOCASE,
	one_day INTEGER NOT NULL DEFAULT 1,
	two_days INTEGER NOT NULL DEFAULT 1,
	one_week INTEGER NOT NULL DEFAULT 1,
	birthday TEXT DEFAULT NULL,
	nameday TEXT DEFAULT NULL,
	UNIQUE(first_name, second_name, last_name)
)`

func CreatePersonTable() (err error) {
	if err = sqlite.SQLite(nil).Exec(createTableCmd); err == nil {
		initContent()
	}
	return
}

func initContent() {
	query := "INSERT INTO person (first_name, second_name, last_name) VALUES (?, ?, ?)"
	_, _ = sqlite.SQLite(nil).InsertQuery(query, "Piotr", "Włodzimierz", "Pszczółkowski")
	_, _ = sqlite.SQLite(nil).InsertQuery(query, "Błażej", "", "Pszczółkowski")
	_, _ = sqlite.SQLite(nil).InsertQuery(query, "Artur", "Piotr", "Pszczółkowski")
}

type Person struct {
	Id         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	OneDay     int    `json:"one_day"`
	TwoDays    int    `json:"two_days"`
	OneWeek    int    `json:"one_week"`
	Birthday   string `json:"birthday"`
	Nameday    string `json:"nameday"`
}

func (p *Person) String() string {
	data, _ := json.MarshalIndent(p, "", "   ")
	return string(data)
}

func (p *Person) OneDayState() bool {
	if p.OneDay == 0 {
		return false
	}
	return true
}

func (p *Person) TwoDaysState() bool {
	if p.TwoDays == 0 {
		return false
	}
	return true
}

func (p *Person) OneWeekState() bool {
	if p.OneWeek == 0 {
		return false
	}
	return true
}

func (p *Person) SetOneDayState(state bool) {
	switch state {
	case false:
		p.OneDay = 0
	case true:
		p.OneDay = 1
	}
}

func (p *Person) SetTwoDaysState(state bool) {
	switch state {
	case false:
		p.TwoDays = 0
	case true:
		p.TwoDays = 1
	}
}

func (p *Person) SetOneWeekState(state bool) {
	switch state {
	case false:
		p.OneWeek = 0
	case true:
		p.OneWeek = 1
	}
}

func NewPerson(firstName, secondName, lastName string) *Person {
	return &Person{
		FirstName:  firstName,
		SecondName: secondName,
		LastName:   lastName,
	}
}

func NewPersonWithId(id int64, firstName, secondName, lastName string) *Person {
	return &Person{
		Id:         id,
		FirstName:  firstName,
		SecondName: secondName,
		LastName:   lastName,
	}
}

func NewPersonForData(data map[string]any) *Person {
	person := new(Person)

	if value, ok := data[`id`]; ok {
		if v, ok := value.(int64); ok {
			person.Id = v
		}
	}
	if value, ok := data[`first_name`]; ok {
		if v, ok := value.(string); ok {
			person.FirstName = v
		}
	}
	if value, ok := data[`second_name`]; ok {
		if v, ok := value.(string); ok {
			person.SecondName = v
		}
	}
	if value, ok := data[`last_name`]; ok {
		if v, ok := value.(string); ok {
			person.LastName = v
		}
	}
	person.OneDay = 1
	if value, ok := data["one_day"]; ok {
		if v, ok := value.(int64); ok {
			if v == 0 {
				person.OneDay = 0
			}
		}
	}
	person.TwoDays = 1
	if value, ok := data["two_days"]; ok {
		if v, ok := value.(int64); ok {
			if v == 0 {
				person.TwoDays = 0
			}
		}
	}
	person.OneWeek = 1
	if value, ok := data["one_week"]; ok {
		if v, ok := value.(int64); ok {
			if v == 0 {
				person.OneWeek = 0
			}
		}
	}

	if value, ok := data["birthday"]; ok {
		if v, ok := value.(string); ok {
			person.Birthday = v
		}
	}
	if value, ok := data["nameday"]; ok {
		if v, ok := value.(string); ok {
			person.Nameday = v
		}
	}

	return person
}

func (p *Person) Save() bool {
	if p.Id == 0 {
		return p.Insert()
	}
	return p.Update()
}

func (p *Person) Update() bool {
	data := sqlite.Row{
		"first_name":  p.FirstName,
		"second_name": p.SecondName,
		"last_name":   p.LastName,
		"one_day":     p.OneDay,
		"two_days":    p.TwoDays,
		"one_week":    p.OneWeek,
		"birthday":    p.Birthday,
		"nameday":     p.Nameday,
	}
	where := sqlite.Where{
		"id": p.Id,
	}
	err := sqlite.SQLite(nil).Update(`person`, data, where)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (p *Person) Insert() bool {
	data := map[string]any{
		"first_name":  p.FirstName,
		"second_name": p.SecondName,
		"last_name":   p.LastName,
		"one_day":     p.OneDay,
		"two_days":    p.TwoDays,
		"one_week":    p.OneWeek,
		"birthday":    p.Birthday,
		"nameday":     p.Nameday,
	}
	id, err := sqlite.SQLite(nil).Insert("person", data)
	if err != nil {
		log.Error(err)
		return false
	}
	p.Id = id
	return true
}

func (p *Person) Delete() bool {
	return DeleteWithId(p.Id)
}

func DeleteWithId(id int64) bool {
	err := sqlite.SQLite(nil).ExecQuery("DELETE FROM person WHERE id=?", id)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func PersonWithId(id int64) *Person {
	data, err := sqlite.SQLite(nil).SelectQuery("SELECT * FROM person WHERE id=?", id)
	if err != nil {
		log.Error(err)
		return nil
	}
	return NewPersonForData(data[0])
}
func AllPersons() []*Person {
	data, err := sqlite.SQLite(nil).SelectQuery("SELECT * FROM person ORDER BY first_name, last_name")
	if err != nil {
		log.Error(err)
		return nil
	}
	retv := make([]*Person, 0, len(data))
	for _, item := range data {
		retv = append(retv, NewPersonForData(item))
	}
	return retv
}

func PersonCount() int {
	data, err := sqlite.SQLite(nil).SelectQuery("SELECT COUNT(*) as count FROM person")
	if err != nil {
		panic(err)
	}
	if value, ok := data[0]["count"]; ok {
		if v, ok := value.(int64); ok {
			return int(v)
		}
	}
	return 0
}
