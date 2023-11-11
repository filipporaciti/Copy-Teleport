package connect

import (
    "net"
	"time"
	"fmt"
	"strings"
    "encoding/json"

    "Copy-Teleport/devices"
    "Copy-Teleport/cipher"
)

func SendConnectionRequest(ip_address string, password string) (bool, error) {
    data := new(ResponseClient)
    data.Type_request = "connection request"
    cipUsername, err := cipher.LocalAESEncrypt([]byte(Username))
    data.Username = string(cipUsername)
    cipPassword, err := cipher.LocalAESEncrypt([]byte(password))
    data.Password = string(cipPassword)

    ris, err := json.Marshal(&data)

    if err != nil {
        fmt.Println("[Errore] codifica json on SendConnectionRequest: " + err.Error())
        return false, err
    }

    return SendData(ip_address, ris, true)
}

func SendAddCopyRequest(text string) (bool, error) {
    data := new(ResponseClient)
    data.Type_request = "add copy"
    cipUsername, err := cipher.LocalAESEncrypt([]byte(Username))
    data.Username = string(cipUsername)
    cipToken, err := cipher.LocalAESEncrypt([]byte(token))
    data.Token = string(cipToken)
    cipText, err := cipher.LocalAESEncrypt([]byte(text))
    data.Data = string(cipText)

    ris, err := json.Marshal(&data)

    if err != nil {
        fmt.Println("[Errore] codifica json on SendAddCopyRequest: " + err.Error())
        return false, err
    }

    out := true
    err = error(nil)

    for _, val := range devices.Values {
        out, err = SendData(val.Ip_address, ris, false)
    }
    return out, err
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
	return SendData(ip_address, ris, true)
}

func SendData(ip_address string, data []byte, response bool) (bool, error) {
	
	fmt.Println("Connecting to " + SERVER_TYPE + " server " + ip_address + ":" + SERVER_PORT)

    conn, err := net.DialTimeout(SERVER_TYPE, ip_address+":"+SERVER_PORT, time.Millisecond * 500)
    if err != nil {
        fmt.Println("[Error creating connect send] ", err.Error())
        return false, err
    }
    defer conn.Close() // send data and stop connection

    fmt.Println("------" + string(data))
    _, err = conn.Write(data) 
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
