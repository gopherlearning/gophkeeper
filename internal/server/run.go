package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gopherlearning/gophkeeper/internal"
	"github.com/gopherlearning/gophkeeper/internal/model"
	"github.com/gopherlearning/gophkeeper/internal/storage/local"
	proto "github.com/gopherlearning/gophkeeper/pkg/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Cmd struct {
	Config  string `short:"c" type:"existingfile" help:"Конфигурационный файл для сервера." default:"config.yaml"`
	Listen  string `yaml:"listen" help:"Порт GRPC сервера" default:":9765"`
	Storage struct {
		Path string `yaml:"path" help:"Путь к локальному хранилищу" default:"gopher-storage"`
		DSN  string `yaml:"dsn" help:"Адрес базы данных (не реализовано)"`
	} `embed:"" prefix:"storage." help:"Параметры хранилища"`
}

func (l *Cmd) Run(ctx *internal.Context) error {
	log.Info().
		Msg("я API сервер для менеджера паролей GophKeeper")

	storage, err := local.NewLocalStorage("secretsecretsecretsecret", ".gopherstorage", "")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterPublicServer(grpcServer, &Server{db: storage})

	lis, err := net.Listen("tcp", l.Listen)
	if err != nil {
		return err
	}
	defer lis.Close()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Err(err)
		}
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	<-terminate

	return nil
}

var _ proto.PublicServer = (*Server)(nil)

type Server struct {
	proto.UnimplementedPublicServer
	db      Repository
	clients sync.Map
}

func (s *Server) Updates(req *proto.Request, srv proto.Public_UpdatesServer) error {
	log.Debug().Msg(req.Owner)

	err := srv.Send(&proto.Secret{Name: "__ping"})

	if err != nil {
		log.Err(err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	updates := make(chan *proto.Secret, 10)
	s.clients.Store(fmt.Sprintf("%s:%s", req.Owner, uuid.NewString()), updates)

	list := s.db.ListKeys()

	reg := regexp.MustCompile(fmt.Sprintf(`^%s:.*`, req.Owner))
	for _, v := range list {
		if !reg.Match([]byte(v)) {
			continue
		}

		secret := s.db.Get(model.Secret{Name: v})
		if secret.Updated < req.Updated {
			continue
		}

		secret.Name = strings.ReplaceAll(secret.Name, req.Owner+":", "")
		err := srv.Send(&proto.Secret{Data: secret.Data, Name: secret.Name, Type: proto.SecretType(secret.Type), Updated: secret.Updated})

		if err != nil {
			log.Err(err)
		}
	}

	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ticker.C:
			err := srv.Send(&proto.Secret{Name: "__ping"})
			if err != nil {
				log.Err(err)
				return status.Error(codes.InvalidArgument, err.Error())
			}
		case <-srv.Context().Done():
			return nil
		case secret := <-updates:
			secret.Name = strings.ReplaceAll(secret.Name, req.Owner+":", "")
			err := srv.Send(secret)

			if err != nil {
				log.Err(err)
				return status.Error(codes.InvalidArgument, err.Error())
			}
		}
	}
}

func (s *Server) Update(ctx context.Context, req *proto.Secret) (*proto.Empty, error) {
	err := s.db.Update(model.Secret{
		Name:  fmt.Sprintf("%s:%s", req.Owner, req.Name),
		Owner: req.Owner,
		Data:  req.Data,
	})

	reg := regexp.MustCompile(fmt.Sprintf(`^%s:.*`, req.Owner))

	s.clients.Range(func(key, value any) bool {
		if reg.Match([]byte(key.(string))) {
			value.(chan *proto.Secret) <- req
		}

		return true
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.Secret) (*proto.Empty, error) {
	err := s.db.Remove(model.Secret{
		Name: fmt.Sprintf("%s:%s", req.Owner, req.Name),
	})

	reg := regexp.MustCompile(fmt.Sprintf(`^%s:.*`, req.Owner))

	s.clients.Range(func(key, value any) bool {
		if reg.Match([]byte(key.(string))) {
			value.(chan *proto.Secret) <- req
		}

		return true
	})

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &proto.Empty{}, nil
}
