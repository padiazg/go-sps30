package sps30

import logger "github.com/d2r2/go-logger"

// You can manage verbosity of log output
// in the package by changing last parameter value.
var lg = logger.NewPackageLogger("sps30",
	logger.DebugLevel,
	// logger.InfoLevel,
)
