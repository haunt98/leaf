package image

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *Service
}

func (h *Handler) HandleReceive(w http.ResponseWriter, req *http.Request) {
	var receiveReq ReceiveRequest
	if err := json.NewDecoder(req.Body).Decode(&receiveReq); err != nil {
		log.Println(err)
		return
	}

	receiveRsp, err := h.Service.Receive(receiveReq)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&receiveRsp); err != nil {
		log.Println(err)
		return
	}

	log.Println("LGTM")
}

func (h *Handler) HandleGetStatus(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	status, err := h.Service.StatusRepo.Read(vars["uuid"])
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&status); err != nil {
		log.Println(err)
		return
	}

	log.Println("LGTM")
}
