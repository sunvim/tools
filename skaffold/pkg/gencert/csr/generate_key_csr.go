package csr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/cloudflare/cfssl/csr"
	"github.com/golang/glog"
)

const csrTpl = `
{
  "hosts": [
    "{{ .SrvName }}",
    "{{ .SrvName }}.{{ .Namespace }}",
    "{{ .SrvName }}.{{ .Namespace }}.svc"
  ],
  "CN": "kubernetes",
  "key": {
    "algo": "ecdsa",
    "size": 256
  },
  "names": [
    {
      "C": "CN",
      "ST": "Shanghai",
      "L": "Shanghai",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
`

const (
	csrFile = "admission/csr.json"
)

func GenerateKeyCsr(srvName, namespace string) (map[string]interface{}, error) {
	tpl, err := template.New("csr").Parse(csrTpl)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	arg := struct {
		SrvName   string
		Namespace string
	}{
		SrvName:   srvName,
		Namespace: namespace,
	}
	if err = tpl.Execute(buffer, arg); err != nil {
		glog.Error(err)
		return nil, err
	}
	// output csr json file
	if err = ioutil.WriteFile(csrFile, buffer.Bytes(), 0666); err != nil {
		glog.Error(err)
		return nil, err
	}
	// output server key and csr file
	body, err := generateCertSrvReq(buffer.Bytes())
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	outFiles(body)
	return body, nil
}

func generateCertSrvReq(csrBodies []byte) (map[string]interface{}, error) {
	req := csr.CertificateRequest{
		KeyRequest: csr.NewKeyRequest(),
	}
	err := json.Unmarshal(csrBodies, &req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	if req.CA != nil {
		glog.Error("ca section only permitted in initca")
		return nil, err
	}
	var key, csrPEM []byte
	g := &csr.Generator{Validator: Validator}
	csrPEM, key, err = g.ProcessRequest(&req)
	if err != nil {
		key = nil
		return nil, err
	}
	return packKeyCsrPem(key, csrPEM)
}

func Validator(req *csr.CertificateRequest) error {
	return nil
}

func packKeyCsrPem(key, csrBytes []byte) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	if key != nil {
		out["key"] = string(key)
	}

	if csrBytes != nil {
		out["csr"] = string(csrBytes)
	}
	return out, nil
}

type outputFile struct {
	Filename string
	Contents string
	IsBinary bool
	Perms    os.FileMode
}

const (
	serverKeyFile = "admission/server-key.pem"
	serverCsrFile = "admission/server.csr"
)

func outFiles(input map[string]interface{}) {
	var (
		keyContent string
		csrContent string
		outs       []outputFile
	)

	if contents, ok := input["key"]; ok {
		keyContent = contents.(string)
	}
	if keyContent != "" {
		outs = append(outs, outputFile{
			Filename: serverKeyFile,
			Contents: keyContent,
			Perms:    0600,
		})
	}

	if contents, ok := input["csr"]; ok {
		csrContent = contents.(string)
	}
	if csrContent != "" {
		outs = append(outs, outputFile{
			Filename: serverCsrFile,
			Contents: csrContent,
			Perms:    0644,
		})
	}

	for _, e := range outs {
		writeFile(e.Filename, e.Contents, e.Perms)
	}

}

func writeFile(fileName, contents string, perms os.FileMode) {
	err := ioutil.WriteFile(fileName, []byte(contents), perms)
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}
