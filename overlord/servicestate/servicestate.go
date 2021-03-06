// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015-2016 Canonical Ltd
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

package servicestate

import (
	"fmt"

	"github.com/snapcore/snapd/client"
	"github.com/snapcore/snapd/overlord/cmdstate"
	"github.com/snapcore/snapd/overlord/snapstate"
	"github.com/snapcore/snapd/overlord/state"
	"github.com/snapcore/snapd/snap"
)

type Instruction struct {
	Action string   `json:"action"`
	Names  []string `json:"names"`
	client.StartOptions
	client.StopOptions
	client.RestartOptions
}

type ServiceActionConflictError struct{ error }

func Change(st *state.State, appInfos []*snap.AppInfo, inst *Instruction) (*state.Change, error) {
	// the argv to call systemctl will need at most one entry per appInfo,
	// plus one for "systemctl", one for the action, and sometimes one for
	// an option. That's a maximum of 3+len(appInfos).
	argv := make([]string, 2, 3+len(appInfos))
	argv[0] = "systemctl"

	argv[1] = inst.Action
	switch inst.Action {
	case "start":
		if inst.Enable {
			argv[1] = "enable"
			argv = append(argv, "--now")
		}
	case "stop":
		if inst.Disable {
			argv[1] = "disable"
			argv = append(argv, "--now")
		}
	case "restart":
		if inst.Reload {
			argv[1] = "reload-or-restart"
		}
	default:
		return nil, fmt.Errorf("unknown action %q", inst.Action)
	}

	snapNames := make([]string, 0, len(appInfos))
	lastName := ""
	names := make([]string, len(appInfos))
	for i, svc := range appInfos {
		argv = append(argv, svc.ServiceName())
		snapName := svc.Snap.Name()
		names[i] = snapName + "." + svc.Name
		if snapName != lastName {
			snapNames = append(snapNames, snapName)
			lastName = snapName
		}
	}

	desc := fmt.Sprintf("%s of %v", inst.Action, names)

	st.Lock()
	defer st.Unlock()
	if err := snapstate.CheckChangeConflictMany(st, snapNames, nil); err != nil {
		return nil, &ServiceActionConflictError{err}
	}

	ts := cmdstate.Exec(st, desc, argv)
	chg := st.NewChange("service-control", desc)
	chg.AddAll(ts)
	return chg, nil
}
