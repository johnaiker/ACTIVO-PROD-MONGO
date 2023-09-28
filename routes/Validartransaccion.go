package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"vueltop2c/helpers"
	JSONSructRS "vueltop2c/request"
	JSONSructRS2 "vueltop2c/responses"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/validator.v9"
	//	"github.com/leekchan/accounting"
)

// RequestVuelto .
type RequestVuelto struct {
	Locator string `json:"locator"`
	Usuario string `json:"usuario"`
	// UsuarioBanco     string `json:"usuarioBanco"`
	Cedula_pagador   string `json:"cedula_pagador"`
	Telefono_pagador string `json:"telefono_pagador"`
	Area_origen      string `json:"area_origen"`
	Tokena           string `json:"tokena" validate:"required"`
	Ip               string `json:"ip" validate:"required"`
	Monto            string `json:"monto" validate:"required"`
	Cedula_benef     string `json:"cedula_benef" validate:"required"`
	Telefono_benef   string `json:"telefono_benef" validate:"required"`
	Bco_benef        string `json:"bco_benef" validate:"required"`
	Motivo           string `json:"motivo" validate:"required"`
	Tipo_cuenta      string `json:"tipo_cuenta" validate:"required"`
}

// Consulta .
func ValidarTransaccionVuelto(c echo.Context) (err error) {
	// var apikey, telefono_comercio, area_origen, rif string
	request := new(JSONSructRS.VueltoRq)
	if err = c.Bind(request); err != nil {
		helpers.ValidError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
		})
	}
	helpers.WriteFile([]byte(fmt.Sprintf("\n%s:::\n Consulta, Pagomovil: Lo que llega del frontend: %+v\n", time.Now().Format("2006-01-02 15:04:05"), request)), "proceso.txt")

	//validaciones
	if err := c.Validate(request); err != nil {
		errArrayValidation := strings.Split(err.(validator.ValidationErrors).Error(), "\n")
		helpers.ValidError(err)
		return c.JSON(http.StatusBadRequest, map[string][]string{"errors": errArrayValidation})
	}
	fmt.Printf("\nConsulta, Pagomovil: Lo que llega del frontend: %+v\n", request)
	//tomar el timepo para la bd
	t := time.Now()
	tiempo := t.Format("03:04:05 PM")
	fmt.Println(tiempo)

	//dandole forma a la fecha del time
	formatted := fmt.Sprintf("%d%02d%02d",
		t.Year(), t.Month(), t.Day())
	fechart := formatted
	Formatfecha := string(fechart[0]) + string(fechart[1]) + string(fechart[2]) + string(fechart[3]) + "-" + string(fechart[4]) + string(fechart[5]) + "-" + string(fechart[6]) + string(fechart[7])

	montoDecimal, err := helpers.ToDecimal(request.Monto)
	fmt.Println("Error al parsear monto?: ", err)

	montoDecimal_mongo, err := primitive.ParseDecimal128(request.Monto)

	if err != nil {
		fmt.Println(err.Error())
		helpers.ValidError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"code":    "001",
			"message": "Posible error al parseo a decimal del campo montoflaot",
			"Error":   err.Error(),
		})
	}

	//Llamado a la conexion de la BD para hacer el insert de los datos recibidos del frontend
	client := helpers.ConnectSessionDBMongo()
	collection := client.Database("pagalofacil").Collection("comercios")

	filter := bson.D{
		{Key: "tokena", Value: request.Tokena},
	}

	query_response := struct {
		Usuario           string
		Apikey            string
		Rif               string
		Telefono_comercio string
		Area_origen       string
	}{}

	if err := collection.FindOne(context.TODO(), filter).Decode(&query_response); err != nil {
		fmt.Println("Error: ", err)
		// return c.JSON(http.StatusBadRequest, map[string]string{
		// 	"message": "Token no asociado",
		// 	"code":    "02",
		// })
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"Error":   err.Error(),
			"code":    "08",
			"message": "error Fidempresa no encontrada",
		})
	}

	collection = client.Database("activo").Collection("vueltop2c")

	// if err := session.Query(`INSERT INTO
	// 	activo.vueltop2c
	// 	(telefono_pagador,
	// 	cedula_pagador,
	// 	locator,
	// 	cedula_benef,
	// 	telefono_benef,
	// 	usuario,
	// 	bco_benef,
	// 	monto,
	// 	area_origen,
	// 	mensaje,
	// 	time,
	// 	date,
	// 	timestamp,
	// 	ip,
	// 	ref,
	// 	status)
	// 	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	// 	rif,
	// 	telefono_comercio,
	// 	request.Locator,
	// 	request.Cedula_benef,
	// 	request.Telefono_benef,
	// 	request.Usuario,
	// 	request.Bco_benef,
	// 	montoDecimal,
	// 	area_origen,
	// 	request.Motivo,
	// 	tiempo,
	// 	Formatfecha,
	// 	t,
	// 	request.Ip,
	// 	"NA",
	// 	false).Exec(); err != nil {
	// 	//inicio de log si hay algun error al insertar los datos en la tabla
	// 	helpers.ValidError(err)
	// 	fmt.Println("Pagomovil: Error al Insertar los datos a la tabla maestra.pagomovil: ", err.Error())
	// 	fmt.Println("|=====================| Fin Pagomovil |========================|")
	// 	helpers.WriteFile([]byte(fmt.Sprintln("|=====================| Fin Pagomovil |========================|", time.Now().Format("2006-01-02 15:04:05"))), "proceso.txt")
	// 	//fin del log
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{
	// 		"code":    "001",
	// 		"message": "Ocurrio un error inesperado, intente mas tarde",
	// 	})
	// }
	// defer session.Close()

	//montoconv := request.Monto
	montoactivo := fmt.Sprintf("%.2f", montoDecimal)
	montoactivo = strings.Replace(montoactivo, ".", ",", 1)
	fmt.Println(montoactivo)
	// fmt.Println(str)

	//Tomando las variables de la request y remplazandolo para un cuerpo para enviar al banco
	Requests := new(RequestVuelto)
	Requests.Ip = request.Ip
	Requests.Locator = request.Locator
	Requests.Usuario = query_response.Usuario
	Requests.Cedula_pagador = query_response.Rif
	Requests.Telefono_pagador = query_response.Telefono_comercio
	Requests.Cedula_benef = request.Cedula_benef
	Requests.Telefono_benef = request.Telefono_benef
	Requests.Bco_benef = request.Bco_benef
	Requests.Monto = montoactivo
	Requests.Area_origen = query_response.Area_origen
	Requests.Motivo = request.Motivo
	Requests.Tipo_cuenta = "001"
	Requests.Tokena = request.Tokena

	fmt.Printf("\nConsulta , Pagomovil: Lo enviado al banco: %+v\n", Requests)
	helpers.WriteFile([]byte(fmt.Sprintf("%s:::\n Consulta, pagomovil: lo enviado al banco: %+v", time.Now().Format("2006-01-02 15:04:05"), Requests)), "proceso.txt")

	//parsear el json a []byte para el envio
	parsear, err := json.Marshal(Requests)
	if err != nil {
		helpers.ValidError(err)
		return c.JSON(http.StatusOK, map[string]string{
			"Error": err.Error(),
		})
	}
	//test
	// var header map[string]string

	header := make(map[string]string)

	header["content-type"] = "application/json"
	header["apikey"] = query_response.Apikey

	//url donde enviamos a ACTIVO
	//	url := "https://bancoactivo-apimanager.siscotel.io:8243/pago/p2c/1.0.0/p2c"
	url := "https://portalapis.bancoactivo.com:8243/pago/p2p/1.0.0/p2p"

	bodyResp, _ := helpers.SendJSON(url, header, parsear)
	s := string(bodyResp)
	fmt.Println("para retornar: ", s)

	responses := new(JSONSructRS2.PayPM)
	jsonResp2 := json.Unmarshal(bodyResp, &responses)
	fmt.Println("HOLA:", jsonResp2)
	fmt.Println("Recibido ACTIVO:", responses)

	if responses.Code == "200" {

		insertQuery := struct {
			Telefono_pagador string
			Cedula_pagador   string
			Locator          string
			Cedula_benef     string
			Telefono_benef   string
			Usuario          string
			Bco_benef        string
			Monto            primitive.Decimal128
			Mensaje          string
			Time             string
			Date             string
			Timestamp        time.Time
			Ip               string
			Ref              string
			Status           bool
		}{query_response.Rif, query_response.Telefono_comercio, request.Locator, request.Cedula_benef, request.Telefono_benef, request.Usuario, request.Bco_benef, montoDecimal_mongo, request.Motivo, tiempo, Formatfecha, t, request.Ip, responses.NroReferencia, true}

		if _, err = collection.InsertOne(context.TODO(), insertQuery); err != nil {
			//inicio de log si hay algun error al insertar los datos en la tabla
			helpers.ValidError(err)
			fmt.Println("Pagomovil: Error al Insertar los datos a la tabla maestra.pagomovil: ", err.Error())
			fmt.Println("|=====================| Fin Pagomovil |========================|")
			helpers.WriteFile([]byte(fmt.Sprintln("|=====================| Fin Pagomovil |========================|", time.Now().Format("2006-01-02 15:04:05"))), "proceso.txt")
			//fin del log

			if err = client.Disconnect(context.TODO()); err != nil {
				fmt.Println("Error: ", err)
			}

			return c.JSON(http.StatusInternalServerError, map[string]string{
				"Code":    "001",
				"Message": "Ocurrio un error inesperado, intente mas tarde",
			})
		}

		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println(err)

	} else {
		collection = client.Database("activo").Collection("vueltop2c_fault")

		insertQuery := struct {
			Telefono_pagador string
			Cedula_pagador   string
			Locator          string
			Cedula_benef     string
			Telefono_benef   string
			Usuario          string
			Bco_benef        string
			Monto            primitive.Decimal128
			Mensaje          string
			Time             string
			Date             string
			Timestamp        time.Time
			Ip               string
			Ref              string
			Status           bool
		}{query_response.Rif, query_response.Telefono_comercio,
			request.Locator,
			request.Cedula_benef,
			request.Telefono_benef,
			request.Usuario,
			request.Bco_benef,
			montoDecimal_mongo,
			request.Motivo,
			tiempo,
			Formatfecha,
			t,
			request.Ip,
			"NA", false}

		if _, err = collection.InsertOne(context.TODO(), insertQuery); err != nil {
			//inicio de log si hay algun error al insertar los datos en la tabla
			helpers.ValidError(err)
			fmt.Println("Pagomovil: Error al Insertar los datos a la tabla maestra.pagomovil: ", err.Error())
			fmt.Println("|=====================| Fin Pagomovil |========================|")
			helpers.WriteFile([]byte(fmt.Sprintln("|=====================| Fin Pagomovil |========================|", time.Now().Format("2006-01-02 15:04:05"))), "proceso.txt")

			if err = client.Disconnect(context.TODO()); err != nil {
				fmt.Println("Error: ", err)
			}
			//fin del log
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"Code":    "001",
				"Message": "Ocurrio un error inesperado, intente mas tarde",
			})
		}

		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println(err)
	}
	return c.JSON(http.StatusOK, responses)
}
