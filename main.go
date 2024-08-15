package main

import (
	"flag"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type InnerData struct {
	Short string
	Long  string
}

var shortUrlSlice []InnerData

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	shortUrlSlice = make([]InnerData, 0, 20)
}

func main() {
	var port string
	var baseUrl string

	defaultPort := ":8080"
	defaultBaseUrl := "http://localhost:8080/"

	envPort := os.Getenv("SERVER_ADDRESS")
	envBaseUrl := os.Getenv("BASE_URL")

	flag.StringVar(&port, "port", defaultPort, "port to run server")
	flag.StringVar(&baseUrl, "base-url", defaultBaseUrl, "base URL for shortened links")
	flag.Parse()

	if envPort != "" {
		port = envPort
	}
	if envBaseUrl != "" {
		baseUrl = envBaseUrl
	}

	r := chi.NewRouter()
	r.Post("/", responceUrl(baseUrl))
	r.Get("/{id}", requestUrl)

	err := http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}

func responceUrl(baseUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error read body", http.StatusInternalServerError)
		}
		shortUrl := doShortUrl(string(body))
		defer r.Body.Close()
		if shortUrl == "" {
			w.Write([]byte("Already Exist"))
			w.WriteHeader(200)
		} else {
			w.Write([]byte(baseUrl + shortUrl))
			w.WriteHeader(201)
		}
	}
}

func requestUrl(w http.ResponseWriter, r *http.Request) {
	var longUrl string
	path := chi.URLParam(r, "id")
	if path == "" {
		http.Error(w, "Empty short Url", http.StatusBadRequest)
		return
	}

	for _, data := range shortUrlSlice {
		if data.Short == path {
			longUrl = data.Long
		}
	}

	if longUrl == "" {
		w.WriteHeader(400)
		w.Write([]byte("Not found url"))
	} else {
		w.WriteHeader(307)
		w.Write([]byte(longUrl))
	}
}

func doShortUrl(url string) string {
	for _, data := range shortUrlSlice {
		if data.Long == url {
			return ""
		}
	}
	shortUrl := generateRandomString(7)
	newDataUrl := InnerData{
		Short: shortUrl,
		Long:  url,
	}
	shortUrlSlice = append(shortUrlSlice, newDataUrl)
	return shortUrl
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
