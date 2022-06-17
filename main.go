package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string
type userObject struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func check(e error, _ int) {
	if e != nil {
		panic(e)
	}
}

func checkOperation(op string) bool {
	switch op {
	case "add", "list", "findById", "remove":
		return true
	}
	return false
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	if validOperation := checkOperation(args["operation"]); !validOperation {
		return fmt.Errorf("Operation %v not allowed!", args["operation"])
	}
	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	f, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	check(err, 0)

	buf, err := ioutil.ReadAll(f)
	check(err, 0)

	var users []userObject
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &users); err != nil {
			panic(err)
		}
	}

	switch arg := args["operation"]; arg {
	case "add":
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}

		var currItem userObject
		if err := json.Unmarshal([]byte(args["item"]), &currItem); err != nil {
			panic(err)
		}
		for _, v := range users {
			if v.Id == currItem.Id {
				n, err := writer.Write([]byte("Item with id " + currItem.Id + " already exists"))
				check(err, n)
			}
		}

		users = append(users, currItem)
		jsonUsers, err := json.Marshal(users)
		check(err, 0)

		if err := ioutil.WriteFile(args["fileName"], jsonUsers, 0644); err != nil {
			panic(err)
		}
	case "list":
		if len(users) > 0 {
			jsonUsers, err := json.Marshal(users)
			check(err, 0)

			n, err := writer.Write(jsonUsers)
			check(err, n)
		}
	case "findById":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		if len(users) > 0 {
			for _, v := range users {
				if v.Id == args["id"] {
					jsonItem, err := json.Marshal(v)
					check(err, 0)

					n, err := writer.Write(jsonItem)
					check(err, n)
				}
			}
		}
	case "remove":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		if len(users) > 0 {
			var isRemove bool
			newUsers := make([]userObject, 0, len(users))

			for i, v := range users {
				if v.Id == args["id"] {
					isRemove = true
					newUsers = append(newUsers, users[i+1:]...)
					break
				} else {
					newUsers = append(newUsers, v)
				}
			}
			if isRemove {
				jsonUsers, err := json.Marshal(newUsers)
				check(err, 0)

				if err := ioutil.WriteFile(args["fileName"], jsonUsers, 0644); err != nil {
					panic(err)
				}
			} else {
				n, err := writer.Write([]byte("Item with id " + args["id"] + " not found"))
				check(err, n)
			}
		} else {
			n, err := writer.Write([]byte("Item with id " + args["id"] + " not found"))
			check(err, n)
		}
	}

	return err
}

func parseArgs() Arguments {
	userId := flag.String("id", "", "user id")
	userObject := flag.String("item", "", "user object")
	operationType := flag.String("operation", "", "type of operation")
	fileName := flag.String("fileName", "", "file name")
	flag.Parse()

	return Arguments{
		"id":        *userId,
		"item":      *userObject,
		"operation": *operationType,
		"fileName":  *fileName,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
