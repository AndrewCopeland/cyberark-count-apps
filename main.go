package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	pasapi "github.com/infamousjoeg/cybr-cli/pkg/cybr/api"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/queries"
	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

var (
	pasHostname = os.Getenv("PAS_BASE_URL")
	pasUsername = os.Getenv("PAS_USERNAME")
	pasPassword = os.Getenv("PAS_PASSWORD")
	pasAuthType = os.Getenv("PAS_AUTH_TYPE")
	pasInsecure = os.Getenv("PAS_INSECURE")

	conjurLogin  = os.Getenv("CONJUR_AUTHN_LOGIN")
	conjurApiKey = os.Getenv("CONJUR_AUTHN_API_KEY")
)

func main() {

	client := pasapi.Client{
		BaseURL:     pasHostname,
		AuthType:    pasAuthType,
		InsecureTLS: strings.Contains(pasInsecure, "yes"),
	}
	err := client.Logon(requests.Logon{
		Username: pasUsername,
		Password: pasPassword,
	})

	if err != nil {
		log.Fatalf("%s", err)
	}

	users, err := client.ListUsers(&queries.ListUsers{})
	if err != nil {
		log.Fatalf("%s", err)
	}

	client.Logoff()

	config, err := conjurapi.LoadConfig()
	if err != nil {
		panic(err)
	}

	conjur, err := conjurapi.NewClientFromKey(config,
		authn.LoginPair{
			Login:  conjurLogin,
			APIKey: conjurApiKey,
		},
	)

	hosts, err := conjur.Resources(&conjurapi.ResourceFilter{
		Kind: "host",
	})

	hostApps := []string{}
	for _, host := range hosts {
		hostApps = append(hostApps, host["id"].(string))
		// fmt.Println(host["id"].(string))
	}

	// Get all of the app providers
	appProviders := []string{}
	for _, user := range users.Users {
		if user.UserType == "AppProvider" {
			appProviders = append(appProviders, user.Username)
			// fmt.Println(user.Username)
		}
	}

	fmt.Println("pas-count, conjur-count")
	fmt.Printf("%d, %d\n", len(appProviders), len(hostApps))
	fmt.Println()

	fmt.Println("id, type")
	for _, prov := range appProviders {
		fmt.Printf("%s, %s\n", prov, "AppProvider")
	}

	for _, host := range hostApps {
		fmt.Printf("%s, %s\n", host, "Conjur")
	}
}
