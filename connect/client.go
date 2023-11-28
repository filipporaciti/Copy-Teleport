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

    err := cipher.RequestAESKeyExchange(ip_address, password)
    return err
}

func SendAddCopyRequest(text string) error {
    out := ResponseClient{}
    out.Type_request = "add copy"

    data := DataResponse{}
    data.Username = Username
    data.Token = token
    data.Data = text

    stringData, err := json.Marshal(&data)
    if err != nil{
            fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
            return err
    }
    
    cipData, err := cipher.LocalAESEncrypt([]byte(stringData))
    if err != nil{
            return err
    }
    b64CipData := cipher.ByteToBase64(cipData)

    out.B64EncData = b64CipData

    ris, err := json.Marshal(&out)

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
	out := ResponseClient{}
	out.Type_request = "beacon request"
	out.B64EncData = Username

	ris, err := json.Marshal(&out)

	if err != nil {
		fmt.Println("\033[31m[Errore] codifica json on SendOneBeaconRequest: " + err.Error(), "\033[0m")
        return err
	}
    err = SendData(ip_address, ris, true)
	return err
}

func SendTokenUpdate(ip_address string) error {


    out := ResponseClient{}
    out.Type_request = "token update request"


    data := DataResponse{}
    data.Username = Username
    data.Token = token
    data.Data = "speriamo che serva a qualcosa questa stringa"

    stringData, err := json.Marshal(&data)
    if err != nil{
        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
        return err
    }
    
    cipData, err := cipher.LocalAESEncrypt([]byte(stringData))
    if err != nil{
        return err
    }
    b64CipData := cipher.ByteToBase64(cipData)

    out.B64EncData = b64CipData



    ris, err := json.Marshal(&out)

    if err != nil {
        fmt.Println("\033[31m[Errore] codifica json on SendTokenUpdate: " + err.Error(), "\033[0m")
        return err
    }
    err = SendData(ip_address, ris, true)

    return err
    
}

func SendTokenUpdateResponse(ip_address string) error {
    out := ResponseClient{}
    out.Type_request = "token update response"

    data := DataResponse{}
    data.Username = Username
    data.Token = token

    stringData, err := json.Marshal(&data)
    if err != nil{
        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
        return err
    }

    cipData, err := cipher.LocalAESEncrypt([]byte(stringData))
    if err != nil{
        return err
    }
    b64CipData := cipher.ByteToBase64(cipData)

    out.B64EncData = b64CipData

    ris, err := json.Marshal(&out)

    if err != nil {
        fmt.Println("\033[31m[Errore] json encoding: " + err.Error(), "\033[0m")
        return err
    }

    err = SendData(ip_address, ris, false)
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
