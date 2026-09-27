package main

import (
	"fmt"
	"net/http"
	"os" // para ler arquivo
	"time"

	"github.com/goccy/go-yaml" // lê arquivo yaml e converte em formato que go entende
)

type Targets struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type Config struct {
	TimeoutSeconds int       `yaml:"timeout_seconds"`
	Targets        []Targets `yaml:"targets"`
}

func loadYaml(filename string) (Config, error) { // devolve os arquivos ou erro caso falhe {

	var config Config
	//em go, uma função pode retornar mais de um valor; o primeiro guarda os dados
	// e o segundo guarda o erro, logo em seguida deve-se realizar tratamento de erro
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil

}

func main() {
	// Carrega as configurações do arquivo
	config, err := loadYaml("config.yaml")
	if err != nil {
		fmt.Printf("Erro ao carregar arquivo YAML: %v\n", err)
		return
	}

	// Percorre a lista de alvos obtida do YAML
	for _, target := range config.Targets {
		fmt.Printf("Testando: %s...\n", target.Name)
		endpoint(target.Url)
	}
}

func endpoint(url string) {
	// Cria um cliente HTTP com limite de tempo de 5 segundos
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	// Dispara a requisição GET
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("[ERRO] Falha ao acessar %s: %v\n", url, err)
		return
	}
	// Garante o fechamento do corpo da resposta após terminar a função
	defer resp.Body.Close()

	// Imprime o resultado formatado no terminal
	fmt.Printf("[OK] URL: %s - Status: %d\n", url, resp.StatusCode)
}
