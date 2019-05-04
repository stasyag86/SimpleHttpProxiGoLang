package main
import (
  "net/http"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "os"
  "strconv"
  "hash/crc32"
  "log"
)


var http_protocol string = "http://"
var default_port = 9087

func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
  //fmt.Println("listening to port 9087: " + r.Method + " " + r.Host)
    if r.Method == http.MethodGet {
      handelGetRequest(r.Host, w)
    }else if r.Method == http.MethodPost{
      handlePostRequest(w, r)
    }
  }

  func handlePostRequest(w http.ResponseWriter, r *http.Request){
    err := r.ParseForm()
    if err != nil {
      w.WriteHeader(http.StatusServiceUnavailable)
      fmt.Fprintf(w, "ParseForm() err: %v", err)
      return
    }

    url := http_protocol + r.Host

    messageBody := map[string]interface{}{
    }

    for k, v := range r.Form {
      messageBody[k] = v
    }
  
    bytesRepresentation, err := json.Marshal(messageBody)
    if err != nil {
      w.WriteHeader(http.StatusServiceUnavailable)
      fmt.Fprintf(w, "json.Marshal() err: %v", err)
      return
    }
  
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
    if err != nil {
      w.WriteHeader(http.StatusServiceUnavailable)
      fmt.Fprintf(w, "http.Post() err: %v", err)
      return
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      w.WriteHeader(http.StatusServiceUnavailable)
      fmt.Fprintf(w, "ioutil.ReadAll() err: %v", err)
      return
    }

    log.Println("CRC32: ", crc32.ChecksumIEEE(bodyBytes))
    bodyString := string(bodyBytes)
    fmt.Fprintf(w, bodyString)

  }

  func handelGetRequest(host string, w http.ResponseWriter){
    url := http_protocol + host
    rs, err := http.Get(url)
    
    // Process response
    if err != nil {
      w.WriteHeader(http.StatusServiceUnavailable)
      fmt.Fprintf(w, "http.Get() err: %v", err)
      return
    }
    defer rs.Body.Close()

    bodyBytes, err := ioutil.ReadAll(rs.Body)
    if err != nil {
      w.WriteHeader(http.StatusServiceUnavailable)
      fmt.Fprintf(w, "ioutil.ReadAll() err: %v", err)
      return
    }

    log.Println("CRC32: ", crc32.ChecksumIEEE(bodyBytes))
    bodyString := string(bodyBytes)
    fmt.Fprintf(w, bodyString) // send data to client side
}


func main() {
  listeningPort := default_port
  if (len(os.Args) > 1){
    i, err := strconv.Atoi(os.Args[1])
    if err != nil {
      fmt.Println("port number is incorrect")
    }else{
      listeningPort = i
    }
  }

  log.Println("Http proxy is ready on port: ", listeningPort)
  http.HandleFunc("/", handleRequestAndRedirect)
  port := fmt.Sprintf(":%d", listeningPort) 
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
    } 
  
}