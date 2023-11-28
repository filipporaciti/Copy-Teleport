package connect

import (
	"fmt"
	"time"
	"crypto/rand"
	"math/big"
	"net"
	"encoding/json"
	"strings"
	"errors"

	"github.com/atotto/clipboard"

	"Copy-Teleport/devices"
	"Copy-Teleport/history"
	"Copy-Teleport/cipher"
)

var Username string
var Password string

var copy string // variable for StartCopyClipboardDaemon
var token string = ""
var SERVER_HOST string = GetIPAddress()
const (
    // SERVER_HOST = "localhost"
    SERVER_PORT = "20917"
    SERVER_TYPE = "tcp"
)

var maleNames = [...]string{"James", "Robert", "John", "Michael", "David", "William", "Richard", "Joseph", "Thomas", "Christopher"}
var femaleNames = [...]string{"Mary", "Patricia", "Jennifer", "Linda", "Elizabeth", "Barbara", "Susan", "Jessica", "Sarah", "Karen"}


// reset Values and call function GetDevices() to scan the network and get aviable devices.
// I use gorutine to avoid freeze in GUI
// Input:
// Output: bool (false if error else true), error (nil if no error)
func DiscoverDevices() (bool, error) {
	Values = make([]AvailableDevice, 0)
	var x = func() {
		GetDevices(SERVER_HOST)
	}
	go x()
	return true, nil // to do
}

// Check every 1000 milliseconds if the program have to update all clipboards (all devices) 
func StartCopyClipboardDaemon() {
	copy = CopyClipboard()
	next := ""

	var x = func() {
		for {
			next = CopyClipboard()
			// fmt.Println(copy, next)
			if next != copy {
				if err:= SendAddCopyRequest(next); err == nil {
					history.Add("me", next)
					copy = next
				}
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
	go x()
}


// Paste a string to operative system's clipboard
// Input: text (string to paste in operative system's clipboard)
// Output: 
func PasteClipboard(text string) {
	if err := clipboard.WriteAll(text); err != nil {
        	fmt.Println("\033[31m[Error] while copy to clipboard:", err.Error(), "\033[0m")
   	}
   	copy = text
}


// Return the last string in the clipboard
func CopyClipboard() string {
	text, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("\033[31m[Error] copy from clipboard:", err.Error(), "\033[0m")
	}
	return text
}


// Update the variable Username with a random one
func SetRandomUsername() {
	usr := ""

	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(maleNames))))
	if err != nil {
		fmt.Println("\033[31m[Error] while getting random username 1:", err.Error(), "\033[0m")
	}
	usr += maleNames[num.Int64()] + " "

	num, err = rand.Int(rand.Reader, big.NewInt(int64(len(maleNames))))
	if err != nil {
		fmt.Println("\033[31m[Error] while getting random username 2:", err.Error(), "\033[0m")
	}
	usr += femaleNames[num.Int64()]

	Username = usr
}


// Update the variable Password with a random one. You have to give the length 
func SetRandomPassword(n int) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i:=0; i<n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			fmt.Println("\033[31m[Error] while getting random password:", err.Error(), "\033[0m")
		}
		ret[i] = letters[num.Int64()]
	}

	Password = string(ret)
}

// Update the variable Token with a random one
func SetRandomToken(){

	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, 30)
	for i := 0; i < 30; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			fmt.Println("\033[31m[Error] while getting random token:", err.Error(), "\033[0m")
		}
		ret[i] = letters[num.Int64()]
	}
	token = string(ret)
}


func GetIPAddress() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        fmt.Println("\033[31m[Error] while getting ip address:", err.Error(), "\033[0m")
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP.String()
}

func checkToken(connection net.Conn, tryToken string) bool {
        if tryToken != token {

                fmt.Println("\033[31m[Error] invalid token:", err.Error(), "\033[0m")
                out := ResponseClient{}

                out.Type_request = "invalid token"

                data := DataResponse{}
                data.Username = Username

                stringData, err := json.Marshal(&data)
                if err != nil{
                        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
                        return false
                }
                
                cipData, err := cipher.LocalAESEncrypt([]byte(stringData))
                if err != nil{
                        return false
                }
                b64CipData := cipher.ByteToBase64(cipData)

                out.B64EncData = b64CipData

                ris, err := json.Marshal(&out)
                if err != nil{
                        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
                        return false
                }
                connection.Write([]byte(ris))
        }

        return tryToken == token
}

func processResponse(connection net.Conn, res ResponseClient) error {
	resData := DataResponse{}
	remote_addr := strings.Split(connection.RemoteAddr().String(), ":")[0]

	if res.Type_request != "beacon response" && res.Type_request != "beacon request" && res.Type_request != "get public key" {

		encResData, _ := cipher.Base64ToByte(res.B64EncData)
		stringResData, err := cipher.LocalAESDecrypt(encResData)
		if err != nil {
			return err
		}

	        json.Unmarshal([]byte(stringResData), &resData)
	}

        switch res.Type_request {
        case "beacon request":
                out := ResponseClient{}
                out.Type_request = "beacon response"
                out.B64EncData = Username

                ris, err := json.Marshal(&out)
                if err != nil{
                        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
                        return err
                }
                connection.Write([]byte(ris))
                return nil
        case "beacon response":
        	// split to remove "port" part in ip address
                Add(res.B64EncData, strings.Split(connection.RemoteAddr().String(), ":")[0])
                return nil
        case "token update request":

        	if resData.Data == "speriamo che serva a qualcosa questa stringa"{

	                token = resData.Token

	                err := SendTokenUpdateResponse(remote_addr)

			return err
		}
		return errors.New("Invalid token")
        case "token update response":

                if checkToken(connection, resData.Token) {
        		devices.Add(resData.Username, "resData.Password", remote_addr)
    			SendUpdateDevices()
    		}


                return err
        case "add copy":

                if checkToken(connection, resData.Token) {
                        history.Add(resData.Username, resData.Data)
                        PasteClipboard(resData.Data)
                }
                return nil

        case "update devices":

                if checkToken(connection, resData.Token) {
                        jsonOut := make([]devices.DevicesElement, 0)

                        byteDev, err := cipher.Base64ToByte(resData.Data)
                        if err != nil {
                        	return err
                        }
                        json.Unmarshal([]byte(byteDev), &jsonOut)
                        devices.Values = jsonOut
                }
                return nil
        case "get public key": // to server

        	passFunc := func(conn net.Conn, pass string) bool {
	        	return pass == Password
        	}

        	err := cipher.ResponseAESKeyExchange(connection, passFunc)

        	if err != nil {
        		return err
        	}


        	SendTokenUpdate(remote_addr)

        	return err
        default:
                fmt.Println("\033[31m[Error] No valid request", "\033[0m")
                return errors.New("No valid request")
        }
}
