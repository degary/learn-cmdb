package tenant

import (
	"github.com/astaxie/beego/logs"
	"github.com/degary/learn-cmdb/cloud"
	"github.com/degary/learn-cmdb/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TenantCloud struct {
	addr       string
	region     string
	accessKey  string
	secretKey  string
	credential *common.Credential
	profile    *profile.ClientProfile
}

func (c *TenantCloud) Type() string {
	return "tenant"
}

func (c *TenantCloud) Name() string {
	return "腾讯云"
}

func (c *TenantCloud) Init(addr, region, accessKey, secretKey string) {
	c.addr = addr
	c.region = region
	c.accessKey = accessKey
	c.secretKey = secretKey
	c.credential = common.NewCredential(c.accessKey, c.secretKey)
	c.profile = profile.NewClientProfile()
	c.profile.HttpProfile.Endpoint = c.addr
}
func (c *TenantCloud) TestConnect() error {
	client, err := cvm.NewClient(c.credential, c.region, c.profile)
	if err != nil {
		return err
	}
	request := cvm.NewDescribeRegionsRequest()
	_, err = client.DescribeRegions(request)

	return err
}

func (c *TenantCloud) GetInstance() []*cloud.Instance {
	var (
		offset int64 = 0
		limit  int64 = 100
		total  int64 = 1
		rt     []*cloud.Instance
	)

	for offset < total {
		var instances []*cloud.Instance
		total, instances = c.getInstanceByOffsetLimit(offset, limit)
		//判断是否为第一次
		if offset == 0 {
			rt = make([]*cloud.Instance, 0, total)
		}
		rt = append(rt, instances...)
		offset += limit
	}
	return rt
}

func (c *TenantCloud) getInstanceByOffsetLimit(offset, limit int64) (int64, []*cloud.Instance) {
	client, err := cvm.NewClient(c.credential, c.region, c.profile)
	if err != nil {
		logs.Error(err)
		return 0, nil
	}
	request := cvm.NewDescribeInstancesRequest()
	request.Offset = common.Int64Ptr(offset)
	request.Limit = common.Int64Ptr(limit)

	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logs.Error("An API error has returned: %s", err)
	}
	if err != nil {
		logs.Error(err)
	}
	total := *response.Response.TotalCount
	instances := response.Response.InstanceSet
	rt := make([]*cloud.Instance, len(instances))

	for index, instance := range instances {
		publicAddrs := make([]string, len(instance.PublicIpAddresses))
		privateAddrs := make([]string, len(instance.PrivateIpAddresses))
		for index, paddr := range instance.PublicIpAddresses {
			publicAddrs[index] = *paddr
		}
		for index, addr := range instance.PrivateIpAddresses {
			privateAddrs[index] = *addr
		}
		rt[index] = &cloud.Instance{
			UUID:         *instance.InstanceId,
			Name:         *instance.InstanceName,
			Os:           *instance.OsName,
			CPU:          int(*instance.CPU),
			Memory:       *instance.Memory,
			PublicAddrs:  publicAddrs,
			PrivateAddrs: privateAddrs,
			Status:       c.transformStatus(*instance.InstanceState),
			CreatedTime:  utils.PtrToString(instance.CreatedTime),
			ExpiredTime:  utils.PtrToString(instance.ExpiredTime),
		}
	}

	return total, rt
}

func (c *TenantCloud) StartInstance(uuid string) error {
	client, err := cvm.NewClient(c.credential, c.region, c.profile)
	if err != nil {
		return err
	}
	request := cvm.NewStartInstancesRequest()
	request.InstanceIds = []*string{&uuid}
	_, err = client.StartInstances(request)
	return err
}

func (c *TenantCloud) StopInstance(uuid string) error {
	client, err := cvm.NewClient(c.credential, c.region, c.profile)
	if err != nil {
		return err
	}
	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{uuid})
	_, err = client.StopInstances(request)
	return err
}

func (c *TenantCloud) RebootInstance(uuid string) error {
	client, err := cvm.NewClient(c.credential, c.region, c.profile)
	if err != nil {
		return err
	}
	request := cvm.NewRebootInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{uuid})
	_, err = client.RebootInstances(request)
	return err

}

func (c *TenantCloud) transformStatus(status string) string {
	smap := map[string]string{
		"PENDING":       cloud.StatusPending,
		"LAUNCH_FAILED": cloud.StatusLaunchFailed,
		"RUNNING":       cloud.StatusRunning,
		"STOPPED":       cloud.StatusStopped,
		"STARTING":      cloud.StatusStarting,
		"STOPPING":      cloud.StatusStopping,
		"REBOOTING":     cloud.StatusRebooting,
		"TERMINATING":   cloud.StatusTerminating,
		"SHUTDOWN":      cloud.StatusShutdown,
	}
	if rt, ok := smap[status]; ok {
		return rt
	}
	return cloud.StatusUnknown
}

func init() {
	cloud.DefaultManager.Register(new(TenantCloud))
}
