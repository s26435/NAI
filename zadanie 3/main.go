package main

// Autorzy:
// Jan Wolski s26435
// Marcin Topolniak s25672

//Polecenie: Zaimplementuj silnik rekomandacji filmów/seriali.
//
//- Przestudiuj materiał	A Comparative Study of Clustering Algorithms | by ishika chatterjee | Analytics Vidhya | Medium
//- Rozważ uzupełnienie ankiety (samodzielnie)
//- Zbuduj silnik rekomendacji filmów i/lub seriali.
//- Zaproponuj 5 filmów interesujących dla wybranego użytkownika, których nie oglądał.
//- Zaproponouj 5 film, których dany użytkownik nie powinnien oglądać (antyrekomendacje).

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const APIKEY string = "api_key_hihi"

// MovieRating to struktura przechowywująca dane uzytkowników i filmów jakie obejrzeli i jak je ocenili wraz z id z bazy danych IM
type MovieRating struct {
	PersonID   int
	MovieTitle string
	IMDBID     string
	Rating     float64
}

// MovieRatings to struktura przechowywująca oceny poszczególnych filmów
type MovieRatings struct {
	Ratings []MovieRating
}

// funkcja która z podanym tytułem zwraca IMDBID
func (movieRatings *MovieRatings) GetIMDBIDByTitle(title string) (string, error) {
	for _, movie := range movieRatings.Ratings {
		if movie.MovieTitle == title {
			return movie.IMDBID, nil
		}
	}
	return "", fmt.Errorf("no sutch a film")
}

// LoadCSV to metoda która ładuje oceny filmów z pliku CSV do struktury MovieRatings
func (mr *MovieRatings) LoadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read CSV: %v", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 3 {
			fmt.Printf("Warning: invalid record at line %d: expected at least 3 fields, got %d\n", i+1, len(record))
			continue
		}

		personID, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Printf("Warning: could not parse person ID at line %d: %v\n", i+1, err)
			continue
		}

		rating := 0.0
		if record[2] != "" {
			rating, err = strconv.ParseFloat(record[2], 64)
			if err != nil {
				fmt.Printf("Warning: could not parse rating for movie '%s' at line %d: %v\n", record[1], i+1, err)
				rating = 0.0
			}
		}

		movieRating := MovieRating{
			PersonID:   personID,
			MovieTitle: record[1],
			Rating:     rating,
		}
		mr.Ratings = append(mr.Ratings, movieRating)
	}

	return nil
}

// LoadIMDBIDs to metoda ładująca IMDBID z pliku csv do struktury MovieRatings
func (mr *MovieRatings) LoadIMDBIDs(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read CSV: %v", err)
	}

	imdbMap := make(map[string]string)
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 2 {
			fmt.Printf("Warning: invalid record at line %d: expected at least 2 fields, got %d\n", i+1, len(record))
			continue
		}
		imdbMap[record[0]] = record[1]
	}
	for i, rating := range mr.Ratings {
		if imdbID, ok := imdbMap[rating.MovieTitle]; ok {
			mr.Ratings[i].IMDBID = imdbID
		}
	}
	return nil
}

// metoda RecommendMovies struktury MovieRatings zwraca dwie tablice stringów zawierajace 5 tytułów filmów wybranych przez algorytm dla urzytkownika o podanym id (dla pana id = 1)
// tworzy mapę filmów i wyłącza te oglądane przez użytkownika
func (mr *MovieRatings) RecommendMovies(personID int, k int) ([]string, []string) {
	userRatings := make(map[int]map[string]float64)
	for _, rating := range mr.Ratings {
		if _, ok := userRatings[rating.PersonID]; !ok {
			userRatings[rating.PersonID] = make(map[string]float64)
		}
		userRatings[rating.PersonID][rating.MovieTitle] = rating.Rating
	}

	users := []int{}
	for user := range userRatings {
		users = append(users, user)
	}

	centroids := make([]int, k)
	for i := 0; i < k; i++ {
		centroids[i] = users[i%len(users)]
	}

	clusters := make(map[int][]int)
	for i := 0; i < 10; i++ {
		for j := 0; j < k; j++ {
			clusters[j] = []int{}
		}

		for _, user := range users {
			closest := 0
			closestDistance := math.MaxFloat64
			for j, centroid := range centroids {
				distance := calculateDistance(userRatings[user], userRatings[centroid])
				if distance < closestDistance {
					closest = j
					closestDistance = distance
				}
			}
			clusters[closest] = append(clusters[closest], user)
		}
		centroids = calculateNewCentroids(clusters, userRatings)
	}

	userCluster := -1
	for j, users := range clusters {
		for _, user := range users {
			if user == personID {
				userCluster = j
				break
			}
		}
		if userCluster != -1 {
			break
		}
	}

	recommendedMovies := make(map[string]float64)
	seenMovies := userRatings[personID]
	for _, user := range clusters[userCluster] {
		for movie, rating := range userRatings[user] {
			if _, seen := seenMovies[movie]; !seen {
				recommendedMovies[movie] += rating
			}
		}
	}

	type movieRecommendation struct {
		Movie  string
		Rating float64
	}
	recommendations := []movieRecommendation{}
	for movie, rating := range recommendedMovies {
		recommendations = append(recommendations, movieRecommendation{Movie: movie, Rating: rating})
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Rating > recommendations[j].Rating
	})

	var bestMovies, worstMovies []string
	for i := 0; i < len(recommendations) && i < 5; i++ {
		bestMovies = append(bestMovies, recommendations[i].Movie)
	}

	for i := len(recommendations) - 1; i >= 0 && len(recommendations)-i <= 5; i-- {
		worstMovies = append(worstMovies, recommendations[i].Movie)
	}

	return bestMovies, worstMovies
}

