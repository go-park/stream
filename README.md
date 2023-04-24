# stream

go stream lib based on generics 

## Installation

```shell
# golang 1.18+ required
go get -u github.com/go-park/stream@latest
```

## Quick Start

The Go Stream API is a powerful tool for processing collections and streams of data in a functional way.

1. Create a stream: To get started with the tool, you need to create a stream. You can create a stream from a slice, an array, or by generating elements dynamically. Here's an example of creating a stream from a slice:

    ```go
    slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
    s := stream.From(slice...)
    ```

2. Use stream operations: Once you have a stream, you can use a variety of stream operations to manipulate and process the data. Stream operations are divided into intermediate and terminal operations. Intermediate operations return a new stream, while terminal operations return a result or a side-effect. Here are some examples of stream operations:

    ```go
	slice := []int{1, 2, 3, 4, 5, 61, 7, 8, 9, 10, 11, 19}
	val := stream.From(slice...).
		Filter(func(t int) bool { return t > 2 }).
		Skip(2).Limit(2).
		Map(func(i int) int {
			return i + 1
		}).
		Reduce(func(i1, i2 int) int { return i1 + i2 })
	val.IfNotEmptyOrElse(
		func(v int) { assert.Equal(t, v, 5+1+61+1) },
		func() { t.Error("empty") })
    ```

3. Close the stream: When you're done processing the stream, you need to close it. This releases any resources associated with the stream. You can close a stream using the Close() method or by using a terminal operation that automatically closes the stream, such as ForEach().