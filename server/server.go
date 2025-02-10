package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"Metamorphoun/config"
	"Metamorphoun/zutil"

	"github.com/tidwall/gjson"
)

type Data struct {
	Name  string `json:"parameter"`
	Value string `json:"value"`
}

func Serve(serverUrl string, serverPort int) bool {
	// Register API handlers first
	http.HandleFunc("/inputApi", formApi)
	http.HandleFunc("/configApi", configApi)
	http.HandleFunc("/imagesFieldChangeApi", imagesFieldChangeApi)
	http.HandleFunc("/textFieldChangeApi", textFieldChangeApi)
	http.HandleFunc("/openLocationApi", openLocationApi)
	http.HandleFunc("/localFontApi", localFontApi)
	http.HandleFunc("/addImagesField", addImagesField)

	// Register static file server
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request was handled by a registered handler
	if r.URL.Path != "/inputapi" && r.URL.Path != "/configApi" {
		http.NotFound(w, r)
		return
	}

	// Check if the response status code has been set
	if w.Header().Get("Content-Type") != "" {
		// Response has already been started
		return
	}

	// Handle other errors (e.g., 500 Internal Server Error)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// func errorHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if the request was handled by a registered handler
// 	if r.URL.Path != "/" {
// 		fmt.Println(r.Host + " Not Found")
// 		http.NotFound(w, r)
// 		return
// 	}

// 	// Check if the response status code has been set
// 	if w.Header().Get("Content-Type") != "" {
// 		fmt.Println("Response has already been started")
// 		// Response has already been started
// 		return
// 	}

// 	// Handle other errors (e.g., 500 Internal Server Error)
// 	fmt.Println("500 Internal Server Error")
// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// }

func formApi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Made it to formapi\r\n")

	// Read the request body
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body to prevent resource leaks

	fmt.Println("formApi-Received JSON:", string(jsonData))

	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	parameterName, ok := data["parameter"].(string)
	if !ok {
		http.Error(w, "Invalid request: 'parameter' field not found or invalid type", http.StatusBadRequest)
		return
	}
	fmt.Println("parameterName:")
	fmt.Println(string(parameterName))

	value, ok1 := data["value"].(interface{})
	if !ok1 {
		http.Error(w, "Invalid request: 'value' field not found or invalid type", http.StatusBadRequest)
		return
	}
	//zutil.AsBool(fmt.Sprintf("%v", value))

	fmt.Println("value:")
	fmt.Println(value)

	// fmt.Println("parameterName:", string(parameterName))
	// fmt.Println("value:", string(value))

	config.UpdateConfigField(parameterName, value)
	// Save the updated configuration
	err = config.SaveConfig(config.GetConfig())
	if err != nil {
		http.Error(w, "Failed to update configuration", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Made it to formApi end\r\n")
	jData, err := json.Marshal(jsonData)
	if err != nil {
		// handle error
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

}

//	func imagesFieldChangeApi(w http.ResponseWriter, r *http.Request) {
//		fmt.Printf("Made it to imagesFieldChangeApi\r\n")
//		jsonData, err1 := ioutil.ReadAll(r.Body)
//		if err1 != nil {
//		}
//		fmt.Println("Received JSON:", string(jsonData))
//		var data Data
//		json.NewDecoder(r.Body).Decode(&data)
//		fmt.Printf("Received data: %+v\n", data)
//		config.UpdateImagesField(data.Name, data.Value)
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte("Data received successfully"))
//	}
// func imagesFieldChangeApi(w http.ResponseWriter, r *http.Request) {
// 	// Read the request body
// 	jsonData, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close() // Close the body to prevent resource leaks

// 	// Print the received JSON for debugging
// 	fmt.Println("Received JSON:", string(jsonData))

// 	// Define a struct to hold the incoming JSON data
// 	var data struct {
// 		Parameter string `json:"parameter"`
// 		Value     bool   `json:"value"`
// 	}

// 	// Unmarshal the JSON data into the struct
// 	err = json.Unmarshal(jsonData, &data)
// 	if err != nil {
// 		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
// 		return
// 	}

// 	config.UpdateImagesField(data.Parameter, data.Value)

// 	// Respond with success
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Configuration updated successfully"))
// }

func imagesFieldChangeApi(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body to prevent resource leaks

	// Print the received JSON for debugging
	fmt.Println("Received JSON:", string(jsonData))

	// Parse the JSON data using gjson
	// parameter := gjson.GetBytes(jsonData, "parameter").String()
	// value := gjson.GetBytes(jsonData, "value").Bool()
	parameterResult := gjson.GetBytes(jsonData, "parameter").String()
	valueResult := gjson.GetBytes(jsonData, "value").Bool()

	// // Debug prints to check parsed values
	// fmt.Println("Parsed parameter result:", parameterResult)
	// fmt.Println("Parsed value result:", valueResult)

	// parameter := parameterResult.String()
	// value := valueResult.Bool()

	// Debug prints to check final values
	// fmt.Println("Final parameter:", parameter)
	// fmt.Println("Final value:", value)

	// Update the configuration based on the received parameter and value
	config.UpdateImagesField(parameterResult, valueResult)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration updated successfully"))
}

func textFieldChangeApi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Made it to textFieldChangeApi\r\n")

	// Read the request body
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body to prevent resource leaks

	fmt.Println("imagesFieldChangeApi-Received JSON:", string(jsonData))

	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	quotesName, ok := data["parameter"].(string)
	if !ok {
		http.Error(w, "Invalid request: 'parameter' field not found or invalid type", http.StatusBadRequest)
		return
	}

	fmt.Printf("pre value:\r\n")
	//fmt.Printf(data["value"].(bool))

	useValue, ok1 := data["value"].(bool)
	if !ok1 {
		http.Error(w, "Invalid request: 'value' field not found or invalid type", http.StatusBadRequest)
		return
	}
	zutil.AsBool(fmt.Sprintf("%v", useValue))
	config.UpdateQuotesField(quotesName, zutil.AsString(useValue))
	// Save the updated configuration
	err = config.SaveConfig(config.GetConfig())
	if err != nil {
		http.Error(w, "Failed to update configuration", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Made it to imagesFieldChangeApi end\r\n")
	jData, err := json.Marshal(jsonData)
	if err != nil {
		// handle error
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func configApi(w http.ResponseWriter, r *http.Request) {
	// Read the JSON file
	usr, err := user.Current()
	if err != nil {
		fmt.Println("failed to get user home directory: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".Metamorphoun", "config.json")
	// Read config file
	jsonData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("Failed to read config file:", err)
		http.Error(w, "Failed to read config file", http.StatusInternalServerError)
		return
	}
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON data to response
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func openLocationApi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Made it to formapi\r\n")

	// Read the request body
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body to prevent resource leaks

	fmt.Println("formApi-Received JSON:", string(jsonData))

	if string(jsonData) == "" {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("failed to get user home directory: %w", err)
		}

		home := usr.HomeDir
		// Open the folder
		errOF := OpenFolder("explorer", home)
		if errOF != nil {
			//http.Error(w, fmt.Sprintf("Error opening folder: %v", errOF), http.StatusInternalServerError)
			//return
			fmt.Printf("Error opening folder: %v", errOF)
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Folder opened successfully! Close this and go back to the form"))
	} else {
		var data map[string]interface{}
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		id, ok := data["id"].(string)
		if !ok {
			http.Error(w, "Invalid request: 'id' field not found or invalid type", http.StatusBadRequest)
			return
		}
		fmt.Println("id:")
		fmt.Println(string(id))
		loc, ok := data["loc"].(string)
		if !ok {
			http.Error(w, "Invalid request: 'loc' field not found or invalid type", http.StatusBadRequest)
			return
		}
		//loc = loc + "\\"
		fmt.Println("Loc:")
		fmt.Println(string(loc))
		usr, err := user.Current()
		if err != nil {
			fmt.Println("failed to get user home directory: %w", err)
		}

		if string(loc) == "" {
			loc = usr.HomeDir
		}
		// Open the folder
		errOF := OpenFolder("explorer", loc)
		if errOF != nil {
			//http.Error(w, fmt.Sprintf("Error opening folder: %v", errOF), http.StatusInternalServerError)
			//return
			fmt.Printf("Error opening folder: %v", errOF)
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Folder opened successfully!"))

	}
}
func OpenFolder(title string, path string) error {
	cmd := exec.Command(title, path) // Assign the command to the cmd variable
	err := cmd.Run()                 // Run the command
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return err
	}
	return nil
}

func localFontApi(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("Made it to localFontApi\r\n")

	// Read the request body
	//jsonData, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Error reading request body", http.StatusBadRequest)
	// 	return
	// }
	defer r.Body.Close() // Close the body to prevent resource leaks

	//fmt.Println("formApi-Received JSON:", string(jsonData))

	// Get the path from the configuration
	fontPath := config.ConfigInstance.TextFontPath

	// Get all font files in the specified path
	fontFiles, err := getFontFiles(fontPath)
	if err != nil {
		http.Error(w, "Error getting font files", http.StatusInternalServerError)
		return
	}

	// Print the font file names for debugging
	//fmt.Println("Found font files:", fontFiles)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Marshal the font file names to JSON
	jsonResponse, err := json.Marshal(fontFiles)
	if err != nil {
		http.Error(w, "Error marshaling font files to JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonResponse)
}

// getFontFiles returns a slice of all font file paths in the given directory
func getFontFiles(dir string) ([]string, error) {
	var fontFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".tts") || strings.HasSuffix(file.Name(), ".ttf") || strings.HasSuffix(file.Name(), ".otf")) {
			fontFiles = append(fontFiles, filepath.Join(dir, file.Name()))
		}
	}

	return fontFiles, nil
}

func getAllFilePaths(root string) ([]string, error) {
	var filePaths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), "jpg") ||
				strings.HasSuffix(info.Name(), "png") ||
				strings.HasSuffix(info.Name(), "bmp") ||
				strings.HasSuffix(info.Name(), "gif") {
				filePaths = append(filePaths, path)
			}
		}
		return nil
	})
	return filePaths, err
}

func addImagesField(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Made it to localFontApi\r\n")

	// Read the request body
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body to prevent resource leaks

	fmt.Println("formApi-Received JSON:", string(jsonData))

	useResult := gjson.GetBytes(jsonData, "use").Bool()
	nameResult := gjson.GetBytes(jsonData, "name").String()
	titleResult := gjson.GetBytes(jsonData, "title").String()
	locationResult := gjson.GetBytes(jsonData, "location").String()
	operationResult := gjson.GetBytes(jsonData, "operation").String()

	result := config.AddImagesField(useResult, nameResult, titleResult, locationResult, operationResult)
	fmt.Println(result)

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON data to response
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON data to response
	_, err = w.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
