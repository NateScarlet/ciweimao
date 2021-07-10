// Code generated from [category.go.gotmpl], DO NOT EDIT.

package book

type Category string

const (
    C全部 Category = "0"
    C灵异未知 Category = "1"
    C女频 Category = "11"
    C免费同人 Category = "24"
    C都市青春 Category = "27"
    C游戏竞技 Category = "3"
    C历史军事 Category = "30"
    C仙侠武侠 Category = "5"
    C科幻无限 Category = "6"
    C玄幻奇幻 Category = "8"
)

func (c Category) String() string {
    switch c {
    case C全部:
        return "全部"
    case C灵异未知:
        return "灵异未知"
    case C女频:
        return "女频"
    case C免费同人:
        return "免费同人"
    case C都市青春:
        return "都市青春"
    case C游戏竞技:
        return "游戏竞技"
    case C历史军事:
        return "历史军事"
    case C仙侠武侠:
        return "仙侠武侠"
    case C科幻无限:
        return "科幻无限"
    case C玄幻奇幻:
        return "玄幻奇幻"
    }
    return ""
}


func CategoryByName(name string) Category {
    switch {
    case name == "全部":
        return C全部
    case name == "灵异未知":
        return C灵异未知
    case name == "女频":
        return C女频
    case name == "免费同人":
        return C免费同人
    case name == "都市青春":
        return C都市青春
    case name == "游戏竞技":
        return C游戏竞技
    case name == "历史军事":
        return C历史军事
    case name == "仙侠武侠":
        return C仙侠武侠
    case name == "科幻无限":
        return C科幻无限
    case name == "玄幻奇幻":
        return C玄幻奇幻
    }
    return ""
}
