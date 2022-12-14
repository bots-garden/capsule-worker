package routes

import (
    "fmt"
    "github.com/bots-garden/capsule-worker/worker/helpers"
    "github.com/bots-garden/capsule-worker/worker/models"
    "github.com/gin-gonic/gin"
    "github.com/go-resty/resty/v2"
    "github.com/google/uuid"
    "net/http"
    "os"
    "os/exec"
    "strconv"
)

// AddRunningWasmModuleToRevision updates the functions map
func AddRunningWasmModuleToRevision(functionName, revisionName string, wasmModule models.RunningWasmModule, functions map[string]models.Function) {
    fmt.Println("🟡[scale] Add a new wasm module to an existing revision of an existing function:", revisionName)
    functions[functionName].Revisions[revisionName].WasmModules[wasmModule.Pid] = wasmModule
}

// AddRevisionWithWasmModuleToFunction updates the functions map
func AddRevisionWithWasmModuleToFunction(functionName, revisionName, wasmModuleUrl string, wasmModule models.RunningWasmModule, functions map[string]models.Function) {
    fmt.Println("🟠️ Creation of the revision to an existing function:", revisionName)
    functions[functionName].Revisions[revisionName] = models.Revision{
        Name:            revisionName,
        WasmRegistryUrl: wasmModuleUrl,
        WasmModules: map[int]models.RunningWasmModule{
            wasmModule.Pid: wasmModule,
        },
    }
}

// AddFunctionWithRevisionWithWasmModule updates the functions map
func AddFunctionWithRevisionWithWasmModule(functionName, revisionName, wasmModuleUrl string, wasmModule models.RunningWasmModule, functions map[string]models.Function) {
    fmt.Println("🟣️ First deployment of the function:", functionName)
    fmt.Println("🟣️ Creation of the revision:", revisionName)

    functions[functionName] = models.Function{
        Name: functionName,
        Revisions: map[string]models.Revision{
            revisionName: models.Revision{
                Name:            revisionName,
                WasmRegistryUrl: wasmModuleUrl,
                WasmModules: map[int]models.RunningWasmModule{
                    wasmModule.Pid: wasmModule,
                },
            },
        },
    }
}

// StartFunction : Start a function
func StartFunction(capsulePath string, wasmEnvVariables map[string]string, wasmModuleUrl string, httpPortCounter int) (pid int, processStatus, tmpFileName string) {

    tmpFileName = uuid.New().String() + ".wasm"

    /*
    	fmt.Println("🎃", "StartFunction")
    	fmt.Println("🎃", "capsulePath", capsulePath)
    	fmt.Println("🎃", "wasmModuleUrl", wasmModuleUrl)
    	fmt.Println("🎃", "tmpFileName", tmpFileName)
    */

    cmd := exec.Command(
        capsulePath,
        "-url="+wasmModuleUrl,
        "-mode=http",
        "-httpPort="+strconv.Itoa(httpPortCounter),
        "-wasm=./tmp/"+tmpFileName) //TODO: record this in the list of modules to clean when undeploy

    //TODO: add this tmpFileName to the function list
    cmd.Env = os.Environ()

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    for varName, varValue := range wasmEnvVariables {
        fmt.Println("🟢Env Variables", varName, varValue)
        cmd.Env = append(cmd.Env, varName+`=`+varValue)
    }

    /*
        Start in http mode:
        MESSAGE="💊 Capsule Rocks 🚀" go run main.go \
       -wasm=../wasm_modules/capsule-hello-post/hello.wasm \
       -mode=http \
       -httpPort=7070

    */

    err := cmd.Start()

    if err != nil {
        fmt.Println("😡 when starting the capsule process", err.Error())
        processStatus = "NOT_STARTED"
    } else {
        processStatus = "STARTED"
    }

    fmt.Println("🚀 service started, process Id:", cmd.Process.Pid)

    return cmd.Process.Pid, processStatus, tmpFileName
}

// RegisterFunction : Register a function to the reverse proxy
func RegisterFunction(functionName, revisionName, moduleServerUrl, reverseProxy, backend, processStatus, reverseProxyAdminToken string) string {
    client := resty.New()
    bodyStr := `{"function":"` + functionName + `", "revision":"` + revisionName + `", "url":"` + moduleServerUrl + `"}`
    resp, err := client.
        R().
        EnableTrace().
        SetHeader("Content-Type", "application/json; charset=utf-8").
        SetHeader("CAPSULE_REVERSE_PROXY_ADMIN_TOKEN", reverseProxyAdminToken).
        SetBody(bodyStr).
        Post(reverseProxy + "/" + backend + "/functions/registration")

    if err != nil {
        fmt.Println("😡 when registering the function to the reverse proxy", err.Error())
        //fmt.Println(bodyStr)
        processStatus += "[NOT_REGISTERED]"
    } else {
        fmt.Println("🖐[RegisterFunction]", resp)
    }
    return processStatus
}

