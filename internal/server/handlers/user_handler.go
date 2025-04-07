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
	Service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{
		Service: service,
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
	err_creation := handler.Service.CreateUser(ctx, &u)
	if err_creation != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, err_creation, r.URL.Path)
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.EncodeDataToJSON(w, r, http.StatusCreated, response.Map{"userCreated": u})
}

func (handler *UserHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := handler.Service.GetAllUsers(ctx)
	if err != nil {
		response.CreateErrorResponse(w, r, http.StatusBadRequest, err, r.URL.Path)
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
	user, error_get := handler.Service.GetUserById(ctx, uint(userId))
	if error_get != nil {
		response.CreateErrorResponse(w, r, http.StatusNotFound, error_get, r.URL.Path)
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
	error_update := handler.Service.UpdateUser(ctx, uint(userId), &u)
	if error_update != nil {
		response.CreateErrorResponse(w, r, http.StatusNotFound, error_update, r.URL.Path)
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
	error_deleting := handler.Service.DeleteUser(ctx, uint(userId))
	if error_deleting != nil {
		response.CreateErrorResponse(w, r, http.StatusNotFound, error_deleting, r.URL.Path)
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}
