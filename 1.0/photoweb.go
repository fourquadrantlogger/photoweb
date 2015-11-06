package main

import (
	"io"
	"os"
	"config"
)

import(
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

const(
	Upload_dir="./uploads"
	View_dir="./views"
)
func listHandler(w http.ResponseWriter,r *http.Request){
	fileInfoArr,err:=ioutil.ReadDir("./uploads")
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		
		return
	}
	
	locals:=make(map[string]interface{})
	images:=[]string{}

	for _,fileInfo:=range fileInfoArr{
		images=append(images,fileInfo.Name())
	}

	locals["images"]=images
	t:=template.Must(template.ParseFiles(View_dir+"/"+"list.html"))
	log.Println(View_dir+"/"+"list.html")
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		
		return
	}
	if err:=t.Execute(w,locals);err!=nil{
		log.Println("Execute"+err.Error())
	}
}
func uploadHandler(w http.ResponseWriter,r *http.Request){
	switch r.Method{
		case "GET":
			bytes,err:=ioutil.ReadFile("./views/upload.html")
			if err!=nil{
				log.Fatal("ioutil.ReadFile():",err.Error())
			}else{
				w.Write(bytes)
			}
		case "POST":
			//寻找表单中名为image的文件域
			file,h,err:=r.FormFile("image")
			if err!=nil{
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
			
			filename:=h.Filename

			defer file.Close()
			t,err:=os.Create(Upload_dir+"/"+filename)
			log.Println("Create"+Upload_dir+"/"+filename)
			if err!=nil{
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}

			defer t.Close()
			_,err=io.Copy(t,file)
			if err!=nil{
				http.Error(w,err.Error(),http.StatusInternalServerError)
				return
			}
			log.Println("upload:"+filename)
			http.Redirect(w,r,"/views?id="+filename,http.StatusFound)
	}
}

func main(){
	config.LoadConfig()
	http.HandleFunc("/upload",uploadHandler)
	http.HandleFunc("/views",listHandler)
	err:=http.ListenAndServe(":8090",nil)
	log.Println("http.ListenAndServe(:8090)")
	if err!=nil{
		log.Fatal("ListenAndServe:",err.Error())
	}
}