package steam

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/vvekic/go-steam/netutil"
)

// Load initial server list from Steam Directory Web API.
// Call InitializeSteamDirectory() before Connect() to use
// steam directory server list instead of static one.
func InitializeSteamDirectory(cellId int) error {
	return steamDirectoryCache.Initialize(cellId)
}

func UpdateSteamDirectory(servers []*netutil.PortAddr) {
	steamDirectoryCache.UpdateServerList(servers)
}

var steamDirectoryCache *steamDirectory = &steamDirectory{}

type steamDirectory struct {
	sync.RWMutex
	servers       []*netutil.PortAddr
	isInitialized bool
}

// Get server list from steam directory and save it for later
func (sd *steamDirectory) Initialize(cellId int) error {
	sd.Lock()
	defer sd.Unlock()
	client := new(http.Client)
	resp, err := client.Get(fmt.Sprintf("https://api.steampowered.com/ISteamDirectory/GetCMList/v1/?cellId=%d", cellId))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r := struct {
		Response struct {
			ServerList []string
			Result     uint32
			Message    string
		}
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	if r.Response.Result != 1 {
		return fmt.Errorf("Failed to get steam directory, result: %v, message: %v\n", r.Response.Result, r.Response.Message)
	}
	if len(r.Response.ServerList) == 0 {
		return fmt.Errorf("Steam returned zero servers for steam directory request\n")
	}
	sd.servers = []*netutil.PortAddr{}
	for _, s := range r.Response.ServerList {
		sd.servers = append(sd.servers, netutil.ParsePortAddr(s))
	}
	sd.isInitialized = true
	return nil
}

func (sd *steamDirectory) GetRandomCM() *netutil.PortAddr {
	sd.RLock()
	defer sd.RUnlock()
	if !sd.isInitialized {
		panic("steam directory is not initialized")
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	addr := sd.servers[rng.Int31n(int32(len(sd.servers)))]
	return addr
}

func (sd *steamDirectory) UpdateServerList(servers []*netutil.PortAddr) {
	sd.Lock()
	defer sd.Unlock()
	log.Printf("Updating Steam CM server list")
	sd.servers = servers
}

func (sd *steamDirectory) IsInitialized() bool {
	sd.RLock()
	defer sd.RUnlock()
	isInitialized := sd.isInitialized
	return isInitialized
}
