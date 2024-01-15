package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
)

var bindIP = "0.0.0.0"
var port = "8084"

var credentialsFilename = "credentials_m-adas.json"
var googleDriveID = "1hmUNiRzApr6RNHTQZqUPwRXqJf1yhiK2"	
var subfolder = "files/"

func main() {

	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)
	r.POST("/upload/:filename", uploadHandler)
	err := r.Run(bindIP + ":" + port)
	if err != nil {
		log.Printf("impossible to start server: %s", err)
	}

}
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func uploadHandler(c *gin.Context) {
	filename := c.Param("filename")
	file, _ := c.FormFile("file")
	err := c.SaveUploadedFile(file, subfolder+file.Filename)
	if err != nil {
		log.Println(err)
	}
	data, _ := file.Open()

	credentials := credentialsFilename
	b, err := os.ReadFile(credentials)
	if err != nil {
		log.Println(err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Println(err)
	}
	client := getClient(config)

	srv, err := drive.NewService(context.Background(),
		option.WithCredentialsFile("credentials.json"),
		option.WithHTTPClient(client),
		option.WithScopes(drive.DriveScope))
	if err != nil {
		log.Println("Unable to access Drive API:", err)
	}
	res, err := srv.Files.Create(
		&drive.File{
			Parents: []string{googleDriveID},
			Name:    filename,
		},
	).Media(data).Do()

	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", res.Id)

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("ok "+filename))
}
