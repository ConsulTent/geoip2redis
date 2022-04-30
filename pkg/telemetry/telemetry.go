package telemetry

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	f "github.com/consultent/geoip2redis/pkg/fstat"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

// This is for testing only

func Send(buildId string, geoipType string, geoipFmt string) {
	var iNoHash string = "0"
	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()

	sysType := hostStat.OS + " " + hostStat.Platform + hostStat.PlatformVersion + hostStat.VirtualizationSystem + " " + cpuStat[0].ModelName
	rfctime := time.Now().Format(time.RFC3339)
	iNoHash = GetMD5Hash(f.FStat(os.Args[0]))
	values := map[string]string{"BuildId": buildId, "sysDate": rfctime, "sysType": sysType, "geoipType": geoipType, "geoipFmt": geoipFmt, "os": hostStat.OS, "vmSys": hostStat.VirtualizationSystem, "iNo": iNoHash}

	json_data, _ := json.Marshal(values)

	resp, _ := http.Post("https://telemetry.consultent.workers.dev/geoip2redis", "application/json",
		bytes.NewBuffer(json_data))

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	// fmt.Println(res["docId"])

}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
