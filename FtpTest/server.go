package main

import (
	// "bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Sqjson struct {
	Jobid    string
	Partname string
	Name     string
	User     string
	St       string
	Time     string
	Nodes    string
	Nodelist string
}
type SBatchreq struct {
	Jobname   string
	Runscript string
	Nodes     string
}

type SBatchjob struct {
	Id string
}

func main() {

	//http.Handle("/", http.FileServer(http.Dir("templates")))
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Respond to preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		// Your handler logic goes here
		w.Write([]byte("This is Slurm job Rest web service !"))
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/ajaxrunFileserver", ajaxrunFileserverHandler)
	http.HandleFunc("/ajaxStatus", ajaxStatusHandler)
	http.HandleFunc("/ajaxRunStar", ajaxRunStarHandler)
	http.HandleFunc("/ajaxGetJobStatus", GetJobStatusHandler)
	http.HandleFunc("/ajaxCancelJob", GetCancelJobHandler)

	http.ListenAndServe(":8000", nil)
}

func ajaxrunFileserverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	cmd := exec.Command("bash", "-c", "wfs-ls")

	// 명령을 백그라운드에서 실행합니다.
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	//fmt.Println("Command is running in the background.")
	//fmt.Printf("Process started with PID: %d\n", cmd.Process.Pid)

	result := "{result: " + strconv.Itoa(cmd.Process.Pid) + "}"

	rstJson, _ := json.Marshal(result)

	w.Write(rstJson)
}

func ajaxStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	//sdata := Sqjson{}
	sarray := []Sqjson{}

	// 외부 프로세스 실행
	cmd := exec.Command("cat", "squeue.log")
	//cmd := exec.Command("squeue")

	// 표준 출력 읽기
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	// 줄 바꿈으로 분할하여 배열로 변환
	lines := strings.Split(string(output), "\n")

	// JSON 배열 생성
	var jsonArray []string
	for _, line := range lines {
		if strings.Contains(line, "JOBID") {
			//if strings.Contains(line, "JOBID") || strings.Contains(line, "Resources") || strings.Contains(line, "Priority") {
			jsonArray = append(jsonArray, line)
		} else {
			words := strings.Fields(line)
			sdata := Sqjson{}
			var i = 0
			for _, word := range words {
				switch i {
				case 0:
					sdata.Jobid = word
				case 1:
					sdata.Partname = word
				case 2:
					sdata.Name = word
				case 3:
					sdata.User = word
				case 4:
					sdata.St = word
				case 5:
					sdata.Time = word
				case 6:
					sdata.Nodes = word
				case 7:
					sdata.Nodelist = word
				}
				i++
			}
			sarray = append(sarray, sdata)
		}
	}

	// JSON으로 인코딩하여 출력
	empJson, err := json.Marshal(sarray)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	//fmt.Println(string(empJson))

	w.Write(empJson)
}

func ajaxRunStarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	sarray := []Sqjson{}

	var job SBatchreq
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	var jobn = "--job-name=" + job.Jobname
	var comd = job.Runscript
	cmd := exec.Command("sbatch", jobn, comd)
	//cmd := exec.Command("sbatch", "--job-name=my_job", "my_script.sh")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("표준 출력 파이프 생성 실패:", err)
		return
	}

	// 명령을 실행합니다.
	if err := cmd.Start(); err != nil {
		fmt.Println("명령 실행 실패:", err)
		return
	}

	// 표준 출력을 읽어 작업 ID를 추출합니다.
	scanner := bufio.NewScanner(stdoutPipe)
	var jobID = ""
	for scanner.Scan() {
		line := scanner.Text()
		// 작업 ID는 "Submitted batch job {job_id}" 형식으로 나옵니다.
		if strings.HasPrefix(line, "Submitted batch job") {
			parts := strings.Fields(line)
			jobID = parts[len(parts)-1]
			fmt.Println("작업 ID:", jobID)
			break
		}
	}

	// 명령이 완료될 때까지 대기합니다.
	if err := cmd.Wait(); err != nil {
		fmt.Println("명령 완료 대기 실패:", err)
	}

	jobid := jobID

	// 외부 프로세스 실행
	//cmd := exec.Command("cat", "squeue.log")
	cmd2 := exec.Command("squeue")

	// 표준 출력 읽기
	output, err := cmd2.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	// 줄 바꿈으로 분할하여 배열로 변환
	lines := strings.Split(string(output), "\n")

	// JSON 배열 생성
	var jsonArray []string
	for _, line := range lines {
		if strings.Contains(line, "JOBID") {
			//if strings.Contains(line, "JOBID") || strings.Contains(line, "Resources") || strings.Contains(line, "Priority") {
			jsonArray = append(jsonArray, line)
		} else {
			words := strings.Fields(line)
			sdata := Sqjson{}
			var i = 0
			for _, word := range words {
				switch i {
				case 0:
					sdata.Jobid = word
				case 1:
					sdata.Partname = word
				case 2:
					sdata.Name = word
				case 3:
					sdata.User = word
				case 4:
					sdata.St = word
				case 5:
					sdata.Time = word
				case 6:
					sdata.Nodes = word
				case 7:
					sdata.Nodelist = word
				}
				i++
			}
			if sdata.Jobid == jobid {
				sarray = append(sarray, sdata)
			}
		}
	}

	// JSON으로 인코딩하여 출력
	empJson, err := json.Marshal(sarray)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	//fmt.Println(string(empJson))

	w.Write(empJson)
}

func GetJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	sarray := []Sqjson{}

	var job SBatchjob
	err := json.NewDecoder(r.Body).Decode(&job)

	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	jobid := job.Id

	// 외부 프로세스 실행
	//cmd := exec.Command("cat", "squeue.log")
	cmd2 := exec.Command("squeue")

	// 표준 출력 읽기
	output, err := cmd2.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	// 줄 바꿈으로 분할하여 배열로 변환
	lines := strings.Split(string(output), "\n")

	// JSON 배열 생성
	var jsonArray []string
	for _, line := range lines {
		if strings.Contains(line, "JOBID") {
			jsonArray = append(jsonArray, line)
		} else {
			words := strings.Fields(line)
			sdata := Sqjson{}
			var i = 0
			for _, word := range words {
				switch i {
				case 0:
					sdata.Jobid = word
				case 1:
					sdata.Partname = word
				case 2:
					sdata.Name = word
				case 3:
					sdata.User = word
				case 4:
					sdata.St = word
				case 5:
					sdata.Time = word
				case 6:
					sdata.Nodes = word
				case 7:
					sdata.Nodelist = word
				}
				i++
			}
			if sdata.Jobid == jobid {
				sarray = append(sarray, sdata)
			}
		}
	}

	// JSON으로 인코딩하여 출력
	empJson, err := json.Marshal(sarray)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	w.Write(empJson)
}

func GetCancelJobHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var job SBatchjob
	err := json.NewDecoder(r.Body).Decode(&job)

	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	jobid := job.Id

	// 외부 프로세스 실행
	cmd2 := exec.Command("scancel", jobid)

	// 표준 출력 읽기
	output, err := cmd2.Output()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	w.Write(output)
}
