package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

const numCPU = 12

const user = "admin"
const password = "Complexpass#123"

const host = "http://192.168.5.52:4080"
const index = "emails"
const _type = "_bulk"

const directory = "./enron_mail_20110402/maildir"
const outputDirectory = "./emails"

// Formatea un string para guardarlo en archivos ndjson
func formatString(str string) (res string) {
	res = strings.Replace(str, "\\", "\\\\", -1)
	res = strings.Replace(res, "\"", "\\\"", -1)
	return
}

// Verifica si un directorio esta vacío
func IsDirEmpty(directory string) (bool, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return true, err
	}
	return len(files) == 0, err
}

// Lee la información de un archivo y separa en partes el email
func readDataFromFile(directory string) (data string) {
	readFile, err := os.Open(directory)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	var notBody = true
	data = "{ "
	var body = "\"body\": \""
	for fileScanner.Scan() {
		if notBody {
			if strings.Contains(fileScanner.Text(), ": ") {
				var textArray = strings.Split(fileScanner.Text(), ": ")
				var key = strcase.ToLowerCamel(textArray[0])
				data += "\"" + key + "\": \"" + formatString(textArray[1]) + "\", "
				if strings.Contains(fileScanner.Text(), "X-FileName") {
					notBody = false
				}
			}
		} else {
			body += formatString(fileScanner.Text()) + "\\n"
		}
	}
	data += body + "\" }\n"

	readFile.Close()

	return
}

// Lee los archivos y los separa en chunks
func readFiles(directories []string, c chan []string) {
	var res []string
	data := ""
	for _, dir := range directories {
		data += "{ \"index\" : { \"_index\" : \"" + index + "\" } }\n" + readDataFromFile(dir)
		//Si la data ya es mayor que 90mb, separarla (Dado que el API acepta request de max. 100mb)
		if len(data)*4 > 90000000 {
			res = append(res, data)
			data = ""
		}
	}
	if len(data) > 0 {
		res = append(res, data)
	}
	c <- res
}

// Crea archivos con la información de los correos electrónicos y el formato adecuado para enviarlos al API
func createData(directory string) {
	var directories []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			//Se agrega el directorio el path de cada email
			directories = append(directories, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total de emails:", len(directories))

	//Se divide el array de todos los directorios en chunks según la cantidad de CPUs
	var dividedDirectories [][]string
	chunkSize := (len(directories) + numCPU - 1) / numCPU
	for i := 0; i < len(directories); i += chunkSize {
		end := i + chunkSize
		if end > len(directories) {
			end = len(directories)
		}
		dividedDirectories = append(dividedDirectories, directories[i:end])
	}

	//Se crean los canales para cada goroutine
	c := make(chan []string)
	for i := 0; i < numCPU; i++ {
		go readFiles(dividedDirectories[i], c)
	}

	ix := 0
	for i := 0; i < numCPU; i++ {
		//Se obtiene la información del sub-proceso
		data := <-c
		//Se guarda la data en un archivo
		for _, element := range data {
			err := ioutil.WriteFile(outputDirectory+"/emails_"+strconv.Itoa(ix)+".ndjson", []byte(element), 0644)
			if err != nil {
				log.Fatal(err)
			}
			ix += 1
		}
	}
}

// Se lee cada archivo creado y cada correo se indexa en ZincSearch
func saveData() {
	files, err := ioutil.ReadDir(outputDirectory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(outputDirectory + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		//Se realiza un POST a ZincSearch para indexar los correos electrónicos
		req, err := http.NewRequest("POST", host+"/api/"+_type, strings.NewReader(string(data)))
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(user, password)
		req.Header.Set("Content-Type", "application/x-ndjson")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		log.Println(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	}
}

func main() {
	//Profiling
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//Si no se ha generado ningún archivo
	if isEmpty, _ := IsDirEmpty(outputDirectory); isEmpty {
		//Se genera la data
		createData(directory)
	}

	//Se guarda la data en ZincSearch
	saveData()
}
