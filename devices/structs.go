package devices

type DevicesElement struct {
	Username string
	Password string
	Ip_address string

}

var Values = make([]DevicesElement, 0)


func Add(username string, password string, ip_address string) {

	for i, _ := range Values {
		if Values[i].Ip_address == ip_address{
			Values = append(Values[:i], Values[i+1:]...)
		}
	}

	Values = append(Values, DevicesElement{
		Username: username,
		Password: password,
		Ip_address: ip_address,
	})

}