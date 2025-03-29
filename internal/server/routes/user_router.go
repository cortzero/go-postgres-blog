package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cortzero/go-postgres-blog/internal/model/user"
	"github.com/cortzero/go-postgres-blog/internal/server/response"
)

type UserRouter struct {
	Repository user.Repository
}

func (router *UserRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	ctx := r.Context()
	err = router.Repository.Create(ctx, &u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u.Password = ""
	w.Header().Add("Location", fmt.Sprintf("%s%d", r.URL.String(), u.ID))
	response.EncodeDataToJSON(w, r, http.StatusCreated, response.Map{"user": u})
}

func (router *UserRouter) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := router.Repository.GetAll(ctx)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"user": users})
}

func (router *UserRouter) GetByIdHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	user, err := router.Repository.GetById(ctx, uint(userId))
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, response.Map{"user": user})
}

func (router *UserRouter) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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
	err = router.Repository.Update(ctx, uint(userId), u)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}

func (router *UserRouter) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("id")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	err = router.Repository.Delete(ctx, uint(userId))
	if err != nil {
		response.CreateHTTPErrorMessage(w, r, http.StatusBadRequest, err.Error())
		return
	}

	response.EncodeDataToJSON(w, r, http.StatusOK, nil)
}
