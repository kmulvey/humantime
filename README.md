# humantime
[![humantime](https://github.com/kmulvey/humantime/actions/workflows/release_build.yml/badge.svg)](https://github.com/kmulvey/humantime/actions/workflows/release_build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kmulvey/humantime)](https://goreportcard.com/report/github.com/kmulvey/humantime) [![Go Reference](https://pkg.go.dev/badge/github.com/kmulvey/humantime.svg)](https://pkg.go.dev/github.com/kmulvey/humantime)

Convert English strings related to time to Go time.Time. This package also implements the [flags.Value](https://pkg.go.dev/flag#Value) and [flag.Getter](https://pkg.go.dev/flag#Getter) interfaces for use in a cli context.

## Nomenclature
- A "date phrase" is text that represents a date and optionally time, examples:
  - May 8, 2009 5:57:51 PM
  - 3/15/2022
  - yesterday
  - yesterday at 4pm
  - tomorrow at 13:34:32
- A complete list of supported date formats can be found [here](https://github.com/araddon/dateparse#extended-example)
  - In addition to this list, "yesterday", "today" and "tomorrow" are also supported
  
## Supported formats
  - since [date phrase]
  - until or til [date phrase]
  - before [date phrase]
  - after [date phrase]
  - [date phrase] ago
  - from [date phrase] to [date phrase]
 
## Example phrases
  - from May 8, 2009 5:57:51 PM to Sep 12, 2021 3:21:22 PM
  - 3 days ago
  - after yesterday at 4pm

## Usage
  [CLI flag example](https://github.com/kmulvey/humantime/blob/main/cmd/main.go)
  ```
    var st, err = NewString2Time(now.Location())
    result, err := st.After("after 3/15/2022")
   
    fmt.Println(result)    // From: 15 Mar 22 00:00 MDT, To: 19 Jul 22 15:02 MDT
  ```
