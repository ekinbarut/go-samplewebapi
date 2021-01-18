package main

func exists(object Object, Objects []Object) bool {
	flag := false
	for _, obj := range Objects {
		if object.Id == obj.Id {
			flag = true
		}
	}

	return flag
}
