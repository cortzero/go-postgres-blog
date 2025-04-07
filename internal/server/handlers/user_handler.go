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
	"github.com/cortzero/go-postgres-blog/internal/service/errors"
)

var (
	usersUrlRegExpNoVars = regexp.MustCompile(`^/api/v1/users$`)
	usersUrlRegExpVars   = regexp.MustCompile(`^/api/v1/users/(\d+)$`)
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
	case r.Method == http.MethodGet && usersUrlRegExpNoVars.Match([]byte(reqURL)):
		handler.GetAllHandler(w, r)
		return
	case r.Method == http.MethodGet && usersUrlRegExpVars.Match([]byte(reqURL)):
		handler.GetByIdHandler(w, r)
		return
	case r.Method == http.MethodPost && usersUrlRegExpNoVars.Match([]byte(reqURL)):
		handler.CreateHandler(w, r)
		return
	case r.Method == http.MethodPut && usersUrlRegExpVars.Match([]byte(reqURL)):
		handler.UpdateHandler(w, r)
		return
	case r.Method == http.MethodDelete && usersUrlRegExpVars.Match([]byte(reqURL)):
		handler.DeleteHandler(w, r)
		return
	default:
		newError := errors.NewCustomError(
			"NOT_FOUND",
			"Could not found the requested URL.",
			fmt.Sprintf("The URL '%s' does not exist.", r.URL.Path),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusNotFound, newError, r.URL.Path)
		return
	}
}

func (handler *UserHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
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
	//u.CreatedAt = time.Now()
	err = handler.Repository.Create(ctx, &u)
	if err != nil {
		newError := errors.NewCustomError(
			"ERROR_CREATING_USER",
			err.Error(),
			"",
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.EncodeDataToJSON(w, r, http.StatusCreated, response.Map{"userCreated": u})
}

func (handler *UserHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := handler.Repository.GetAll(ctx)
	if err != nil {
		newError := errors.NewCustomError(
			"ERROR",
			err.Error(),
			"",
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusBadRequest, newError, r.URL.Path)
		return
	}
	if users != nil {
		response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"users": users})
	} else {
		response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"users": []user.User{}})
	}
}

func (handler *UserHandler) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, err.Error(), r.URL.Path)
		return
	}

	ctx := r.Context()
	user, err := handler.Repository.GetById(ctx, uint(userId))
	if err != nil {
		newError := errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			"The requested resource was not found.",
			fmt.Sprintf("The user with ID '%d' does not exist.", userId),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusNotFound, newError, r.URL.Path)
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"user": user})
}

func (handler *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, err.Error(), r.URL.Path)
		return
	}

	var u user.User
	err = json.NewDecoder(r.Body).Decode(&u)
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
	u.UpdatedAt = time.Now()
	err = handler.Repository.Update(ctx, uint(userId), u)
	if err != nil {
		newError := errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			"The requested resource was not found.",
			err.Error(),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusNotFound, newError, r.URL.Path)
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}

func (handler *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, err.Error(), r.URL.Path)
		return
	}

	ctx := r.Context()
	err = handler.Repository.Delete(ctx, uint(userId))
	if err != nil {
		newError := errors.NewCustomError(
			"RESOURCE_NOT_FOUND",
			"The requested resource was not found.",
			err.Error(),
			time.Now())
		response.CreateErrorResponse(w, r, http.StatusNotFound, newError, r.URL.Path)
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}
