package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "path/filepath"
  )

const (
  port = :8080,
  uploadDir = "./data"
  )

func main(){
  _ = os.MakeDir(uploadDir, 0755)
  http.HandleFunc("/api/upload", uploadHandler)
  fmt.Printf("API is running on http://localhost:%s\n", port)
  log.Fatal(http.ListenAndServe(port, nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request){
  if r.Method != http.MethodPost {
    http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
  }
  
  if err := r.ParseMultipartForm(32 << 20); err != nil{
    http.Error(w, "Bad request", http.StatusBadRequest)
    return 
  }
  
  file, header, err := r.FormFile("file")
  if err != nil {
    http.Error(w, "No file uploaded. Please upload a file", http.StatusBadRequest)
    return 
  }
  defer file.Close()
  question := r.FormValue("question")
  if question == ""{
    http.Error(w, "Question is required", http.StatusBadRequest)
    return 
  }
  
  // to save the file
  filePath := filePath.Join(uploadDir, "sample"+filepath.Ext(header.Filename))
  out, err := os.Create(filePath)
  if err != nil{
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
  defer out.Close()
  
  _, err := io.Copy(out, file)
  if err != nil{
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
  
  // process rag
  content, err := processFileWithRag(filePath, question)
  if err != nil{
    http.Error(w, "Rag Processing Failed", http.StatusInternalServerError)
    return
  }
  
  _, = os.Remove(filePath)

  if err := os.Remove(filePath); err != nil {
    log.Printf("Warning: failed to remove file: %v", err)
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Printf(w, `{
    "message": "Content has been generated successfully",
    "data":{
      "content" %q
    }
  }`, content)
}