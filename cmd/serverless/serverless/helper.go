package serverless

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

// RunCMD run command would run the command
func RunCMD(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("COMMIT_SHA", os.Getenv("COMMIT_SHA"))
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	commandToInvoke := os.Getenv("COMMAND_TO_RUN")
	if commandToInvoke == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("COMMAND_TO_RUN env variable is not set"))
		return
	}
	cmd := exec.Command(fmt.Sprintf("./" + commandToInvoke))
	cmd.Stdout = w
	params := r.URL.Query()
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REQUEST_LINK=%s/%s", r.Host, r.URL.Path),
		fmt.Sprintf("LANGUAGE=%s", params.Get("language")),
		fmt.Sprintf("SINCE=%s", params.Get("since")),
	)
	err := cmd.Run()
	if err != nil {
		os.Stderr.WriteString("Error while processing the request: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
