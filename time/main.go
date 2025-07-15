package main

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//rawURL := "https://xxx/click?offer_id=123&ua=Mozilla/5.0%20(iPhone;%20CPU%20iPhone%20OS%2017_7_1%20like%20Mac%20OS%20X)%20AppleWebKit/605.1.15%20(KHTML,%20like%20Gecko)%20Version/17.7%20Mobile/15E148%20Safari/604.1"
	//
	//parsedURL, err := url.Parse(rawURL)
	//if err != nil {
	//	fmt.Println("解析失败:", err)
	//	return
	//}
	//
	//queryParams := parsedURL.Query()
	//fmt.Println("ua:", queryParams.Get("ua")) // 输出空值

	fmt.Printf("Parsed URL: %s\n", generateKey("com.amazon.mp3", 251))
	//
	//key := fmt.Sprintf("%s-%s", "aaa.bbb", "8c3a5a58-1515-485a-9d57-fb32468c14fb")
	//shortKey := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	//fmt.Printf("Short Key: %s\n", shortKey)

	//fmt.Println(GetRandomDateStr("2020-01-01T00:00:00Z", "2025-01-01T00:00:00Z"))

	//fmt.Println(time.UnixMilli(1751155800000))
	//beginTime, err := ParseTimeField("29/06/25 2:48")
	//if err != nil {
	//	fmt.Printf("Error parsing time: %v\n", err)
	//	return
	//}
	//fmt.Printf("Parsed Time: %s\n", beginTime)

	fmt.Printf("Generated S3 Prefix: %s\n", generatorS3Prefix(time.Now()))
}

func generateKey(bundleId string, dmpId int) string {
	key := fmt.Sprintf("%s:%d", bundleId, dmpId)
	return fmt.Sprintf("%s:%x", "dmp", md5.Sum([]byte(key)))
}

func GetRandomDateStr(startStr, endStr string) string {
	layout := "2006-01-02T15:04:05Z"
	start, _ := time.Parse(layout, startStr)
	end, _ := time.Parse(layout, endStr)

	duration := end.Sub(start)
	randomTime := start.Add(time.Duration(rand.Int63n(int64(duration))))

	return randomTime.Format(layout)
}

func ParseTimeField(timeStr string) (time.Time, error) {
	layouts := []string{
		"2006/1/2 15:04:05",
		"2006/1/2 15:04",
		"2006/01/02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-01-2006 15:04",
		"02/01/2006 15:04",
		"02/01/06 15:04",    // 添加 DD/MM/YY 格式
		"02/01/06 15:04:05", // 添加 DD/MM/YY 格式带秒
	}

	for _, layout := range layouts {
		parsedTime, err := time.Parse(layout, timeStr)
		if err != nil {
			continue
		}
		return parsedTime, nil
	}
	return time.Time{}, fmt.Errorf("failed to parse time: %s", timeStr)
}

func generatorS3Prefix(t time.Time) string {
	date := t.Format("2006-01-02")
	return fmt.Sprintf("dmp/%d/%s/%s/%s/", 1, date, "", "ios")
}
