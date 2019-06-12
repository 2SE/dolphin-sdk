package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	addr := r.FormValue("addr")
	if addr == "" {
		addr = defaultDolphinAddr
	}
	err := registerManager.RegisterServerOnDolpin(addr)
	if err != nil {
		logrus.Errorf("Registration failed because: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		logrus.Info("Registered successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Registered successfully!!!"))
	}
}

func withRetryHttp(port string) {
	go func() {
		addr := "0.0.0.0:"
		if port == "" {
			port = "9999"
		}
		addr = addr + port
		http.HandleFunc("/register", handler)
		logrus.Infof("dolphin register retry prot ready to start and listen on %s", addr)
		server := &http.Server{
			Addr:         addr,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}
		err := server.ListenAndServe()
		if err != nil {
			panic(fmt.Errorf("ListenAndServe: %v", err))
		}
	}()
}
