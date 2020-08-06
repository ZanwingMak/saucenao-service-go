package main

import (
	"io"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"encoding/json"
	"net/url"
	"io/ioutil"
	"path"
	"strings"
	"strconv"
)

type Data struct{
	Result string `json:"result"`
}

type Response struct {
	Data Data `json:"data"`
	Code int `json:"code"`
	Success bool `json:"success"`
}

type ErrorResponse struct {
	Msg string `json:"msg"`
	Code int `json:"code"`
	Success bool `json:"success"`
}

type Configuration struct {
	UploadDir string
	ImagesPathPrefix string
}

type SaucenaoParams struct {
	api_key string
	image_url string
	output_type string
	testmode string
	numres string
	db string
	minsim string
}



func saucenaoSearch(w http.ResponseWriter, req *http.Request, data SaucenaoParams) {
	// req.Header.Add("Cookie","token=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println("data:", data);
	params := url.Values{}
	Url, err := url.Parse("https://saucenao.com/search.php")
	params.Set("api_key",data.api_key)
	params.Set("url",data.image_url)
	params.Set("output_type",data.output_type)
	params.Set("testmode",data.testmode)
	params.Set("numres",data.numres)
	params.Set("db",data.db)
	params.Set("minsim",data.minsim)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)

	resp, err := http.Get(urlPath)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http get error.")
		// fmt.Fprintf(w, `{msg: "获取失败",code: 0,success: false}`)

		err_response := ErrorResponse{Code: 2,Success:false,Msg:"获取失败"}			
		err_response_json, _ := json.Marshal(err_response)
		io.WriteString(w, string(err_response_json))

	} else {

		src := string(body)

		// 模拟延时
		// time.Sleep(time.Second * 2)

		// 写法1
		var data Data
		data.Result = src
		// 写法2
		// data := Data{Result: src}

		var res Response
		res.Code = 1
		res.Success = true
		res.Data = data

		res_json, _ := json.Marshal(res)

		// fmt.Fprintf(w, string(bytes))
		// fmt.Fprintf(w, `{data:%+v, success:true, code:1}`, src)
		io.WriteString(w, string(res_json))
	}
}


