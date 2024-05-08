package crds

type Crd struct {
	Resource string
	Group    string
	Version  string
}

func GetCRDList() []Crd {
	crdList := []Crd{
		{
			Resource: "apdoslogconfs",
			Group:    "appprotectdos.f5.com",
			Version:  "v1beta1",
		},
		{
			Resource: "apdospolicies",
			Group:    "appprotectdos.f5.com",
			Version:  "v1beta1",
		},
		{
			Resource: "dosprotectedresources",
			Group:    "appprotectdos.f5.com",
			Version:  "v1beta1",
		},
		{
			Resource: "aplogconfs",
			Group:    "appprotect.f5.com",
			Version:  "v1beta1",
		},
		{
			Resource: "appolicies",
			Group:    "appprotect.f5.com",
			Version:  "v1beta1",
		},
		{
			Resource: "apusersigs",
			Group:    "appprotect.f5.com",
			Version:  "v1beta1",
		},
		{
			Resource: "globalconfigurations",
			Group:    "k8s.nginx.org",
			Version:  "v1",
		},
		{
			Resource: "policies",
			Group:    "k8s.nginx.org",
			Version:  "v1",
		},
		{
			Resource: "transportservers",
			Group:    "k8s.nginx.org",
			Version:  "v1",
		},
		{
			Resource: "virtualserverroutes",
			Group:    "k8s.nginx.org",
			Version:  "v1",
		},
		{
			Resource: "virtualservers",
			Group:    "k8s.nginx.org",
			Version:  "v1",
		},
		{
			Resource: "cosarara",
			Group:    "k8s.nginx.org",
			Version:  "v1",
		},
	}
	return crdList
}

