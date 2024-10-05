package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-golang/database"
	"test-golang/models"

	"github.com/gorilla/mux"
)

func GetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	articles := []models.Article{}
	rows, err := database.DB.Query("SELECT id, title, content FROM articles ORDER BY id ASC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}

	json.NewEncoder(w).Encode(articles)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	var article models.Article
	err := database.DB.QueryRow("SELECT id, title, content FROM articles WHERE id=?", id).Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var article models.Article
	_ = json.NewDecoder(r.Body).Decode(&article)

	result, err := database.DB.Exec("INSERT INTO articles (title, content) VALUES (?, ?)", article.Title, article.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the inserted article
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	article.ID = strconv.FormatInt(id, 10)

	json.NewEncoder(w).Encode(article)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	var article models.Article
	_ = json.NewDecoder(r.Body).Decode(&article)
	article.ID = id

	_, err := database.DB.Exec("UPDATE articles SET title=?, content=? WHERE id=?", article.Title, article.Content, article.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	_, err := database.DB.Exec("DELETE FROM articles WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Article deleted successfully"})
}
