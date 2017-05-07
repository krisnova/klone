package local

import "os/user"

func Home() string {
	usr, err := user.Current()
	if err != nil {
		Printf("unable to find user: %v", err)
		return ""
	}
	return usr.HomeDir
}
