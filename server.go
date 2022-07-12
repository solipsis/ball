package ball

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type Server struct {
	state        state
	clientInputs chan playerUpdate
	clients      []*websocket.Conn
	inputBuffer  [][]input
	stateBuffer  []state
}

type playerUpdate struct {
	Frame int   `json:"frame"`
	Input input `json:"input"`
}

func NewServer() *Server {

	s := Server{
		clientInputs: make(chan playerUpdate),
		inputBuffer:  make([][]input, 2),
		stateBuffer:  make([]state, 60),
		state:        newState(),
	}
	s.inputBuffer[0] = make([]input, 60)
	s.inputBuffer[1] = make([]input, 60)

	return &s
}

func (s *Server) Listen() {

	l, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Fatalf("server net.Listen(): %v", err)
	}
	srv := &http.Server{
		Handler:      s,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Printf("listening on http://%v", l.Addr())
	if err := srv.Serve(l); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			fmt.Println("normal closure")
			return
		}
		if err != nil {
			log.Printf("failed to read update from client %v: %v", r.RemoteAddr, err)
			//return
			continue
		}
		/*
			if pu.Input.Dir != NONE {
				fmt.Println(pu)
			}
		*/
		//fmt.Println(pu)
		// TODO(IMPLEMENT)
		s.clientInputs <- pu
	}
}

func (s *Server) readUpdate(ctx context.Context, c *websocket.Conn, l *rate.Limiter) (playerUpdate, error) {
	//	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	//	defer cancel()

	_, buf, err := c.Read(context.TODO())
	if err != nil {
		return playerUpdate{}, fmt.Errorf("failed to read update: %v", err)
	}
	pu := playerUpdate{}
	err = json.Unmarshal(buf, &pu)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(buf))
	}
	//fmt.Println("received update:", pu)
	//	s.state[pu.ID].X = pu.X
	//	s.state[pu.ID].Y = pu.Y

	return pu, nil
}

// implements ebiten.Game
func (s *Server) Update() error {
	s.update()
	return nil
}

func (s *Server) update() {
	// read all pending inputs
	userUpdate := false
	for {
		select {
		case up := <-s.clientInputs:
			s.handleRollback(up)
			userUpdate = true
			continue
		default:
			break
		}
		break
	}
	if userUpdate {

	}

	// apply all inputs to buffer
	// rollback to oldest of updates

	s.state = step(s.state, s.inputBuffer)
}

func (s *Server) handleRollback(up playerUpdate) {
	//prevFrame := up.Frame
	//prevState := s.stateBuffer[prevFrame%len(stateBuffer)]

}

// implements ebiten.Game
func (s *Server) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.state.ball.pos.x-s.state.ball.radius, s.state.ball.pos.y-s.state.ball.radius)
	screen.DrawImage(ballImage, op)
	op = &ebiten.DrawImageOptions{}
	DrawRect(screen, 100, 200, 20, 50, color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0x00})
	s.state.players[0].draw(screen)
}

// implements ebiten.Game
func (s *Server) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
