// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// simple does nothing except block while running the service.
package main

import (
	"time"
	"os"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/kardianos/service"
)

var logger zerolog.Logger

type program struct{
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	logger.Info().Msg("Starting named pipe listeners...")
	
	logger.Info().Msg("Starting main loop...")
	go p.run()
	return nil
}
func (p *program) run() {

	ticker := time.NewTicker(1 * time.Second)
	for {
		logger.Info().Msg("Going to sleep for 1 second")
		select {
		case tm := <-ticker.C:

			logger.Info().Msg("Recieved tick on ticker channel at " + tm.String())

				
				

		case <-p.exit:
			logger.Info().Msg("scamjam-dns-watchdog has recieved exit signal!")
		
			ticker.Stop()
		}
	}

	// Do work here
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	close(p.exit)
	return nil
}

func main() {

	useLogFile, logPath := initLogger()
		if useLogFile == true {
			logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error opening log file! logging will be console only.")
				logger = zerolog.New(os.Stderr).With().Caller().Logger()
			} else {
				defer logFile.Close()
				logger = zerolog.New(zerolog.MultiLevelWriter(os.Stdout, logFile)).With().Caller().Logger()
			}
			
		}


	svcConfig := &service.Config{
		Name:        "scamjam-service",
		DisplayName: "ScamJam Service",
		Description: "Backround service for ScamJam protection",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	err = s.Run()
	if err != nil {
		log.Error().Msg(err.Error())
	}
}