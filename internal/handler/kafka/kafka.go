package kafka

import (
	"analytic-service/internal/logger"
	"errors"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type (
	Header struct {
		Name  string
		Value []byte
	}
	Message struct {
		Key     string
		Value   []byte
		Headers []Header
	}
	Handler = func(ctx context.Context, message Message) error

	TopicSet struct {
		State  string
		Event  string
		Notify string
	}

	KafkaClient struct {
		Wg        sync.WaitGroup
		cancelCtx context.Context
		cancelFn  context.CancelFunc
		groupId   string
		brokers   []string
		logger    *logger.Logger
		handler   *HandlerKafka
		topicSet  TopicSet
	}
)

func (c *KafkaClient) handle(ctx context.Context, m kafka.Message, fn Handler) error {
	headers := make([]Header, len(m.Headers))
	for i, header := range m.Headers {
		headers[i] = Header{Name: header.Key, Value: header.Value}
	}
	return fn(ctx, Message{Key: string(m.Key), Value: m.Value, Headers: headers})
}
func (c *KafkaClient) Consume(ctx context.Context, topic string, fn Handler) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		Topic:    topic,
		GroupID:  c.groupId,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})
	c.Wg.Add(1)
	go func(r *kafka.Reader, fn Handler) {
		defer c.Wg.Done()
		for c.cancelCtx.Err() != nil {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}

			//Не коммитимся о получении, пока не убедились, что логика на стороне сервиса отработала
			err = c.handle(ctx, m, fn)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}

			err = r.CommitMessages(ctx, m)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}
		}
		err := r.Close()
		if err != nil {
			c.logger.Error("ошибка в момент чтения топика кафкой ", zap.String("topic", r.Config().Topic))
		}
	}(reader, fn)
	return nil
}

func (c *KafkaClient) Dispose() {
	c.cancelFn()
	c.Wg.Wait()
}

//Создаем топики, если их еще нет
//Проверить наличие в списке через ps
//docker run -d --rm -p 9000:9000 -e KAFKA_BROKERCONNECT=91.185.95.87:9094 -e JVM_OPTS="-Xms32M -Xmx64M" -e SERVER_SERVLET_CONTEXTPATH="/" obsidiandynamics/kafdrop
func (c *KafkaClient) InitTopics(ctx context.Context) error {
	controllerConn, err := kafka.Dial("tcp", c.brokers[0])
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicCfgs := make([]kafka.TopicConfig, 3)
	topicCfgs[0] = kafka.TopicConfig{Topic: c.topicSet.State, NumPartitions: 1, ReplicationFactor: 1}
	topicCfgs[1] = kafka.TopicConfig{Topic: c.topicSet.Event, NumPartitions: 1, ReplicationFactor: 1}
	topicCfgs[2] = kafka.TopicConfig{Topic: c.topicSet.Notify, NumPartitions: 1, ReplicationFactor: 1}
	return controllerConn.CreateTopics(topicCfgs...)
}

func (c *KafkaClient) StartReaders(ctx context.Context) {
	var err error

	err = c.InitTopics(ctx)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	err = c.Consume(ctx, c.topicSet.State, c.handler.ProcessState)
	if err != nil {
		c.logger.Error(err.Error())
	}

	err = c.Consume(ctx, c.topicSet.Event, c.handler.ProcessEvent)
	if err != nil {
		c.logger.Error(err.Error())
	}

	err = c.Consume(ctx, c.topicSet.Notify, c.handler.ProcessNotify)
	if err != nil {
		c.logger.Error(err.Error())
	}
}

func NewKafka(brokers []string, groupId string, logger *logger.Logger, handler *HandlerKafka, topicSet TopicSet) (*KafkaClient, error) {
	if len(brokers) == 0 || brokers[0] == "" || groupId == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}
	ctx, fn := context.WithCancel(context.Background())
	c := KafkaClient{groupId: groupId, logger: logger, brokers: brokers, handler: handler, cancelCtx: ctx, cancelFn: fn, topicSet: topicSet}
	return &c, nil
}
