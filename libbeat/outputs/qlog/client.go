package qlog

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/outputs/codec"
	"github.com/elastic/beats/libbeat/outputs/outil"
	"github.com/elastic/beats/libbeat/publisher"
)

type client struct {
	observer outputs.Observer
	hosts    []string
	topic    outil.Selector
	index    string
	codec    codec.Codec
}

func newQlogClient(
	observer outputs.Observer,
	hosts []string,
	index string,
	topic outil.Selector,
	codec codec.Codec,
) (*client, error) {
	c := &client{
		observer: observer,
		hosts:    hosts,
		topic:    topic,
		index:    index,
		codec:    codec,
	}
	return c, nil
}

func (c *client) Close() error {
	Close()
	return nil
}

func (c *client) Publish(batch publisher.Batch) error {
	st := c.observer
	events := batch.Events()
	st.NewBatch(len(events))
	dropped := 0
	sucessed := 0
	for i := range events {
		ok := c.publishEvent(&events[i])
		if !ok {
			dropped++
			log.Println("publishing faild:\n", dropped, &events[i])
		} else {
			sucessed++
			log.Println("sucessed:", sucessed)
		}
	}
	batch.ACK()
	st.Dropped(dropped)
	st.Acked(len(events) - dropped)
	return nil
}

func (c *client) Publishs(batch publisher.Batch) error {
	st := c.observer
	events := batch.Events()
	st.NewBatch(len(events))
	rest, err := c.publishEvents(events)
	if len(rest) == 0 {
		batch.ACK()
	} else {
		batch.RetryEvents(rest)
	}
	return err
}

func (c *client) publishEvents(data []publisher.Event) ([]publisher.Event, error) {

	if len(data) == 0 {
		return nil, nil
	}
	var topic string
	var buffer bytes.Buffer
	for _, event := range data {
		serializedEvent, err := c.codec.Encode(c.index, &event.Content)
		if err != nil {
			return nil, err
		}
		if config.Container {
			topic, err = c.getEventTopic(&event)
			if err != nil {
				return data, err
			}
			controllerKind, err := c.getEventControllerKind(&event)
			if err != nil {
				return nil, err
			}
			if strings.ToLower(controllerKind) != "deployment" {
				topic = strings.Replace(topic, "k8s", "qlog", -1)
			} else {
				topic = strings.Replace(topic, "k8s_docker", "qlog", -1)
			}

		} else {
			topic = config.Topic
		}

		buffer.Write(serializedEvent)

	}
	msg := buffer.Bytes()

	qlog := NewQlog(
		hostName,
		topic,
		20000,
		config.IDC,
		"",
		"0",
		time.Now(),
		"filebeat",
		1,
		msg,
	)
	qlog.Pack()

SENDDATA:
	conn := Get()
	_, connErr := conn.Write(qlog.GetBin())
	if connErr != nil {
		Drop(conn)
		goto SENDDATA
	}
	Put(conn)

	return nil, nil

}

func (c *client) publishEvent(event *publisher.Event) bool {
	var err error
	var topic string
	if config.Container {
		topic, err = c.getEventTopic(event)
		if err != nil {
			return false
		}
		controllerKind, err := c.getEventControllerKind(event)
		if err != nil {
			return false
		}
		if strings.ToLower(controllerKind) != "deployment" {
			topic = strings.Replace(topic, "k8s", "qlog", -1)
		} else {
			topic = strings.Replace(topic, "k8s_docker", "qlog", -1)
		}

		err = c.makeEventTopic(event, topic)
		if err != nil {
			return false
		}

		//log.Println("发送日志对象为:", controllerKind, "本次日志发送topic为:", topic)

	} else {
		topic = config.Topic
	}

	serializedEvent, err := c.codec.Encode(c.index, &event.Content)
	if err != nil {
		return false
	}

	qlog := NewQlog(
		hostName,
		topic,
		20000,
		config.IDC,
		"",
		"0",
		time.Now(),
		"filebeat",
		1,
		serializedEvent,
	)
	qlog.Pack()

SENDDATA:
	conn := Get()
	_, connErr := conn.Write(qlog.GetBin())
	//log.Println("Conn status :", conn.RemoteAddr(), conn, connStatus, connErr)
	if connErr != nil {
		log.Println("Conn error :", conn.RemoteAddr(), connErr)
		Drop(conn)
		goto SENDDATA
	}
	//log.Println("Send qlog data :", string(qlog.GetBin()))
	//log.Println("Send data to qlog server succeed :", conn.RemoteAddr())
	Put(conn)

	return true
}

func (c *client) getEventTopic(data *publisher.Event) (string, error) {
	var topic string
	event := &data.Content
	if event.Fields != nil {
		if value, ok := event.Fields["topic"]; ok {
			if v, ok := value.(string); ok {
				topic = v
			}
		}
	}
	return topic, nil
}

func (c *client) makeEventTopic(data *publisher.Event, newTopic string) error {
	event := &data.Content
	_, err := c.topic.Select(event)
	if err != nil {
		return fmt.Errorf("setting qlog topic failed with %v", err)
	}
	if newTopic == "" {
		return fmt.Errorf(" qlog newTopic empty", err)
	}
	if event.Fields == nil {
		event.Fields = map[string]interface{}{}
	}
	event.Fields["topic"] = newTopic
	return nil

}

func (c *client) getEventControllerKind(data *publisher.Event) (string, error) {
	var controllerKind string
	event := &data.Content
	if event.Fields != nil {
		if value, ok := event.Fields["controller_kind"]; ok {
			if v, ok := value.(string); ok {
				controllerKind = v
			}
		}
	}

	return controllerKind, nil
}

func (out *client) String() string {
	return "qlog(" + out.index + ")"
}
