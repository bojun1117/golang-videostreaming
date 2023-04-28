# golang-videostreaming


http://golangvideostreaming-env.eba-zvpmmryb.ap-northeast-1.elasticbeanstalk.com/


## 基本介紹


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
* cache : redis/Amazon ElastiCache
```
"github.com/go-redis/redis/v8"
```
* deploy : Elastic Beanstalk
