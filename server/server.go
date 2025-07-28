package server

import (
	"Metamorphoun/config"
	"Metamorphoun/enum"
	"Metamorphoun/shared"
	"Metamorphoun/zutil"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"

	//"net/http/pprof"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings" // Import the pprof package explicitly

	// Import the pprof package explicitly
	//_ "net/http/pprof"

	"github.com/tidwall/gjson"
)

var GetFolderPath func(string) string

type PathLocType string

type Data struct {
	Name  string `json:"parameter"`
	Value string `json:"value"`
}

func Serve(cfg config.Config) bool { //serverUrl string, serverPort int
	//mux := http.NewServeMux()
	// Register API handlers first
	http.HandleFunc("/inputApi", formApi)
	http.HandleFunc("/configApi", configApi)
	http.HandleFunc("/imagesFieldChangeApi", imagesFieldChangeApi)
	http.HandleFunc("/textFieldChangeApi", textFieldChangeApi)
	http.HandleFunc("/openLocationApi", openLocationApi)
	http.HandleFunc("/localFontApi", localFontApi)
	http.HandleFunc("/addImagesField", addImagesField)
	http.HandleFunc("/editImagesField", editImagesField)
	http.HandleFunc("/currentInfoApi", currentInfoApi)
	http.HandleFunc("/fileUploadForm", fileUploadForm)
	//http.HandleFunc("/uploadFile", uploadFile)

	// Register pprof handlers on the custom mux
	// mux.HandleFunc("/debug/pprof/", pprof.Index)
	// mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// For heap profile, you might want to use the 'heap' name explicitly
	//mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	// Create a sub-FS rooted at "static" inside the embedded FS
	staticFiles, err := fs.Sub(shared.StaticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub FS: %v", err)
	}
	fmt.Println("staticFiles:")
	fmt.Println(staticFiles)

	// Serve embedded static files
	http.Handle("/", http.FileServer(http.FS(staticFiles)))

	// Register static file server
	//fs := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fs)
	//mux.Handle("/", fs)

	log.Print("Listening on :3000...")
	log.Printf("Listening on %s:%d...", cfg.ServerAddress, cfg.ServerPort)
	//serverAddress := fmt.Sprintf("%s:%d", cfg.ServerAddress, cfg.ServerPort)
	//err := http.ListenAndServe(serverAddress, mux)
	err = http.ListenAndServe(":3000", nil)
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

// Save Changes made on web page for configuration
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
	configPath := config.GetFolderPath(enum.PathLoc.ConfigFile)
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
	fmt.Printf("Made it to openLocationApi\r\n")

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

// func OpenFolder(title string, path string) error {
// 	cmd := exec.Command(title, path) // Assign the command to the cmd variable
// 	err := cmd.Run()                 // Run the command
// 	if err != nil {
// 		fmt.Println("Error opening folder:", err)
// 		return err
// 	}
// 	return nil
// }

func OpenFolder(title string, path string) error {
	urlRegex := `^(http|https)://`
	matched, err := regexp.MatchString(urlRegex, path)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		if matched {
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
		} else {
			cmd = exec.Command("explorer", path)
		}
	case "darwin":
		// macOS uses `open` for both folders and URLs
		cmd = exec.Command("open", path)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error opening folder or URL: %v\n", err)
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
	fontPath := GetFolderPath(enum.PathLoc.Fonts)

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
	fmt.Printf("Made it to addImagesField\r\n")

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
func editImagesField(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Made it to editImagesField\r\n")

	// Read the request body
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // Close the body to prevent resource leaks

	fmt.Println("formApi-Received JSON:", string(jsonData))

	useResult := gjson.GetBytes(jsonData, "use").Bool()
	_ = useResult
	nameResult := gjson.GetBytes(jsonData, "name").String()
	// Check if the nameResult is empty
	if nameResult == "" {
		http.Error(w, "Name field cannot be empty", http.StatusBadRequest)
		return
	}
	imageItem := config.GetImageByName(nameResult)
	_ = imageItem
	titleResult := gjson.GetBytes(jsonData, "title").String()
	// locationResult := gjson.GetBytes(jsonData, "location").String()
	// operationResult := gjson.GetBytes(jsonData, "operation").String()

	_ = titleResult

	//result := config.editImagesField(useResult, nameResult, titleResult, locationResult, operationResult)
	//fmt.Println(result)

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

func currentInfoApi(w http.ResponseWriter, r *http.Request) {
	// Read the JSON file
	var rtnJson = config.ConfigInstance.PicHistories[0]

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Write JSON data to response
	jsonBytes, err := json.Marshal(rtnJson)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		fmt.Println("Failed to write response:", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Define your request structure
type FileUploadRequest struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Creators string `json:"creators"`
	Citation string `json:"citation"`
	Info     string `json:"info"`
	FilePath string `json:"filePath"` // Assume this is the path to the file to be copied
}

// func uploadFile(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var request FileUploadRequest

// 	// Decode the JSON body into our struct
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	// Define the target directory path
// 	configFolderPath := filepath.Join("path_to_config_folder") // Update with your config path
// 	targetDir := filepath.Join(configFolderPath, "Quotes")
// 	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
// 		http.Error(w, "Failed to create directories", http.StatusInternalServerError)
// 		return
// 	}

// 	// Define the target file path
// 	fileName := filepath.Base(request.FilePath) // Extract just the fileName from the filePath
// 	targetFilePath := filepath.Join(targetDir, fileName)

// 	// Copy the file
// 	if err := zutil.CopyFile(request.FilePath, targetFilePath); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to copy file: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Update request with new location
// 	request.FilePath = targetFilePath

// 	// Here you would update `TextLibraries` of config to add that new record as needed.
// 	fmt.Println("File copied and updated:", request)

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode("File uploaded successfully")
// }

// Receives upload of form to add a new quote file
func fileUploadForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// 1. Validate required fields (except hidden)
	requiredFields := []string{"fileInput", "name", "title", "creators", "citation", "info"} // Example required fields
	missing := []string{}
	for _, field := range requiredFields {
		if field != "fileInput" {
			val := r.FormValue(field)
			if val == "" {
				missing = append(missing, field)
			}
		}
	}

	// 2. Handle file upload
	file, fileHeader, err := r.FormFile("fileInput")
	if err != nil {
		missing = append(missing, "fileInput")
	} else {
		defer file.Close()
	}

	// 3. If any required fields are missing, return error
	if len(missing) > 0 {
		http.Error(w, "Missing required fields: "+strings.Join(missing, ", "), http.StatusBadRequest)
		return
	}

	// 4. Read and validate file content as JSON
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading uploaded file", http.StatusInternalServerError)
		return
	}
	var js interface{}
	if err := json.Unmarshal(fileContent, &js); err != nil {
		http.Error(w, "Uploaded file is not valid JSON", http.StatusBadRequest)
		return
	}

	// 5. Save file
	quotesDir := filepath.Join(config.GetFolderPath(enum.PathLoc.Config), "Quotes")
	if err := os.MkdirAll(quotesDir, 0755); err != nil {
		http.Error(w, "Error creating quotes directory", http.StatusInternalServerError)
		return
	}
	outFilePath := filepath.Join(quotesDir, fileHeader.Filename)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		http.Error(w, "Error creating output file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	if _, err := outFile.Write(fileContent); err != nil {
		http.Error(w, "Error saving uploaded file", http.StatusInternalServerError)
		return
	}

	// 6. Set hidden fields to defaults, except location
	newLib := config.TextLibrary{
		Name:     r.FormValue("name"),
		Title:    r.FormValue("title"),
		Creators: r.FormValue("creators"),
		Citation: r.FormValue("citation"),
		Info:     r.FormValue("info"),
		Location: outFilePath,
		Use:      true,
	}
	config.ConfigInstance.TextLibraries = append(config.ConfigInstance.TextLibraries, newLib)
	hiddenDefaults := map[string]string{
		//"hidden1": "default1",
		//"hidden2": "default2",
		// Add more as needed
		"location": outFilePath, // Set location to saved file path
	}

	// You can now use hiddenDefaults as needed (e.g., save to DB, etc.)
	_ = hiddenDefaults
	fmt.Fprintf(w, "File uploaded and validated as JSON: %s", outFilePath)
}
