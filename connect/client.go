package connect

import (
    "net"
	"time"
	"fmt"
	"strings"
    "encoding/json"

    "Copy-Teleport/devices"
)

func SendConnectionRequest(ip_address string, password string) (bool, error) {
    data := new(ResponseClient)
    data.Type_request = "connection request"
    data.Username = Username
    data.Password = password

    ris, err := json.Marshal(&data)

    if err != nil {
        fmt.Println("[Errore] codifica json on SendConnectionRequest: " + err.Error())
        return false, err
    }

    return SendData(ip_address, string(ris), true)
}

func SendAddCopyRequest(text string) (bool, error) {
    data := new(ResponseClient)
    data.Type_request = "add copy"
    data.Username = Username
    data.Token = token
    data.Data = text

    ris, err := json.Marshal(&data)

    if err != nil {
        fmt.Println("[Errore] codifica json on SendAddCopyRequest: " + err.Error())
        return false, err
    }

    out := true
    e := error(nil)
    
    // out, e = SendData("192.168.1.57", string(ris), false) // da togliere

    for _, val := range devices.Values {
        out, e = SendData(val.Ip_address, string(ris), false)
    }
    return out, e
}

func SendOneBeaconRequest(ip_address string) (bool, error) {
	data := new(ResponseClient)
	data.Type_request = "beacon request"
	data.Username = Username

	ris, err := json.Marshal(&data)

	if err != nil {
		fmt.Println("[Errore] codifica json on SendOneBeaconRequest: " + err.Error())
        return false, err
	}
	return SendData(ip_address, string(ris), true)
}

func SendData(ip_address string, data string, response bool) (bool, error) {
	
	fmt.Println("Connecting to " + SERVER_TYPE + " server " + ip_address + ":" + SERVER_PORT)

    conn, err := net.DialTimeout(SERVER_TYPE, ip_address+":"+SERVER_PORT, time.Millisecond * 500)
    if err != nil {
        fmt.Println("[Error creating connect send] ", err.Error())
        return false, err
    }
    defer conn.Close() // send data and stop connection

    fmt.Println("------" + data)
    _, err = conn.Write([]byte(data)) // prima cera uno \n
    if err != nil {
        fmt.Println("[Error send data]" + err.Error())
        return false, err
    }

    if response {
        buffer := make([]byte, 4096)
        mLen, err := conn.Read(buffer)
        if err != nil {
                fmt.Println("[Error reading after send] ", err.Error())
                return false, err
        }
        out := strings.Trim(string(buffer[:mLen]), "\n")

        fmt.Println("Received after send: ", out)

        ris := ResponseClient{}
        json.Unmarshal([]byte(out), &ris)
        processResponse(conn, ris)

    }

    return true, nil

}
