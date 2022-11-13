package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 5
const delay = 5

func main() {
	introducao()

	for {
		comando := getComando()
		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			exibirLog()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Digite um comando válido")
			main()
		}
	}
}

func introducao() {
	var nome string
	fmt.Println("Olá, qual é o seu nome?")
	fmt.Scan(&nome)
	fmt.Println("Bem vindo, ao programa Sr", nome)
}

func getComando() int {
	var comando int
	fmt.Println("\n1 - Iniciar monitoramento")
	fmt.Println("2 - Exibir  Logs")
	fmt.Println("0 - Sair do Programa")
	fmt.Scan(&comando)
	return comando
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	e := os.Remove("log.txt")
	if e != nil {
		fmt.Println(e)
	}

	logs := []*http.Response{}

	sites := leSitesDoArquivo()

	for i := 0; i < monitoramento; i++ {
		for i, site := range sites {
			tests := testarSite(site, i)
			logs = append(logs, tests)
		}
		time.Sleep(delay * time.Second)
	}
}
func leSitesDoArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}
	arquivo.Close()

	return sites
}

func exibirLog() {
	fmt.Println("Exibindo logs....")
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	fmt.Println(string(arquivo))
}

func registrarLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

		fmt.Println(err)
	}
	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " Online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func testarSite(site string, i int) *http.Response {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("")
		fmt.Print((i + 1), "-", "Site: ", site, " funciona perfeitamente\n ")
		registrarLog(site, true)
	} else {
		fmt.Println("")
		fmt.Print((i + 1), "-", "Site: ", site, " não está funcionando\n ")
		registrarLog(site, false)
	}
	return resp
}
