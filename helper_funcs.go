package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/tunein/dshareiff-playground/mrDeets/metadata"
)

func prettyPrintJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}

func generateStrings(values []string) []*string {
	var response []*string
	for _, value := range values {
		if value != "" {
			response = append(response, aws.String(value))
		}
	}
	return response

}

func createFilters(opts *metadata.Ec2Options) []*ec2.Filter {
	filters := []*ec2.Filter{}
	msValuePtr := reflect.ValueOf(opts)
	msValue := msValuePtr.Elem()

	for i := 0; i < msValue.NumField(); i++ {
		field := msValue.Field(i)
		if strArray, ok := field.Interface().([]string); ok {
			if len(strArray) == 0 {
				continue
			}
			if parseTag(msValue.Type().Field(i).Tag) == "out" {
				continue
			}
			// Tags are special sauce.  If the caller wants to pass in specific flags they must be in the form of
			//  [{
			//  Name: "tag:TAG-NAME",
			//  Values: ["TAG-VALUE"]
			//  }]
			//
			if parseTag(msValue.Type().Field(i).Tag) == "tag" {
				k := reflect.Indirect(msValuePtr).FieldByName("Tag")
				values := []string(k.Interface().([]string))
				for _, value := range values {
					pair := strings.Split(value, "=")
					fmt.Println(pair)
					n := fmt.Sprintf("tag:%s", pair[0])
					f := ec2.Filter{
						Name:   aws.String(n),
						Values: []*string{aws.String(pair[1])},
					}
					filters = append(filters, &f)
				}
				continue
			}
			if parseTag(msValue.Type().Field(i).Tag) == "security-groups" {
				f := ec2.Filter{
					Name:   aws.String("instance.group-id"),
					Values: generateStrings(strArray),
				}
				filters = append(filters, &f)
				continue
			}

			f := ec2.Filter{
				Name:   aws.String(parseTag(msValue.Type().Field(i).Tag)),
				Values: generateStrings(strArray),
			}
			filters = append(filters, &f)
		}
	}
	fmt.Println(filters)
	return filters

}

func createOutFilter(opts *metadata.Ec2Options, instance *metadata.MinimalInstance) (results []string) {

	msValuePtr := reflect.ValueOf(opts)
	msValue := msValuePtr.Elem()
	for i := 0; i < msValue.NumField(); i++ {

		for _, s := range opts.Out {
			if (parseTag(msValue.Type().Field(i).Tag) == "tag-key") && (s == "tag-key") {
				tags := getTags(instance)
				for key := range tags {
					t := fmt.Sprintf("tag-key: %s", key)
					results = append(results, t)

				}
				continue
			}
			if (parseTag(msValue.Type().Field(i).Tag) == "tag-value") && (s == "tag-value") {
				tags := getTags(instance)
				for _, value := range tags {
					t := fmt.Sprintf("tag-value: %s", value)
					results = append(results, t)

				}
				continue
			}
			if (parseTag(msValue.Type().Field(i).Tag) == "tag-key") && (s == "tag" || s == "tags") {
				tags := getTags(instance)
				for key, value := range tags {
					t := fmt.Sprintf("tag: %s=%s", key, value)
					results = append(results, t)
				}
				continue
			}
			if parseTag(msValue.Type().Field(i).Tag) == "security-groups" {
				sg := getSecurityGroups(instance)
				for name, id := range sg {
					i := fmt.Sprintf("instance.group-id: %s", id)
					n := fmt.Sprintf("instance.group-name: %s", name)
					results = append(results, i)
					results = append(results, n)
				}
				continue
				results = append(results, getField(instance, msValue.Type().Field(i).Name))
			}
			if parseTag(msValue.Type().Field(i).Tag) == s {

				results = append(results, getField(instance, msValue.Type().Field(i).Name))
			}

		}
	}
	return
}
func getTags(v *metadata.MinimalInstance) map[string]string {
	tagFields := make(map[string]string)
	for _, tag := range v.Tags {
		r := reflect.ValueOf(tag)
		k := reflect.Indirect(r).FieldByName("Key")
		v := reflect.Indirect(r).FieldByName("Value")
		tagFields[string(k.String())] = string(v.String())
	}
	return tagFields
}
func getSecurityGroups(v *metadata.MinimalInstance) map[string]string {
	sgs := make(map[string]string)
	for _, sg := range v.SecurityGroups {
		r := reflect.ValueOf(sg)
		k := reflect.Indirect(r).FieldByName("GroupName")
		v := reflect.Indirect(r).FieldByName("GroupId")
		sgs[string(k.String())] = string(v.String())
	}
	return sgs
}

func getField(v *metadata.MinimalInstance, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

func removeQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func parseTag(tag reflect.StructTag) string {
	r := regexp.MustCompile(`name:(?P<Value>.*)\s`)
	resp := r.FindStringSubmatch(strings.TrimSpace(string(tag)))
	return removeQuotes(resp[1])
}
