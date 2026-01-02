/**************************************/
/*                                    */
/*    Game Library Manager - Go       */
/*     Frutiger Aero + Y2K Edition    */
/*           Programmed by            */
/*            Sertaç Ataç             */
/*            02.01.2026              */
/*                                    */
/**************************************/

package library

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/**************************************************/
/*                                                */
/*             GAME INFO STRUCTURE                */
/*        Represents a game in the library        */
/*                                                */
/**************************************************/

type GameInfo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Platform    string    `json:"platform"`
	Path        string    `json:"path"`
	CoverPath   string    `json:"cover_path"`
	LastPlayed  time.Time `json:"last_played"`
	PlayCount   int       `json:"play_count"`
	Favorite    bool      `json:"favorite"`
	Category    string    `json:"category"`
}

/**************************************************/
/*                                                */
/*             LIBRARY STRUCTURE                  */
/*         Manages the game collection            */
/*                                                */
/**************************************************/

type Library struct {
	Games      []GameInfo     `json:"games"`
	ScanPaths  []string       `json:"scan_paths"`
	Categories map[string]int `json:"categories"`
	Platforms  map[string]int `json:"platforms"`
	LastScan   time.Time      `json:"last_scan"`
	mu         sync.RWMutex
	configPath string
}

/*    Supported ROM extensions by platform       */
var platformExtensions = map[string][]string{
	"NES":   {".nes", ".unf", ".unif"},
	"SNES":  {".sfc", ".smc"},
	"N64":   {".n64", ".z64", ".v64"},
	"GBA":   {".gba"},
	"GB":    {".gb", ".gbc"},
	"ATARI": {".a26", ".bin"},
}

/**************************************************/
/*                                                */
/*            LIBRARY CONSTRUCTOR                 */
/*                                                */
/**************************************************/

func NewLibrary(configPath string) *Library {
	lib := &Library{
		Games:      make([]GameInfo, 0),
		ScanPaths:  make([]string, 0),
		Categories: make(map[string]int),
		Platforms:  make(map[string]int),
		configPath: configPath,
	}
	lib.Load()
	return lib
}

/**************************************************/
/*                                                */
/*              ADD SCAN PATH                     */
/*                                                */
/**************************************************/

func (lib *Library) AddScanPath(path string) error {
	lib.mu.Lock()
	defer lib.mu.Unlock()

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return os.ErrNotExist
	}

	for _, p := range lib.ScanPaths {
		if p == path {
			return nil
		}
	}

	lib.ScanPaths = append(lib.ScanPaths, path)
	return nil
}

/**************************************************/
/*                                                */
/*         SCAN FOR GAMES IN PATHS                */
/*                                                */
/**************************************************/

func (lib *Library) Scan() error {
	lib.mu.Lock()
	defer lib.mu.Unlock()

	lib.Games = make([]GameInfo, 0)
	lib.Categories = make(map[string]int)
	lib.Platforms = make(map[string]int)

	for _, scanPath := range lib.ScanPaths {
		filepath.Walk(scanPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			platform := lib.detectPlatform(ext)

			if platform != "" {
				game := GameInfo{
					ID:       generateID(path),
					Title:    cleanGameTitle(filepath.Base(path)),
					Platform: platform,
					Path:     path,
					Category: "Uncategorized",
				}
				lib.Games = append(lib.Games, game)
				lib.Platforms[platform]++
				lib.Categories[game.Category]++
			}
			return nil
		})
	}

	lib.LastScan = time.Now()
	return lib.Save()
}

/**************************************************/
/*                                                */
/*          DETECT PLATFORM BY EXTENSION          */
/*                                                */
/**************************************************/

func (lib *Library) detectPlatform(ext string) string {
	for platform, extensions := range platformExtensions {
		for _, e := range extensions {
			if ext == e {
				return platform
			}
		}
	}
	return ""
}

