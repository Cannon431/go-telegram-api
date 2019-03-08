package telegram_api

func intInSlice(item int, slice []int) bool {
	found := false

	for _, v := range slice {
		if v == item {
			found = true
			break
		}
	}

	return found
}
