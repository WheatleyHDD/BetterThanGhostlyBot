package utils

import (
    
)

func RemoveItemString(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}