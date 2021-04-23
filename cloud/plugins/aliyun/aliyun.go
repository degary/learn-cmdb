package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/astaxie/beego/logs"
	"github.com/degary/learn-cmdb/cloud"
)

type Aliyun struct {
	addr      string
	region    string
	accessKey string
	secretKey string
}

func (c *Aliyun) Type() string {
	return "aliyun"
}

func (c *Aliyun) Name() string {
	return "阿里云"
}

func (c *Aliyun) Init(addr, region, accessKey, secretKey string) {
	c.addr = addr
	c.region = region
	c.accessKey = accessKey
	c.secretKey = secretKey
}

func (c *Aliyun) TestConnect() error {
	client, err := ecs.NewClientWithAccessKey(c.region, c.accessKey, c.secretKey)
	if err != nil {
		return err
	}
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"
	_, err = client.DescribeRegions(request)
	return err
}

func (c *Aliyun) GetInstance() []*cloud.Instance {
	var (
		offset int = 1
		limit  int = 100
		total  int = 2
		rt     []*cloud.Instance
	)

	for offset < total {
		var instances []*cloud.Instance
		total, instances = c.getInstanceByOffsetLimit(offset, limit)
		if offset == 1 {
			//初始化rt
			rt = make([]*cloud.Instance, 0, total)
		}
		rt = append(rt, instances...)
		offset += limit
	}
	return rt

}

func (c *Aliyun) getInstanceByOffsetLimit(offset, limit int) (int, []*cloud.Instance) {
	client, err := ecs.NewClientWithAccessKey(c.region, c.accessKey, c.secretKey)
	if err != nil {
		logs.Error("get new aliyun client err:%s", err)
		return 0, nil
	}

	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"

	request.PageNumber = requests.NewInteger(offset)
	request.PageSize = requests.NewInteger(limit)

	response, err := client.DescribeInstances(request)
	if err != nil {
		logs.Error(err)
		return 0, nil
	}
	total := response.TotalCount
	instances := response.Instances.Instance
	//fmt.Println("=============total",total)
	//fmt.Println("=============instances",instances)
	//fmt.Println("=============len instances",len(instances))
	rt := make([]*cloud.Instance, len(instances))
	for index, instance := range instances {
		privateAddrs := make([]string, 0)
		privateAddrs = append(privateAddrs, instance.InnerIpAddress.IpAddress...)
		privateAddrs = append(privateAddrs, instance.VpcAttributes.PrivateIpAddress.IpAddress...)

		rt[index] = &cloud.Instance{
			UUID:         instance.InstanceId,
			Name:         instance.InstanceName,
			Os:           instance.OSName,
			CPU:          instance.Cpu,
			Memory:       int64(instance.Memory),
			CreatedTime:  instance.CreationTime,
			ExpiredTime:  instance.ExpiredTime,
			Status:       c.transformStatus(instance.Status),
			PrivateAddrs: privateAddrs,
			PublicAddrs:  instance.PublicIpAddress.IpAddress,
		}
	}
	return total, rt

}

func (c *Aliyun) StartInstance(uuid string) error {
	client, err := ecs.NewClientWithAccessKey(c.region, c.accessKey, c.secretKey)

	request := ecs.CreateStartInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = uuid

	_, err = client.StartInstance(request)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (c *Aliyun) StopInstance(uuid string) error {
	client, err := ecs.NewClientWithAccessKey(c.region, c.accessKey, c.secretKey)
	if err != nil {
		logs.Error(err)
		return err
	}
	request := ecs.CreateStopInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = uuid
	_, err = client.StopInstance(request)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (c *Aliyun) RebootInstance(uuid string) error {
	client, err := ecs.NewClientWithAccessKey(c.region, c.accessKey, c.secretKey)
	if err != nil {
		logs.Error(err)
		return err
	}
	request := ecs.CreateRebootInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = uuid
	_, err = client.RebootInstance(request)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (c *Aliyun) transformStatus(status string) string {
	smap := map[string]string{
		"Running":  cloud.StatusRunning,
		"Stopped":  cloud.StatusStopped,
		"Starting": cloud.StatusStarting,
		"Stopping": cloud.StatusStopping,
	}
	if rt, ok := smap[status]; ok {
		return rt
	}
	return cloud.StatusUnknown
}

func init() {
	cloud.DefaultManager.Register(new(Aliyun))
}
