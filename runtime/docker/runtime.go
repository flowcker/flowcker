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

package docker

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
	"github.com/fsouza/go-dockerclient"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("flowcker_runtime")

// Creates a docker client from the enviroment
func getClient() (client *docker.Client, err error) {
	endpoint := os.Getenv("DOCKER_HOST")

	if os.Getenv("DOCKER_TLS_VERIFY") == "1" {
		path := os.Getenv("DOCKER_CERT_PATH")
		ca := fmt.Sprintf("%s/ca.pem", path)
		cert := fmt.Sprintf("%s/cert.pem", path)
		key := fmt.Sprintf("%s/key.pem", path)
		client, err = docker.NewTLSClient(endpoint, cert, key, ca)
	} else {
		client, err = docker.NewClient(endpoint)
	}

	return client, err
}

func pullImage(client *docker.Client, image string) error {
	log.Info("Pulling image %s", image)

	var buff bytes.Buffer
	err := client.PullImage(docker.PullImageOptions{Repository: image, OutputStream: &buff, RawJSONStream: true}, docker.AuthConfiguration{})
	if err != nil {
		return err
	}

	log.Debug(buff.String())

	log.Info("Image %s pulled", image)

	return nil
}

func LaunchAtom(atom *fc.Atom) (string, error) {
	if os.Getenv("DOCKER_HOST") == "" {
		log.Fatal("DOCKER_HOST not set")
		return "", errors.New("DOCKER_HOST not set")
	}

	containersIP := os.Getenv("FLOWCKER_CONTAINERS_IP")
	if containersIP == "" {
		log.Warning("FLOWCKER_CONTAINERS_IP not set, setting to default 127.0.0.1")
		containersIP = "127.0.0.1"
	}

	element := strings.Split(atom.Element, " ")

	client, err := getClient()
	if err != nil {
		panic(err)
	}

	log.Notice("Launching Atom with docker image %s", element[0])

	_, err = client.InspectImage(element[0])
	if err == docker.ErrNoSuchImage {
		log.Info("Image not found, pulling")
		err = pullImage(client, element[0])
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else {
		log.Info("Image already present")
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{Config: &docker.Config{Image: element[0], Cmd: element[1:], ExposedPorts: map[docker.Port]struct{}{"3000": {}}}})
	if err != nil {
		spew.Dump(err)
		panic(err)
	}

	err = client.StartContainer(container.ID, &docker.HostConfig{PublishAllPorts: true})
	if err != nil {
		spew.Dump(err)
		panic(err)
	}

	container, err = client.InspectContainer(container.ID)
	if err != nil {
		spew.Dump(err)
		panic(err)
	}

	containerPort := container.NetworkSettings.Ports["3000"][0].HostPort

	return containersIP + ":" + containerPort, nil
}
