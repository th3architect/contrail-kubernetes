/*
Copyright 2015 Juniper Networks, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package network

import (
	"io"
	"net"
	"time"

	"github.com/scalingdata/gcfg"
)

type Config struct {
	KubeUrl        string        `gcfg:"master"`
	ResyncPeriod   time.Duration `gcfg:"resync-interval"`
	ClusterIpRange string        `gcfg:"service-cluster-ip-range"`
}

type configWrapper struct {
	Default Config
}

func ReadConfiguration(reader io.Reader, config *Config) error {
	wrapper := configWrapper{Default: *config}
	wrapper.Default.ResyncPeriod = 0

	err := gcfg.ReadInto(&wrapper, reader)
	if err != nil {
		return err
	}

	if wrapper.Default.ResyncPeriod != 0 {
		wrapper.Default.ResyncPeriod = wrapper.Default.ResyncPeriod * time.Second
	} else {
		wrapper.Default.ResyncPeriod = config.ResyncPeriod
	}

	if clusterIp := wrapper.Default.ClusterIpRange; clusterIp != "" {
		if _, _, err := net.ParseCIDR(clusterIp); err != nil {
			return err
		}
	}

	*config = wrapper.Default
	return nil
}
