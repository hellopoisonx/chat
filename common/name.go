package common

import (
	"math/rand"
)

func GetName() string {
	// 定义一些常见的名字部分
	firstNames := []string{"John", "Emma", "Michael", "Olivia", "William", "Sophia", "James", "Ava", "Ethan", "Isabella"}
	lastNames := []string{"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor"}

	// 随机选择名字和姓氏部分
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	// 返回拼接后的名字
	return firstName + " " + lastName
}
