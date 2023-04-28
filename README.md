# golang-videostreaming


http://golangvideostreaming-env.eba-zvpmmryb.ap-northeast-1.elasticbeanstalk.com/


### 使用套件


* database : Postgresql/AWS RDS
```
"database/sql"
"github.com/lib/pq"
```
* authentication : Session
```
"github.com/gorilla/sessions"
```
* stream : net/http package
```
"net/http"
```
* videos : AWS S3
```
"github.com/aws/aws-sdk-go-v2/aws"
"github.com/aws/aws-sdk-go-v2/config"
"github.com/aws/aws-sdk-go-v2/service/s3"
```
* cache : Redis/Amazon ElastiCache
```
"github.com/go-redis/redis/v8"
```
* deploy : Elastic Beanstalk


### 基本功能

!( [圖片網址](https://user-images.githubusercontent.com/98681/91365119-402cdc00-e7b5-11ea-9a2c-e1a03aed21c3.png) "收藏")
