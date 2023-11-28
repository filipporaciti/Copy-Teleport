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
                        fmt.Println("\033[31m[Error] listening:", err.Error(), "\033[0m")
                        os.Exit(1)
                }
                defer server.Close()
                fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
                fmt.Println("Waiting for client...")
                for {
                        connection, err := server.Accept()

                        if err != nil {
                                fmt.Println("\033[31m[Error] server accept connection:", err.Error(), "\033[0m")
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
                if err != nil{
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


func SendUpdateDevices() error {
        fmt.Println("[Info] send update devices")
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
                        fmt.Println("\033[31m[Errore] codifica json on SendUpdateDevices (Values): " + err.Error(), "\033[0m")
                        return err
                }

                out := ResponseClient{}
                out.Type_request = "update devices"

                data := DataResponse{}
                data.Username = Username
                data.Token = token
                data.Data = cipher.ByteToBase64(dev)

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
                        fmt.Println("\033[31m[Errore] codifica json on SendUpdateDevices (output): " + err.Error(), "\033[0m")
                        return err
                }

                

                err = SendData(val.Ip_address, ris, false)
        }

    return err
}

