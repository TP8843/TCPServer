package main

import (
	v1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"net/http"
)

const runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Room struct {
	Address string
	Port    int32
}

func generateRoomCode(length int) string {
	out := make([]byte, length)
	for i := range length {
		out[i] = runes[rand.Intn(len(runes))]
	}

	return string(out)
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	//roomId := generateRoomCode(4)

	options := &v1.GameServerAllocation{
		Spec: v1.GameServerAllocationSpec{},
	}

	allocation, err := agonesClient.AllocationV1().GameServerAllocations("default").Create(context.TODO(), options, metav1.CreateOptions{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error allocating room"))
		return
	}

	roomInfo := Room{
		Address: allocation.Status.Address,
		Port:    allocation.Status.Ports[0].Port,
	}

	render.JSON(w, r, roomInfo)
}

func getRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "room")

	w.Write([]byte("It's room getting time " + roomId))
}
