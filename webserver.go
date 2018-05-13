package main


import (
  "net/http"      //http handling and socket
  "io/ioutil"     //io functions like R/W files
  "strings"       //string handling
  "log"           //logging
  "flag"          //command-line flag
)



func main(){
  //read user input in pointer
  portPtr := flag.String("port", "80", "Listen and serve on this port")
  pathPtr := flag.String("path", "/var/www/", "Folder which contains web-files")
  //parse user input
  flag.Parse()

  //log user input
  log.Printf("Serving and listening on port: " + *portPtr)
  log.Printf("Serving files from directory: " + *pathPtr)

  //handler for everything located in the given path
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    rootHandler(w, r, pathPtr)
  })

  //start server and bind to specific port, check also for errors
  log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
}



//HTTP-Handler for the root folder (defined as destination where main.go is launched)
func rootHandler(w http.ResponseWriter, r *http.Request, pathPtr *string) {
  logRequest(r)
  serveStaticContent(w, r, pathPtr)
}



//Function is loging basically all HTTP-Requests
//It logs the requested path, and passed parameters
func logRequest(r *http.Request) {
  path := r.URL.Path[1:]
  args := r.URL.Query()
  log.Println(path, args)
}



//Function to server the requested content
//Checks also if requested file has a defined suffix
func serveStaticContent(w http.ResponseWriter, r *http.Request, pathPtr *string) {
  path := r.URL.Path[1:]
  data, err := ioutil.ReadFile(*pathPtr + string(path))

  if err == nil {
    var contentType string
    if strings.HasSuffix(path, ".css") {
      contentType = "text/css"
    } else if strings.HasSuffix(path, ".html") {
      contentType = "text/html"
    } else if strings.HasSuffix(path, ".png") {
      contentType = "image/png"
    } else if strings.HasSuffix(path, ".js") {
      contentType = "application/javascript"
    } else if strings.HasSuffix(path, ".svg") {
      contentType = "image/svg+xml"
    } else {
      contentType = "text/plain"
    }

    //allow any origin to access the resource
    w.Header().Add("Access-Control-Allow-Origin", "*")
    //add content type to header
    w.Header().Add("Content-Type", contentType)
    //write data to browser and fulfill the request
    w.Write(data)
  } else {
    //tell the header that no such file found
    w.WriteHeader(404)
    //write text which is displayed to the user
    w.Write([]byte("404 Site " + http.StatusText(404)))
  }
}
