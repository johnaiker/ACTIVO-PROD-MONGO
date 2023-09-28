package responses

//PayPM .

//type Error struct {
//    Code   string    `json:"code"`
//  Descripcion string    `json:"descripcion"`
//}

type PayPM struct {
	Code          string `json:"code" validate:"required"`
	NroReferencia string `json:"nroReferencia" `
	//	Error Error `json:"error"`
	Descripcion string `json:"descripcion"`
	Description string `json:"description"`
	Message     string `json:"message"`
}

//Airavailrs KiuAiravailrs
//type Error struct {
//	Code   string    `json:"code"`
//	Descripcion string    `json:"descripcion"`
//}
