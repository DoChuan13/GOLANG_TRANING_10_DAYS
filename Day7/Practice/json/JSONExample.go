package json

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address string  `json:"address"`
	Phone   []Phone `json:"phone"`
}
type Phone struct {
	No string `json:"phone_no"`
}

func TestJsonExample() {
	var phones = []Phone{{No: "0393.xxx.xxx"}, {No: "0393.xxx.xxy"}}
	user := User{Name: "Chuan", Age: 30, Address: "Nam Dinh", Phone: phones}
	val, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Convert Failed")
		return
	}
	fmt.Println("Original Value ===>", user)
	fmt.Println("JSON Content ===>", string(val))

	valIndent, _ := json.MarshalIndent(user, "", "\t")
	fmt.Println("JSON Content ===>")
	fmt.Println(string(valIndent))

	var unUser User
	err = json.Unmarshal([]byte(val), &unUser)
	if err != nil {
		fmt.Println("Un Marshal Failed")
	}
	fmt.Println("Un Marshal JSON Content ===>", unUser)
}