// calculateDistance oblicza odległość Euklidesową między dwoma użytkownikami
func calculateDistance(user1, user2 map[string]float64) float64 {
	sum := 0.0
	for movie, rating1 := range user1 {
		rating2, ok := user2[movie]
		if ok {
			sum += (rating1 - rating2) * (rating1 - rating2)
		} else {
			sum += rating1 * rating1
		}
	}
	for movie, rating2 := range user2 {
		if _, ok := user1[movie]; !ok {
			sum += rating2 * rating2
		}
	}
	return math.Sqrt(sum)
}

// calculateNewCentroid oblicza nowy centroid jako średnią ocen użytkowników w klastrze użycie algorytmu "Mean od Nearest Points"
func calculateNewCentroids(clusters map[int][]int, userRatings map[int]map[string]float64) []int {
	centroids := make([]int, len(clusters))
	for clusterIdx, users := range clusters {
		if len(users) == 0 {
			continue
		}

		minDistanceSum := math.MaxFloat64
		newCentroid := users[0]

		for _, candidate := range users {
			distanceSum := 0.0
			for _, user := range users {
				distanceSum += calculateDistance(userRatings[candidate], userRatings[user])
			}

			if distanceSum < minDistanceSum {
				minDistanceSum = distanceSum
				newCentroid = candidate
			}
		}

		centroids[clusterIdx] = newCentroid
	}

	return centroids
}

// getMovieDetails tworzy zapytania dotyczące filmów na podstawie imdbID
func getMovieDetails(imdbID string) (string, error) {
	imdbID = strings.TrimPrefix(imdbID, "tt")
	imdbID = strings.TrimSpace(imdbID)
	url := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&i=%s", APIKEY, imdbID)
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not fetch movie details: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	// Print the raw response for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %v", err)
	}

	//fmt.Println("API Response:", string(bodyBytes))

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("could not decode response: %v", err)
	}

	if response, ok := result["Response"].(string); !ok || response == "False" {
		return "", fmt.Errorf("movie not found: %v", result["Error"])
	}

	details := fmt.Sprintf("Title: %s\nYear: %s\nRated: %s\nReleased: %s\nRuntime: %s\nGenre: %s\nDirector: %s\nActors: %s\nPlot: %s\nIMDB Rating: %s\n",
		getString(result, "Title"), getString(result, "Year"), getString(result, "Rated"),
		getString(result, "Released"), getString(result, "Runtime"), getString(result, "Genre"),
		getString(result, "Director"), getString(result, "Actors"), getString(result, "Plot"),
		getString(result, "imdbRating"))

	return details, nil
}

// getString zwraca dane z jsona, N/A jeśli są niepoprawne albo ich nie ma
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return "N/A"
}

func main() {
	var movieRatings MovieRatings
	err := movieRatings.LoadCSV("dane.csv")
	if err != nil {
		fmt.Println("Error loading CSV:", err)
		return
	}
	err = movieRatings.LoadIMDBIDs("imdb.csv")
	if err != nil {
		fmt.Println("Error loading IMDB CSV:", err)
		return
	}

	recommendations, antiRecommendations := movieRatings.RecommendMovies(1, 2)

	fmt.Println("Top 5 Recommended Movies:")
	for _, movie := range recommendations {
		fmt.Println(movie)
		imdbID, err := movieRatings.GetIMDBIDByTitle(movie)
		//fmt.Println(imdbID)
		if err != nil {
			fmt.Println("Error fetching movie details:", err)
		}
		if details, err := getMovieDetails(imdbID); err == nil {
			fmt.Println(details)
		} else {
			fmt.Println("Error fetching movie details:", err)
		}
	}

	fmt.Println("\nTop 5 Anti-Recommended Movies:")
	for _, movie := range antiRecommendations {
		fmt.Println(movie)
		imdbID, err := movieRatings.GetIMDBIDByTitle(movie)
		if err != nil {
			fmt.Println("Error fetching movie details:", err)
		}
		if details, err := getMovieDetails(imdbID); err == nil {
			fmt.Println(details)
		} else {
			fmt.Println("Error fetching movie details:", err)
		}
	}

}