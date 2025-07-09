package utils

var PanicOnError = panicOnError

func panicOnError(err error) {
  if err != nil {
    panic(err)
  }
}
