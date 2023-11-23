package connect


import (
    "fmt"
    "net"
    "sync"
    "strings"
    "os/exec"
    "runtime"
    "os"
)

func GetDevices(targetIP string) {
    var wg sync.WaitGroup

    network := strings.Split(targetIP, ".")

    switch network[0] {
    case "192": // mask 255.255.255.0
    	for i := 1; i < 255; i++ {
	        ip := fmt.Sprintf("%s.%s.%s.%d", network[0], network[1], network[2], i)
	        // fmt.Printf(ip)
	        addr := net.ParseIP(ip)
	        if addr == nil {
	            fmt.Println("\033[31m[Error] Invalid IP address:", ip, "\033[0m")
	            continue
	        }
	        wg.Add(1)
	        x := func(wg *sync.WaitGroup) { 
    			defer wg.Done()
	        	if err := pingIP(ip); err == nil && ip != targetIP {
	            	SendOneBeaconRequest(ip)
	        		// go scanPort(ip, port)
	        	}
	        }
	        go x(&wg)
	    }
	case "172": // mask 255.255.0.0
    	for i := 1; i < 255; i++ {
    		for j := 1; j < 255; j++ {
		        ip := fmt.Sprintf("%s.%s.%d.%d", network[0], network[1], i, j)
		        // fmt.Printf(ip)
		        addr := net.ParseIP(ip)
		        if addr == nil {
		            fmt.Println("\033[31m[Error] Invalid IP address:", ip, "\033[0m")
		            continue
		        }
		        wg.Add(1)
		        x := func(wg *sync.WaitGroup) { 
	    			defer wg.Done()
		        	if err := pingIP(ip); err == nil && ip != targetIP {
		            	SendOneBeaconRequest(ip)
		        		// go scanPort(ip, port)
		        	}
		        }
		        go x(&wg)
		    }
	    }
    default:
    	return
    }
    wg.Wait()
    return 
}

// func scanPort(ip string, port int) {
//     // defer wg.Done()

//     target := fmt.Sprintf("%s:%d", ip, port)
//     conn, err := net.DialTimeout("tcp", target, 1*time.Second)
//     if err != nil {
//         return // Port is closed
//     }
//     conn.Close()
//     // fmt.Printf("Port %d is open on %s\n", port, ip)
// }


// Make ping request from terminal
func pingIP(ip string) error {
	var cmd *exec.Cmd

    command := fmt.Sprintf("ping -c 1 %s", ip)

    switch runtime.GOOS{
    case "windows": //Windows
    	fmt.Println("\033[31m[Error] Project can't work on Windows", "\033[0m")
        os.Exit(1)
    default://Mac & Linux
        cmd = exec.Command("bash", "-c", command)
    }

    return cmd.Run()
}
