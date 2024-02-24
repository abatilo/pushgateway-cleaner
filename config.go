package main

import "time"

const (
	FlagVerbose               = "verbose"
	FlagVerboseDefault        = false
	FlagTTL                   = "ttl"
	FlagTTLDefault            = 10 * time.Minute
	FlagSyncPeriod            = "sync-period"
	FlagSyncPeriodDefault     = 3 * time.Minute
	FlagPushgatewayURL        = "pushgateway-url"
	FlagPushgatewayURLDefault = "http://localhost:9091"

	FlagHelp        = "help"
	FlagHelpDefault = false
)
