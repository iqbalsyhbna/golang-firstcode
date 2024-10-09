// internal/handlers/articles.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test-golang/internal/models"
	"test-golang/internal/service"
	"test-golang/pkg/common"

	"github.com/gorilla/mux"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(service *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		service: service,
	}
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.service.GetAll()
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "Failed to get articles")
		return
	}
	common.WriteJSON(w, http.StatusOK, articles, "Successfully fetched articles")
}

func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	article, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("article with ID %d not found", id) {
			common.WriteError(w, http.StatusNotFound, "Article not found")
		} else {
			common.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	common.WriteJSON(w, http.StatusOK, article, "Successfully fetched article")
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		fmt.Println(err)
		common.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdArticle, err := h.service.Create(article)
	if err != nil {
		fmt.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "Failed to create article")
		return
	}

	common.WriteJSON(w, http.StatusCreated, createdArticle, "Successfully created article")
}

func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println(err)
		common.WriteError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		fmt.Println(err)
		common.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	article.ID = id
	updatedArticle, err := h.service.Update(article)
	if err != nil {
		fmt.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "Failed to update article")
		return
	}

	common.WriteJSON(w, http.StatusOK, updatedArticle, "Article updated successfully")
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		common.WriteError(w, http.StatusInternalServerError, "Failed to delete article")
		return
	}

	common.WriteJSON(w, http.StatusOK, nil, "Article deleted successfully")
}
