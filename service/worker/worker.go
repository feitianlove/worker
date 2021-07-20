package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/feitianlove/worker/config"
	w_pb "github.com/feitianlove/worker/rpc/worker/w_pb"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Worker struct {
	sync.Mutex
	task chan *Task
	cond *sync.Cond
}
type Task struct {
	RequestId string
	Data      string //需要安装的IP
	Module    string //需要装的module
}

func NewWorker() *Worker {
	mr := &Worker{
		task: make(chan *Task, 10),
	}
	mr.cond = sync.NewCond(mr)
	return mr
}

func RunWorker(conf *config.Config, worker *Worker) error {
	addr := fmt.Sprintf("%s:%d", conf.Worker.Domain, conf.Worker.ListenPort)
	fmt.Println(addr)
	listener, err := net.Listen("tcp", addr)
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
	worker.Lock()
	defer worker.Unlock()
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
	worker.task <- &Task{
		RequestId: request.RequestId,
		Data:      request.Data,
		Module:    request.Module,
	}
	worker.cond.Broadcast()
	response.Message = "success"
	return response, nil
}

func (worker *Worker) Schedule() {
	for {
		select {
		case tk := <-worker.task:
			fmt.Println(tk)
		default:
			fmt.Println("default")
			worker.Lock()
			// TODO 这里注意，Wait方法是直接调用unlock的，不加锁使用会报错
			worker.cond.Wait()
			worker.Unlock()
		}
	}
}
