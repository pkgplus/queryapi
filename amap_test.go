package queryapi

import (
	"fmt"
	"os"
	"testing"
)

func TestGetLocationDesc(t *testing.T) {
	geoKey := os.Getenv("GAODE_KEY")
	if geoKey == "" {
		t.Skip("the env \"GAODE_KEY\" is null!")
	}

	loc := "116.4931397606772,39.92179726719849"
	client := &GDClient{GAODE_KEY: geoKey}
	gdgeo, err := client.GetAddrByLocation2(loc, TYPE_COORDSYS_GPS)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(gdgeo.GetLocationDesc())
}
