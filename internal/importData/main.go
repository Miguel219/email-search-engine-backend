package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"

	services "email-search-engine-backend/internal/server/services"

	"github.com/iancoleman/strcase"
)

const numCPU = 4

const index = "emails"

const directory = "./enron_mail_20110402/maildir"
const outputDirectory = "./emails"

var keys = [...][]byte{
	[]byte("Message-ID"),
	[]byte("Subject"),
	[]byte("X-From"),
	[]byte("From"),
	[]byte("To"),
	[]byte("Date"),
}

const bodySeparator = "\n\n"

// Formatea un string para guardarlo en archivos ndjson
func formatString(str string) (res string) {
	res = strings.Replace(str, "\\", "\\\\", -1)
	res = strings.Replace(res, "\"", "\\\"", -1)
	res = strings.Replace(res, "\n", "\\n", -1)
	return
}

// Verifica si un directorio esta vacío
func isDirEmpty(directory string) (bool, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return true, err
	}
	return len(files) == 0, err
}

// Lee la información de un archivo y separa en partes el email
func readDataFromFile(directory string) (data string) {

	// V2 - No se guarda toda la información, solo:
	// {
	// 	messageID
	// 	subject
	// 	xFrom
	// 	from
	// 	to
	// 	date
	// 	body
	// }
	// Se utiliza "ioutil.ReadFile" en lugar de "os.Open"
	file, err := ioutil.ReadFile(directory)
	if err != nil {
		fmt.Println(err)
	}
	headerFile, bodyFile, _ := bytes.Cut(file, []byte(bodySeparator))
	data = "{ "

	lenHeaderFile := len(headerFile)
	for _, key := range keys {
		i := bytes.Index(headerFile, key)
		if i != -1 {
			var x int = i + len(key) + 2
			var y int
			for y = x; y < lenHeaderFile; y++ {
				if headerFile[y] == byte('\n') {
					y--
					break
				}
			}

			var value string = ""
			if y > x {
				value = string(headerFile[x:y])
			}
			data += "\"" + strcase.ToLowerCamel(string(key)) +
				"\": \"" + formatString(value) +
				"\", "
		}
	}

	data += "\"body\": \"" + formatString(string(bodyFile)) + "\" }\n"

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

	fmt.Println("Total emails:", len(directories))

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
		resp, err := services.IndexEmails(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(resp))
	}
}

func ImportData() {
	//Profiling
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//Si no se ha generado ningún archivo
	if isEmpty, _ := isDirEmpty(outputDirectory); isEmpty {
		//Se genera la data
		createData(directory)
	}

	//Se guarda la data en ZincSearch
	saveData()
}