/**************************************************/
/*                                                */
/*               GET GAMES                        */
/*                                                */
/**************************************************/

func (lib *Library) GetGames(platform, category string) []GameInfo {
	lib.mu.RLock()
	defer lib.mu.RUnlock()

	if platform == "" && category == "" {
		return lib.Games
	}

	filtered := make([]GameInfo, 0)
	for _, game := range lib.Games {
		if platform != "" && game.Platform != platform {
			continue
		}
		if category != "" && game.Category != category {
			continue
		}
		filtered = append(filtered, game)
	}
	return filtered
}

/**************************************************/
/*                                                */
/*             GET GAME BY ID                     */
/*                                                */
/**************************************************/

func (lib *Library) GetGameByID(id string) *GameInfo {
	lib.mu.RLock()
	defer lib.mu.RUnlock()

	for i := range lib.Games {
		if lib.Games[i].ID == id {
			return &lib.Games[i]
		}
	}
	return nil
}

/**************************************************/
/*                                                */
/*            TOGGLE FAVORITE                     */
/*                                                */
/**************************************************/

func (lib *Library) ToggleFavorite(id string) error {
	lib.mu.Lock()
	defer lib.mu.Unlock()

	for i := range lib.Games {
		if lib.Games[i].ID == id {
			lib.Games[i].Favorite = !lib.Games[i].Favorite
			return lib.saveUnlocked()
		}
	}
	return nil
}

/**************************************************/
/*                                                */
/*              GET FAVORITES                     */
/*                                                */
/**************************************************/

func (lib *Library) GetFavorites() []GameInfo {
	lib.mu.RLock()
	defer lib.mu.RUnlock()

	favorites := make([]GameInfo, 0)
	for _, game := range lib.Games {
		if game.Favorite {
			favorites = append(favorites, game)
		}
	}
	return favorites
}

/**************************************************/
/*                                                */
/*          GET RECENTLY PLAYED                   */
/*                                                */
/**************************************************/

func (lib *Library) GetRecentlyPlayed(limit int) []GameInfo {
	lib.mu.RLock()
	defer lib.mu.RUnlock()

	games := make([]GameInfo, len(lib.Games))
	copy(games, lib.Games)

	for i := 0; i < len(games)-1; i++ {
		for j := 0; j < len(games)-i-1; j++ {
			if games[j].LastPlayed.Before(games[j+1].LastPlayed) {
				games[j], games[j+1] = games[j+1], games[j]
			}
		}
	}

	if limit > 0 && limit < len(games) {
		return games[:limit]
	}
	return games
}

/**************************************************/
/*                                                */
/*              SAVE / LOAD                       */
/*                                                */
/**************************************************/

func (lib *Library) Save() error {
	lib.mu.Lock()
	defer lib.mu.Unlock()
	return lib.saveUnlocked()
}

func (lib *Library) saveUnlocked() error {
	data, err := json.MarshalIndent(lib, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(lib.configPath)
	os.MkdirAll(dir, 0755)
	return os.WriteFile(lib.configPath, data, 0644)
}

func (lib *Library) Load() error {
	lib.mu.Lock()
	defer lib.mu.Unlock()

	data, err := os.ReadFile(lib.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, lib)
}

/**************************************************/
/*                                                */
/*            HELPER FUNCTIONS                    */
/*                                                */
/**************************************************/

func generateID(path string) string {
	hash := uint32(0)
	for _, c := range path {
		hash = hash*31 + uint32(c)
	}
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if len(base) > 8 {
		base = base[:8]
	}
	return strings.ToUpper(base) + string(rune('A'+hash%26))
}

func cleanGameTitle(filename string) string {
	title := strings.TrimSuffix(filename, filepath.Ext(filename))
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")
	for _, pattern := range []string{"(USA)", "(Europe)", "(Japan)"} {
		title = strings.ReplaceAll(title, pattern, "")
	}
	return strings.TrimSpace(title)
}
