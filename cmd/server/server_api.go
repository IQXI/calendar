package main

import (
	cfg "calendar/internal/config"
	"calendar/internal/interfaces/postgres"
	lg "calendar/internal/logger"
	pb "calendar/internal/proto"
	"calendar/internal/services"
	"google.golang.org/grpc"
	"net"
)

func main() {
	logger := lg.GetLogger(cfg.GetConfig())

	logger.Info("Service loading!")

	psql, err := postgres.NewPSQL(logger, cfg.GetConfig())
	if err != nil {
		logger.Error(err.Error())
	}
	defer psql.Close()

	//создаем структуру
	sch := services.NewAPI(logger, psql)

	//обьявляем TCP листенер на 50051 порту
	netListener, err := net.Listen("tcp", ":50051")
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
