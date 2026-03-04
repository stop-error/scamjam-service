package main

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

)

func initLogger() (useLogFile bool, logFilePath string) {
	executable, err := os.Executable()
	if err != nil { 
		log.Error().Msg("Could not get root path of executable file! Logging will be console only." + err.Error())
		return false, ""
	} 

	currentLogPath := filepath.Dir(executable) + "\\scamjam-service.log"
	oldLogPath := filepath.Dir(executable) + "\\scamjam-service.old"
	err = runLogCleanup(currentLogPath, oldLogPath)
	if err != nil {
		log.Error().Msg("Error running log cleanup! logging will be console-only." + err.Error())
		return false, ""
	}

	return true, currentLogPath
}

func runLogCleanup(currentLogPath string, oldLogPath string) (error) { //TODO: Oh my god clean this up

	log.Info().Msg("Running log cleanup")

	if _, err := os.Stat(oldLogPath); err == nil {
		log.Info().Msg("Log cleanup: Deleting current .old file")
		err := os.Remove(oldLogPath)
		if err != nil {
			log.Error().Msg("Error deleting .old file!" + err.Error())
			return err
		}
	}

	if _, err := os.Stat(currentLogPath); err == nil {
		log.Info().Msg("Log cleanup: .log file becomes .old file")
		err = CopyFile(currentLogPath, oldLogPath)
		if err != nil {
			log.Error().Msg("Error copying .log file to .old file!" + err.Error())
			return err
		}
		log.Info().Msg("Log cleanup: deleting .log file")
		err = os.Remove(currentLogPath)
		if err != nil {
			log.Error().Msg("Error deleting .log file! Logs will not rotate correctly" + err.Error())
			return err
		}
	}	
	return nil
}
