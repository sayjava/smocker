package server

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Thiht/smocker/server/config"
	"github.com/Thiht/smocker/server/handlers"
	"github.com/Thiht/smocker/server/services"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func createTLSConfig () (*tls.Config, error) {
	var certByte []byte
	var keyByte []byte
	var err error
	var cert tls.Certificate

	
	if certByte, err = ioutil.ReadFile("ssl/cert.pem"); err != nil {
		return nil, err
	}
	
	if keyByte, err = ioutil.ReadFile("ssl/key.pem"); err != nil {
		return nil, err
	}

	if cert, err = tls.X509KeyPair(certByte, keyByte); err != nil {
		return nil, err
	}

	return  &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func NewMockServer(cfg config.Config) (*http.Server, services.Mocks) {
	mockServerEngine := echo.New()
	persistence := services.NewPersistence(cfg.PersistenceDirectory)
	sessions, err := persistence.LoadSessions()
	if err != nil {
		log.Error("Unable to load sessions: ", err)
	}
	mockServices := services.NewMocks(sessions, cfg.HistoryMaxRetention, persistence)

	mockServerEngine.HideBanner = true
	mockServerEngine.HidePort = true
	mockServerEngine.Use(recoverMiddleware(), loggerMiddleware(), HistoryMiddleware(mockServices))

	handler := handlers.NewMocks(mockServices)
	mockServerEngine.Any("/*", handler.GenericHandler)

	mockServerEngine.Server.Addr = ":" + strconv.Itoa(cfg.MockServerListenPort)

	if cfg.TLSEnabled {
		if tlsConfig, err := createTLSConfig(); err != nil {
			log.Fatal(err)
		} else {
			mockServerEngine.Server.TLSConfig = tlsConfig
		}
	}

	return mockServerEngine.Server, mockServices
}
