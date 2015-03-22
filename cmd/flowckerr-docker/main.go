/*  Copyright (C) 2015 Leopoldo Lara Vazquez.

    This program is free software: you can redistribute it and/or  modify
    it under the terms of the GNU Affero General Public License, version 3,
	  as published by the Free Software Foundation.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// Command line interface to flowckerr-docker
//
// Runs a molecule as Docker containers. Library name of element should be
// a docker image tag.
// The environment should be configured correctly:
//  - DOCKER_HOST: where to find docker
//  - DOCKER_TLS_VERIFY=1: if you Docker is configured to use TLS
//  - DOCKER_CERT_PATH: Path to TLS certificate
//  - FLOWCKER_CONTAINERS_IP: IP in which the docker containers will be listening to,
//     important for example if you are using boot2docker in a mac
//  - FLOWCKER_CONTROL_IP: IP that is accesible in your host and that is in the same network as
//     FLOWCKER_CONTAINERS_IP, important if youa re using boot2docker in a mac
package main

import (
	"fmt"
	"os"

	fc "github.com/flowcker/flowcker/common"
	"github.com/flowcker/flowcker/runtime"
	"github.com/flowcker/flowcker/runtime/docker"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Must pass molecule")
		return
	}
	runtime.Run(fc.NewMoleculeFromFile(os.Args[1]), docker.LaunchAtom)
}
