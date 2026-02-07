package v1

import (
	"io"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	statusrpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (h *NoteHandler) Chat(st grpc.BidiStreamingServer[pb.Message, pb.Message]) error {
	ctx := st.Context()

	errSystemChan := make(chan error, 1)
	defer close(errSystemChan)

	messagesChan := h.initMessages()
	replyChan := make(chan *pb.Message, 10)
	errChan := make(chan *pb.Message, 10)

	var wg sync.WaitGroup

	// Reciever goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(replyChan)
		defer close(errChan)

		for {
			message, err := st.Recv()

			if err != nil {
				if err == io.EOF {
					h.log.Info("client stop sending messages")
				} else if err := errorRecieveHandling(h.log, err); err != nil {
					errSystemChan <- err
				}
				return
			}

			if msg, ok := message.Payload.(*pb.Message_Text); ok {
				if stat := h.validateMessage(msg); stat == nil {
					select {
					case replyChan <- message:
					default:
					}
				} else {
					message.Payload = &pb.Message_Error{Error: stat}
					select {
					case errChan <- message:
					default:
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	// Sender goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		for errChan != nil || replyChan != nil || messagesChan != nil {
			var message *pb.Message
			var ok bool

			select {
			case <-ctx.Done():
				errorSendHandling(h.log, ctx.Err())
				return
			case message, ok = <-errChan:
				if !ok {
					errChan = nil
					continue
				}
			case message, ok = <-replyChan:
				if !ok {
					replyChan = nil
					continue
				}
				message = h.prepareResponse(message)
			case message, ok = <-messagesChan:
				if !ok {
					messagesChan = nil
					continue
				}
			}

			if err := st.Send(message); err != nil {
				errorSendHandling(h.log, err)
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	var err error

	select {
	case <-ctx.Done():
		h.log.Info("context done")
	case <-done:
		h.log.Info("reciever and sender has been stopped")
	case err = <-errSystemChan:
		h.log.Error("chat error", logger.NewField("error", err))
	}

	h.log.Info("bidirectional stream finished")
	return err
}

func (h *NoteHandler) initMessages() chan *pb.Message {
	messages := []*pb.Message{
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "shark"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "lion"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "elephant"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "giraffe"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "zebra"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "monkey"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "lion"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "kangaroo"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Error{Error: &statusrpc.Status{Code: int32(codes.Internal), Message: "internal error"}},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "rabbit"},
		},
	}
	messagesChan := make(chan *pb.Message, len(messages))
	defer close(messagesChan)

	for _, msg := range messages {
		messagesChan <- msg
	}
	return messagesChan
}

func (h *NoteHandler) validateMessage(msg *pb.Message_Text) *statusrpc.Status {
	text := []rune(msg.Text)
	slices.Reverse(text)
	if msg.Text == string(text) {
		return &statusrpc.Status{
			Code:    int32(codes.InvalidArgument),
			Message: "the word reads the same on both sides",
		}
	}
	return nil
}

func (h *NoteHandler) prepareResponse(msg *pb.Message) *pb.Message {
	if i, ok := msg.Payload.(*pb.Message_Text); ok {
		text := []rune(i.Text)
		slices.Reverse(text)
		msg.Payload = &pb.Message_Text{Text: string(text)}
	}
	return msg
}
