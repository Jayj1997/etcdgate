/*
 * @Author       : jayj
 * @Date         : 2021-06-26 23:03:32
 * @Description  :
 */
package res

type ErrorCode int

const (
	_ int = iota + 1000
	// OK 成功
	OK
	ParamsInvalid
	InternalServerError
	RegisterFailed
	TokenInvalid
	TokenExpired
	UserNotExists
	UrlNotFound
	UnauthorizedError
	NotRoot
	InvalidAccountOrPassword
	DuplicatedName
)

var Msg = map[int]string{
	OK:                       "ok",
	ParamsInvalid:            "invalid params",
	InternalServerError:      "server internal error",
	RegisterFailed:           "register failed",
	TokenInvalid:             "token invalid",
	TokenExpired:             "token expired",
	UserNotExists:            "user not exists",
	UrlNotFound:              "url not found",
	UnauthorizedError:        "unauthorized error",
	NotRoot:                  "not root account or auth is not enabled",
	InvalidAccountOrPassword: "invalid account or password",
	DuplicatedName:           "duplicated name",
}

func GetMsg(code int) string {
	return Msg[code]
}
