package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/Masa4240/go-mission-catechdojo/service"
)

// A HealthzHandler implements health check endpoint.
type UserHandler struct {
	svc *service.UserService
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) Create(ctx context.Context, req *model.UserResistrationRequest) (*model.UserResistrationResponse, error) {
	//	fmt.Fprint(nil, "Hello TODO!")
	userinfo, err := h.svc.CreateUser(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	// return &model.CreateTODOResponse{TODO: *todo}, nil
	return &model.UserResistrationResponse{Token: *&userinfo.Token}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *model.UserGetRequest) (*model.UserGetReponse, error) {
	//	fmt.Fprint(nil, "Hello TODO!")
	userinfo, err := h.svc.GetUser(ctx, int(req.ID))
	if err != nil {
		return nil, err
	}
	// return &model.CreateTODOResponse{TODO: *todo}, nil
	return &model.UserGetReponse{Name: *&userinfo.Name}, nil
}

func (h *UserHandler) Update(ctx context.Context, req *model.UserUpdateRequest) (*model.UserUpdateReponse, error) {
	//	fmt.Fprint(nil, "Hello TODO!")
	_, err := h.svc.UpdateUser(ctx, req.Newname, int(req.ID))
	if err != nil {
		return nil, err
	}
	//	return &model.CreateTODOResponse{TODO: *todo}, nil
	// return &model.UserGetReponse{Name: *&userinfo.Name}, nil
	return nil, nil
}

// ServeHTTP implements http.Handler interface.
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("user handler Start")
	// fmt.Println(r.URL)
	// fmt.Println(r.Header)
	// fmt.Println(r.Host)
	if r.Method == "POST" && r.URL.String() == "/user/create" {
		var req = model.UserResistrationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "decode error")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "empty Name")
			//fmt.Fprintln(w, r)
			return
		}
		//fmt.Fprintln(w, req.Name)
		resp, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//jsonout, _ := json.Marshal(err)
			fmt.Fprintln(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//fmt.Fprintln(w, err, resp)
		fmt.Println("Finish Create User process")
		return
	}
	if r.Method == "GET" && r.URL.String() == "/user/get" {
		//fmt.Println(r.Context().Value("name"))
		fmt.Println("Start Get User process")
		fmt.Println(r.Context().Value("id"))
		var req = model.UserGetRequest{}
		req.ID = r.Context().Value("id").(int64)
		//req.Name = r.Context().Value("name").(string)

		resp, err := h.GetUser(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//jsonout, _ := json.Marshal(err)
			fmt.Fprintln(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//fmt.Fprintln(w, err, resp)
		fmt.Println("Finish Get User process")
		return

	}
	if r.Method == "PUT" && r.URL.String() == "/user/update" {
		fmt.Println("Start Update User process")
		var req = model.UserUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "decode error")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Newname == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "empty Name")
			//fmt.Fprintln(w, r)
			return
		}
		req.ID = r.Context().Value("id").(int64)
		fmt.Println(req)
		//fmt.Fprintln(w, req.Name)
		_, err := h.Update(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//jsonout, _ := json.Marshal(err)
			fmt.Fprintln(w, err)
			return
		}
		//fmt.Fprintln(w, err, resp)
		fmt.Println("Finish Update User process")
		return

	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Incorrect request")
}
