package main

import (
	tLog "com.csion/tlog"
)

func main() {
	log := tLog.GetTLog()
	log.Trace("Trace test...")
	log.Debug("Debug test...")
	log.Info("Info test...")
	log.Warn("Warn test...")
	log.Error("Error test...")

}
