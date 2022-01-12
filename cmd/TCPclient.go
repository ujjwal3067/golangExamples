package main
import( 
  "bufio"
  "fmt"
  "net"
  "os"
  "strings"
)

func main() { 
  arguments := os.Args
  if len(arguments)  == 1 { 
    fmt.Println("Please Provide host:port")
    return // stop the client process 
  }
  CONNECT := arguments[1] // address
  c ,err := net.Dial("tcp", CONNECT) // try connecting to server
  if err != nil { 
    fmt.Println(err)
    return 
  }
  for { 
    reader := bufio.NewReader(os.Stdin) 
    fmt.Print(">> ")
    text , _ := reader.ReadString('\n')
    fmt.Fprintf(c, text + "\n")
    message, _ := bufio.NewReader(c).ReadString('\n')
    fmt.Print("->:" + message)
    if strings.TrimSpace(string(text)) == "STOP" { 
      fmt.Println("TCP client exiting ...")
      return // stop the client process 
    }
  }
}
