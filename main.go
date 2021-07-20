package main

import (
	"github.com/feitianlove/worker/config"
	"github.com/feitianlove/worker/logger"
	"github.com/feitianlove/worker/service/worker"
	"github.com/sirupsen/logrus"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	conf, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = logger.InitLog(conf)
	logger.CtrlLog.WithFields(logrus.Fields{}).Info("init log success")
	if err != nil {
		logger.CtrlLog.WithFields(logrus.Fields{
			"err": err,
		}).Error("init log fail")
		panic(err)
	}
	//运行worker
	w := worker.NewWorker()
	go func() {
		w.Schedule()
	}()
	wg.Add(1)
	go func() {
		err := worker.RunWorker(conf, w)
		if err != nil {
			logger.CtrlLog.WithFields(logrus.Fields{
				"Error": err,
			}).Error("worker run fail")
			panic(err)
		}
		logger.CtrlLog.WithFields(logrus.Fields{}).Info("worker run success")
	}()
	wg.Wait()
}
