package connect


var Values = make([]AvailableDevice, 0)

type AvailableDevice struct {
	Username string
	Ip_address string
}

func Add(username string, ip_address string) {

        Values = append(Values, AvailableDevice{
                Username: username,
                Ip_address: ip_address,
        })

}

func (d *AvailableDevice) Connect() {

}


type ResponseClient struct {

        Type_request string     `json:"type_request"`
        B64EncData string          `json:"b64encdata"`
        
}

type DataResponse struct {
        Username string         `json:"username"`
        Password string         `json:"password"`
        Token string            `json:"token"`
        Errore string           `json:"errore"`

        Data string             `json:"data"`
        Key string              `json:"key"`

}