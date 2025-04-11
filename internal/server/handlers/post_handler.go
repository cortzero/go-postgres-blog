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
	"github.com/cortzero/go-postgres-blog/internal/server/response"
	"github.com/cortzero/go-postgres-blog/internal/service/errors"
)

var (
	postsUrlRegExpNoVars = regexp.MustCompile(`^/api/v1/posts$`)
	postsUrlRegExpVars   = regexp.MustCompile(`^/api/v1/posts/(\d+)$`)
)

type PostHandler struct {
	Service post.Service
}

func NewPostHandler(service post.Service) *PostHandler {
	return &PostHandler{
		Service: service,
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

	posts, err := handler.Service.GetAllPosts(ctx)
	if err != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, err, r.URL.Path)
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
		newError := errors.NewCustomError(
			"ERROR",
			err.Error(),
			fmt.Sprintf("Error parsing the path variable '%s' on the URL", postIdStr),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}

	ctx := r.Context()
	post, error_get := handler.Service.GetPostById(ctx, uint(postId))
	if error_get != nil {
		response.CreateErrorResponse(w, r, http.StatusNotFound, error_get, r.URL.Path)
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"post": post})
}

func (handler *PostHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		newError := errors.NewCustomError(
			"BAD_REQUEST",
			"The request is malformed.",
			"The body of the request may have an incorrect format.",
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	error_creating := handler.Service.CreatePost(ctx, &p)
	if error_creating != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, error_creating, r.URL.Path)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.Path, p.ID))
	response.EncodeDataToJSON(w, r, http.StatusCreated, response.Map{"postCreated": p})
}

func (handler *PostHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func (handler *PostHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

}
