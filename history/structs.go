package history


type HistoryElement struct {
	Username string
	Value string
}

var Values = make([]HistoryElement, 0)


func Add(username string, value string) {

	Values = append(Values, HistoryElement{
		Username: username,
		Value: value,
	})

}

func Get(index int) HistoryElement {

	return Values[index]

}

// func (e HistoryElement) Get() []HistoryElement {

// 	return 

// }

