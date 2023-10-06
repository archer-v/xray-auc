package main

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Operation uint

const (
	OpAdd    Operation = 1
	OpRemove Operation = 2
)

func (o Operation) String() string {
	if o == OpAdd {
		return "add"
	}
	if o == OpRemove {
		return "remove"
	}
	return "undef"
}

func ApplyToXray(host string, op Operation, records []UserRecord) (done int, err error) {

	if len(records) == 0 {
		fmt.Printf("Nothing to do")
		return
	}

	var xrayApi XrayAPI

	if err := xrayApi.Init(host); err != nil {
		fmt.Printf("Couldn't init xray connection: %v", err)
	}

	defer func() {
		xrayApi.Close()
	}()

	for _, u := range records {
		if err := u.Validate(); err != nil {
			continue
		}
		var e error
		if op == OpAdd {
			e = xrayApi.AddUser(u.Proto, u.Tag, map[string]interface{}{
				"email":    u.Email,
				"id":       u.Password,
				"flow":     u.Flow,
				"password": u.Password,
				"cipher":   u.Cipher,
			})
		} else if op == OpRemove {
			e = xrayApi.RemoveUser(u.Tag, u.Email)
		} else {
			fmt.Println("Unknown operation")
			continue
		}

		if e == nil {
			fmt.Printf("UserRecord [%v] has been %ved by api\n", u.Email, op.String())
			done++
		} else {
			if s, ok := status.FromError(e); ok {
				if s.Code() == codes.Unavailable {
					err = e
					return
				} else {
					if op == OpRemove && s.Code() == codes.Unknown {
						fmt.Printf("UserRecord [%v] doesn't exists or already removed\n", u.Email, op.String())
						done++
					}
				}
			}
			fmt.Printf("Error in %ving the user [%v] by api: %v\n", op.String(), u.Email, e)
		}
	}
	return
}
