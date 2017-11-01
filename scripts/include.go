package main

import (
    "os"
    "fmt"
    "path/filepath"
    "io/ioutil"
    "encoding/base64"
    "strings"
    
)
var WebContent string
var WebContentURL string

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func VisitFile(fp string, fi os.FileInfo, err error) error {
		var fileContent []byte
	    check(err)
	    
	    if !!fi.IsDir() {
	        fmt.Println("directory: "+fp)
	        return nil // not a file.
	    }
	    
	    fmt.Println("file: "+fp)

	    fmt.Println("open file: "+".\\"+fp)
	    fileContent, err = ioutil.ReadFile(fp)
	    check(err)
	    URL:=strings.Replace(fp,"\\","/",-1)
	    
   	    WebContentURL += "`" + URL + "`,\n  "
   	    WebContent += "`" + base64.StdEncoding.EncodeToString(fileContent) + "`,\n  "
   	    // data, err := base64.StdEncoding.DecodeString(imgBase64) // data is of type []byte

	    return nil
	}


func main() {
    OutWebContent, _ := os.Create("./src/WebServerDefault/webContent.go")
    OutWebContent.Write([]byte("package main \n\n"))

	filepath.Walk("./web/",VisitFile)
	
    fmt.Println("WebContent")
    fmt.Println(WebContent)
    
    OutWebContent.Write([]byte("var WebContentURL = []string{"+WebContentURL+"}\n"))
    OutWebContent.Write([]byte("var WebContent = []string{"+WebContent+"}\n"))
    fmt.Println("Fertig!")

}