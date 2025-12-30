package scanner

import "os"

func countDir(path string) (files int, dirs int) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return 0, 0
    }

    for _, e := range entries {
        if e.IsDir() {
            dirs++
        } else {
            files++
        }
    }
    return
}

