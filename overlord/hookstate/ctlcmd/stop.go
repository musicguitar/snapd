// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package ctlcmd

import (
	"github.com/snapcore/snapd/client"
	"github.com/snapcore/snapd/i18n"
	"github.com/snapcore/snapd/overlord/servicestate"
)

type stopCommand struct {
	baseCommand
	Positional struct {
		ServiceNames []string `positional-arg-name:"<service>" required:"yes"`
	} `positional-args:"yes" required:"yes"`
	Disable bool `long:"disable"`
}

var (
	shortStopHelp = i18n.G("Stop services")
)

func init() {
	addCommand("stop", shortStopHelp, "", func() command { return &stopCommand{} })
}

func (c *stopCommand) Execute(args []string) error {
	inst := servicestate.Instruction{
		Action: "stop",
		Names:  c.Positional.ServiceNames,
		StopOptions: client.StopOptions{
			Disable: c.Disable,
		},
	}
	return runServiceCommand(c.context(), &inst, c.Positional.ServiceNames)
}
