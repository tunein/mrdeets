package metadata

import "github.com/aws/aws-sdk-go/service/ec2"

type Ec2Options struct {
	InstanceId    []string `name:"instance-id" from:"*" `
	PublicDNSName []string `name:"public-dns-name" from:"*" `
	//LaunchTime       time.Time `name:"launch-time" from:"*" `
	PublicIPAddress  []string `name:"public-ip-address" from:"*" `
	PrivateIPAddress []string `name:"private-ip-address" from:"*" `
	VpcID            []string `name:"vpc-id" from:"*" `
	AvailabilityZone []string `name:"availability-zone" from:"*" `
	State            []string `name:"state" from:"*" `
	Out              []string `name:"out" from:"*" `
	TagKeys          []string `name:"tag-key" from:"*" `
	TagValues        []string `name:"tag-value" from:"*" `
	Tag              []string `name:"tag" from:"*" `
	SecurityGroups   []string `name:"security-groups" from:"*" `
	Region           string   `name:"aws-region" validate:"nonzero "from:"*" `
}

const TimeStamp = "2006-01-02T15:04:05.000"

type MinimalInstance struct {
	InstanceId       string          `json:"instance-id"`
	PublicIPAddress  string          `json:"public-ip-address"`
	PublicDNSName    string          `json:"public-dns-name"`
	PrivateIPAddress string          `json:"private-ip-address"`
	PrivateDNSName   string          `json:"private-dns-name"`
	LaunchTime       string          `json:"launch-time"`
	SecurityGroups   []SecurityGroup `json:"security-groups"`
	VpcID            string          `json:"vpc-id"`
	Tags             []Tag           `json:"tags"`
}

type SecurityGroup struct {
	GroupName string `json:"group_name"`
	GroupId   string `json:"group_id"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewMinimalEC2Instance(instance *ec2.Instance) *MinimalInstance {

	var secGroups []SecurityGroup
	var tags []Tag
	for _, secGroup := range instance.SecurityGroups {
		s := SecurityGroup{
			GroupName: *secGroup.GroupName,
			GroupId:   *secGroup.GroupId,
		}
		secGroups = append(secGroups, s)
	}

	for _, tag := range instance.Tags {
		t := Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		}
		tags = append(tags, t)
	}
	p := MinimalInstance{
		InstanceId:       *instance.InstanceId,
		PublicIPAddress:  *instance.PublicIpAddress,
		PublicDNSName:    *instance.PublicDnsName,
		PrivateIPAddress: *instance.PrivateIpAddress,
		PrivateDNSName:   *instance.PrivateDnsName,
		LaunchTime:       instance.LaunchTime.Format(TimeStamp),
		SecurityGroups:   secGroups,
		Tags:             tags,
		VpcID:            *instance.VpcId,
	}
	return &p
}
