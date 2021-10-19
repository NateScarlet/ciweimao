# ciweimao

[![godev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)](https://pkg.go.dev/github.com/NateScarlet/ciweimao/pkg)

刺猬猫 Go 客户端，使用安卓客户端 API。

个人用于通过 RSS 发现新书，所以没有动力实现章节相关内容，欢迎 PR。

API 解密实现基于 [Cirno-go](https://github.com/zsakvo/Cirno-go) 的解包成果。

- [x] 账号密码登录
- [x] 获取排行榜
- [x] 搜索书籍
- [x] 书籍详情
- [ ] 章节内容

## 使用方法

```go
package main

import (
    "context"

    "github.com/NateScarlet/ciweimao/pkg/client"
    "github.com/NateScarlet/ciweimao/pkg/book"
)

// 默认客户端用环境变量 `CIWEIMAO_ACCOUNT` `CIWEIMAO_LOGIN_TOKEN` `CIWEIMAO_DEVICE_TOKEN` 登录。
// 这些值获取方法见下方 `通过环境变量配置登录凭据`
// 并且 User-Agent 使用 `CIWEIMAO_USER_AGENT` 或库内置的默认值。
client.Default

var c = new(client.Client)
c.ApplyDefaultConfig()

// 账号密码登录（每次进程启动都需要重新登录，或许会触发风控）
_, err = c.Login(username, password)

// 登录凭据登录（推荐）
// 数值可通过账号密码登录或者使用 ciweimao-login 命令获得（见下方）
c.Account = account
c.LoginToken = loginToken
c.DeviceToken = deviceToken

// 所有查询从 context 获取客户端设置, 如未设置将使用默认客户端。
var ctx = context.Background()
ctx = client.With(ctx, c)

// 搜索书籍
result, err := book.Search(ctx, book.SearchOptionQuery("搜索关键词"))
result.JSON // JSON 响应数据.
result.Books() // []book.Book
book.Search(ctx, book.SearchOptionQuery("搜索关键词"), book.SearchOptionPageIndex(1)) // 获取第二页

// 书籍详情
var book = &book.Book{ID: "12345678"}
_, err = i.Fetch(ctx) // 获取书籍详情

// 书籍排行
rank, err := book.Rank(ctx, book.RTClick, book.RPWeek) // 周点击榜
require.NoError(t, err)
rank.Books() // []book.Book
book.Rank(ctx, book.RTClick, book.RPWeek, book.RankOptionPageIndex(2)) // 第二页
book.Rank(ctx, book.RTClick, book.RPWeek, book.RankOptionPageCategory(book.C免费同人)) // 分类筛选
```

### 通过环境变量配置登录凭据

获取登录凭据，设置环境变量然后运行 ciweimao-login

```bash
CIWEIMAO_USERNAME={username} CIWEIMAO_PASSWORD={password} go run github.com/NateScarlet/ciweimao/cmd/ciweimao-login
```

成功登录后会出现如下输出，设置这些环境变量即可使默认客户端自动登录

```shell
CIWEIMAO_ACCOUNT=书客*********
CIWEIMAO_LOGIN_TOKEN=******************************
CIWEIMAO_DEVICE_TOKEN=ciweimao_***************
```

凭据会在一定时间后过期，如果要长期运行可额外设置账号密码使客户端自动更新令牌

```shell
CIWEIMAO_USERNAME=*******
CIWEIMAO_PASSWORD=*******
```

如果需要存储更新的令牌可以通过自定义 Client.TokenRefresher 实现

## 类似项目

都是用来找书的

- [pixiv](https://github.com/NateScarlet/pixiv)
- [qidian](https://github.com/NateScarlet/qidian)
