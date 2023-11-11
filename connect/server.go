package connect

import (
	"net"
	"os"
        "fmt"
        "strings"
        "encoding/json"

        "Copy-Teleport/devices"
        "Copy-Teleport/cipher"

)

var server, err = net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

func ServerStart(){

        SetRandomToken()
        StartCopyClipboardDaemon()

        go func() {
        	fmt.Println("Server Running...")
                
                if err != nil {
                        fmt.Println("Error listening:", err.Error())
                        os.Exit(1)
                }
                defer server.Close()
                fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
                fmt.Println("Waiting for client...")
                for {
                        connection, err := server.Accept()

                        if err != nil {
                                fmt.Println("[Error] ", err.Error())
                                continue
                        }
                        fmt.Println("client connected")
                        go processClient(connection)   
                } 
        }()   
}


func processClient(connection net.Conn) {
        buffer := make([]byte, 4096)
        for {
                mLen, err := connection.Read(buffer)
                if err != nil {
                        fmt.Println("Error reading:", err.Error())
                        return
                }
                out := strings.Trim(string(buffer[:mLen]), "\n")

                fmt.Println("Received: ", out)

                if mLen != 1 {

                        ris := ResponseClient{}
                        json.Unmarshal([]byte(out), &ris)
                        processResponse(connection, ris)
                }

        }
        // connection.Close()
}


func SendUpdateDevices() (bool, error) {
        out := true
        e := error(nil)
        for _, val := range devices.Values {

                v := devices.Values 
                for i, _ := range v {
                        if v[i].Ip_address == SERVER_HOST{
                                v = append(v[:i], v[i+1:]...)
                        }
                }
                v = append(v, devices.DevicesElement{
                        Username: Username,
                        Password: Password,
                        Ip_address: SERVER_HOST,
                })

                index := 0
                for i, _ := range v{
                        if v[i].Ip_address == val.Ip_address{
                                index = i
                        }
                }
                v = append(v[:index], v[index+1:]...)
                dev, err := json.Marshal(&v)

                if err != nil {
                        fmt.Println("[Errore] codifica json on SendUpdateDevices (Values): " + err.Error())
                        return false, err
                }

                data := new(ResponseClient)
                data.Type_request = "update devices"
                cipUsername, err := cipher.LocalAESEncrypt([]byte(Username))
                data.Username = string(cipUsername)
                cipToken, err := cipher.LocalAESEncrypt([]byte(token))
                data.Token = string(cipToken)
                cipDev, err := cipher.LocalAESEncrypt(dev)
                data.Data = string(cipDev)

                ris, err := json.Marshal(&data)

                if err != nil {
                        fmt.Println("[Errore] codifica json on SendUpdateDevices (output): " + err.Error())
                        return false, err
                }

                

                out, e = SendData(val.Ip_address, ris, false)
        }

    return out, e
}

