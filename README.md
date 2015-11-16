# Config

Library for reading and writing hierarchical configuration data.

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/twoleds-golang/config)

## Release

This library is not ready for production usage. It works, but 
contains bugs.

## Example

```plain
BoolValue true
FloatValue 3.14
IntValue 123
StringValue "Test String 123"

Section One {
    BoolValue true
    FloatValue 2.745
    IntValue 124
    StringValue "Test String - First Section"
}

Nested One {
    Nested Two {
        Nested Three {
            BoolValue false
            FloatValue -3.14
            IntValue -123456
            StringValue "Test String - Nested Section"
        }
    }
}

Section Two {
    BoolValue false
    FloatValue -2.745
    IntValue 123456
    StringValue "Test String 123"
}

Section Three {
    BoolValue true
    FloatValue 3.14
    IntValue -150
    StringValue "Test String - Last Section"
}
```

## Usage

```go

package main

import "github.com/twoleds-golang/config"
import "fmt"

func main() {
	
	// Parse configuration file

	c, err := config.ParseFromFile("/path/to/config/file")
	if err != nil {
		panic(err.Error())
	}

	// Read configuration value

	fmt.Printf(
		"Config value for Section:Two/FloatValue is %f\n",
		c.FloatOrDefault("Section:Two/FloatValue", 0.0),
	)

}

```
