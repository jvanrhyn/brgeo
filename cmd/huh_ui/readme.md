# huh_ui

`huh_ui` is a Go application designed to fetch and display geographical information based on IP addresses. It interacts with a local geolocation API and presents the results in a user-friendly format.

## Functionality

The core function of the application, `getGeoInfo`, takes an IP address as an input and makes a GET request to a local API endpoint (`http://localhost:3000/api/lookup/{ipAddress}`). The API is expected to return geographical information related to the provided IP address.

The geographical information is then unmarshalled from the JSON response into a `GeoInfo` struct, which is returned by the function along with any potential error.

The application then formats and displays this information in a user-friendly way using the `lipgloss` package to create a styled output.

## Usage

To use the `huh_ui` application, you need to have a running instance of the geolocation API on your local machine.

Once the API is running, you can call the `getGeoInfo` function with an IP address to get the geographical information related to that IP address.

```go
geoInfo, err := getGeoInfo("8.8.8.8")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("GeoInfo: %+v\n", geoInfo)
```

## Error Handling

The `getGeoInfo` function is designed to handle potential errors that might occur during the HTTP request or the JSON unmarshalling process. If an error occurs, the function will return `nil` for the `GeoInfo` struct and the error.

## Dependencies

The `huh_ui` application depends on the following Go packages:

- `net/http`: for making HTTP requests
- `fmt`: for formatting strings
- `io`: for reading the HTTP response body
- `encoding/json`: for unmarshalling the JSON response into a struct
- `lipgloss`: for creating styled terminal layouts

## Future Improvements

Future versions of the `huh_ui` application might include more robust error handling, support for different geolocation APIs, and additional features based on user feedback.