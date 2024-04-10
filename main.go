package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const (
	clientID     = "5406ffed7c33489b915321267e3ca75f"
	clientSecret = "bbaefa27380049a39c30e53c9bddc0c6"
)

type spotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type spotifyTrack struct {
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	TrackURL string `json:"track_url"`
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/play", playHandler)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Spotify Music Player</title>
</head>
<body>
    <h1>Search for a song</h1>
    <form action="/search" method="get">
        <input type="text" name="query" placeholder="Enter song name">
        <button type="submit">Search</button>
    </form>
</body>
</html>
`
	fmt.Fprintf(w, tmpl)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	accessToken, err := getAccessToken()
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", query)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to perform request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var searchResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	tracks := extractTracks(searchResponse)

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Search Results</title>
</head>
<body>
    <h1>Search Results</h1>
    <ul>
        {{range .}}
        <li>{{.Name}} - {{.Artist}} <a href="/play?track_url={{.TrackURL}}">Play</a></li>
        {{end}}
    </ul>
</body>
</html>
`
	t := template.Must(template.New("searchResults").Parse(tmpl))
	t.Execute(w, tracks)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	trackURL := r.URL.Query().Get("track_url")
	if trackURL == "" {
		http.Error(w, "Missing track URL parameter", http.StatusBadRequest)
		return
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Play Music</title>
</head>
<body>
    <h1>Now Playing</h1>
    <audio controls autoplay>
        <source src="{{.}}" type="audio/mpeg">
    </audio>
</body>
</html>
`
	fmt.Fprintf(w, tmpl, trackURL)
}

func getAccessToken() (string, error) {
	authURL := "https://accounts.spotify.com/api/token"
	body := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", authURL, body)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse spotifyTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func extractTracks(response map[string]interface{}) []spotifyTrack {
	var tracks []spotifyTrack

	if tracksData, ok := response["tracks"].(map[string]interface{})["items"].([]interface{}); ok {
		for _, trackData := range tracksData {
			if track, ok := trackData.(map[string]interface{}); ok {
				name := track["name"].(string)
				artist := track["artists"].([]interface{})[0].(map[string]interface{})["name"].(string)
				trackURL := track["external_urls"].(map[string]interface{})["spotify"].(string)

				tracks = append(tracks, spotifyTrack{Name: name, Artist: artist, TrackURL: trackURL})
			}
		}
	}

	return tracks
}
