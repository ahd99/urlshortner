package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	fmt.Println("test")
	test1()
	fmt.Println("------------------------- test 2 --------")
	test2()
}

// simplest use
func test1() {
	logger := zap.NewExample()
	//logger := zap.NewProduction()
	//logger := zap.NewDevelopment()
	defer logger.Sync()
	
	logger.Debug("THis is ,essage", 
		zap.String("strField", "ali"),
		zap.Int("weight", 3),
		zap.Duration("time", 10 * time.Second))	
	
	
	sugar := logger.Sugar()
	sugar.Infow("log message", "time", 3, "weight", 10)	// {"level":"info","msg":"log message","time":3,"weight":10}
	sugar.Infof("log message %d %d", 3 ,10)  // {"level":"info","msg":"log message 3 10"}
	sugar.Info("time: ", 3, "  weight: ", 10)  //{"level":"info","msg":"time: 3  weight: 10"}

}


// custom logger
func test2() {

	//encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	

	syncer := zapcore.AddSync(os.Stdout)
	syncer = zapcore.Lock(syncer)
	fmt.Printf("syncer type: %T\n", syncer)

	core := zapcore.NewCore(encoder, syncer, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))

	logger.Info("message 1", 
		zap.String("name", "ali"),
		zap.Int("age", 30))	

	logger.Error("message 2", 
	zap.String("err", "not found"),
	zap.Int("severity", 30))	

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i:= 1; i<10000; i++ {
			logger.Info("||||||||||||||||||||||||||||||||||||||||||||||||||")
		}
		wg.Done()
	}()

	go func() {
		for i:= 1; i<10000; i++ {
			logger.Info("--------------------------------------------")
		}
		wg.Done()
	}()

	wg.Wait()
}