func saucenaoServer(w http.ResponseWriter, req *http.Request) {
	headerType := req.Header.Get("Content-Type")
	fmt.Println("header-type : ", headerType)

	w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域

	w.Header().Set("Access-Control-Allow-Headers", "*") // 允许的header的类型
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	// w.Header().Add("Access-Control-Allow-Headers", "locale")
	// w.Header().Add("Access-Control-Allow-Headers", "token")

	w.Header().Set("content-type", "application/json")	// 返回json数据

	switch req.Method {
		case "GET" :
			query := req.URL.Query()

			// params := url.Values{}

			// api_key := query.Get("api_key")
			// image_url := query.Get("url")
			// output_type := query.Get("output_type")
			// testmode := query.Get("testmode")
			// numres := query.Get("numres")
			// db := query.Get("db")
			// minsim := query.Get("minsim")

			var params SaucenaoParams

			params.api_key = query.Get("api_key")
			params.image_url = query.Get("url")
			params.output_type = query.Get("output_type")
			params.testmode = query.Get("testmode")
			params.numres = query.Get("numres")
			params.db = query.Get("db")
			params.minsim = query.Get("minsim")

			fmt.Println("params:", params)

			saucenaoSearch(w, req, params)

		case "POST" :
			// json

			// // 根据 body 创建一个 json 解析器实例
			// decoder := json.NewDecoder(req.Body)
			// // 存放参数
			// var params map[string]string
			// // 解析参数 存入 map
			// decoder.Decode(&params)

			// api_key := params["api_key"]
			// output_type := params["output_type"]
			// testmode := params["testmode"]
			// numres := params["numres"]
			// db := params["db"]
			// minsim := params["minsim"]
			// file := params["file"]

			// fmt.Printf("POST json: api_key=%s , output_type=%s , testmode=%s , numres=%s , db=%s , minsim=%s , file=%s\n", api_key, output_type, testmode, numres, db, minsim,file)
			// fmt.Fprintf(w, `{"code":0}`)

			// FormData
			reader, err := req.MultipartReader()
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				err_response := ErrorResponse{Code: 2,Success:false,Msg:err.Error()}			
				err_response_json, _ := json.Marshal(err_response)
				io.WriteString(w, string(err_response_json))
				return
			}

			file, _ := os.Open("./config/saucenao.json")
			defer file.Close()
			decoder := json.NewDecoder(file)
			configuration := Configuration{}
			file_err := decoder.Decode(&configuration)
			curDir,_ := os.Getwd() // 获得当前路径
			fmt.Println("current_path:", curDir)
			// execpath, err := os.Executable() // 获得程序路径
			// fmt.Println("exe_path:", execpath)
			fmt.Println("uploadDir:", configuration.UploadDir)
			if file_err != nil {
				fmt.Println("error:", err)
				err_response := ErrorResponse{Code: 2,Success:false,Msg:"系统异常"}			
				err_response_json, _ := json.Marshal(err_response)
				io.WriteString(w, string(err_response_json))
				return
			}

			var params SaucenaoParams

			// pl := []string{
			// 	"api_key",
			// 	"image_url",
			// 	"output_type",
			// 	"testmode",
			// 	"numres",
			// 	"db",
			// 	"minsim",
			// }
			
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
		
				// fmt.Printf("FileName=[%s], FormName=[%s]\n", part.FileName(), part.FormName())
				if part.FileName() == "" {  // this is FormData
					// fmt.Println(part);
					curData, _ := ioutil.ReadAll(part)
					fmt.Printf("FormName=[%s] FormData=[%s]\n", part.FormName(), curData)
					// formname := string(part.FormName())

					// for i, value := range pl {
					// 	fmt.Println(i, value)
						// if value == "api_key" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.api_key = string(curData)
						// } else if value == "image_url" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.image_url = string(curData)
						// } else if value == "output_type" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.output_type = string(curData)
						// } else if value == "testmode" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.testmode = string(curData)
						// } else if value == "numres" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.numres = string(curData)
						// } else if value == "db" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.db = string(curData)
						// } else if value == "minsim" {
						// 	fmt.Printf("FormName=[%s]", part.FormName())
						// 	params.minsim = string(curData)
						// }
					// }

						if part.FormName() == "api_key" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.api_key = string(curData)
						} else if part.FormName() == "image_url" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.image_url = string(curData)
						} else if part.FormName() == "output_type" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.output_type = string(curData)
						} else if part.FormName() == "testmode" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.testmode = string(curData)
						} else if part.FormName() == "numres" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.numres = string(curData)
						} else if part.FormName() == "db" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.db = string(curData)
						} else if part.FormName() == "minsim" {
							fmt.Printf("FormName=[%s]", part.FormName())
							params.minsim = string(curData)
						}

				} else {    // This is FileData
					fmt.Printf("FormName=[%s]", part.FormName())
					fmt.Print(part)
					fullFilename := part.FileName()
					fmt.Println(fullFilename)
					filenameWithSuffix := path.Base(fullFilename)
					// fmt.Println(filenameWithSuffix)
					fileSuffix := path.Ext(filenameWithSuffix)
					// fmt.Println(fileSuffix)
					filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
					// fmt.Println(filenameOnly)

					// string转成int：
					// 	int, err := strconv.Atoi(string)
					// string转成int64：
					// 	int64, err := strconv.ParseInt(string, 10, 64)
					// int转成string：
					// 	string := strconv.Itoa(int)
					// int64转成string：
					// 	string := strconv.FormatInt(int64,10)

					newFilename := filenameOnly + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + fileSuffix

					createFullpath := curDir + configuration.UploadDir + newFilename

					saucenaoFullpath := configuration.ImagesPathPrefix + newFilename

					fmt.Println("newFilename:", newFilename)
					fmt.Println("createFullpath", createFullpath)
					fmt.Println("saucenaoFullpath", saucenaoFullpath)

					params.image_url = saucenaoFullpath

					dst, _ := os.Create(createFullpath)
					defer dst.Close()
					io.Copy(dst, part)
				}
			}

			fmt.Println("params:", params)

			saucenaoSearch(w, req, params)

	}

}


func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/saucenao", saucenaoServer)

	timeOut := time.Second * 30

	srv := &http.Server{
		Addr:           ":23232",
		Handler:        mux,
		ReadTimeout:    timeOut,
		WriteTimeout:   timeOut,
		IdleTimeout:    timeOut * 2,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		fmt.Println("开始监听接口")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf(" listen and serve http server fail:\n %v ", err)
			fmt.Println(" listen and serve http server fail:\n %v ", err)
		}
		fmt.Println("监听结束")
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)
	<-exit
	ctx, cacel := context.WithTimeout(context.Background(), timeOut)
	defer cacel()
	err := srv.Shutdown(ctx)
	log.Println("shutting down now. ", err)
	os.Exit(0)
}