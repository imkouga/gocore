package guid

func GetID() string {
	return newObjectId().Hex()
}
