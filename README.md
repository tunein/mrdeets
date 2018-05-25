# MrDeets
Wrapper around the AWS cli for quick queries of AWS resources (currently only EC2) 

#Overview
When trying to access an AWS resource, most developers have to open the console, search for their instance, copy and paste the info they need then proceed.  Mr. Deets is a small binary to avoid the small but mondane task.  

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
