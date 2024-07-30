package main

import (
	// "bytes"
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
	Simfile   string
	Unixpath  string
}

type SBatchjob struct {
	Id      string
	Fileurl string
}

type Users struct {
	User     string
	Fileurl  string
	Unixpath string
}
type Job struct {
	Jobid string
}

type User struct {
	Id     string `json:"id"`
	Passwd string `json:"passwd"`
	Name   string `json:"name"`
	Rdate  string `json:"rdate"`
	Home   string `json:"home"`
	Status string `json:"status"`
	Server string `json:"server"`
	Role   string `json:"role"`
	Userid string `json:"userid"`
}

type Data struct {
	VALUE   string
	PROJECT string
}

func dbConn(server int) (db *sql.DB) {

	dnDriver := "mysql"
	var constr string

	switch server {
	case 1:
		constr = "root:manager@tcp(10.15.10.148:3306)/leadsship_db"
	case 2:
		constr = "root:leadship!@tcp(10.15.10.129:3306)/manhour"
	case 3:
		constr = "root:leadship!@tcp(10.15.20.108:3306)/leadsship_db"
	default:

	}

	db, err := sql.Open(dnDriver, constr)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {

	//http.Handle("/", http.FileServer(http.Dir("templates")))
	// handler := func(w http.ResponseWriter, r *http.Request) {
	// 	// Set CORS headers to allow requests from any origin
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 	// Respond to preflight requests
	// 	if r.Method == "OPTIONS" {
	// 		return
	// 	}

	// 	// Your handler logic goes here
	// 	w.Write([]byte("This is Slurm job Rest web service !"))
	// }
	//http.HandleFunc("/", handler)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/leadLogin", ajaxleadLoginHandler)

	http.HandleFunc("/ajaxrunFileserver", ajaxrunFileserverHandler)
	http.HandleFunc("/ajaxStopfileserver", ajaxStopfileserverHandler)

	http.HandleFunc("/ajaxStatus", ajaxStatusHandler)
	http.HandleFunc("/ajaxRunStar", ajaxRunStarHandler)
	http.HandleFunc("/ajaxGetJobStatus", GetJobStatusHandler)
	http.HandleFunc("/ajaxCancelJob", GetCancelJobHandler)
	http.HandleFunc("/ajaxFileUpload", ajaxFileUploadHandler)
	http.HandleFunc("/Company", Company)

	http.ListenAndServe(":9022", nil)
}

func ajaxleadLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	hash := md5.New()
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)

	db := dbConn(1)
	selDB, err := db.Query("SELECT id, passwd, name, home, server, role, userid, status from hpc_users where id = ?", u.Id)

	if err != nil {
		panic(err.Error())
	}

	user := User{}
	res := []User{}

	for selDB.Next() {
		var id, passwd, name, home, server, role, userid, status string

		err = selDB.Scan(&id, &passwd, &name, &home, &server, &role, &userid, &status)
		if err != nil {
			panic(err.Error())
		}

		user.Id = id
		user.Passwd = passwd
		user.Name = name
		user.Home = home
		user.Server = server
		user.Role = role
		user.Userid = userid
		user.Status = status

		res = append(res, user)
	}

	var result string

	hash.Write([]byte(user.Id))
	hashedBytes := hash.Sum(nil)
	decid := hex.EncodeToString(hashedBytes)

	if decid == user.Passwd {
		result = `{"result": "success"}`
	} else {
		result = `{"result": "false"}`
	}

	if result == `{"result": "success"}` {
		rst, err := db.Exec("update hpc_users set status='Y', rdate = now() where id =? ", u.Id)
		if err != nil {
			panic(err.Error())
		}
		nRow, err := rst.RowsAffected()
		fmt.Println("update count: ", nRow)
	}
	defer db.Close()

	if result == `{"result": "success"}` {
		empJson, _ := json.Marshal(res)
		w.Write(empJson)
	} else {
		empJson, _ := json.Marshal(result)
		w.Write(empJson)
	}
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
	//workDir := "/home/leadship/HPJobs/10-SimFileTest/00-createSimFile/01-slurm"
	workDir := job.Unixpath
	//cmd := exec.Command("cd ", "/home/leadship/HPJobs/10-SimFileTest/00-createSimFile/01-slurm")

	err = os.Chdir(workDir)
	if err != nil {
		fmt.Println("작업 디렉토리 변경 중 오류 발생: ", err)
	}

	//var jobn = "--job-name=" + job.Jobname
	comd := "/home/leadship/999platform/Works/" + job.Runscript
	jobTime := "02:00:00"
	var jobArgument = ""
	var jobArgument1 = ""
	if job.Simfile == "" {
		jobArgument = "starccm11"
		jobArgument1 = ""
	} else {
		jobArgument = "starccm17"
		jobArgument1 = job.Simfile
	}

	cmd := exec.Command("sbatch",
		"--job-name", job.Jobname,
		"--nodes", job.Nodes,
		"--time", jobTime,
		comd, jobArgument, jobArgument1)
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

func ajaxrunFileserverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var user Users
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	cmd := exec.Command("cd ", user.User)
	//workDir0 := "/home/leadship/999platform/LeadsSHIP"
	workDir := user.Unixpath
	//workDir := "/home/leadship/HPJobs/10-SimFileTest/00-createSimFile/01-slurm"

	err = os.Chdir(workDir)
	if err != nil {
		fmt.Println("작업 디렉토리 변경 중 오류 발생: ", err)
	}
	fileserver := "~/999platform/Works/wfs-ls " + workDir
	//fileserver := "~/999platform/Works/wfs-ls " + workDir0
	cmd = exec.Command("bash", "-c", fileserver)
	//cmd = exec.Command("bash", "-c", "./wfs-ls", workDir0)
	//cmd = exec.Command("bash", "-c", "./"+user.User+"/wfs-ls")
	//cmd := exec.Command("sudo", "-u", "mgsim", "bash", "-c", "./"+"mgsim"+"/wfs-ls")
	//cmd := exec.Command("sudo", "-u", "leadship", "bash", "-c", "./wfs-ls")

	// 명령을 백그라운드에서 실행합니다.
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	cmd = exec.Command("pwd")
	//output, err := cmd.Output()
	// 명령을 백그라운드에서 실행합니다.
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	//fmt.Println(output)
	//fmt.Println("Command is running in the background.")
	//fmt.Printf("Process started with PID: %d\n", cmd.Process.Pid)

	//result := "{result: " + strconv.Itoa(cmd.Process.Pid) + "}"

	rstJson, _ := json.Marshal(cmd.Process.Pid)

	w.Write(rstJson)
}

func ajaxStopfileserverHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var job Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//cmd := exec.Command("sudo", "-u", user.User, "bash", "-c", "./"+user.User+"/wfs-ls")
	cmd := exec.Command("bash", "-c", "kill -9 "+job.Jobid)
	//cmd := exec.Command("sudo", "-u", "leadship", "bash", "-c", "./wfs-ls")

	fmt.Println("Stop Process", job.Jobid)

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
	//cmd := exec.Command("cat", "squeue.log")
	cmd := exec.Command("squeue")

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
	cmd2 := exec.Command("cat", "squeue.log")
	//cmd2 := exec.Command("squeue")

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

	var jobs []Job
	err := json.NewDecoder(r.Body).Decode(&jobs)

	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}

	// Log the received items

	for _, item := range jobs {
		fmt.Printf("Received item: ID=%d, Name=%s\n", item.Jobid)
		jobid := item.Jobid

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

}

func ajaxFileUploadHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	err := r.ParseMultipartForm(10 << 20) // Max upload size ~10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file0")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a new file in the current directory with the same name as the uploaded file
	dst, err := os.Create(handler.Filename)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error copying file: ", err)
		return
	}

	// Display success message
	fmt.Fprintf(w, "File %s uploaded successfully!", handler.Filename)
}

func Company(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	req, _ := io.ReadAll(r.Body)
	uid := string(req)

	db := dbConn(3)
	selDB, err := db.Query("SELECT company FROM leads_users where id = ?", uid)

	check(err)

	emp := Data{}
	res := []Data{}
	var data string

	for selDB.Next() {

		err = selDB.Scan(&data)

		check(err)

		emp.VALUE = data

		//res = append(res, emp)
	}

	query2 := "SELECT PROJECT_NO FROM json_file WHERE status = 1"
	rows2, err := db.Query(query2)
	check(err)
	defer rows2.Close()

	// 두 번째 쿼리 결과 처리
	for rows2.Next() {
		var project_no string
		err := rows2.Scan(&project_no)
		check(err)
		emp.PROJECT = project_no
		res = append(res, emp)
	}

	defer db.Close()

	empJson, err := json.Marshal(res)
	check(err)
	w.Write(empJson)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
