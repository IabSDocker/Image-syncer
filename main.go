package main

import (
    "flag"
    "fmt"
    "image-syncer/utils"
    "image-syncer/syncer"
    "os"
)

func main() {
    // Define command-line arguments
    configPath := flag.String("config", "config/config.yaml", "Path to the configuration file")
    flag.Parse()

    // Load configuration file
    config, err := utils.LoadConfig(*configPath)
    if err != nil {
        fmt.Printf("Error loading config: %v\n", err)
        os.Exit(1)
    }

    // Create logger
    logger := utils.CreateLogger()

    // Create and start syncer
    syncer := syncer.NewSyncer(config, logger)
    syncer.StartSync()
}

