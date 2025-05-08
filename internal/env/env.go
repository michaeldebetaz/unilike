package env

type Env struct {
	BASE_PATH string
}

const BASE_PATH = "https://applicationspub.unil.ch/interpub/noauth/php/Ud/"

func GetEnv(key string) Env {
	return Env{BASE_PATH: BASE_PATH}
}
