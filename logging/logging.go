package logging

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger
var LogFile *os.File
 

func GetLoggerAsPointer() *zerolog.Logger {
	return &Logger
}

func InitLogger(logFileName string) error {
	executable, err := os.Executable()
	currentLogPath := filepath.Dir(executable) + "\\" + logFileName + ".log"
	oldLogPath := filepath.Dir(executable) + "\\" + logFileName + ".old"

	err = RunLogCleanup(currentLogPath, oldLogPath)
	if err != nil {
		log.Error().Err(err).Msg("Error running log cleanup! logging will be console-only.")
		return err
	}


	logFile, err := os.OpenFile(currentLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil { 
		log.Error().Err(err).Msg("Could not get root path of executable file! Logging will be console only.")
		Logger = zerolog.New(os.Stderr).With().Caller().Logger()
	} else {
		defer logFile.Close()
		Logger = zerolog.New(zerolog.MultiLevelWriter(os.Stdout, logFile)).With().Caller().Logger()
	}

	return nil
}


func RunLogCleanup(currentLogPath string, oldLogPath string) (error) { //TODO: Oh my god clean this up

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

func CopyFile(src string, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Error().Msg("Error opening file at " + src + "Error: " + err.Error())
		return err
	}

	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		log.Error().Msg("Error creating file at " + dst + "! Error: " + err.Error())
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(sourceFile, destFile)
	if err != nil {
		log.Error().Msg("Error copying file " + dst + "! Error: " + err.Error())
		return err
	}

	return nil

}