// Server Script
package main
import ( 
  "bufio"
  "fmt"
  "net"
  "os" 
  "strings"
  "time"
)

func args() -> string { 
	arguments := os.Args
	if len(arguments) == 1 { 
		fmt.Println("please provide port number" ) 
		return "" 
	} 
	return arguments[1]

} 

func main() { 
  arguments := os.Args 
  if len(arguments) ==1  { 
    fmt.Println("Please provide port number")
    return // stop the server process
  }
  PORT := ":" + arguments[1]
  l , err := net.Listen("tcp", PORT)
  if err != nil { 
    fmt.Println(err)
    return 
  }
  defer l.Close() // closing connection before leaving this method 
  c ,err := l.Accept() 
  if err != nil  { 
    fmt.Println(err)
    return 
  }
  for { 
    netData, err := bufio.NewReader(c).ReadString('\n')
    if err != nil { 
      fmt.Println(err)
      return 
    }
    if strings.TrimSpace(string(netData)) == "STOP" { 
      fmt.Println("Exiting TCP server!")
      return 
    }
    fmt.Println("-> ", string(netData))
    t := time.Now()
    myTime := t.Format(time.RFC3339) + "\n"
    c.Write([]byte(myTime)) // write to the output aka. os.stdin
  }
}
