/*
Copyright (c) 2015 VMware, Inc. All Rights Reserved.

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

package portgroup

import (
	"flag"
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

type add struct {
	*flags.DatacenterFlag

	types.DVPortgroupConfigSpec

	path string
}

func init() {
	cli.Register("dvs.portgroup.add", &add{})
}

func (cmd *add) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.DatacenterFlag, ctx = flags.NewDatacenterFlag(ctx)
	cmd.DatacenterFlag.Register(ctx, f)

	f.StringVar(&cmd.path, "dvs", "", "DVS path")

	ptypes := []string{
		string(types.DistributedVirtualPortgroupPortgroupTypeEarlyBinding),
		string(types.DistributedVirtualPortgroupPortgroupTypeLateBinding),
		string(types.DistributedVirtualPortgroupPortgroupTypeEphemeral),
	}

	f.StringVar(&cmd.DVPortgroupConfigSpec.Type, "type", ptypes[0],
		fmt.Sprintf("Portgroup type (%s)", strings.Join(ptypes, "|")))

	cmd.DVPortgroupConfigSpec.NumPorts = 128 // default
	f.Var(flags.NewInt32(&cmd.DVPortgroupConfigSpec.NumPorts), "nports", "Number of ports")
}

func (cmd *add) Process(ctx context.Context) error {
	if err := cmd.DatacenterFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *add) Usage() string {
	return "NAME"
}

func (cmd *add) Run(ctx context.Context, f *flag.FlagSet) error {
	if f.NArg() == 0 {
		return flag.ErrHelp
	}

	name := f.Arg(0)

	finder, err := cmd.Finder()
	if err != nil {
		return err
	}

	net, err := finder.Network(ctx, cmd.path)
	if err != nil {
		return err
	}

	dvs, ok := net.(*object.DistributedVirtualSwitch)
	if !ok {
		return fmt.Errorf("%s (%T) is not of type %T", cmd.path, net, dvs)
	}

	cmd.DVPortgroupConfigSpec.Name = name

	task, err := dvs.AddPortgroup(ctx, []types.DVPortgroupConfigSpec{cmd.DVPortgroupConfigSpec})
	if err != nil {
		return err
	}

	logger := cmd.ProgressLogger(fmt.Sprintf("adding %s portgroup to dvs %s... ", name, dvs.InventoryPath))
	defer logger.Wait()

	_, err = task.WaitForResult(ctx, logger)
	return err
}
