# MrDeets
Wrapper around the AWS cli for quick queries of AWS resources (currently only EC2) 

# Overview
When trying to access an AWS resource, most developers have to open the console, search for their instance, copy and paste the info they need then proceed.  Mr. Deets is a small binary to avoid the small but mundane task.  

# Install


# Usage 
` ./mrDeets -option []{value} ` 
` ./mrDeets -option []{value} -out -option []{value} ` 

# Options
```shell
-availability-zone value
    	the availability zone of the instance
  -aws-region string
    	what region should should we search in
  -instance-id value
    	the instance-id of the instance
  -out value
    	what properties to display.  use --help to see properties
  -private-ip-address value
    	the private IP address
  -public-dns-name value
    	the public DNS name of the instance
  -public-ip-address value
    	the public IP address
  -security-groups value
    	the security groups of the instance
  -state value
    	the current state of the instance (running, pending, etc)
  -tag value
    	the tag key(s)/value(s) of the instance
  -tag-key value
    	the tag key(s) of the instance
  -tag-value value
    	the tag values(s) of the instance
  -vpc-id value
    	the VPC Id of the instance
```

# Examples 
*search by tag key/value*
```shell 
./mrDeets -tag Name=nonprod::streaming*  -aws-region us-west-2
```
*search by tag value*
```shell 
./mrDeets -tag-value nonprod*  -aws-region us-west-2
```
*get the IP address back from an instance id* 
```shell  
./mrDeets -instance-id i-05502a7f472395473 -out public-ip-address -aws-region us-west-2
54.201.42.177
```
*think I know part of the name and want the IP address* 
```shell 
./mrDeets -tag Name=*soundwave*development*  -aws-region us-west-2 -out public-ip-address
```
*ssh into the box using the instance id*
```shell 
ssh ec2-user@`./mrDeets -instance-id i-05502a7f472395473 -out public-ip-address -aws-region us-west-2`
Last login: Fri May 25 22:12:43 2018 from 38.140.202.59
...
[ec2-user@ip-10-82-28-212 ~]$
```
*combine filters* 
```shell 
./mrDeets -tag-value nonprod* -vpc-id vpc-c4717aa0  -aws-region us-west-2
``` 
*combine output filters*
``` shell 
./mrDeets -instance-id i-05502a7f472395473 -out "public-ip-address,vpc-id,tag" -aws-region us-west-2

54.201.42.177
vpc-c4717aa0
tag: aws:autoscaling:groupName=streaming_service-1-staging
tag: Name=nonprod::streaming_service::staging::1
tag: datadog=monitored
tag: Application=streaming_service
tag: Environment=staging
```
 
