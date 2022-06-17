package shoset

import (
	"time"

	"github.com/ditrit/shoset/msg"
	"github.com/rs/zerolog/log"
)

type EventHandler struct{}

// GetEvent :
func (eh *EventHandler) Get(c *ShosetConn) (msg.Message, error) {
	var evt msg.Event
	err := c.ReadMessage(&evt)
	return evt, err
}

// HandleEvent :
func (eh *EventHandler) HandleDoubleWay(c *ShosetConn, message msg.Message) error {
	evt := message.(msg.Event)
	if state := c.GetCh().Queue["evt"].Push(evt, c.GetRemoteShosetType(), c.GetLocalAddress()); state {
		eh.Send(c.ch, evt)
	}
	return nil
}

// SendEventConn :
func (eh *EventHandler) SendEventConn(c *ShosetConn, evt interface{}) {
	_, err := c.WriteString("evt")
	if err != nil {
		log.Warn().Msg("couldn't write string evt : " + err.Error())
		return
	}
	err = c.WriteMessage(evt)
	if err != nil {
		log.Warn().Msg("couldn't write message evt : " + err.Error())
		return
	}
}

// SendEvent :
func (eh *EventHandler) Send(c *Shoset, evt msg.Message) {
	c.ConnsByName.Iterate(
		func(key string, conn interface{}) {
			err := conn.(*ShosetConn).SendMessage(evt)
			if err != nil {
				log.Warn().Msg("couldn't send evt : " + err.Error())
			}
		},
	)
}

// WaitEvent :
func (eh *EventHandler) Wait(c *Shoset, replies *msg.Iterator, args map[string]string, timeout int) *msg.Message {
	topicName, ok := args["topic"]
	if !ok {
		return nil
	}
	eventName := args["event"]
	term := make(chan *msg.Message, 1)
	cont := true //??

	go func() {
		for cont {
			//Check message presence in two steps to avoid accessing attributs of <nil>
			cell := replies.Get()
			if cell == nil {
				time.Sleep(time.Duration(10) * time.Millisecond)
				continue
			}
			message := cell.GetMessage()
			if message == nil {
				time.Sleep(time.Duration(10) * time.Millisecond)
				continue
			}
			event := message.(msg.Event)
			if event.GetTopic() == topicName && (eventName == VOID || event.GetEvent() == eventName) {
				term <- &message
			}
		}
	}()
	select {
	case res := <-term:
		cont = false
		return res
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil
	}
}
