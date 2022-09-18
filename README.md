# oof
Errors with stack traces for golang

## Usage

```
import "github.com/hashibuto/oof"
```

Wrap an error and include a stack trace

```
err := method(xyz)
if err != nil {
    return oof.Trace(err)
}
```

Wrap an error and include a stack trace, with an annotation
```
err := method(xyz)
if err != nil {
    return oof.Tracef("This is my error: %w", err)
}
```

Get the original error (safe to call on any error interface) - it will return the supplied error if `err` is not an `OofError` type.
```
err := method(xyz)
origErr := oof.GetOrigError(err)
```

## Example output
```
2022/09/18 05:27:17 Hello, this is my error: Special error occurred: SpecialError
goroutine 7 [running]:
runtime/debug.Stack()
        /usr/local/go/src/runtime/debug/stack.go:24 +0x65
github.com/hashibuto/oof.Tracef({0x52e6b9, 0x1b}, {0xc00009fee8, 0x1, 0x1})
        /home/me/oof/oof.go:67 +0xb8
github.com/hashibuto/oof.ApplicationLevelCaller1()
        /home/me/oof/oof_test.go:36 +0x99
github.com/hashibuto/oof.TestTracef(0xc000007ba0)
        /home/me/oof/oof_test.go:57 +0x25
testing.tRunner(0xc000007ba0, 0x534b20)
        /usr/local/go/src/testing/testing.go:1446 +0x10b
created by testing.(*T).Run
        /usr/local/go/src/testing/testing.go:1493 +0x35f
```