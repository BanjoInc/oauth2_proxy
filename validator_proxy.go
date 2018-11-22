package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"unsafe"
)

type PolicyMap struct {
	policiesFile string
	m            unsafe.Pointer
}

type Policies struct {
	Policies []Policy `json:"policies"`
}

type Policy struct {
	Hosts      []string `json:"hosts"`
	Operations []string `json:"operations"`
	Users      []string `json:"users"`
}

func NewPolicyMap(policyFile string, done <-chan bool, onUpdate func()) *PolicyMap {
	pm := &PolicyMap{policiesFile: policyFile}
	m := make(map[string]bool)
	atomic.StorePointer(&pm.m, unsafe.Pointer(&m))
	if policyFile != "" {
		log.Printf("using policies file %s", policyFile)
		WatchForUpdates(policyFile, done, func() {
			pm.LoadPoliciesFile()
			onUpdate()
		})
		pm.LoadPoliciesFile()
	}
	return pm
}

func (pm *PolicyMap) LoadPoliciesFile() {
	fmt.Println("trying to load json config")
	r, err := os.Open(pm.policiesFile)
	if err != nil {
		log.Fatalf("failed opening polycies-file=%q, %s", pm.policiesFile, err)
	}
	defer r.Close()
	byteValue, _ := ioutil.ReadAll(r)
	var records Policies
	json.Unmarshal([]byte(byteValue), &records)
	fmt.Println("trying to unmarshall json config")
	if err != nil {
		log.Printf("error reading policies-file=%q, %s", pm.policiesFile, err)
		return
	}
	updated := make(map[string]bool)
	for _, policy := range records.Policies {
		fmt.Println("- policy")
		for _, host := range policy.Hosts {
			fmt.Println("- host")
			for _, operation := range policy.Operations {
				fmt.Println("- operation")
				for _, user := range policy.Users {
					fmt.Println("- user")
					address := strings.ToLower(host + ":" + operation + ":" + user)
					updated[address] = true
					fmt.Println("-->" + address)
				}
			}
		}
	}

	atomic.StorePointer(&pm.m, unsafe.Pointer(&updated))
}

func (pm *PolicyMap) IsValid(hostOperationEmail string) (result bool) {
	m := *(*map[string]bool)(atomic.LoadPointer(&pm.m))
	_, result = m[hostOperationEmail]
	return
}
