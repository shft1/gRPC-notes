package v1

import (
	"context"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"google.golang.org/grpc/status"
)

func (gw *noteGateway) Chat(ctx context.Context, errChan chan<- error) {
	chatCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	st, err := gw.client.Chat(chatCtx)

	if err != nil {
		gw.log.Error("failed to init bi-di stream for chat", logger.NewField("error", err))
		if stat, ok := status.FromError(err); ok {
			err = mapErrorRPC(gw.log, stat)
		} else {
			err = mapError(gw.log, err)
		}
		errChan <- err
		return
	}

	errChan <- nil
	gw.log.Info("successfully create bi-di stream for chat")

	errSystemChan := make(chan error, 1)
	defer close(errSystemChan)

	messagesChan := gw.initMessages()

	var wg sync.WaitGroup

	// Reciever goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, err := st.Recv()

			if err != nil {
				if err == io.EOF {
					gw.log.Warn("server stop sending messages")
				} else if err := errorRecieveHandling(gw.log, err); err != nil {
					errSystemChan <- err
				}
				return
			}
			time.Sleep(3 * time.Second)
		}
	}()

	// Sender goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer st.CloseSend()

		for msg := range messagesChan {
			if err := st.Send(msg); err != nil {
				errorSendHandling(gw.log, err)
				break
			}
			time.Sleep(3 * time.Second)
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-chatCtx.Done():
		gw.log.Info("context done")
	case <-done:
		gw.log.Info("reciever and sender has been stopped")
	case err = <-errSystemChan:
		gw.log.Error("chat error", logger.NewField("error", err))
	}
	gw.log.Info("bidirectional stream finished")
}

func (gw *noteGateway) initMessages() chan *pb.Message {
	messages := []*pb.Message{
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "tiger"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "wolf"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "bear"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "fox"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "deer"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "otter"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "penguin"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "dolphin"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "madam"},
		},
		{
			CorrelationId: int64(rand.Intn(1000) + 1),
			Payload:       &pb.Message_Text{Text: "eagle"},
		},
	}
	messagesChan := make(chan *pb.Message, len(messages))
	defer close(messagesChan)

	for _, msg := range messages {
		messagesChan <- msg
	}
	return messagesChan
}
