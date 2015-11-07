# photoweb 
## 使用原生html/template
    
	go语言最简单模板照片上传网站
## 1.0
	这个是最简化的go语言demo
	使用了以下包中的函数
	
~~~go
"os"
"io/ioutil"
"log"
"net/http"
"html/template"
~~~
### http/template
	双大括号{{}}是区分模板代码和HTML的分隔符
	括号里边可以是要显示输 出的数据,或者是控制语句,比如if判断式或者range循环体等
	.|formatter表示对当前这个元素的值以 formatter 方式进行格式化输出
	.|urlquery}即表示对当前元素的值进行转换以适合作为URL一部 分
	.|html 表示对当前元素的值进行适合用于HTML 显示的字符转化,比如">"会被转义 成"&gt;"
## 2.0
 	在这个版本中，增加了viewHandler
### 
http.ServeFile()方法将该路径下的文件从磁盘中 读取并作为服务端的返回信息输出给客户端
