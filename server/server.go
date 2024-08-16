package main

import (
	"ascii_renderer/renderer"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func upload_image(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Image received!")
	if r.Method != http.MethodPost {
		http.Error(w, "Only post calls are allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	widthStr := r.FormValue("width")
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		http.Error(w, "Invalid width value", http.StatusBadRequest)
		return
	}

	ascii_array := renderer.Parse_and_return_ascii(file, width)
	fmt.Println("Image parsed into ascii!")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ascii_array)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/image_upload", upload_image)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}

}
