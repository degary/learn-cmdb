package cloud

const (
	StatusPending      = "创建中"
	StatusLaunchFailed = "创建失败"
	StatusRunning      = "运行中"
	StatusStopped      = "已停止"
	StatusStarting     = "开机中"
	StatusStopping     = "停止中"
	StatusRebooting    = "重启中"
	StatusTerminating  = "销毁中"
	StatusUnknown      = "未知状态"
	StatusShutdown     = "停止待销毁"
)

type Instance struct {
	UUID         string
	Name         string
	Os           string
	CPU          int
	Memory       int64
	PublicAddrs  []string
	PrivateAddrs []string
	CreatedTime  string
	ExpiredTime  string
	Status       string
}

type ICloud interface {
	Type() string
	Name() string
	Init(string, string, string, string)
	TestConnect() error
	GetInstance() []*Instance
	StartInstance(string) error
	StopInstance(string) error
	RebootInstance(string) error
}
