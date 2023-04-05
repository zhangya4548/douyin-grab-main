package logger

import (
	"douyin-grab/constv"
	"fmt"
	"time"
)

// getLogFilePath get the log file save path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", constv.RuntimeRootPath, constv.LogSavePath)
}

// getLogFileName get the save name of the log file
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		constv.LogSaveName,
		time.Now().Format(constv.TimeFormat),
		constv.LogFileExt,
	)
}
