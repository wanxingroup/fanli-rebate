module dev-gitlab.wanxingrowth.com/fanli/rebate

go 1.13

require (
	dev-gitlab.wanxingrowth.com/fanli/card v0.0.0-20200922154313-ad42e368223a
	dev-gitlab.wanxingrowth.com/fanli/merchant v0.0.0-20200927035010-cff6f54a4fe5
	dev-gitlab.wanxingrowth.com/fanli/user v0.0.2
	dev-gitlab.wanxingrowth.com/wanxin-go-micro/base v0.2.26
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.2.2
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/shomali11/util v0.0.0-20190608141102-c39c2521a2ab // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/ini.v1 v1.60.2 // indirect
)

replace dev-gitlab.wanxingrowth.com/fanli/merchant => github.com/wanxingroup/fanli-merchant v0.0.0

replace dev-gitlab.wanxingrowth.com/fanli/user => github.com/wanxingroup/fanli-user v0.0.2

replace dev-gitlab.wanxingrowth.com/wanxin-go-micro/base => github.com/wanxingroup/base v0.2.27

replace dev-gitlab.wanxingrowth.com/fanli/card => github.com/wanxingroup/fanli-card v0.0.0
