package domain

type EncryptedImg struct {
	Id       string
	UserId   string
	EncryUrl string
}

type EncryptedIv struct {
	Id      string
	EncryId string
	Ivs     string
}

type DataEncryRes struct {
	Id       string
	UserId   string
	EncryUrl string
	Ivs      string
}
