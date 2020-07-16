hammer
====
提供**网易云音乐** Golang API 服务接口，提供直接调用方式  

[接口加密算法参考00](https://github.com/darknessomi/musicbox/wiki)  
[接口加密算法参考01](https://github.com/Binaryify/NeteaseCloudMusicApi/blob/master/util/crypto.js)  
[接口地址参考02](https://github.com/Binaryify/NeteaseCloudMusicApi)  

**快速上手**  
----
创建文件夹  

    mkdir hammer
    cd hammer
下载并安装  

    go get github.com/io24m/hammer
创建hammer.go文件  

    package main
    
    import r "github.com/io24m/hammer"
    
    func main() {
    	r.Run()
    }
