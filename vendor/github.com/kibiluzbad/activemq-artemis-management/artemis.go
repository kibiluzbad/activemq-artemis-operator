package artemis

import (
	"strconv"
	"strings"

	"github.com/kibiluzbad/activemq-artemis-management/jolokia"
)

// IArtemis interface
type IArtemis interface {
	NewArtemis(_ip string, _jolokiaPort string, _name string) *Artemis
	Uptime() (*jolokia.ReadData, error)
	CreateAddress(addressName string, routingType string) (*jolokia.ExecData, error)
	CreateDivert(name string, routingName string, address string, forwardingAddress string, exclusive bool, filterString string, transformerClass string) (*jolokia.ExecData, error)
	CreateQueue(addressName string, queueName string, routingType string) (*jolokia.ExecData, error)
	DeleteAddress(addressName string) (*jolokia.ExecData, error)
	DeleteAddressForce(addressName string, force bool) (*jolokia.ExecData, error)
	DeleteDivert(name string) (*jolokia.ExecData, error)
	DeleteQueue(queueName string) (*jolokia.ExecData, error)
	ListBindingsForAddress(addressName string) (*jolokia.ExecData, error)
}

// Artemis struct
type Artemis struct {
	ip          string
	jolokiaPort string
	name        string
	jolokia     *jolokia.Jolokia
}

// NewArtemis new instance of Artemis struct.
func NewArtemis(_ip string, _jolokiaPort string, _name string) *Artemis {

	artemis := Artemis{
		ip:          _ip,
		jolokiaPort: _jolokiaPort,
		name:        _name,
		jolokia:     jolokia.NewJolokia(_ip, _jolokiaPort, "/console/jolokia"),
	}

	return &artemis
}

// Uptime show broker uptime.
func (artemis *Artemis) Uptime() (*jolokia.ReadData, error) {

	uptimeURL := "org.apache.activemq.artemis:broker=\"" + artemis.name + "\"/Uptime"
	data, err := artemis.jolokia.Read(uptimeURL)

	return data, err
}

// CreateAddress create an address.
func (artemis *Artemis) CreateAddress(addressName string, queueName string, routingType string) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	routingType = strings.ToUpper(routingType)
	parameters := `"` + addressName + `","` + routingType + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"createAddress(java.lang.String,java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

// CreateQueue create a queue.
func (artemis *Artemis) CreateQueue(addressName string, queueName string, routingType string) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	routingType = strings.ToUpper(routingType)
	parameters := `"` + addressName + `","` + queueName + `",` + `"` + routingType + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"createQueue(java.lang.String,java.lang.String,java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

// CreateDivert create a divert.
func (artemis *Artemis) CreateDivert(name string, routingName string, address string, forwardingAddress string, exclusive bool) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + name + `","` + routingName + `","` + address + `","` + forwardingAddress + `","` + strconv.FormatBool(exclusive) + `","",""`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"createDivert(java.lang.String,java.lang.String,java.lang.String,java.lang.String,boolean,java.lang.String,java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err

}

// DeleteQueue delte a queue.
func (artemis *Artemis) DeleteQueue(queueName string) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + queueName + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"destroyQueue(java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

// ListBindingsForAddress list bindings for address.
func (artemis *Artemis) ListBindingsForAddress(addressName string) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + addressName + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"listBindingsForAddress(java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

// DeleteAddress delete an address.
func (artemis *Artemis) DeleteAddress(addressName string) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + addressName + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"deleteAddress(java.lang.String)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}

// DeleteAddressForce force delete an address.
func (artemis *Artemis) DeleteAddressForce(addressName string, force bool) (*jolokia.ExecData, error) {

	url := "org.apache.activemq.artemis:broker=\\\"" + artemis.name + "\\\""
	parameters := `"` + addressName + `","` + strconv.FormatBool(force) + `"`
	jsonStr := `{ "type":"EXEC","mbean":"` + url + `","operation":"deleteAddress(java.lang.String, boolean)","arguments":[` + parameters + `]` + ` }`
	data, err := artemis.jolokia.Exec(url, jsonStr)

	return data, err
}
