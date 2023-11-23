package cipher

import(
    "bytes"
    "errors"
    "net"
    "encoding/json"
    "encoding/base64"
    "fmt"
    "time"
    "strings"

)

const (
    SERVER_PORT = "20917"
    SERVER_TYPE = "tcp"
)

type RequestResponse struct {
    Type_request string     `json:"type_request"`
    Key string              `json:"key"`
    Data string             `json:"data"`
}

type CheckPassword func(net.Conn, string) bool // Function to get check password function (check password from input to local account password)


func RequestAESKeyExchange(ip_address string, password string) error {

    fmt.Println("Connecting to " + SERVER_TYPE + " server " + ip_address + ":" + SERVER_PORT)

    conn, err := net.DialTimeout(SERVER_TYPE, ip_address+":"+SERVER_PORT, time.Millisecond * 2000)
    if err != nil {
        fmt.Println("\033[31m[Error] creating connect send:", err.Error(), "\033[0m")
        return err
    }
    defer conn.Close() // send data and stop connection


    SendGetPublicKey(conn)
    out, err := ReciveData(conn) // from SendPublicKey()
    if err != nil {
        return err
    }

    ris, err := DecodeRSAPublicKeyPEM(out.Key)
    if err != nil {
        return err
    }
    RemotePublicRSAKey = ris


    SendPasswordPublicKey(conn, password)
    out, err = ReciveData(conn) // from SendAESKey()
    if err != nil {
        return err
    }
    if out.Type_request == "wrong password" {
        return errors.New("Wrong password")
    }
    
    cipAESkey, err := Base64ToByte(out.Key)
    if err != nil {
        return err
    }

    AESRemoteKey, err := LocalRSADecrypt(cipAESkey)
    if err != nil {
        return err
    }

    SaveAESKey(AESRemoteKey)


    return err
}

func ResponseAESKeyExchange(conn net.Conn, cp CheckPassword) (error) {

    SendPublicKey(conn)
    out, err := ReciveData(conn) // from SendPasswordPublicKey()
    if err != nil {
        return err
    }

    ris, err := DecodeRSAPublicKeyPEM(out.Key)
    if err != nil {
        return err
    }
    RemotePublicRSAKey = ris

    password, err := Base64ToByte(out.Data)
    if err != nil {
        return err
    }
    plainPassword, _ := LocalRSADecrypt(password)

    SendAESKey(conn, string(plainPassword), cp)

    return err
}


func SaveAESKey(key []byte) {
    privateAESKey = key
}

func SendAESKey(conn net.Conn, password string, cp CheckPassword) error {
    var x = RequestResponse{}


    if cp(conn, password) {
        x.Type_request = "aes key"

        ris, _ := RSAEncrypt(RemotePublicRSAKey, privateAESKey)

        x.Key = ByteToBase64(ris)
    } else {
        x.Type_request = "wrong password"
    }


    data, err := json.Marshal(&x)
    if err != nil{
            fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
            return err
    }

    fmt.Println("\nSend: " + string(data))
    _, err = conn.Write([]byte(data)) 
    if err != nil {
        fmt.Println("\033[31m[Error] send data:", err.Error(), "\033[0m")
        return err
    }

    return err

}

func SendPasswordPublicKey(conn net.Conn, password string) error {
    var x = RequestResponse{}
    x.Type_request = "password public key"
    x.Key = EncodeRSAPublicKeyPEM(&privateRSAKey.PublicKey)

    res, err := RSAEncrypt(RemotePublicRSAKey, []byte(password))
    if err != nil {
        return err
    }
    x.Data = ByteToBase64(res)

    data, err := json.Marshal(&x)
    if err != nil{
            fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
            return err
    }

    
    fmt.Println("\nSend: " + string(data))
    _, err = conn.Write(data)
    if err != nil {
        fmt.Println("\033[31m[Error] send data:", err.Error(), "\033[0m")
        return err
    }

    return err
}

func SendPublicKey(conn net.Conn) error {
    var x = RequestResponse{}
    x.Type_request = "public key"
    x.Key = EncodeRSAPublicKeyPEM(&privateRSAKey.PublicKey)

    data, err := json.Marshal(&x)
    if err != nil{
            fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
            return err
    }

    
    fmt.Println("\nSend: " + string(data))
    _, err = conn.Write([]byte(data)) 
    if err != nil {
        fmt.Println("\033[31m[Error] send data:", err.Error(), "\033[0m")
        return err
    }

    return err
}


func SendGetPublicKey(conn net.Conn) error {
    var x = RequestResponse{}
    x.Type_request = "get public key"

    data, err := json.Marshal(&x)
    if err != nil{
            fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
            return err
    }

    
    fmt.Println("\nSend: " + string(data))
    _, err = conn.Write([]byte(data)) 
    if err != nil {
        fmt.Println("\033[31m[Error] send data:" + err.Error(), "\033[0m")
        return err
    }

    return err
}


func ReciveData(conn net.Conn) (RequestResponse, error) {
    buffer := make([]byte, 4096)
    mLen, err := conn.Read(buffer)
    if err != nil {
            fmt.Println("\033[31m[Error] reading after send:", err.Error(), "\033[0m")
            return RequestResponse{}, err
    }
    out := strings.Trim(string(buffer[:mLen]), "\n")

    fmt.Println("Received after send: ", out)

    ris := RequestResponse{}
    json.Unmarshal([]byte(out), &ris)

    return ris, err
}


// Add padding to input
//
// Input: src (input bytes), size (padding block size)
//
// Output: (src+padding)
func Pad(src []byte, size int) []byte {
    padding := size - len(src)%size
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
}


// Remove padding from src (with pad)
//
// Input: src (with pad)
//
// Output: src (without padding), error (nil if no error)
func Unpad(src []byte) ([]byte, error) {
    length := len(src)
    unpadding := int(src[length-1])

    if unpadding > length {
        return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
    }

    return src[:(length - unpadding)], nil
}


func ByteToBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

func Base64ToByte(data string) ([]byte, error) {
    ris, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        fmt.Println("\033[31m[Error] base64 decode:", err.Error(), "\033[0m")
    }
    return ris, err
}