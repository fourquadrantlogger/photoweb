package main

import (
	"path"
	"io"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

var	root=""
const(
	Assert_dir="/assert"
	Upload_dir="/uploads"
	View_dir="/views"
)
var MyTemplates=make(map[string] *template.Template)
//给templates加载所有views文件夹下的模版文件
func Loadtmpl(){	
	fileInfoArr,err:=ioutil.ReadDir(root+View_dir)
	check(err)
	
	var temlateName,temlatePath string
	for _,fileInfo:=range fileInfoArr{
		temlateName=fileInfo.Name()
		log.Println(temlateName)
		//检查后缀名
		if ext:=path.Ext(temlateName);ext!=".html"{
			continue
		}
		temlatePath=root+View_dir+"/"+temlateName
		log.Println("Loadtmpl()",temlatePath)
		t:=template.Must(template.ParseFiles(temlatePath))
		MyTemplates[temlateName]=t;
		log.Println(len(MyTemplates),MyTemplates[temlateName])
	}
}
func check(err error){
	if err!=nil{
		panic(err)
	}
}
func renderHtml(w http.ResponseWriter,tmpl string,locals map[string] interface{}){	
	err:=MyTemplates[tmpl+".html"].Execute(w,locals)
	check(err)
}

func checkUploadDir(){
	Uploadinfo, err := os.Stat(root+Upload_dir)
	if err != nil {
	    os.Mkdir(root+Upload_dir,os.ModePerm)
		log.Println("os.Mkdir("+root+Upload_dir)
	    return
	}
	if Uploadinfo.IsDir() {
	    // it's a file
	} else {
	    os.Mkdir(root+Upload_dir,os.ModePerm)
	}

	_, err= os.Stat(root+View_dir)
	if err != nil {
		log.Fatal("views file not found! require template file")
	    // no such file or dir
	    return
	}

}

func uploadHandler(w http.ResponseWriter,r *http.Request){
	
	switch r.Method{
		case "GET":
			renderHtml(w,"upload",nil)
		case "POST":
			//寻找表单中名为image的文件域
			f,h,err:=r.FormFile("image")
			check(err)			
			filename:=h.Filename

			defer f.Close()
			//
			t,err:=os.Create(root+"/"+Upload_dir+"/"+filename)
			log.Println("Create"+Upload_dir+"/"+filename)
			check(err)

			defer t.Close()
			_,err=io.Copy(t,f)
			check(err)
			log.Println("upload:"+filename)
			http.Redirect(w,r,"/list",http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter,r *http.Request){
	imageid:=r.FormValue("id")
	imagepath:=root+Upload_dir+"/"+imageid
	if _,err:=os.Stat(imagepath);err!=nil{
		http.NotFound(w,r)
	}
	w.Header().Set("Content-Type","image")
	http.ServeFile(w,r,imagepath)
}
func listHandler(w http.ResponseWriter,r *http.Request){
	fileInfoArr,err:=ioutil.ReadDir(root+"/"+"uploads")
	check(err)
	
	locals:=make(map[string]interface{})
	images:=[]string{}

	for _,fileInfo:=range fileInfoArr{
		images=append(images,fileInfo.Name())
	}

	locals["images"]=images
	
	renderHtml(w,"list",locals);
}
func staticDirHandler(mux *http.ServeMux,prefix string,staticDir string,flags int){
	mux.HandleFunc(prefix,
		func(w http.ResponseWriter,r *http.Request){
			log.Println(r.URL.Path)
			file:=root+r.URL.Path
			log.Println(file)
			http.ServeFile(w,r,file)
		})
}
func main(){
    var mux=http.NewServeMux()
	
	root,_=os.Getwd()
	checkUploadDir()
	Loadtmpl()
	staticDirHandler(mux,"/assets/",root+"/assets",0)
	mux.HandleFunc("/upload",uploadHandler)
	mux.HandleFunc("/list",listHandler)
	mux.HandleFunc("/views",viewHandler)
	err:=http.ListenAndServe(":8090",mux)
	log.Println("http.ListenAndServe(:8090)")
	if err!=nil{
		log.Fatal("ListenAndServe:",err.Error())
	}
}