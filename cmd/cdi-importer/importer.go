package main

// importer.go implements a data fetching service capable of pulling objects from remote object
// stores and writing to a local directory. It utilizes the minio-go client sdk for s3 remotes,
// https for public remotes, and "file" for local files. The main use-case for this importer is
// to copy VM images to a "golden" namespace for consumption by kubevirt.
// This process expects several environmental variables:
//    IMPORTER_ENDPOINT       Endpoint url minus scheme, bucket/object and port, eg. s3.amazon.com.
//			      Access and secret keys are optional. If omitted no creds are passed
//			      to the object store client.
//    IMPORTER_ACCESS_KEY_ID  Optional. Access key is the user ID that uniquely identifies your
//			      account.
//    IMPORTER_SECRET_KEY     Optional. Secret key is the password to your account.

import (
	"flag"
	"os"

	"github.com/golang/glog"

	"kubevirt.io/containerized-data-importer/pkg/importer"
)

// Constants we use as env variables for the importer
const (
	importerEndpoint    = "IMPORTER_ENDPOINT"
	importerAccessKeyID = "IMPORTER_ACCESS_KEY_ID"
	importerSecretKey   = "IMPORTER_SECRET_KEY"
)

func init() {
	flag.Parse()
}

func main() {
	defer glog.Flush()

	glog.V(1).Infoln("Starting importer")
	ep, _ := importer.ParseEnvVar(importerEndpoint, false)
	acc, _ := importer.ParseEnvVar(importerAccessKeyID, false)
	sec, _ := importer.ParseEnvVar(importerSecretKey, false)

	glog.V(1).Infoln("begin import process")
	err := importer.CopyImage(importer.IMPORTER_WRITE_PATH, ep, acc, sec)
	if err != nil {
		glog.Errorf("%+v", err)
		os.Exit(1)
	}
	glog.V(1).Infoln("import complete")

	// temporary local import deprecation notice
	glog.Warningf("\nDEPRECATION NOTICE:\n   Support for local (file://) endpoints will be removed from CDI in the next release.\n   There is no replacement and no work-around.\n   All import endpoints must reference http(s) or s3 endpoints\n")
}
