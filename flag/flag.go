package flag

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	GitlabAddr  string
	GitlabToken string
	SonarAddr   string
	SonarToken  string
)

type Config struct {
	GitlabAddr    string `yaml:"GitlabAddr"`
	GitlabToken   string `yaml:"GitlabToken"`
	SonarAddr     string `yaml:"SonarAddr"`
	SonarToken    string `yaml:"SonarToken"`
	PostgreSQLDsn string `yaml:"PostgreSQLDsn"`
}

var Configuration *Config

func Flag() {
	configFile := flag.String("c", "config/Configuration.yaml", "Path to YAML configuration file")

	//flag.StringVar(&GitlabAddr, "gitlabAddr", "", "gitlab address")
	//flag.StringVar(&GitlabToken, "gitlabToken", "", "gitlab token")
	//flag.StringVar(&SonarAddr, "sonarAddr", "", "sonar address")
	//flag.StringVar(&SonarToken, "sonarToken", "", "sonar token")

	flag.StringVar(&GitlabAddr, "gitlabAddr", "", "gitlab address")
	flag.StringVar(&GitlabToken, "gitlabToken", "", "gitlab token")
	flag.StringVar(&SonarAddr, "sonarAddr", "", "sonar address")
	flag.StringVar(&SonarToken, "sonarToken", "", "sonar token")

	flag.Parse()

	// Read configuration file
	var err error
	Configuration, err = readConfig(*configFile)
	if err != nil {
		log.Fatalf("Error reading Configuration file: %v", err)
	}
	fmt.Println("Configuration: ", *Configuration)

	// If command line arguments are provided, override the values from the configuration file
	if GitlabAddr != "" {
		Configuration.GitlabAddr = GitlabAddr
	}
	if GitlabToken != "" {
		Configuration.GitlabToken = GitlabToken
	}
	if SonarAddr != "" {
		Configuration.SonarAddr = SonarAddr
	}
	if SonarToken != "" {
		Configuration.SonarToken = SonarToken
	}
	//fmt.Println("Configuration: ", *Configuration)

	GitlabAddr = Configuration.GitlabAddr
	GitlabToken = Configuration.GitlabToken
	SonarAddr = Configuration.SonarAddr
	SonarToken = Configuration.SonarToken

}

func readConfig(filename string) (*Config, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	return &config, nil
}
