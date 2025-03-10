# JSON Structure CLI
[![Go](https://github.com/ptrkrlsrd/json-structure/actions/workflows/go.yml/badge.svg)](https://github.com/ptrkrlsrd/json-structure/actions/workflows/go.yml)
## TODO
- Switch to Cobra


## Example 
Note: The output is a not valid JSON. I might switch to JSON or another format at a later stage.
``` json
[
  {
    "image": string,
    "id": number,
    "winery": string,
    "wine": string,
    "rating": {
      "average": string,
      "reviews": string
    },
    "location": string
  }
]
```