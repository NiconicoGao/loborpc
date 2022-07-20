package lobo

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Console *zap.Logger
	Kafka   *zap.Logger
}

func NewLogger() *Logger {
	var err error
	c := new(Logger)
	c.Console, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	writeSyncer := zapcore.AddSync(NewKafkaWriter())
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.EpochNanosTimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), writeSyncer, zapcore.InfoLevel)
	c.Kafka = zap.New(core)
	return c
}

func (c *Logger) Init() {
	var err error
	c.Console, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	writeSyncer := zapcore.AddSync(NewKafkaWriter())
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.EpochNanosTimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), writeSyncer, zapcore.InfoLevel)
	c.Kafka = zap.New(core)

}

func (c *Logger) Info(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	c.Console.Info(str)
	c.Kafka.Info(str)
}

func (c *Logger) Warn(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	c.Console.Warn(str)
	c.Kafka.Warn(str)
}

func (c *Logger) Error(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	c.Console.Error(str)
	c.Kafka.Error(str)
}

func (c *Logger) Panic(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	c.Console.Panic(str)
}
