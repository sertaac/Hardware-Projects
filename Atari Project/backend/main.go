/**************************************/
/*                                    */
/*     Backend Service Entry Point    */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"retro-gaming-ui/backend/library"
	"retro-gaming-ui/backend/server"
)

/**************************************************/
/*                                                */
/*                   CONSTANTS                    */
/*                                                */
/**************************************************/

const (
	IPC_PORT    = 9847
	CONFIG_FILE = "library.json"
)

/**************************************************/
/*                                                */
/*                 MAIN FUNCTION                  */
/*                                                */
/**************************************************/

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║    RETRO GAMING HUB - BACKEND          ║")
	fmt.Println("║    Frutiger Aero • Y2K Edition         ║")
	fmt.Println("╚════════════════════════════════════════╝")

	/*          Get config directory              */
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	configPath := filepath.Join(configDir, "retro-gaming-hub", CONFIG_FILE)

	/*           Initialize library               */
	lib := library.NewLibrary(configPath)
	fmt.Printf("Library loaded from: %s\n", configPath)

	/*           Create IPC server                */
	ipcServer := server.NewIPCServer(IPC_PORT)

	/*          Set up message handler            */
	ipcServer.SetHandler(func(req server.Request) server.Response {
		return handleRequest(lib, req)
	})

	/*              Start server                  */
	if err := ipcServer.Start(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}

	/*        Wait for shutdown signal            */
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\nBackend running. Press Ctrl+C to stop.")
	<-sigChan

	fmt.Println("\nShutting down...")
	ipcServer.Stop()
	lib.Save()
	fmt.Println("Goodbye!")
}

/**************************************************/
/*                                                */
/*               REQUEST HANDLER                  */
/*                                                */
/**************************************************/

func handleRequest(lib *library.Library, req server.Request) server.Response {
	switch req.Type {

	case server.MsgTypeListGames:
		var payload server.GameListPayload
		if req.Payload != nil {
			json.Unmarshal(req.Payload, &payload)
		}
		games := lib.GetGames(payload.Platform, payload.Category)
		return server.Response{
			Type: server.MsgTypeSuccess, ID: req.ID,
			Success: true, Data: games,
		}

	case server.MsgTypeGetGame:
		var id string
		json.Unmarshal(req.Payload, &id)
		game := lib.GetGameByID(id)
		if game == nil {
			return server.Response{
				Type: server.MsgTypeError, ID: req.ID,
				Success: false, Error: "Game not found",
			}
		}
		return server.Response{
			Type: server.MsgTypeSuccess, ID: req.ID,
			Success: true, Data: game,
		}

	case server.MsgTypeGetFavorites:
		return server.Response{
			Type: server.MsgTypeSuccess, ID: req.ID,
			Success: true, Data: lib.GetFavorites(),
		}

	case server.MsgTypeToggleFavorite:
		var id string
		json.Unmarshal(req.Payload, &id)
		if err := lib.ToggleFavorite(id); err != nil {
			return server.Response{
				Type: server.MsgTypeError, ID: req.ID,
				Success: false, Error: err.Error(),
			}
		}
		return server.Response{
			Type: server.MsgTypeSuccess, ID: req.ID, Success: true,
		}

	case server.MsgTypeScan:
		if err := lib.Scan(); err != nil {
			return server.Response{
				Type: server.MsgTypeError, ID: req.ID,
				Success: false, Error: err.Error(),
			}
		}
		return server.Response{
			Type: server.MsgTypeSuccess, ID: req.ID,
			Success: true, Data: fmt.Sprintf("Found %d games", len(lib.Games)),
		}

	case server.MsgTypeStatus:
		return server.Response{
			Type: server.MsgTypeStatus, ID: req.ID,
			Success: true, Data: map[string]interface{}{
				"status":  "ready",
				"version": "1.0.0",
			},
		}

	default:
		return server.Response{
			Type: server.MsgTypeError, ID: req.ID,
			Success: false, Error: "Unknown message type",
		}
	}
}
