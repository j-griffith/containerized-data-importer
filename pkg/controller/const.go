package controller

import "time"

const (
	CDI_LABEL_KEY         = "app"                                 //controller only
	CDI_LABEL_VALUE       = "containerized-data-importer"         //controller only
	CDI_LABEL_SELECTOR    = CDI_LABEL_KEY + "=" + CDI_LABEL_VALUE //controller and tests *can access from controller for all of them
	IMPORTER_PODNAME      = "importer"                            //Controller only
	IMPORTER_DATA_DIR     = "/data"                               //Controller only
	DEFAULT_PULL_POLICY   = "ifNotPresent"                        //Controller and cmd/controller suck it in from controller
	CLONING_LABEL_KEY     = "cloning"                             //controller only
	CLONING_LABEL_VALUE   = "host-assisted-cloning"               //controller only
	CLONING_TOPOLOGY_KEY  = "kubernetes.io/hostname"              //controller only
	CLONER_SOURCE_PODNAME = "clone-source-pod"                    //controller and tests
	CLONER_TARGET_PODNAME = "clone-target-pod"                    //controller and tests
	CLONER_IMAGE_PATH     = "/tmp/clone/image"                    //controller only
	CLONER_SOCKET_PATH    = "/tmp/clone/socket"                   //controller only
	KeyAccess             = "accessKeyId"                         //controller only
	KeySecret             = "secretKey"                           //controller only
	DEFAULT_RESYNC_PERIOD = 10 * time.Minute                      // controller only

	DataVolName         = "cdi-data-vol"
	ImagePathName       = "image-path"
	socketPathName      = "socket-path"
	controllerAgentName = "datavolume-controller"

	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by DataVolume"
	MessageResourceSynced = "DataVolume synced successfully"

	// pvc annotations
	AnnEndpoint  = "cdi.kubevirt.io/storage.import.endpoint"
	AnnSecret    = "cdi.kubevirt.io/storage.import.secretName"
	AnnImportPod = "cdi.kubevirt.io/storage.import.importPodName"
	// importer pod annotations
	AnnCreatedBy   = "cdi.kubevirt.io/storage.createdByController"
	AnnPodPhase    = "cdi.kubevirt.io/storage.import.pod.phase"
	LabelImportPvc = "cdi.kubevirt.io/storage.import.importPvcName"

	IMPORTER_ENDPOINT      = "IMPORTER_ENDPOINT"
	IMPORTER_ACCESS_KEY_ID = "IMPORTER_ACCESS_KEY_ID"
	IMPORTER_SECRET_KEY    = "IMPORTER_SECRET_KEY"
)
