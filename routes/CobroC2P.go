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
	// "github.com/leekchan/accounting"
)

// RequestCobro .
type RequestCobro struct {
	Locator           string `json:"locator"`
	Usuario           string `json:"usuario"`
	Rif_comercio      string `json:"rif_comercio"`
	Telefono_comercio string `json:"telefono_comercio"`
	Area_origen       string `json:"area_origen"`
	Ip                string `json:"ip" validate:"required"`
	Monto             string `json:"monto" validate:"required"`
	Cedula_pago       string `json:"cedula_pago" validate:"required"`
	Telefono_pago     string `json:"telefono_pago" validate:"required"`
	Bco_pago          string `json:"bco_pago" validate:"required"`
	Motivo            string `json:"motivo" validate:"required"`
	Otp               string `json:"otp" validate:"required"`
	// Tokena            string `json:"tokena" validate:"required"`
}

// Consulta .
func ValidarTransaccionCobroC2P(c echo.Context) (err error) {
	// var apikey, telefono_comercio, area_origen, rif string
	request := new(JSONSructRS.CobroC2PRq)
	if err = c.Bind(request); err != nil {
		helpers.ValidError(err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
		})
	}
	helpers.WriteFile([]byte(fmt.Sprintf("\n%s:::\n Consulta, Pagomovil: Lo que llega del frontend: %+v\n", time.Now().Format("2006-01-02 15:04:05"), request)), "proceso.txt")
	// helpers.WriteFile([]byte(fmt.Sprintln("%s:::%s\n", "Consulta, Pagomovil: Lo que llega del frontend: ", request, time.Now().Format("2006-01-02 15:04:05"))), "proceso.txt")

	//validaciones
	if err := c.Validate(request); err != nil {
		errArrayValidation := strings.Split(err.(validator.ValidationErrors).Error(), "\n")
		helpers.ValidError(err)
		return c.JSON(http.StatusBadRequest, map[string][]string{"errors": errArrayValidation})
	}
	fmt.Println("Consulta, Pagomovil: Lo que llega del frontend: ", request)
	//tomar el timepo para la bd
	t := time.Now()
	tiempo := t.Format("03:04:05 PM")
	fmt.Println(tiempo)

	//dandole forma a la fecha del time
	formatted := fmt.Sprintf("%d%02d%02d",
		t.Year(), t.Month(), t.Day())
	fechart := formatted
	Formatfecha := string(fechart[0]) + string(fechart[1]) + string(fechart[2]) + string(fechart[3]) + "-" + string(fechart[4]) + string(fechart[5]) + "-" + string(fechart[6]) + string(fechart[7])

	// montoDecimal, err := helpers.ToDecimal(request.Monto)
	// fmt.Println("Error al parsear monto?: ", err)

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
		// fmt.Println("Error: ", err)

		// return c.JSON(http.StatusInternalServerError, map[string]string{
		// 	"Error":   err.Error(),
		// 	"code":    "08",
		// 	"message": "error Fidempresa no encontrada",
		// })
		fmt.Println(err, "Error en la consulta")
		helpers.ValidError(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"Error":   err.Error(),
			"code":    "08",
			"message": "error Fidempresa no encontrada",
		})
	}

	montoconv := request.Monto
	montoactivo := strings.Replace(montoconv, ".", ",", 1)
	fmt.Println(montoactivo)

	//Tomando las variables de la request y remplazandolo para un cuerpo para enviar al banco
	Requests := new(RequestCobro)
	Requests.Usuario = request.Usuario
	Requests.Rif_comercio = query_response.Rif
	Requests.Telefono_comercio = query_response.Telefono_comercio
	Requests.Cedula_pago = request.Cedula_pago
	Requests.Telefono_pago = request.Telefono_pago
	Requests.Bco_pago = request.Bco_pago
	Requests.Monto = montoactivo
	Requests.Area_origen = query_response.Area_origen
	Requests.Motivo = request.Motivo
	Requests.Otp = request.Otp

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
	//cabeza de lo que enviamos al banco
	header := map[string]string{
		"content-type": "application/json",
		"apikey":       query_response.Apikey,
	}

	//url donde enviamos a BOD
	//	url := "https://bancoactivo-apimanager.siscotel.io:8243/cobro/c2p/1.0.0/c2p"
	url := "https://portalapis.bancoactivo.com:8243/cobro/c2p/1.0.0/c2p"

	bodyResp, _ := helpers.SendJSON(url, header, parsear)
	s := string(bodyResp)
	fmt.Println("para retornar: ", s)

	responses := new(JSONSructRS2.PayPM)
	jsonResp2 := json.Unmarshal(bodyResp, &responses)
	fmt.Println("jsonResp2:", jsonResp2)
	fmt.Println("Recibido ACTIVO:", responses)

	if responses.Code == "200" {
		collection := client.Database("activo").Collection("cobroc2p")

		insert_query := struct {
			Telefono_comercio string
			Locator           string
			Cedula_pago       string
			Telefono_pago     string
			Usuario           string
			Bco_pago          string
			Monto             primitive.Decimal128
			Mensaje           string
			Time              string
			Date              string
			Timestamp         time.Time
			Ip                string
			Ref               string
			Status            bool
		}{query_response.Telefono_comercio, request.Locator, request.Cedula_pago, request.Telefono_pago, request.Usuario, request.Bco_pago, montoDecimal_mongo, request.Motivo, tiempo, Formatfecha, t, request.Ip, responses.NroReferencia, true}

		if _, err = collection.InsertOne(context.TODO(), insert_query); err != nil {
			//inicio de log si hay algun error al insertar los datos en la tabla
			helpers.ValidError(err)
			fmt.Println("Pagomovil: Error al Insertar los datos a la tabla activo.cobroc2p: ", err.Error())
			fmt.Println("|=====================| Fin Pagomovil |========================|")
			helpers.WriteFile([]byte(fmt.Sprintln("|=====================| Fin Pagomovil |========================|", time.Now().Format("2006-01-02 15:04:05"))), "proceso.txt")
			//fin del log
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"code":    "001",
				"message": "Ocurrio un error inesperado, intente mas tarde",
			})
		}
	} else {
		collection := client.Database("activo").Collection("cobroc2p_fault")

		insert_query := struct {
			Telefono_comercio string
			Locator           string
			Cedula_pago       string
			Telefono_pago     string
			Usuario           string
			Bco_pago          string
			Monto             primitive.Decimal128
			Mensaje           string
			Time              string
			Date              string
			Timestamp         time.Time
			Ip                string
			Ref               string
			Status            bool
		}{query_response.Telefono_comercio, request.Locator, request.Cedula_pago, request.Telefono_pago, request.Usuario, request.Bco_pago, montoDecimal_mongo, request.Motivo, tiempo, Formatfecha, t, request.Ip, "N/A", false}

		if _, err = collection.InsertOne(context.TODO(), insert_query); err != nil {
			//inicio de log si hay algun error al insertar los datos en la tabla
			helpers.ValidError(err)
			fmt.Println("Pagomovil: Error al Insertar los datos a la tabla activo.cobroc2p: ", err.Error())
			fmt.Println("|=====================| Fin Pagomovil |========================|")
			helpers.WriteFile([]byte(fmt.Sprintln("|=====================| Fin Pagomovil |========================|", time.Now().Format("2006-01-02 15:04:05"))), "proceso.txt")
			//fin del log
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"Code":    "004",
				"Message": "Ocurrio un error inesperado, intente mas tarde",
			})
		}

	}
	return c.JSON(http.StatusOK, responses)
}
