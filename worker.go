package worker

import (
	"context"
	"errors"
	"fmt"
	w_pb "github.com/feitianlove/worker/rpc/worker/w_pb"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Worker struct {
	Lock sync.RWMutex
	task []*Task
}

type Task struct {
	RequestId string
	Data      string //需要安装的IP
	Module    string //需要装的module
}

func NewWorker() *Worker {
	return &Worker{
		Lock: sync.RWMutex{},
		task: make([]*Task, 0),
	}
}

func RunWorker(worker *Worker) error {
	listener, err := net.Listen("tcp", ":9003")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	w_pb.RegisterTaskServer(grpcServer, worker)
	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}
func (worker *Worker) DistributeTask(ctx context.Context, request *w_pb.TaskRequest) (*w_pb.TaskResponse, error) {
	worker.Lock.Lock()
	defer worker.Lock.Unlock()
	var response = &w_pb.TaskResponse{
		Code:    0,
		Message: "",
	}
	// 判断参数
	if len(request.Module) == 0 || len(request.Data) == 0 {
		response.Code = -1
		response.Message = fmt.Sprintf("the param is invalid, Module [%s] Data [%s]", request.Module, request.Data)
		return response, errors.New("the param is invalid")
	}
	worker.task = append(worker.task, &Task{
		RequestId: request.RequestId,
		Data:      request.Data,
		Module:    request.Module,
	})
	response.Message = "success"
	return response, nil
}
