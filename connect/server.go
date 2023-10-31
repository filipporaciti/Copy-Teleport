package connect

import (
	"net"
	"os"
        "fmt"
        "strings"
        "encoding/json"

        "Copy-Teleport/devices"

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

func checkPassword(connection net.Conn, res ResponseClient) bool {
        //return true // da togliere
        if res.Password != Password {

                fmt.Println("[Error] invalid password")
                out := ResponseClient{}

                out.Type_request = "invalid password"
                out.Username = Username
                

                ris, err := json.Marshal(&out)
                if err != nil{
                        fmt.Println("[Errore] json decoder: " + err.Error())
                }
                connection.Write([]byte(ris))
        }

        return res.Password == Password
}

func checkToken(connection net.Conn, res ResponseClient) bool {
        if res.Token != token {

                fmt.Println("[Error] invalid token")
                out := ResponseClient{}

                out.Type_request = "invalid token"
                out.Username = Username
                

                ris, err := json.Marshal(&out)
                if err != nil{
                        fmt.Println("[Errore] json decoder: " + err.Error())
                }
                connection.Write([]byte(ris))
        }

        return res.Token == token
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
                data.Username = Username
                data.Token = token
                data.Data = string(dev)

                ris, err := json.Marshal(&data)

                if err != nil {
                        fmt.Println("[Errore] codifica json on SendUpdateDevices (output): " + err.Error())
                        return false, err
                }

                

                out, e = SendData(val.Ip_address, string(ris), false)
        }

    return out, e
}

