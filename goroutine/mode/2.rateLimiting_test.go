package mode

import (
	"testing"
)

func Test2_RateLimiting(t *testing.T) {
	burstLimit(2, 2) // 每秒最多调用2次，初始允许2次突发请求，每500ms补充1个令牌
}
