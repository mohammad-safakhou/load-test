package models

import (
	"fmt"
	"strconv"
)

type AccountID string

func (a AccountID) String() string {
	return string(a)
}

func (a AccountID) ID() int {
	aa, _ := strconv.Atoi(string(a[2:]))
	return aa
}

func (a AccountID) Other(num int) AccountID {
	temp := num - a.ID() + 1
	if temp == a.ID() {
		temp += 1
	}
	return NewAccountID(temp)
}

func NewAccountID(id int) AccountID {
	return AccountID(fmt.Sprintf("00%d", id))
}
