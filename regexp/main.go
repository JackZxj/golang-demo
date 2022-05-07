package main

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("\\s+")

func match(expectedSQL, actualSQL string) error {
	expect := stripQuery(expectedSQL)
	actual := stripQuery(actualSQL)
	re, err := regexp.Compile(expect)
	if err != nil {
		return err
	}
	if !re.MatchString(actual) {
		return fmt.Errorf(`could not match actual sql: "%s" with expected regexp "%s"`, actual, re.String())
	}
	return nil
}

func stripQuery(q string) (s string) {
	return strings.TrimSpace(re.ReplaceAllString(q, " "))
}

func main() {
	fmt.Println(match("INSERT INTO `ocm_karmada_host`(.*)", "INSERT INTO `ocm_karmada_host` (`id`, `created_time`, `updated_time`, `is_deleted`, `deleted_time`, `creator`, `account_id`, `bill_type`, `bill_start_time`, `bill_end_time`, `bill_status`, `name`, `version`, `cluster_id`, `kubeconfig_id`, `karmada_host_kubeconfig_id`, `namespace`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)") == nil)

	b, e := regexp.MatchString("^[a-z0-9]([-a-z0-9]*[a-z0-9])?$", "aaaa")
	fmt.Println("#", b, e)
}
