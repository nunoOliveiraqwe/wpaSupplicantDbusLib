package wpaSuppDBusLib

import "sync"

var mapMutex sync.Mutex
var eapProviders = make(map[string]eapMethod)

func register(name string, provider eapMethod) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	eapProviders[name] = provider
}

func lookup(name string) *eapMethod {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	provider := eapProviders[name]
	return &provider
}

func eapRegistryGetMatchingKeys(nameSlice []string) []string {
	filteredSlice := make([]string, 0)
	for i := 0; i < len(nameSlice); i++ {
		_, exists := eapProviders[nameSlice[i]]
		if exists {
			filteredSlice = append(filteredSlice, nameSlice[i])
		}
	}
	return filteredSlice
}
