package main

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

type Proxy struct {
	Username  string
	Password  string
	IpAddress string
	Port      int
}

func main() {
	var resources = []string{
		"https://ifconfig.co",
		"https://ifconfig.me",
		"https://api.ipify.org",
		"https://www.google.com",
		"https://www.youtube.com",
		"https://www.yandex.ru",
		"https://www.vk.com",
		"https://www.ok.ru",
		"https://www.mail.ru",
		"https://www.rambler.ru",
		"https://www.avito.ru",
		"https://www.yandex.com",
		"https://www.ozon.ru",
		"https://www.drom.ru",
		"https://ipinfo.io/",
		"https://www.ipaddressguide.com/",
		"https://www.maxmind.com/",
		"https://ip-api.com/",
		"https://iplocation.io/",
		"https://www.iptrackeronline.com/",
		"https://www.ip-tracker.org/",
		"https://www.hostip.info/",
		"https://www.site24x7.com/ip-address/",
		"https://ip-address.org/",
		"https://iplocation.co/",
		"https://www.ip-api.com/",
		"https://ipinfo.co/",
		"https://whatismyip.live/",
		"https://ipinfo.info/",
		"https://www.find-ip-address.org/",
		"https://ipchicken.com/",
		"https://iplocate.io/",
		"https://www.ipaddresslocation.org/",
		"https://www.geobytes.com/",
		"https://2ip.io/",
		"https://www.yougetsignal.com/",
		"https://www.ipligence.com/",
		"https://freegeoip.app/",
		"https://www.ip-address.cc/",
		"https://ipfind.com/",
		"https://www.ipligence.com/",
		"https://www.locationof.com/",
		"https://www.ip-address.org/",
		"https://www.geoiptool.com/",
		"https://www.hostip.info/",
		"https://www.ipaddresslocation.org/",
	}

	proxies, err := Read()

	var errors []Proxy

	if err != nil {
		panic(err)

		return
	}

	times := time.Now()

	wg := sync.WaitGroup{}

	for {
		for id, proxy := range proxies {
			go func() {
				wg.Add(1)

				randomIndex := rand.Intn(len(resources))
				resourceUrl := resources[randomIndex]

				time.Sleep(1 * time.Second)

				code := SendQuery(&proxy, resourceUrl)

				if code != 200 {
					errors = append(errors, proxy)

					delete(proxies, id)
				}

				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

				wg.Done()
			}()
		}

		wg.Wait()

		if len(proxies) == 0 {
			break
		}

		if times.After(time.Now().Add(20 * time.Minute)) {
			break
		}
	}

	fmt.Println(errors)
}

func SendQuery(proxy *Proxy, endpoint string) int {
	urlAddress, err := url.Parse("http://" + proxy.IpAddress + ":" + strconv.Itoa(proxy.Port))

	if err != nil {
		return http.StatusInternalServerError
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlAddress),
		},
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxy.Username+":"+proxy.Password))

	client.Transport.(*http.Transport).ProxyConnectHeader = http.Header{"Proxy-Authorization": []string{basicAuth}}

	resp, err := client.Get(endpoint)

	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err.Error())

		return http.StatusInternalServerError
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return resp.StatusCode
}

func Read() (map[int]Proxy, error) {
	workDirectory, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	file, err := os.Open(path.Join(workDirectory, "resources", "data.csv"))

	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	mapped := make(map[int]Proxy)

	for _, record := range records {
		if len(record) < 5 {
			continue
		}

		port, err := strconv.Atoi(record[4])

		if err != nil {
			fmt.Println(err.Error())

			continue
		}

		proxy := Proxy{
			Username:  record[1],
			Password:  record[2],
			IpAddress: record[3],
			Port:      port,
		}

		key, _ := strconv.Atoi(record[0])

		mapped[key] = proxy
	}

	return mapped, nil
}
