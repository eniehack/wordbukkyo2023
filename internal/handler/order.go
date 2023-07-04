package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/WORD-COINS/wordbukkyo2/internal/models"
)

type NewOrderRequestPayload struct {
	ItemID    string `json:"item_id"`
	StudentID string `json:"student_id"`
	Number    uint   `json:"number,omitempty"`
}

func (h *Handler) NewOrder(w http.ResponseWriter, r *http.Request) {
	order := new(NewOrderRequestPayload)

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	var users map[string]uint
	userfp, err := os.Open(h.UserDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	syscall.Flock(int(userfp.Fd()), syscall.LOCK_SH)
	if err = json.NewDecoder(userfp).Decode(&users); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	syscall.Flock(int(userfp.Fd()), syscall.LOCK_UN)
	if _, ok := users[order.StudentID]; ok == false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userfp.Close()

	var item models.Item
	itemfp, err := os.Open(h.ItemDB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	syscall.Flock(int(itemfp.Fd()), syscall.LOCK_SH)
	items := new(models.Items)
	if err = json.NewDecoder(itemfp).Decode(items); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	syscall.Flock(int(itemfp.Fd()), syscall.LOCK_UN)
	itemfp.Close()

	for _, val := range *items {
		if val.ID == order.ItemID {
			item = val
			break
		}
	}
	if item.ID == "" {
		log.Println("item.ID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if users[order.StudentID] < uint(item.Price) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users[order.StudentID] -= uint(item.Price)
	userwfp, err := os.OpenFile(h.UserDB, os.O_WRONLY|os.O_TRUNC, 644)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(userwfp).Encode(&users); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
