package walker

func arrayToMap[C comparable](arr []C) map[C]bool {
    m := make(map[C]bool)
    for _, v := range arr {
        m[v] = true
    }
    return m
}

