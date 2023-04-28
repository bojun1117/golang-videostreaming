# golang-videostreaming


http://golangvideostreaming-env.eba-zvpmmryb.ap-northeast-1.elasticbeanstalk.com/

----

### 使用套件


* database : Postgresql/AWS RDS
```
"database/sql"
"github.com/lib/pq"
```
* authentication : Session + C
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

----

### 基本功能

* thumbs-up : 收藏/收藏影片
* chat-square-text : 留言
* file-person : 我的影片
* file-earmark-arrow-up : 上傳
* trash : 刪除影片
* heartbreak : 移除收藏
