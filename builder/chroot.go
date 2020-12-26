//
// Copyright © 2016-2018 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package builder

import (
	"fmt"
	"github.com/getsolus/libosdev/commands"
	log "github.com/sirupsen/logrus"
	"os"
)

// Chroot will attempt to spawn a chroot in the overlayfs system
func (p *Package) Chroot(notif PidNotifier, pman *EopkgManager, overlay *Overlay) error {
	log.WithFields(log.Fields{
		"profile": overlay.Back.Name,
		"version": p.Version,
		"package": p.Name,
		"type":    p.Type,
		"release": p.Release,
	}).Debug("Beginning chroot")

	var env []string
	if p.Type == PackageTypeXML {
		env = SaneEnvironment("root", "/root")
	} else {
		env = SaneEnvironment(BuildUser, BuildUserHome)
	}
	ChrootEnvironment = env

	if err := p.ActivateRoot(overlay); err != nil {
		return err
	}

	// Now kill networking
	if p.Type == PackageTypeYpkg {
		if !p.CanNetwork {
			if err := DropNetworking(); err != nil {
				return err
			}

			// Ensure the overlay can network on localhost only
			if err := overlay.ConfigureNetworking(); err != nil {
				return err
			}
		} else {
			log.Warning("Package has explicitly requested networking, sandboxing disabled")
		}
	}

	log.Debug("Spawning login shell")
	// Allow bash to work
	commands.SetStdin(os.Stdin)

	// Legacy package format requires root, stay as root.
	user := BuildUser
	if p.Type == PackageTypeXML {
		user = "root"
	}

	loginCommand := fmt.Sprintf("/bin/su - %s -s %s", user, BuildUserShell)
	err := ChrootExecStdin(notif, overlay.MountPoint, loginCommand)
	commands.SetStdin(nil)
	notif.SetActivePID(0)
	return err
}