func RegisterRevision(functionName, revisionName, moduleServerUrl, reverseProxy, backend, processStatus, reverseProxyAdminToken string) string {
    //TODO: prevent creating "default" revision
    client := resty.New()
    bodyStr := `{"function":"` + functionName + `", "revision":"` + revisionName + `", "url":"` + moduleServerUrl + `"}`
    resp, err := client.
        R().
        EnableTrace().
        SetHeader("Content-Type", "application/json; charset=utf-8").
        SetHeader("CAPSULE_REVERSE_PROXY_ADMIN_TOKEN", reverseProxyAdminToken).
        SetBody(bodyStr).
        Post(reverseProxy + "/" + backend + "/functions/" + functionName + "/revision")

    // http://localhost:8888/memory/functions/hola/revision
    if err != nil {
        fmt.Println("😡 when registering the revision to the reverse proxy", err.Error())
        //fmt.Println(bodyStr)
        processStatus += "[NOT_REGISTERED]"
    } else {
        fmt.Println("🖐[RegisterRevision]", resp)
    }
    return processStatus
}

func RegisterURL(functionName, revisionName, moduleServerUrl, reverseProxy, backend, processStatus, reverseProxyAdminToken string) string {
    client := resty.New()
    bodyStr := `{"url":"` + moduleServerUrl + `"}`
    resp, err := client.
        R().
        EnableTrace().
        SetHeader("Content-Type", "application/json; charset=utf-8").
        SetHeader("CAPSULE_REVERSE_PROXY_ADMIN_TOKEN", reverseProxyAdminToken).
        SetBody(bodyStr).
        Post(reverseProxy + "/" + backend + "/functions/" + functionName + "/" + revisionName + "/url")

    if err != nil {
        fmt.Println("😡 when registering the url to the reverse proxy", err.Error())
        //fmt.Println(bodyStr)
        processStatus += "[NOT_REGISTERED]"
    } else {
        fmt.Println("🖐[RegisterURL]", resp)
    }
    return processStatus
}

