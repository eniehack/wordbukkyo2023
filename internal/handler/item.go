package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/WORD-COINS/wordbukkyo2/internal/models"
	"github.com/oklog/ulid/v2"
)

type AddItemRequestPayload struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func (h *Handler) ListItem(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)

	items, err := os.Open(h.ItemDB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer items.Close()
	syscall.Flock(int(items.Fd()), syscall.LOCK_SH)
	defer syscall.Flock(int(items.Fd()), syscall.LOCK_UN)
	buf.ReadFrom(items)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(buf.String()))
	w.WriteHeader(http.StatusOK)
	return
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var (
		newItem AddItemRequestPayload
		items   models.Items
	)
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		log.Println("JSON Decode failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	read_fp, err := os.Open("items.json")
	if err != nil {
		log.Println("Open failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	syscall.Flock(int(read_fp.Fd()), syscall.LOCK_SH)
	if err := json.NewDecoder(read_fp).Decode(&items); err != nil {
		log.Println("JSON Decode failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	syscall.Flock(int(read_fp.Fd()), syscall.LOCK_UN)
	read_fp.Close()

	now := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)

	items = append(items, models.Item{
		Name:  newItem.Name,
		ID:    ulid.MustNew(ulid.Timestamp(now), entropy).String(),
		Price: newItem.Price,
	})

	wfp, err := os.OpenFile("items.json", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("OpenFile failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer wfp.Close()
	syscall.Flock(int(wfp.Fd()), syscall.LOCK_EX)
	defer syscall.Flock(int(wfp.Fd()), syscall.LOCK_UN)
	if err := json.NewEncoder(wfp).Encode(items); err != nil {
		log.Println("JSON Encode failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
