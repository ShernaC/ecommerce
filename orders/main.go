package main

import "products/config"

func init() {
	config.ConnectDB()
	config.SyncDB()
}
