package mode

import "testing"

func Test4_ConcurrentDownload(t *testing.T) {
	urls := []string{
		"http://example.com/file1",
		"http://example.com/file2",
		"http://example.com/file3",
		"http://example.com/file4",
		"http://example.com/file5",
	}

	// 3个worker，每秒2个请求，突发3个
	concurrentDownload(urls, 2, 3, 3)
}
