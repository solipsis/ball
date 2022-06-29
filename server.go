package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type server struct {
	state   state
	updates chan playerUpdate
}

type playerUpdate struct {
	Frame int   `json:"frame"`
	Input input `json:"input"`
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		log.Fatalf("establishing websocket: %v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	l := rate.NewLimiter(rate.Every(time.Millisecond*10), 10)
	for {
		pu, err := s.readUpdate(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Printf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
		s.updates <- pu
	}
}

func (s *server) readUpdate(ctx context.Context, c *websocket.Conn, l *rate.Limiter) (playerUpdate, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, buf, err := c.Read(ctx)
	if err != nil {
		return playerUpdate{}, fmt.Errorf("failed to read update")
	}
	pu := playerUpdate{}
	err = json.Unmarshal(buf, &pu)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(buf))
	}
	fmt.Println("received update:", pu)
	//	s.state[pu.ID].X = pu.X
	//	s.state[pu.ID].Y = pu.Y

	return pu, nil
}
