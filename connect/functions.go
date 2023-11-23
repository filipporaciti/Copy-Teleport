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
                out.Username = Username
                

                ris, err := json.Marshal(&out)
                if err != nil{
                        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
                }
                connection.Write([]byte(ris))
        }

        return tryToken == token
}

func checkPassword(connection net.Conn, tryPassword string) bool {
        //return true // da togliere
        // if tryPassword != Password {

        //         fmt.Println("[Error] invalid password")
        //         out := ResponseClient{}

        //         out.Type_request = "invalid password"
        //         out.Username = Username
                

        //         ris, err := json.Marshal(&out)
        //         if err != nil{
        //                 fmt.Println("[Error] json decoder: " + err.Error())
        //         }
        //         connection.Write([]byte(ris))
        // }

        return tryPassword == Password
}


func processResponse(connection net.Conn, res ResponseClient) error {

        switch res.Type_request {
        case "beacon request":
                out := ResponseClient{}
                out.Type_request = "beacon response"
                out.Username = Username

                ris, err := json.Marshal(&out)
                if err != nil{
                        fmt.Println("\033[31m[Error] json decoder:", err.Error(), "\033[0m")
                        return err
                }
                connection.Write([]byte(ris))
                return nil
        case "beacon response":

        	// plainUsername, err := cipher.LocalAESDecrypt([]byte(res.Username))
        	// if err != nil{
                //         fmt.Println("[Errore] local AES decrypt: " + err.Error())
                //         return err
                // }

        	// split to remove "port" part in ip address
                Add(res.Username, strings.Split(connection.RemoteAddr().String(), ":")[0])
                return nil
        case "token update":
                

                b64CipToken, err := cipher.Base64ToByte(res.Token)
                if err != nil {
                	return err
                }
                
                plainToken, err := cipher.LocalAESDecrypt(b64CipToken)
                if err != nil {
                	return err
                }

                token = string(plainToken)

                return err

        case "add copy":

        	plainToken, err := cipher.LocalAESDecrypt([]byte(res.Token))
        	if err != nil{
                        fmt.Println("\033[31m[Error] local AES decrypt:", err.Error(), "\033[0m")
                        return err
                }

        	plainData, err := cipher.LocalAESDecrypt([]byte(res.Data))
        	if err != nil{
                        fmt.Println("\033[31m[Error] local AES decrypt:", err.Error(), "\033[0m")
                        return err
                }

                if checkToken(connection, string(plainToken)) {
                        history.Add(res.Username, string(plainData))
                        PasteClipboard(string(plainData))
                }
                return nil

        case "update devices":

        	byteEncToken, _ := cipher.Base64ToByte(res.Token)
        	plainToken, err := cipher.LocalAESDecrypt(byteEncToken)
        	if err != nil{
                        fmt.Println("\033[31m[Error] local AES decrypt:", err.Error(), "\033[0m")
                        return err
                }

                byteEncData, _ := cipher.Base64ToByte(res.Data)
        	plainData, err := cipher.LocalAESDecrypt(byteEncData)
        	if err != nil{
                        fmt.Println("\033[31m[Error] local AES decrypt:", err.Error(), "\033[0m")
                        return err
                }
                

                if checkToken(connection, string(plainToken)) {
                        jsonOut := make([]devices.DevicesElement, 0)
                        json.Unmarshal(plainData, &jsonOut)
                        devices.Values = jsonOut
                }
                return nil
        case "get public key":

        	passFunc := func(conn net.Conn, pass string) bool {
	        	return pass == Password
        	}

        	err := cipher.ResponseAESKeyExchange(connection, passFunc)

        	if err != nil {
        		return err
        	}

        	SendTokenUpdate(strings.Split(connection.RemoteAddr().String(), ":")[0], "Pippo", "plainPassword")

        	return err
        default:
                fmt.Println("\033[31m[Error] No valid request", "\033[0m")
                return errors.New("No valid request")
        }
}
