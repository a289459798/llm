package utils

func GetSuanLi(t string) int {
	suanli := 1
	switch t {
	case "image/create":
	case "image/edit":
		suanli = 3
		break
	case "creation/article":
	case "code/generate":
	case "code/exam":
		suanli = 2
		break
	}
	return suanli
}