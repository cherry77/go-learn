# Go 并发模式练习题：Actor 模型实现

以下是几个基于 Actor 模型的 Go 并发编程练习题，帮助你掌握这种通过通道序列化状态访问的模式。

## 基础练习题

### 1. 计数器服务
实现一个并发安全的计数器服务：
- 提供 Add(int) 方法增加计数值
- 提供 Get() int 方法获取当前值
- 所有操作必须通过通道序列化处理

```go
// 你的实现代码
type Counter struct {
    // 需要哪些字段?
}

func NewCounter() *Counter {
    // 如何初始化?
}

func (c *Counter) Add(n int) {
    // 如何通过通道发送增加请求?
}

func (c *Counter) Get() int {
    // 如何通过通道获取当前值?
}
```

## 中级练习题

### 2. 银行账户系统
实现一个银行账户系统，支持：
- 存款(Deposit)
- 取款(Withdraw)
- 查询余额(Balance)
- 转账(Transfer)

要求：
- 每个账户由一个"actor"管理
- 转账操作需要原子性（要么全成功，要么全失败）
- 避免死锁

```go
type Account struct {
    // 你的实现
}

func NewAccount(balance int) *Account {
    // 初始化
}

// 实现各种方法...
```

### 3. 聊天室系统
实现一个简单的聊天室：
- 用户可以加入(Join)、离开(Leave)
- 用户可以发送消息(Send)
- 所有消息广播给所有用户
- 用户列表和消息分发由中央"actor"管理

```go
type ChatRoom struct {
    // 你的实现
}

type User struct {
    Name string
    Recv chan string
    // 其他字段
}

// 实现各种方法...
```

## 高级练习题

### 4. 分布式任务队列
实现一个任务队列系统：
- 生产者可以提交任务
- 多个工作者可以消费任务
- 支持任务优先级
- 支持任务取消
- 所有任务状态由中央"actor"管理

```go
type TaskQueue struct {
    // 你的实现
}

type Task struct {
    ID       int
    Priority int
    Payload  interface{}
    // 其他字段
}

// 实现各种方法...
```

### 5. 股票交易撮合引擎
实现一个简化的股票交易撮合引擎：
- 接收买单和卖单
- 按照价格优先、时间优先原则撮合
- 维护订单簿状态
- 所有订单处理由单一"actor"序列化处理

```go
type OrderBook struct {
    // 你的实现
}

type Order struct {
    ID     int
    Symbol string
    Price  float64
    Amount int
    IsBuy  bool
    // 其他字段
}

// 实现各种方法...
```

## 解答思路提示

1. **基本模式**：
    - 创建一个管理状态的goroutine
    - 定义操作请求和响应的结构体
    - 使用通道传递请求和接收响应

2. **进阶技巧**：
    - 使用 `select` 处理多个通道
    - 使用 `context` 实现优雅关闭
    - 对于复杂操作，可以使用回调通道

3. **错误处理**：
    - 在响应结构中包含错误字段
    - 使用超时机制避免阻塞

4. **性能优化**：
    - 批量处理请求
    - 使用缓冲通道减少阻塞

需要具体哪个练习的详细解答可以告诉我，我可以为你提供完整的实现代码和解释。