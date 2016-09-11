package queryapi

import (
	"encoding/json"
	"errors"
	"fmt"
	json_sample "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	UPLOAD_LOC_CREATE_URL = "http://yuntuapi.amap.com/datamanage/data/create"
	UPLOAD_LOC_UPDATE_URL = "http://yuntuapi.amap.com/datamanage/data/update"

	TYPE_COORDSYS_GPS      = "gps"
	TYPE_COORDSYS_MAPBAR   = "mapbar"
	TYPE_COORDSYS_BAIDU    = "baidu"
	TYPE_COORDSYS_AUTONAVI = "autonavi"
)

type GDClient struct {
	GAODE_KEY string
}

type GDGeo struct {
	*BaseResp
	Regeocode struct {
		F_address     string `json:"formatted_address"`
		AddrComponent struct {
			Province     string      `json:"province"`
			City         interface{} `json:"city"`
			Citycode     string      `json:"citycode"`
			District     string      `json:"district"`
			Adcode       string      `json:"adcode"`
			Township     string      `json:"township"`
			Neighborhood struct {
				Name interface{} `json:"name"`
				Type interface{} `json:"type"`
			} `json:"neighborhood"`
			Building struct {
				Name interface{} `json:"name"`
				Type interface{} `json:"type"`
			} `json:"building"`
			StreetNumber struct {
				Street    string `json:"street"`
				Number    string `json:"number"`
				Location  string `json:"location"`
				Direction string `json:"direction"`
				Distance  string `json:"distance"`
			} `json:"streetNumber"`
			BusinessAreas interface{} `json:"businessAreas"`
		} `json:"addressComponent"`
		Aois []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"aois"`
	} `json:"regeocode"`
}

type BusinessAreas struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Id       string `json:"id"`
}

type UserLocation struct {
	Key     string      `json:"key"`
	TableID string      `json:"tableid"`
	Data    interface{} `json:"data"`
}

/*type UserLocationData struct {
    UserID     string `json:"userid"`
    Name       string `json:"_name"`
    Location   string `json:"_location"`
    Address    string `json:"_address"`
    DetailAddr string `json:"detailaddress"`
    Aoi        string `json:"aoi"`
    PicUrl     string `json:"pic_url"`
}*/

type BaseResp struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
}

type LocationConvert struct {
	*BaseResp
	Locations string `json:"locations"`
}

func (client *GDClient) GetAddrByLocation2(loc, ct string) (gdgeo *GDGeo, err error) {
	loc = client.ConvertLocation(loc, ct)
	return client.GetAddrByLocation(loc)
}

func (client *GDClient) GetAddrByLocation(loc string) (gdgeo *GDGeo, err error) {
	gaode_url := fmt.Sprintf("http://restapi.amap.com/v3/geocode/regeo?output=json&location=%s&radius=0&extensions=all&key=%s", loc, client.GAODE_KEY)
	//fmt.Print(gaode_url + "\n")

	http_resp, err1 := http.Get(gaode_url)
	if err1 != nil {
		return nil, err1
	}
	defer http_resp.Body.Close()

	body, err2 := ioutil.ReadAll(http_resp.Body)
	if err2 != nil {
		return nil, err2
	}

	//fmt.Print(string(body) + "\n")

	gdgeo = &GDGeo{}
	err = json.Unmarshal(body, gdgeo)
	if err != nil {
		return gdgeo, err
	}

	return gdgeo, nil
}

func (client *GDClient) ConvertLocation(loc, ct string) string {
	if ct == "" || ct == TYPE_COORDSYS_AUTONAVI {
		return loc
	}

	gaode_url := fmt.Sprintf("http://restapi.amap.com/v3/assistant/coordinate/convert?key=%s&locations=%s&coordsys=%s", client.GAODE_KEY, loc, ct)
	http_resp, err1 := http.Get(gaode_url)
	if err1 != nil {
		//fmt.Println(err1)
		return loc
	}
	defer http_resp.Body.Close()

	body, err2 := ioutil.ReadAll(http_resp.Body)
	if err2 != nil {
		//fmt.Println(err2)
		return loc
	}

	lcr := &LocationConvert{}
	err := json.Unmarshal(body, lcr)
	if err != nil {
		//fmt.Println(err)
		return loc
	}

	return lcr.Locations
}

func (gdgeo *GDGeo) GetLocationDesc() string {
	var desc string

	aoi := gdgeo.GetFirstLocationAoi()
	if aoi != "" {
		desc = fmt.Sprintf("%s\n", aoi)
	}

	desc = desc + fmt.Sprintf(
		"地址:%s\n街道:%s",
		gdgeo.Regeocode.F_address,
		gdgeo.GetDetailLocation())

	return desc
}

func (gdgeo *GDGeo) GetDetailLocation() string {
	return fmt.Sprintf(
		"%s%s向%s%s米",
		gdgeo.Regeocode.AddrComponent.StreetNumber.Street,
		gdgeo.Regeocode.AddrComponent.StreetNumber.Number,
		gdgeo.Regeocode.AddrComponent.StreetNumber.Direction,
		gdgeo.Regeocode.AddrComponent.StreetNumber.Distance)
}

func (gdgeo *GDGeo) GetFirstLocationAoi() string {
	if len(gdgeo.Regeocode.Aois) == 1 {
		return gdgeo.Regeocode.Aois[0].Name
	}
	return ""
}

func (user_location *UserLocation) UploadLocation(amapid string) (string, error) {
	data_bytes, _ := json.Marshal(user_location.Data)
	//fmt.Println(string(data_bytes))

	v := url.Values{}
	v.Set("key", user_location.Key)
	v.Set("tableid", user_location.TableID)
	v.Set("data", string(data_bytes))

	var post_url string
	if amapid == "" {
		post_url = UPLOAD_LOC_CREATE_URL
	} else {
		post_url = UPLOAD_LOC_UPDATE_URL
	}

	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	resp, err := http.Post(post_url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resp_body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return "", err1
	}

	ret_json, j_err := json_sample.NewJson(resp_body)
	if j_err != nil {
		return "", j_err
	}
	json_map, map_err := ret_json.Map()
	if map_err != nil {
		return "", map_err
	}
	//fmt.Println(json_map)

	// 新增位置
	if amapid == "" {
		amap_id, ok := json_map["_id"]
		if !ok {
			return "", errors.New("not found _id!")
		}

		_id_str, parse_ok := amap_id.(string)
		if !parse_ok {
			return _id_str, errors.New("_id is not string!")
		}

		return _id_str, nil
	} else { //更新位置
		info, ok := json_map["info"]
		if !ok {
			return amapid, errors.New("info not found!")
		}

		info_str, parse_ok := info.(string)
		if parse_ok && info_str == "OK" {
			return amapid, nil
		} else {
			return amapid, errors.New("update localtion error!")
		}
	}

}
