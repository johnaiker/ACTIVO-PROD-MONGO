package request

// VueltoRq .
type VueltoRq struct {
	Locator string `json:"locator" validate:"required"`
	Usuario string `json:"usuario"`
	// UsuarioBanco string `json:"usuarioBanco"`
	//	Rif_comercio           string `json:"rif_comercio"`
	//	Telefono_comercio      string `json:"telefono_comercio"`
	//	Area_origen            string `json:"area_origen"`
	Tokena         string `json:"tokena" validate:"required"`
	Ip             string `json:"ip" validate:"required"`
	Monto          string `json:"monto" validate:"required"`
	Cedula_benef   string `json:"cedula_benef" validate:"required"`
	Telefono_benef string `json:"telefono_benef" validate:"required"`
	Bco_benef      string `json:"bco_benef" validate:"required"`
	Motivo         string `json:"motivo" validate:"required"`
}

// PayRP .
type CobroC2PRq struct {
	Locator       string `json:"locator" validate:"required"`
	Usuario       string `json:"usuario"`
	Tokena        string `json:"tokena" validate:"required"`
	Ip            string `json:"ip" validate:"required"`
	Monto         string `json:"monto" validate:"required"`
	Cedula_pago   string `json:"cedula_pago" validate:"required"`
	Telefono_pago string `json:"telefono_pago" validate:"required"`
	Bco_pago      string `json:"bco_pago" validate:"required"`
	Motivo        string `json:"motivo" validate:"required"`
	Otp           string `json:"otp" validate:"required"`
}
