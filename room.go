package main

import (
	agonesV1 "agones.dev/agones/pkg/apis/agones/v1"
	allocationV1 "agones.dev/agones/pkg/apis/allocation/v1"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"net/http"
)

const runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Room struct {
	Code    string
	Address string
	Port    int32
}

type Response struct {
	success bool
	data    interface{}
	error   string
}

func generateRoomCode(length int) string {
	out := make([]byte, length)
	for i := range length {
		out[i] = runes[rand.Intn(len(runes))]
	}

	return string(out)
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, errorCode int, error string) {
	w.WriteHeader(errorCode)
	render.JSON(w, r, Response{
		success: false,
		data:    nil,
		error:   error,
	})
}

func sendSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.WriteHeader(200)
	render.JSON(w, r, Response{
		success: true,
		data:    data,
		error:   "",
	})
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	roomId := generateRoomCode(4)

	options := &allocationV1.GameServerAllocation{
		Spec: allocationV1.GameServerAllocationSpec{
			Selectors: []allocationV1.GameServerSelector{
				{
					LabelSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"agones.dev/fleet": "pigeon-project-fleet",
						},
					},
				},
			},
			MetaPatch: allocationV1.MetaPatch{
				Labels: map[string]string{
					"room": roomId,
				},
			},
		},
	}

	allocation, err := agonesClient.AllocationV1().GameServerAllocations("default").Create(context.TODO(), options, metav1.CreateOptions{})

	if err != nil {
		sendErrorResponse(w, r, http.StatusInternalServerError, "Error allocating game server")
		return
	}

	roomInfo := Room{
		Code:    roomId,
		Address: allocation.Status.Address,
		Port:    allocation.Status.Ports[0].Port,
	}

	sendSuccessResponse(w, r, roomInfo)
}

func getRoom(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "room")

	allocated := agonesV1.GameServerStateAllocated

	options := &allocationV1.GameServerAllocation{
		Spec: allocationV1.GameServerAllocationSpec{
			Selectors: []allocationV1.GameServerSelector{
				{
					GameServerState: &allocated,
					LabelSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{
							"agones.dev/fleet": "pigeon-project-fleet",
							"room":             roomId,
						},
					},
				},
			},
		},
	}

	allocation, err := agonesClient.AllocationV1().GameServerAllocations("default").Create(context.TODO(), options, metav1.CreateOptions{})

	if err != nil {
		sendErrorResponse(w, r, http.StatusInternalServerError, "Error allocating game server")
		return
	}

	// If no available ports, then could not find a room with that code
	if len(allocation.Status.Ports) == 0 {
		sendErrorResponse(w, r, http.StatusNotFound, "Could not find room with code "+roomId)
		return
	}

	roomInfo := Room{
		Code:    roomId,
		Address: allocation.Status.Address,
		Port:    allocation.Status.Ports[0].Port,
	}

	sendSuccessResponse(w, r, roomInfo)
}
