package domain

import (
	"time"
)

type UserProof = int

const (
	UserProofNone  UserProof = 0
	UserProofEmail UserProof = 1 << iota
	UserProofPhone
	UserProofBoris
)

type User struct {
	Id             int64
	Email          string
	Surname        string
	GivenNames     string
	Phone164       uint64
	BornAt         time.Time
	HasProof       UserProof
	PasswordDigest []byte
	CreatedAt      time.Time
}

func (u *User) Reset() {
	u.Id = 0
	u.Email = ""
	u.Surname = ""
	u.GivenNames = ""
	u.Phone164 = 0
	u.BornAt = time.Time{}
	u.HasProof = 0
	u.PasswordDigest = u.PasswordDigest[:0]
	u.CreatedAt = time.Time{}
}
