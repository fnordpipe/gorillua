package logger

import (
  "fmt"
  "os"
)

func Debug(msg string) {
  fmt.Fprintf(os.Stdout, "%s\n", msg)
}

func Error(msg string) {
  fmt.Fprintf(os.Stdout, "%s\n", msg)
}

func Stdout(msg string) {
  fmt.Fprintf(os.Stdout, "%s\n", msg)
}
