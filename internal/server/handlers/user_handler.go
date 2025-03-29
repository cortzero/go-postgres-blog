package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cortzero/go-postgres-blog/internal/model/user"
	"github.com/cortzero/go-postgres-blog/internal/server/response"
)

var (
	urlReWithoutParams = regexp.MustCompile(`^/api/v1/users$`)
	urlReWithParams    = regexp.MustCompile(`^/api/v1/users/(\d+)$`)
)

type UserHandler struct {
	Repository user.Repository
}

func NewUserHandler(repository user.Repository) *UserHandler {
	return &UserHandler{
		Repository: repository,
	}
}

func (handler *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqURL := strings.TrimSuffix(r.URL.Path, "/")
	switch {
	case r.Method == http.MethodGet && urlReWithoutParams.Match([]byte(reqURL)):
		handler.GetAllHandler(w, r)
		return
	case r.Method == http.MethodGet && urlReWithParams.Match([]byte(reqURL)):
		handler.GetByIdHandler(w, r)
		return
	case r.Method == http.MethodPost && urlReWithoutParams.Match([]byte(reqURL)):
		handler.CreateHandler(w, r)
		return
	case r.Method == http.MethodPut && urlReWithParams.Match([]byte(reqURL)):
		handler.UpdateHandler(w, r)
		return
	case r.Method == http.MethodDelete && urlReWithParams.Match([]byte(reqURL)):
		handler.DeleteHandler(w, r)
		return
	default:
		fmt.Println("404 Not Found")
		return
	}
}

func (handler *UserHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	u.CreatedAt = time.Now()
	err = handler.Repository.Create(ctx, &u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.EncodeDataToJSON(w, r, http.StatusCreated, response.Map{"user": u})
}

func (handler *UserHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := handler.Repository.GetAll(ctx)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"users": users})
}

func (handler *UserHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	user, err := handler.Repository.GetById(ctx, uint(userId))
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusNotFound, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"user": user})
}

func (handler *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var u user.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	u.UpdatedAt = time.Now()
	err = handler.Repository.Update(ctx, uint(userId), u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}

func (handler *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = handler.Repository.Delete(ctx, uint(userId))
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}
