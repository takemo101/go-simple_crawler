# シンプルなクローラ（練習）  
指定したページのaタグのhrefを取得して、エクセルファイルとして保存する簡単なクローラを作りました。コマンド引数でURIと保存するファイル名を設定します。  
  
### 必要なパッケージ  
外部パッケージが必要となりますので、下記コマンドで取得してください。    
```  
go get github.com/PuerkitoBio/goquery  
go get github.com/tealeg/xlsx  
```  
  
### 実行方法  
対象ページのURIと保存するエクセルファイル名を指定します。  
```  
go run main.go -uri=http://xxx.com -fn=links
```  