func DefineDeployRoute(router *gin.Engine, functions map[string]models.Function, capsulePath string, httpPortCounter int, workerDomain, reverseProxy, backend, reverseProxyAdminToken, workerAdminToken string) {

    /*
        ==============================================================
        Deploy a new function:
          - download (from the registry) and start the wasm module
          - register to the reverse proxy
        🖐🖐🖐 Revisions and Tags are not necessarily the same thing
        ==============================================================
       Curl Query:
          curl -v -X POST \
          http://localhost:9999/functions/deploy \
          -H 'content-type: application/json; charset=utf-8' \
          -d '{"function": "hello", "revision": "default", "downloadUrl": "http://localhost:4999/k33g/hello/0.0.1/hello.wasm"}'
          echo ""

       How to pass the environment variables???
       If I call it 2 times, it scales
    */
    router.POST("functions/deploy", func(c *gin.Context) {
        //TODO: check if there is a better practice to handle authentication token
        if len(workerAdminToken) == 0 || CheckWorkerAdminToken(c, workerAdminToken) {

            // check json payload parameters
            jsonMap := make(map[string]interface{})

            if err := c.Bind(&jsonMap); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code":    "JSON_PARSE_ERROR",
                    "message": err.Error()})
            } else {

                //TODO: check if the values are empty or not
                functionName := jsonMap["function"].(string)
                revisionName := jsonMap["revision"].(string)
                wasmModuleUrl := jsonMap["downloadUrl"].(string) // the downloadUrl to download the module from the registry
                envVariables := jsonMap["envVariables"]

                var wasmEnvVariables = make(map[string]string)
                if envVariables == nil {
                    wasmEnvVariables = map[string]string{}
                } else {
                    for key, value := range envVariables.(map[string]interface{}) {
                        strKey := fmt.Sprintf("%v", key)
                        strValue := fmt.Sprintf("%v", value)
                        wasmEnvVariables[strKey] = strValue
                    }

                    //wasmEnvVariables = envVariables.(map[string]string)
                }
                fmt.Println("🚀[envVariables]:", wasmEnvVariables)

                fmt.Println("⏳ downloading from:", wasmModuleUrl)
                fmt.Println("🚀 starting on http port:", httpPortCounter)

                // Increment the httpPort counter
                httpPortCounter += 1

                //🖐 we need the IP address of the worker (for the registration with the reverse proxy)
                //or domain name
                //because the worker and the module are on the same machine
                //but not necessarily the reverse proxy
                //🤔 how to start a module with https??? (or only the reverse proxy???)

                moduleServerPort := "http" //TODO: handle the case of https
                moduleServerUrl := helpers.GetModuleServerUrl(workerDomain, moduleServerPort, httpPortCounter)
                moduleRemoteUrl := reverseProxy + "/functions/" + functionName + "/" + revisionName

                fmt.Println("🌍 the 💊 Capsule module is running at:", moduleServerUrl)
                fmt.Println("📝 registering to the reverse proxy:", reverseProxy)
                fmt.Println("🎉 you can call the function at:", moduleRemoteUrl)

                fmt.Println("🗄 updating the list of the functions")

                // 🚀 Start a function
                pid, processStatus, tmpFileName := StartFunction(capsulePath, wasmEnvVariables, wasmModuleUrl, httpPortCounter)

                if models.IsFunctionExist(functionName, functions) == true {

                    //TODO: prevent the creation of "default" revision

                    // the function already exists
                    if models.IsRevisionExist(functionName, revisionName, functions) == true {
                        // 📦 Register a new **Wasm Module**
                        processStatus = RegisterURL(functionName, revisionName, moduleServerUrl, reverseProxy, backend, processStatus, reverseProxyAdminToken)

                        // The revision already exists
                        // Then we will add a new running wasm module (scale)
                        wasmModule := models.RunningWasmModule{
                            Pid:          pid,
                            Status:       processStatus,
                            LocalUrl:     moduleServerUrl,
                            RemoteUrl:    moduleRemoteUrl,
                            EnvVariables: wasmEnvVariables,
                            TmpFileName:  tmpFileName,
                        }
                        AddRunningWasmModuleToRevision(functionName, revisionName, wasmModule, functions)

                    } else {
                        // 📝 Register a **Revision** to the reverse proxy
                        processStatus = RegisterRevision(functionName, revisionName, moduleServerUrl, reverseProxy, backend, processStatus, reverseProxyAdminToken)

                        // The revision does not exist
                        // Then Create a new revision for the function
                        // With a running wasm module
                        wasmModule := models.RunningWasmModule{
                            Pid:          pid, // at the end pid counter will be the id of the process
                            Status:       processStatus,
                            LocalUrl:     moduleServerUrl,
                            RemoteUrl:    moduleRemoteUrl,
                            EnvVariables: wasmEnvVariables,
                            TmpFileName:  tmpFileName,
                        }
                        AddRevisionWithWasmModuleToFunction(functionName, revisionName, wasmModuleUrl, wasmModule, functions)

                    }
                } else {

                    // 📝 Register a **Function** to the reverse proxy
                    processStatus = RegisterFunction(functionName, revisionName, moduleServerUrl, reverseProxy, backend, processStatus, reverseProxyAdminToken)

                    // The function does not exist: this is the first deployment of the function
                    // The revision does not exist
                    // Then, create the function and the revision
                    wasmModule := models.RunningWasmModule{
                        Pid:          pid,
                        Status:       processStatus,
                        LocalUrl:     moduleServerUrl,
                        RemoteUrl:    moduleRemoteUrl,
                        EnvVariables: wasmEnvVariables,
                        TmpFileName:  tmpFileName,
                    }
                    AddFunctionWithRevisionWithWasmModule(functionName, revisionName, wasmModuleUrl, wasmModule, functions)
                }

                //delete(functions[functionName].Revisions, revisionName)

                c.JSON(http.StatusAccepted, gin.H{
                    "code":      "FUNCTION_DEPLOYED",
                    "message":   "Function deployed",
                    "function":  functionName,
                    "revision":  revisionName,
                    "localUrl":  moduleServerUrl,
                    "remoteUrl": moduleRemoteUrl})

            }
        } else {
            c.JSON(http.StatusForbidden, gin.H{
                "code":    "KO",
                "from":    "worker",
                "message": "Forbidden"})
        }

    })
}
