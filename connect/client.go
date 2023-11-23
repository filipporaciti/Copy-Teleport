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

func SendConnectionRequest(ip_address string, password string) error {

    // cipPassword, err := cipher.LocalAESEncrypt([]byte(password))
    // data.Password = string(cipPassword)

    err := cipher.RequestAESKeyExchange(ip_address, password)
    return err
}

func SendAddCopyRequest(text string) error {
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
        fmt.Println("\033[31m[Errore] codifica json on SendAddCopyRequest: " + err.Error(), "\033[0m")
        return err
    }

    err = nil

    for _, val := range devices.Values {
        err = SendData(val.Ip_address, ris, false)
    }
    return err
}

func SendOneBeaconRequest(ip_address string) error {
	data := new(ResponseClient)
	data.Type_request = "beacon request"
	data.Username = Username

	ris, err := json.Marshal(&data)

	if err != nil {
		fmt.Println("\033[31m[Errore] codifica json on SendOneBeaconRequest: " + err.Error(), "\033[0m")
        return err
	}
    err = SendData(ip_address, ris, true)
	return err
}

func SendTokenUpdate(ip_address string, user string, pass string) error {


    data := new(ResponseClient)
    data.Type_request = "token update"
    data.Username = Username
 
    cipToken, err := cipher.LocalAESEncrypt([]byte(token))
    if err != nil {
        return err
    }
    b64CipToken := cipher.ByteToBase64(cipToken)
    
    data.Token = b64CipToken

    ris, err := json.Marshal(&data)

    if err != nil {
        fmt.Println("\033[31m[Errore] codifica json on SendTokenUpdate: " + err.Error(), "\033[0m")
        return err
    }
    err = SendData(ip_address, ris, false)

                
    // plainPassword, err := cipher.LocalAESEncrypt([]byte(res.Password))
    // if err != nil{
        //         fmt.Println("[Errore] local AES encrypt: " + err.Error())
        //         return false
        // }

    devices.Add(user, pass, ip_address)
    SendUpdateDevices()

    return err
}

func SendData(ip_address string, data []byte, response bool) error {
	
	fmt.Println("Connecting to " + SERVER_TYPE + " server " + ip_address + ":" + SERVER_PORT)

    conn, err := net.DialTimeout(SERVER_TYPE, ip_address+":"+SERVER_PORT, time.Millisecond * 500)
    if err != nil {
        fmt.Println("\033[31m[Error creating connect send] ", err.Error(), "\033[0m")
        return err
    }
    defer conn.Close() // send data and stop connection

    fmt.Println("Send: ", string(data))
    _, err = conn.Write(data) 
    if err != nil {
        fmt.Println("\033[31m[Error send data]" + err.Error(), "\033[0m")
        return err
    }

    if response {
        buffer := make([]byte, 4096)
        mLen, err := conn.Read(buffer)
        if err != nil {
                fmt.Println("\033[31m[Error reading after send] ", err.Error(), "\033[0m")
                return err
        }

        out := strings.Trim(string(buffer[:mLen]), "\n")

        fmt.Println("Received after send: ", out)

        ris := ResponseClient{}
        json.Unmarshal([]byte(out), &ris)
        processResponse(conn, ris)

    }

    return err

}
