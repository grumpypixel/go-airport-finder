# go-airport-finder

This is a Golang library which retrieves all sorts of information about airports around the world. The underlying data is based on data collected and provided by [OurAirports.com](https://ourairports.com).

An initial version of this library was built to power the "Airport Finder" for the [GoPilot](https://github.com/grumpypixel/msfs2020-gopilot).

## The API

This should be pretty straightforward and hopefully easy-to-use.

The following snippets are based on the [provided example code](https://github.com/grumpypixel/go-airport-finder/blob/main/example/main.go).

```golang
// Create an AirportFinder instance
finder := alphafoxtrot.NewAirportFinder()

// Create LoadOptions with preset filepaths
dataDir := "./data"
options := alphafoxtrot.PresetLoadOptions(dataDir)

// Specify a filter what airport types should be loaded, in this case: all airports, please.
filter := alphafoxtrot.AirportTypeAll

// Load the airport data into memory
if err := finder.Load(options, filter); len(err) > 0 {
	log.Println("errors:", err)
}
```

```golang
// Alternatively, you can specify the LoadOptions "manually".
// So in this case, only airports, frequencies and runways will be loaded.
// Regions, countries, navaids will be omitted.
options := alphafoxtrot.LoadOptions{
	AirportsFilename:    "./data/airports.csv",            // *required*
	FrequenciesFilename: "./data/airport-frequencies.csv", // optional
	RunwaysFilename      "./data/runways.csv",             // optional
}
```

Please note that the code above assumes that the necessary .CSV files are located in a subdirectory named "data".

You can download the latest version of the .CSV files directly from [OurAirports](https://ourairports.com/data).


```golang
// A hacky way to download all needed files is by calling DownloadDatabase()
dataDir := "./data"
alphafoxtrot.DownloadDatabase(dataDir) // assuming that the given directory exists...
```

So much for the initialization part.

```golang
// Find an airport by its ICAO code
if airport := finder.FindAirportByICAOCode("KLAX"); airport != nil {
	fmt.Println(*airport)
}
```

```golang
// Find an airport by its IATA code
if airport := finder.FindAirportByIATACode("DUS"); airport != nil {
	fmt.Println(*airport)
}
```

```golang
// Find the nearest active airport within a given radius
latitude := 33.942501
longitude := -118.407997
radiusInMeters := alphafoxtrot.NauticalMilesToMeters(25)
airportTypeFilter := alphafoxtrot.AirportTypeActive
if airport := finder.FindNearestAirport(latitude, longitude, radiusInMeters, airportTypeFilter); airport != nil {
	fmt.Println(*airport)
}
```

```golang
// Find the active airports with a runway (so no heliports and no seaplane based airports) within a given radius
latitude := 33.942501
longitude := -118.407997
radiusInMeters := alphafoxtrot.MilesToMeters(61)
maxResults := 10
airportTypeFilter := alphafoxtrot.AirportTypeActive|alphafoxtrot.AirportTypeRunways
airports := finder.FindNearestAirports(latitude, longitude, radiusInMeters, maxResults, airportTypeFilter)
for i, airport := range airports {
	fmt.Println(i, *airport)
}
```

```golang
// Find all large airports in a specific region
regionISOCode := "US-CA"
latitude := 33.942501
longitude := -118.407997
radiusInMeters := alphafoxtrot.NauticalMilesToMeters(100)
maxResults := 10
airportTypeFilter := alphafoxtrot.AirportTypeMedium|alphafoxtrot.AirportTypeSmall
airports := finder.FindNearestAirportsByRegion(regionISOCode, latitude, longitude, radiusInMeters, maxResults, airportTypeFilter)
for i, airport := range airports {
	fmt.Println(i, *airport)
}
```

```golang
// Find all large and medium airports in a country
countryISOCode := "IS"
latitude := 64.1299972534
longitude := -21.9405994415
radiusInMeters := alphafoxtrot.NauticalMilesToMeters(250)
maxResults := 10
airportTypeFilter := alphafoxtrot.AirportTypeLarge|alphafoxtrot.AirportTypeMedium
airports := finder.FindNearestAirportsByCountry(countryISOCode, latitude, longitude, radiusInMeters, maxResults, airportTypeFilter)
for i, airport := range airports {
	fmt.Println(i, *airport)
}
```

```golang
// Find the nearest navaids within a given radius
latitude := 33.942501
longitude := -118.407997
radiusInMeters := alphafoxtrot.KilometersToMeters(50)
maxResults := 10
navaids := finder.FindNearestNavaids(latitude, longitude, radiusInMeters, maxResults)
for i, navaid := range navaids {
	fmt.Println(i, *navaid)
}
```

```golang
// Find the navaids associated with an airport (by ICAO code)
airportICAOCode := "CYYC"
navaids := finder.FindNavaidsByAirportICAOCode(airportICAOCode)
for i, navaid := range navaids {
	fmt.Println(i, *navaid)
}
```

## OurAirports

### Terms of use for the data
From OurAirports:

"DOWNLOAD AND USE AT YOUR OWN RISK! We hereby release all of these files into the [Public Domain](https://en.wikipedia.org/wiki/Public_domain), with no warranty of any kind — By downloading any of these files, you agree that OurAirports.com, Megginson Technologies Ltd., and anyone involved with the web site or company hold no liability for anything that happens when you use the data, including (but not limited to) computer damage, lost revenue, flying into cliffs, or a general feeling of drowsiness that persists more than two days.

Do you agree with the above conditions? If so, then download away!"

## Shoutouts

Big thanks go out to [David Megginson](http://ourairports.com/about.html#credits) and all contributors of OurAirports.com. What an impressive piece of collected awesomeness!