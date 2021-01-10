package config

//TODO: better way https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
// https://eltonminetto.dev/en/post/2018-06-25-golang-usando-build-tags/

//GeneralConfig GeneralConfig
type GeneralConfig struct {
	DatabaseHost string
	DatabaseName string
	APIPort      string
}

//DevConfig DevConfig
var DevConfig = GeneralConfig{DatabaseHost: "mongodb://localhost:27017", DatabaseName: "products-service", APIPort: ":8080"}
