package database

import "github.com/idprm/go-yellowclinic/src/model"

var configs = []model.Config{
	{
		Name:  "AUTO_MESSAGE_SENDBIRD",
		Value: "Hi, Saya @v1 silahkan jelaskan keluhan kamu",
	},
	{
		Name:  "NOTIF_MESSAGE_DOCTOR",
		Value: "Hi *@v1*, User *@v2* menunggu konfirmasi untuk konsultasi online. Klik disini untuk memulai chat @v3",
	},
	{
		Name:  "NOTIF_MESSAGE_USER",
		Value: "Hi *@v1* pembayaran Anda sudah terkonfirmasi. Untuk Chat dengan *@v2* klik disini @v3 (Add to contact agar link bisa diklik)",
	},
	{
		Name:  "NOTIF_OTP_USER",
		Value: "Berikut adalah kode OTP kamu : *@v1* untuk mulai konsultasi dokter di sehatcepat.com",
	},
	{
		Name:  "PRICE",
		Value: "20000",
	},
	{
		Name:  "DISCOUNT",
		Value: "10000",
	},
	{
		Name:  "PAGE_FINISH",
		Value: "<p>Hi @v1 pastikan Anda sudah membayar Rp @v2 untuk berkonsultasi dengan @v3</p><p>(Abaikan apabila sudah membayar)</p><p>Kami akan mengirimkan Whatsapp notifikasi apabila pembayaran sudah terkonfirmasi</p>",
	},
	{
		Name:  "PAGE_UNFINISH",
		Value: "-",
	},
	{
		Name:  "PAGE_ERROR",
		Value: "-",
	},
}

var doctors = []model.Doctor{
	{
		Name:                 "dr. Ernita Rosyanti Dewi",
		Photo:                "dr-ernita.png",
		Number:               "STR 3121100220145699",
		Experience:           5,
		GraduatedFrom:        "Universitas Indonesia, 2007",
		ConsultationSchedule: "06.00 - 23.00 WIB",
		PlacePractice:        "Jakarta Timur, DKI Jakarta",
		Phone:                "6281776736076",
		UserId:               "dr-ernita",
	},
	{
		Name:                 "dr. Ayu A. Istiana",
		Photo:                "dr-ayu.png",
		Number:               "STR 3121100220145544",
		Experience:           7,
		GraduatedFrom:        "Universitas Yarsi, 2013",
		ConsultationSchedule: "06.00 - 23.00 WIB",
		PlacePractice:        "Bogor, Jawa Barat",
		Phone:                "6281212480644",
		UserId:               "dr-ayu",
	},
	{
		Name:                 "dr. Peter Fernando",
		Photo:                "dr-peter.png",
		Number:               "STR 6111100120221435",
		Experience:           3,
		GraduatedFrom:        "Universitas Tanjungpura, 2019",
		ConsultationSchedule: "06.00 - 23.00 WIB",
		PlacePractice:        "Ngabang, Kalimantan Timur",
		Phone:                "6281776736076",
		UserId:               "dr-peter",
	},
}
