/*
Copyright (c) 2016-2017 Parallels International GmbH.

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

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
)

var port = flag.Uint("port", 2000, "A port to serve requests on.")
var certFile = flag.String("certfile", "", "A PEM encoded certificate file.")
var certKey = flag.String("keyfile", "", "A PEM encoded private key file.")

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s --certfile <path> --keyfile <path> \n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	if *certFile == "" || *certKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	server := &http.Server{Addr: fmt.Sprintf(":%d", *port)}
	server.TLSConfig = &tls.Config{
		// Change default from SSLv3 to TLSv1.0 (because of POODLE vulnerability)
		MinVersion: tls.VersionTLS10,
	}

	http.Handle("/authenticate", newAuthWebhook())

	glog.Infof("Serving on :%d", *port)
	glog.Fatal(server.ListenAndServeTLS(*certFile, *certKey))

	os.Exit(0)
}
