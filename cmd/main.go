package main

import (
	cfg "calendar/internal/config"
	"calendar/internal/interfaces"
	lg "calendar/internal/logger"
	pb "calendar/internal/proto"
	"calendar/internal/services"
	"google.golang.org/grpc"
	"net"
)

func main() {
	logger := lg.GetLogger(cfg.GetConfig())

	logger.Info("Service loading!")

	psql := interfaces.NewPSQL(logger, cfg.GetConfig())

	//создаем структуру
	sch := services.NewScheduler(logger, psql)

	//обьявляем TCP листенер на 5454 порту
	netListener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		logger.Error(err.Error())
	}

	//создаем grpc сервер и регистрируем его через функцию в прото файлике
	grpcServer := grpc.NewServer()
	pb.RegisterAPIServer(grpcServer, sch)

	logger.Info("Service started!")

	//связываем grpc сервер и tcp листенер. Затем запускаем сервер
	err = grpcServer.Serve(netListener)
	if err != nil {
		logger.Error(err.Error())
	}

}
