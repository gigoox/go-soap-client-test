package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const payload = `
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cli="http://mdwcorp.falabella.com/common/schema/clientserviceFIF" xmlns:cli1="http://mdwcorp.falabella.com/common/schema/clientservice" xmlns:req="http://mdwcorp.falabella.com/FIF/CORP/OSB/schema/Cliente/SolicitudesAdmision/Obtener/Req-v2018.01">
   <soapenv:Header>
      <cli:ClientServiceFIF/>
      <cli1:ClientService>
         <cli1:country>CO</cli1:country>
         <cli1:commerce>Banco</cli1:commerce>
         <cli1:channel>Web</cli1:channel>
      </cli1:ClientService>
   </soapenv:Header>
   <soapenv:Body>
      <req:clienteSolicitudesAdmisionObtenerExpReq>
         <documentoIdentidad>
            <tipoDocumento>1</tipoDocumento>
            <numeroDocumento>NUMERO-DOCUMENTO</numeroDocumento>
         </documentoIdentidad>
         <!--Optional:-->
         <situacion>
            <estado></estado>
         </situacion>
      </req:clienteSolicitudesAdmisionObtenerExpReq>
   </soapenv:Body>
</soapenv:Envelope>`

type Body  struct {
	XMLName  xml.Name
	R Response `xml:"clienteSolicitudesAdmisionObtenerExpResp"`
}

type Response  struct {
	XMLName xml.Name
	E EO  `xml:"estadoOperacion"`
	Solicitudes Solicitudes `xml:"solicitudes"`
}

type Solicitudes struct {

	Solicitud []Solicitud `xml:"solicitud"`
}

type Solicitud struct {

	NumeroSolicitud string `xml:"numeroSolicitud"`
	FechaSolicitud string `xml:"fechaSolicitud"`
	Situacion Situacion `xml:"situacion"`
	Producto Producto `xml:"producto"`
}

type Situacion struct {

	Estado string `xml:"estado"`
}

type Producto struct {

	CodigoProducto string `xml:"codigoProducto"`
}


type EO struct {

	CO string `xml:"codigoOperacion"`
	GO string `xml:"glosaOperacion"`
}

type Historico struct {
	XMLName xml.Name
	B Body    `xml:"Body"`
}

type Request struct {
	NumeroDocumento string
}

func (r *Request)getPayload() []byte{

	return []byte(strings.TrimSpace(strings.Replace(payload, "NUMERO-DOCUMENTO", r.NumeroDocumento,1)))
}


func main() {
	request := &Request{NumeroDocumento:"1618795"}
	url := "http://middlewaretest.falabella.cl/bco/co/12c/bus/Cliente/SolicitudesAdmision/Obtener/v1.0"
	soapAction := "urn:ClienteSolicitudesAdmisionObtenerOp" // The format is `urn:<soap_action>`
	httpMethod := "POST"
	req, err := http.NewRequest(httpMethod,url,bytes.NewReader(request.getPayload()))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}

	/*buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	fmt.Printf(newStr)*/

	result := new(Historico)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}

	fmt.Printf("%+v\n", result.B.R.Solicitudes)
	/*b, err := json.Marshal(result.B.R)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))*/

	request.NumeroDocumento = "1324123"

	req, err = http.NewRequest(httpMethod,url,bytes.NewReader(request.getPayload()))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err = client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}

	/*buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	newStr := buf.String()

	fmt.Printf(newStr)*/

	result = new(Historico)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}

	fmt.Printf("%+v\n", result.B.R.Solicitudes)


}


