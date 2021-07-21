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
	task chan *w_pb.TaskRequest
	cond *sync.Cond
}

func NewWorker() *Worker {
	mr := &Worker{
		task: make(chan *w_pb.TaskRequest, 10),
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
	if len(request.Ip) == 0 || len(request.Data) == 0 {
		response.Code = -1
		response.Message = fmt.Sprintf("the param is invalid, Module [%s] Data [%s]", request.Ip, request.Data)
		return response, errors.New("the param is invalid")
	}
	worker.task <- &w_pb.TaskRequest{
		RequestId: request.RequestId,
		Data:      request.Data,
		Ip:        request.Ip,
	}
	worker.cond.Broadcast()
	response.Message = "success"
	return response, nil
}

func (worker *Worker) Schedule() {
	for {
		select {
		case tk := <-worker.task:
			for i := 0; i < len(tk.Data); i++ {
				//TODO agent去执行
				fmt.Println(tk.Ip, tk.Data[i])
			}
		default:
			fmt.Println("default")
			worker.Lock()
			// TODO 这里注意，Wait方法是直接调用unlock的，不加锁使用会报错
			worker.cond.Wait()
			worker.Unlock()
		}
	}
}
