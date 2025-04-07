package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cortzero/go-postgres-blog/internal/model/post"
	"github.com/cortzero/go-postgres-blog/internal/server/errors"
	"github.com/cortzero/go-postgres-blog/internal/server/response"
)

var (
	postsUrlRegExpNoVars = regexp.MustCompile(`^/api/v1/posts$`)
	postsUrlRegExpVars   = regexp.MustCompile(`^/api/v1/posts/(\d+)$`)
)

type PostHandler struct {
	Repository post.Repository
}

func NewPostHandler(repository post.Repository) *PostHandler {
	return &PostHandler{
		Repository: repository,
	}
}

func (handler *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqURL := strings.TrimSuffix(r.URL.Path, "/")
	switch {
	case r.Method == http.MethodGet && postsUrlRegExpNoVars.MatchString(reqURL):
		handler.GetAllHandler(w, r)
		return
	case r.Method == http.MethodGet && postsUrlRegExpVars.MatchString(reqURL):
		handler.GetByIdHandler(w, r)
		return
	case r.Method == http.MethodPost && postsUrlRegExpNoVars.MatchString(reqURL):
		handler.CreateHandler(w, r)
		return
	case r.Method == http.MethodPut && postsUrlRegExpVars.MatchString(reqURL):
		handler.UpdateHandler(w, r)
		return
	case r.Method == http.MethodDelete && postsUrlRegExpVars.MatchString(reqURL):
		handler.DeleteHandler(w, r)
		return
	default:
		fmt.Println("Not Found: " + r.URL.Path)
		return
	}
}

func (handler *PostHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := handler.Repository.GetAll(ctx)
	if err != nil {
		newError := errors.NewErrorObject(
			"ERROR",
			err.Error(),
			"",
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}
	if posts != nil {
		response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"posts": posts})
	} else {
		response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"posts": []post.Post{}})
	}
}

func (handler *PostHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.PathValue("id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		newError := errors.NewErrorObject(
			"ERROR",
			err.Error(),
			fmt.Sprintf("Error parsing the path variable '%s' on the URL", postIdStr),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}

	ctx := r.Context()
	post, err := handler.Repository.GetById(ctx, uint(postId))
	if err != nil {
		newError := errors.NewErrorObject(
			"RESOURCE_NOT_FOUND",
			"The requested resource was not found.",
			fmt.Sprintf("The post with ID '%d' does not exist.", postId),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusNotFound, newError, r.URL.Path)
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"post": post})
}

func (handler *PostHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		newError := errors.NewErrorObject(
			"BAD_REQUEST",
			"The request is malformed.",
			"The body of the request may have an incorrect format.",
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	p.CreatedAt = time.Now()
	err = handler.Repository.Create(ctx, &p)
	if err != nil {
		newError := errors.NewErrorObject(
			"ERROR_CREATING_POST",
			"An error occurred while creating the post.",
			err.Error(),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.Path, p.ID))
	response.EncodeDataToJSON(w, r, http.StatusCreated, response.Map{"postCreated": p})
}

func (handler *PostHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func (handler *PostHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

}
